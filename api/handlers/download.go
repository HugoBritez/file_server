package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"file-server-sofmar/config"
	"file-server-sofmar/middleware"
	"file-server-sofmar/models"

	"github.com/gorilla/mux"
)

// DownloadFile maneja la descarga de archivos
func DownloadFile(w http.ResponseWriter, r *http.Request) {
	// Obtener fileId de la URL
	vars := mux.Vars(r)
	fileID := vars["fileId"]
	if fileID == "" {
		sendErrorResponse(w, "ID de archivo requerido", http.StatusBadRequest)
		return
	}

	// Obtener client ID del contexto
	clientID := middleware.GetClientFromContext(r.Context())
	clientConfig, exists := config.GetClientConfig(clientID)
	if !exists {
		sendErrorResponse(w, "Cliente no configurado", http.StatusBadRequest)
		return
	}

	// Encontrar archivo por ID
	fileInfo, err := findFileByID(fileID, clientID, clientConfig.StoragePath)
	if err != nil {
		sendErrorResponse(w, "Archivo no encontrado: "+err.Error(), http.StatusNotFound)
		return
	}

	// Verificar que el archivo existe
	if _, err := os.Stat(fileInfo.Path); os.IsNotExist(err) {
		sendErrorResponse(w, "Archivo no existe en el filesystem", http.StatusNotFound)
		return
	}

	// Abrir archivo
	file, err := os.Open(fileInfo.Path)
	if err != nil {
		sendErrorResponse(w, "Error al abrir archivo: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Obtener información del archivo
	stat, err := file.Stat()
	if err != nil {
		sendErrorResponse(w, "Error al obtener información del archivo", http.StatusInternalServerError)
		return
	}

	// Headers para descarga
	w.Header().Set("Content-Type", fileInfo.MimeType)
	w.Header().Set("Content-Length", strconv.FormatInt(stat.Size(), 10))
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, fileInfo.OriginalName))
	w.Header().Set("Cache-Control", "public, max-age=31536000") // Cache por 1 año

	// Manejar range requests para streaming
	rangeHeader := r.Header.Get("Range")
	if rangeHeader != "" {
		handleRangeRequest(w, r, file, stat.Size(), fileInfo.MimeType)
		return
	}

	// Transferir archivo completo
	w.WriteHeader(http.StatusOK)
	_, err = io.Copy(w, file)
	if err != nil {
		// Log error but can't send response as headers are already sent
		fmt.Printf("Error streaming file: %v\n", err)
	}
}

// handleRangeRequest maneja requests con Range header para streaming
func handleRangeRequest(w http.ResponseWriter, r *http.Request, file *os.File, fileSize int64, mimeType string) {
	rangeHeader := r.Header.Get("Range")
	
	// Parse range header (formato: "bytes=start-end")
	if !strings.HasPrefix(rangeHeader, "bytes=") {
		w.WriteHeader(http.StatusRequestedRangeNotSatisfiable)
		return
	}

	rangeSpec := strings.TrimPrefix(rangeHeader, "bytes=")
	rangeParts := strings.Split(rangeSpec, "-")
	
	var start, end int64
	var err error

	// Parse start
	if rangeParts[0] != "" {
		start, err = strconv.ParseInt(rangeParts[0], 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusRequestedRangeNotSatisfiable)
			return
		}
	} else {
		start = 0
	}

	// Parse end
	if len(rangeParts) > 1 && rangeParts[1] != "" {
		end, err = strconv.ParseInt(rangeParts[1], 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusRequestedRangeNotSatisfiable)
			return
		}
	} else {
		end = fileSize - 1
	}

	// Validar range
	if start < 0 || start >= fileSize || end < start || end >= fileSize {
		w.WriteHeader(http.StatusRequestedRangeNotSatisfiable)
		return
	}

	// Mover a la posición inicial
	_, err = file.Seek(start, 0)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Headers para partial content
	contentLength := end - start + 1
	w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, fileSize))
	w.Header().Set("Content-Length", strconv.FormatInt(contentLength, 10))
	w.Header().Set("Content-Type", mimeType)
	w.WriteHeader(http.StatusPartialContent)

	// Transferir el rango solicitado
	_, err = io.CopyN(w, file, contentLength)
	if err != nil {
		fmt.Printf("Error streaming range: %v\n", err)
	}
}

// findFileByID busca un archivo por su ID en el filesystem
func findFileByID(fileID, clientID, storagePath string) (*models.FileMetadata, error) {
	// Construir ruta de búsqueda
	searchPath := filepath.Join("/app", storagePath)
	
	// Buscar archivos que empiecen con el fileID
	files, err := filepath.Glob(filepath.Join(searchPath, fileID+".*"))
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, fmt.Errorf("archivo no encontrado")
	}

	// Tomar el primer archivo encontrado
	filePath := files[0]
	
	// Obtener información del archivo
	stat, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}

	fileName := filepath.Base(filePath)
	extension := filepath.Ext(fileName)
	
	// Intentar extraer el nombre original (remover UUID prefix)
	originalName := fileName
	if len(fileName) > 36 && fileName[36] == '.' {
		// Formato: uuid + extension
		originalName = "file" + extension
	}

	// Crear metadata básica
	metadata := &models.FileMetadata{
		FileID:       fileID,
		OriginalName: originalName,
		FileName:     fileName,
		Client:       clientID,
		Size:         stat.Size(),
		Extension:    extension,
		Path:         filePath,
		URL:          fmt.Sprintf("/static/%s/%s", clientID, fileName),
	}

	return metadata, nil
}
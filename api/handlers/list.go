package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"file-server-sofmar/config"
	"file-server-sofmar/middleware"
	"file-server-sofmar/models"

	"github.com/gorilla/mux"
)

// ListFiles maneja el listado de archivos por cliente
func ListFiles(w http.ResponseWriter, r *http.Request) {
	// Obtener client ID de la URL
	vars := mux.Vars(r)
	clientID := vars["client"]
	if clientID == "" {
		// Fallback al contexto
		clientID = middleware.GetClientFromContext(r.Context())
	}

	if clientID == "" {
		sendErrorResponse(w, "Client ID requerido", http.StatusBadRequest)
		return
	}

	// Verificar que el cliente existe
	clientConfig, exists := config.GetClientConfig(clientID)
	if !exists {
		sendErrorResponse(w, "Cliente no configurado", http.StatusBadRequest)
		return
	}

	// Parámetros de query opcionales
	query := r.URL.Query()
	limit, _ := strconv.Atoi(query.Get("limit"))
	offset, _ := strconv.Atoi(query.Get("offset"))
	sortBy := query.Get("sort")
	order := query.Get("order")
	filter := query.Get("filter")
	folderFilter := query.Get("folder") // Filtrar por subcarpeta específica

	// Valores por defecto
	if limit <= 0 || limit > 1000 {
		limit = 100
	}
	if sortBy == "" {
		sortBy = "uploadedAt"
	}
	if order == "" {
		order = "desc"
	}

	// Buscar archivos en el directorio del cliente
	files, err := scanClientFiles(clientID, clientConfig.StoragePath)
	if err != nil {
		sendErrorResponse(w, "Error al listar archivos: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Filtrar archivos si se especifica
	if filter != "" {
		files = filterFiles(files, filter)
	}

	// Filtrar por subcarpeta si se especifica
	if folderFilter != "" {
		files = filterFilesByFolder(files, folderFilter)
	}

	// Ordenar archivos
	sortFiles(files, sortBy, order)

	// Aplicar paginación
	total := len(files)
	if offset >= total {
		files = []models.FileMetadata{}
	} else {
		end := offset + limit
		if end > total {
			end = total
		}
		files = files[offset:end]
	}

	// Preparar respuesta
	response := models.ListResponse{
		Success: true,
		Data:    files,
		Count:   total,
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Total-Count", strconv.Itoa(total))
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// scanClientFiles escanea el directorio de un cliente y retorna metadata de archivos
func scanClientFiles(clientID, storagePath string) ([]models.FileMetadata, error) {
	searchPath := filepath.Join("/app", storagePath)
	
	// Verificar que el directorio existe
	if _, err := os.Stat(searchPath); os.IsNotExist(err) {
		// Directorio no existe, retornar lista vacía
		return []models.FileMetadata{}, nil
	}

	var files []models.FileMetadata

	// Caminar por el directorio
	err := filepath.Walk(searchPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Saltear directorios
		if info.IsDir() {
			return nil
		}

		// Extraer fileID del nombre del archivo
		fileName := info.Name()
		fileID := extractFileID(fileName)
		
		// Detectar subcarpeta: obtener la ruta relativa desde el directorio del cliente
		relativePath, _ := filepath.Rel(searchPath, filepath.Dir(path))
		var folder string
		var fileURL string
		
		if relativePath == "." || relativePath == "" {
			// Archivo en el directorio raíz del cliente
			folder = ""
			fileURL = fmt.Sprintf("/static/%s/%s", clientID, fileName)
		} else {
			// Archivo en una subcarpeta
			folder = relativePath
			fileURL = fmt.Sprintf("/static/%s/%s/%s", clientID, folder, fileName)
		}
		
		// Crear metadata del archivo
		metadata := models.FileMetadata{
			FileID:       fileID,
			OriginalName: extractOriginalName(fileName),
			FileName:     fileName,
			Client:       clientID,
			Folder:       folder,
			Size:         info.Size(),
			MimeType:     getMimeTypeFromExtension(filepath.Ext(fileName)),
			Extension:    filepath.Ext(fileName),
			UploadedAt:   info.ModTime(),
			URL:          fileURL,
			Path:         path,
		}

		files = append(files, metadata)
		return nil
	})

	return files, err
}

// extractFileID extrae el UUID del nombre del archivo
func extractFileID(fileName string) string {
	// Formato esperado: uuid.extension
	parts := strings.Split(fileName, ".")
	if len(parts) > 0 && len(parts[0]) == 36 {
		return parts[0]
	}
	// Fallback: usar el nombre completo sin extensión
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}

// extractOriginalName intenta extraer el nombre original del archivo
func extractOriginalName(fileName string) string {
	// Por ahora retornamos el nombre actual
	// En una implementación más avanzada, esto vendría de una base de datos
	return fileName
}

// getMimeTypeFromExtension obtiene el MIME type de una extensión
func getMimeTypeFromExtension(ext string) string {
	ext = strings.ToLower(ext)
	mimeTypes := map[string]string{
		".pdf":  "application/pdf",
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".gif":  "image/gif",
		".txt":  "text/plain",
		".doc":  "application/msword",
		".docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		".xls":  "application/vnd.ms-excel",
		".xlsx": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		".zip":  "application/zip",
		".rar":  "application/x-rar-compressed",
	}

	if mimeType, exists := mimeTypes[ext]; exists {
		return mimeType
	}
	return "application/octet-stream"
}

// filterFiles filtra archivos por nombre o extensión
func filterFiles(files []models.FileMetadata, filter string) []models.FileMetadata {
	filter = strings.ToLower(filter)
	var filtered []models.FileMetadata

	for _, file := range files {
		if strings.Contains(strings.ToLower(file.OriginalName), filter) ||
			strings.Contains(strings.ToLower(file.FileName), filter) ||
			strings.Contains(strings.ToLower(file.Extension), filter) {
			filtered = append(filtered, file)
		}
	}

	return filtered
}

// filterFilesByFolder filtra archivos por subcarpeta específica
func filterFilesByFolder(files []models.FileMetadata, targetFolder string) []models.FileMetadata {
	var filtered []models.FileMetadata

	for _, file := range files {
		// Si targetFolder está vacío, mostrar solo archivos sin subcarpeta
		if targetFolder == "" && file.Folder == "" {
			filtered = append(filtered, file)
		} else if file.Folder == targetFolder {
			filtered = append(filtered, file)
		}
	}

	return filtered
}

// sortFiles ordena la lista de archivos
func sortFiles(files []models.FileMetadata, sortBy, order string) {
	sort.Slice(files, func(i, j int) bool {
		var result bool
		
		switch sortBy {
		case "name":
			result = files[i].OriginalName < files[j].OriginalName
		case "size":
			result = files[i].Size < files[j].Size
		case "extension":
			result = files[i].Extension < files[j].Extension
		case "uploadedAt":
			fallthrough
		default:
			result = files[i].UploadedAt.Before(files[j].UploadedAt)
		}

		if order == "desc" {
			result = !result
		}

		return result
	})
}
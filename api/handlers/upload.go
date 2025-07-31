package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"file-server-sofmar/config"
	"file-server-sofmar/middleware"
	"file-server-sofmar/models"

	"github.com/google/uuid"
)

// UploadFile maneja la subida de archivos
func UploadFile(w http.ResponseWriter, r *http.Request) {
	// Obtener client ID del contexto
	clientID := middleware.GetClientFromContext(r.Context())
	clientConfig, exists := config.GetClientConfig(clientID)
	if !exists {
		sendErrorResponse(w, "Cliente no configurado", http.StatusBadRequest)
		return
	}

	// Verificar método
	if r.Method != "POST" {
		sendErrorResponse(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Limitar tamaño del request
	r.Body = http.MaxBytesReader(w, r.Body, clientConfig.MaxFileSize)

	// Parse del form multipart
	err := r.ParseMultipartForm(clientConfig.MaxFileSize)
	if err != nil {
		sendErrorResponse(w, "Error al procesar archivo: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Obtener archivo del form
	file, header, err := r.FormFile("file")
	if err != nil {
		sendErrorResponse(w, "Archivo no encontrado en el formulario", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Validar tamaño
	if header.Size > clientConfig.MaxFileSize {
		sendErrorResponse(w, fmt.Sprintf("Archivo demasiado grande. Máximo: %d bytes", clientConfig.MaxFileSize), http.StatusBadRequest)
		return
	}

	// Validar tipo de archivo
	if !isAllowedFileType(header.Filename, clientConfig.AllowedTypes) {
		sendErrorResponse(w, "Tipo de archivo no permitido", http.StatusBadRequest)
		return
	}

	// Obtener subcarpeta (opcional)
	folder := r.FormValue("folder")
	if folder != "" {
		// Sanitizar nombre de carpeta (solo permitir letras, números, guiones y guiones bajos)
		folder = strings.ReplaceAll(folder, " ", "_")
		folder = strings.ReplaceAll(folder, "..", "") // Prevenir path traversal
		if len(folder) > 50 {
			folder = folder[:50] // Limitar longitud
		}
	}

	// Generar ID único para el archivo
	fileID := uuid.New().String()
	extension := filepath.Ext(header.Filename)
	fileName := fileID + extension

	// Crear ruta con subcarpeta si se especifica
	var uploadPath string
	if folder != "" {
		uploadPath = filepath.Join("/app", clientConfig.StoragePath, folder)
	} else {
		uploadPath = filepath.Join("/app", clientConfig.StoragePath)
	}

	// Crear directorio (incluyendo subcarpetas) si no existe
	if err := os.MkdirAll(uploadPath, 0755); err != nil {
		sendErrorResponse(w, "Error al crear directorio: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Ruta completa del archivo
	filePath := filepath.Join(uploadPath, fileName)

	// Crear archivo destino
	destFile, err := os.Create(filePath)
	if err != nil {
		sendErrorResponse(w, "Error al crear archivo: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer destFile.Close()

	// Copiar contenido con hash calculation
	hasher := sha256.New()
	writer := io.MultiWriter(destFile, hasher)
	_, err = io.Copy(writer, file)
	if err != nil {
		os.Remove(filePath) // Cleanup en caso de error
		sendErrorResponse(w, "Error al guardar archivo: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Calcular hash
	fileHash := hex.EncodeToString(hasher.Sum(nil))

	// Detectar MIME type
	mimeType := mime.TypeByExtension(extension)
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	// Crear URL con subcarpeta si existe
	var fileURL string
	if folder != "" {
		fileURL = fmt.Sprintf("/static/%s/%s/%s", clientID, folder, fileName)
	} else {
		fileURL = fmt.Sprintf("/static/%s/%s", clientID, fileName)
	}

	// Crear metadata del archivo
	metadata := models.FileMetadata{
		FileID:       fileID,
		OriginalName: header.Filename,
		FileName:     fileName,
		Client:       clientID,
		Folder:       folder,
		Size:         header.Size,
		MimeType:     mimeType,
		Extension:    extension,
		UploadedAt:   time.Now(),
		URL:          fileURL,
		Path:         filePath,
		Hash:         fileHash,
	}

	// TODO: Guardar metadata en base de datos si es necesario
	// Por ahora solo guardamos en filesystem

	// Respuesta exitosa
	response := models.UploadResponse{
		Success: true,
		Data:    metadata,
		Message: "Archivo subido exitosamente",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// isAllowedFileType verifica si el tipo de archivo está permitido
func isAllowedFileType(filename string, allowedTypes []string) bool {
	if len(allowedTypes) == 0 {
		return true // Sin restricciones
	}

	extension := strings.ToLower(filepath.Ext(filename))
	mimeType := mime.TypeByExtension(extension)

	for _, allowedType := range allowedTypes {
		allowedType = strings.ToLower(strings.TrimSpace(allowedType))
		
		// Permitir todos los tipos
		if allowedType == "*/*" {
			return true
		}

		// Verificar por MIME type
		if strings.Contains(allowedType, "/") {
			if allowedType == mimeType {
				return true
			}
			// Verificar wildcards como "image/*"
			if strings.HasSuffix(allowedType, "/*") {
				typePrefix := strings.TrimSuffix(allowedType, "/*")
				if strings.HasPrefix(mimeType, typePrefix+"/") {
					return true
				}
			}
		} else {
			// Verificar por extensión
			if allowedType == extension || allowedType == strings.TrimPrefix(extension, ".") {
				return true
			}
		}
	}

	return false
}

// sendErrorResponse envía una respuesta de error estandarizada
func sendErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	response := models.ErrorResponse{
		Success: false,
		Error:   message,
		Code:    statusCode,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}
package handlers

import (
	"encoding/json"
	"net/http"
	"os"

	"file-server-sofmar/config"
	"file-server-sofmar/middleware"
	"file-server-sofmar/models"

	"github.com/gorilla/mux"
)

// DeleteFile maneja la eliminación de archivos
func DeleteFile(w http.ResponseWriter, r *http.Request) {
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

	// Buscar archivo por ID
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

	// Verificar permisos (opcional - implementar según reglas de negocio)
	userID := middleware.GetUserFromContext(r.Context())
	if !canDeleteFile(userID, clientID, fileInfo) {
		sendErrorResponse(w, "No tienes permisos para eliminar este archivo", http.StatusForbidden)
		return
	}

	// Eliminar archivo del filesystem
	err = os.Remove(fileInfo.Path)
	if err != nil {
		sendErrorResponse(w, "Error al eliminar archivo: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO: Si usamos base de datos, eliminar metadata también
	// db.DeleteFileMetadata(fileID)

	// Respuesta exitosa
	response := map[string]interface{}{
		"success": true,
		"message": "Archivo eliminado exitosamente",
		"fileId":  fileID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// canDeleteFile verifica si un usuario puede eliminar un archivo
func canDeleteFile(userID, clientID string, fileInfo *models.FileMetadata) bool {
	// Implementar lógica de permisos según reglas de negocio
	// Por ejemplo:
	// - Admins pueden eliminar cualquier archivo
	// - Usuarios solo pueden eliminar sus propios archivos
	// - Algunos clientes tienen restricciones especiales
	
	// Por ahora permitimos a todos eliminar
	// En producción implementar con base de datos de usuarios
	return true
}

// BulkDelete maneja la eliminación múltiple de archivos
func BulkDelete(w http.ResponseWriter, r *http.Request) {
	// Estructura para recibir lista de IDs
	var request struct {
		FileIDs []string `json:"fileIds"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		sendErrorResponse(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	if len(request.FileIDs) == 0 {
		sendErrorResponse(w, "Lista de archivos vacía", http.StatusBadRequest)
		return
	}

	// Limitar número de archivos por operación
	if len(request.FileIDs) > 100 {
		sendErrorResponse(w, "Máximo 100 archivos por operación", http.StatusBadRequest)
		return
	}

	clientID := middleware.GetClientFromContext(r.Context())
	clientConfig, exists := config.GetClientConfig(clientID)
	if !exists {
		sendErrorResponse(w, "Cliente no configurado", http.StatusBadRequest)
		return
	}

	var successFiles []string
	var errorFiles []map[string]string

	// Procesar cada archivo
	for _, fileID := range request.FileIDs {
		fileInfo, err := findFileByID(fileID, clientID, clientConfig.StoragePath)
		if err != nil {
			errorFiles = append(errorFiles, map[string]string{
				"fileId": fileID,
				"error":  "Archivo no encontrado: " + err.Error(),
			})
			continue
		}

		// Verificar permisos
		userID := middleware.GetUserFromContext(r.Context())
		if !canDeleteFile(userID, clientID, fileInfo) {
			errorFiles = append(errorFiles, map[string]string{
				"fileId": fileID,
				"error":  "Sin permisos para eliminar",
			})
			continue
		}

		// Eliminar archivo
		err = os.Remove(fileInfo.Path)
		if err != nil {
			errorFiles = append(errorFiles, map[string]string{
				"fileId": fileID,
				"error":  "Error al eliminar: " + err.Error(),
			})
			continue
		}

		successFiles = append(successFiles, fileID)
	}

	// Respuesta con resultados
	response := map[string]interface{}{
		"success":      true,
		"deletedFiles": successFiles,
		"errors":       errorFiles,
		"total":        len(request.FileIDs),
		"deleted":      len(successFiles),
		"failed":       len(errorFiles),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
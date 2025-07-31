package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"

	"file-server-sofmar/config"
	"file-server-sofmar/middleware"
	"file-server-sofmar/models"

	"github.com/gorilla/mux"
)

// GetMetadata obtiene la metadata de un archivo específico
func GetMetadata(w http.ResponseWriter, r *http.Request) {
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

	// Verificar que el archivo existe y obtener información actualizada
	stat, err := os.Stat(fileInfo.Path)
	if os.IsNotExist(err) {
		sendErrorResponse(w, "Archivo no existe en el filesystem", http.StatusNotFound)
		return
	} else if err != nil {
		sendErrorResponse(w, "Error al acceder al archivo: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Actualizar información del archivo con datos del filesystem
	fileInfo.Size = stat.Size()
	fileInfo.UploadedAt = stat.ModTime()

	// Obtener información adicional si está disponible
	// TODO: Si usamos base de datos, obtener metadata extendida
	
	// Respuesta con metadata completa
	response := map[string]interface{}{
		"success": true,
		"data":    fileInfo,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// SearchFiles maneja la búsqueda de archivos con filtros avanzados
func SearchFiles(w http.ResponseWriter, r *http.Request) {
	// Obtener client ID de la URL
	vars := mux.Vars(r)
	clientID := vars["client"]
	if clientID == "" {
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

	// Decodificar request de búsqueda
	var searchReq models.SearchRequest
	if err := json.NewDecoder(r.Body).Decode(&searchReq); err != nil {
		sendErrorResponse(w, "JSON inválido: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validar request
	if searchReq.Query == "" {
		sendErrorResponse(w, "Query de búsqueda requerido", http.StatusBadRequest)
		return
	}

	// Obtener todos los archivos del cliente
	files, err := scanClientFiles(clientID, clientConfig.StoragePath)
	if err != nil {
		sendErrorResponse(w, "Error al buscar archivos: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Aplicar filtros de búsqueda
	filteredFiles := applySearchFilters(files, searchReq)

	// Preparar respuesta
	response := models.ListResponse{
		Success: true,
		Data:    filteredFiles,
		Count:   len(filteredFiles),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// applySearchFilters aplica los filtros de búsqueda a la lista de archivos
func applySearchFilters(files []models.FileMetadata, searchReq models.SearchRequest) []models.FileMetadata {
	var filtered []models.FileMetadata

	for _, file := range files {
		// Filtro por query (nombre del archivo)
		if !matchesQuery(file, searchReq.Query) {
			continue
		}

		// Filtro por tipos de archivo
		if len(searchReq.Types) > 0 && !matchesTypes(file, searchReq.Types) {
			continue
		}

		// Filtro por tamaño mínimo
		if searchReq.MinSize > 0 && file.Size < searchReq.MinSize {
			continue
		}

		// Filtro por tamaño máximo
		if searchReq.MaxSize > 0 && file.Size > searchReq.MaxSize {
			continue
		}

		// Filtro por fecha desde
		if searchReq.DateFrom != "" {
			if dateFrom, err := time.Parse("2006-01-02", searchReq.DateFrom); err == nil {
				if file.UploadedAt.Before(dateFrom) {
					continue
				}
			}
		}

		// Filtro por fecha hasta
		if searchReq.DateTo != "" {
			if dateTo, err := time.Parse("2006-01-02", searchReq.DateTo); err == nil {
				// Agregar 24 horas para incluir todo el día
				dateTo = dateTo.Add(24 * time.Hour)
				if file.UploadedAt.After(dateTo) {
					continue
				}
			}
		}

		// Si pasa todos los filtros, incluir en resultados
		filtered = append(filtered, file)
	}

	return filtered
}

// matchesQuery verifica si un archivo coincide con la query de búsqueda
func matchesQuery(file models.FileMetadata, query string) bool {
	query = strings.ToLower(query)
	return strings.Contains(strings.ToLower(file.OriginalName), query) ||
		strings.Contains(strings.ToLower(file.FileName), query) ||
		strings.Contains(strings.ToLower(file.Extension), query)
}

// matchesTypes verifica si un archivo coincide con los tipos especificados
func matchesTypes(file models.FileMetadata, types []string) bool {
	for _, allowedType := range types {
		if allowedType == "*" || allowedType == "*/*" {
			return true
		}

		// Verificar MIME type exacto
		if file.MimeType == allowedType {
			return true
		}

		// Verificar wildcards como "image/*"
		if strings.HasSuffix(allowedType, "/*") {
			typePrefix := strings.TrimSuffix(allowedType, "/*")
			if strings.HasPrefix(file.MimeType, typePrefix+"/") {
				return true
			}
		}

		// Verificar extensión
		if strings.HasPrefix(allowedType, ".") && file.Extension == allowedType {
			return true
		}
	}

	return false
}
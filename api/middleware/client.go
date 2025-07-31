package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"file-server-sofmar/config"
	"file-server-sofmar/models"
)

// ClientValidation middleware para validar clientes
func ClientValidation() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extraer client ID de la URL o header
			clientID := extractClientID(r)
			
			// Si no hay client ID, usar el default
			if clientID == "" {
				cfg := config.Load()
				clientID = cfg.DefaultClient
			}

			// Validar que el cliente existe
			if !config.IsValidClient(clientID) {
				errorResponse := models.ErrorResponse{
					Success: false,
					Error:   "Cliente no válido: " + clientID,
					Code:    http.StatusBadRequest,
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(errorResponse)
				return
			}

			// Añadir client ID al contexto
			ctx := context.WithValue(r.Context(), "clientID", clientID)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

// extractClientID extrae el ID del cliente de la URL o headers
func extractClientID(r *http.Request) string {
	// 1. Intentar obtener de la URL path (/api/files/list/{client})
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	
	// Para rutas como /api/files/list/{client}
	if len(pathParts) >= 4 && pathParts[0] == "api" && pathParts[1] == "files" {
		if pathParts[2] == "list" && len(pathParts) >= 4 {
			return pathParts[3]
		}
		if pathParts[2] == "search" && len(pathParts) >= 4 {
			return pathParts[3]
		}
	}

	// 2. Intentar obtener del header X-Client-Id
	if clientID := r.Header.Get("X-Client-Id"); clientID != "" {
		return clientID
	}

	// 3. Intentar obtener de query parameter
	if clientID := r.URL.Query().Get("client"); clientID != "" {
		return clientID
	}

	return ""
}

// GetClientFromContext obtiene el client ID del contexto
func GetClientFromContext(ctx context.Context) string {
	if clientID, ok := ctx.Value("clientID").(string); ok {
		return clientID
	}
	return ""
}
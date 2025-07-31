package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"file-server-sofmar/config"
	"file-server-sofmar/models"

	"github.com/golang-jwt/jwt/v5"
)

// JWTAuth middleware para autenticación JWT opcional
func JWTAuth() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Obtener cliente del contexto
			clientID := GetClientFromContext(r.Context())
			
			// Verificar si el cliente requiere autenticación
			clientConfig, exists := config.GetClientConfig(clientID)
			if !exists || !clientConfig.RequiresAuth {
				// Cliente no requiere auth, continuar
				next.ServeHTTP(w, r)
				return
			}

			// Cliente requiere autenticación, verificar token
			token := extractToken(r)
			if token == "" {
				unauthorizedResponse(w, "Token de autenticación requerido")
				return
			}

			// Validar JWT token
			userID, err := validateJWTToken(token)
			if err != nil {
				unauthorizedResponse(w, "Token inválido: "+err.Error())
				return
			}

			// Añadir user ID al contexto
			ctx := context.WithValue(r.Context(), "userID", userID)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

// extractToken extrae el token JWT del header Authorization
func extractToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	// Format: "Bearer <token>"
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}

	return parts[1]
}

// validateJWTToken valida un token JWT y retorna el user ID
func validateJWTToken(tokenString string) (string, error) {
	cfg := config.Load()
	
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verificar el método de signing
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(cfg.JWTSecret), nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Extraer user ID del token
		if userID, exists := claims["sub"]; exists {
			return userID.(string), nil
		}
		if userID, exists := claims["user_id"]; exists {
			return userID.(string), nil
		}
		return "unknown", nil
	}

	return "", jwt.ErrSignatureInvalid
}

// unauthorizedResponse envía una respuesta 401
func unauthorizedResponse(w http.ResponseWriter, message string) {
	errorResponse := models.ErrorResponse{
		Success: false,
		Error:   message,
		Code:    http.StatusUnauthorized,
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(errorResponse)
}

// GetUserFromContext obtiene el user ID del contexto
func GetUserFromContext(ctx context.Context) string {
	if userID, ok := ctx.Value("userID").(string); ok {
		return userID
	}
	return ""
}
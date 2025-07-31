package middleware

import (
	"log"
	"net/http"
	"time"
)

// Logging middleware para registrar todas las requests
func Logging() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Crear wrapper para capturar status code
			wrapper := &responseWrapper{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			// Procesar request
			next.ServeHTTP(wrapper, r)

			// Log de la request
			duration := time.Since(start)
			clientID := GetClientFromContext(r.Context())
			userID := GetUserFromContext(r.Context())

			log.Printf(
				"[%s] %s %s | Status: %d | Duration: %v | Client: %s | User: %s | Size: %d bytes",
				r.Method,
				r.URL.Path,
				r.RemoteAddr,
				wrapper.statusCode,
				duration,
				clientID,
				userID,
				wrapper.bytesWritten,
			)
		})
	}
}

// responseWrapper para capturar información de la respuesta
type responseWrapper struct {
	http.ResponseWriter
	statusCode   int
	bytesWritten int
}

func (rw *responseWrapper) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWrapper) Write(b []byte) (int, error) {
	rw.bytesWritten += len(b)
	return rw.ResponseWriter.Write(b)
}

// RateLimit middleware básico para prevenir abuso
func RateLimit() func(http.Handler) http.Handler {
	// Implementación básica - en producción usar redis o similar
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// TODO: Implementar rate limiting real con redis/memoria
			// Por ahora solo pasamos la request
			next.ServeHTTP(w, r)
		})
	}
}
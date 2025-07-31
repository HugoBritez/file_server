package main

import (
	"fmt"
	"log"
	"net/http"

	"file-server-sofmar/config"
	"file-server-sofmar/handlers"
	"file-server-sofmar/middleware"

	gorrillaHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	// Cargar configuración
	cfg := config.Load()

	// Crear router principal
	r := mux.NewRouter()

	// Middleware global
	r.Use(middleware.CORS())
	r.Use(middleware.Logging())

	// API Routes
	api := r.PathPrefix("/api").Subrouter()

	// Auth endpoint (sin autenticación)
	api.HandleFunc("/login", handlers.Login).Methods("POST")

	// File endpoints (con autenticación)
	files := api.PathPrefix("/files").Subrouter()
	files.Use(middleware.ClientValidation())
	files.Use(middleware.JWTAuth()) // 🔐 Autenticación JWT
	files.HandleFunc("/upload", handlers.UploadFile).Methods("POST")
	files.HandleFunc("/download/{fileId}", handlers.DownloadFile).Methods("GET")
	files.HandleFunc("/list/{client}", handlers.ListFiles).Methods("GET")
	files.HandleFunc("/{fileId}", handlers.DeleteFile).Methods("DELETE")
	files.HandleFunc("/metadata/{fileId}", handlers.GetMetadata).Methods("GET")
	files.HandleFunc("/search/{client}", handlers.SearchFiles).Methods("POST")

	// Health check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","service":"file-server","version":"1.0.0"}`))
	}).Methods("GET")

	// Configurar CORS
	corsHandler := gorrillaHandlers.CORS(
		gorrillaHandlers.AllowedOrigins(cfg.AllowedOrigins),
		gorrillaHandlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		gorrillaHandlers.AllowedHeaders([]string{"Content-Type", "Authorization", "X-Client-Id"}),
	)(r)

	port := cfg.Port
	if port == "" {
		port = "3000"
	}

	fmt.Printf("🚀 Servidor de archivos iniciado en puerto %s\n", port)
	fmt.Printf("📁 Directorio de uploads: %s\n", cfg.UploadDir)
	fmt.Printf("🌐 Health check: http://localhost:%s/health\n", port)
	fmt.Printf("📊 API endpoints: http://localhost:%s/api/files/\n", port)

	log.Fatal(http.ListenAndServe(":"+port, corsHandler))
}
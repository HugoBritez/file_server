package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Port           string
	UploadDir      string
	MaxFileSize    int64
	AllowedOrigins []string
	JWTSecret      string
	Environment    string
	DefaultClient  string
	AdminUser      string
	AdminPassword  string
}

func Load() *Config {
	// Tamaño máximo por defecto: 100MB
	maxSize := int64(100 * 1024 * 1024)
	if envSize := os.Getenv("MAX_FILE_SIZE"); envSize != "" {
		if size, err := parseSize(envSize); err == nil {
			maxSize = size
		}
	}

	// Orígenes permitidos
	origins := []string{"http://localhost:3000", "http://localhost:5173"}
	if envOrigins := os.Getenv("ALLOWED_ORIGINS"); envOrigins != "" {
		origins = strings.Split(envOrigins, ",")
		// Limpiar espacios en blanco
		for i, origin := range origins {
			origins[i] = strings.TrimSpace(origin)
		}
	}

	return &Config{
		Port:           getEnv("PORT", "3000"),
		UploadDir:      getEnv("UPLOAD_DIR", "/app/uploads"),
		MaxFileSize:    maxSize,
		AllowedOrigins: origins,
		JWTSecret:      getEnv("JWT_SECRET", "default_secret_change_in_production"),
		Environment:    getEnv("GO_ENV", "development"),
		DefaultClient:  getEnv("DEFAULT_CLIENT", "shared"),
		AdminUser:      getEnv("USER", "admin"),
		AdminPassword:  getEnv("PASSWORD", "admin123"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func parseSize(size string) (int64, error) {
	// Convertir "100MB" a bytes
	size = strings.ToUpper(strings.TrimSpace(size))
	
	if strings.HasSuffix(size, "MB") {
		size = strings.TrimSuffix(size, "MB")
		if val, err := strconv.ParseInt(size, 10, 64); err == nil {
			return val * 1024 * 1024, nil
		}
	} else if strings.HasSuffix(size, "GB") {
		size = strings.TrimSuffix(size, "GB")
		if val, err := strconv.ParseInt(size, 10, 64); err == nil {
			return val * 1024 * 1024 * 1024, nil
		}
	} else if strings.HasSuffix(size, "KB") {
		size = strings.TrimSuffix(size, "KB")
		if val, err := strconv.ParseInt(size, 10, 64); err == nil {
			return val * 1024, nil
		}
	}
	
	// Tratar como bytes si no tiene sufijo
	return strconv.ParseInt(size, 10, 64)
}
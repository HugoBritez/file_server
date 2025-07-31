package config

// ClientConfig define la configuración específica para cada cliente
type ClientConfig struct {
	MaxFileSize        int64    `json:"maxFileSize"`
	AllowedTypes       []string `json:"allowedTypes"`
	StoragePath        string   `json:"storagePath"`
	RequiresAuth       bool     `json:"requiresAuth"`
	CompressionEnabled bool     `json:"compressionEnabled"`
	Description        string   `json:"description"`
}

// ClientConfigs - Configuración basada en tu scripts/updateConfig.js existente
var ClientConfigs = map[string]ClientConfig{
	"acricolor": {
		MaxFileSize:        50 * 1024 * 1024, // 50MB
		AllowedTypes:       []string{"image/*", "application/pdf", "text/*"},
		StoragePath:        "uploads/acricolor",
		RequiresAuth:       true,
		CompressionEnabled: true,
		Description:        "Acricolor - Archivos de catálogos y documentos",
	},
	"lobeck": {
		MaxFileSize:        100 * 1024 * 1024, // 100MB
		AllowedTypes:       []string{"*/*"},
		StoragePath:        "uploads/lobeck",
		RequiresAuth:       true,
		CompressionEnabled: true,
		Description:        "Lobeck - Documentos técnicos y manuales",
	},
	"gaesa": {
		MaxFileSize:        200 * 1024 * 1024, // 200MB
		AllowedTypes:       []string{"*/*"},
		StoragePath:        "uploads/gaesa",
		RequiresAuth:       true,
		CompressionEnabled: true,
		Description:        "Gaesa - Archivos de ingeniería y proyectos",
	},

	"shared": {
		MaxFileSize:        10 * 1024 * 1024, // 10MB
		AllowedTypes:       []string{"image/*", "application/pdf", "text/*"},
		StoragePath:        "uploads/shared",
		RequiresAuth:       false,
		CompressionEnabled: false,
		Description:        "Archivos compartidos - Sin autenticación requerida",
	},
}

// GetClientConfig obtiene la configuración para un cliente específico
func GetClientConfig(clientID string) (ClientConfig, bool) {
	config, exists := ClientConfigs[clientID]
	return config, exists
}

// GetAllClients retorna la lista de todos los clientes configurados
func GetAllClients() []string {
	clients := make([]string, 0, len(ClientConfigs))
	for clientID := range ClientConfigs {
		clients = append(clients, clientID)
	}
	return clients
}

// IsValidClient verifica si un cliente está configurado
func IsValidClient(clientID string) bool {
	_, exists := ClientConfigs[clientID]
	return exists
}
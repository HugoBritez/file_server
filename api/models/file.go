package models

import (
	"time"
)

// FileMetadata representa la información de un archivo en el sistema
type FileMetadata struct {
	FileID      string    `json:"fileId"`
	OriginalName string   `json:"originalName"`
	FileName    string    `json:"fileName"`
	Client      string    `json:"client"`
	Folder      string    `json:"folder,omitempty"`
	Size        int64     `json:"size"`
	MimeType    string    `json:"mimeType"`
	Extension   string    `json:"extension"`
	UploadedAt  time.Time `json:"uploadedAt"`
	URL         string    `json:"url"`
	Path        string    `json:"path"`
	Hash        string    `json:"hash,omitempty"`
}

// UploadResponse representa la respuesta de una subida exitosa
type UploadResponse struct {
	Success bool         `json:"success"`
	Data    FileMetadata `json:"data,omitempty"`
	Error   string       `json:"error,omitempty"`
	Message string       `json:"message,omitempty"`
}

// ListResponse representa la respuesta de listado de archivos
type ListResponse struct {
	Success bool           `json:"success"`
	Data    []FileMetadata `json:"data,omitempty"`
	Count   int            `json:"count"`
	Error   string         `json:"error,omitempty"`
}

// ErrorResponse representa una respuesta de error
type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Code    int    `json:"code,omitempty"`
}

// SearchRequest representa una solicitud de búsqueda
type SearchRequest struct {
	Query    string   `json:"query"`
	Types    []string `json:"types,omitempty"`
	MinSize  int64    `json:"minSize,omitempty"`
	MaxSize  int64    `json:"maxSize,omitempty"`
	DateFrom string   `json:"dateFrom,omitempty"`
	DateTo   string   `json:"dateTo,omitempty"`
}

// HealthResponse representa la respuesta del health check
type HealthResponse struct {
	Status  string            `json:"status"`
	Service string            `json:"service"`
	Version string            `json:"version"`
	Uptime  string            `json:"uptime,omitempty"`
	Stats   map[string]interface{} `json:"stats,omitempty"`
}
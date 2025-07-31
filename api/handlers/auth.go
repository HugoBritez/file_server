package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"file-server-sofmar/config"

	"github.com/golang-jwt/jwt/v5"
)

// LoginRequest representa una solicitud de login
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse representa la respuesta del login
type LoginResponse struct {
	Success bool   `json:"success"`
	Token   string `json:"token,omitempty"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

// Login maneja la autenticación básica
func Login(w http.ResponseWriter, r *http.Request) {
	cfg := config.Load()

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"success":false,"error":"Datos inválidos"}`, http.StatusBadRequest)
		return
	}

	// Verificar credenciales
	if req.Username != cfg.AdminUser || req.Password != cfg.AdminPassword {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(LoginResponse{
			Success: false,
			Error:   "Usuario o contraseña incorrectos",
		})
		return
	}

	// Generar JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user":     req.Username,
		"client":   "shared",
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // 24 horas
		"iat":      time.Now().Unix(),
	})

	tokenString, err := token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(LoginResponse{
			Success: false,
			Error:   "Error generando token",
		})
		return
	}

	// Respuesta exitosa
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(LoginResponse{
		Success: true,
		Token:   tokenString,
		Message: "Login exitoso",
	})
}
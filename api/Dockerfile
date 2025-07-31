# Dockerfile para el API de archivos
FROM golang:1.21-alpine AS builder

# Instalar dependencias del sistema
RUN apk add --no-cache git ca-certificates

# Directorio de trabajo
WORKDIR /app

# Copiar archivos de configuración Go
COPY go.mod go.sum ./

# Descargar dependencias
RUN go mod download

# Copiar código fuente
COPY . .

# Compilar la aplicación
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Imagen final ligera
FROM alpine:latest

# Instalar ca-certificates para HTTPS
RUN apk --no-cache add ca-certificates

# Crear directorio de trabajo
WORKDIR /app

# Crear directorio para uploads
RUN mkdir -p /app/uploads

# Copiar binario compilado
COPY --from=builder /app/main .

# Exponer puerto
EXPOSE 3000

# Comando por defecto
CMD ["./main"]
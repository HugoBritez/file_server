#!/bin/bash

# Script de setup inicial para el servidor de archivos
echo "🚀 Configurando servidor de archivos Sofmar..."

# Verificar que Docker esté instalado
if ! command -v docker &> /dev/null; then
    echo "❌ Docker no está instalado. Por favor instala Docker primero."
    exit 1
fi

if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose no está instalado. Por favor instala Docker Compose primero."
    exit 1
fi

# Crear archivo .env si no existe
if [ ! -f .env ]; then
    echo "📝 Creando archivo .env..."
    cat > .env << EOF
# Configuración del Servidor de Archivos Sofmar
JWT_SECRET=$(openssl rand -base64 32)
MAX_FILE_SIZE=100MB
DEFAULT_CLIENT=shared
ALLOWED_ORIGINS=https://*.sofmar.com.py,https://*.gaesa.com.py,http://localhost:3000,http://localhost:5173
GO_ENV=production
PORT=3000
EOF
    echo "✅ Archivo .env creado con valores por defecto"
else
    echo "ℹ️ Archivo .env ya existe, saltando..."
fi

# Crear directorios de uploads si no existen
echo "📁 Creando directorios de uploads..."
mkdir -p uploads/{acricolor,lobeck,gaesa,shared}

# Configurar permisos
echo "🔒 Configurando permisos..."
chmod 755 uploads/
chmod 755 uploads/*

# Verificar la configuración de Docker
echo "🔍 Verificando configuración..."
if docker-compose config > /dev/null 2>&1; then
    echo "✅ Configuración de Docker Compose válida"
else
    echo "❌ Error en la configuración de Docker Compose"
    exit 1
fi

echo ""
echo "🎉 Setup completado!"
echo ""
echo "📖 Próximos pasos:"
echo "1. Ejecutar: docker-compose up -d"
echo "2. Verificar: curl http://localhost:4040/health"
echo "3. Ver logs: docker-compose logs -f"
echo ""
echo "📊 Endpoints disponibles:"
echo "   - Health check: http://localhost:4040/health"
echo "   - API base: http://localhost:4040/api/files/"
echo "   - Archivos estáticos: http://localhost:4040/static/{cliente}/{archivo}"
echo ""
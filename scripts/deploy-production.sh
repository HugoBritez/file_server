#!/bin/bash

# 🚀 Script de Deploy para Producción - Servidor de Archivos Sofmar
# Uso: ./scripts/deploy-production.sh

set -e

echo "🚀 Iniciando deploy en producción..."

# Verificar que estamos en el directorio correcto
if [ ! -f "docker-compose.yml" ]; then
    echo "❌ Error: Ejecutar desde el directorio raíz del proyecto"
    exit 1
fi

# Verificar que existe el archivo .env
if [ ! -f ".env" ]; then
    echo "❌ Error: Archivo .env no encontrado"
    echo "Crea el archivo .env con las variables necesarias"
    exit 1
fi

echo "📋 Verificando configuración..."

# Mostrar configuración actual (sin passwords)
echo "✅ Variables de entorno configuradas:"
echo "   - USER: $(grep '^USER=' .env | cut -d'=' -f2)"
echo "   - DEFAULT_CLIENT: $(grep '^DEFAULT_CLIENT=' .env | cut -d'=' -f2 || echo 'shared')"
echo "   - MAX_FILE_SIZE: $(grep '^MAX_FILE_SIZE=' .env | cut -d'=' -f2 || echo '100MB')"

# Detener servicios existentes
echo "🛑 Deteniendo servicios anteriores..."
docker-compose down 2>/dev/null || true

# Limpiar imágenes viejas (opcional)
echo "🧹 Limpiando imágenes Docker antiguas..."
docker system prune -f 2>/dev/null || true

# Construir y levantar servicios
echo "🔨 Construyendo servicios..."
docker-compose build --no-cache

echo "🚀 Iniciando servicios en producción..."
docker-compose up -d

# Esperar a que los servicios estén listos
echo "⏳ Esperando que los servicios estén listos..."
sleep 10

# Verificar estado de los servicios
echo "📊 Estado de los servicios:"
docker-compose ps

# Verificar conectividad
echo "🔍 Verificando conectividad..."

# Health check
if curl -f http://localhost:4040/health >/dev/null 2>&1; then
    echo "✅ Servidor respondiendo correctamente en puerto 4040"
else
    echo "❌ Error: Servidor no responde en puerto 4040"
    echo "📋 Logs del nginx:"
    docker-compose logs --tail=20 nginx
    echo "📋 Logs del API:"
    docker-compose logs --tail=20 api
    exit 1
fi

# Verificar si el puerto está abierto externamente
echo "🌐 Verificando acceso externo..."
PUBLIC_IP=$(curl -s ifconfig.me 2>/dev/null || echo "No disponible")
echo "   - IP Pública: $PUBLIC_IP"
echo "   - URL Local: http://localhost:4040/"
echo "   - URL Producción: https://node.sofmar.com.py:4040/"

echo ""
echo "🎉 ¡Deploy completado exitosamente!"
echo ""
echo "📋 URLs disponibles:"
echo "   🌍 Aplicación: https://node.sofmar.com.py:4040/"
echo "   🔍 Health Check: https://node.sofmar.com.py:4040/health"
echo "   📁 API: https://node.sofmar.com.py:4040/api/"
echo ""
echo "🔐 Credenciales de acceso:"
echo "   Usuario: $(grep '^USER=' .env | cut -d'=' -f2)"
echo "   Contraseña: [configurada en .env]"
echo ""
echo "📊 Para ver logs en tiempo real:"
echo "   docker-compose logs -f"
echo ""
echo "🛑 Para detener servicios:"
echo "   docker-compose down"
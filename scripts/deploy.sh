#!/bin/bash

# Script de deploy para el servidor de archivos
set -e

echo "🚀 Iniciando deploy del servidor de archivos..."

# Verificar que estamos en el directorio correcto
if [ ! -f "docker-compose.yml" ]; then
    echo "❌ Error: docker-compose.yml no encontrado"
    echo "   Ejecuta este script desde el directorio raíz del proyecto"
    exit 1
fi

# Crear backup de archivos si existen
if [ -d "uploads" ] && [ "$(ls -A uploads)" ]; then
    BACKUP_DATE=$(date +%Y%m%d_%H%M%S)
    echo "💾 Creando backup de archivos existentes..."
    tar -czf "backups/files_backup_$BACKUP_DATE.tar.gz" uploads/ 2>/dev/null || {
        mkdir -p backups
        tar -czf "backups/files_backup_$BACKUP_DATE.tar.gz" uploads/
    }
    echo "✅ Backup creado: backups/files_backup_$BACKUP_DATE.tar.gz"
fi

# Detener servicios existentes
echo "🛑 Deteniendo servicios existentes..."
docker-compose down

# Limpiar imágenes viejas (opcional)
if [ "$1" = "--clean" ]; then
    echo "🧹 Limpiando imágenes Docker viejas..."
    docker system prune -f
    docker-compose build --no-cache
else
    # Build de las imágenes
    echo "🏗️ Construyendo imágenes..."
    docker-compose build
fi

# Iniciar servicios
echo "▶️ Iniciando servicios..."
docker-compose up -d

# Esperar a que los servicios estén listos
echo "⏳ Esperando a que los servicios estén listos..."
sleep 10

# Verificar que los servicios estén funcionando
echo "🔍 Verificando servicios..."

# Health check nginx
if curl -f http://localhost:4040/health > /dev/null 2>&1; then
    echo "✅ nginx: OK"
else
    echo "❌ nginx: FALLO"
    echo "📋 Logs de nginx:"
    docker-compose logs nginx
    exit 1
fi

# Health check API
if curl -f http://localhost:4040/api/files/list/shared > /dev/null 2>&1; then
    echo "✅ API Go: OK"
else
    echo "❌ API Go: FALLO"
    echo "📋 Logs del API:"
    docker-compose logs api
    exit 1
fi

echo ""
echo "🎉 Deploy completado exitosamente!"
echo ""
echo "📊 Estado de servicios:"
docker-compose ps

echo ""
echo "🌐 URLs disponibles:"
echo "   - Health check: http://localhost:4040/health"
echo "   - API: http://localhost:4040/api/files/"
echo "   - Archivos estáticos: http://localhost:4040/static/"
echo ""
echo "📋 Comandos útiles:"
echo "   - Ver logs: docker-compose logs -f"
echo "   - Reiniciar: docker-compose restart"
echo "   - Detener: docker-compose down"
echo ""
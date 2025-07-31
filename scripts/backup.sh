#!/bin/bash

# Script de backup automático para el servidor de archivos
BACKUP_DIR="backups"
DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="files_backup_$DATE.tar.gz"

echo "💾 Iniciando backup de archivos..."

# Crear directorio de backups si no existe
mkdir -p $BACKUP_DIR

# Verificar que exista el directorio uploads
if [ ! -d "uploads" ]; then
    echo "❌ Directorio uploads no encontrado"
    exit 1
fi

# Crear backup
echo "📦 Creando archivo: $BACKUP_FILE"
tar -czf "$BACKUP_DIR/$BACKUP_FILE" uploads/

if [ $? -eq 0 ]; then
    echo "✅ Backup creado exitosamente"
    echo "📁 Archivo: $BACKUP_DIR/$BACKUP_FILE"
    echo "📊 Tamaño: $(du -h $BACKUP_DIR/$BACKUP_FILE | cut -f1)"
else
    echo "❌ Error al crear backup"
    exit 1
fi

# Limpiar backups viejos (mantener solo los últimos 7)
echo "🧹 Limpiando backups antiguos..."
find $BACKUP_DIR -name "files_backup_*.tar.gz" -mtime +7 -delete
REMAINING=$(find $BACKUP_DIR -name "files_backup_*.tar.gz" | wc -l)
echo "📊 Backups restantes: $REMAINING"

echo "✅ Backup completado"
# 📁 **Servidor de Archivos Sofmar**

Microservicio de archivos estáticos con API REST, usando Docker + nginx + Go. Compatible con tu infraestructura multi-cliente existente.

## 🚀 **Inicio Rápido**

### **1. Setup inicial**
```bash
# Clonar y entrar al directorio
git clone <repo> # o usar directorio actual
cd file-server

# Configurar el proyecto
./scripts/setup.sh

# Iniciar servicios
docker-compose up -d
```

### **2. Verificar funcionamiento**
```bash
# Health check
curl http://localhost:4040/health

# Listar archivos (cliente shared)
curl http://localhost:4040/api/files/list/shared
```

### **3. Test de upload**
```bash
# Subir un archivo de prueba
curl -X POST \
  -H "X-Client-Id: shared" \
  -F "file=@/path/to/your/file.pdf" \
  http://localhost:4040/api/files/upload
```

---

## 🏗️ **Arquitectura**

```
Cliente → nginx:4040 → {
  /static/* → nginx (archivos estáticos)
  /api/*    → Go API:3000 (gestión)
}
```

### **Servicios:**
- **nginx**: Proxy reverso + archivos estáticos (Puerto 4040)
- **go-api**: API REST para gestión de archivos (Puerto 3000)

---

## 📡 **API REST**

Ver documentación completa en: **[API-REST-DOCUMENTATION.md](./API-REST-DOCUMENTATION.md)**

### **Endpoints principales:**
```http
POST   /api/files/upload          # Subir archivo
GET    /api/files/download/{id}   # Descargar archivo  
GET    /api/files/list/{client}   # Listar archivos
DELETE /api/files/{id}            # Eliminar archivo
GET    /static/{client}/{file}    # Acceso directo (nginx)
```

### **Clientes configurados:**
- `acricolor`, `lobeck`, `gaesa` (con auth JWT requerido)
- `shared` (sin auth)

---

## 🔧 **Configuración**

### **Variables de entorno (.env):**
```bash
JWT_SECRET=tu_jwt_secret_super_seguro
MAX_FILE_SIZE=100MB
ALLOWED_ORIGINS=https://*.sofmar.com.py,https://*.gaesa.com.py
DEFAULT_CLIENT=shared
```

### **Límites por cliente:**
- **acricolor**: 50MB, imágenes/PDF/texto (requiere JWT)
- **lobeck**: 100MB, todos los tipos (requiere JWT)
- **gaesa**: 200MB, todos los tipos (requiere JWT)
- **shared**: 10MB, imágenes/PDF/texto (sin auth)

---

## 🛠️ **Scripts Disponibles**

```bash
# Setup inicial
./scripts/setup.sh

# Deploy/actualizar
./scripts/deploy.sh

# Backup de archivos
./scripts/backup.sh

# Deploy con limpieza completa
./scripts/deploy.sh --clean
```

---

## 📊 **Comandos Docker**

```bash
# Iniciar servicios
docker-compose up -d

# Ver logs
docker-compose logs -f

# Reiniciar servicios
docker-compose restart

# Detener servicios
docker-compose down

# Ver estado
docker-compose ps

# Reconstruir imágenes
docker-compose build
```

---

## 🌐 **URLs en Producción**

Configurar en tu servidor (puerto 4040):
```
https://node.sofmar.com.py:4040/api/files/    # API
https://node.sofmar.com.py:4040/static/       # Archivos
https://node.sofmar.com.py:4040/health        # Health check
```

---

## 📁 **Estructura del Proyecto**

```
file-server/
├── docker-compose.yml           # Orquestación de servicios
├── .env                        # Variables de entorno
├── api/                        # API Go
│   ├── main.go                 # Servidor principal
│   ├── handlers/               # Controladores HTTP
│   ├── middleware/             # Middleware (CORS, auth)
│   ├── config/                 # Configuración multi-cliente
│   └── models/                 # Modelos de datos
├── nginx/
│   └── nginx.conf              # Configuración nginx optimizada
├── uploads/                    # Archivos por cliente
│   ├── acricolor/              # Requiere JWT
│   ├── lobeck/                 # Requiere JWT
│   ├── gaesa/                  # Requiere JWT
│   └── shared/                 # Sin autenticación
└── scripts/                    # Scripts de deploy y backup
```

---

## 🔒 **Seguridad**

- ✅ Validación de tipos de archivo por cliente
- ✅ Límites de tamaño configurables
- ✅ Rate limiting en nginx
- ✅ Headers de seguridad
- ✅ CORS configurado
- ✅ Autenticación JWT (opcional por cliente)

---

## 📈 **Características**

- ✅ **Ultra-rápido**: nginx sirve archivos estáticos directamente
- ✅ **Escalable**: Docker + nginx + Go
- ✅ **Multi-cliente**: Configuración específica por cliente
- ✅ **API REST**: Fácil integración con cualquier frontend
- ✅ **Streaming**: Soporte para archivos grandes (range requests)
- ✅ **Backup**: Scripts automáticos de respaldo
- ✅ **Monitoreo**: Logs estructurados y health checks

---

## 🔗 **Integración con tus Clientes**

### **JavaScript/Fetch:**
```javascript
// Subir archivo
const formData = new FormData();
formData.append('file', file);

const response = await fetch('http://localhost:4040/api/files/upload', {
  method: 'POST',
  headers: { 'X-Client-Id': 'acricolor' },
  body: formData
});

const result = await response.json();
console.log('URL del archivo:', result.data.url);
```

### **React/Vue/Angular:**
Ver ejemplos completos en [API-REST-DOCUMENTATION.md](./API-REST-DOCUMENTATION.md)

---

## 🚨 **Troubleshooting**

### **Servicios no inician:**
```bash
# Verificar logs
docker-compose logs

# Verificar configuración
docker-compose config

# Reiniciar limpio
docker-compose down && docker-compose up -d
```

### **Error de permisos:**
```bash
# Arreglar permisos de uploads
chmod 755 uploads/
chmod 755 uploads/*
```

### **Puerto ocupado:**
```bash
# Verificar que puerto 4040 esté libre
netstat -tlnp | grep 4040

# O cambiar puerto en docker-compose.yml
```

---

## 📞 **Soporte**

- **Health Check**: `http://localhost:4040/health`
- **Logs**: `docker-compose logs -f`
- **API Docs**: [API-REST-DOCUMENTATION.md](./API-REST-DOCUMENTATION.md)

---

**¡Listo para usar! 🎉**# file_server

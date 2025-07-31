# 🚀 **Guía de Deploy en Servidor - node.sofmar.com.py**

## 📋 **Checklist de Deploy**

### **1. 📁 Crear archivo .env en el servidor**

En el directorio raíz del proyecto, crear el archivo `.env`:

```bash
# 🔐 Configuración de autenticación
USER=Sofmar
PASSWORD=s17052006

# 🔑 JWT Secret (PRODUCCIÓN)
JWT_SECRET=sofmar_file_server_jwt_secret_2024_prod_secure

# 📁 Configuración del servidor
MAX_FILE_SIZE=100MB
DEFAULT_CLIENT=shared

# 🌐 CORS - Orígenes permitidos para PRODUCCIÓN
ALLOWED_ORIGINS=https://node.sofmar.com.py,https://*.sofmar.com.py,https://*.gaesa.com.py
```

### **2. 🌐 Verificar DNS**

Verificar que el dominio apunte al servidor:
```bash
# Verificar resolución DNS
nslookup node.sofmar.com.py

# Debe devolver la IP de tu servidor
```

### **3. 🔥 Configurar Firewall**

Abrir los puertos necesarios:
```bash
# Ubuntu/Debian
sudo ufw allow 4040/tcp
sudo ufw allow 4043/tcp
sudo ufw reload

# CentOS/RHEL
sudo firewall-cmd --permanent --add-port=4040/tcp
sudo firewall-cmd --permanent --add-port=4043/tcp
sudo firewall-cmd --reload
```

### **4. 🚀 Ejecutar Deploy**

```bash
# Dar permisos al script
chmod +x scripts/deploy-production.sh

# Ejecutar deploy
./scripts/deploy-production.sh
```

### **5. ✅ Verificar Funcionamiento**

Después del deploy, verificar:

```bash
# 1. Estado de contenedores
docker-compose ps

# 2. Health check local
curl http://localhost:4040/health

# 3. Logs en tiempo real
docker-compose logs -f
```

## 🌍 **URLs de Acceso**

- **Aplicación Web:** https://node.sofmar.com.py:4040/
- **Health Check:** https://node.sofmar.com.py:4040/health
- **API:** https://node.sofmar.com.py:4040/api/

## 🔐 **Credenciales de Acceso**

- **Usuario:** `Sofmar`
- **Contraseña:** `s17052006`

## 🛠️ **Comandos Útiles**

```bash
# Ver logs
docker-compose logs -f

# Reiniciar servicios
docker-compose restart

# Detener servicios
docker-compose down

# Reconstruir y reiniciar
docker-compose down && docker-compose build --no-cache && docker-compose up -d

# Ver estado del sistema
docker system df
docker-compose ps
```

## ⚠️ **Troubleshooting**

### **ERR_CONNECTION_REFUSED**

Si obtienes este error:

1. **Verificar servicios:**
   ```bash
   docker-compose ps
   ```

2. **Verificar puertos:**
   ```bash
   netstat -tlnp | grep 4040
   ```

3. **Verificar logs:**
   ```bash
   docker-compose logs nginx
   docker-compose logs api
   ```

4. **Verificar firewall:**
   ```bash
   sudo ufw status
   ```

### **Servicios no inician:**

```bash
# Limpiar y reiniciar
docker-compose down
docker system prune -f
docker-compose build --no-cache
docker-compose up -d
```

### **Problemas de permisos:**

```bash
# Arreglar permisos de uploads
sudo chown -R 1000:1000 uploads/
chmod -R 755 uploads/
```

## 🔄 **Actualizaciones Futuras**

Para actualizar el código:

```bash
# 1. Hacer pull del código
git pull origin main

# 2. Reiniciar servicios
./scripts/deploy-production.sh
```

---

**¡Listo para producción! 🎉**
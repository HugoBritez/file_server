# 🔐 **Guía de Autenticación - Servidor de Archivos**

## **Cómo Funciona la Autenticación**

El sistema usa **autenticación por cliente** con tokens JWT. Cada cliente puede configurarse para requerir o no autenticación.

---

## **🏷️ Configuración por Cliente**

| Cliente | Requiere Auth | Tamaño Max | Descripción |
|---------|--------------|------------|-------------|
| `acricolor` | ✅ SÍ | 50MB | Archivos de catálogos y documentos |
| `lobeck` | ✅ SÍ | 100MB | Documentos técnicos y manuales |
| `gaesa` | ✅ SÍ | 200MB | Archivos de ingeniería y proyectos |
| `shared` | ❌ NO | 10MB | Archivos compartidos públicos |

---

## **🔑 Cómo Autenticarse**

### **1. Clientes SIN Autenticación (`shared`)**
```javascript
// No necesita token, solo especificar cliente
const response = await fetch('/api/files/upload', {
  method: 'POST',
  headers: {
    'X-Client-Id': 'shared'
  },
  body: formData
});
```

### **2. Clientes CON Autenticación (`acricolor`, `lobeck`, `gaesa`)**
```javascript
// Necesita token JWT en Authorization header
const response = await fetch('/api/files/upload', {
  method: 'POST',
  headers: {
    'Authorization': 'Bearer YOUR_JWT_TOKEN',
    'X-Client-Id': 'acricolor'
  },
  body: formData
});
```

---

## **🛠️ Formato del Token JWT**

El token JWT debe contener el **user ID** en uno de estos campos:

```json
{
  "sub": "user123",           // Campo estándar (preferred)
  "user_id": "user123",       // Campo alternativo
  "iat": 1640995200,          // Timestamp de emisión
  "exp": 1641081600           // Timestamp de expiración
}
```

### **Ejemplo de creación de token (Node.js):**
```javascript
const jwt = require('jsonwebtoken');

const token = jwt.sign(
  { 
    sub: 'user123',           // ID del usuario
    client: 'acricolor',      // Cliente (opcional)
    permissions: ['upload', 'download'] // Permisos (opcional)
  },
  process.env.JWT_SECRET,     // Mismo secret que en .env del servidor
  { expiresIn: '24h' }        // Expiración
);
```

---

## **🔄 Flujo de Autenticación**

### **Paso a paso:**

1. **Cliente envía request** con headers:
   ```
   X-Client-Id: acricolor
   Authorization: Bearer jwt_token_here
   ```

2. **Sistema verifica**:
   - ¿El cliente `acricolor` existe? ✅
   - ¿El cliente `acricolor` requiere auth? ✅ SÍ
   - ¿Hay token en Authorization header? ✅
   - ¿El token es válido? ✅

3. **Si todo OK**: Request procesado
4. **Si falta algo**: Error 401 Unauthorized

---

## **❌ Respuestas de Error**

### **Cliente no requiere auth pero se envía token:**
✅ **Se ignora el token** - Request procede normalmente

### **Cliente requiere auth pero no hay token:**
```json
{
  "success": false,
  "error": "Token de autenticación requerido",
  "code": 401
}
```

### **Token inválido o expirado:**
```json
{
  "success": false,
  "error": "Token inválido: signature invalid",
  "code": 401
}
```

### **Cliente no existe:**
```json
{
  "success": false,
  "error": "Cliente no válido: cliente_inexistente",
  "code": 400
}
```

---

## **🧪 Testing de Autenticación**

### **Test 1: Cliente sin auth (shared)**
```bash
# Debería funcionar SIN token
curl -X POST \
  -H "X-Client-Id: shared" \
  -F "file=@test.pdf" \
  http://localhost:4040/api/files/upload
```

### **Test 2: Cliente con auth SIN token (debería fallar)**
```bash
# Debería dar error 401
curl -X POST \
  -H "X-Client-Id: acricolor" \
  -F "file=@test.pdf" \
  http://localhost:4040/api/files/upload
```

### **Test 3: Cliente con auth CON token válido**
```bash
# Debería funcionar
curl -X POST \
  -H "X-Client-Id: acricolor" \
  -H "Authorization: Bearer YOUR_VALID_JWT_TOKEN" \
  -F "file=@test.pdf" \
  http://localhost:4040/api/files/upload
```

---

## **⚙️ Configuración del JWT Secret**

### **En tu .env:**
```bash
# IMPORTANTE: Cambiar en producción
JWT_SECRET=tu_jwt_secret_super_seguro_de_al_menos_32_caracteres
```

### **Generar secret seguro:**
```bash
# Opción 1: OpenSSL
openssl rand -base64 32

# Opción 2: Node.js
node -e "console.log(require('crypto').randomBytes(32).toString('base64'))"

# Opción 3: Online
# https://generate-secret.now.sh/32
```

---

## **🔧 Integración con tu Sistema de Auth**

### **Si tienes un sistema de autenticación existente:**

1. **Usa el mismo JWT_SECRET** en ambos sistemas
2. **Genera tokens** con el mismo formato
3. **Incluye user_id** en el payload del token

### **Ejemplo de integración:**
```javascript
// En tu sistema de auth existente
const generateFileServerToken = (userId, permissions = []) => {
  return jwt.sign(
    {
      sub: userId,                    // ID del usuario
      permissions: permissions,       // Permisos específicos
      iss: 'tu-sistema-principal',   // Emisor
      aud: 'file-server'             // Audiencia
    },
    process.env.JWT_SECRET,
    { expiresIn: '24h' }
  );
};

// Uso
const fileToken = generateFileServerToken('user123', ['upload', 'download']);
```

---

## **🚨 Consideraciones de Seguridad**

### **✅ Buenas Prácticas:**
- Usa HTTPS en producción
- Rota el JWT_SECRET regularmente  
- Tokens de corta duración (max 24h)
- Valida permisos específicos si es necesario
- Logs de acceso y auditoría

### **⚠️ Importante:**
- El cliente `shared` NO requiere auth (por diseño)
- Los archivos estáticos (`/static/`) no tienen auth (nginx los sirve directamente)
- Para proteger archivos estáticos también, usar solo `/api/files/download/`

---

## **💡 Casos de Uso**

### **Caso 1: App pública con archivos compartidos**
```javascript
// Usar cliente 'shared' - sin auth
const uploadPublic = (file) => {
  return fetch('/api/files/upload', {
    method: 'POST',
    headers: { 'X-Client-Id': 'shared' },
    body: formData
  });
};
```

### **Caso 2: App privada con usuarios autenticados**
```javascript
// Usar cliente específico con auth
const uploadPrivate = (file, userToken) => {
  return fetch('/api/files/upload', {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${userToken}`,
      'X-Client-Id': 'acricolor'
    },
    body: formData
  });
};
```

---

**🎯 ¡El sistema de auth está listo y funcionando!**
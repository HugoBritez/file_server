# File Server Microservice

**Descripción:** Microservicio de archivos estáticos con API REST multi-cliente, arquitectura Docker + nginx + Go para gestión y distribución de archivos.

**Tecnologías:** Go, nginx, Docker, Docker Compose, JWT, REST API, HTML/CSS/JavaScript, SSL/HTTPS

**Características principales:**
- Arquitectura multi-cliente con configuración específica por cliente
- API REST completa para upload/download/listado de archivos
- Proxy reverso con nginx para archivos estáticos de alto rendimiento
- Autenticación JWT opcional por cliente (acricolor, lobeck, gaesa, shared)
- Límites de tamaño y tipos de archivo configurables por cliente
- Streaming de archivos grandes con range requests
- Scripts automatizados de deploy y backup

**Logros técnicos:** Sistema escalable que soporta múltiples clientes con diferentes niveles de seguridad, optimizado para archivos estáticos con nginx, API REST completa con autenticación granular, contenedorización completa con Docker y scripts de automatización. 
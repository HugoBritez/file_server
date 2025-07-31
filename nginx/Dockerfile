FROM nginx:alpine

# Eliminar la configuración por defecto
RUN rm -f /etc/nginx/conf.d/default.conf

# Copiar nuestra configuración
COPY nginx.conf /etc/nginx/nginx.conf

# Crear directorios necesarios
RUN mkdir -p /var/www/frontend /var/www/static /etc/nginx/ssl

EXPOSE 4040 4043

CMD ["nginx", "-g", "daemon off;"] 
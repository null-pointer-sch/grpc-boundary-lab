#!/bin/sh
set -e

TLS_ENABLED=${TLS:-0}

if [ "$TLS_ENABLED" = "1" ]; then
    echo "Frontend TLS is enabled. Generating self-signed certificate..."
    openssl req -x509 -nodes -days 1 -newkey rsa:2048 \
        -keyout /etc/nginx/cert.key -out /etc/nginx/cert.crt \
        -subj "/C=US/ST=State/L=City/O=Organization/CN=localhost" 2>/dev/null

    # Switch Nginx to HTTPS and update the proxy_pass to use https
    sed -i 's/listen 80;/listen 443 ssl;\n    ssl_certificate \/etc\/nginx\/cert.crt;\n    ssl_certificate_key \/etc\/nginx\/cert.key;/g' /etc/nginx/conf.d/default.conf
    sed -i 's/proxy_pass http:\/\/gateway:8080\/api\/;/proxy_pass https:\/\/gateway:8080\/api\/;/g' /etc/nginx/conf.d/default.conf
else
    echo "Frontend TLS is disabled."
fi

echo "Starting Nginx..."
exec nginx -g "daemon off;"

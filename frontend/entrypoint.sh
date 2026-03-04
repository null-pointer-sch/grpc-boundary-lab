#!/bin/sh
set -e

TLS_ENABLED=${TLS:-0}

if [ "$TLS_ENABLED" = "1" ]; then
    echo "Frontend TLS is enabled. Configuring Nginx with shared certificates..."

    # Switch Nginx to HTTPS and update the proxy_pass to use https
    sed -i 's/listen 80;/listen 443 ssl;\n    ssl_certificate \/certs\/frontend.crt;\n    ssl_certificate_key \/certs\/frontend.key;/g' /etc/nginx/conf.d/default.conf
    sed -i 's/proxy_pass http:\/\/gateway:8080\/api\/;/proxy_pass https:\/\/gateway:8080\/api\/;/g' /etc/nginx/conf.d/default.conf
else
    echo "Frontend TLS is disabled."
fi

echo "Starting Nginx..."
exec nginx -g "daemon off;"

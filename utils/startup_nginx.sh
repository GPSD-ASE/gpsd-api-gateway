#!/bin/sh

# Ensure SSL certificates exist (for local testing or fallback)
if [ ! -f /etc/nginx/certs/tls.crt ] || [ ! -f /etc/nginx/certs/tls.key ]; then
    echo "Generating self-signed SSL certificates..."
    mkdir -p /etc/nginx/certs
    openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
        -keyout /etc/nginx/certs/tls.key \
        -out /etc/nginx/certs/tls.crt \
        -subj "/CN=api.gpsd.com"
fi

nginx -g "daemon off;"
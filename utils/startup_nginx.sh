#!/bin/sh

# # Ensure SSL certificates exist (for local testing or fallback)
# if [ ! -f /etc/nginx/certs/api.gpsd.com.crt ] || [ ! -f /etc/nginx/certs/api.gpsd.com.key ]; then
#     echo "Generating self-signed SSL certificates..."
#     mkdir -p /etc/nginx/certs
#     openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
#         -keyout /etc/nginx/certs/api.gpsd.com.key \
#         -out /etc/nginx/certs/api.gpsd.com.crt \
#         -subj "/CN=api.gpsd.com"
# fi

nginx -g "daemon off;"
#!/bin/sh
set -e

# If LEGO_DOMAIN is not set, run nginx in plain HTTP mode (development)
if [ -z "$LEGO_DOMAIN" ]; then
    echo "[ssl-entrypoint] LEGO_DOMAIN not set, running in HTTP mode."
    exec nginx -g "daemon off;"
fi

# --- SSL mode ---
DOMAIN="$LEGO_DOMAIN"
EMAIL="${LEGO_EMAIL:-admin@example.com}"
LEGO_DATA="/etc/lego"
CERT_DIR="${LEGO_DATA}/certificates"
CERT_FILE="${CERT_DIR}/${DOMAIN}.crt"
ACME_DIR="/var/www/acme"

mkdir -p "$ACME_DIR"

# Switch to SSL nginx config (render template with envsubst)
envsubst '${LEGO_DOMAIN}' < /etc/nginx/ssl.conf > /etc/nginx/conf.d/default.conf

# Obtain certificate if it doesn't exist yet
if [ ! -f "$CERT_FILE" ]; then
    echo "[ssl-entrypoint] Obtaining certificate for ${DOMAIN}..."
    lego --accept-tos \
         --email="$EMAIL" \
         --domains="$DOMAIN" \
         --path="$LEGO_DATA" \
         --http \
         --http.port ":80" \
         run
    echo "[ssl-entrypoint] Certificate obtained."
else
    # Attempt renewal on startup in case cert is near expiry
    echo "[ssl-entrypoint] Checking certificate renewal on startup..."
    lego --accept-tos \
         --email="$EMAIL" \
         --domains="$DOMAIN" \
         --path="$LEGO_DATA" \
         --http \
         --http.webroot "$ACME_DIR" \
         renew --days 30 || true
fi

# Background renewal loop (every 7 days, renew if <30 days remaining)
(
    while true; do
        sleep 604800
        echo "[ssl-entrypoint] Attempting certificate renewal..."
        if lego --accept-tos \
                --email="$EMAIL" \
                --domains="$DOMAIN" \
                --path="$LEGO_DATA" \
                --http \
                --http.webroot "$ACME_DIR" \
                renew --days 30; then
            echo "[ssl-entrypoint] Reloading nginx with new certificate..."
            nginx -s reload
        fi
    done
) &

exec nginx -g "daemon off;"

#!/bin/bash
# Aktifkan mode maintenance (Nginx akan balas 503 + halaman maintenance.html)
# Sekaligus generate token bypass baru supaya Anda tetap bisa mengetes aplikasi selagi maintenance aktif.
set -e

FLAG_FILE="/var/www/pos-web/maintenance.flag"
TOKEN_FILE="/etc/nginx/maintenance_token.conf"
DOMAIN="${1:-pos.domain-anda.com}"

TOKEN=$(openssl rand -hex 16)
echo "set \$maintenance_token \"$TOKEN\";" | sudo tee "$TOKEN_FILE" > /dev/null

sudo touch "$FLAG_FILE"

sudo nginx -t
sudo systemctl reload nginx

echo "Maintenance mode: ON"
echo "Flag file: $FLAG_FILE"
echo ""
echo "URL bypass (buka SEKALI di browser Anda untuk mulai testing, cookie berlaku 24 jam):"
echo "https://${DOMAIN}/maintenance-bypass?token=${TOKEN}"

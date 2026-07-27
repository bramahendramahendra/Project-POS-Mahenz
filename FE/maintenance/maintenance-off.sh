#!/bin/bash
# Matikan mode maintenance (Nginx kembali melayani request normal)
# Sekaligus merotasi token bypass supaya cookie mnt_bypass lama otomatis tidak berlaku lagi.
set -e

FLAG_FILE="/var/www/pos-web/maintenance.flag"
TOKEN_FILE="/etc/nginx/maintenance_token.conf"

sudo rm -f "$FLAG_FILE"

NEW_TOKEN=$(openssl rand -hex 16)
echo "set \$maintenance_token \"$NEW_TOKEN\";" | sudo tee "$TOKEN_FILE" > /dev/null

sudo nginx -t
sudo systemctl reload nginx

echo "Maintenance mode: OFF"
echo "Token bypass lama sudah di-rotasi (cookie mnt_bypass yang lama otomatis tidak berlaku lagi)."

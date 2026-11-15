#!/bin/bash
# Aktifkan mode maintenance (Nginx membalas 503 + halaman maintenance.html).
#
# Pemakaian:
#   maintenance-on.sh                 -> SEMUA aplikasi
#   maintenance-on.sh toko            -> hanya Toko
#   maintenance-on.sh portal          -> hanya Portal
#   maintenance-on.sh all DOMAIN      -> semua, dengan domain untuk URL bypass
#   maintenance-on.sh toko DOMAIN
#
# Sekaligus membuat token bypass baru supaya Anda tetap bisa mengetes aplikasi
# selagi maintenance aktif (cookie bypass berlaku untuk semua aplikasi).
set -e

APP="${1:-all}"
DOMAIN="${2:-mahenz.duckdns.org}"
WEB_ROOT="/var/www/pos-web"
TOKEN_FILE="/etc/nginx/maintenance_token.conf"

case "$APP" in
  all)    FLAG_FILE="$WEB_ROOT/maintenance.flag" ;;
  portal|toko|coffee) FLAG_FILE="$WEB_ROOT/maintenance-$APP.flag" ;;
  *) echo "Aplikasi tidak dikenal: $APP (pilihan: all, portal, toko, coffee)"; exit 1 ;;
esac

TOKEN=$(openssl rand -hex 16)
echo "set \$maintenance_token \"$TOKEN\";" | sudo tee "$TOKEN_FILE" > /dev/null

sudo touch "$FLAG_FILE"

sudo nginx -t
sudo systemctl reload nginx

echo "Maintenance mode: ON ($APP)"
echo "Flag file: $FLAG_FILE"
echo ""
echo "URL bypass (buka SEKALI di browser Anda untuk mulai testing, cookie berlaku 24 jam):"
echo "https://${DOMAIN}/maintenance-bypass?token=${TOKEN}"

#!/bin/bash
# Matikan mode maintenance (Nginx kembali melayani request normal).
#
# Pemakaian:
#   maintenance-off.sh                -> matikan maintenance SEMUA aplikasi (hapus semua flag)
#   maintenance-off.sh toko           -> hanya Toko
#   maintenance-off.sh portal         -> hanya Portal
#
# Token bypass dirotasi hanya bila sudah tidak ada flag maintenance yang tersisa,
# supaya bypass untuk aplikasi lain yang masih maintenance tidak ikut putus.
set -e

APP="${1:-all}"
WEB_ROOT="/var/www/pos-web"
TOKEN_FILE="/etc/nginx/maintenance_token.conf"

case "$APP" in
  all)
    sudo rm -f "$WEB_ROOT/maintenance.flag" "$WEB_ROOT"/maintenance-*.flag
    ;;
  portal|toko|coffee)
    sudo rm -f "$WEB_ROOT/maintenance-$APP.flag"
    ;;
  *)
    echo "Aplikasi tidak dikenal: $APP (pilihan: all, portal, toko, coffee)"; exit 1 ;;
esac

# Masih ada flag lain? Jika tidak, rotasi token bypass.
if ls "$WEB_ROOT"/maintenance.flag "$WEB_ROOT"/maintenance-*.flag > /dev/null 2>&1; then
  echo "Masih ada maintenance aktif untuk aplikasi lain; token bypass tidak dirotasi."
else
  NEW_TOKEN=$(openssl rand -hex 16)
  echo "set \$maintenance_token \"$NEW_TOKEN\";" | sudo tee "$TOKEN_FILE" > /dev/null
  echo "Token bypass lama sudah di-rotasi (cookie mnt_bypass yang lama otomatis tidak berlaku lagi)."
fi

sudo nginx -t
sudo systemctl reload nginx

echo "Maintenance mode: OFF ($APP)"

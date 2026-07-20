#!/bin/bash
# Aktifkan mode maintenance (Nginx akan balas 503 + halaman maintenance.html)
set -e

FLAG_FILE="/var/www/pos-web/maintenance.flag"

touch "$FLAG_FILE"
echo "Maintenance mode: ON"
echo "Flag file: $FLAG_FILE"

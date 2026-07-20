#!/bin/bash
# Matikan mode maintenance (Nginx kembali melayani request normal)
set -e

FLAG_FILE="/var/www/pos-web/maintenance.flag"

rm -f "$FLAG_FILE"
echo "Maintenance mode: OFF"

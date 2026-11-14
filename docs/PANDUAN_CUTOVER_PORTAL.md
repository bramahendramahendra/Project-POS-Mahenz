# Cutover VPS: BE User + Portal + Toko + Coffee

Sekali jalan. Copas per blok, cek baris `# harus:`. Kalau hasil beda: berhenti, kirim hasilnya.

---

**0. Cek dulu (hanya membaca, boleh sehari sebelumnya)**

```bash
go version                      # harus: go1.24 atau lebih baru
node -v                         # harus: v20.19+ atau v22.12+
systemctl status pos-backend --no-pager
ls /opt/pos-mahenz              # harus: BE  FE
ls -l /etc/nginx/sites-enabled/ # harus: ada pos-web (kalau beda nama, ganti "pos-web" di langkah 9)
sudo mysql -u root -p -e "SELECT user, host FROM mysql.user WHERE user='pos_user';"
```

---

**1. Maintenance + backup**

```bash
TGL=$(date +%Y%m%d)
sudo maintenance-on.sh mahenz.duckdns.org
mkdir -p ~/cutover_$TGL
mysqldump -u pos_user -p --no-tablespaces pos_retail_db > ~/cutover_$TGL/pos_retail_db_sebelum.sql
ls -lh ~/cutover_$TGL/          # harus: file .sql tidak 0 byte
```

**2. Matikan BE lama, ambil kode baru**

```bash
sudo systemctl stop pos-backend
sudo systemctl disable pos-backend
cd /opt
sudo mv pos-mahenz pos-mahenz_lama_$TGL
sudo git clone https://github.com/bramahendramahendra/Project-POS-Mahenz.git pos-mahenz
sudo chown -R $USER:$USER /opt/pos-mahenz
ls /opt/pos-mahenz              # harus: AUTH  COFFEE  PORTAL  TOKO  nginx ...

OLD=/opt/pos-mahenz_lama_$TGL
cp -r $OLD/BE/logs    /opt/pos-mahenz/TOKO/BE/logs    2>/dev/null || true
cp -r $OLD/BE/backups /opt/pos-mahenz/TOKO/BE/backups 2>/dev/null || true

diff $OLD/BE/config/config_prod.json /opt/pos-mahenz/TOKO/BE/config/config_prod.json
# harus: SecretKey dan Database TIDAK muncul di diff. Kalau muncul, samakan pakai nano.
```

**3. Database baru**

```bash
sudo mysql -u root -p
```

```sql
CREATE DATABASE IF NOT EXISTS pos_auth_db   CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE DATABASE IF NOT EXISTS pos_portal_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE DATABASE IF NOT EXISTS pos_coffee_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
GRANT ALL PRIVILEGES ON pos_auth_db.*   TO 'pos_user'@'localhost';
GRANT ALL PRIVILEGES ON pos_portal_db.* TO 'pos_user'@'localhost';
GRANT ALL PRIVILEGES ON pos_coffee_db.* TO 'pos_user'@'localhost';
-- kalau di langkah 0 pos_user juga punya host '%', ulangi 3 GRANT di atas dengan 'pos_user'@'%'
FLUSH PRIVILEGES;
EXIT;
```

**4. File .env**

```bash
for a in AUTH PORTAL TOKO COFFEE; do cp /opt/pos-mahenz/$a/BE/.env.prod /opt/pos-mahenz/$a/BE/.env; done
grep -H APP_PORT /opt/pos-mahenz/*/BE/.env   # harus: AUTH 8080, PORTAL 8081, TOKO 8082, COFFEE 8083
grep -H '"Path"' /opt/pos-mahenz/*/BE/config/config_prod.json   # harus: relatif, mis. "logs" (tanpa / di depan)
```

**5. Build BE**

```bash
# wajib build di server, di folder ini (biner mencari folder migrasi dari lokasi build). Jangan salin biner dari laptop.
cd /opt/pos-mahenz/AUTH/BE   && go build -o auth_api .
cd /opt/pos-mahenz/PORTAL/BE && go build -o portal_api .
cd /opt/pos-mahenz/TOKO/BE   && go build -o pos_api .
cd /opt/pos-mahenz/COFFEE/BE && go build -o coffee_api .
sudo chown -R www-data:www-data /opt/pos-mahenz/AUTH/BE /opt/pos-mahenz/PORTAL/BE /opt/pos-mahenz/TOKO/BE /opt/pos-mahenz/COFFEE/BE
```

**6. Service systemd (copas sekaligus)**

```bash
buat_service() {
sudo tee /etc/systemd/system/$1.service > /dev/null <<EOF
[Unit]
Description=$1
After=network.target mysql.service

[Service]
User=www-data
Group=www-data
WorkingDirectory=/opt/pos-mahenz/$2/BE
ExecStart=/opt/pos-mahenz/$2/BE/$3
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF
}
buat_service pos-auth-backend   AUTH   auth_api
buat_service pos-portal-backend PORTAL portal_api
buat_service pos-toko-backend   TOKO   pos_api
buat_service pos-coffee-backend COFFEE coffee_api
sudo systemctl daemon-reload
```

**7. Nyalakan (urutan jangan diubah)**

```bash
sudo systemctl enable --now pos-auth-backend
sleep 3; curl -s http://127.0.0.1:8080/api/health    # harus: balasan JSON

sudo systemctl enable --now pos-toko-backend
sleep 3; curl -s http://127.0.0.1:8082/api/health    # harus: balasan JSON
sudo journalctl -u pos-toko-backend -n 30 --no-pager | grep -i error   # harus: kosong
```

Salin akun kasir Toko ke BE User:

```bash
mysqldump -u pos_user -p --no-tablespaces pos_retail_db > ~/cutover_$TGL/pos_retail_db_sebelum_akun.sql
cd /opt/pos-mahenz/TOKO/BE
go run ./cmd/migrate_users_to_auth
# harus (hasil gladi): BUAT 2 (Banana, BramaGanteng), AMBIL ALIH 1 (admin), KONFLIK 0
# kalau ada KONFLIK: berhenti, kirim hasilnya.
go run ./cmd/migrate_users_to_auth --apply --i-have-backup
go run ./cmd/migrate_users_to_auth   # harus: SUDAH TERTAUT 3
```

```bash
sudo systemctl enable --now pos-portal-backend
sudo systemctl enable --now pos-coffee-backend
sleep 3; curl -s http://127.0.0.1:8081/api/health; curl -s http://127.0.0.1:8083/api/health   # harus: 2 balasan JSON
```

**8. Kunci aplikasi**

Kunci `pak_...` hanya tampil sekali. Langsung tempel ke perintah di bawahnya.

```bash
cd /opt/pos-mahenz/AUTH/BE

go run ./cmd/set_app_key --app portal --confirm
sudo mysql -u root -p pos_portal_db -e "INSERT INTO auth_credentials (id, app_id, app_key) VALUES (1, 'portal', 'TEMPEL_KUNCI') ON DUPLICATE KEY UPDATE app_key=VALUES(app_key);"

go run ./cmd/set_app_key --app toko --confirm
sudo mysql -u root -p pos_retail_db -e "INSERT INTO auth_credentials (id, app_id, app_key) VALUES (1, 'toko', 'TEMPEL_KUNCI') ON DUPLICATE KEY UPDATE app_key=VALUES(app_key);"

go run ./cmd/set_app_key --app coffee --confirm
sudo mysql -u root -p pos_coffee_db -e "INSERT INTO auth_credentials (id, app_id, app_key) VALUES (1, 'coffee', 'TEMPEL_KUNCI') ON DUPLICATE KEY UPDATE app_key=VALUES(app_key);"

for a in portal toko coffee; do go run ./cmd/set_app_key --app $a | grep 'sudah punya kunci'; done   # harus: 3 baris "sudah punya kunci=true"
```

**9. FE + maintenance baru + nginx**

```bash
cd /opt/pos-mahenz/PORTAL/FE && npm ci && npm run build
cd /opt/pos-mahenz/TOKO/FE   && npm ci && npm run build
cd /opt/pos-mahenz/COFFEE/FE && npm ci && npm run build

sudo mkdir -p /var/www/pos-web/portal /var/www/pos-web/apps/toko /var/www/pos-web/apps/coffee /var/www/pos-web/maintenance
sudo cp -r /opt/pos-mahenz/PORTAL/FE/dist/* /var/www/pos-web/portal/
sudo cp -r /opt/pos-mahenz/TOKO/FE/dist/*   /var/www/pos-web/apps/toko/
sudo cp -r /opt/pos-mahenz/COFFEE/FE/dist/* /var/www/pos-web/apps/coffee/
sudo chown -R www-data:www-data /var/www/pos-web/portal /var/www/pos-web/apps

sudo cp /opt/pos-mahenz/nginx/maintenance/maintenance-on.sh  /usr/local/bin/
sudo cp /opt/pos-mahenz/nginx/maintenance/maintenance-off.sh /usr/local/bin/
sudo chmod +x /usr/local/bin/maintenance-*.sh
sudo cp /opt/pos-mahenz/nginx/maintenance/maintenance.html /var/www/pos-web/maintenance/

[ -f /etc/nginx/maintenance_token.conf ] || echo 'set $maintenance_token "";' | sudo tee /etc/nginx/maintenance_token.conf
sudo touch /var/www/pos-web/maintenance.flag

sudo cp /etc/nginx/sites-available/pos-web ~/cutover_$TGL/pos-web.nginx.lama
sudo cp /opt/pos-mahenz/nginx/pos-web.conf /etc/nginx/sites-available/pos-web
sudo nginx -t
# harus: syntax is ok + test is successful
# kalau GAGAL: jangan lanjut, jalankan:
#   sudo cp ~/cutover_$TGL/pos-web.nginx.lama /etc/nginx/sites-available/pos-web
```

```bash
sudo systemctl reload nginx
sudo certbot --nginx
sudo nginx -t && sudo systemctl reload nginx
sudo maintenance-on.sh all mahenz.duckdns.org   # catat URL bypass yang muncul
```

**10. Cek di browser (buka URL bypass dulu)**

- [ ] `https://mahenz.duckdns.org/` → login Portal pakai `admin`
- [ ] Kalau password `admin` masih `admin123` → ganti di `/account/password`
- [ ] Kartu POS Toko → masuk Toko tanpa login ulang
- [ ] `/toko/login` → kasir lama login pakai password lama
- [ ] Toko: transaksi kecil + cetak struk, lalu void
- [ ] Portal `/admin/users` → `admin` → **Akses aplikasi** → pastikan POS Coffee tercentang (bawaan sudah)
- [ ] Kartu POS Coffee → masuk Coffee → buat shift di menu Shift
- [ ] Jangan beri akses Coffee ke kasir dulu

```bash
sudo journalctl -u pos-auth-backend -u pos-portal-backend -u pos-toko-backend -u pos-coffee-backend -n 50 --no-pager | grep -i error   # harus: kosong
sudo maintenance-off.sh
```

---

**Rollback (kalau gagal parah)**

```bash
sudo maintenance-on.sh all mahenz.duckdns.org
sudo systemctl disable --now pos-coffee-backend pos-portal-backend pos-toko-backend pos-auth-backend
sudo cp ~/cutover_$TGL/pos-web.nginx.lama /etc/nginx/sites-available/pos-web
sudo nginx -t && sudo systemctl reload nginx
cd /opt && sudo mv pos-mahenz pos-mahenz_gagal_$TGL && sudo mv pos-mahenz_lama_$TGL pos-mahenz
sudo systemctl enable --now pos-backend
sudo maintenance-off.sh
# Database Toko tidak perlu dipulihkan. Kalau data rusak: restore ~/cutover_$TGL/pos_retail_db_sebelum.sql
```

Catatan: `$TGL` hilang kalau sesi SSH tertutup. Isi ulang dengan tanggal hari cutover, misalnya `TGL=20261120`.

---

**Beberapa hari setelah stabil**

```bash
sudo rm -rf /opt/pos-mahenz_lama_$TGL /var/www/pos-web/dist
sudo rm /etc/systemd/system/pos-backend.service && sudo systemctl daemon-reload
```

- Beri akses Coffee ke kasir di Portal (`/admin/users` → Akses aplikasi), lalu atur role-nya di Coffee.
- Hapus blok `location ^~ /api/` di `nginx/pos-web.conf` setelah semua perangkat memakai `/toko/api/`.
- Deploy berikutnya: [SOP_DEPLOY.md](SOP_DEPLOY.md) bagian "SOP Deploy per Aplikasi".

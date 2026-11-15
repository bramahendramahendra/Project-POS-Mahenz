> Catatan lama (struktur sebelum cutover: folder `BE/` dan `FE/`, Toko di root). Setelah cutover pakai bagian
> "SOP Deploy per Aplikasi" di bawah.

**1. Aktifkan maintenance mode**

```bash
sudo maintenance-on.sh pos.domain-anda.com
sudo maintenance-on.sh 139.180.214.187
```

**2. Duplicate database ke nama baru bertanggal**

```bash
sudo mysql -u root -e "SHOW DATABASES;"
sudo mysql -u pos_user -p -e "SHOW DATABASES;"
```

<!-- Jika ada db yang ingin dihapus : -->

```bash
sudo mysql -u root -p -e "DROP DATABASE pos_retail_db_20260907;"
```

<!-- Cek host yang dipakai apa aja -->
```bash
sudo mysql -u root -p

# -- Cek dulu host mana saja yang dipakai pos_user:
SELECT user, host FROM mysql.user WHERE user = 'pos_user';

# Create table backup
CREATE DATABASE pos_retail_db_20260930 CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
GRANT ALL PRIVILEGES ON pos_retail_db_20260930.* TO 'pos_user'@'localhost';
GRANT ALL PRIVILEGES ON pos_retail_db_20260930.* TO 'pos_user'@'%';

# -- Verifikasi kedua host sudah punya izin ke database baru
SHOW GRANTS FOR 'pos_user'@'localhost';
SHOW GRANTS FOR 'pos_user'@'%';

# Hapus izin di kedua host
REVOKE ALL PRIVILEGES ON `pos_retail_db_20260907`.* FROM 'pos_user'@'localhost';
REVOKE ALL PRIVILEGES ON `pos_retail_db_20260907`.* FROM 'pos_user'@'%';

EXIT;
```

> Memindahkan data db di db utama ke db backup.
```bash
read -s -p "Password pos_user: " DBPASS && echo
mysqldump -u pos_user -p"$DBPASS" --no-tablespaces pos_retail_db | mysql -u pos_user -p"$DBPASS" pos_retail_db_20260930
unset DBPASS
```

Cek data
```bash
sudo mysql -u pos_user -p
```
```sql
-- jalankan di prompt mysql untuk membandingkan langsung dalam satu query
SELECT
  (SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'pos_retail_db')          AS tabel_lama,
  (SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'pos_retail_db_20260930') AS tabel_baru;
```

## Langkah 2 — Backup Database (WAJIB)

```bash
sudo mkdir -p /opt/pos-mahenz/BE/backups

mysqldump -u pos_user -p --no-tablespaces pos_retail_db > ~/backup_sebelum_laba_rugi_$(date +%Y%m%d_%H%M%S).sql
sudo mv ~/backup_sebelum_laba_rugi_*.sql /opt/pos-mahenz/BE/backups/
```

Pastikan file backup benar-benar terisi (ukuran tidak 0 byte):

```bash
ls -lh /opt/pos-mahenz/BE/backups/
```

**4. Rename folder induk lama, clone ulang penuh yang baru**

# jika mau rename saja
```bash
sudo mv pos-mahenz pos-mahenz_20260930

sudo mkdir -p /opt/pos-mahenz
sudo git clone https://github.com/bramahendramahendra/Project-POS-Mahenz.git /opt/pos-mahenz
sudo chown -R $USER:$USER /opt/pos-mahenz
```

# JIka mau di copy aja
```bash
sudo cp -a pos-mahenz pos-mahenz_20260930

# ambil kodingan 
git satsh # jika tidak bisa pull
git pull
sudo chown -R $USER:$USER /opt/sipaduke-testing
```

**5. Setup ulang BE seperti deploy pertama kali**

```bash
cd /opt/pos-mahenz/TOKO/BE
go mod tidy
go build -o pos_api .
sudo systemctl restart pos-toko-backend
sudo systemctl status pos-toko-backend    # harus "active (running)", port 8081
```
**Frontend (tampilan):**

```bash
cd /opt/pos-mahenz/TOKO/FE
npm install
npm run type-check
npm run lint
npm run build


sudo rm -rf /var/www/pos-web/dist
sudo cp -r dist /var/www/pos-web/
sudo chown -R $USER:$USER /var/www/pos-web
```

**7. Verifikasi & matikan maintenance mode**
```bash
sudo maintenance-off.sh
```


---

# SOP Deploy per Aplikasi (setelah cutover)

Deploy satu aplikasi saja, yang lain tetap jalan. Copas bagian aplikasi yang di-update.

```bash
sudo maintenance-on.sh toko mahenz.duckdns.org     # ganti: portal | toko | coffee | all
# ... jalankan bagian aplikasinya di bawah ...
sudo maintenance-off.sh toko
```

BE User dipakai semua aplikasi untuk login: kalau update BE User, pakai `maintenance-on.sh all`.

## POS Toko

### Backend (service `pos-toko-backend`, port 8082, database `pos_retail_db`)

```bash
mysqldump -u pos_user -p --no-tablespaces pos_retail_db > ~/pos_retail_db_$(date +%Y%m%d_%H%M%S).sql
cd /opt/pos-mahenz && git pull
cd TOKO/BE && go build -o pos_api .
sudo systemctl restart pos-toko-backend
sudo systemctl status pos-toko-backend --no-pager     # harus "active (running)"
curl -s http://127.0.0.1:8082/api/health
```

### Frontend (alamat `/toko/`, folder `/var/www/pos-web/apps/toko`)

```bash
cd /opt/pos-mahenz/TOKO/FE
npm ci && npm run build   # type-check dan lint sudah dijalankan di laptop sebelum push
sudo rm -rf /var/www/pos-web/apps/toko/*
sudo cp -r dist/* /var/www/pos-web/apps/toko/
sudo chown -R www-data:www-data /var/www/pos-web/apps/toko
```

## POS Coffee

### Backend (service `pos-coffee-backend`, port 8083, database `pos_coffee_db`)

```bash
mysqldump -u pos_user -p --no-tablespaces pos_coffee_db > ~/pos_coffee_db_$(date +%Y%m%d_%H%M%S).sql
cd /opt/pos-mahenz && git pull
cd COFFEE/BE && go build -o coffee_api .
sudo systemctl restart pos-coffee-backend
sudo systemctl status pos-coffee-backend --no-pager     # harus "active (running)"
curl -s http://127.0.0.1:8083/api/health
```

### Frontend (alamat `/coffee/`, folder `/var/www/pos-web/apps/coffee`)

```bash
cd /opt/pos-mahenz/COFFEE/FE
npm ci && npm run build   # type-check dan lint sudah dijalankan di laptop sebelum push
sudo rm -rf /var/www/pos-web/apps/coffee/*
sudo cp -r dist/* /var/www/pos-web/apps/coffee/
sudo chown -R www-data:www-data /var/www/pos-web/apps/coffee
```

## Portal

### Backend (service `pos-portal-backend`, port 8081, database `pos_portal_db`)

```bash
mysqldump -u pos_user -p --no-tablespaces pos_portal_db > ~/pos_portal_db_$(date +%Y%m%d_%H%M%S).sql
cd /opt/pos-mahenz && git pull
cd PORTAL/BE && go build -o portal_api .
sudo systemctl restart pos-portal-backend
sudo systemctl status pos-portal-backend --no-pager     # harus "active (running)"
curl -s http://127.0.0.1:8081/api/health
```

### Frontend (alamat `/`, folder `/var/www/pos-web/portal`)

```bash
cd /opt/pos-mahenz/PORTAL/FE
npm ci && npm run build   # type-check dan lint sudah dijalankan di laptop sebelum push
sudo rm -rf /var/www/pos-web/portal/*
sudo cp -r dist/* /var/www/pos-web/portal/
sudo chown -R www-data:www-data /var/www/pos-web/portal
```

## BE User

### Backend (service `pos-auth-backend`, port 8080, database `pos_auth_db`)

```bash
mysqldump -u pos_user -p --no-tablespaces pos_auth_db > ~/pos_auth_db_$(date +%Y%m%d_%H%M%S).sql
cd /opt/pos-mahenz && git pull
cd AUTH/BE && go build -o auth_api .
sudo systemctl restart pos-auth-backend
sudo systemctl status pos-auth-backend --no-pager     # harus "active (running)"
curl -s http://127.0.0.1:8080/api/health
```

## Nginx (bila `nginx/pos-web.conf` berubah)

```bash
sudo cp /etc/nginx/sites-available/pos-web ~/pos-web.nginx.$(date +%Y%m%d)
sudo cp /opt/pos-mahenz/nginx/pos-web.conf /etc/nginx/sites-available/pos-web
sudo nginx -t                      # bila gagal: JANGAN reload, kembalikan file cadangan
sudo systemctl reload nginx
sudo certbot --nginx               # memasang kembali blok HTTPS
sudo nginx -t && sudo systemctl reload nginx
```

## Ganti kunci aplikasi (kalau perlu)

```bash
cd /opt/pos-mahenz/AUTH/BE
go run ./cmd/set_app_key --app toko --confirm      # portal | toko | coffee
sudo mysql -u root -p pos_retail_db -e "INSERT INTO auth_credentials (id, app_id, app_key) VALUES (1, 'toko', 'TEMPEL_KUNCI') ON DUPLICATE KEY UPDATE app_key=VALUES(app_key);"
# database: portal -> pos_portal_db, toko -> pos_retail_db, coffee -> pos_coffee_db
```

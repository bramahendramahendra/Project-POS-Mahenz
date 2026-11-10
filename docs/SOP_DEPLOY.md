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

# SOP Deploy per Aplikasi (setelah Portal dipasang)

Berlaku setelah cutover ([PANDUAN_CUTOVER_PORTAL.md](PANDUAN_CUTOVER_PORTAL.md)). Tiap aplikasi di-deploy sendiri;
deploy Toko tidak menyentuh Portal, begitu juga sebaliknya. Langkah umum: backup database aplikasi itu,
maintenance **hanya aplikasi itu**, tarik kode, build, restart service / salin FE, verifikasi, maintenance off.

Perintah maintenance per aplikasi (skrip ada di `nginx/maintenance/`, terpasang di `/usr/local/bin`):

```bash
sudo maintenance-on.sh toko mahenz.duckdns.org     # hanya Toko (alternatif: portal, all)
sudo maintenance-off.sh toko
```

## POS Toko

Backend (service `pos-toko-backend`, port 8081):

```bash
cd /opt/pos-mahenz
mysqldump -u pos_user -p --no-tablespaces pos_retail_db > ~/pos_retail_db_$(date +%Y%m%d_%H%M%S).sql
git pull
cd TOKO/BE
go mod tidy
go build -o pos_api .
sudo systemctl restart pos-toko-backend
sudo systemctl status pos-toko-backend      # harus "active (running)"
```

Frontend (disajikan di `/toko/`, folder `/var/www/pos-web/apps/toko`):

```bash
cd /opt/pos-mahenz/TOKO/FE
npm ci
npm run type-check
npm run lint
npm run build                               # memakai .env.production (base /toko/)
sudo rm -rf /var/www/pos-web/apps/toko/*
sudo cp -r dist/* /var/www/pos-web/apps/toko/
sudo chown -R www-data:www-data /var/www/pos-web/apps/toko
```

## Portal

Backend (service `pos-portal-backend`, port 8080; database `pos_portal_db`):

```bash
cd /opt/pos-mahenz
mysqldump -u portal_user -p --no-tablespaces pos_portal_db > ~/pos_portal_db_$(date +%Y%m%d_%H%M%S).sql
git pull
cd PORTAL/BE
go mod tidy
go build -o portal_api .
sudo systemctl restart pos-portal-backend
sudo systemctl status pos-portal-backend
```

Frontend (root domain, folder `/var/www/pos-web/portal`):

```bash
cd /opt/pos-mahenz/PORTAL/FE
npm ci
npm run type-check
npm run lint
npm run build
sudo rm -rf /var/www/pos-web/portal/*
sudo cp -r dist/* /var/www/pos-web/portal/
sudo chown -R www-data:www-data /var/www/pos-web/portal
```

## Kunci SSO aplikasi (diisi lewat SQL)

Kunci dibuat di Portal (Kelola Aplikasi -> Buat kunci baru, tampil sekali), lalu dipasang di database aplikasi tujuan:

```sql
INSERT INTO portal_credentials (id, app_id, app_key) VALUES (1, 'toko', '<KUNCI>')
ON DUPLICATE KEY UPDATE app_id = VALUES(app_id), app_key = VALUES(app_key);
```

Setelah kunci diganti di Portal, kunci lama langsung tidak berlaku: pasang kunci baru di aplikasi tujuan segera,
kalau tidak login lewat Portal ke aplikasi itu gagal. File backup SQL aplikasi memuat kunci ini, jadi perlakukan sebagai rahasia.

# Panduan Cutover VPS — Portal + POS Toko

Panduan ini memindahkan production dari "Toko di root domain" menjadi "Portal di root, Toko di `/toko/`".
Rancangan dan alasannya ada di [RANCANGAN_PORTAL_TOKO_COFFEE.md](RANCANGAN_PORTAL_TOKO_COFFEE.md).

> **Status:** panduan ini disusun dan **belum pernah dijalankan di VPS**. `nginx/pos-web.conf` belum
> divalidasi dengan `nginx -t` (nginx tidak tersedia di mesin pengembangan). Jalankan `sudo nginx -t` sebelum
> reload dan ikuti langkah verifikasi di tiap bagian. Jadwalkan di jam sepi: ada downtime beberapa menit.

## Gambaran akhir

| Bagian | Alamat | Service / folder | Port |
|--------|--------|------------------|------|
| FE Portal | `/` | `/var/www/pos-web/portal` | - |
| BE Portal | `/portal-api/` | `pos-portal-backend` | 8080 |
| FE Toko | `/toko/` | `/var/www/pos-web/apps/toko` | - |
| BE Toko | `/toko/api/` (dan `/api/` lama) | `pos-toko-backend` | 8081 |
| Database | - | `pos_portal_db` (baru), `pos_retail_db` (tetap) | - |

Perubahan penting dibanding sekarang: BE Toko pindah port **8080 ke 8081** dan berganti nama service.
FE lama di `/var/www/pos-web/dist` **dibiarkan** sebagai jalur rollback.

## 0. Prasyarat (siapkan sebelum jadwal)

1. Backup database Toko sekarang dan simpan di luar server juga:
   `mysqldump -u USER -p pos_retail_db > ~/pos_retail_db_sebelum_portal.sql`
2. Rahasia baru (jangan di-commit, simpan di password manager):
   - `SECRETKEY` JWT **Portal** (acak panjang, berbeda dari Toko).
   - Password MySQL untuk user `portal_user`.
3. Buat database dan user Portal (sekali):
   ```sql
   CREATE DATABASE pos_portal_db CHARACTER SET utf8mb4;
   CREATE USER 'portal_user'@'localhost' IDENTIFIED BY '<PASSWORD_KUAT>';
   GRANT ALL PRIVILEGES ON pos_portal_db.* TO 'portal_user'@'localhost';
   FLUSH PRIVILEGES;
   ```
4. Pastikan Go dan Node terpasang di VPS (seperti deploy biasa).
5. Pelajari rollback (bagian 7) sebelum mulai.

## 1. Mulai: maintenance dan cadangan

```bash
# Salin skrip maintenance versi baru (mendukung per aplikasi) dari repo yang sudah di-clone ulang nanti.
# Untuk langkah ini pakai skrip lama yang sudah terpasang:
sudo maintenance-on.sh mahenz.duckdns.org
```

Cadangkan juga file di luar git yang akan dipakai ulang: `.env`, `config/config_prod.json`, folder `logs/`
dan `backups/` milik BE lama.

## 2. Ganti kode (rename lalu clone ulang)

```bash
sudo systemctl stop pos-backend
sudo systemctl disable pos-backend

cd /opt
sudo mv pos-mahenz pos-mahenz_$(date +%Y%m%d)
git clone <URL_REPO> pos-mahenz
sudo chown -R $USER:$USER /opt/pos-mahenz
```

Struktur baru: `/opt/pos-mahenz/{PORTAL,TOKO}/{BE,FE}`.

### 2.1 Pindahkan file non-git Toko

```bash
OLD=/opt/pos-mahenz_$(date +%Y%m%d)      # sesuaikan dengan nama folder lama
cp $OLD/BE/.env                 /opt/pos-mahenz/TOKO/BE/.env
cp $OLD/BE/config/config_prod.json /opt/pos-mahenz/TOKO/BE/config/config_prod.json
cp -r $OLD/BE/logs    /opt/pos-mahenz/TOKO/BE/logs      2>/dev/null || true
cp -r $OLD/BE/backups /opt/pos-mahenz/TOKO/BE/backups   2>/dev/null || true
```

Lalu ubah di `/opt/pos-mahenz/TOKO/BE/.env`:

```env
APP_PORT=8081
```

Di `TOKO/BE/config/config_prod.json` pastikan ada (tambahkan bila belum, karena file lama dari server):

```json
"PortalUrl": "http://127.0.0.1:8080/api",
"DirectLoginForceEnabled": false,
```

### 2.2 Siapkan file Portal

```bash
cd /opt/pos-mahenz/PORTAL/BE
# .env.prod sudah ada di repo; salin sebagai .env
cp .env.prod .env
```

Edit `PORTAL/BE/config/config_prod.json`: ganti `Database.Password` dengan password `portal_user`
(`SecretKey` dibiarkan kosong: diisi lewat environment variable di service, langkah 3).

## 3. Build dan service

```bash
# BE Toko
cd /opt/pos-mahenz/TOKO/BE && go mod tidy && go build -o pos_api .
# BE Portal
cd /opt/pos-mahenz/PORTAL/BE && go mod tidy && go build -o portal_api .
sudo chown -R www-data:www-data /opt/pos-mahenz/TOKO/BE /opt/pos-mahenz/PORTAL/BE
```

Buat dua unit systemd (ganti nilai `SECRETKEY`; untuk Toko pakai **nilai yang sama seperti sebelumnya** agar
token yang sudah dipegang kasir tetap berlaku):

`/etc/systemd/system/pos-portal-backend.service`
```ini
[Unit]
Description=POS Portal API
After=network.target mysql.service

[Service]
User=www-data
Group=www-data
WorkingDirectory=/opt/pos-mahenz/PORTAL/BE
Environment="SECRETKEY=<SECRETKEY_PORTAL>"
ExecStart=/opt/pos-mahenz/PORTAL/BE/portal_api
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
```

`/etc/systemd/system/pos-toko-backend.service`
```ini
[Unit]
Description=POS Toko API
After=network.target mysql.service

[Service]
User=www-data
Group=www-data
WorkingDirectory=/opt/pos-mahenz/TOKO/BE
Environment="SECRETKEY=<SECRETKEY_TOKO_YANG_LAMA>"
ExecStart=/opt/pos-mahenz/TOKO/BE/pos_api
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
```

```bash
sudo systemctl daemon-reload
sudo systemctl enable --now pos-portal-backend      # migrasi Portal berjalan otomatis
sudo systemctl enable --now pos-toko-backend        # migrasi 011 Toko berjalan otomatis
sudo systemctl status pos-portal-backend pos-toko-backend
curl -s http://127.0.0.1:8080/api/health
curl -s http://127.0.0.1:8081/api/health
```

Migrasi Toko 011 bersifat **aditif** (kolom `portal_user_id`, role `menunggu_akses`, tabel
`portal_credentials`, setting `direct_login_enabled`): BE lama tetap bisa berjalan di atas skema ini bila perlu rollback.

## 4. Salin user Toko ke Portal

Dari `/opt/pos-mahenz/TOKO/BE`:

```bash
# Backup kedua database (WAJIB)
mysqldump -u USER -p pos_retail_db > ~/pos_retail_db_sebelum_tautan.sql
mysqldump -u USER -p pos_portal_db > ~/pos_portal_db_sebelum_tautan.sql

# Lihat rencana dulu (tidak menulis apa pun)
go run ./cmd/migrate_users_to_portal --portal-db pos_portal_db \
    # bila user Portal berbeda: tambahkan --portal-dsn 'portal_user:PASS@tcp(127.0.0.1:3306)/pos_portal_db?charset=utf8&parseTime=True&loc=Local'
```

Baca hasilnya:
- `BUAT`: akun akan dibuat di Portal (password sama dengan di Toko).
- `AMBIL ALIH`: akun seed Portal (`admin`) diganti dengan data akun Toko dengan username sama.
- `KONFLIK`: username sudah ada di Portal dengan password lain; selesaikan manual dulu.
- `PERINGATAN`: akun Portal yang masih memakai password bawaan; ganti atau hapus.

Aturan role: `admin` Toko menjadi `admin` Portal, semua role lain menjadi `user` Portal.

```bash
go run ./cmd/migrate_users_to_portal --portal-db pos_portal_db --apply --i-have-backup
```

Jika tidak ada user `admin` di Toko, akun seed Portal `admin` (password `admin123`) tetap ada:
**login dan ganti passwordnya segera**.

## 5. Kunci aplikasi, FE, dan nginx

Urutan bagian ini: 5.1 build FE, 5.2 nginx, 5.3 kunci SSO (butuh Portal sudah terbuka).

### 5.1 Build FE

```bash
cd /opt/pos-mahenz/PORTAL/FE && npm ci && npm run type-check && npm run build
cd /opt/pos-mahenz/TOKO/FE   && npm ci && npm run type-check && npm run build

sudo mkdir -p /var/www/pos-web/portal /var/www/pos-web/apps/toko
sudo rm -rf /var/www/pos-web/portal/* /var/www/pos-web/apps/toko/*
sudo cp -r /opt/pos-mahenz/PORTAL/FE/dist/* /var/www/pos-web/portal/
sudo cp -r /opt/pos-mahenz/TOKO/FE/dist/*   /var/www/pos-web/apps/toko/
sudo chown -R www-data:www-data /var/www/pos-web
```

`.env.production` Toko memakai `VITE_BASE_PATH=/toko/`, `VITE_API_URL=/toko/api`, `VITE_PORTAL_URL=/`;
Portal memakai `VITE_API_URL=/portal-api`. **Jangan** build Toko dengan env lain untuk server ini.

### 5.2 Nginx

```bash
sudo cp /etc/nginx/sites-available/pos-web ~/pos-web.nginx.backup     # cadangan konfigurasi lama
sudo cp /opt/pos-mahenz/nginx/pos-web.conf /etc/nginx/sites-available/pos-web
# skrip maintenance versi baru (mendukung per aplikasi)
sudo cp /opt/pos-mahenz/nginx/maintenance/maintenance-on.sh  /usr/local/bin/
sudo cp /opt/pos-mahenz/nginx/maintenance/maintenance-off.sh /usr/local/bin/
sudo chmod +x /usr/local/bin/maintenance-*.sh

sudo nginx -t
sudo systemctl reload nginx
sudo certbot --nginx          # memasang kembali blok HTTPS (sertifikat lama dipakai ulang)
sudo nginx -t && sudo systemctl reload nginx
```

Bila `nginx -t` gagal, **jangan reload**; kembalikan cadangan (`~/pos-web.nginx.backup`).

### 5.3 Kunci SSO Toko

Saat ini maintenance masih aktif. Buka URL bypass yang dicetak `maintenance-on.sh` (sekali) agar bisa
membuka Portal, lalu:

1. Login ke Portal sebagai admin.
2. Menu **Kelola Aplikasi** -> aplikasi **POS Toko** -> **Buat kunci baru** -> salin kunci (tampil sekali).
3. Pastikan alamat aplikasi POS Toko adalah `/toko/` (menu Ubah).
4. Pasang di database Toko:
   ```sql
   INSERT INTO portal_credentials (id, app_id, app_key) VALUES (1, 'toko', '<KUNCI_YANG_DISALIN>')
   ON DUPLICATE KEY UPDATE app_id = VALUES(app_id), app_key = VALUES(app_key);
   ```
   (tidak perlu restart; dibaca tiap login lewat Portal)

## 6. Verifikasi (dengan URL bypass maintenance)

Dengan cookie bypass dari 5.3, cek:

- [ ] `https://mahenz.duckdns.org/` menampilkan halaman login Portal; login admin berhasil.
- [ ] Halaman Aplikasi menampilkan kartu POS Toko; klik membuka `/toko/` dan masuk tanpa login ulang.
- [ ] Akun Toko lama (mis. kasir) login langsung di `/toko/login` masih berhasil (login langsung masih menyala).
- [ ] Kasir yang sebelumnya login tetap login (migrasi key `auth-session` -> `toko.auth-session`); keranjang belum selesai masih ada.
- [ ] Kasir baru dari Portal masuk ke halaman "Menunggu pengaturan akses" sampai admin Toko mengatur role.
- [ ] Logout dari Toko kembali ke Portal.
- [ ] Cetak struk (browser, dan APK/Bluetooth bila dipakai) berfungsi di `/toko/`.
- [ ] Tablet APK: dibuka ulang, tampil Portal, lalu masuk ke Toko dan cetak Bluetooth tetap bekerja (alamat printer terjaga).
- [ ] `sudo journalctl -u pos-toko-backend -u pos-portal-backend -n 50` tanpa error.

Bila semua oke:

```bash
sudo maintenance-off.sh
```

## 7. Rollback (bila ada masalah serius)

```bash
sudo maintenance-on.sh all mahenz.duckdns.org
sudo systemctl stop pos-toko-backend pos-portal-backend
sudo systemctl disable pos-toko-backend pos-portal-backend

# Nginx lama
sudo cp ~/pos-web.nginx.backup /etc/nginx/sites-available/pos-web
sudo nginx -t && sudo systemctl reload nginx

# Kode lama
cd /opt && sudo mv pos-mahenz pos-mahenz_gagal && sudo mv pos-mahenz_<TANGGAL> pos-mahenz

sudo systemctl enable --now pos-backend       # BE lama di port 8080
sudo maintenance-off.sh
```

FE lama di `/var/www/pos-web/dist` tidak disentuh, jadi langsung berfungsi lagi. Database Toko tidak perlu
dipulihkan: migrasi 011 aditif. Pulihkan dari backup hanya bila data rusak.

## 8. Setelah stabil (beberapa hari)

- Matikan **login langsung** di Toko lewat Manajemen User (khusus admin) bila semua sudah memakai Portal.
  Bila Portal mati saat login langsung mati, nyalakan paksa lewat SSH: set `"DirectLoginForceEnabled": true`
  di `TOKO/BE/config/config_prod.json` lalu `sudo systemctl restart pos-toko-backend`.
- Hapus blok `location ^~ /api/` di nginx setelah semua perangkat memakai `/toko/api/`.
- Pindahkan `SecretKey` dan password database dari `config_prod.json` Toko ke environment variable
  (saat ini tersimpan di file repo).
- Hapus akun Portal yang tidak dipakai dan pastikan tidak ada akun berpassword bawaan.
- Setelah yakin, hapus folder lama (`/opt/pos-mahenz_<TANGGAL>`, `/var/www/pos-web/dist`).

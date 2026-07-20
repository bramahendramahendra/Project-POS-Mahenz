# Panduan Instalasi Server Production — POS Mahenz

Dokumen ini menjelaskan langkah-langkah lengkap instalasi Backend (Go) dan Frontend (Vite/React) ke server production, beserta penjelasan **kenapa** setiap langkah dilakukan — supaya bisa dipakai sebagai bahan belajar, bukan sekadar "copy-paste perintah".

> Ditulis berdasarkan kondisi kode saat ini (Juli 2026). Lihat bagian [Catatan Kondisi Kode Saat Ini](#catatan-kondisi-kode-saat-ini) untuk hal-hal yang perlu diperbaiki sebelum benar-benar deploy.

---

## Daftar Isi

1. [Gambaran Arsitektur](#1-gambaran-arsitektur)
2. [Prasyarat Server](#2-prasyarat-server)
3. [Instalasi Dependensi Server](#3-instalasi-dependensi-server)
4. [Setup Database MySQL](#4-setup-database-mysql)
5. [Deploy Backend (Go)](#5-deploy-backend-go)
6. [Menjalankan Backend sebagai Service (systemd)](#6-menjalankan-backend-sebagai-service-systemd)
7. [Deploy Frontend (Vite/React)](#7-deploy-frontend-vitereact)
8. [Konfigurasi Nginx (Reverse Proxy + Static Hosting)](#8-konfigurasi-nginx-reverse-proxy--static-hosting)
9. [HTTPS dengan Let's Encrypt](#9-https-dengan-lets-encrypt)
10. [Checklist Deploy](#10-checklist-deploy)
11. [Update / Redeploy Selanjutnya](#11-update--redeploy-selanjutnya)
12. [Troubleshooting](#12-troubleshooting)
13. [Catatan Kondisi Kode Saat Ini](#catatan-kondisi-kode-saat-ini)

---

## 1. Gambaran Arsitektur

```
                         ┌─────────────────────────┐
  Browser  ── HTTPS ──▶  │  Nginx (port 80/443)     │
                         │  - Serve FE static (dist)│
                         │  - Proxy /api → BE :8080 │
                         └───────────┬──────────────┘
                                     │
                                     ▼
                         ┌─────────────────────────┐
                         │  Backend Go (pos_api)    │
                         │  systemd service :8080   │
                         └───────────┬──────────────┘
                                     │
                                     ▼
                         ┌─────────────────────────┐
                         │  MySQL 8.0               │
                         │  database: pos_retail_db │
                         └─────────────────────────┘
```

**Kenapa arsitektur seperti ini?**

- Backend adalah **satu binary Go** hasil compile (`pos_api`), tidak butuh runtime Node/PHP di server — cukup jalankan langsung. Ini salah satu keunggulan Go untuk deployment: tidak ada dependency runtime.
- Frontend adalah **Vite SPA** (Single Page Application) — bukan Next.js dengan SSR. Artinya hasil `npm run build` hanya berupa file statis (HTML/CSS/JS) di folder `dist/`. Tidak perlu Node.js berjalan terus-menerus di server, cukup di-serve oleh web server statis seperti Nginx.
- Nginx berperan ganda: (1) menyajikan file statis FE, dan (2) sebagai **reverse proxy** meneruskan request `/api/*` ke backend Go yang berjalan di port 8080. Ini menghindari masalah CORS karena dari sisi browser, FE dan API terlihat berasal dari domain yang sama.

---

## 2. Prasyarat Server

Server production (disarankan Ubuntu 22.04 LTS atau sejenis) perlu:

| Software | Versi Minimal | Kegunaan |
|---|---|---|
| Go | Versi terbaru (≥ 1.24.5) | Compile backend (hanya dibutuhkan saat build, tidak wajib di server jika build dilakukan di tempat lain / CI) |
| MySQL | 8.0 | Database utama |
| Node.js | 18+ | Build frontend (hanya saat build) |
| Nginx | 1.18+ | Web server / reverse proxy |
| Git | terbaru | Menarik source code |
| Certbot (opsional) | terbaru | Sertifikat HTTPS gratis dari Let's Encrypt |

**Catatan penting:** Go dan Node.js sebenarnya hanya dibutuhkan untuk proses **build**. Kalau Anda build binary/artifact di mesin lain (misalnya laptop atau CI server) lalu meng-upload hasilnya, server production tidak wajib punya Go/Node terinstall. Tapi untuk pemula, lebih mudah build langsung di server yang sama supaya tidak ada perbedaan arsitektur (misal ARM vs x86).

---

## 3. Instalasi Dependensi Server

Contoh untuk Ubuntu/Debian:

```bash
# Update package list
sudo apt update && sudo apt upgrade -y

# Install Nginx
sudo apt install -y nginx
systemctl status nginx

# Install MySQL Server
sudo apt install -y mysql-server
systemctl status mysql
sudo mysql_secure_installation   # ikuti wizard: set root password, hapus anonymous user, dll

Pertanyaan | Jawaban
VALIDATE PASSWORD COMPONENT? | Y (opsional, tapi disarankan untuk cek kekuatan password)
Kalau muncul pilihan level password policy | Pilih 2 (STRONG) kalau ditanya, atau 0/1 kalau mau lebih fleksibel — pilih 1 (MEDIUM) kalau ragu
Set root password? / Change the root password? | Y — masukkan password baru yang kuat untuk database
Remove anonymous users? | Y
Disallow root login remotely? | Y
Remove test database and access to it? | Y
Reload privilege tables now? | Y

# Install Git
sudo apt install -y git

# Install Go versi terbaru (ikuti panduan resmi https://go.dev/learn/)
# Ambil nama file tarball terbaru secara otomatis dari go.dev, lalu unduh & pasang
GO_TARBALL=$(curl -s https://go.dev/VERSION?m=text | head -n1)
wget https://go.dev/dl/${GO_TARBALL}.linux-amd64.tar.gz
sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf ${GO_TARBALL}.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
go version   # pastikan versi ≥ go1.24.5 (versi minimal yang diminta go.mod project ini)

# Install Node.js 20 LTS (via NodeSource)
curl -fsSL https://deb.nodesource.com/setup_20.x | sudo -E bash -
sudo apt install -y nodejs
node -v   # pastikan v20.x
```

---

## 4. Setup Database MySQL

Backend menggunakan MySQL 8.0 dengan **migrasi otomatis** (dijalankan sendiri oleh aplikasi saat start, tidak perlu tool migrasi eksternal seperti `goose`/`migrate`). Anda hanya perlu menyiapkan database kosong dan user MySQL.

```bash
sudo mysql -u root -p
```

```sql
CREATE DATABASE pos_retail_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- Buat user khusus aplikasi (jangan pakai root untuk aplikasi!)
CREATE USER 'pos_user'@'localhost' IDENTIFIED BY 'PASSWORD_KUAT_DISINI';
GRANT ALL PRIVILEGES ON pos_retail_db.* TO 'pos_user'@'localhost';
FLUSH PRIVILEGES;
EXIT;
```

**Kenapa buat user terpisah, bukan pakai `root`?** Prinsip *least privilege* — kalau kredensial aplikasi bocor, penyerang hanya bisa mengakses database `pos_retail_db`, bukan seluruh instance MySQL.

Saat backend pertama kali dijalankan, ia otomatis membaca semua file di `BE/database/migrations/` secara berurutan (`001_init_schema.sql` → `006_sync_id_map.sql` saat ini) dan mencatat progres di tabel `migrations_history`. Anda **tidak perlu** menjalankan file SQL secara manual.

---

## 5. Deploy Backend (Go)

### 5.1 Clone / upload source code

```bash
sudo mkdir -p /opt/pos-mahenz
sudo chown -R $USER:$USER /opt/pos-mahenz
git clone <URL_REPO_ANDA> /opt/pos-mahenz
cd /opt/pos-mahenz/BE
```

### 5.2 Siapkan file konfigurasi

Ada dua lapis konfigurasi backend:

1. **`.env`** — menentukan mode aplikasi & port dasar.
2. **`config/config_prod.json`** — konfigurasi detail (database, JWT, CORS, logging) untuk mode `prod`.

```bash
cp .env.example .env   # jika ada; jika tidak, buat manual seperti di bawah
```

Isi `BE/.env`:

```env
GIN_MODE=release
APP_NAME=POS Retail API
APP_AUTHOR=MAHENZ
APP_VERSION=1.0.0
APP_HOST=https://api.domain-anda.com/
APP_PORT=8080

RELEASE_MODE=prod
```

**Kenapa `GIN_MODE=release`?** Mode `debug` (default) mencetak log verbose setiap request dan menampilkan stack trace — cocok untuk development tapi boros resource dan berpotensi membocorkan detail internal di production. `RELEASE_MODE=prod` adalah variabel custom aplikasi ini yang menentukan file config JSON mana yang dipakai (`config_prod.json`).

Edit `BE/config/config_prod.json`, sesuaikan bagian `Database` dan `CorsAllowOrigins`:

```json
{
  "Database": {
    "Type": "mysql",
    "Host": "127.0.0.1",
    "Port": "3306",
    "User": "pos_user",
    "Password": "PASSWORD_KUAT_DISINI",
    "Database": "pos_retail_db",
    "MaxOpenConns": 50,
    "MaxIdleConns": 10,
    "MaxLifetime": 300,
    "MaxIdleTime": 600
  },
  "CorsAllowOrigins": [
    "https://pos.domain-anda.com"
  ]
}
```

`CorsAllowOrigins` **wajib** diisi dengan domain frontend production Anda yang sebenarnya (bukan `example.com`) — kalau tidak, browser akan memblokir request FE ke API karena kebijakan CORS.

### 5.3 Set `SECRETKEY` (JWT) via environment variable — WAJIB

Perhatikan bahwa `SecretKey` di `config_prod.json` **sengaja dikosongkan** di dalam repo. Ini bukan bug, melainkan langkah keamanan yang disengaja: kode aplikasi akan `panic` (menolak start) jika `SecretKey` kosong, memaksa Anda menyuplainya lewat environment variable saat deploy — supaya secret JWT **tidak pernah tersimpan di git**.

```bash
export SECRETKEY="ganti-dengan-string-acak-panjang-dan-rahasia"
```

Generate string acak yang aman:

```bash
openssl rand -base64 48
```

> Simpan `SECRETKEY` ini di tempat aman. Jika berubah, semua token JWT yang sudah terbit (sesi login user) akan langsung invalid.

### 5.4 Install dependency & build binary

```bash
cd /opt/pos-mahenz/BE
go mod tidy
go build -o pos_api main.go
```

Hasilnya: satu file binary `pos_api` yang bisa langsung dieksekusi, sudah berisi seluruh dependency Go ter-compile (statically linked). Tidak perlu `go` terinstall lagi untuk *menjalankan* binary ini setelah build selesai — hanya dibutuhkan saat build.

### 5.5 Test jalankan manual dulu

```bash
./pos_api
```

Pastikan muncul log server listening di port 8080 tanpa error, dan tabel-tabel ter-migrasi otomatis (cek dengan `SHOW TABLES;` di MySQL). Setelah yakin jalan, hentikan dengan `Ctrl+C` — selanjutnya kita jalankan lewat systemd (langkah berikutnya) supaya otomatis restart jika crash dan otomatis start saat server reboot.

---

## 6. Menjalankan Backend sebagai Service (systemd)

Menjalankan binary langsung di terminal (`./pos_api &`) tidak cukup untuk production — begitu SSH terputus atau server reboot, aplikasi akan mati. Solusinya: bungkus sebagai **systemd service**.

Buat file `/etc/systemd/system/pos-backend.service`:

```bash
sudo nano /etc/systemd/system/pos-backend.service
```

```ini
[Unit]
Description=POS Retail Backend API
After=network.target mysql.service

[Service]
Type=simple
User=www-data
Group=www-data
WorkingDirectory=/opt/pos-mahenz/BE
Environment="SECRETKEY=ganti-dengan-string-acak-panjang-dan-rahasia"
ExecStart=/opt/pos-mahenz/BE/pos_api
Restart=on-failure
RestartSec=5
StandardOutput=append:/var/log/pos-backend/stdout.log
StandardError=append:/var/log/pos-backend/stderr.log

[Install]
WantedBy=multi-user.target
```

**Penjelasan tiap bagian:**
- `After=network.target mysql.service` → pastikan MySQL sudah aktif dulu sebelum backend dicoba start.
- `User=www-data` → jangan jalankan sebagai `root`; prinsip least privilege lagi.
- `Environment="SECRETKEY=..."` → cara systemd menyuntikkan env var yang dibutuhkan langkah 5.3. (Alternatif lebih aman: gunakan `EnvironmentFile=/opt/pos-mahenz/BE/.env.secret` yang permission-nya dibatasi `600`, supaya secret tidak tampil di `systemctl status` atau `ps aux`.)
- `Restart=on-failure` → kalau aplikasi crash, systemd otomatis restart.

Siapkan folder log & set kepemilikan:

```bash
sudo mkdir -p /var/log/pos-backend
sudo chown www-data:www-data /var/log/pos-backend /opt/pos-mahenz/BE -R
```

Aktifkan service:

```bash
sudo systemctl daemon-reload
sudo systemctl enable pos-backend
sudo systemctl start pos-backend
sudo systemctl status pos-backend   # pastikan "active (running)"
```

Cek log jika ada masalah:

```bash
journalctl -u pos-backend -f
```

---

## 7. Deploy Frontend (Vite/React)

Frontend project ini bernama `web-v2`, dibangun dengan Vite + React 19 + React Router v7 (bukan Next.js — jadi tidak ada server-side rendering, murni SPA statis).

### 7.1 Siapkan environment variable production

```bash
cd /opt/pos-mahenz/FE   # atau lokasi clone FE Anda
```

Edit `FE/.env.production`:

```env
VITE_API_URL=https://api.domain-anda.com/api
VITE_APP_NAME=POS System
VITE_PLATFORM=web
```

**Kenapa perlu di-set sebelum build, bukan saat runtime?** Vite meng-*inline* semua variabel `VITE_*` ke dalam bundle JavaScript **saat proses build**, bukan dibaca saat runtime seperti aplikasi backend biasa. Artinya jika Anda ganti `VITE_API_URL` nanti, Anda **wajib build ulang** — tidak cukup edit file `.env` di server lalu restart, karena tidak ada proses yang "restart" untuk file statis.

### 7.2 Build

```bash
npm install
npm run type-check   # pastikan 0 TypeScript error
npm run lint          # pastikan 0 ESLint error
npm run build         # output ke FE/dist/
```

Hasil build ada di `FE/dist/` — kumpulan file HTML/CSS/JS statis siap disajikan Nginx.

### 7.3 Copy ke folder yang akan di-serve Nginx

```bash
sudo mkdir -p /var/www/pos-web
sudo cp -r dist/* /var/www/pos-web/dist/ 2>/dev/null || sudo cp -r dist /var/www/pos-web/
sudo chown -R www-data:www-data /var/www/pos-web
```

(Sesuaikan path persis dengan yang dipakai di `nginx.conf`, lihat langkah 8.)

---

## 8. Konfigurasi Nginx (Reverse Proxy + Static Hosting)

Project sudah menyediakan template di `FE/nginx.conf`. Salin dan sesuaikan:

```bash
sudo cp /opt/pos-mahenz/FE/nginx.conf /etc/nginx/sites-available/pos-web
sudo nano /etc/nginx/sites-available/pos-web   # ganti server_name sesuai domain
```

Isi konfigurasi (sudah ada di repo, `FE/nginx.conf`):

```nginx
server {
    listen 80;
    server_name pos.domain-anda.com;
    root /var/www/pos-web/dist;
    index index.html;

    gzip on;
    gzip_types text/plain text/css application/json application/javascript text/xml application/xml application/xml+rss text/javascript;

    location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg|woff|woff2)$ {
        expires 1y;
        add_header Cache-Control "public, immutable";
    }

    location / {
        try_files $uri $uri/ /index.html;
    }

    location /api {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
}
```

**Penjelasan poin-poin penting:**
- `try_files $uri $uri/ /index.html;` → ini **wajib** untuk SPA berbasis React Router. Tanpa baris ini, refresh browser di URL seperti `/dashboard` akan menghasilkan **404** karena Nginx mencari file fisik `dashboard` yang tidak ada — padahal routing `/dashboard` hanya dikenali oleh JavaScript React Router di sisi client. Baris ini memberitahu Nginx: "kalau file/folder tidak ditemukan, kembalikan `index.html` saja", lalu React Router yang menangani routing di browser.
- `expires 1y` aman dipakai karena Vite menambahkan **content hash** ke nama file build (misal `index-a1b2c3.js`) — jika isi file berubah, nama filenya juga berubah otomatis, sehingga cache lama tidak akan pernah menyajikan konten usang.
- `location /api { proxy_pass http://localhost:8080; }` → inilah yang menyatukan FE dan BE dari sudut pandang browser: request ke `https://pos.domain-anda.com/api/...` diteruskan Nginx ke backend Go di `localhost:8080`. Karena itu `VITE_API_URL` di langkah 7.1 diarahkan ke path relatif domain yang sama (atau ke subdomain `api.` jika Anda memisahkan domain FE/BE — sesuaikan setup mana yang dipakai).

Aktifkan dan reload:

```bash
sudo ln -s /etc/nginx/sites-available/pos-web /etc/nginx/sites-enabled/pos-web
sudo nginx -t          # test syntax config sebelum reload
sudo systemctl reload nginx
```

---

## 9. HTTPS dengan Let's Encrypt

Production **wajib** HTTPS — terutama karena aplikasi ini mengirim kredensial login dan JWT token.

```bash
sudo apt install -y certbot python3-certbot-nginx
sudo certbot --nginx -d pos.domain-anda.com -d api.domain-anda.com
```

Certbot otomatis mengedit config Nginx untuk redirect HTTP→HTTPS dan mengatur renewal otomatis (cek dengan `sudo certbot renew --dry-run`).

Setelah HTTPS aktif, update:
- `FE/.env.production` → `VITE_API_URL=https://api.domain-anda.com/api`, lalu **build ulang FE**.
- `BE/config/config_prod.json` → `CorsAllowOrigins` pakai `https://` bukan `http://`.

---

## 10. Checklist Deploy

**Backend:**
- [ ] `SECRETKEY` sudah di-set (bukan string kosong) via systemd env
- [ ] `config_prod.json` → `Database` sudah diisi kredensial production yang benar
- [ ] `config_prod.json` → `CorsAllowOrigins` sudah diisi domain FE production yang benar
- [ ] `.env` → `RELEASE_MODE=prod` dan `GIN_MODE=release`
- [ ] `go build` sukses tanpa error, binary bisa dijalankan manual dulu untuk verifikasi migrasi DB berjalan
- [ ] Service systemd `pos-backend` aktif dan `enable` (auto-start saat reboot)
- [ ] User MySQL aplikasi bukan `root`, password kuat

**Frontend:**
- [ ] `.env.production` → `VITE_API_URL` mengarah ke domain API production yang benar
- [ ] `npm run type-check` dan `npm run lint` → 0 error
- [ ] `npm run build` sukses, folder `dist/` ter-generate
- [ ] Nginx `try_files ... /index.html` sudah ada (cegah 404 saat refresh route)
- [ ] Test manual: buka semua halaman utama, refresh di URL non-root (misal `/dashboard`) tidak 404

**Infrastruktur:**
- [ ] HTTPS aktif (Let's Encrypt) di kedua domain FE & API
- [ ] Firewall server hanya membuka port 80/443 (dan 22 untuk SSH) ke publik — port 8080 (backend) dan 3306 (MySQL) **tidak** perlu terbuka ke internet, cukup diakses `localhost`
- [ ] Backup database terjadwal (lihat folder `BE/backups/` sebagai referensi mekanisme backup yang sudah ada di kode)

---

## 11. Update / Redeploy Selanjutnya

**Backend** (ada perubahan kode):

```bash
cd /opt/pos-mahenz/BE
git pull
go build -o pos_api main.go
sudo systemctl restart pos-backend
journalctl -u pos-backend -f   # pastikan start normal & migrasi baru (jika ada) sukses
```

Migrasi database baru (file SQL baru di `database/migrations/`) akan otomatis dijalankan saat service restart — tidak perlu langkah manual tambahan, asalkan nomor urut file migrasi lebih besar dari yang terakhir tercatat di tabel `migrations_history`.

**Frontend** (ada perubahan kode):

```bash
cd /opt/pos-mahenz/FE
git pull
npm install
npm run build
sudo rm -rf /var/www/pos-web/dist
sudo cp -r dist /var/www/pos-web/
sudo chown -R www-data:www-data /var/www/pos-web
```

Tidak perlu reload Nginx untuk update FE (Nginx hanya membaca file dari disk setiap request), kecuali Anda juga mengubah `nginx.conf`.

---

## 12. Troubleshooting

| Gejala | Kemungkinan Penyebab | Solusi |
|---|---|---|
| Backend gagal start, log `panic: SecretKey is empty` | Env var `SECRETKEY` belum di-set | Cek `systemctl cat pos-backend`, pastikan baris `Environment="SECRETKEY=..."` ada dan tidak kosong |
| Backend gagal start, error koneksi database | Kredensial/host MySQL salah di `config_prod.json`, atau MySQL belum jalan | `sudo systemctl status mysql`; cek `Database.Host/User/Password` |
| FE menampilkan halaman tapi API call gagal (CORS error di console browser) | `CorsAllowOrigins` di `config_prod.json` tidak cocok dengan domain FE | Tambahkan domain FE yang benar (termasuk `https://`), restart backend |
| Refresh di `/dashboard` (atau route lain) muncul 404 dari Nginx | `try_files` belum ada di config Nginx | Tambahkan `try_files $uri $uri/ /index.html;`, `nginx -t` lalu reload |
| Perubahan `VITE_API_URL` tidak berpengaruh setelah edit `.env.production` | Lupa build ulang — Vite meng-inline env var saat build, bukan runtime | Jalankan `npm run build` lagi lalu copy ulang `dist/` |
| Service backend restart terus-menerus (crash loop) | Cek `journalctl -u pos-backend -f` untuk stack trace asli | Biasanya error koneksi DB atau file migrasi SQL yang gagal dieksekusi |

---

## Catatan Kondisi Kode Saat Ini

Beberapa hal di repo yang perlu Anda ketahui/perbaiki sebelum mengandalkan file existing untuk deploy:

1. **`BE/deploy/Dockerfile` dan `BE/deploy/docker-compose.yml` sudah usang (stale)** — dokumen ini sengaja tidak memakai Docker karena:
   - Dockerfile mereferensikan `./cmd/main.go`, padahal struktur project saat ini `main.go` ada langsung di root `BE/`.
   - Dockerfile meng-copy satu `config.json`, padahal konfigurasi sekarang terpisah per environment (`config_dev.json` / `config_prod.json`) di folder `config/`.
   - Base image `golang:1.21-alpine` lebih lama dari `go 1.24.5` yang diminta `go.mod`.
   - `docker-compose.yml` memakai nama env var (`DB_HOST`, `JWT_SECRET`, dst) yang **tidak dibaca** oleh `config.go` (yang sebenarnya pakai `SECRETKEY` dan file JSON, bukan env var per-field database).
   - **Jika suatu saat ingin containerize deployment**, kedua file ini perlu ditulis ulang menyesuaikan struktur project sekarang — bisa saya bantu jika diperlukan.
2. **Tidak ada CI/CD** (tidak ditemukan `.github/workflows/` atau sejenisnya) — deployment saat ini sepenuhnya manual sesuai panduan di atas. Jika ke depan volume deploy makin sering, pertimbangkan setup GitHub Actions untuk build + deploy otomatis.
3. **Redis ada di `go.mod`** dan ada implementasinya di `BE/pkg/redis/`, tapi menurut `BE/config/CONFIG_MIGRATION_NOTES.md`, fitur cache berbasis Redis **saat ini tidak dipakai** (field Redis sudah dihapus dari config JSON). Jadi server production **tidak perlu** menginstall Redis untuk saat ini.
4. **Tidak ada Dockerfile untuk frontend** — deployment FE didokumentasikan sebagai build lokal lalu copy `dist/` ke server (sesuai `FE/DEPLOYMENT.md`), bukan containerized. Panduan di atas mengikuti pendekatan ini.

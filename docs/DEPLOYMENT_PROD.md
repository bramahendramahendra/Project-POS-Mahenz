# Panduan Instalasi Server Production — POS Mahenz

Dokumen ini menjelaskan langkah-langkah lengkap instalasi Backend (Go) dan Frontend (Vite/React) ke server production, beserta penjelasan **kenapa** setiap langkah dilakukan — supaya bisa dipakai sebagai bahan belajar, bukan sekadar "copy-paste perintah".

> Ditulis berdasarkan kondisi kode saat ini (Juli 2026). Lihat bagian [Catatan Kondisi Kode Saat Ini](#catatan-kondisi-kode-saat-ini) untuk hal-hal yang perlu diperbaiki sebelum benar-benar deploy.

---

## Daftar Isi

1. [Gambaran Arsitektur](#1-gambaran-arsitektur)
2. [Prasyarat Server](#2-prasyarat-server)
3. [Instalasi Dependensi Server](#3-instalasi-dependensi-server)
4. [Setup Database MySQL](#4-setup-database-mysql)
    - 4.1 [Akses Database dari Navicat (Remote)](#41-akses-database-dari-navicat-remote)
5. [Deploy Backend (Go)](#5-deploy-backend-go)
6. [Menjalankan Backend sebagai Service (systemd)](#6-menjalankan-backend-sebagai-service-systemd)
7. [Deploy Frontend (Vite/React)](#7-deploy-frontend-vitereact)
8. [Konfigurasi Nginx (Reverse Proxy + Static Hosting)](#8-konfigurasi-nginx-reverse-proxy--static-hosting)
9. [HTTPS dengan Let's Encrypt](#9-https-dengan-lets-encrypt)
10. [Checklist Deploy](#10-checklist-deploy)
11. [Update / Redeploy Selanjutnya](#11-update--redeploy-selanjutnya)
12. [Maintenance Database](#12-maintenance-database)
13. [Maintenance Mode (Aplikasi)](#13-maintenance-mode-aplikasi)
14. [Troubleshooting](#14-troubleshooting)
15. [Full Backup & Redeploy (FE, BE, DB)](#15-full-backup--redeploy-fe-be-db)
16. [Catatan Kondisi Kode Saat Ini](#catatan-kondisi-kode-saat-ini)

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

### 4.1 Akses Database dari Navicat (Remote)

Secara default MySQL di server hanya bind ke `127.0.0.1` (lihat checklist §10: port 3306 **tidak** dibuka ke publik) — ini benar untuk keamanan production, tapi berarti Navicat di laptop Anda tidak bisa connect langsung ke `IP_SERVER:3306`. Ada dua cara mengaksesnya, pilih salah satu:

**Opsi A — SSH Tunnel (disarankan, tidak perlu ubah firewall/MySQL sama sekali)**

Cara ini paling aman karena port 3306 tetap tertutup ke internet; koneksi dienkripsi lewat SSH yang memang sudah terbuka.

1. Di Navicat, buat koneksi baru **MySQL** → tab **General** isi seperti biasa (Host `127.0.0.1`, Port `3306`, User `pos_user`, Password sesuai `config_prod.json`).
2. Buka tab **SSH**, centang **Use SSH Tunnel**, lalu isi:
   ```
   Host: <IP_SERVER_ANDA>
   Port: 22
   User: <user SSH Anda>
   Authentication Method: Password atau Private Key (sesuaikan dengan cara Anda login SSH ke server)
   ```
3. Klik **Test Connection** — Navicat akan connect ke MySQL lewat "terowongan" SSH tersebut, seolah-olah MySQL ada di `localhost` Anda sendiri.
4. Save & Connect.

**Opsi B — Buka MySQL untuk remote access (kurang disarankan)**

Hanya lakukan ini kalau Opsi A tidak memungkinkan, dan **wajib** batasi akses hanya dari IP Anda (jangan `0.0.0.0/0`) kecuali IP Anda memang dinamis (lihat varian "tanpa IP" di bawah).

1. Ubah bind address MySQL agar menerima koneksi dari luar `localhost`:
   ```bash
   sudo nano /etc/mysql/mysql.conf.d/mysqld.cnf
   ```
   Cari baris `bind-address = 127.0.0.1`, ubah jadi:
   ```
   bind-address = 0.0.0.0
   ```
   Lalu restart: `sudo systemctl restart mysql`
2. Buat/izinkan user MySQL untuk connect dari host tertentu, misalnya dari IP publik laptop/kantor Anda `203.0.113.10` (bukan `'localhost'`):
   ```sql
   CREATE USER 'pos_user'@'203.0.113.10' IDENTIFIED BY 'PASSWORD_KUAT_DISINI';
   GRANT ALL PRIVILEGES ON pos_retail_db.* TO 'pos_user'@'203.0.113.10';
   FLUSH PRIVILEGES;
   ```
3. Buka port 3306 di firewall **hanya** untuk IP Anda:
   ```bash
   sudo ufw allow from 203.0.113.10 to any port 3306
   ```

  Cek port :
   ```bash
   sudo ss -tulnp | grep 3306
   ```
4. Di Navicat, buat koneksi MySQL biasa (tanpa tab SSH): Host `IP_SERVER_ANDA`, Port `3306`, User `pos_user`, Password sesuai.

**Kalau IP Anda tidak tetap (dinamis, misal ganti-ganti WiFi/tethering) sehingga tidak bisa dipatok ke satu IP:**

Anda tetap bisa buka sementara tanpa membatasi IP tertentu — tapi ini **membuka MySQL ke seluruh internet**, jadi hanya untuk keperluan cepat dan **wajib** segera ditutup lagi (lihat langkah "Menutup kembali" di bawah).

```sql
-- User yang boleh connect dari IP mana pun (tanpa batasan IP)
CREATE USER 'pos_user'@'%' IDENTIFIED BY 'PASSWORD_KUAT_DISINI';
GRANT ALL PRIVILEGES ON pos_retail_db.* TO 'pos_user'@'%';
FLUSH PRIVILEGES;
```

```bash
# Buka port 3306 untuk semua IP (bukan hanya IP tertentu)
sudo ufw allow 3306
```

Cek port :
```bash
sudo ss -tulnp | grep 3306
```

Di Navicat: Host `IP_SERVER_ANDA`, Port `3306`, User `pos_user`, Password sesuai — bisa connect dari IP mana pun selama aturan ini masih aktif.

> ⚠️ Selama aturan ini aktif, siapa pun di internet yang tahu IP server + password bisa mencoba connect ke MySQL Anda. Pastikan password `pos_user` kuat, dan jangan dibiarkan menyala lebih dari sesi kerja Anda saat itu.

**Menutup kembali setelah selesai (penting kalau niatnya cuma buka sebentar):**

Sesuaikan perintah dengan varian yang tadi dipakai — dengan IP tertentu atau tanpa IP.

```bash
# 1. Tutup port 3306 di firewall lagi
sudo ufw status numbered          # lihat nomor aturan 3306 yang aktif
sudo ufw delete allow from 203.0.113.10 to any port 3306   # kalau tadi pakai IP tertentu
# ATAU, kalau tadi pakai varian tanpa IP:
sudo ufw delete allow 3306
sudo ufw status   # pastikan aturan 3306 sudah hilang

# 2. Kembalikan bind-address MySQL ke localhost saja
sudo nano /etc/mysql/mysql.conf.d/mysqld.cnf
# ubah balik: bind-address = 127.0.0.1
sudo systemctl restart mysql

# 3. (opsional tapi disarankan) hapus user remote yang tadi dibuat
mysql -u root -p
```
```sql
DROP USER 'pos_user'@'203.0.113.10';   -- atau: DROP USER 'pos_user'@'%';
FLUSH PRIVILEGES;
```

Jika tidak pakai IP 

```sql
DROP USER 'pos_user'@'%';
FLUSH PRIVILEGES;
```


Verifikasi sudah tertutup dari luar (jalankan dari komputer Anda, bukan dari server):
```bash
telnet IP_SERVER_ANDA 3306   # atau: nc -zv IP_SERVER_ANDA 3306
# harus gagal connect / connection refused / timeout
```

> Urutan penting: tutup firewall dulu (langkah 1) baru ubah `bind-address` (langkah 2) — supaya tidak ada jendela waktu di mana MySQL sudah bind ke `0.0.0.0` tapi firewall belum sempat diaktifkan ulang (untuk kasus ini urutannya tidak terlalu kritis karena keduanya dilakukan berurutan cepat, tapi kebiasaan ini penting kalau prosesnya dipisah/didelegasikan).

**Kenapa Opsi A lebih disarankan?** Membuka 3306 ke internet (meski dibatasi IP) menambah *attack surface* — kalau IP Anda dinamis (berubah-ubah, misal WiFi rumah/kafe), aturan firewall harus terus diperbarui atau malah dilonggarkan jadi rentan. SSH tunnel tidak butuh perubahan firewall/MySQL sama sekali dan tetap aman walau IP Anda berubah, karena yang dibuka publik hanya port 22 (SSH) yang memang sudah ada.

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

### Prasyarat: Wajib Punya Domain (Tidak Bisa Pakai IP)

**Let's Encrypt tidak menerbitkan sertifikat untuk alamat IP telanjang** (misal `139.180.214.187`), hanya untuk domain (FQDN — Fully Qualified Domain Name). Certbot akan gagal kalau dicoba langsung ke IP. Alasannya: proses verifikasi kepemilikan Let's Encrypt (disebut *ACME challenge*) bekerja lewat DNS/domain — IP address tidak punya mekanisme pembuktian kepemilikan seperti itu.

Jadi kalau server Anda **masih diakses via IP** (belum ada domain), HTTPS lewat Let's Encrypt **belum bisa dilakukan** — lewati dulu langkah ini, lanjut pakai HTTP untuk sementara, dan kembali ke sini setelah domain siap.

### Langkah Sebelum Bisa Menjalankan Certbot

1. **Beli domain** (contoh registrar: Niagahoster, Domainesia, Namecheap, Cloudflare Registrar)
2. **Arahkan domain ke IP server** — buat **A record** di pengaturan DNS domain tersebut:
   ```
   Type: A
   Name: pos (atau @ untuk root domain)
   Value: 139.180.214.187   (IP server Anda)
   TTL: default / 3600
   ```
3. **Tunggu propagasi DNS** — biasanya beberapa menit, bisa sampai 24 jam tergantung registrar
4. **Verifikasi domain sudah resolve ke IP server** sebelum lanjut:
   ```bash
   ping pos.domain-anda.com
   # pastikan IP yang muncul = IP server Anda
   ```
5. **Update `server_name` di Nginx** dari IP ke domain baru (lihat §8), lalu `nginx -t` dan `systemctl reload nginx`, dan pastikan situs masih bisa diakses via `http://pos.domain-anda.com` sebelum lanjut ke Certbot

### Menjalankan Certbot

Setelah domain terverifikasi mengarah ke server, baru jalankan:

```bash
sudo apt install -y certbot python3-certbot-nginx
sudo certbot --nginx -d pos.domain-anda.com -d api.domain-anda.com
```

Certbot otomatis mengedit config Nginx untuk redirect HTTP→HTTPS dan mengatur renewal otomatis (cek dengan `sudo certbot renew --dry-run`).

### Setelah HTTPS Aktif

Update konfigurasi berikut supaya konsisten pakai `https://`, lalu build ulang FE:

- `FE/.env.production` → `VITE_API_URL=https://api.domain-anda.com/api` (atau tetap `/api` kalau pakai path relatif, lihat §7.1) — kalau diubah, **wajib `npm run build` ulang** karena Vite meng-inline env var saat build, bukan runtime.
- `BE/config/config_prod.json` → `CorsAllowOrigins` pakai `https://` bukan `http://`, lalu `sudo systemctl restart pos-backend`.

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

## 12. Maintenance Database

### Backup Database

Selalu backup dulu sebelum melakukan operasi destruktif (drop tabel, reset database, dsb). Ada dua cara:

**Opsi A — Lewat fitur bawaan aplikasi (disarankan untuk pemakaian rutin)**

Backend sudah punya endpoint backup bawaan (`BE/domain/backup/`) yang menjalankan `mysqldump` dan menyimpan hasilnya ke `BE/backups/*.sql` — kemungkinan besar juga sudah ada tombolnya di menu Sistem/Pengaturan pada FE (cek menu terkait "Backup" di aplikasi). Kelebihannya: terproteksi permission (`sistem.backup`, `can_create`/`can_view`), dan filenya bisa langsung di-download lewat aplikasi tanpa perlu SSH ke server.

**Opsi B — Manual lewat `mysqldump` di server (untuk backup di luar aplikasi, mis. sebelum reset total)**

```bash
mkdir -p /opt/pos-mahenz/BE/backups   # folder yang sama dipakai fitur backup bawaan aplikasi
mysqldump -u pos_user -p pos_retail_db > /opt/pos-mahenz/BE/backups/backup_manual_$(date +%Y%m%d_%H%M%S).sql
```

Verifikasi file backup tidak kosong/corrupt:
```bash
ls -lh /opt/pos-mahenz/BE/backups/
tail -n 5 /opt/pos-mahenz/BE/backups/backup_manual_*.sql   # harus diakhiri baris normal SQL, bukan terpotong
```

**Kenapa perlu backup dulu sebelum reset?** Langkah drop tabel/database di bawah ini **destruktif dan permanen** — begitu tabel di-drop dan backend restart (migrasi otomatis membuat ulang skema kosong), seluruh data lama (produk, transaksi, user, dll) tidak bisa dikembalikan kecuali dari file `.sql` yang sudah di-backup.

**Cara restore dari file backup** (kalau suatu saat perlu mengembalikan data):
```bash
mysql -u pos_user -p pos_retail_db < /opt/pos-mahenz/BE/backups/NAMA_FILE_BACKUP.sql
```
> Restore akan menimpa data yang ada saat ini di tabel-tabel yang sama persis dengan yang ada di file backup — pastikan database dalam kondisi yang Anda inginkan sebelum menjalankan ini (biasanya dijalankan tepat setelah drop tabel/database di bawah, pada database yang masih kosong).

---

### Drop Semua Tabel (Reset Skema, User & Database Tetap Ada)

Dipakai kalau Anda ingin mengosongkan seluruh skema database (misal sebelum re-migrasi dari awal saat ada perubahan besar) **tanpa** menghapus database atau user MySQL-nya.

> ⚠️ **Destruktif dan permanen.** Semua data (produk, transaksi, user, dll) hilang tanpa bisa dikembalikan kecuali ada backup. Jangan jalankan di production kecuali benar-benar sengaja reset total.

```bash
mysql -u pos_user -p pos_retail_db
```

```sql
SET FOREIGN_KEY_CHECKS = 0;

SET @tables = NULL;
SELECT GROUP_CONCAT('`', table_name, '`') INTO @tables
FROM information_schema.tables
WHERE table_schema = 'pos_retail_db';

SET @tables = CONCAT('DROP TABLE IF EXISTS ', @tables);
PREPARE stmt FROM @tables;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET FOREIGN_KEY_CHECKS = 1;
```

**Penjelasan:**
- `SET FOREIGN_KEY_CHECKS = 0` — mematikan sementara pengecekan foreign key, supaya tabel bisa di-drop dalam urutan apa pun tanpa error "cannot drop table referenced by foreign key"
- Query `GROUP_CONCAT` mengumpulkan semua nama tabel di database `pos_retail_db` menjadi satu string
- `PREPARE` + `EXECUTE` menjalankan `DROP TABLE IF EXISTS tabel1, tabel2, ...` sekaligus untuk seluruh tabel yang ditemukan
- `FOREIGN_KEY_CHECKS` dinyalakan lagi di akhir

Verifikasi database benar-benar kosong:
```sql
SHOW TABLES;   -- harus muncul: Empty set
EXIT;
```

Setelah tabel kosong, restart backend supaya migrasi otomatis membuat ulang semua tabel dari `database/migrations/`:
```bash
sudo systemctl restart pos-backend
sudo systemctl status pos-backend
sudo tail -n 30 /var/log/pos-backend/stdout.log   # pastikan migrasi jalan sukses
```

### Alternatif: Drop & Buat Ulang Database (Kalau Ingin Reset User/Privilege Juga)

Kalau Anda juga ingin database dibuat ulang dari nol (bukan cuma tabel), lihat §4 [Setup Database MySQL](#4-setup-database-mysql) — gunakan `DROP DATABASE IF EXISTS pos_retail_db;` lalu `CREATE DATABASE` lagi. Cara ini juga menghapus tabel `migrations_history`, jadi semua migrasi ikut berjalan ulang dari awal.

### Ringkasan Langkah: Reset Total Aplikasi dari Awal

Urutan lengkap kalau Anda ingin menjalankan ulang aplikasi seolah baru pertama kali install (misal untuk demo ulang, ganti data uji, atau ada perubahan skema besar):

1. **(Opsional tapi sangat disarankan)** Aktifkan [maintenance mode](#13-maintenance-mode-aplikasi) dulu supaya tidak ada user lain yang sedang transaksi saat proses reset berjalan: `sudo maintenance-on.sh`
2. **Backup data lama** (lihat bagian [Backup Database](#backup-database) di atas) — jangan lewati langkah ini kalau ada kemungkinan data lama masih dibutuhkan.
3. **Drop tabel/database** — pilih salah satu:
   - Drop semua tabel saja (user & privilege MySQL tetap) — lihat langkah "Drop Semua Tabel" di atas, atau
   - Drop & buat ulang database dari nol — lihat "Alternatif" di atas.
4. **Restart backend** supaya migrasi otomatis membuat ulang seluruh skema kosong dari `BE/database/migrations/`:
   ```bash
   sudo systemctl restart pos-backend
   sudo tail -n 30 /var/log/pos-backend/stdout.log   # pastikan semua file migrasi jalan sukses tanpa error
   ```
5. Verifikasi tabel-tabel dasar sudah terbentuk kembali (`SHOW TABLES;`), lalu login ke aplikasi untuk memastikan alur dasar (login, buat data awal, dst) berjalan normal.
6. Matikan maintenance mode: `sudo maintenance-off.sh`

---

## 13. Maintenance Mode (Aplikasi)

Dipakai saat butuh menutup akses ke seluruh aplikasi sementara (deploy backend, migrasi database besar, dsb) tanpa mematikan Nginx atau service lain. Mekanismenya berbasis **flag file**: Nginx mengecek keberadaan file tersebut di setiap request, kalau ada maka semua request (termasuk `/api`) langsung dibalas `503` dan diarahkan ke halaman statis `maintenance.html`.

Ada juga **mode bypass** berbasis cookie + token, supaya Anda (developer/tester) tetap bisa membuka aplikasi secara normal untuk mengetes selagi maintenance aktif untuk user lain.

File terkait ada di `FE/maintenance/`:

| File | Fungsi |
|---|---|
| `maintenance.html` | Halaman statis yang ditampilkan ke user saat maintenance aktif |
| `maintenance-on.sh` | Mengaktifkan mode maintenance (`touch` flag file) + generate token bypass baru |
| `maintenance-off.sh` | Menonaktifkan mode maintenance (`rm` flag file) + rotasi token (invalidate cookie bypass lama) |
| `maintenance_token.conf.example` | Template file token bypass, di-copy sekali ke `/etc/nginx/maintenance_token.conf` saat setup awal |

### Cara kerja bypass

1. Setiap kali `maintenance-on.sh` dijalankan, script generate **token acak baru** (`openssl rand -hex 16`), menulisnya ke `/etc/nginx/maintenance_token.conf`, lalu reload Nginx (graceful, tanpa downtime) dan mencetak URL bypass lengkap dengan token tersebut.
2. Anda buka URL itu **sekali** di browser, mis. `https://pos.domain-anda.com/maintenance-bypass?token=abcd1234...` — Nginx mencocokkan token di URL dengan token aktif, lalu men-set cookie `mnt_bypass` (masa berlaku 24 jam) dan redirect ke `/`.
3. Selama cookie itu masih ada & masih cocok dengan token aktif, **semua request Anda selanjutnya** (halaman apa pun, tanpa perlu menambahkan token lagi di URL) dilewatkan dari blok `503` — sementara user lain yang tidak punya cookie ini tetap melihat halaman maintenance.
4. Saat `maintenance-off.sh` dijalankan, token **otomatis dirotasi ulang** — jadi cookie bypass lama (kalau masih tersisa di browser Anda) langsung basi. Kalau minggu depan Anda `maintenance-on.sh` lagi, Anda wajib buka ulang URL bypass dengan token yang baru dicetak saat itu.
5. Sebagai lapis pengaman kedua, cookie juga otomatis kedaluwarsa setelah **24 jam** meski Anda lupa menjalankan `maintenance-off.sh`.

### Setup awal (sekali saja di server)

```bash
sudo mkdir -p /var/www/pos-web/maintenance
sudo cp /opt/pos-mahenz/FE/maintenance/maintenance.html /var/www/pos-web/maintenance/
sudo cp /opt/pos-mahenz/FE/maintenance/maintenance-on.sh /opt/pos-mahenz/FE/maintenance/maintenance-off.sh /usr/local/bin/
sudo chmod +x /usr/local/bin/maintenance-on.sh /usr/local/bin/maintenance-off.sh

# Siapkan file token bypass awal (kosong, akan diisi otomatis saat maintenance-on.sh pertama kali dijalankan)
sudo cp /opt/pos-mahenz/FE/maintenance/maintenance_token.conf.example /etc/nginx/maintenance_token.conf
```

Pastikan `nginx.conf` yang aktif di server (`/etc/nginx/sites-available/...`) sudah memuat konfigurasi `maintenance.flag`, `include /etc/nginx/maintenance_token.conf;`, dan `location = /maintenance-bypass` seperti di `FE/nginx.conf` pada repo ini (lihat [bagian 8](#8-konfigurasi-nginx-reverse-proxy--static-hosting)).

### Mengaktifkan maintenance

```bash
# argumen opsional: domain Anda, dipakai untuk mencetak URL bypass yang siap pakai
sudo maintenance-on.sh pos.domain-anda.com
```

Semua request ke domain akan langsung mendapat halaman maintenance (503). Script juga mencetak URL bypass — buka sekali di browser Anda untuk mulai testing normal:

```
Maintenance mode: ON
Flag file: /var/www/pos-web/maintenance.flag

URL bypass (buka SEKALI di browser Anda untuk mulai testing, cookie berlaku 24 jam):
https://pos.domain-anda.com/maintenance-bypass?token=3f9a7c2e1b8d4560a1c2e3f4a5b6c7d8
```

### Menonaktifkan maintenance

```bash
sudo maintenance-off.sh
```

Aplikasi langsung kembali normal begitu flag file dihapus — tidak ada delay/cache karena Nginx mengecek keberadaan file ini di setiap request baru. Token bypass juga otomatis dirotasi, jadi cookie lama tidak berlaku lagi untuk maintenance berikutnya.

### Ganti token secara manual (opsional)

Kalau suatu saat ingin mengganti token tanpa menunggu siklus on/off (misal token bocor), edit langsung:

```bash
echo 'set $maintenance_token "TOKEN_BARU_ANDA";' | sudo tee /etc/nginx/maintenance_token.conf
sudo nginx -t
sudo systemctl reload nginx
```

**Catatan:**
- Halaman `maintenance.html` di server (`/var/www/pos-web/maintenance/`) **terpisah** dari folder `dist/` FE — supaya tetap bisa diakses walau proses build FE sedang berjalan/gagal.
- Kalau desain `maintenance.html` diubah di repo, perlu `cp` ulang manual ke server (tidak otomatis ikut proses redeploy FE di [bagian 11](#11-update--redeploy-selanjutnya)).
- File `/etc/nginx/maintenance_token.conf` sengaja **di luar** repo git (mirip `SECRETKEY` di §5.3) — supaya token aktif tidak pernah tersimpan/ter-commit ke source control.
- Mekanisme ini menutup **seluruh** aplikasi (FE + API) untuk yang tidak punya cookie bypass. Untuk mengunci sebagian fitur saja tanpa menutup akses total, perlu pendekatan berbeda (middleware di level backend) — belum diimplementasikan di project ini.

---

## 14. Troubleshooting

| Gejala | Kemungkinan Penyebab | Solusi |
|---|---|---|
| Backend gagal start, log `panic: SecretKey is empty` | Env var `SECRETKEY` belum di-set | Cek `systemctl cat pos-backend`, pastikan baris `Environment="SECRETKEY=..."` ada dan tidak kosong |
| Backend gagal start, error koneksi database | Kredensial/host MySQL salah di `config_prod.json`, atau MySQL belum jalan | `sudo systemctl status mysql`; cek `Database.Host/User/Password` |
| FE menampilkan halaman tapi API call gagal (CORS error di console browser) | `CorsAllowOrigins` di `config_prod.json` tidak cocok dengan domain FE | Tambahkan domain FE yang benar (termasuk `https://`), restart backend |
| Refresh di `/dashboard` (atau route lain) muncul 404 dari Nginx | `try_files` belum ada di config Nginx | Tambahkan `try_files $uri $uri/ /index.html;`, `nginx -t` lalu reload |
| Perubahan `VITE_API_URL` tidak berpengaruh setelah edit `.env.production` | Lupa build ulang — Vite meng-inline env var saat build, bukan runtime | Jalankan `npm run build` lagi lalu copy ulang `dist/` |
| Service backend restart terus-menerus (crash loop) | Cek `journalctl -u pos-backend -f` untuk stack trace asli | Biasanya error koneksi DB atau file migrasi SQL yang gagal dieksekusi |
| `nginx -t` gagal setelah setup maintenance dengan pesan file `/etc/nginx/maintenance_token.conf` tidak ditemukan | Belum copy file token awal | Jalankan langkah "Setup awal" §13: `sudo cp .../maintenance_token.conf.example /etc/nginx/maintenance_token.conf` |
| URL `/maintenance-bypass?token=...` membalas `403` | Token di URL tidak cocok dengan token aktif (sudah kedaluwarsa/dirotasi, atau salah copy-paste) | Jalankan ulang `maintenance-on.sh` untuk generate token baru, pakai URL yang baru dicetak |
| Sudah buka URL bypass tapi tetap kena halaman maintenance | Cookie `mnt_bypass` belum ter-set (browser blokir cookie / third-party cookie disabled) atau sudah lewat 24 jam | Cek di DevTools → Application → Cookies apakah `mnt_bypass` ada; kalau tidak, ulangi buka URL bypass |

---

## 15. Full Backup & Redeploy (FE, BE, DB)

Dipakai kalau Anda ingin melakukan **redeploy total** (FE, BE, dan database sekaligus) dengan strategi *rename-lalu-clone-ulang*: versi lama (kode + database) tidak dihapus, hanya di-*rename*/duplicate jadi arsip bertanggal, lalu semuanya di-*clone*/dibuat ulang dari nol dengan nama folder & database yang sama seperti semula. Konsepnya mirip snapshot manual — kalau ada yang salah setelah redeploy, versi lama masih utuh di server dan bisa dipakai untuk rollback.

**Kapan dipakai:** situasi yang lebih berat dari update biasa di [§11](#11-update--redeploy-selanjutnya) — misal ada perubahan besar di skema/kode yang ingin dites dari kondisi benar-benar bersih, tapi Anda tetap mau punya jejak/arsip penuh dari versi sebelumnya (kode maupun data) tanpa harus rely sepenuhnya ke git history atau file `.sql` backup saja.

> ⚠️ Proses ini menimbulkan downtime (BE mati sesaat, DB di-drop & dibuat ulang kosong). **Wajib** aktifkan [maintenance mode](#13-maintenance-mode-aplikasi) dulu.

### Ringkasan strategi

| Komponen | Yang lama diapakan | Yang baru |
|---|---|---|
| Database `pos_retail_db` | Di-duplicate ke database baru bertanggal (mis. `pos_retail_db_20260727`), lalu `pos_retail_db` asli di-**drop & dibuat ulang kosong** | Migrasi otomatis jalan lagi dari `BE/database/migrations/` saat backend restart |
| Folder `BE/` | Di-**rename** jadi `BE_20260727` (arsip, dibiarkan ada di server) | `git clone` ulang ke `BE/` (nama sama seperti semula) |
| Folder `FE/` (source) | Di-**rename** jadi `FE_20260727` (arsip) | `git clone` ulang ke `FE/` (nama sama seperti semula) |
| `/var/www/pos-web/dist` (hasil build FE yang disajikan Nginx) | Di-**rename** jadi `dist_20260727` (arsip) | Hasil `npm run build` baru di-copy ke `dist/` (nama sama seperti semula) |

Karena semua folder baru memakai **nama identik** dengan sebelumnya (`BE`, `FE`, `dist`), Anda **tidak perlu** mengubah `WorkingDirectory` di systemd (`pos-backend.service`) maupun `root` di `nginx.conf` — keduanya tetap menunjuk ke path yang sama seperti biasa.

### Langkah-langkah

```bash
# Variabel tanggal, dipakai konsisten di semua langkah di bawah
TODAY=$(date +%Y%m%d)
echo "Tanggal arsip: $TODAY"
```

**1. Aktifkan maintenance mode**

```bash
sudo maintenance-on.sh pos.domain-anda.com
```

**2. Duplicate database ke nama baru bertanggal**

```bash
sudo mysql -u root -p
```

```sql
CREATE DATABASE pos_retail_db_20260727 CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
EXIT;
```

```bash
# Dump dari database lama, langsung import ke database baru (tanpa file perantara)
mysqldump -u pos_user -p pos_retail_db | mysql -u pos_user -p pos_retail_db_20260727
```

Verifikasi jumlah tabel di database baru sama dengan yang lama:
```sql
-- jalankan di masing-masing database untuk dibandingkan
SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'pos_retail_db';
SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'pos_retail_db_20260727';
```

**Kenapa dump+import langsung lewat pipe, bukan lewat file `.sql` dulu?** Lebih cepat dan tidak perlu ruang disk ekstra untuk file perantara — cocok untuk database yang belum terlalu besar. Kalau database sudah besar (ratusan MB+) dan ingin ada file `.sql` fisik sebagai arsip tambahan juga, bisa gabungkan dengan pendekatan di [Backup Database](#backup-database) (§12) terlebih dahulu, baru `mysql ... < file.sql` ke database baru.

**3. Reset database aktif (`pos_retail_db`) jadi kosong**

Ikuti langkah **Drop Semua Tabel** di [§12](#12-maintenance-database) — bukan `DROP DATABASE`, karena user `pos_user` dan privilege-nya ingin tetap dipakai apa adanya untuk deploy baru.

**4. Rename folder BE & FE lama, clone ulang yang baru**

```bash
cd /opt/pos-mahenz

# Hentikan service backend dulu sebelum folder BE dipindah, supaya tidak ada proses yang masih mengunci file
sudo systemctl stop pos-backend

sudo mv BE BE_20260727
sudo mv FE FE_20260727

git clone <URL_REPO_ANDA> /tmp/pos-mahenz-fresh
sudo mv /tmp/pos-mahenz-fresh/BE /opt/pos-mahenz/BE
sudo mv /tmp/pos-mahenz-fresh/FE /opt/pos-mahenz/FE
rm -rf /tmp/pos-mahenz-fresh
```

**5. Setup ulang BE seperti deploy pertama kali**

Folder `BE/` baru hasil clone masih kosong dari file konfigurasi (`.env`, `config_prod.json` tidak ikut ter-commit ke git). Salin dari arsip lama supaya tidak perlu isi ulang dari nol, lalu build:

```bash
cd /opt/pos-mahenz/BE
cp /opt/pos-mahenz/BE_20260727/.env .env
cp /opt/pos-mahenz/BE_20260727/config/config_prod.json config/config_prod.json

go mod tidy
go build -o pos_api main.go
```

> Cek ulang `config_prod.json` → `Database.Database` tetap `pos_retail_db` (bukan nama bertanggal) — aplikasi tetap terhubung ke database aktif yang sudah dikosongkan di langkah 3.

Jalankan ulang service (`WorkingDirectory` di systemd tidak berubah, karena path `BE/` namanya sama seperti sebelumnya):

```bash
sudo chown -R www-data:www-data /opt/pos-mahenz/BE
sudo systemctl start pos-backend
sudo systemctl status pos-backend
sudo tail -n 30 /var/log/pos-backend/stdout.log   # pastikan migrasi jalan sukses membuat ulang skema kosong
```

**6. Setup ulang FE seperti deploy pertama kali**

```bash
cd /opt/pos-mahenz/FE
cp /opt/pos-mahenz/FE_20260727/.env.production .env.production

npm install
npm run type-check
npm run lint
npm run build

sudo mv /var/www/pos-web/dist /var/www/pos-web/dist_20260727
sudo cp -r dist /var/www/pos-web/
sudo chown -R www-data:www-data /var/www/pos-web
```

Tidak perlu ubah/reload Nginx — `root /var/www/pos-web/dist` di `nginx.conf` tetap menunjuk ke path yang sama, isinya saja yang baru.

**7. Verifikasi & matikan maintenance mode**

- Buka URL bypass maintenance ([§13](#13-maintenance-mode-aplikasi)) untuk mengetes aplikasi normal dulu sebelum dibuka ke semua user: login, cek halaman-halaman utama, pastikan API jalan.
- Kalau semua sudah oke: `sudo maintenance-off.sh`

### Membersihkan arsip lama (opsional, manual)

Folder & database bertanggal (`BE_20260727`, `FE_20260727`, `dist_20260727`, `pos_retail_db_20260727`) **sengaja dibiarkan** di server sebagai rollback point, tidak dihapus otomatis oleh langkah-langkah di atas. Kalau sudah yakin tidak diperlukan lagi (biasanya setelah beberapa hari/minggu aplikasi baru terbukti stabil), hapus manual:

```bash
sudo rm -rf /opt/pos-mahenz/BE_20260727 /opt/pos-mahenz/FE_20260727
sudo rm -rf /var/www/pos-web/dist_20260727
```

```sql
DROP DATABASE pos_retail_db_20260727;
```

> Pantau juga penggunaan disk (`df -h`) kalau proses redeploy total ini dilakukan berkala — arsip yang menumpuk tanpa pernah dibersihkan lama-lama bisa memenuhi storage server, terutama untuk database yang besar.

### Rollback (kalau redeploy baru ternyata bermasalah)

```bash
sudo maintenance-on.sh

sudo systemctl stop pos-backend
sudo mv /opt/pos-mahenz/BE /opt/pos-mahenz/BE_gagal_20260727
sudo mv /opt/pos-mahenz/BE_20260727 /opt/pos-mahenz/BE
sudo mv /opt/pos-mahenz/FE /opt/pos-mahenz/FE_gagal_20260727
sudo mv /opt/pos-mahenz/FE_20260727 /opt/pos-mahenz/FE

sudo mv /var/www/pos-web/dist /var/www/pos-web/dist_gagal_20260727
sudo mv /var/www/pos-web/dist_20260727 /var/www/pos-web/dist

# Kembalikan database: drop yang baru, restore dari arsip
mysql -u root -p -e "DROP DATABASE pos_retail_db;"
mysql -u root -p -e "CREATE DATABASE pos_retail_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
mysqldump -u pos_user -p pos_retail_db_20260727 | mysql -u pos_user -p pos_retail_db

sudo systemctl start pos-backend
sudo maintenance-off.sh
```

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

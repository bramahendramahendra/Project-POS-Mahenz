# Panduan Deploy POS Mahenz (Versi Simpel)

Panduan ini ditulis untuk orang yang **tidak biasa koding**. Tinggal ikuti langkah dari
atas ke bawah, salin-tempel perintahnya ke terminal server. Setiap perintah diberi
keterangan singkat "untuk apa", bukan teori panjang.

> Butuh detail lebih dalam? Ada dua dokumen lama yang lebih lengkap:
> `DEPLOYMENT_PROD.md` (instalasi lengkap) dan `DEPLOYMENT_REDEPLOY_FULL.md` (redeploy total).
> Untuk pemakaian sehari-hari, cukup file ini saja.

---

## Isi Panduan

- [Bagian 1 — Istilah Penting (baca sekali)](#bagian-1--istilah-penting-baca-sekali)
- [Bagian 2 — Pasang Pertama Kali (server baru)](#bagian-2--pasang-pertama-kali-server-baru)
- [Bagian 3 — Update Aplikasi (rutin, paling sering dipakai)](#bagian-3--update-aplikasi-rutin-paling-sering-dipakai)
- [Bagian 4 — Backup Database](#bagian-4--backup-database)
- [Bagian 5 — Kalau Ada Masalah](#bagian-5--kalau-ada-masalah)

---

## Bagian 1 — Istilah Penting (baca sekali)

Biar tidak bingung, ini 4 kata yang sering muncul:

| Istilah | Artinya sederhana |
|---|---|
| **Backend (BE)** | "Mesin" aplikasi. Program Go yang mengolah data & terhubung ke database. |
| **Frontend (FE)** | "Tampilan" aplikasi. Yang dilihat user di browser. |
| **Database** | Tempat semua data disimpan (produk, transaksi, user). Namanya `pos_retail_db`. |
| **Nginx** | "Pintu masuk". Mengarahkan pengunjung ke tampilan (FE) atau ke mesin (BE). |

Lokasi standar di server (jangan diubah):
- Kode aplikasi: `/opt/pos-mahenz` (di dalamnya ada folder `BE` dan `FE`)
- Tampilan yang tampil ke user: `/var/www/pos-web/dist`

---

## Bagian 2 — Pasang Pertama Kali (server baru)

> Lakukan ini **hanya sekali** saat server benar-benar baru. Kalau aplikasi sudah pernah
> jalan dan Anda cuma mau update, **lompat ke [Bagian 3](#bagian-3--update-aplikasi-rutin-paling-sering-dipakai)**.

### 2.1 Pasang program yang dibutuhkan

Salin-tempel semua ini (untuk Ubuntu). Ini memasang web server, database, dan alat build.

```bash
sudo apt update && sudo apt upgrade -y
sudo apt install -y nginx mysql-server git
```

Pasang Go dan Node.js (dipakai untuk membangun aplikasi):

```bash
# Go
GO_TARBALL=$(curl -s https://go.dev/VERSION?m=text | head -n1)
wget https://go.dev/dl/${GO_TARBALL}.linux-amd64.tar.gz
sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf ${GO_TARBALL}.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc && source ~/.bashrc

# Node.js 20
curl -fsSL https://deb.nodesource.com/setup_20.x | sudo -E bash -
sudo apt install -y nodejs
```

### 2.2 Siapkan database kosong

Masuk ke MySQL:

```bash
sudo mysql -u root -p
```

Lalu salin-tempel ini (ganti `PASSWORD_KUAT_DISINI` dengan password buatan Anda sendiri):

```sql
CREATE DATABASE pos_retail_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER 'pos_user'@'localhost' IDENTIFIED BY 'PASSWORD_KUAT_DISINI';
GRANT ALL PRIVILEGES ON pos_retail_db.* TO 'pos_user'@'localhost';
FLUSH PRIVILEGES;
EXIT;
```

> Tabel-tabelnya **dibuat otomatis** oleh aplikasi nanti. Anda tidak perlu bikin tabel manual.

### 2.3 Ambil kode aplikasi

```bash
sudo mkdir -p /opt/pos-mahenz
sudo git clone <URL_REPO_ANDA> /opt/pos-mahenz
sudo chown -R $USER:$USER /opt/pos-mahenz
```

### 2.4 Atur & nyalakan Backend (mesin)

Buat file pengaturan `BE/.env`:

```bash
cd /opt/pos-mahenz/BE
nano .env
```

Isi dengan (simpan dengan Ctrl+O lalu Ctrl+X):

```env
GIN_MODE=release
APP_NAME=POS Retail API
APP_PORT=8080
RELEASE_MODE=prod
```

Isi pengaturan database di `config/config_prod.json` — bagian `Database` diisi password
yang tadi Anda buat, dan `CorsAllowOrigins` diisi alamat website Anda:

```bash
nano config/config_prod.json
```

```json
{
  "Database": {
    "Host": "127.0.0.1",
    "Port": "3306",
    "User": "pos_user",
    "Password": "PASSWORD_KUAT_DISINI",
    "Database": "pos_retail_db"
  },
  "CorsAllowOrigins": ["https://pos.domain-anda.com"]
}
```

Buat kunci rahasia (untuk keamanan login). Jalankan, lalu **simpan hasilnya**:

```bash
openssl rand -base64 48
```

Bangun & jalankan mesinnya sebagai layanan otomatis:

```bash
go mod tidy
go build -o pos_api main.go
```

Buat layanan systemd (biar aplikasi auto-jalan & auto-restart):

```bash
sudo nano /etc/systemd/system/pos-backend.service
```

Isi (ganti bagian `SECRETKEY` dengan hasil `openssl` tadi):

```ini
[Unit]
Description=POS Retail Backend API
After=network.target mysql.service

[Service]
Type=simple
User=www-data
Group=www-data
WorkingDirectory=/opt/pos-mahenz/BE
Environment="SECRETKEY=tempel-hasil-openssl-disini"
ExecStart=/opt/pos-mahenz/BE/pos_api
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
```

Nyalakan:

```bash
sudo chown -R www-data:www-data /opt/pos-mahenz/BE
sudo systemctl daemon-reload
sudo systemctl enable pos-backend
sudo systemctl start pos-backend
sudo systemctl status pos-backend   # harus muncul "active (running)"
```

### 2.5 Bangun & pasang Frontend (tampilan)

```bash
cd /opt/pos-mahenz/FE
nano .env.production
```

Isi (arahkan ke alamat API Anda):

```env
VITE_API_URL=https://api.domain-anda.com/api
VITE_APP_NAME=POS System
```

Bangun tampilannya lalu salin ke folder yang disajikan Nginx:

```bash
npm install
npm run build
sudo mkdir -p /var/www/pos-web
sudo cp -r dist /var/www/pos-web/
sudo chown -R www-data:www-data /var/www/pos-web
```

### 2.6 Atur Nginx (pintu masuk)

```bash
sudo cp /opt/pos-mahenz/FE/nginx.conf /etc/nginx/sites-available/pos-web
sudo nano /etc/nginx/sites-available/pos-web   # ganti server_name jadi domain Anda
sudo ln -s /etc/nginx/sites-available/pos-web /etc/nginx/sites-enabled/pos-web
sudo nginx -t                    # cek konfigurasi tidak error
sudo systemctl reload nginx
```

### 2.7 Pasang HTTPS (kunci gembok hijau) — butuh domain

> Hanya bisa kalau sudah punya **domain** (mis. `pos.domain-anda.com`), tidak bisa pakai IP saja.

```bash
sudo apt install -y certbot python3-certbot-nginx
sudo certbot --nginx -d pos.domain-anda.com -d api.domain-anda.com
```

Selesai. Buka website Anda di browser untuk mengecek.

---

## Bagian 3 — Update Aplikasi (rutin, paling sering dipakai)

Ini yang Anda pakai setiap ada perubahan kode baru. **Backup dulu** (lihat
[Bagian 4](#bagian-4--backup-database)) sebelum mulai, biar aman.

### Aktifkan maintenance:
```bash
sudo maintenance-on.sh 139.180.214.187
```

### Update Database:
> Jika ada perubahan pada database
Lakukan Backup databse **Bagian 4 — Backup Database** ikuti instruksi bab tersebut.

### Update Backend (mesin):

```bash
cd /opt/pos-mahenz/BE
git pull
go build -o pos_api main.go
sudo systemctl restart pos-backend
sudo systemctl status pos-backend   # pastikan "active (running)"
```

> Kalau ada perubahan struktur database, itu **jalan otomatis** saat restart. Tidak perlu langkah manual.

### Update Frontend (tampilan):

```bash
cd /opt/pos-mahenz/FE
git pull
npm install
npm run type-check
npm run lint
npm run build
sudo rm -rf /var/www/pos-web/dist
sudo cp -r dist /var/www/pos-web/
sudo chown -R www-data:www-data /var/www/pos-web
```

Selesai. Refresh browser (Ctrl+Shift+R) untuk melihat perubahan.

---

## Bagian 4 — Backup Database

**2. Duplicate database ke nama baru bertanggal**

```bash
sudo mysql -u root -p
```

```sql
CREATE DATABASE pos_retail_db_20260913 CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

SELECT user, host FROM mysql.user WHERE user = 'pos_user';

GRANT ALL PRIVILEGES ON pos_retail_db_20260913.* TO 'pos_user'@'localhost';
GRANT ALL PRIVILEGES ON pos_retail_db_20260913.* TO 'pos_user'@'%';
FLUSH PRIVILEGES;

-- Verifikasi kedua host sudah punya izin ke database baru
SHOW GRANTS FOR 'pos_user'@'localhost';
SHOW GRANTS FOR 'pos_user'@'%';
EXIT;
```


**Selalu backup sebelum update besar.** Perintah ini menyimpan salinan seluruh data ke satu file:

```bash
mkdir -p /opt/pos-mahenz/BE/backups
mysqldump -u pos_user -p pos_retail_db > /opt/pos-mahenz/BE/backups/backup_$(date +%Y%m%d_%H%M%S).sql
```

Cek file backup benar-benar terisi (jangan kosong):

```bash
ls -lh /opt/pos-mahenz/BE/backups/
```

**Kalau perlu mengembalikan data** dari sebuah file backup:

```bash
mysql -u pos_user -p pos_retail_db < /opt/pos-mahenz/BE/backups/NAMA_FILE_BACKUP.sql
```

---

## Bagian 5 — Kalau Ada Masalah

**Aplikasi tidak bisa dibuka / error 500?** Cek log mesin:

```bash
sudo journalctl -u pos-backend -f
```

(Tekan Ctrl+C untuk berhenti melihat log.)

**Habis update tampilan tapi masih terlihat versi lama?**
Refresh keras di browser: tekan **Ctrl+Shift+R**.

**Backend mati / tidak "running"?** Coba nyalakan ulang:

```bash
sudo systemctl restart pos-backend
sudo systemctl status pos-backend
```

**Website 404 saat refresh di halaman selain beranda?**
Pastikan Nginx pakai konfigurasi dari repo (`FE/nginx.conf`), lalu:

```bash
sudo nginx -t && sudo systemctl reload nginx
```

**Perubahan gagal total dan ingin balik ke data lama?**
Kembalikan dari file backup terakhir (lihat [Bagian 4](#bagian-4--backup-database)).

---

> Checklist singkat sebelum bilang "beres":
> 1. `sudo systemctl status pos-backend` → active (running)
> 2. Website kebuka di browser & bisa login
> 3. Coba buka beberapa halaman, tidak ada error

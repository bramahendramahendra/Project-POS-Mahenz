# Full Backup & Redeploy (FE, BE, DB) — POS Mahenz

Dokumen terpisah dari [DEPLOYMENT_PROD.md](DEPLOYMENT_PROD.md) khusus untuk skenario **redeploy total**: FE, BE, dan database sekaligus, dengan strategi *rename-lalu-clone-ulang* — versi lama tidak dihapus, hanya di-*rename*/duplicate jadi arsip bertanggal, lalu semuanya di-*clone*/dibuat ulang dari nol dengan nama folder & database yang sama seperti semula.

**Kapan dipakai:** situasi yang lebih berat dari update biasa di [DEPLOYMENT_PROD.md §11](DEPLOYMENT_PROD.md#11-update--redeploy-selanjutnya) — misal ada perubahan besar di skema/kode yang ingin dites dari kondisi benar-benar bersih, tapi Anda tetap mau punya jejak/arsip penuh dari versi sebelumnya (kode maupun data) tanpa harus rely sepenuhnya ke git history atau file `.sql` backup saja.

Ada dua mode:

- **Mode A — Rename & Ganti Langsung**: versi lama di-arsipkan (mati), versi baru langsung jadi satu-satunya yang aktif. Lebih simpel, tapi kalau versi baru ternyata error, rollback berarti "matikan dulu, pasang lagi yang lama" (ada downtime tambahan saat rollback).
- **Mode B — Jalan Bersamaan (Side-by-Side)**: versi lama **tetap hidup** di service & port terpisah selama masa testing, versi baru juga hidup di service utama. Anda bisa berpindah-pindah menguji keduanya lewat cookie bypass, baru putuskan: tutup yang lama (kalau baru aman) atau alihkan balik ke lama (kalau baru bermasalah) — keduanya tanpa downtime tambahan karena yang lama tidak pernah benar-benar mati selama proses berlangsung.

> ⚠️ Kedua mode tetap butuh downtime singkat di awal (drop & bikin ulang database aktif) dan **wajib** aktifkan [maintenance mode](DEPLOYMENT_PROD.md#13-maintenance-mode-aplikasi) dulu sebelum mulai.

**Struktur folder yang dipakai di panduan ini:** repo `BE/` dan `FE/` berada dalam **satu folder induk yang sama**, `/opt/pos-mahenz/` (persis seperti hasil `git clone` di [DEPLOYMENT_PROD.md §5.1](DEPLOYMENT_PROD.md#51-clone--upload-source-code)). Karena itu, proses rename-arsip di sini dilakukan **satu kali di level folder induk** (`pos-mahenz` → `pos-mahenz_20260727`), bukan me-rename `BE` dan `FE` terpisah — jadi proses clone ulang juga cukup **satu kali clone penuh**, tidak perlu clone ke folder sementara lalu pindah `BE`/`FE` satu-satu.

---

## Daftar Isi

1. [Mode A — Rename & Ganti Langsung](#mode-a--rename--ganti-langsung)
2. [Mode B — Jalan Bersamaan (Side-by-Side)](#mode-b--jalan-bersamaan-side-by-side)
3. [Membersihkan Arsip Lama](#membersihkan-arsip-lama)

---

## Mode A — Rename & Ganti Langsung

### Ringkasan strategi

| Komponen | Yang lama diapakan | Yang baru |
|---|---|---|
| Database `pos_retail_db` | Di-duplicate ke database baru bertanggal (mis. `pos_retail_db_20260727`), lalu `pos_retail_db` asli di-**drop & dibuat ulang kosong** | Migrasi otomatis jalan lagi dari `BE/database/migrations/` saat backend restart |
| Folder induk `/opt/pos-mahenz` (berisi `BE/` + `FE/`) | Di-**rename** jadi `/opt/pos-mahenz_20260727` (arsip, dibiarkan ada di server) | `git clone` ulang penuh ke `/opt/pos-mahenz` (nama sama seperti semula, isi `BE/` & `FE/` ikut baru) |
| `/var/www/pos-web/dist` (hasil build FE yang disajikan Nginx) | Di-**rename** jadi `dist_20260727` (arsip) | Hasil `npm run build` baru di-copy ke `dist/` (nama sama seperti semula) |

Karena folder baru memakai **nama identik** dengan sebelumnya (`/opt/pos-mahenz`, `dist`), Anda **tidak perlu** mengubah `WorkingDirectory` di systemd (`pos-backend.service`) maupun `root` di `nginx.conf` — keduanya tetap menunjuk ke path yang sama seperti biasa.

### Langkah-langkah

```bash
# Variabel tanggal, dipakai konsisten di semua langkah di bawah
TODAY=$(date +%Y%m%d)
echo "Tanggal arsip: $TODAY"
```

**1. Aktifkan maintenance mode**

```bash
sudo maintenance-on.sh pos.domain-anda.com
sudo maintenance-on.sh 139.180.214.187
```

**2. Duplicate database ke nama baru bertanggal**

```bash
sudo mysql -u root -p
```

```sql
CREATE DATABASE pos_retail_db_20260727 CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- pos_user hanya punya privilege di pos_retail_db (lihat DEPLOYMENT_PROD.md §4) — tanpa GRANT ini,
-- langkah import di bawah akan gagal "Access denied ... to database pos_retail_db_20260727"
GRANT ALL PRIVILEGES ON pos_retail_db_20260727.* TO 'pos_user'@'localhost';
FLUSH PRIVILEGES;
EXIT;
```

```bash
# Dump dari database lama, langsung import ke database baru (tanpa file perantara)
# --no-tablespaces: user pos_user cuma diberi privilege ALL PRIVILEGES per-database,
# bukan privilege PROCESS (global) yang dibutuhkan mysqldump untuk dump info tablespace — tanpa
# flag ini akan gagal dengan error "Access denied ... PROCESS privilege(s)"
#
# Password diambil SEKALI lalu disuplai ke kedua command lewat variabel — JANGAN pakai "-p" polos
# di kedua sisi pipe sekaligus, karena dua prompt password interaktif di terminal yang sama akan
# saling tertukar/kececer (salah satu proses akan menerima password kosong, gagal "using password: NO")
read -s -p "Password pos_user: " DBPASS && echo
mysqldump -u pos_user -p"$DBPASS" --no-tablespaces pos_retail_db | mysql -u pos_user -p"$DBPASS" pos_retail_db_20260727
unset DBPASS
```

> Perhatikan **tidak ada spasi** setelah `-p` (`-p"$DBPASS"`) — kalau ada spasi, MySQL client akan salah membaca `$DBPASS` sebagai nama database, bukan sebagai password.

Verifikasi jumlah tabel di database baru sama dengan yang lama:
```sql
-- jalankan di masing-masing database untuk dibandingkan
SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'pos_retail_db';
SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'pos_retail_db_20260727';
```

**Kenapa dump+import langsung lewat pipe, bukan lewat file `.sql` dulu?** Lebih cepat dan tidak perlu ruang disk ekstra untuk file perantara — cocok untuk database yang belum terlalu besar. Kalau database sudah besar (ratusan MB+) dan ingin ada file `.sql` fisik sebagai arsip tambahan juga, bisa gabungkan dengan pendekatan di [Backup Database](DEPLOYMENT_PROD.md#backup-database) (§12) terlebih dahulu, baru `mysql ... < file.sql` ke database baru.

**3. Reset database aktif (`pos_retail_db`) jadi kosong**

Ikuti langkah **Drop Semua Tabel** di [DEPLOYMENT_PROD.md §12](DEPLOYMENT_PROD.md#12-maintenance-database) — bukan `DROP DATABASE`, karena user `pos_user` dan privilege-nya ingin tetap dipakai apa adanya untuk deploy baru.

### Drop Semua Tabel (Reset Skema, User & Database Tetap Ada)

Dipakai kalau Anda ingin mengosongkan seluruh skema database (misal sebelum re-migrasi dari awal saat ada perubahan besar) **tanpa** menghapus database atau user MySQL-nya.

> ⚠️ **Destruktif dan permanen.** Semua data (produk, transaksi, user, dll) hilang tanpa bisa dikembalikan kecuali ada backup. Jangan jalankan di production kecuali benar-benar sengaja reset total.

```bash
mysql -u pos_user -p
```

```sql
-- Wajib: tanpa USE, PREPARE/EXECUTE di bawah akan gagal "No database selected" karena
-- query DROP TABLE yang dihasilkan tidak menyertakan prefix nama database di depan nama tabelnya
USE pos_retail_db;

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


**4. Rename folder induk lama, clone ulang penuh yang baru**

```bash
# Hentikan service backend dulu sebelum folder dipindah, supaya tidak ada proses yang masih mengunci file
sudo systemctl stop pos-backend

# cek Status 
sudo systemctl status pos-backend

cd /opt
sudo mv pos-mahenz pos-mahenz_20260727

# Satu kali clone penuh — otomatis membawa BE/ dan FE/ sekaligus, tidak perlu pindah folder satu-satu
sudo git clone <URL_REPO_ANDA> /opt/pos-mahenz
sudo chown -R $USER:$USER /opt/pos-mahenz
```

**5. Setup ulang BE seperti deploy pertama kali**

Folder `BE/` baru hasil clone masih kosong dari file konfigurasi (`.env`, `config_prod.json` tidak ikut ter-commit ke git). Salin dari arsip lama supaya tidak perlu isi ulang dari nol, lalu build:

Generate string acak yang aman:

```bash
openssl rand -base64 48
```

```bash
cd /opt/pos-mahenz/BE
sudo cp /opt/pos-mahenz_20260727/BE/.env .env
sudo cp /opt/pos-mahenz_20260727/BE/config/config_prod.json config/config_prod.json

go mod tidy
go build -o pos_api main.go

# Jalankan manual 
./pos_api
```

> Cek ulang `config_prod.json` → `Database.Database` tetap `pos_retail_db` (bukan nama bertanggal) — aplikasi tetap terhubung ke database aktif yang sudah dikosongkan di langkah 3.

Jalankan ulang service (`WorkingDirectory` di systemd tidak berubah, karena path `/opt/pos-mahenz/BE` namanya sama seperti sebelumnya):

```bash
sudo chown -R $USER:$USER /opt/pos-mahenz
sudo systemctl start pos-backend
sudo systemctl status pos-backend
sudo tail -n 30 /var/log/pos-backend/stdout.log   # pastikan migrasi jalan sukses membuat ulang skema kosong
```

**6. Setup ulang FE seperti deploy pertama kali**

```bash
cd /opt/pos-mahenz/FE
sudo cp /opt/pos-mahenz_20260727/FE/.env.production .env.production

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

- Buka URL bypass maintenance ([DEPLOYMENT_PROD.md §13](DEPLOYMENT_PROD.md#13-maintenance-mode-aplikasi)) untuk mengetes aplikasi normal dulu sebelum dibuka ke semua user: login, cek halaman-halaman utama, pastikan API jalan.
- Kalau semua sudah oke: `sudo maintenance-off.sh`

### Rollback Mode A (kalau redeploy baru ternyata bermasalah)

```bash
sudo maintenance-on.sh

sudo systemctl stop pos-backend
sudo mv /opt/pos-mahenz /opt/pos-mahenz_gagal_20260727
sudo mv /opt/pos-mahenz_20260727 /opt/pos-mahenz

sudo mv /var/www/pos-web/dist /var/www/pos-web/dist_gagal_20260727
sudo mv /var/www/pos-web/dist_20260727 /var/www/pos-web/dist

# Kembalikan database: drop yang baru, restore dari arsip
mysql -u root -p -e "DROP DATABASE pos_retail_db;"
mysql -u root -p -e "CREATE DATABASE pos_retail_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
mysqldump -u pos_user -p pos_retail_db_20260727 | mysql -u pos_user -p pos_retail_db

sudo systemctl start pos-backend
sudo maintenance-off.sh
```

**Kelemahan Mode A:** selama proses rollback ini (drop database, restore, restart service), aplikasi kembali kena downtime tambahan — karena versi lama benar-benar sudah "mati" sejak langkah redeploy, bukan sekadar berhenti dipakai. Kalau Anda ingin rollback yang jauh lebih cepat (tanpa restore database dari awal, tanpa restart service), pakai **Mode B** di bawah.

---

## Mode B — Jalan Bersamaan (Side-by-Side)

Ide utamanya: **versi lama tidak pernah benar-benar dimatikan** selama masa testing. Backend lama tetap jalan sebagai service terpisah di port lain, memakai database arsip yang utuh (bukan yang sedang di-reset). Backend baru jalan seperti biasa di port default. Nginx diberi cookie tambahan untuk memilih mau diarahkan ke versi lama atau versi baru — di domain & URL yang sama persis.

### Ringkasan strategi

| Komponen | Versi Lama (`old`) | Versi Baru (default) |
|---|---|---|
| Database | `pos_retail_db_20260727` (arsip, **tidak diubah**, hasil duplicate seperti Mode A langkah 2) | `pos_retail_db` (di-reset kosong, migrasi ulang) |
| Backend | Service **baru** `pos-backend-old`, `WorkingDirectory=/opt/pos-mahenz_20260727/BE`, port **8081** | Service `pos-backend` seperti biasa, `WorkingDirectory=/opt/pos-mahenz/BE`, port 8080 |
| Frontend (hasil build) | `/var/www/pos-web/dist_old` (nama tetap, bukan bertanggal — lihat langkah 7) | `/var/www/pos-web/dist` |
| Cara akses | Cookie `ver_bypass=old` (didapat lewat URL bypass dengan `&version=old`) | Default — tidak perlu cookie tambahan apa pun |

### Persiapan (sama seperti Mode A langkah 1–3)

Ikuti [Mode A langkah 1–3](#langkah-langkah) apa adanya: aktifkan maintenance, duplicate database ke `pos_retail_db_20260727`, lalu reset `pos_retail_db` jadi kosong. Bedanya mulai dari sini: folder induk lama (`/opt/pos-mahenz_20260727`, hasil rename) **tidak dihapus fungsinya**, isinya (`BE/`) **tetap dijalankan** sebagai service kedua.

**4. Rename folder induk lama, clone ulang penuh yang baru** — identik dengan Mode A langkah 4.

**5. Setup BE baru** — identik dengan Mode A langkah 5 (build & start `pos-backend` seperti biasa, terhubung ke `pos_retail_db` yang sudah kosong).

**6. Hidupkan BE lama sebagai service kedua di port 8081**

Folder arsip `/opt/pos-mahenz_20260727/BE` sudah berisi binary lama (`pos_api`) hasil build sebelumnya — tidak perlu build ulang, cukup ubah port & pastikan tetap menunjuk ke database arsipnya sendiri:

```bash
cd /opt/pos-mahenz_20260727/BE

# Ubah APP_PORT di .env supaya tidak bentrok dengan backend baru yang pakai 8080
sudo sed -i 's/^APP_PORT=.*/APP_PORT=8081/' .env

# Pastikan config_prod.json versi lama TETAP mengarah ke database arsip, BUKAN pos_retail_db yang baru direset
grep '"Database"' -A 8 config/config_prod.json   # cek manual: "Database": "pos_retail_db_20260727"
```

> Kalau field `Database.Database` di `config_prod.json` arsip ternyata masih tertulis `pos_retail_db` (bukan nama bertanggal — karena memang itu nama database saat versi lama masih aktif), **wajib** diubah dulu ke `pos_retail_db_20260727` di file arsip ini. Kalau tidak, backend lama akan ikut menulis ke database yang sama dengan backend baru dan mengacaukan hasil testing.

Buat service systemd kedua, `/etc/systemd/system/pos-backend-old.service`:

```bash
sudo nano /etc/systemd/system/pos-backend-old.service
```

```ini
[Unit]
Description=POS Retail Backend API (versi lama, arsip 20260727 - side-by-side testing)
After=network.target mysql.service

[Service]
Type=simple
User=www-data
Group=www-data
WorkingDirectory=/opt/pos-mahenz_20260727/BE
Environment="SECRETKEY=ganti-dengan-string-acak-panjang-dan-rahasia"
ExecStart=/opt/pos-mahenz_20260727/BE/pos_api
Restart=on-failure
RestartSec=5
StandardOutput=append:/var/log/pos-backend-old/stdout.log
StandardError=append:/var/log/pos-backend-old/stderr.log

[Install]
WantedBy=multi-user.target
```

> `SECRETKEY` boleh sama dengan yang dipakai backend baru — tidak masalah dua service pakai secret JWT yang sama, karena keduanya bagian dari aplikasi yang sama, hanya beda versi kode & data.

```bash
sudo mkdir -p /var/log/pos-backend-old
sudo chown -R www-data:www-data /var/log/pos-backend-old /opt/pos-mahenz_20260727/BE

sudo systemctl daemon-reload
sudo systemctl start pos-backend-old
sudo systemctl status pos-backend-old   # pastikan "active (running)" di port 8081
```

> Service ini **sengaja tidak** di-`enable` (auto-start saat reboot) — service ini bersifat sementara selama masa testing, bukan bagian permanen dari infrastruktur.

**7. Setup FE baru, siapkan FE lama dengan nama tetap `dist_old`**

```bash
cd /opt/pos-mahenz/FE
sudo cp /opt/pos-mahenz_20260727/FE/.env.production .env.production

npm install
npm run type-check
npm run lint
npm run build

# Folder lama JANGAN dinamai bertanggal di sini — beri nama tetap "dist_old" supaya cocok
# dengan mapping $spa_root di nginx.conf (lihat catatan di bawah), tidak perlu edit nginx.conf tiap redeploy
sudo mv /var/www/pos-web/dist /var/www/pos-web/dist_old
sudo cp -r dist /var/www/pos-web/
sudo chown -R www-data:www-data /var/www/pos-web
```

### Konfigurasi Nginx untuk routing dua versi (sudah baku di `FE/nginx.conf`)

`FE/nginx.conf` di repo **sudah menyertakan** mapping `$backend_port`/`$spa_root` berbasis cookie `ver_bypass` (map `default` → 8080/`dist`, `old` → 8081/`dist_old`) beserta `location = /maintenance-bypass` yang bisa men-set cookie itu lewat parameter `&version=old`. Jadi **tidak perlu edit `nginx.conf` manual** setiap kali menjalankan Mode B — asalkan:

1. Server sudah memakai `nginx.conf` versi yang sudah menyertakan bagian ini (cek dengan `grep spa_root /etc/nginx/sites-available/pos-web` — kalau kosong, berarti server masih pakai versi lama dan perlu di-`cp` ulang dari repo).
2. Folder FE lama dinamai **tepat** `dist_old` (langkah 7 di atas), dan backend lama dijalankan **tepat** di port `8081` (langkah 6 di atas) — sesuai nilai yang di-hardcode di `map`.

Kalau kedua syarat itu terpenuhi, langsung lanjut cek konfigurasi:
```bash
sudo rm /etc/nginx/sites-enabled/default
sudo nginx -t
sudo systemctl reload nginx
```

### Cara testing kedua versi

1. **Versi baru (default, tidak perlu cookie tambahan)** — buka URL bypass biasa seperti sebelumnya:
   ```
   https://pos.domain-anda.com/maintenance-bypass?token=TOKEN_AKTIF
   ```
2. **Versi lama** — buka URL bypass dengan tambahan `&version=old`:
   ```
   https://pos.domain-anda.com/maintenance-bypass?token=TOKEN_AKTIF&version=old
   ```
   Cookie `ver_bypass=old` ter-set, dan selama cookie itu ada, semua request Anda diarahkan ke `dist_20260727` + backend port 8081 (data lama, kode lama) — tanpa perlu URL berbeda sama sekali.
3. Untuk kembali ke versi baru dari browser yang sama, hapus cookie `ver_bypass` (DevTools → Application → Cookies), atau buka lagi URL bypass **tanpa** `&version=old`.

> **Disarankan pakai dua browser/profile berbeda** (atau satu normal + satu mode incognito) supaya bisa buka versi lama dan versi baru **bersamaan** di layar yang berbeda untuk dibandingkan langsung, tanpa harus bolak-balik hapus cookie.

### Setelah testing: dua kemungkinan keputusan

**A. Versi baru aman → tutup versi lama**

```bash
sudo systemctl stop pos-backend-old
sudo systemctl disable pos-backend-old 2>/dev/null || true   # jaga-jaga kalau sempat ke-enable manual

# Kembalikan nginx.conf ke bentuk semula (hapus blok map, kembalikan location ke root/proxy_pass tetap)
# atau cukup biarkan — map dengan cookie yang sudah tidak pernah di-set lagi otomatis selalu jatuh ke "default"
sudo nginx -t
sudo systemctl reload nginx
```

Lanjut ke [Membersihkan Arsip Lama](#membersihkan-arsip-lama) untuk folder & database arsip.

**B. Versi baru bermasalah → alihkan balik ke versi lama (rollback cepat, tanpa restore database)**

Karena versi lama **masih hidup penuh** dengan datanya sendiri, rollback di sini cuma soal mengubah **default** routing Nginx — jauh lebih cepat dibanding Mode A. Bedanya dengan konfigurasi baku di `FE/nginx.conf`: untuk kasus darurat ini, edit **langsung di file aktif server** (`/etc/nginx/sites-available/pos-web`), sementara — nilai `default`/`old` di dua `map` dibalik urutannya:

```bash
sudo nano /etc/nginx/sites-available/pos-web
```
```nginx
map $cookie_ver_bypass $backend_port {
    default 8081;   # dibalik: default sekarang ke versi LAMA
    new     8080;   # versi baru sekarang jadi opsional, diakses lewat ?version=new
}
map $cookie_ver_bypass $spa_root {
    default /var/www/pos-web/dist_old;
    new     /var/www/pos-web/dist;
}
```

```bash
sudo nginx -t
sudo systemctl reload nginx
sudo maintenance-off.sh
```

Traffic publik otomatis kembali ke versi lama (yang sudah terbukti stabil) begitu maintenance dimatikan, **tanpa perlu drop/restore database sama sekali** — karena `pos_retail_db_20260727` tidak pernah disentuh selama proses ini. Versi baru yang bermasalah tetap ada di `/opt/pos-mahenz/BE` & `dist` untuk diperbaiki lebih lanjut kapan saja, lalu diuji ulang lewat `?version=new` sebelum dicoba jadi default lagi. Setelah versi baru akhirnya diperbaiki dan siap dicoba lagi, **kembalikan dulu** `nginx.conf` server ke isi baku `FE/nginx.conf` di repo (`default`→8080/`dist`, `old`→8081/`dist_old`) sebelum mengulang proses testing.

---

## Membersihkan Arsip Lama

Folder & database arsip (`/opt/pos-mahenz_20260727`, `/var/www/pos-web/dist_old`, `pos_retail_db_20260727`) **sengaja dibiarkan** di server sebagai rollback point, tidak dihapus otomatis oleh langkah-langkah di atas. Kalau sudah yakin tidak diperlukan lagi (biasanya setelah beberapa hari/minggu aplikasi baru terbukti stabil), hapus manual:

```bash
# Kalau sempat pakai Mode B, matikan dulu service lama kalau masih hidup
sudo systemctl stop pos-backend-old 2>/dev/null || true
sudo systemctl disable pos-backend-old 2>/dev/null || true
sudo rm -f /etc/systemd/system/pos-backend-old.service
sudo systemctl daemon-reload

sudo rm -rf /opt/pos-mahenz_20260727
sudo rm -rf /var/www/pos-web/dist_old
sudo rm -rf /var/log/pos-backend-old
```

```sql
DROP DATABASE pos_retail_db_20260727;
```

`map`/routing dua-versi di `nginx.conf` **tidak perlu dibersihkan** — bagian itu sudah baku di `FE/nginx.conf` dan aman dibiarkan permanen: tanpa cookie `ver_bypass=old` dan tanpa folder `dist_old`/service `pos-backend-old`, jalur "old" itu memang tidak pernah tersentuh siapa pun. Kecuali Anda sempat melakukan edit manual darurat (rollback Mode B di atas) yang belum dikembalikan ke baku — cek dulu `grep spa_root /etc/nginx/sites-available/pos-web` sebelum lanjut.

> Pantau juga penggunaan disk (`df -h`) kalau proses redeploy total ini dilakukan berkala — arsip yang menumpuk tanpa pernah dibersihkan lama-lama bisa memenuhi storage server, terutama untuk database yang besar.

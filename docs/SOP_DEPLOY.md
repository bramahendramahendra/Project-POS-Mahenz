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
git pull
sudo chown -R $USER:$USER /opt/sipaduke-testing
```

**5. Setup ulang BE seperti deploy pertama kali**

```bash
cd /opt/pos-mahenz/BE
go mod tidy
go build -o pos_api .
sudo systemctl restart pos-backend
sudo systemctl status pos-backend    # harus "active (running)"
```
**Frontend (tampilan):**

```bash
cd /opt/pos-mahenz/FE
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
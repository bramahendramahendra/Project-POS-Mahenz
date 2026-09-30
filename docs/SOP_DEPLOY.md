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
sudo git clone <URL_REPO_ANDA> /opt/pos-mahenz
sudo chown -R $USER:$USER /opt/pos-mahenz



sudo chown -R $USER:$USER /opt/sipaduke-testing
git clone https://<GITHUB_TOKEN>@github.com/bramahendramahendra/sipaduke.git /opt/sipaduke-testing


# JIka mau di cpy aja
sudo cp -a pos-mahenz pos-mahenz_20260930

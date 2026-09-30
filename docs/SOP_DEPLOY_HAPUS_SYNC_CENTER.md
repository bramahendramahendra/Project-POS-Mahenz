# SOP Deploy — Hapus Sync Center

Menghapus fitur Sync Center: kode BE/FE + 4 tabel DB (sync_conflicts, sync_queue,
sync_history, sync_id_map). Tabel di prod TIDAK terhapus otomatis → drop manual (Langkah 4).

**1. Aktifkan maintenance mode**

```bash
sudo maintenance-on.sh 139.180.214.187
```

**2. Backup database (WAJIB)**

```bash
sudo mkdir -p /opt/pos-mahenz/BE/backups
mysqldump -u pos_user -p --no-tablespaces pos_retail_db > ~/backup_sebelum_hapus_sync_$(date +%Y%m%d_%H%M%S).sql
sudo mv ~/backup_sebelum_hapus_sync_*.sql /opt/pos-mahenz/BE/backups/
ls -lh /opt/pos-mahenz/BE/backups/
```

**3. Update kode BE & FE**

```bash
cd /opt/pos-mahenz
git pull

cd /opt/pos-mahenz/BE
go mod tidy
go build -o pos_api .
sudo systemctl restart pos-backend
sudo systemctl status pos-backend    # harus "active (running)"

cd /opt/pos-mahenz/FE
npm install
npm run type-check
npm run build
sudo rm -rf /var/www/pos-web/dist
sudo cp -r dist /var/www/pos-web/
sudo chown -R $USER:$USER /var/www/pos-web
```

**4. Hapus tabel Sync di database**

Cek dulu isi tabel (semua harus 0). Kalau ada yang tidak 0 → STOP, konsultasi.

```bash
sudo mysql -u pos_user -p pos_retail_db
```
```sql
SELECT
  (SELECT COUNT(*) FROM sync_conflicts) AS conflicts,
  (SELECT COUNT(*) FROM sync_queue)     AS queue,
  (SELECT COUNT(*) FROM sync_history)   AS history,
  (SELECT COUNT(*) FROM sync_id_map)    AS id_map;
EXIT;
```

Kalau semua 0, jalankan script penghapusan:

```bash
cd /opt/pos-mahenz/BE
mysql -u pos_user -p pos_retail_db < database/oneoff/remove_sync_center_20261001.sql
```

Output verifikasi `menu_sync`, `route_sync`, `tabel_sync` harus **0** semua.

**5. Verifikasi & matikan maintenance mode**

Cek: menu Sync Center hilang, Kasir bisa transaksi, Transaksi tampil normal.

```bash
sudo maintenance-off.sh
```

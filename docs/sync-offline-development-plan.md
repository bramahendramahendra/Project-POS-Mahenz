# Prompt Pengembangan Sync Offline — Bertahap per Fase

Rencana ini adalah hasil diskusi desain untuk fitur sync offline aplikasi Desktop (dan
Android) ke backend `Project_pos_mahenz`. Backend baru ini sengaja dipisah dari proyek lama
(`Destop\Project_POS`) agar jadi fondasi bersih — Desktop/Android akan direvamp mengikuti
backend ini, bukan sebaliknya. Jadi kontrak/mekanisme sync lama di `Destop\Project_POS`
(termasuk entity type seperti `cash_drawer` versi lama) **diabaikan**, tidak jadi acuan.

## Cakupan fitur yang disepakati

Awalnya 3 fitur (Kasir, Buka/Tutup Shift, Kas Harian), tapi setelah Fase 1 selesai ditemukan
bahwa **"Shift" di aplikasi ini adalah data master/referensi** (mis. "Shift Pagi 08:00-16:00"),
dikelola CRUD biasa lewat web oleh admin — **bukan** sesi kerja yang dibuka/ditutup tiap kasir.
Sesi kerja kasir yang sesungguhnya (modal awal, setoran akhir, buka/tutup) itu sepenuhnya
direpresentasikan oleh domain **Kas Harian (`cash_drawer`)** — `shift_id` di situ cuma
menunjuk ke salah satu shift-template yang SUDAH ADA (selalu online, tidak pernah dibuat
baru dari device offline). Jadi cakupan final:

1. **Kasir (Checkout)** — sudah selesai & teruji sebelum rencana ini dibuat (dedupe via
   `sync_id_map` selesai di Fase 1).
2. ~~Buka/Tutup Shift~~ — **di-skip**, tidak relevan (lihat penjelasan di atas). `shift_id`
   yang dipakai transaksi/kas harian offline selalu berupa ID yang sudah valid, tidak pernah
   perlu di-resolve dari local_id.
3. **Kas Harian** (buka/tutup kas, setor) — satu-satunya entity baru yang perlu apply-logic
   sync offline.

Fitur lain (produk, pelanggan, pengeluaran, laporan, pengaturan) sengaja **tidak** dibuatkan
sync offline — cukup online-only (lewat web) atau read-only cache di Desktop nanti (di luar
cakupan backend ini).

## Mekanisme inti yang disepakati

- **`sync_id_map`** (tabel baru, sudah dibuat di Fase 1): pemetaan
  `(device_id, local_id, entity_type) → server_id`. Dipakai SEMUA entity yang sync offline
  (transaksi — sudah migrasi dari mekanisme lama di Fase 1; kas harian — menyusul di Fase 2)
  untuk dedupe/idempotency: cegah proses dobel saat device retry push. Mekanisme
  resolve-ID-lintas-entity yang tadinya dirancang untuk kasus shift **tidak jadi dipakai**
  (tidak ada entity yang perlu merujuk shift secara lokal), tapi tabel ini tetap berguna
  murni untuk dedupe kas harian.
- **`sync_queue`**: tetap jadi log/jejak audit untuk SEMUA entity yang di-push (termasuk
  transaksi, sejak Fase 1). Penerapan ke tabel asli tetap terjadi LANGSUNG saat item
  diproses (tidak ditunda ke proses batch terpisah).
- **`sync_conflicts`**: tidak berubah, tetap jadi tempat data server vs lokal dibandingkan
  saat terjadi konflik nyata, diselesaikan lewat UI Sync Center yang sudah ada.

Jalankan **satu fase per sesi/percakapan**, tunggu selesai + terverifikasi (lewat curl,
bukan cuma baca kode), baru lanjut fase berikutnya.

---

## ✅ Fase 1 — Infrastruktur `sync_id_map` + migrasi dedupe transaksi (SELESAI)

Status: **selesai & terverifikasi**. Ringkasan hasil:
- Tabel `sync_id_map` dibuat, kolom `local_id` ditambah ke `sync_queue`, kolom
  `sync_device_id`/`sync_local_id` di `transactions` (migrasi 005) dibatalkan/di-drop
  (migrasi 006).
- `BE/pkg/syncmap/syncmap.go` dibuat (`Resolve`, `Record`), dipakai `ApplySyncTransaction`
  menggantikan mekanisme dedupe lama.
- `CreateQueueItem` sekarang terima parameter `status` eksplisit dan menyimpan `local_id`;
  transaksi ikut tercatat ke `sync_queue` (status langsung `"synced"`, diterapkan seketika).
- Terverifikasi via curl: push transaksi offline 2x (device_id+local_id identik) →
  `server_id` sama persis, stok cuma terpotong sekali, `sync_id_map` cuma 1 baris,
  `sync_queue` 2 baris (jejak tiap push). Regresi checkout online normal dicek, tetap jalan.

<details>
<summary>Prompt asli Fase 1 (arsip, tidak perlu dijalankan ulang)</summary>

```
Lanjut pengembangan fitur sync offline untuk aplikasi POS ini (d:\Develop\Project_pos_mahenz).
Ini FASE 1 dari pengembangan "sync offline untuk Kasir/Kas Harian" yang sudah didiskusikan
dan disepakati sebelumnya. Baca dulu riwayat diskusi terkait di percakapan sebelumnya kalau
tersedia, atau minta ringkasan ke user kalau mulai sesi baru.

KONTEKS MEKANISME (wajib dipahami sebelum mulai):
Saat ini dedupe transaksi offline (mencegah transaksi dobel kalau device retry push karena
network flaky) dilakukan lewat kolom khusus di tabel `transactions`: `sync_device_id` +
`sync_local_id` + unique constraint (lihat migrasi 005_transactions_sync_idempotency.sql).
Fase ini MENGGANTI mekanisme itu dengan tabel pemetaan generik `sync_id_map`, supaya SEMUA
entity yang nanti sync offline (transaksi, kas harian) pakai satu mekanisme yang sama —
bukan kolom khusus per tabel.

`sync_id_map` menyimpan pemetaan (device_id, local_id, entity_type) → server_id. Ini dipakai
untuk dedupe/idempotency: sebelum apply item apapun, cek dulu ke tabel ini — kalau kombinasi
(device_id, local_id, entity_type) sudah ada, langsung balikin server_id yang lama, JANGAN
proses ulang (mencegah transaksi dobel / stok kepotong dua kali saat retry).

Cakupan teknis Fase 1 ini:

1. Migrasi database (file baru, JANGAN edit migrasi 005 yang sudah jalan):
   a. Buat tabel `sync_id_map`:
      CREATE TABLE sync_id_map (
          id          INT AUTO_INCREMENT PRIMARY KEY,
          device_id   VARCHAR(100) NOT NULL,
          local_id    VARCHAR(36)  NOT NULL,
          entity_type VARCHAR(50)  NOT NULL,
          server_id   INT          NOT NULL,
          created_at  DATETIME     DEFAULT CURRENT_TIMESTAMP,
          UNIQUE KEY unique_sync_origin (device_id, local_id, entity_type)
      ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
   b. Tambah kolom `local_id VARCHAR(36) NULL` ke tabel `sync_queue` (setelah kolom device_id).
   c. Batalkan migrasi 005: DROP kolom `sync_device_id`, `sync_local_id`, dan unique key
      `unique_sync_origin` dari tabel `transactions`.

2. Kode BE:
   a. Buat fungsi baru (paket berdiri sendiri, mis. `BE/pkg/syncmap`, supaya tidak membuat
      domain sync dan domain transaction saling bergantung):
      `Resolve(db, deviceID, localID, entityType) (serverID int, found bool, err error)`
      `Record(db, deviceID, localID, entityType, serverID) error`
   b. Ubah `ApplySyncTransaction` di `BE/domain/transaction/repo/transaction_repo.go` untuk
      pakai `Resolve`/`Record` alih-alih kolom khusus.
   c. Pastikan sync_queue TETAP mencatat item transaksi juga (untuk jejak audit), meskipun
      penerapannya tetap langsung/atomik seperti sekarang (BUKAN ditunda).

3. Verifikasi WAJIB via curl: push transaksi offline baru, push ULANG dengan device_id+
   local_id sama (retry) → pastikan server_id sama, stok cuma berkurang sekali. Cek
   `sync_id_map` cuma 1 baris, `sync_queue` ada baris untuk tiap push. Cek kolom sync lama
   di `transactions` sudah hilang.
```

</details>

---

## Fase 2 — Kas Harian offline (buka/tutup/setor) + perbaiki gap update total penjualan

```
Lanjut pengembangan fitur sync offline untuk aplikasi POS ini (d:\Develop\Project_pos_mahenz).
Ini FASE 2 — pastikan Fase 1 (infrastruktur sync_id_map, migrasi dedupe transaksi) sudah
selesai dan terverifikasi sebelum mulai (lihat docs/sync-offline-development-plan.md untuk
ringkasannya). Fase "Buka/Tutup Shift" yang sebelumnya direncanakan terpisah SUDAH DI-SKIP —
"Shift" di aplikasi ini adalah data master/referensi (dikelola CRUD biasa oleh admin lewat
web), BUKAN sesi kerja yang dibuka-tutup tiap kasir. Sesi kerja kasir sepenuhnya
direpresentasikan oleh Kas Harian (`cash_drawer`) — `shift_id` di situ cuma menunjuk ke
salah satu shift-template yang SUDAH ADA (selalu online), jadi TIDAK PERNAH perlu di-resolve
dari local_id. Jangan bangun apa pun terkait "shift offline".

KONTEKS MEKANISME:
Kas Harian (cash_drawer) punya field `shift_id` (opsional, menunjuk shift master data yang
sudah pasti valid — tinggal dipakai langsung, tidak perlu resolve-ID lintas-entity apa pun).
Jadi entity ini SEBENARNYA LEBIH SEDERHANA dari yang tadinya direncanakan: cukup pakai
`sync_id_map` untuk dedupe (pola identik dengan transaksi di Fase 1), tidak perlu mekanisme
resolve-ID tambahan.

Baca dulu domain cash_drawer yang sudah ada (`BE/domain/cash_drawer/**`) untuk memahami:
skema tabel `cash_drawer`, validasi bisnis di jalur online (mis. `GetOpenCashDrawer`, apakah
user boleh buka >1 kas harian sekaligus, cara hitung `expected_balance`/`difference` saat
tutup), dan MULTI method service yang perlu dipakai ulang.

🔴 TEMUAN PENTING YANG WAJIB DIPERBAIKI DI FASE INI (ditemukan saat diskusi perencanaan,
bukan bagian "buat baru" tapi "perbaiki gap yang sudah ada"):
Jalur checkout ONLINE (`BE/domain/transaction/service/transaction_service.go`) melakukan 2
hal saat transaksi dibuat: (1) `GetOpenCashDrawer(userID)` — cari kas harian yang sedang
terbuka milik kasir itu (dicari dinamis by user_id, BUKAN disimpan sebagai FK di tabel
transactions), (2) `cashDrawerRepo.UpdateSales(drawer.ID, ...)` — update total penjualan
berjalan di kas harian itu. `ApplySyncTransaction` (jalur OFFLINE, sudah ada sejak sebelum
Fase 1) **TIDAK melakukan dua hal ini sama sekali** — transaksi yang masuk lewat sync
offline tidak pernah menambah total penjualan di kas harian manapun, membuat laporan Kas
Harian salah/kurang untuk setiap transaksi yang datang dari sync. INI HARUS DIPERBAIKI
sebagai bagian dari fase ini, sebelum atau sepaket dengan membangun apply-logic kas harian
itu sendiri — supaya begitu Kas Harian offline selesai dibangun, angkanya juga BENAR, bukan
cuma "bisa buka/tutup" tapi totalnya kosong.

Cakupan teknis:

1. Perbaiki `ApplySyncTransaction` (`BE/domain/transaction/repo/transaction_repo.go`):
   tambahkan pemanggilan `GetOpenCashDrawer(userID)` + `UpdateSales(drawer.ID, ...)` setelah
   transaksi berhasil dibuat, PERSIS seperti yang dilakukan `transaction_service.go` di
   jalur online (pakai ulang cash_drawer_repo yang sama, jangan tulis ulang logicnya).
   Pastikan ini tidak mengubah alur untuk kasus stok gagal/conflict yang sudah ada (cabang
   itu return sebelum sampai ke langkah update kas harian).

2. Tentukan payload shape untuk sync item entity_type="cash_drawer":
   - action="create" (buka kas): opening_balance, shift_id (integer biasa, SELALU sudah
     valid — tidak ada shift_local_id), open_time, dst sesuai field asli.
   - action="update" (tutup/setor kas): closing_balance, notes, dan cara menunjuk kas
     harian mana yang ditutup (server_id — SELALU sudah ada karena kas harian yang mau
     ditutup pasti sudah dibuka lebih dulu, baik online maupun offline dalam sesi yang sama
     — kalau dalam sesi yang sama, resolve lewat sync_id_map seperti dedupe transaksi biasa,
     BUKAN mekanisme resolve-ID lintas-entity yang rumit).

3. Di `sync_service.go` PushSync: tambah cabang baru untuk `item.EntityType == "cash_drawer"`,
   terpisah dari cabang "transaction" yang sudah ada (JANGAN ubah cabang transaction selain
   perbaikan gap di poin 1).
   - action="create": cek dedupe via `syncmap.Resolve`, kalau belum ada panggil ulang logic
     buka-kas dari cash_drawer_service (pastikan validasi bisnis yang sama tetap jalan,
     mis. tidak bisa buka 2 kas harian aktif sekaligus untuk user yang sama), lalu
     `syncmap.Record(deviceID, localID, "cash_drawer", newDrawerID)`.
   - action="update"/tutup: resolve ID kas harian yang dimaksud (server_id langsung, ATAU
     via syncmap kalau dibuka di batch offline yang sama), lalu panggil logic tutup-kas dari
     cash_drawer_service (pakai ulang perhitungan expected_balance/difference yang sudah ada).
   - Item tetap dicatat ke sync_queue (audit), diterapkan langsung (tidak ditunda).

4. Tangani kasus konflik yang masuk akal, misalnya: kas harian yang mau ditutup ternyata
   sudah ditutup duluan oleh device lain, atau user mencoba buka kas harian baru padahal
   masih ada yang aktif (race condition, device offline belum tahu kondisi terbaru). Rancang
   pesan error/conflict yang jelas — diskusikan dulu ke user kalau perilaku "benar" untuk
   kasus ini tidak jelas/perlu keputusan bisnis.

5. Verifikasi via curl:
   a. Push "buka kas harian" offline (action=create) → cek record cash_drawer beneran
      tercipta dengan data benar, tercatat di sync_id_map.
   b. Push LAGI dengan device_id+local_id sama (retry) → pastikan idempotent, TIDAK bikin
      kas harian dobel.
   c. Push transaksi offline (entity_type="transaction") SETELAH kas harian di atas dibuka
      → cek `total_sales`/`total_cash_sales` di cash_drawer BENAR bertambah sesuai transaksi
      itu (ini pembuktian utama perbaikan gap di poin 1).
   d. Push "tutup kas harian" merujuk ke kas harian dari langkah (a) → cek status berubah,
      `expected_balance`/`difference` terhitung benar berdasarkan transaksi-transaksi yang
      masuk di langkah (c).

Di akhir: laporkan hasil verifikasi dengan data konkret (ID, angka saldo/total sebelum-
sesudah). Bug jelas/kecil perbaiki langsung, temuan besar/keputusan desain didiskusikan
dulu ke user.
```

---

## Fase 3 — Uji integrasi penuh + rapikan UI Sync Center

```
Lanjut pengembangan fitur sync offline untuk aplikasi POS ini (d:\Develop\Project_pos_mahenz).
Ini FASE 3 (terakhir) — pastikan Fase 1 & 2 sudah selesai dan terverifikasi.

Cakupan:

1. Simulasi skenario penuh "sehari kerja offline" lewat SATU push berisi item berurutan
   (device_id sama untuk semua item, meniru satu device yang beneran offline seharian):
   a. Buka kas harian (local_id="kas-1", shift_id=<ID shift master yang sudah ada>)
   b. 3-5 transaksi penjualan
   c. Tutup/setor kas harian (merujuk kas harian dari langkah a)
   Verifikasi SEMUA langkah berhasil, dan angka-angka di kas harian yang tertutup (total
   penjualan, selisih kas) MATEMATIS BENAR berdasarkan transaksi-transaksi yang masuk di
   langkah (b) — cross-check manual seperti yang biasa dilakukan di sub-fase testing laporan
   sebelumnya.

2. Uji idempotency SELURUH BATCH: kirim ULANG persis batch yang sama (semua device_id+
   local_id identik) → pastikan TIDAK ADA yang dobel sama sekali (kas harian & transaksi
   idempotent), dan hasil akhirnya identik dengan push pertama.

3. Uji skenario konflik nyata: coba buat kondisi kas harian yang BENAR-benar bentrok (mis.
   kas harian yang sama coba ditutup dua kali dari dua push terpisah dengan data penutupan
   berbeda) → pastikan masuk ke sync_conflicts dengan benar dan bisa diselesaikan lewat UI
   resolve yang sudah ada (Terima Server / Pakai Data Lokal) — verifikasi hasilnya masuk
   akal untuk kasus kas harian, bukan cuma expense seperti yang sudah dites sebelumnya.

4. Update FE Sync Center:
   - Tambah label "Kas Harian" ke `ENTITY_TYPE_LABEL` di
     `FE/src/features/operational/sync/components/ConflictList.tsx` dan
     `ENTITY_TYPE_LABEL` di `SyncQueueTableColumns.tsx`.
   - Playwright: pastikan tabel Antrian Sync dan Riwayat Sinkronisasi menampilkan entity
     cash_drawer dengan rapi (bukan raw string type), tidak ada console error.

5. Login sebagai owner/owner123, jalankan Playwright untuk screenshot akhir halaman Sync
   Center menampilkan data hasil simulasi di atas, pastikan semua tampil benar.

Di akhir: berikan ringkasan lengkap fitur sync offline yang sudah selesai (cakupan final:
Kasir + Kas Harian, dedupe seragam), daftar bug yang ditemukan+diperbaiki sepanjang 3 fase,
dan konfirmasi tidak ada regresi ke fitur yang sudah ada sebelumnya (transaksi checkout
normal, laporan-laporan yang bergantung ke data transaksi/kas harian).
```

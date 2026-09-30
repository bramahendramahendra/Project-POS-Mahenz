# Rancangan Penghapusan Menu Sync Center

> Status: **RANCANGAN — belum ada penghapusan.**
> Tujuan: daftar lengkap & berurutan untuk menghapus TOTAL fitur "Sync Center"
> (kode BE + FE + tabel DB) tanpa merusak fitur lain (terutama Transaksi/Kasir).
> Ditulis untuk direview dulu sebelum eksekusi.

---

## 1. Latar Belakang

Menu **Sync Center** (`/sync`, key `operasional.sync`) berfungsi menyinkronkan data dari
aplikasi **desktop/Android offline** ke server (push data, deteksi konflik, riwayat sync).
User tidak memakai aplikasi offline, jadi menu ini tidak terpakai dan akan dihapus total
beserta kode BE, FE, dan tabel DB-nya.

---

## 2. Hasil Investigasi — Tingkat Kesulitan

Ini **bukan sekadar hapus folder**. Fitur sync punya 3 titik keterkaitan (dependency) ke
kode inti yang WAJIB dilepas hati-hati, kalau tidak aplikasi gagal build / fitur lain rusak:

1. **Transaksi bergantung ke domain sync.** `transaction_repo.go` meng-import
   `domain/sync/dto` (`SyncTransactionPayload`) dan `pkg/syncmap`, serta punya 3 metode
   khusus sync: `ApplySyncTransaction`, `UpdateFromSync`, `ReturnStockForRejectSync`.
   Menghapus folder sync tanpa mencabut ini → BE gagal build.
2. **Navbar (header FE) pakai badge sync.** `Navbar.tsx` meng-import `useSyncStatus` dan
   menampilkan badge jumlah konflik di ikon lonceng. Menghapus folder FE sync tanpa
   merapikan ini → FE gagal build.
3. **`device_source` BUKAN milik sync — JANGAN disentuh.** Kolom
   `transactions.device_source` ('desktop'/'web'/'android') dipakai transaksi normal,
   penomoran kode, backdate, & laporan. Tetap dipertahankan.

Catatan aman: tidak ada cron/background job sync di server; jalur checkout online (kasir)
terpisah total dari jalur sync. Jadi penghapusan aman ASALKAN berurutan.

---

## 3. Daftar Lengkap Yang Dihapus / Diubah

### A. Backend — hapus file penuh
- `BE/domain/sync/` (seluruh folder, 7 file: dto/dto_sync.go, handler/sync_handler.go,
  model/sync_model.go, repo/sync_repo.go, repo/sync_repo_interface.go,
  service/sync_service.go, service/sync_service_interface.go)
- `BE/routes/segment/sync_routes.go`
- `BE/pkg/syncmap/syncmap.go` (folder `pkg/syncmap`) — SETELAH pemanggilnya dilepas

### B. Backend — edit (lepas benang sync)
- `BE/routes/protected_routes.go` (baris 46): hapus `segment.SyncRoutes(r)`
- `BE/domain/transaction/repo/transaction_repo_interface.go`: hapus 3 metode dari
  interface: `UpdateFromSync`, `ReturnStockForRejectSync`, `ApplySyncTransaction` +
  import `cash_drawer_repo` bila jadi tak terpakai (cek dulu).
- `BE/domain/transaction/repo/transaction_repo.go`: hapus implementasi metode sync
  (`ApplySyncTransaction`, `applySyncTransactionOnce`, `ReturnStockForRejectSync`,
  `UpdateFromSync`) + import `dto_sync "pos_api/domain/sync/dto"`, `pkg/syncmap`, dan
  `cash_drawer_repo` (import cash_drawer HANYA dipakai metode sync di file ini — WAJIB ikut
  dihapus, kalau tidak error "imported and not used").
  **⚠️ KOREKSI PENTING (hasil cek ulang):**
  - `isDuplicateTransactionCodeError` **JANGAN DIHAPUS** — dipakai juga oleh jalur ONLINE
    `createOnce()` (retry kode transaksi duplikat), bukan cuma sync.
  - **JANGAN sentuh** `Create()`, `createOnce()`, `generateTransactionCodeQuery`,
    `device_source`, prefixMap.
  - Cek: import `cash_drawer_repo` di transaction_repo.go & _interface.go hanya untuk
    signature/isi `ApplySyncTransaction`. Setelah metode itu hilang, hapus import-nya di
    KEDUA file. (transaction_SERVICE tetap pakai cashDrawerRepo — itu file berbeda, JANGAN
    disentuh.)
  - **IMPORT YANG HARUS DIHAPUS** dari transaction_repo.go (hasil baca baris demi baris —
    ketiganya HANYA dipakai jalur sync): `encoding/json`, `dto_sync "pos_api/domain/sync/dto"`,
    `pos_api/pkg/syncmap`, `cash_drawer_repo "pos_api/domain/cash_drawer/repo"`.
  - **IMPORT YANG TETAP** (dipakai fungsi non-sync, JANGAN hapus): `stderrors`, `fmt`,
    `request_helper` (GetAll), `time_helper`, `model_product`, `product_repo`,
    `transaction/dto`, `transaction/model`, `mysql`, `gorm`.
  - `ReturnStockForRejectSync` & `UpdateFromSync` (transaction) memakai import yang juga
    dipakai fungsi lain → menghapusnya tidak menyisakan import menggantung.
- `BE/domain/expense/repo/expense_repo.go` + `expense_repo_interface.go`: hapus
  `UpdateFromSync` (dead code setelah sync hilang).

### C. Frontend — hapus file penuh
- `FE/src/features/operational/sync/` (seluruh folder: index.ts, sync.api.ts,
  sync.types.ts, SyncCenterPage.tsx, components/*, hooks/useSyncStatus.ts, .gitkeep)

### D. Frontend — edit (lepas benang sync)
- `FE/src/shared/components/layouts/Navbar.tsx`: hapus import `useSyncStatus` (baris 6),
  hapus `const { conflictCount } = useSyncStatus()` (baris 32), hapus badge konflik di
  ikon lonceng (baris 61–65). **Keputusan:** ikon lonceng dipertahankan tanpa badge, atau
  dihapus sekalian bila lonceng memang hanya untuk sync (akan dicek saat eksekusi —
  default: pertahankan lonceng polos tanpa badge, paling aman & minim perubahan).
- `FE/src/app/router.tsx`: hapus baris 50 (lazy import SyncCenterPage) & baris 114
  (route `{ path: ROUTES.SYNC, menuKey: 'operasional.sync', ... }`).
- `FE/src/shared/constants/routes.ts` (baris 43): hapus `SYNC: '/sync'`.
- `FE/src/shared/constants/queryKeys.ts` (baris ~132–138): hapus blok `sync: {...}`.
- `FE/src/shared/components/StatusBadge/StatusBadge.tsx` (opsional): status
  'synced'/'unsynced'/'conflict' hanya dipakai sync. Boleh dihapus, tapi RENDAH prioritas &
  bukan penyebab error. **Keputusan: pertahankan** (tidak mengganggu, hindari risiko).

### E. Database — drop tabel (khusus sync)
- `sync_conflicts`, `sync_queue`, `sync_history`, `sync_id_map` + index terkait.
- Untuk DB yang sudah ada (lokal/prod): jalankan `DROP TABLE` manual (script disiapkan).
- Untuk instalasi baru: hapus definisi tabel + index dari
  `BE/database/migrations/001_init_schema.sql` (baris ~540–599 tabel, ~665–673 index).

### F. Seed menu — `BE/database/migrations/002_seed_data.sql`
- Hapus INSERT menu `operasional.sync` (baris ~182–184).
- Hapus route_registry `('/sync', 'Sync Center')` (baris ~244).
- Rapikan daftar exclude Owner grant: keluarkan `'operasional.sync'` dari `NOT IN (...)`
  (baris ~280) + komentar (baris ~257).
- Untuk DB existing: DELETE manual dari `role_menu_access`, `menus`, `route_registry`
  untuk key `operasional.sync` / path `/sync` (script disiapkan).

---

## 4. JANGAN DISENTUH (bukan sync / dipakai fitur lain) — hasil verifikasi ulang
- Kolom `transactions.device_source` + seluruh jalur `Create()` / `createOnce()` / backdate.
- Fungsi `isDuplicateTransactionCodeError` (dipakai jalur online, BUKAN cuma sync).
- FE `device_source` di `transactions.types.ts`, `PaymentModal`, `BackdatePaymentModal`.
- `pkg/pricing` (dipakai checkout online juga).
- `cash_drawer` service/repo (dipakai kas harian). Di `transaction_SERVICE` import
  `cash_drawer_repo` TETAP dipakai `Create()`/`Void()` — JANGAN disentuh. Yang dihapus
  hanya import cash_drawer di transaction_REPO (+interface) yang khusus untuk metode sync.
- Tabel `log_schedulers` (bukan sync).
- `StatusBadge` status 'synced/unsynced/conflict' (dipertahankan, tak mengganggu).

## 4b. Hasil Verifikasi Ulang (cek teliti sebelum eksekusi)
- ✅ TIDAK ada file test (Go/FE) yang menyinggung fitur sync → tak ada test yang gagal.
- ✅ `GetEntitySnapshot` (sync_repo) hanya MEMBACA `updated_at` tabel transaction/expense
  (read-only) → drop tabel sync tidak memengaruhi data transaksi/expense.
- ✅ Semua referensi FE ke `/sync` terkurung di: folder sync, Navbar.tsx, router.tsx,
  routes.ts, queryKeys.ts. Tidak ada pemakaian tersembunyi di menu lain.
- ✅ `useSyncStatus` punya efek samping toast "Ada X konflik sinkronisasi" + localStorage
  key `sync_conflict_notified_count` — hilang total saat hook dihapus (memang diinginkan).
- ✅ Metode sync di BE (`ApplySyncTransaction`, `UpdateFromSync`, `ReturnStockForRejectSync`
  di transaction; `UpdateFromSync` di expense) HANYA dipanggil dari `sync_service.go`.
  Setelah domain sync hilang → dead code, aman dihapus bersama.

---

## 5. Urutan Eksekusi (bertahap, verifikasi tiap tahap)

**Tahap A — Frontend dulu (lepas benang + hapus folder):**
1. Rapikan `Navbar.tsx` (buang badge sync).
2. Bersihkan `router.tsx`, `routes.ts`, `queryKeys.ts`.
3. Hapus folder `FE/src/features/operational/sync/`.
4. Verifikasi: `npm run type-check` + `npm run lint` + `npm run build` BERSIH.

**Tahap B — Backend (lepas benang + hapus folder):**
5. Hapus `segment.SyncRoutes(r)` di protected_routes.go.
6. Lepas metode sync dari transaction repo (+interface) & expense repo (+interface).
7. Hapus folder `BE/domain/sync/` + `sync_routes.go`.
8. Hapus `BE/pkg/syncmap/`.
9. Verifikasi: `go build ./...` BERSIH.

**Tahap C — Database:**
10. Jalankan script DROP tabel sync + DELETE seed menu (di DB lokal dulu).
11. Update file migrasi (001 & 002) untuk instalasi baru.

**Tahap D — Verifikasi menyeluruh:**
12. Nyalakan BE+FE, tes browser: login, buka Kasir & buat transaksi (pastikan checkout
    normal jalan), buka Transaksi (device_source tampil), buka header (lonceng normal),
    pastikan menu Sync Center HILANG dari sidebar, tidak ada error konsol.
13. Bersihkan file tes.

---

## 6. Risiko & Mitigasi
- **Risiko:** transaksi rusak karena benang sync tercabut sembarangan.
  **Mitigasi:** lepas metode sync satu per satu, build BE setelah tiap perubahan besar;
  jalur `Create()` tidak disentuh.
- **Risiko:** FE gagal build karena Navbar/queryKeys.
  **Mitigasi:** rapikan referensi SEBELUM hapus folder; type-check tiap langkah.
- **Risiko:** DB production. **Mitigasi:** kerjakan & tes di DB lokal (RELEASE_MODE=dev)
  dulu; untuk prod, sediakan script terpisah + backup (masuk ke catatan deploy).
- **Rollback:** semua via git (kode) + `DROP TABLE` bisa dibuat ulang dari migrasi 001
  bila perlu (tabel sync kosong/tak terpakai, jadi drop aman).

---

## 7. Catatan Deploy
Perubahan ini menyentuh BE + FE + DB. Untuk rilis ke produksi nanti perlu langkah DB
khusus (DROP tabel sync + DELETE seed menu). Akan disiapkan script SQL terpisah +
disarankan backup DB sebelum eksekusi (pola sama SOP deploy Laba Rugi).

---

## 8. Aturan Kerja
- Tidak ada penghapusan sampai rancangan ini disepakati.
- Bertahap; verifikasi (build/type-check/lint + browser) tiap tahap.
- Kerjakan di lokal (RELEASE_MODE=dev). File tes di folder `testing`.

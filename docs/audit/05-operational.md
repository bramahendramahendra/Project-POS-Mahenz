# Audit: Modul Operational (shifts / sync)

Scope: `FE/src/features/operational/**`, `BE/domain/shift/**`, `BE/domain/sync/**`, `BE/routes/segment/{shift_routes.go,sync_routes.go}`

## Temuan

### 1. [KRITIS] Endpoint `POST /push` sync tanpa permission gate
- **File**: `BE/routes/segment/sync_routes.go` (baris 34)
- **Masalah**: Semua route lain di file ini (conflicts, conflicts/count, resolve, queue, history) di-gate `perm(...)`, tapi `PushSync` — endpoint paling powerful di domain ini — tidak.
- **Dampak**: Siapa pun yang terautentikasi bisa push data sync tanpa pengecekan izin sama sekali.

### 2. [Tinggi] Payload transaksi sync dipercaya mentah dari client
- **File**: `BE/domain/sync/dto/dto_sync.go` (baris 104-118, `SyncTransactionPayload`)
- **Masalah**: `Subtotal`, `Discount`, `Tax`, `TotalAmount`, `PaymentAmount`, `ChangeAmount` diteruskan sebagai JSON opaque tanpa rekomputasi server-side.
- **Pembanding**: `cash_drawer_service.go` (baris 198-199) menghitung ulang `difference` server-side.
- **Dampak**: Client yang di-crafting/berbahaya bisa push total transaksi sembarangan lewat jalur sync.

### 3. [Tinggi] Payload non-transaksi ditandai "synced" walau tidak benar-benar diterapkan
- **File**: `BE/domain/sync/service/sync_service.go` (baris 84-91)
- **Masalah**: Untuk entity non-transaksi (product, customer, stock, dll), queue item langsung ditandai `"synced"` (baris 91) tanpa ada kode yang benar-benar menerapkan payload ke tabel target.
- **Dampak**: Client mendapat respons sukses padahal data (misal edit produk offline) hilang diam-diam.

### 4. [Medium] Tidak ada dedupe by (device_id, local_id)
- **File**: `BE/domain/sync/repo/sync_repo.go` (baris 254-262, `CreateQueueItem`), juga cabang `ServerID == 0` di `ApplySyncTransaction` (`sync_service.go` baris 61-82)
- **Dampak**: Push yang di-retry (misal karena network flaky) bisa membuat queue item duplikat atau menerapkan transaksi/pengurangan stok dua kali.

### 5. [Tinggi] Route baca Shift tidak konsisten di-gate
- **File**: `BE/routes/segment/shift_routes.go` (baris 25, 26, 28) — `GetAll`, `GetOptions`, `GetByID` tidak di-gate, sementara Create/Update/Delete/ToggleStatus/GetSummary di-gate.
- **Pembanding**: Semua GET di sync routes di-gate.
- **Dampak**: Siapa pun yang login bisa membaca semua data shift terlepas dari permission yang diberikan.

### 6. [Low] Kolom `start_time` di tabel Shift tidak `sortable`
- **File**: `FE/src/features/operational/shifts/components/ShiftTableColumns.tsx` (baris 29-35)
- **Pembanding**: Kolom sibling `name` (24) dan `is_active` (41) sudah `sortable`. BE (`shift_repo.go` baris 54-59) sudah mendukung sort by `start_time`, jadi ini gap nyata bukan keterbatasan BE.

### 7. [Low] Fitur sync status FE tidak pernah berfungsi (dead code)
- **File**: `FE/src/features/operational/sync/sync.types.ts` (baris 5-11, `SyncStatusData`), `useSyncStatus.ts` (baris 24-30, hardcode `status: null`)
- **Dampak**: `SyncStatusCard.tsx` (baris 39) selalu render state `idle`; state `syncing`/`success`/`error` tidak pernah terpakai.

### 8. [Low] Zustand store sync tidak dipakai
- **File**: `FE/src/features/operational/sync/sync.store.ts` — didefinisikan tapi tidak dipakai; `ConflictList.tsx` implementasi ulang state secara lokal.

### 9. [Low] Cache key pagination sync history tidak lengkap
- **File**: `sync.api.ts` (baris 20) — `queryKeys.sync.history()` tidak menyertakan `page`/`page_size`
- **Pembanding**: `queryKeys.shifts.list(filter...)` (`shifts.api.ts` baris 12) menyertakan filter lengkap.
- **Dampak**: History sync yang dipaginasi bisa menyajikan halaman cache basi.

## Yang Sudah Baik
Shift create/update payload memakai satu type via spread (tidak duplikasi manual). Validasi field FE/BE shift selaras. Mutation shift konsisten menampilkan toast sukses.

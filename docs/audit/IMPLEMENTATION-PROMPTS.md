# Prompt Implementasi Perbaikan — Bertahap per Fase

Temuan audit terlalu banyak (~52) untuk dikerjakan dalam satu sesi. File ini berisi prompt siap pakai untuk tiap fase. Jalankan satu fase per sesi/percakapan baru dengan Claude Code, tunggu selesai + verifikasi, baru lanjut fase berikutnya.

Setiap prompt sudah merujuk ke file laporan spesifik di `docs/audit/` supaya Claude tidak perlu membaca ulang seluruh codebase dari nol.

---

## Fase 1 — Security Kritis (kerjakan lebih dulu, sebelum fase lain)

### 1a. Privilege escalation di Settings/Roles/Users
```
Baca docs/audit/08-settings.md temuan #1 sampai #5 (privilege escalation di SetRoleAccess,
UserService.Update, UserService.ToggleStatus, RoleService.Update tanpa guard IsSystem, dan
tidak adanya proteksi last-admin). Lakukan perbaikan untuk kelima temuan tersebut:
1. Tambahkan guard di SetRoleAccess (BE/domain/access/service/service_access.go) agar
   request dengan RoleID = role milik pemanggil sendiri, atau role dengan is_system=true,
   ditolak (kecuali dilakukan oleh role "owner").
2. Tambahkan pengecekan id != currentUserID di UserService.Update ketika field role_id
   berubah, dan di UserService.ToggleStatus. Pastikan currentUserID diteruskan dari handler
   (ambil dari context/JWT claims sesuai pola yang sudah ada di endpoint lain).
3. Tambahkan guard IsSystem di RoleService.Update, konsisten dengan Delete dan ToggleStatus
   yang sudah benar.
4. Tambahkan pengecekan "last active admin" sebelum user dengan role admin/owner terakhir
   dihapus, di-toggle nonaktif, atau di-demote.
Setelah selesai, jelaskan test manual apa yang perlu saya lakukan untuk verifikasi (skenario
self-escalation harus gagal, admin terakhir tidak bisa dihapus/dinonaktifkan/didemote).
```

### 1b. Checkout tidak menghitung ulang uang di server (Sales)
```
Baca docs/audit/06-sales.md temuan #1, #2, #3, #6 (BE tidak menghitung ulang subtotal/
diskon/pajak/total transaksi, endpoint create tanpa permission gate, shift_id tidak
divalidasi terhadap user, dan tidak ada upper-bound validasi diskon). Perbaiki:
1. Di BE/domain/transaction/service atau repo, hitung ulang subtotal per item dari harga
   produk asli (bukan dari req.Price) dikali quantity, kurangi diskon, tambah pajak, dan
   bandingkan/timpa TotalAmount yang dikirim client dengan hasil hitungan server. Tolak
   request jika ada selisih signifikan (atau selalu pakai hasil hitungan server, abaikan
   nilai dari client).
2. Tambahkan permission middleware (can_create) pada route POST transaction create di
   BE/routes/segment/transaction_routes.go, konsisten dengan pola supplier_purchase_routes.go.
3. Validasi bahwa shift_id yang dikirim benar milik user yang sedang login sebelum
   memproses transaksi.
4. Tambahkan validasi max pada Discount/DiscountItem/Tax di dto_transaction.go agar tidak
   bisa melebihi subtotal.
Jangan ubah kontrak response API yang sudah dipakai FE kecuali benar-benar perlu — cek
FE/src/features/sales/cashier/cashier.api.ts untuk memastikan kompatibilitas.
```

### 1c. Void transaksi tidak membalikkan cash-drawer (Sales)
```
Baca docs/audit/06-sales.md temuan #4 dan #5 (void tidak reverse cash-drawer sales total,
dan update cash-drawer tidak atomik dengan create/void transaksi). Perbaiki:
1. Pada BE/domain/transaction/repo/transaction_repo.go fungsi Void, tambahkan pemanggilan
   balik ke cash-drawer repo untuk mengurangi total_cash_sales sejumlah nilai transaksi yang
   di-void (mirror dari UpdateSales yang dipanggil saat create).
2. Bungkus update cash-drawer dan create/void transaksi dalam satu DB transaction (gunakan
   pola transaction Go yang sudah dipakai di bagian lain codebase ini, misal supplier_purchase
   create) supaya keduanya konsisten atau sama-sama rollback saat gagal.
Setelah selesai, jelaskan skenario test manual: create transaksi cash lalu void, pastikan
total_cash_sales cash-drawer kembali ke nilai semula.
```

### 1d. Auth: refresh token, JWT secret, rate limiting
```
Baca docs/audit/09-auth.md temuan #1, #2, #3, #4, #5 (refresh token dead-config, field
mismatch access_token vs token, JWT secret placeholder di config_prod.json, tidak ada rate
limit di login, logout fail-open). Perbaiki:
1. Perbaiki BE/domain/auth/dto/dto_auth.go RefreshResponse agar field name konsisten dengan
   yang dibaca FE di FE/src/services/api.client.ts (pilih salah satu nama, `access_token`
   direkomendasikan untuk konsistensi, lalu selaraskan kedua sisi).
2. Perbaiki BE/domain/auth/service/auth_service.go agar expiry refresh token benar-benar
   memakai config.Cfg.RefreshTokenExpire, bukan TokenExpire yang sama dengan access token.
3. Ganti nilai SecretKey di BE/config/config_prod.json agar dibaca dari environment variable,
   bukan hardcoded string di file yang ter-commit. Jangan expose nilai secret asli di kode/
   commit message.
4. Tambahkan rate limiting sederhana (in-memory atau middleware) pada endpoint POST
   /auth/login, misal maksimal N percobaan gagal per username/IP dalam window waktu tertentu.
5. Perbaiki alur logout agar tidak menghapus session lokal FE sebelum server konfirmasi
   sukses menghapus session, ATAU pastikan best-effort retry agar session BE tetap terhapus.
Jelaskan implikasi setiap perubahan terhadap sesi user yang sedang aktif saat deploy.
```

### 1e. Operational: sync push tanpa permission, payload sync dipercaya mentah
```
Baca docs/audit/05-operational.md temuan #1, #2, #3, #5. Perbaiki:
1. Tambahkan permission middleware pada route POST /push di BE/routes/segment/sync_routes.go,
   konsisten dengan route lain di file yang sama.
2. Tambahkan permission middleware pada Shift GetAll/GetOptions/GetByID, konsisten dengan
   Create/Update/Delete yang sudah di-gate.
3. Untuk payload transaksi via sync (SyncTransactionPayload), terapkan rekomputasi nilai
   uang server-side yang sama seperti perbaikan checkout di Fase 1b (reuse logic yang sama
   bila memungkinkan agar tidak duplikasi aturan bisnis).
4. Untuk payload non-transaksi (product/customer/stock) di sync_service.go, jangan tandai
   status "synced" sebelum benar-benar diterapkan ke tabel target — implementasikan apply
   logic-nya atau tandai sebagai "pending_apply"/error jika belum didukung.
Laporkan bagian mana yang butuh keputusan produk (misal: entity apa saja yang benar-benar
perlu didukung sync-nya) sebelum melanjutkan implementasi penuh.
```

---

## Fase 2 — Correctness Tinggi (uang & data tidak sinkron)

### 2a. Fitur pembayaran piutang rusak (Customers)
```
Baca docs/audit/02-customers.md temuan #1. Perbaiki mismatch DTO antara FE dan BE untuk
pembayaran piutang:
1. Putuskan kontrak final: apakah payment_method wajib dikirim FE, dan apakah payment_date
   perlu didukung BE (atau BE tetap pakai waktu server saat itu, dan FE tidak perlu
   mengirim payment_date sama sekali — pilih salah satu, jangan setengah-setengah).
2. Sesuaikan FE/src/features/customers/receivables/receivables.types.ts (CreatePaymentPayload)
   dan FE/src/features/customers/receivables/components/PaymentRecordModal.tsx dengan
   kontrak final tersebut.
3. Sesuaikan BE/domain/receivable/dto/dto_receivable.go (PayRequest) dan
   BE/domain/receivable/repo/receivable_repo.go (CreatePayment) agar konsisten.
Setelah selesai, test manual: buat piutang lalu lakukan pembayaran dari UI, pastikan tidak
ada error validasi dan payment tercatat dengan benar.
```

### 2b. Update purchase order salah hitung total (Procurement)
```
Baca docs/audit/04-procurement.md temuan #1. Perbaiki BE/domain/supplier_purchase/repo/
purchase_repo.go fungsi Update agar:
1. Menghitung total_amount dengan cara yang sama seperti Create (subtotal dikurangi
   discount_amount).
2. Menyertakan discount_amount, payment_status, paid_amount, payment_method ke dalam
   statement UPDATE (saat ini hilang).
Setelah selesai, test manual: buat purchase order dengan diskon, edit datanya, pastikan
total_amount dan field pembayaran tersimpan benar setelah edit.
```

### 2c. Finance: authorization expense + transaksi cash-drawer update
```
Baca docs/audit/03-finance.md temuan #1 dan #2. Perbaiki:
1. Tambahkan pengecekan role/kepemilikan di BE/domain/expense/service/expense_service.go
   Update dan Delete, konsisten dengan pola di BE/domain/cash_drawer/service/
   cash_drawer_service.go (owner/admin atau pemilik record sendiri saja yang boleh edit/hapus).
2. Bungkus penulisan expense dan update cash-drawer (UpdateExpenses) dalam satu DB
   transaction, dan jangan buang error dari UpdateExpenses — propagate agar caller tahu
   dan bisa rollback bila gagal.
Test manual: login sebagai cashier A, coba edit/hapus expense milik cashier B, pastikan
ditolak (403). Simulasikan kegagalan update cash-drawer, pastikan seluruh operasi rollback.
```

### 2d. Reporting: permission gate + konsistensi revenue
```
Baca docs/audit/07-reporting.md temuan #1 dan #2. Perbaiki:
1. Tambahkan permission middleware (can_view) pada seluruh route sales & stock report di
   BE/routes/segment/report_routes.go, konsisten dengan profit-loss dan cashier-performance.
2. Selaraskan basis perhitungan revenue antara profit-loss dan sales report (pertimbangkan
   apakah profit-loss seharusnya pakai t.total_amount seperti sales report, atau sales
   report yang perlu breakdown pre/post-diskon) — sebelum implementasi, tanyakan ke saya
   definisi bisnis yang benar jika tidak jelas dari kode yang ada.
```

---

## Fase 3 — Konsistensi Menengah (tidak mendesak, tapi berulang)

```
Baca semua file docs/audit/01-products.md sampai 09-auth.md, kumpulkan seluruh temuan
berkategori "Low"/"Medium" yang berpola sama: "Update...Payload di-duplikasi manual instead
of Partial<Create...Payload>". Untuk setiap modul yang disebutkan (products/units, expenses,
purchases, suppliers, roles, users, customers), ubah type Update menjadi turunan dari type
Create memakai Partial<> atau Omit<>, sesuai pola yang sudah benar di modul lain (categories,
shifts). Jangan ubah behavior runtime, ini murni refactor type-level. Setelah selesai jalankan
type-check FE (tsc --noEmit) untuk memastikan tidak ada breakage.
```

```
Baca semua file docs/audit/*.md, kumpulkan seluruh temuan "kolom sortable tidak konsisten"
di berbagai TableColumns.tsx. Tambahkan sortable: true pada kolom yang disebutkan di setiap
temuan, dan pastikan BE endpoint terkait memang sudah mendukung sort_by untuk kolom tersebut
sebelum menambahkannya di FE (jika BE belum mendukung, catat sebagai follow-up terpisah,
jangan tambahkan sortable palsu).
```

```
Baca docs/audit/01-products.md temuan #1 (field description produk setengah jadi). Putuskan
salah satu: (a) tambahkan input description di ProductFormModal + kolom description di BE
model/repo/migration, atau (b) hapus seluruh referensi description dari FE schema/types/
ProductDetailModal karena memang tidak didukung BE. Tanyakan preferensi saya sebelum memilih
arah, karena ini mengubah scope fitur bukan sekadar bug fix.
```

---

## Fase 4 — Polish UX/Konsistensi Rendah

```
Baca seluruh temuan berkategori "Low" di docs/audit/*.md yang belum masuk fase 1-3
(missing toast sukses, RoleGuard tidak konsisten dipakai, label required (*) tidak konsisten,
export report tidak konsisten, dead code seperti sync.store.ts yang tidak dipakai atau
BearerAuthMiddleware yang tidak dipasang). Perbaiki satu per satu, screenshot/jelaskan
sebelum-sesudah untuk perubahan UI jika memungkinkan dijalankan di dev server.
```

---

## Catatan Pemakaian

- Jalankan fase secara berurutan (1 → 2 → 3 → 4). Fase 1 wajib selesai lebih dulu karena
  menyangkut keamanan aktif (privilege escalation, manipulasi harga checkout, secret bocor).
- Setiap prompt bisa langsung di-paste sebagai pesan baru ke Claude Code.
- Setelah tiap sub-fase selesai, jalankan test manual yang diminta di masing-masing prompt
  sebelum lanjut ke sub-fase berikutnya, terutama untuk fase 1 dan 2 yang menyentuh alur uang.
- Gunakan git commit terpisah per sub-fase supaya mudah di-review/rollback jika ada regresi.

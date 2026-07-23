# Rencana Debug Testing Menyeluruh — POS Mahenz

Dokumen ini berisi rencana pengujian debug seluruh menu dan fitur aplikasi (FE + BE terkait), dipecah menjadi fase-fase kecil supaya bisa dicicil dan tetap terarah. Setiap fase punya prompt siap-pakai — cukup rujuk fase & file ini ke asisten (mis. "Jalankan Fase 1.1 dari docs\DEBUG_TESTING_PLAN.md"), tidak harus berurutan dalam satu sesi.

---

## Prasyarat & Cara Menjalankan (wajib dibaca sebelum fase manapun)

**1. Jalankan BE & FE (dua proses terpisah, biarkan tetap hidup selama fase berjalan):**
```
# Terminal 1 — Backend
cd BE
go run main.go
# → API di http://localhost:8080/api, health check di http://localhost:8080/health

# Terminal 2 — Frontend
cd FE
npm run dev
# → aplikasi di http://localhost:3000
```
Kalau salah satu proses gagal start (port bentrok, error koneksi DB, dsb), itu sendiri sudah temuan Fase 0 — jangan lanjut ke fase lain sebelum keduanya bisa jalan normal.

**2. Pastikan database (`pos_retail_db`, MySQL/MariaDB lokal) sudah ada skema & seed-nya.** Kalau kosong (misal habis di-drop untuk test "dari kosong"), jalankan dulu:
```
BE/database/migrations/001_init_schema.sql
BE/database/migrations/002_seed_data.sql
```

**3. Kredensial login untuk tiap role (dari seed data):**
- Owner: `owner` / `owner123`
- Admin: `admin` / `admin123`
- Kasir: **tidak ada user seeded** — kalau fase butuh role Kasir, buat dulu user baru dengan role Kasir lewat Manajemen User (login sebagai Owner/Admin dulu), baru login sebagai user itu.

**4. Cara "men-drive browser asli" secara konkret:** tidak ada tool GUI browser interaktif di sini — cara yang dipakai adalah menulis script **Playwright** (Node.js, headless Chromium) ke folder scratchpad, dijalankan via `node nama-script.js`, dan mengambil screenshot di titik-titik penting (`page.screenshot(...)`) untuk diperiksa. Alur standarnya per fase:
1. Tulis script Playwright: buka `http://localhost:3000/login`, login dengan kredensial sesuai role yang dites, lalu navigasi ke menu terkait dan jalankan skenario (isi form, klik tombol, dst — sesuai kondisi normal & adversarial yang wajib dicoba di Aturan Umum).
2. Rekam: screenshot di tiap langkah penting, log `console` browser (tangkap event `page.on('console', ...)` untuk error JS), dan log response API relevan (`page.on('response', ...)`) untuk verifikasi status code & body — ini yang menggantikan "lihat langsung di browser" secara manual.
3. Jalankan script (`node ...`), lalu baca screenshot yang dihasilkan (pakai tool baca file/gambar) untuk verifikasi visual — jangan hanya percaya log teks tanpa cek visual untuk hal yang sifatnya tampilan (mis. pesan error muncul di UI, tombol hilang/muncul, layout struk).
4. Untuk kasus yang butuh akses API langsung tanpa lewat UI (poin keamanan di Aturan Umum), boleh pakai `curl`/PowerShell `Invoke-RestMethod` dengan token dari login, bukan lewat Playwright.
5. Kalau Playwright belum pernah dipakai di lingkungan kerja saat itu (folder scratchpad baru/kosong), install dulu (`npm init -y && npm install playwright && npx playwright install chromium`) sebelum mulai — ini satu kali saja per lingkungan kerja baru, bukan per fase.

---

## Aturan Umum (berlaku di SETIAP fase)

**Mindset wajib — bertindak sebagai QA/tester/bug bounty hunter, bukan developer yang mendemokan fitur:**
- Tugas testing **bukan** membuktikan fitur bisa jalan (happy path), tapi **berusaha aktif membuatnya gagal/rusak/tidak konsisten**. Kalau setelah dicoba dengan berbagai cara tetap tidak ketemu bug, itu baru hasil yang valid untuk fase tsb — bukan karena cuma dicoba jalur normalnya sekali lalu dianggap selesai.
- Perlakukan setiap input dari user sebagai **tidak bisa dipercaya**: apapun yang bisa diketik/dikirim, coba kirim yang paling "nakal", bukan cuma yang wajar.
- Perlakukan FE **hanya sebagai lapisan kenyamanan**, bukan sumber kebenaran validasi — semua aturan bisnis penting harus diverifikasi ulang tetap ditolak BE walau FE-nya dilewati (lewat curl/devtools/API langsung).
- Kalau suatu kondisi negatif "kebetulan" tidak bisa dicoba lewat UI (misal tombol disabled), itu **bukan alasan untuk skip** — coba tetap dari sisi API langsung, karena constraint UI bisa dilewati siapa pun yang tahu endpoint-nya.

**Kategori kondisi yang WAJIB dicoba di setiap fitur (bukan cuma yang relevan sekilas — evaluasi semua, catat yang di-skip beserta alasannya):**
- **Validasi & boundary:** field kosong/wajib tidak diisi, nilai batas (0, negatif, desimal aneh, angka sangat besar/overflow, teks sangat panjang/emoji/karakter aneh), tipe data salah dikirim langsung ke API (string ke field angka, dsb).
- **Duplikasi & konsistensi data:** SKU/barcode/nama/kode yang sama, dua entitas yang saling mereferensikan secara sirkular kalau modelnya memungkinkan.
- **Variasi data uji:** jangan pakai 1 pola data yang sama berulang-ulang (misal selalu "Produk Test 1" dengan harga bulat). Variasikan: kategori/satuan berbeda-beda, harga desimal vs bulat vs sangat murah/mahal, nama pendek vs panjang vs mengandung karakter khusus, qty kecil vs besar, kombinasi diskon/pajak berbeda, tanggal transaksi berbeda (hari ini, lintas bulan/tahun kalau relevan untuk laporan). Tujuannya supaya bug yang cuma muncul di kondisi data tertentu (mis. pembulatan desimal, sorting, format tanggal) ikut ketahuan — bukan cuma bug yang kelihatan di 1 skenario data yang itu-itu saja.
- **Konkurensi & timing:** double-click submit, klik cepat berulang pada tombol aksi (approve/void/bayar), dua tab/sesi browser berbeda memproses data yang sama bersamaan (race condition — misal dua kasir approve retur yang sama, dua device bayar PO yang sama), refresh/reload di tengah proses async.
- **State & alur yang tidak semestinya:** loncat langsung ke state akhir tanpa lewat state antara (misal bayar PO yang belum dibuat lunas dari status apapun), aksi pada data yang sudah di status final (approve 2x, void yang sudah void, hapus yang sudah dihapus), kombinasi status yang harusnya saling eksklusif.
- **Keamanan & otorisasi (perlakukan aplikasi seperti target bug bounty):** akses endpoint API langsung tanpa lewat UI, akses resource milik/scope role lain lewat manipulasi ID di URL/payload (IDOR — misal user kasir akses detail transaksi kasir lain atau produk yang bukan haknya), coba aksi yang UI-nya disembunyikan tapi endpoint BE-nya tetap hidup, uji role/permission berbeda (owner/admin/kasir) untuk tiap aksi CRUD, cek apakah token/session yang sudah logout/expired benar-benar ditolak BE (bukan cuma di-redirect FE), cek pesan error tidak membocorkan info sensitif (stack trace, query SQL, keberadaan akun).
- **Integritas lintas modul:** setelah suatu aksi (checkout, void, retur, bayar), verifikasi SEMUA sisi yang seharusnya berubah benar-benar konsisten (stok, mutasi stok, kas, piutang, laporan) — bukan cuma modul yang sedang difokuskan.
- **UX kegagalan:** matikan network sesaat/perlambat koneksi (throttle) saat submit penting, cek apakah UI stuck loading/blank atau kasih pesan jelas; cek console browser & tab Network setiap langkah — tidak boleh ada error JS atau request gagal yang didiamkan tanpa feedback ke user.

**Cara testing:**
- Semua pengujian dilakukan dengan menjalankan BE (`go run main.go`) dan FE (`npm run dev`) sungguhan, lalu men-drive browser asli (bukan cuma baca kode atau curl API sepihak) — curl/devtools dipakai sebagai **pelengkap** untuk kasus yang tidak bisa dipicu murni dari UI (lihat poin keamanan di atas), bukan pengganti testing browser.
- **Data uji (termasuk data dummy/adversarial) TIDAK perlu dihapus** — biarkan tetap ada di database development setelah dipakai, karena ini berguna juga bagi user (bisa dicek/dipakai ulang di luar sesi testing). Tidak ada langkah cleanup data di akhir kasus maupun akhir fase.
- Setelah selesai satu fase, cukup matikan proses BE/FE — tidak perlu membersihkan data.

**Cara memperbaiki bug (kalau ditemukan):**
1. **Analisis dulu** akar masalahnya sebelum ubah kode — baca kode terkait secara utuh, jangan asumsi dari gejala saja.
2. **Perbaikan seminimal mungkin** menyasar akar masalah — tidak menambah fitur/refactor di luar yang dibutuhkan untuk memperbaiki bug tsb.
3. **1 perbaikan → langsung 1 kali testing browser** untuk verifikasi perbaikan itu benar bekerja — bukan menumpuk banyak perbaikan dulu baru ditest di akhir.
4. Kalau perbaikan di satu file ternyata memicu error/warning baru di file lain (misalnya karena linter/compiler baru bisa menganalisis penuh setelah suatu pola diperbaiki), itu **tetap dibereskan sampai tuntas** di fase yang sama — dilaporkan dulu ke user kalau cakupannya jadi melebar signifikan.
5. Setelah semua perbaikan di fase itu selesai, jalankan:
   ```
   npm run type-check
   npm run lint
   npm run build
   ```
   Ketiganya harus **0 error dan 0 warning** sebelum fase dianggap selesai.
6. Laporan akhir tiap fase: daftar bug yang ditemukan (penyebab, cara perbaikan, bukti hasil testing browser setelah perbaikan) **dan** daftar kondisi negatif/adversarial yang sudah dicoba tapi ternyata aman (supaya kelihatan cakupan testingnya, bukan cuma yang gagal) — jangan laporkan fase "selesai tanpa bug" tanpa merinci apa saja yang sudah dicoba untuk membuatnya gagal.

---

## Daftar Fase

### Fase 0 — Baseline & Auth

**0.1 Health check & login**
```
Jalankan Fase 0.1 (Health check & login) dari docs\DEBUG_TESTING_PLAN.md :

Jalankan Fase 0.1: start BE dan FE, cek endpoint health BE merespon normal dan FE bisa diakses. Login via browser sebagai owner, admin, dan kasir (3 kali login terpisah) — pastikan masing-masing berhasil redirect ke dashboard tanpa error, dan sidebar menu yang muncul sesuai hak akses role masing-masing (bandingkan dengan seed data role_menu_access). Coba juga login dengan password salah dan username tidak terdaftar, pastikan pesan error jelas dan tidak bocor informasi (misal tidak bilang "user tidak ada" vs "password salah" secara eksplisit kalau itu bukan by design). Fix bug yang ditemukan sesuai Aturan Umum, lalu type-check+lint+build sampai bersih.
```

**0.2 Session & logout**
```
Jalankan Fase 0.2 (Session & logout) dari docs\DEBUG_TESTING_PLAN.md :

Jalankan Fase 0.2: test logout dari browser, pastikan token/session benar-benar invalid setelah logout (coba akses halaman butuh auth setelah logout, harus ke-redirect ke login). Test juga apa yang terjadi kalau token kedaluwarsa/expired (kalau bisa disimulasikan) — pastikan tidak stuck di halaman blank atau infinite loading. Fix bug yang ditemukan, lalu type-check+lint+build sampai bersih.
```

---

### Fase 1 — Produk & Inventori

**1.1 Produk — CRUD dasar**
```
Jalankan Fase 1.1 (Produk — CRUD dasar) dari docs\DEBUG_TESTING_PLAN.md :

Jalankan Fase 1.1: debug CRUD Produk. Test tambah produk (semua field wajib & opsional, termasuk kondisi field kosong dan barcode/SKU duplikat), edit produk, nonaktifkan/aktifkan produk, hapus produk (termasuk coba hapus produk yang masih punya stok/riwayat — pastikan ditolak dengan pesan jelas kalau memang ada aturan begitu). Test juga search & filter di list produk. Fix bug yang ditemukan sesuai Aturan Umum, lalu type-check+lint+build sampai bersih.
```

**1.2 Produk — gating stok berdasarkan role**
```
Jalankan Fase 1.2 (Produk — gating stok berdasarkan role) dari docs\DEBUG_TESTING_PLAN.md :

Jalankan Fase 1.2: regresi khusus fitur "field Stok hanya bisa diedit Admin" yang sudah dibuat sebelumnya. Login sebagai Admin (stok harus editable) dan Owner (stok harus readonly) di form Tambah & Edit Produk, verifikasi lewat browser. Juga coba bypass langsung ke API sebagai Owner mengirim nilai stok berbeda — pastikan BE tetap menolak/override. Fix bug yang ditemukan, lalu type-check+lint+build sampai bersih.
```

**1.3 Produk — import massal & generate kode**
```
Jalankan Fase 1.3 (Produk — import massal & generate kode) dari docs\DEBUG_TESTING_PLAN.md :

Jalankan Fase 1.3: debug fitur Import Produk (upload file, preview hasil parsing, submit import — termasuk file dengan baris invalid/duplikat untuk lihat penanganan errornya) dan fitur Generate Barcode/SKU otomatis saat tambah produk baru. Fix bug yang ditemukan, lalu type-check+lint+build sampai bersih.
```

**1.4 Produk — satuan/paket dinamis (graph qty/ref_qty) & tier harga**
```
Jalankan Fase 1.4 (Produk — satuan/paket dinamis (graph qty/ref_qty) & tier harga) dari docs\DEBUG_TESTING_PLAN.md :

Jalankan Fase 1.4: debug fitur satuan/paket produk dengan model dinamis saat ini (product_packages: unit_id, ref_package_id, qty, ref_qty — bukan lagi single "konversi ke basis"). Test di form Tambah Produk (mode draft, 1 form sekaligus): buat produk dengan beberapa paket berantai (mis. Slop -> Pack -> Batang) dalam satu kali submit, verifikasi resolved_factor tiap paket benar dan anchor (paket dasar) otomatis ter-set is_default. Test di form Edit Produk (mode live): tambah paket baru, edit rasio qty/ref_qty paket yang sudah ada (bukan anchor), coba hapus anchor (harus ditolak), coba hapus paket yang masih direferensikan paket lain (harus ditolak). Test juga tier harga per satuan tetap konsisten dipakai saat kasir/pembelian. Coba juga buat kombinasi rasio yang membentuk siklus (cycle) antar paket — harus ditolak validasi. Fix bug yang ditemukan, lalu type-check+lint+build sampai bersih.
```

**1.5 Kategori Produk**
```
Jalankan Fase 1.5 (Kategori Produk) dari docs\DEBUG_TESTING_PLAN.md :

Jalankan Fase 1.5: debug menu Kategori. CRUD lengkap, coba hapus kategori yang masih dipakai produk (pastikan ditangani dengan benar — ditolak atau produk ikut ter-uncategorized sesuai desain). Fix bug yang ditemukan, lalu type-check+lint+build sampai bersih.
```

**1.6 Unit/Satuan**
```
Jalankan Fase 1.6 (Unit/Satuan) dari docs\DEBUG_TESTING_PLAN.md :

Jalankan Fase 1.6: debug menu Unit. CRUD lengkap, coba hapus unit yang masih dipakai produk. Fix bug yang ditemukan, lalu type-check+lint+build sampai bersih.
```

---

### Fase 2 — Supplier

```
Jalankan Fase 2 (Supplier) dari docs\DEBUG_TESTING_PLAN.md :

Jalankan Fase 2: debug menu Supplier. Test tambah/edit supplier (field kosong, nomor telepon/email format aneh), toggle aktif/nonaktif, hapus supplier (termasuk coba hapus supplier yang masih punya riwayat pembelian — pastikan ditangani benar), search/filter di list. Fix bug yang ditemukan sesuai Aturan Umum, lalu type-check+lint+build sampai bersih.
```

---

### Fase 3 — Pengadaan: Pembelian Supplier

**3.1 Create**
```
Jalankan Fase 3.1 (Create) dari docs\DEBUG_TESTING_PLAN.md :

Jalankan Fase 3.1: debug create Pembelian Supplier. Test dengan 1 item, banyak item, produk yang sama dipilih 2x (harus ditolak validasi), diskon lebih besar dari subtotal (harus ditolak), submit tanpa pilih supplier/produk, double-klik tombol simpan (pastikan tidak submit 2x). Verifikasi stok & mutasi stok bertambah benar di DB setelah submit sukses. Fix bug yang ditemukan sesuai Aturan Umum, lalu type-check+lint+build sampai bersih.
```

**3.2 Edit & Delete**
```
Jalankan Fase 3.2 (Edit & Delete) dari docs\DEBUG_TESTING_PLAN.md :

Jalankan Fase 3.2: debug edit dan delete Pembelian Supplier. Pastikan tombol Edit/Delete cuma muncul untuk PO yang belum ada pembayaran. Test edit ubah qty/harga/item (verifikasi stok lama di-rollback dan stok baru benar), test delete (verifikasi stok & mutasi stok ikut dihapus/dikembalikan). Coba akses endpoint edit/delete langsung via API untuk PO yang sudah dibayar — pastikan ditolak backend juga. Fix bug yang ditemukan, lalu type-check+lint+build sampai bersih.
```

**3.3 Bayar (Pay)**
```
Jalankan Fase 3.3 (Bayar (Pay)) dari docs\DEBUG_TESTING_PLAN.md :

Jalankan Fase 3.3: debug fitur Bayar PO. Test bayar sebagian (verifikasi status jadi "Sebagian" & sisa hutang benar), bayar lunas via tombol "Bayar Lunas" maupun input manual, coba bayar lebih dari sisa hutang (harus ditolak), coba bayar PO yang sudah lunas (harus ditolak). Fix bug yang ditemukan, lalu type-check+lint+build sampai bersih.
```

**3.4 Void**
```
Jalankan Fase 3.4 (Void) dari docs\DEBUG_TESTING_PLAN.md :

Jalankan Fase 3.4: debug fitur Void PO. Test void PO yang belum dibayar (berhasil, stok balik, status jadi Dibatalkan), coba void PO yang sudah dibayar (harus ditolak), coba void 2x berturut-turut cepat (pastikan tidak dobel rollback stok), coba void PO yang sudah punya retur terkait (harus ditolak). Fix bug yang ditemukan, lalu type-check+lint+build sampai bersih.
```

**3.5 Tambah Item ke PO Lunas**
```
Jalankan Fase 3.5: debug fitur Tambah Item pada PO yang sudah ada pembayaran. Test tambah 1 dan banyak item baru, verifikasi item lama tidak berubah, total & status pembayaran ter-update otomatis dengan benar. Coba tambah item ke PO yang masih unpaid (harus diarahkan pakai Edit) dan ke PO yang sudah void (harus ditolak). Fix bug yang ditemukan, lalu type-check+lint+build sampai bersih.
```

**3.6 Filter, sort, pagination**
```
Jalankan Fase 3.6: debug filter tanggal, filter supplier, filter status pembayaran, sorting kolom, dan pagination di list Pembelian Supplier. Coba kombinasi filter yang menghasilkan data kosong (pastikan empty state tampil benar, bukan error/crash). Fix bug yang ditemukan, lalu type-check+lint+build sampai bersih.
```

---

### Fase 4 — Pengadaan: Retur Pembelian

**4.1 Create Retur**
```
Jalankan Fase 4.1: debug create Retur Pembelian. Pilih PO, centang item, ubah qty retur (test qty 0, qty negatif, qty melebihi qty pembelian asli — semua harus ditolak validasi), ganti PO yang dipilih di tengah proses (pastikan checklist ke-reset, bukan bawa data PO sebelumnya), submit tanpa isi alasan. Fix bug yang ditemukan sesuai Aturan Umum, lalu type-check+lint+build sampai bersih.
```

**4.2 Approve & Reject**
```
Jalankan Fase 4.2: debug approve dan reject Retur. Approve: verifikasi stok produk berkurang sesuai qty retur dan utang ke supplier (remaining_amount PO terkait) tersesuaikan otomatis. Reject: verifikasi reserved stock terlepas. Coba approve/reject retur yang sudah diproses sebelumnya (harus ditolak, tidak bisa diproses 2x). Fix bug yang ditemukan, lalu type-check+lint+build sampai bersih.
```

---

### Fase 5 — Pelanggan

```
Jalankan Fase 5: debug menu Pelanggan. CRUD lengkap termasuk field Limit Kredit (test nilai 0, negatif, sangat besar), search, hapus pelanggan yang masih punya piutang aktif (pastikan ditangani benar — ditolak atau diperbolehkan sesuai desain). Fix bug yang ditemukan sesuai Aturan Umum, lalu type-check+lint+build sampai bersih.
```

---

### Fase 6 — Kasir & Transaksi (jalur paling kritis)

**6.1 Buka/Tutup Kas**
```
Jalankan Fase 6.1: debug buka dan tutup kas dari halaman Kas Saya/Kasir. Test buka kas dengan saldo awal 0 dan nominal tertentu, coba buka kas 2x tanpa tutup dulu (harus ditolak), coba akses Kasir sebelum kas dibuka (tombol Bayar harus disabled dengan pesan jelas), tutup kas dan verifikasi ringkasan penjualan/pengeluaran hari itu benar. Fix bug yang ditemukan sesuai Aturan Umum, lalu type-check+lint+build sampai bersih.
```

**6.2 Cari & tambah produk ke keranjang**
```
Jalankan Fase 6.2: debug pencarian produk (nama, SKU, scan barcode manual input) dan penambahan ke keranjang di Kasir. Test tambah produk stok habis (harus ditolak atau diberi peringatan), ubah qty di keranjang (termasuk qty melebihi stok tersedia), hapus item dari keranjang, kosongkan keranjang. Fix bug yang ditemukan, lalu type-check+lint+build sampai bersih.
```

**6.3 Checkout — Tunai**
```
Jalankan Fase 6.3: debug checkout metode Tunai. Test bayar pas, bayar lebih (verifikasi kembalian dihitung benar), bayar kurang (harus ditolak submit), tombol nominal cepat (uang pas tersedia). Verifikasi setelah sukses: stok berkurang, saldo kas harian bertambah, struk muncul benar. Fix bug yang ditemukan sesuai Aturan Umum, lalu type-check+lint+build sampai bersih.
```

**6.4 Checkout — Non-tunai (Transfer/QRIS/Kartu)**
```
Jalankan Fase 6.4: debug checkout metode Transfer, QRIS, dan Kartu. Verifikasi tidak perlu input jumlah bayar/kembalian (atau sesuai desain masing-masing), transaksi tersimpan dengan payment_method yang benar, stok & kas ter-update sesuai. Fix bug yang ditemukan, lalu type-check+lint+build sampai bersih.
```

**6.5 Checkout — Kredit ke Pelanggan**
```
Jalankan Fase 6.5: debug checkout metode Kredit. Coba checkout kredit tanpa pilih pelanggan (harus ditolak), pilih pelanggan yang limit kreditnya sudah/hampir terlampaui (verifikasi validasi limit kredit), verifikasi setelah sukses otomatis muncul piutang baru di menu Piutang dengan nominal benar. Fix bug yang ditemukan, lalu type-check+lint+build sampai bersih.
```

**6.6 Struk & Preview**
```
Jalankan Fase 6.6: debug preview & cetak struk setelah transaksi (ukuran kertas 58mm/80mm, header/footer sesuai Pengaturan Printer, auto-print jika aktif). Fix bug yang ditemukan sesuai Aturan Umum, lalu type-check+lint+build sampai bersih.
```

**6.7 Riwayat Transaksi & Detail**
```
Jalankan Fase 6.7: debug halaman Transaksi (list riwayat). Test filter tanggal/kasir/metode pembayaran, buka detail transaksi (verifikasi item, harga, diskon, total cocok dengan yang tersimpan), pagination. Fix bug yang ditemukan, lalu type-check+lint+build sampai bersih.
```

**6.8 Void Transaksi**
```
Jalankan Fase 6.8: debug void transaksi dari halaman Transaksi. Verifikasi stok balik, saldo kas harian ikut dikurangi (untuk transaksi tunai), piutang terkait ikut jadi void (untuk transaksi kredit). Coba void transaksi yang sudah void (harus ditolak). Fix bug yang ditemukan sesuai Aturan Umum, lalu type-check+lint+build sampai bersih.
```

---

### Fase 7 — Keuangan: Kas Harian & Kas Saya

```
Jalankan Fase 7: debug menu Kas Harian (rekap semua kasir) dan Kas Saya (kasir yang login). Test lihat ringkasan per shift/tanggal, filter per kasir (Kas Harian), riwayat transaksi tunai/non-tunai dan pengeluaran dalam 1 sesi kas. Verifikasi angka summary konsisten dengan transaksi & pengeluaran dari fase sebelumnya. Fix bug yang ditemukan sesuai Aturan Umum, lalu type-check+lint+build sampai bersih.
```

---

### Fase 8 — Keuangan: Pengeluaran

```
Jalankan Fase 8: debug menu Pengeluaran. Test tambah pengeluaran (kategori, nominal 0/negatif harus ditolak), edit, hapus, verifikasi saldo kas harian berkurang sesuai nominal pengeluaran, coba tambah pengeluaran saat kas belum dibuka (harus ditolak). Fix bug yang ditemukan sesuai Aturan Umum, lalu type-check+lint+build sampai bersih.
```

---

### Fase 9 — Piutang

```
Jalankan Fase 9: debug menu Piutang. Verifikasi piutang dari transaksi kredit (Fase 6.5) muncul dengan data benar, test pelunasan piutang (sebagian & lunas), filter status (belum lunas/lunas), verifikasi piutang dari transaksi yang di-void (Fase 6.8) berubah status jadi void juga. Fix bug yang ditemukan sesuai Aturan Umum, lalu type-check+lint+build sampai bersih.
```

---

### Fase 10 — Dashboard & Dashboard Keuangan

```
Jalankan Fase 10: debug halaman Dashboard dan Dashboard Keuangan. Pastikan semua widget/grafik (statistik, tren penjualan, produk terlaris, kategori terlaris, metode pembayaran) load tanpa error di kondisi data kosong maupun ada data, dan angkanya konsisten dengan data transaksi/kas dari fase-fase sebelumnya. Fix bug yang ditemukan sesuai Aturan Umum, lalu type-check+lint+build sampai bersih.
```

---

### Fase 11 — Pelaporan

**11.1 Laporan Penjualan**
```
Jalankan Fase 11.1: debug Laporan Penjualan. Test filter rentang tanggal, filter kasir/produk kalau ada, export kalau ada fiturnya, verifikasi total cocok dengan data transaksi. Fix bug yang ditemukan sesuai Aturan Umum, lalu type-check+lint+build sampai bersih.
```

**11.2 Laporan Laba Rugi**
```
Jalankan Fase 11.2: debug Laporan Laba Rugi. Verifikasi perhitungan pendapatan, HPP, pengeluaran, dan laba bersih cocok dengan data transaksi & pengeluaran yang sudah diuji sebelumnya. Fix bug yang ditemukan sesuai Aturan Umum, lalu type-check+lint+build sampai bersih.
```

**11.3 Laporan Stok**
```
Jalankan Fase 11.3: debug Laporan Stok. Verifikasi mutasi stok (masuk dari pembelian, keluar dari transaksi, penyesuaian dari retur/void) tercatat dan angka stok akhir konsisten dengan yang ada di menu Produk. Fix bug yang ditemukan sesuai Aturan Umum, lalu type-check+lint+build sampai bersih.
```

**11.4 Kinerja Kasir**
```
Jalankan Fase 11.4: debug Laporan Kinerja Kasir. Test filter per kasir/tanggal, verifikasi jumlah transaksi & total penjualan per kasir cocok dengan data riil. Fix bug yang ditemukan sesuai Aturan Umum, lalu type-check+lint+build sampai bersih.
```

---

### Fase 12 — Shift

```
Jalankan Fase 12: debug menu Shift. Test CRUD shift (jam mulai/selesai, nama shift), assign shift ke user/kasir, verifikasi shift yang aktif muncul benar saat buka kas. Fix bug yang ditemukan sesuai Aturan Umum, lalu type-check+lint+build sampai bersih.
```

---

### Fase 13 — Sync Center

```
Jalankan Fase 13: debug menu Sync Center — modul integrasi paling kompleks (transaksi, pengeluaran, produk, kas). Test alur sinkronisasi normal, dan kalau ada mekanisme simulasi data pending/offline, test skenario konflik data (dua perubahan bersaing pada data yang sama) dan pastikan penyelesaian konfliknya sesuai desain, bukan data hilang/dobel. Fix bug yang ditemukan sesuai Aturan Umum, lalu type-check+lint+build sampai bersih.
```

---

### Fase 13b — PIN Kasir (verifikasi dulu sebelum di-debug penuh)

> Catatan: BE punya domain `pin` lengkap dengan endpoint `/pin/check|set|verify|change` (`BE/domain/pin`), tapi saat pengecekan terakhir **tidak ditemukan pemanggilannya di kode FE manapun**. Kemungkinan fitur ini belum diwire ke UI (rencana ke depan) atau memang sudah tidak dipakai. Jalankan sub-fase ini duluan sebagai investigasi, baru putuskan apakah perlu jadi fase debug penuh atau cukup dicatat sebagai dead code.

```
Jalankan Fase 13b: cek apakah fitur PIN (endpoint /pin/check, /pin/set, /pin/verify, /pin/change di BE/domain/pin) benar-benar dipakai di suatu alur UI (misal ganti kasir cepat, konfirmasi void/diskon, dsb) dengan grep menyeluruh di FE dan cek juga halaman-halaman yang mungkin memicunya lewat network request di browser. Kalau ternyata dipakai, laporkan di menu/alur mana dan baru susulkan skenario testingnya. Kalau ternyata tidak dipakai sama sekali, laporkan sebagai temuan (bukan fix) untuk didiskusikan — apakah mau diwire ke FE, atau dihapus sebagai dead code — jangan hapus/ubah apapun tanpa konfirmasi eksplisit.
```

---

### Fase 14 — Sistem: User, Role, Akses

**14.1 Manajemen User**
```
Jalankan Fase 14.1: debug menu Manajemen User. CRUD user (username duplikat harus ditolak, password lemah kalau ada validasi), assign role, nonaktifkan user, coba login pakai user yang baru dinonaktifkan (harus ditolak). Fix bug yang ditemukan sesuai Aturan Umum, lalu type-check+lint+build sampai bersih.
```

**14.2 Manajemen Role**
```
Jalankan Fase 14.2: debug menu Manajemen Role. CRUD role, coba hapus role yang masih dipakai user (harus ditolak). Fix bug yang ditemukan sesuai Aturan Umum, lalu type-check+lint+build sampai bersih.
```

**14.3 Role Access**
```
Jalankan Fase 14.3: debug menu Role Access. Ubah permission (can_view/can_create/can_edit/can_delete) suatu role untuk beberapa menu, logout-login ulang sebagai user dengan role itu, verifikasi perubahan permission langsung berefek ke sidebar dan tombol aksi yang muncul di halaman terkait. Fix bug yang ditemukan sesuai Aturan Umum, lalu type-check+lint+build sampai bersih.
```

**14.4 Manajemen Menu**
```
Jalankan Fase 14.4: debug menu Manajemen Menu. Tambah menu baru (pilih parent, path dari route registry, icon), edit, hapus, ubah urutan. Verifikasi menu baru bisa langsung di-assign permission-nya di Role Access dan muncul di sidebar untuk role yang diberi akses. Fix bug yang ditemukan sesuai Aturan Umum, lalu type-check+lint+build sampai bersih.
```

---

### Fase 15 — Sistem: Printer, Profil Toko, Versi, Backup

**15.1 Pengaturan Printer**
```
Jalankan Fase 15.1: debug Pengaturan Printer. Ubah ukuran kertas, header/footer, toggle logo & auto-print, verifikasi preview struk live-update dan tersimpan benar setelah reload halaman. Fix bug yang ditemukan sesuai Aturan Umum, lalu type-check+lint+build sampai bersih.
```

**15.2 Profil Toko**
```
Jalankan Fase 15.2: debug Profil Toko. Edit nama toko, alamat, kontak, logo. Verifikasi perubahan tercermin di tempat lain yang memakainya (misal struk, header aplikasi jika ada). Fix bug yang ditemukan sesuai Aturan Umum, lalu type-check+lint+build sampai bersih.
```

**15.3 Versi Aplikasi**
```
Jalankan Fase 15.3: debug menu Versi Aplikasi. Tambah versi baru, toggle wajib update, verifikasi versi terbaru yang tampil benar. Fix bug yang ditemukan sesuai Aturan Umum, lalu type-check+lint+build sampai bersih.
```

**15.4 Backup & Restore**
```
Jalankan Fase 15.4: debug menu Backup & Restore. Jalankan backup manual, verifikasi file backup ter-generate benar. HATI-HATI dengan fitur restore — kalau mau ditest, konfirmasi dulu ke user sebelum benar-benar menjalankan restore karena berpotensi menimpa data development yang sedang dipakai untuk testing fase-fase lain. Fix bug yang ditemukan sesuai Aturan Umum, lalu type-check+lint+build sampai bersih.
```

---

### Fase 16 — Cross-check Permission/RBAC

```
Jalankan Fase 16: debug permission lintas role secara menyeluruh. Untuk role Owner, Admin, dan Kasir: cek satu-satu menu apa saja yang muncul di sidebar sesuai role_menu_access, lalu untuk tiap menu yang TIDAK seharusnya bisa diakses role tsb, coba akses langsung via URL (bukan cuma sidebar disembunyikan) dan via API langsung — pastikan backend menolak dengan status 403/pesan jelas, bukan cuma disembunyikan di UI. Fix bug yang ditemukan sesuai Aturan Umum, lalu type-check+lint+build sampai bersih.
```

---

### Fase 17 — Regresi Akhir

```
Jalankan Fase 17: regresi akhir menyeluruh. Jalankan ulang secara singkat (smoke test) skenario inti tiap modul dari Fase 1 sampai 16 untuk memastikan tidak ada yang rusak akibat perbaikan-perbaikan di fase sebelumnya. Jalankan type-check, lint, dan build final — harus 0 error 0 warning. Susun laporan ringkasan: daftar semua bug yang ditemukan sepanjang seluruh fase, penyebabnya, dan cara perbaikannya.
```

---

## Catatan Penggunaan

- Fase tidak harus dijalankan berurutan dalam satu sesi — bisa dicicil kapan saja, tapi **urutan antar fase tetap disarankan mengikuti nomor di atas** karena ada dependency data (misal Fase 6 butuh produk dari Fase 1 & supplier dari Fase 2 sudah beres).
- Kalau satu fase menemukan banyak bug dan terasa kepanjangan, boleh dihentikan di tengah dan dilanjutkan nanti — cukup sebutkan fase & sub-bagian mana yang terakhir dikerjakan.
- Data uji yang dibuat selama testing (produk/PO/transaksi dummy, termasuk yang adversarial) **dibiarkan saja di database development, tidak perlu dihapus** — ini berguna juga bagi user untuk dicek/dipakai kembali di luar sesi testing.
- Fase 6.6 & 15.1 (struk/Pengaturan Printer): bug "struk asli tidak memakai Pengaturan Printer" sudah diperbaiki di luar rencana ini (ReceiptPrint.tsx sekarang konsumsi header/footer/paper_size/logo/auto_print, plus default BE untuk paper_size & footer saat setting kosong). Saat fase ini dijalankan, perlakukan sebagai **regresi** — pastikan perbaikan itu masih benar, bukan menemukan ulang dari nol.

# Rencana Testing Debug via Browser — Bertahap per Fase (AI-driven)

Beda dengan `docs/audit/` (audit statis baca-kode/code review), dokumen ini untuk
**testing fungsional dinamis yang dijalankan sepenuhnya oleh AI**: AI yang menjalankan
BE+FE, membuka browser lewat Playwright, mengisi form, klik tombol, dan mengamati hasilnya
sendiri (screenshot, console error, response API) — **bukan testing manual oleh Anda**.
Peran Anda di tiap fase: memicu prompt-nya, meninjau laporan bug yang ditemukan, dan
mengambil keputusan untuk temuan besar yang butuh diskusi. Ini sesuai kesepakatan kita
sebelumnya: testing debug lewat browser dikerjakan AI, bukan klik manual Anda.

Supaya hasilnya maksimal, tiap fase dipecah jadi **sub-fase kecil per fitur/menu**, bukan
satu fase besar mencakup banyak menu sekaligus — dengan sub-fase yang sempit, AI bisa lebih
teliti menguji tiap sudut satu fitur (form, validasi, edge case) daripada terburu-buru
menyapu banyak menu dalam satu putaran.

Jalankan **satu sub-fase per sesi**, tunggu bug ditemukan+diperbaiki, verifikasi, baru lanjut
sub-fase berikutnya. Urutannya mengikuti ketergantungan data: data master dulu, baru
transaksi yang butuh data master itu, baru laporan yang butuh transaksi.

## Cara pakai

1. Buka sesi baru (atau lanjutkan sesi ini).
2. Copy salah satu prompt sub-fase di bawah, tempel apa adanya.
3. Setelah sub-fase selesai dan Anda puas dengan hasilnya, lanjut ke prompt sub-fase
   berikutnya.
4. Kalau ada temuan besar yang perlu keputusan desain (bukan sekadar bug kecil), Claude akan
   berhenti dan diskusi dulu sebelum implementasi — sama seperti pola kerja kita sejauh ini.

## Persiapan Lingkungan Teknis (rujukan semua sub-fase)

Ini persis langkah yang sudah terbukti jalan di sesi-sesi sebelumnya. AI yang menjalankan
tiap sub-fase harus ikuti pola ini, bukan berimprovisasi:

**1. Jalankan Backend (Go)**
```bash
cd d:/Develop/Project_pos_mahenz/BE
go build -o /tmp/pos_be_test.exe .
# pastikan MySQL sudah jalan di 127.0.0.1:3306, database pos_retail_db (lihat BE/config/config_dev.json)
(/tmp/pos_be_test.exe > /tmp/be_test.log 2>&1 &)
sleep 3
curl -s http://localhost:8080/api/health   # harus balas {"code":"00",...,"status":true}
```

**2. Jalankan Frontend (Vite)**
```bash
cd d:/Develop/Project_pos_mahenz/FE
(npm run dev > /tmp/fe_test.log 2>&1 &)
sleep 6   # tunggu Vite selesai "ready in ...ms" di /tmp/fe_test.log, biasanya jalan di :3000
```

**3. Siapkan Playwright (tidak terpasang otomatis, harus di-install manual tiap sesi baru)**

Gunakan scratchpad directory (bukan folder proyek), contoh path sesi:
`C:\Users\brama\AppData\Local\Temp\claude\...\scratchpad`

```bash
cd <scratchpad-dir>
npm init -y
npm install playwright
npx playwright install chromium
```

**4. Pola script login (Node + Playwright)** — dipakai berulang di semua sub-fase, cukup
ganti kredensial dan langkah setelah login:
```js
const { chromium } = require('playwright');
(async () => {
  const browser = await chromium.launch();
  const page = await browser.newPage();
  const consoleErrors = [];
  const httpErrors = [];
  page.on('console', m => { if (m.type() === 'error') consoleErrors.push(m.text()); });
  page.on('pageerror', e => consoleErrors.push('PAGEERROR: ' + e.message));
  page.on('response', res => { if (res.status() >= 400 && res.url().includes('/api/')) httpErrors.push(res.status() + ' ' + res.url()); });

  await page.goto('http://localhost:3000/login', { waitUntil: 'networkidle' });
  await page.getByPlaceholder('Masukkan username').fill('owner');   // ganti sesuai role yang diuji
  await page.getByPlaceholder('Masukkan password').fill('owner123');
  await page.getByRole('button', { name: 'Masuk' }).click();
  await page.waitForURL('**/dashboard');

  // ... navigasi ke halaman yang diuji, isi form, screenshot, dst ...
  await page.screenshot({ path: __dirname + '/nama_screenshot.png', fullPage: true });

  console.log('console/page errors:', JSON.stringify(consoleErrors));
  console.log('http 4xx/5xx:', JSON.stringify(httpErrors));
  await browser.close();
})();
```
Jalankan dengan `node nama_file.js`. Simpan tiap screenshot dengan nama yang jelas per
langkah supaya mudah ditinjau ulang.

**5. Verifikasi langsung ke API (bukan cuma tampilan FE)** — kalau perlu mengecek data
mentah di database tanpa lewat UI (mis. cek apakah subtotal transaksi benar-benar dihitung
ulang server), pakai curl langsung dengan token dari login:
```bash
TOKEN=$(curl -s -X POST http://localhost:8080/api/auth/login -H "Content-Type: application/json" \
  -d '{"username":"owner","password":"owner123"}' \
  | node -e "let d='';process.stdin.on('data',c=>d+=c);process.stdin.on('end',()=>{const j=JSON.parse(d);console.log(j.data.access_token||j.data.token)})")
curl -s -X POST http://localhost:8080/api/<endpoint> -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" -d '{...}'
```

**6. Setelah selesai testing**, matikan proses BE/FE yang dijalankan di background supaya
tidak menumpuk port terpakai di sesi berikutnya (cari proses dengan `Get-NetTCPConnection`
di PowerShell, atau `pkill` kalau di bash, sesuai port 8080/3000).

## Kredensial yang tersedia

- `owner` / `owner123` — akses penuh (dari seed data resmi `002_seed_data.sql`)
- `admin` / `admin123` — akses manajemen tanpa pengaturan sistem (dari seed data resmi)
- Role **kasir**: belum ada kredensial resmi dari seed data migrasi. User `kasir1`/`kasir2`/
  `kasir3` yang mungkin terlihat di halaman Manajemen User itu dibuat manual di sesi
  dev/testing sebelumnya (bukan dari migration file), sehingga passwordnya **tidak
  diketahui**. Di Sub-fase 1.2, buat/reset user kasir baru sendiri lewat halaman Manajemen
  User supaya ada kredensial kasir yang pasti untuk dipakai di sub-fase berikutnya.

## Catatan silang-referensi

Beberapa bug yang mungkin muncul saat testing browser di area Auth/Sistem dan Sales sudah
pernah ditemukan lewat audit statis di `docs/audit/09-auth.md` dan `docs/audit/06-sales.md`
(mis. refresh token tidak jalan, checkout tidak menghitung ulang uang di server). Kalau
Claude menemukan gejala yang cocok saat testing, sebaiknya baca file itu dulu supaya tidak
menganalisis dari nol.

## Template pembuka tiap prompt (sudah disisipkan otomatis di tiap sub-fase di bawah)

Kalimat pembuka standar tiap prompt sudah menegaskan: **AI yang menjalankan seluruh
testing** (build+run BE/FE, kendalikan browser lewat Playwright, isi form, verifikasi
API) — user tidak melakukan klik manual, hanya meninjau hasil dan memutuskan untuk temuan
besar. Ini konsisten di semua 24 sub-fase di bawah.

---

# FASE 1 — Fondasi: Auth & Sistem

## Sub-fase 1.1 — Login, Logout, Sesi & Refresh Token

```
Jalankan testing debug fungsional secara otomatis (Anda yang mengeksekusi lewat browser
Playwright, saya tidak testing manual — sesuai kesepakatan kita) untuk alur Auth di aplikasi
POS ini (d:\Develop\Project_pos_mahenz). Ikuti persis bagian "Persiapan Lingkungan Teknis"
di docs/browser-testing-plan.md untuk menjalankan BE+FE dan menyiapkan Playwright.

Cakupan sub-fase ini (fokus HANYA login/logout/sesi, jangan melebar ke menu lain dulu):
1. Login dengan owner/owner123 — pastikan redirect ke /dashboard, token tersimpan, menu
   sidebar yang tampil sesuai izin role owner (harus lengkap semua menu termasuk sistem.users).
2. Login dengan admin/admin123 — pastikan menu sistem.users TIDAK muncul (sesuai desain
   RBAC yang sudah dikonfirmasi sebelumnya: admin tidak boleh akses manajemen user).
3. Logout dari masing-masing role — pastikan token benar-benar dihapus (coba akses ulang
   halaman yang butuh auth setelah logout, harus terlempar ke /login, dan request API lama
   dengan token itu harus ditolak 401 kalau dicoba ulang via curl).
4. Coba login dengan kredensial salah (password salah, username tidak ada) — pastikan
   pesan error jelas dan tidak membocorkan info (misal jangan beda pesan antara "user tidak
   ada" vs "password salah").
5. Amati mekanisme refresh token: biarkan sesi berjalan, lihat apakah ada pemanggilan
   /api/auth/refresh otomatis dari FE sebelum token kedaluwarsa, dan apakah access_token
   baru benar-benar dipakai (bukan `undefined`). Area ini sudah dicurigai bermasalah di
   docs/audit/09-auth.md temuan #1 dan #2 — baca file itu dulu, lalu verifikasi apakah
   gejala yang dijelaskan di sana benar-benar terjadi di browser sungguhan.

Cek console error dan response API 4xx/5xx yang tidak wajar di sepanjang proses. Di akhir,
berikan daftar bug dengan langkah reproduksi singkat. Perbaiki bug kecil/jelas langsung.
Untuk temuan soal keamanan sesi/token (bukan sekadar UI), berhenti dan diskusikan dulu
dengan saya sebelum implementasi.
```

## Sub-fase 1.2 — Manajemen User

```
Lanjut testing debug fungsional otomatis (Anda yang mengeksekusi lewat Playwright, bukan
saya manual) untuk halaman Manajemen User (/settings/users) di aplikasi POS ini
(d:\Develop\Project_pos_mahenz). Ikuti "Persiapan Lingkungan Teknis" di
docs/browser-testing-plan.md untuk menjalankan BE+FE. Login sebagai owner/owner123 (hanya
owner yang boleh akses halaman ini).

Cakupan sub-fase ini (fokus HANYA menu ini):
1. Tambah user baru dengan role "kasir" (ini akan dipakai sebagai kredensial kasir resmi
   untuk sub-fase-sub-fase berikutnya yang butuh login sebagai kasir — catat username/
   password yang dipakai di laporan akhir Anda).
2. Edit user yang baru dibuat (ubah nama, role).
3. Ganti password user tsb, lalu verifikasi bisa login dengan password baru.
4. Toggle status aktif/nonaktif user tsb, verifikasi user nonaktif tidak bisa login lagi.
5. Coba skenario ganjil:
   - Ubah role diri sendiri (sedang login sebagai owner, coba ubah role owner sendiri jadi
     kasir) — apakah ada guard yang mencegah, atau malah lolos dan owner kehilangan akses.
   - Nonaktifkan/hapus admin terakhir yang tersisa (kalau cuma ada 1 admin) — apakah ada
     proteksi "last admin", atau malah bisa membuat sistem tidak punya admin sama sekali.
   - Coba hapus user yang masih dipakai sebagai kasir di suatu shift/transaksi (kalau sudah
     ada data transaksi dari testing sebelumnya) — apakah ada guard integritas data.
6. Coba input tidak valid: username kosong/duplikat, password terlalu pendek, karakter aneh
   di nama — cek konsistensi validasi FE vs BE.

Cek console error dan response API 4xx/5xx tidak wajar. Di akhir, berikan daftar bug dengan
langkah reproduksi, dan SEBUTKAN username/password kasir yang baru dibuat di poin 1 supaya
bisa dipakai di sub-fase lain. Perbaiki bug kecil/jelas langsung. Untuk temuan privilege
escalation atau proteksi last-admin yang hilang, berhenti dan diskusikan dulu dengan saya.
```

## Sub-fase 1.3 — Manajemen Role & Hak Akses

```
Lanjut testing debug fungsional otomatis (Anda yang mengeksekusi lewat Playwright) untuk
halaman Manajemen Role (/settings/roles) di aplikasi POS ini (d:\Develop\Project_pos_mahenz).
Ikuti "Persiapan Lingkungan Teknis" di docs/browser-testing-plan.md. Login sebagai
owner/owner123.

Cakupan sub-fase ini (fokus HANYA menu ini):
1. Lihat daftar role (owner/admin/kasir), lihat detail konfigurasi masing-masing.
2. Buka halaman atur akses menu per role (/settings/roles/:id/access) untuk role "kasir",
   coba ubah izin akses (misal cabut akses ke satu menu, tambahkan akses ke menu lain),
   simpan.
3. Verifikasi perubahan benar-benar berlaku: login sebagai user kasir (pakai kredensial dari
   Sub-fase 1.2), cek menu sidebar berubah sesuai konfigurasi baru — bukan cuma tersimpan di
   database tapi tidak berefek ke FE.
4. Coba skenario ganjil: hapus/ubah role sistem (owner/admin/kasir yang is_system=true) —
   pastikan ada guard yang menolak (role sistem seharusnya tidak bisa dihapus/diubah nama
   dasarnya, hanya izin aksesnya yang bisa diatur).
5. Coba beri role "kasir" akses ke menu sistem.users lewat halaman ini, simpan, lalu login
   sebagai kasir — apakah dia benar-benar bisa akses Manajemen User sekarang (menguji apakah
   role access assignment ini konsisten dengan guard hardcoded "owner only" yang mungkin ada
   di kode, atau malah bertentangan/membingungkan).

Cek console error dan response API 4xx/5xx tidak wajar. Di akhir, berikan daftar bug dengan
langkah reproduksi. Perbaiki bug kecil/jelas langsung. Untuk temuan soal proteksi role
sistem atau konflik antara role-access-assignment vs guard hardcoded, berhenti dan
diskusikan dulu dengan saya.
```

## Sub-fase 1.4 — Manajemen Menu

```
Lanjut testing debug fungsional otomatis (Anda yang mengeksekusi lewat Playwright) untuk
halaman Manajemen Menu (/settings/menus) di aplikasi POS ini (d:\Develop\Project_pos_mahenz).
Ikuti "Persiapan Lingkungan Teknis" di docs/browser-testing-plan.md. Login sebagai
owner/owner123.

Cakupan sub-fase ini (fokus HANYA menu ini):
1. Lihat daftar menu (36 data per temuan sebelumnya), coba fitur reorder/urutan menu (drag
   atau tombol naik-turun, tergantung UI-nya), simpan, refresh halaman, pastikan urutan
   tersimpan permanen dan sidebar benar-benar berubah urutannya.
2. Toggle status aktif/nonaktif satu menu (pilih menu yang bukan menu inti seperti dashboard),
   verifikasi menu itu hilang dari sidebar semua role yang tadinya punya akses ke situ.
3. Coba nonaktifkan menu "beranda.dashboard" (menu utama) — apakah ada guard yang mencegah
   menonaktifkan menu esensial, atau malah bisa membuat semua orang kehilangan akses
   dashboard.
4. Coba edit label/path menu — pastikan tidak merusak routing FE (path harus tetap valid,
   dicocokkan dengan route yang benar-benar ada di FE/src/app/router.tsx).

Cek console error dan response API 4xx/5xx tidak wajar. Di akhir, berikan daftar bug dengan
langkah reproduksi. Perbaiki bug kecil/jelas langsung. Untuk temuan yang bisa membuat menu
esensial hilang tanpa guard, berhenti dan diskusikan dulu dengan saya.
```

## Sub-fase 1.5 — Profil Toko

```
Lanjut testing debug fungsional otomatis (Anda yang mengeksekusi lewat Playwright) untuk
halaman Profil Toko (/settings/store) di aplikasi POS ini (d:\Develop\Project_pos_mahenz).
Ikuti "Persiapan Lingkungan Teknis" di docs/browser-testing-plan.md. Login sebagai
owner/owner123 atau admin/admin123 (cek siapa saja yang punya akses tulis ke halaman ini).

Cakupan sub-fase ini (fokus HANYA menu ini):
1. Edit nama toko, alamat, telepon, dan field lain yang ada di form ini, simpan.
2. Refresh halaman (full reload, bukan navigasi client-side) — pastikan data yang tadi
   disimpan benar-benar persist dari database, bukan cuma state FE yang hilang saat reload.
3. Coba input tidak valid (nama kosong, format telepon aneh, teks sangat panjang) — cek
   konsistensi validasi FE vs BE.
4. Cek apakah perubahan profil toko ini muncul di tempat lain yang memakainya (misal struk
   kasir/nota, kalau ada preview cetak).

Cek console error dan response API 4xx/5xx tidak wajar. Di akhir, berikan daftar bug dengan
langkah reproduksi, perbaiki yang jelas/kecil langsung, diskusikan dulu untuk temuan besar.
```

---

# FASE 2 — Data Master Produk

## Sub-fase 2.1 — Kategori Produk

```
Lanjut testing debug fungsional otomatis (Anda yang mengeksekusi lewat Playwright, bukan
saya manual) untuk halaman Kategori Produk (/products/categories) di aplikasi POS ini
(d:\Develop\Project_pos_mahenz). Ikuti "Persiapan Lingkungan Teknis" di
docs/browser-testing-plan.md. Login sebagai admin/admin123.

Cakupan sub-fase ini (fokus HANYA menu ini):
1. Tambah kategori baru, edit nama/deskripsi, toggle aktif/nonaktif, lihat jumlah produk
   per kategori ter-update dengan benar.
2. Coba nonaktifkan/hapus kategori yang masih dipakai produk aktif — pastikan ada guard
   yang jelas (pesan error yang masuk akal), bukan error 500 mentah atau produk jadi
   "yatim" (kategori hilang tapi produk masih menunjuk ke situ).
3. Coba input tidak valid: nama kosong, nama duplikat, kode kategori duplikat, teks sangat
   panjang — cek konsistensi validasi FE vs BE.
4. Cek filter/search kategori (aktif/nonaktif/semua, pencarian nama) berfungsi benar.

Cek console error dan response API 4xx/5xx tidak wajar. Di akhir, berikan daftar bug dengan
langkah reproduksi, perbaiki yang jelas/kecil langsung, diskusikan dulu untuk temuan besar.
```

## Sub-fase 2.2 — Satuan/Unit Produk

```
Lanjut testing debug fungsional otomatis (Anda yang mengeksekusi lewat Playwright) untuk
halaman Satuan Produk (/products/units) di aplikasi POS ini (d:\Develop\Project_pos_mahenz).
Ikuti "Persiapan Lingkungan Teknis" di docs/browser-testing-plan.md. Login sebagai
admin/admin123.

Cakupan sub-fase ini (fokus HANYA menu ini):
1. Tambah satuan baru (nama + singkatan), edit, toggle aktif/nonaktif.
2. Coba nonaktifkan/hapus satuan yang masih dipakai produk aktif — pastikan ada guard yang
   jelas, bukan error 500 mentah.
3. Coba input tidak valid: nama/singkatan kosong atau duplikat.

Cek console error dan response API 4xx/5xx tidak wajar. Di akhir, berikan daftar bug dengan
langkah reproduksi, perbaiki yang jelas/kecil langsung, diskusikan dulu untuk temuan besar.
```

## Sub-fase 2.3 — Produk: CRUD Dasar & Generate Barcode/SKU

```
Lanjut testing debug fungsional otomatis (Anda yang mengeksekusi lewat Playwright) untuk
halaman Produk (/products) di aplikasi POS ini (d:\Develop\Project_pos_mahenz). Ikuti
"Persiapan Lingkungan Teknis" di docs/browser-testing-plan.md. Login sebagai admin/admin123.
Pastikan sudah ada minimal 1 kategori dan 1 satuan aktif dari Sub-fase 2.1/2.2 sebelum mulai.

Cakupan sub-fase ini (fokus HANYA CRUD dasar produk, JANGAN dulu ke import/paket/multi-harga
— itu di sub-fase terpisah):
1. Tambah produk baru lengkap: nama, kategori, satuan, harga beli, harga jual, stok awal.
   Coba fitur generate barcode otomatis dan generate SKU otomatis di form — klik berkali-kali
   dan pastikan tidak ada nilai duplikat yang ter-generate.
2. Edit produk (ubah harga, stok, kategori), toggle status aktif/nonaktif, lihat detail,
   coba fitur cetak (print) kalau ada tombolnya.
3. Hapus produk — coba hapus produk yang belum pernah ada transaksi (harus berhasil bersih)
   vs produk yang sudah pernah dipakai transaksi/pembelian dari testing sebelumnya (kalau
   ada) — pastikan ada guard yang jelas, bukan data transaksi lama jadi rusak/yatim.
4. Coba input tidak valid: harga beli/jual negatif, harga jual lebih murah dari harga beli
   (apakah ada warning margin negatif), stok negatif, nama kosong/sangat panjang.

Cek console error dan response API 4xx/5xx tidak wajar. Di akhir, berikan daftar bug dengan
langkah reproduksi, perbaiki yang jelas/kecil langsung, diskusikan dulu untuk temuan besar.
```

## Sub-fase 2.4 — Produk: Paket, Multi-Harga, Search & Barcode Lookup

```
Lanjut testing debug fungsional otomatis (Anda yang mengeksekusi lewat Playwright) untuk
fitur lanjutan halaman Produk di aplikasi POS ini (d:\Develop\Project_pos_mahenz). Ikuti
"Persiapan Lingkungan Teknis" di docs/browser-testing-plan.md. Login sebagai admin/admin123.
Pakai produk yang sudah dibuat di Sub-fase 2.3.

Cakupan sub-fase ini (fokus HANYA fitur ini):
1. Kelola paket produk (packages) di detail produk kalau ada UI-nya — tambah beberapa
   varian paket (misal isi per dus/pack), edit, hapus paket, pastikan konversi
   stok/harga antar paket masuk akal.
2. Kelola multi-harga (prices) di detail produk kalau ada UI-nya — tambah harga untuk
   tingkat pelanggan berbeda (grosir/eceran, dsb kalau ada), edit, hapus.
3. Search produk lewat kotak pencarian (nama parsial, barcode parsial) di halaman /products
   dan/atau di Kasir — pastikan hasil relevan dan cepat.
4. Cari produk by barcode exact match (endpoint by-barcode) — coba barcode yang ada dan
   yang tidak ada, pastikan pesan "tidak ditemukan" jelas bukan error mentah.

Cek console error dan response API 4xx/5xx tidak wajar. Di akhir, berikan daftar bug dengan
langkah reproduksi, perbaiki yang jelas/kecil langsung, diskusikan dulu untuk temuan besar.
```

## Sub-fase 2.5 — Import Produk (CSV/Excel)

```
Lanjut testing debug fungsional otomatis (Anda yang mengeksekusi lewat Playwright) untuk
fitur Import Produk di halaman /products di aplikasi POS ini (d:\Develop\Project_pos_mahenz).
Ikuti "Persiapan Lingkungan Teknis" di docs/browser-testing-plan.md. Login sebagai
admin/admin123.

Cakupan sub-fase ini (fokus HANYA fitur import):
1. Download template import, lihat formatnya.
2. Buat file import dengan semua baris valid, upload, cek preview (import-preview) sesuai,
   lalu proses import-bulk, verifikasi semua produk masuk dengan benar ke halaman /products
   dan ke database.
3. Buat file import dengan CAMPURAN baris valid dan baris error (misal harga bukan angka,
   kategori tidak ada, kolom wajib kosong) — cek preview memisahkan mana yang valid/error
   dengan jelas, dan pastikan proses bulk HANYA memasukkan baris valid (baris error tidak
   ikut masuk maupun merusak baris valid di sekitarnya).
4. Coba upload file kosong, file format salah (bukan csv/xlsx), file sangat besar (banyak
   baris) — lihat penanganannya.

Cek console error dan response API 4xx/5xx tidak wajar. Di akhir, berikan daftar bug dengan
langkah reproduksi, perbaiki yang jelas/kecil langsung, diskusikan dulu untuk temuan besar
(terutama kalau ada risiko data korup akibat baris error ikut ter-proses).
```

---

# FASE 3 — Pengadaan

## Sub-fase 3.1 — Supplier

```
Lanjut testing debug fungsional otomatis (Anda yang mengeksekusi lewat Playwright) untuk
halaman Supplier (/suppliers) di aplikasi POS ini (d:\Develop\Project_pos_mahenz). Ikuti
"Persiapan Lingkungan Teknis" di docs/browser-testing-plan.md. Login sebagai admin/admin123.

Cakupan sub-fase ini (fokus HANYA menu ini):
1. Tambah supplier baru (nama, kontak, telepon, alamat), edit, toggle aktif/nonaktif, lihat
   detail.
2. Coba nonaktifkan/hapus supplier yang nanti (di Sub-fase 3.2) sudah punya riwayat
   pembelian — kalau sub-fase ini dikerjakan sebelum ada riwayat, catat dan tandai untuk
   dicoba ulang setelah Sub-fase 3.2 selesai.
3. Coba input tidak valid: nama kosong/duplikat, format telepon aneh.

Cek console error dan response API 4xx/5xx tidak wajar. Di akhir, berikan daftar bug dengan
langkah reproduksi, perbaiki yang jelas/kecil langsung, diskusikan dulu untuk temuan besar.
```

## Sub-fase 3.2 — Pembelian Supplier (Create & Edit)

```
Lanjut testing debug fungsional otomatis (Anda yang mengeksekusi lewat Playwright) untuk
halaman Pembelian Supplier (/suppliers/purchases) di aplikasi POS ini
(d:\Develop\Project_pos_mahenz). Ikuti "Persiapan Lingkungan Teknis" di
docs/browser-testing-plan.md. Login sebagai admin/admin123. Pakai supplier dari Sub-fase 3.1
dan produk dari Fase 2.

Cakupan sub-fase ini (fokus HANYA membuat & edit pembelian, JANGAN dulu ke pembayaran/retur
— itu di sub-fase terpisah):
1. Buat pembelian baru dengan beberapa item produk berbeda dan quantity berbeda-beda, coba
   fitur generate kode PO otomatis (klik berkali-kali, pastikan tidak duplikat).
2. Verifikasi ke BE langsung (via API): setelah pembelian dibuat, apakah stok produk terkait
   BENAR-BENAR bertambah sesuai quantity yang dibeli (bandingkan stok sebelum vs sesudah).
3. Edit pembelian yang baru dibuat (ubah quantity/item) — verifikasi stok ikut terkoreksi
   dengan benar (bukan malah dobel nambah atau tidak ter-update).
4. Coba input tidak valid: quantity 0/negatif, supplier tidak dipilih, tidak ada item sama
   sekali.

Cek console error dan response API 4xx/5xx tidak wajar. Di akhir, berikan daftar bug dengan
langkah reproduksi. Perbaiki bug kecil/jelas langsung. Untuk temuan soal stok yang tidak
akurat, berhenti dan diskusikan dulu dengan saya karena ini menyangkut integritas data.
```

## Sub-fase 3.3 — Pembayaran Pembelian

```
Lanjut testing debug fungsional otomatis (Anda yang mengeksekusi lewat Playwright) untuk
fitur pembayaran di halaman Pembelian Supplier (/suppliers/purchases) di aplikasi POS ini
(d:\Develop\Project_pos_mahenz). Ikuti "Persiapan Lingkungan Teknis" di
docs/browser-testing-plan.md. Login sebagai admin/admin123. Pakai pembelian dari Sub-fase 3.2.

Cakupan sub-fase ini (fokus HANYA pembayaran):
1. Bayar pembelian secara PARTIAL (sebagian), verifikasi status berubah jadi "Sebagian" dan
   sisa hutang terhitung benar.
2. Bayar sisa hutangnya hingga LUNAS, verifikasi status berubah jadi "Lunas".
3. Lihat riwayat pembayaran (GetPayments) untuk pembelian itu — pastikan semua pembayaran
   partial+lunas tadi tercatat lengkap dan totalnya cocok dengan total pembelian.
4. Coba skenario ganjil: bayar melebihi sisa hutang, bayar pembelian yang sudah lunas —
   pastikan ada guard yang jelas (ditolak dengan pesan jelas), bukan data pembayaran jadi
   negatif atau melebihi total pembelian.

Cek console error dan response API 4xx/5xx tidak wajar. Di akhir, berikan daftar bug dengan
langkah reproduksi. Perbaiki bug kecil/jelas langsung. Untuk temuan soal perhitungan uang
yang salah, berhenti dan diskusikan dulu dengan saya.
```

## Sub-fase 3.4 — Retur Pembelian

```
Lanjut testing debug fungsional otomatis (Anda yang mengeksekusi lewat Playwright) untuk
halaman Retur Pembelian (/suppliers/returns) di aplikasi POS ini
(d:\Develop\Project_pos_mahenz). Ikuti "Persiapan Lingkungan Teknis" di
docs/browser-testing-plan.md. Login sebagai admin/admin123. Pakai pembelian dari Sub-fase 3.2.
Baca docs/audit/04-procurement.md dulu sebelum mulai kalau ada gejala mencurigakan soal
stok/uang.

Cakupan sub-fase ini (fokus HANYA menu ini):
1. Buat retur berdasarkan pembelian yang sudah ada, pilih sebagian item/quantity untuk
   diretur (bukan semua), verifikasi ke BE: stok produk BENAR-BENAR berkurang sesuai
   quantity yang diretur.
2. Coba retur MELEBIHI quantity yang pernah dibeli — pastikan ditolak dengan guard yang
   jelas, bukan stok jadi negatif atau salah hitung.
3. Update status retur (kalau ada alur approve/proses), cek konsistensi status di seluruh
   halaman.
4. Coba retur dari pembelian yang belum pernah ada (ID tidak valid) atau retur dobel untuk
   quantity yang sama — pastikan ada guard.

Cek console error dan response API 4xx/5xx tidak wajar. Di akhir, berikan daftar bug dengan
langkah reproduksi. Perbaiki bug kecil/jelas langsung. Untuk temuan soal stok/uang yang
salah, berhenti dan diskusikan dulu dengan saya.
```

---

# FASE 4 — Operasional & Penjualan

Area paling kritis menurut audit statis sebelumnya (checkout, void, cash-drawer) — dipecah
paling detail, kerjakan dengan hati-hati dan teliti, urut sesuai nomor.

## Sub-fase 4.1 — Manajemen Shift

```
Lanjut testing debug fungsional otomatis (Anda yang mengeksekusi lewat Playwright) untuk
halaman Manajemen Shift (/shifts) di aplikasi POS ini (d:\Develop\Project_pos_mahenz). Ikuti
"Persiapan Lingkungan Teknis" di docs/browser-testing-plan.md. Login sebagai admin/admin123.

Cakupan sub-fase ini (fokus HANYA menu ini, JANGAN dulu ke alur buka-shift-dari-kasir — itu
di sub-fase berikutnya):
1. Buat shift baru (nama shift, jam mulai/selesai kalau ada field itu), edit, toggle status
   aktif/nonaktif.
2. Lihat shift yang sedang aktif (endpoint "active"), pastikan datanya benar.
3. Coba input tidak valid: nama kosong, jam tidak masuk akal (selesai sebelum mulai kalau
   ada validasi itu).

Cek console error dan response API 4xx/5xx tidak wajar. Di akhir, berikan daftar bug dengan
langkah reproduksi, perbaiki yang jelas/kecil langsung, diskusikan dulu untuk temuan besar.
```

## Sub-fase 4.2 — Buka Shift & Kas Awal dari Kasir

```
Lanjut testing debug fungsional otomatis (Anda yang mengeksekusi lewat Playwright) untuk
alur BUKA SHIFT dari halaman Kasir (/kasir) di aplikasi POS ini
(d:\Develop\Project_pos_mahenz). Ikuti "Persiapan Lingkungan Teknis" di
docs/browser-testing-plan.md. Login sebagai user kasir yang dibuat di Sub-fase 1.2 (kalau
belum ada, buat dulu lewat Manajemen User sebagai owner).

Cakupan sub-fase ini (fokus HANYA proses buka shift/kas, JANGAN dulu bertransaksi — itu di
sub-fase berikutnya):
1. Coba buka halaman Kasir SEBELUM ada shift/kas dibuka — pastikan ada pesan jelas ("Buka
   shift terlebih dahulu") dan tombol bayar disabled, bukan bisa checkout tanpa kas dibuka.
2. Buka kas dengan modal awal (starting cash) tertentu, verifikasi ke BE: data cash-drawer
   baru benar-benar dibuat dengan modal awal yang sesuai, terhubung ke shift dan user yang
   benar.
3. Cek halaman Kas Harian (/finance/cash-drawer) menampilkan kas yang baru dibuka ini dengan
   status "Buka".
4. Coba buka kas KEDUA KALINYA saat kas pertama masih terbuka (belum ditutup) — pastikan ada
   guard yang mencegah (satu user tidak boleh punya 2 kas terbuka bersamaan), bukan malah
   membuat 2 cash-drawer aktif yang membingungkan.

Cek console error dan response API 4xx/5xx tidak wajar. Di akhir, berikan daftar bug dengan
langkah reproduksi, perbaiki yang jelas/kecil langsung, diskusikan dulu untuk temuan besar.
```

## Sub-fase 4.3 — Kasir: Alur Keranjang & Checkout Dasar

```
Lanjut testing debug fungsional otomatis (Anda yang mengeksekusi lewat Playwright) untuk
alur keranjang & checkout di halaman Kasir (/kasir) di aplikasi POS ini
(d:\Develop\Project_pos_mahenz). Ikuti "Persiapan Lingkungan Teknis" di
docs/browser-testing-plan.md. Login sebagai user kasir, pastikan kas sudah dibuka (Sub-fase
4.2). Pakai produk dari Fase 2.

Cakupan sub-fase ini (fokus HANYA alur UI keranjang→bayar, verifikasi keamanan angka ada di
sub-fase 4.4 terpisah):
1. Cari produk lewat search dan lewat input barcode, tambah ke keranjang, ubah quantity
   (naik/turun), hapus item dari keranjang.
2. Tambah pelanggan ke transaksi (pilih dari daftar pelanggan kalau sudah ada dari testing
   lain, atau lewati kalau belum ada).
3. Coba checkout dengan metode TUNAI: bayar pas, bayar lebih (cek kembalian dihitung benar),
   bayar KURANG dari total (harus ditolak, tidak boleh checkout).
4. Coba checkout dengan metode TRANSFER dan metode lain yang tersedia.
5. Coba transaksi dengan diskon (kalau ada field diskon di kasir) — item diskon dan/atau
   diskon total.
6. Setelah bayar, pastikan keranjang kosong lagi dan siap untuk transaksi berikutnya, dan
   struk/nota (kalau ada preview cetak) menampilkan rincian yang benar.

Cek console error dan response API 4xx/5xx tidak wajar. Di akhir, berikan daftar bug dengan
langkah reproduksi, perbaiki yang jelas/kecil langsung, diskusikan dulu untuk temuan besar.
```

## Sub-fase 4.4 — Kasir: Verifikasi Keamanan Perhitungan Uang di Server

```
Lanjut testing debug fungsional otomatis (Anda yang mengeksekusi lewat Playwright + curl
langsung ke API) untuk memverifikasi KEAMANAN perhitungan transaksi kasir di aplikasi POS
ini (d:\Develop\Project_pos_mahenz). Baca docs/audit/06-sales.md dulu (temuan #1, #2, #3, #6
soal checkout tidak menghitung ulang uang di server) sebelum mulai — sub-fase ini murni
untuk memverifikasi apakah temuan itu benar dan masih terjadi.

Ikuti "Persiapan Lingkungan Teknis" di docs/browser-testing-plan.md untuk menjalankan BE+FE.
Ini sub-fase yang butuh manipulasi request langsung (bukan cuma isi form normal), jadi boleh
pakai kombinasi Playwright (untuk dapat token & konteks yang sah) dan curl manual (untuk
mengirim payload yang dimodifikasi seolah dari browser yang di-tamper).

Cakupan sub-fase ini:
1. Login sebagai kasir, buka shift (kalau belum), lalu lewat curl (bukan lewat form UI),
   kirim POST /api/transactions/create dengan payload yang harga/subtotal/total_amount-nya
   SENGAJA dibuat tidak cocok dengan harga asli produk di database (misal produk asli Rp
   50.000 tapi kirim price: 1). Amati: apakah BE menolak/mengoreksi, atau malah menyimpan
   nilai palsu itu apa adanya.
2. Kirim payload dengan discount yang lebih besar dari subtotal (diskon melebihi harga) —
   amati apakah ada validasi batas atas.
3. Kirim payload dengan shift_id milik user LAIN (bukan user yang sedang login) — amati
   apakah BE memvalidasi kepemilikan shift atau menerima begitu saja.
4. Cek apakah endpoint POST /api/transactions/create punya permission middleware (can_create)
   terpasang, konsisten dengan endpoint lain seperti supplier-purchases/create.

Laporkan semua temuan dengan detail: request yang dikirim, response yang diterima, dan data
akhir yang tersimpan di database (cek lewat GET detail transaksi). JANGAN perbaiki apa pun
di sub-fase ini — ini murni untuk memastikan lengkap tidaknya bukti sebelum kita diskusikan
bersama cara perbaikannya, karena ini menyangkut keputusan desain penting soal validasi
server-side untuk data uang.
```

## Sub-fase 4.5 — Transaksi: List & Detail

```
Lanjut testing debug fungsional otomatis (Anda yang mengeksekusi lewat Playwright) untuk
halaman Transaksi (/transactions) di aplikasi POS ini (d:\Develop\Project_pos_mahenz). Ikuti
"Persiapan Lingkungan Teknis" di docs/browser-testing-plan.md. Login sebagai admin/admin123.
Pakai transaksi dari Sub-fase 4.3.

Cakupan sub-fase ini (fokus HANYA list & detail, JANGAN dulu void — itu di sub-fase
berikutnya):
1. Lihat daftar transaksi, coba semua filter (tanggal, kode/pelanggan, metode bayar, status),
   coba sorting per kolom.
2. Lihat detail transaksi (item, harga, diskon, pajak, kasir, pelanggan) — cocokkan dengan
   yang sungguh diinput di Sub-fase 4.3.
3. Cek pagination bekerja benar kalau data lebih dari 1 halaman.

Cek console error dan response API 4xx/5xx tidak wajar. Di akhir, berikan daftar bug dengan
langkah reproduksi, perbaiki yang jelas/kecil langsung, diskusikan dulu untuk temuan besar.
```

## Sub-fase 4.6 — Void Transaksi & Dampak ke Cash-Drawer

```
Lanjut testing debug fungsional otomatis (Anda yang mengeksekusi lewat Playwright + curl)
untuk fitur VOID transaksi di aplikasi POS ini (d:\Develop\Project_pos_mahenz). Baca
docs/audit/06-sales.md temuan #4 dan #5 (void tidak reverse cash-drawer, update cash-drawer
tidak atomik) dulu sebelum mulai. Ikuti "Persiapan Lingkungan Teknis" di
docs/browser-testing-plan.md.

Cakupan sub-fase ini:
1. Catat total_cash_sales di cash-drawer AKTIF (lewat GET /api/cash-drawer/current) SEBELUM
   void.
2. Void salah satu transaksi TUNAI yang dibuat di Sub-fase 4.3 (lewat halaman /transactions).
3. Cek lagi total_cash_sales di cash-drawer SETELAH void — pastikan nilainya BERKURANG
   sejumlah nilai transaksi yang di-void. Kalau tidak berkurang (tetap sama seperti sebelum
   void), itu bug nyata yang cocok dengan dugaan di audit.
4. Coba void transaksi yang sama DUA KALI — pastikan ada guard (tidak bisa void yang sudah
   void), dan cash-drawer tidak double-berkurang.
5. Coba void transaksi lalu tutup shift — cek ringkasan penutupan shift tidak ikut menghitung
   transaksi yang sudah di-void sebagai penjualan sah.

Laporkan temuan dengan angka konkret (before/after). Untuk bug soal cash-drawer tidak ter-
reverse, JANGAN perbaiki langsung — berhenti dan diskusikan dulu dengan saya karena
menyangkut keputusan desain (perlu transaksi DB atomik antara void + update cash-drawer).
```

## Sub-fase 4.7 — Tutup Shift & Ringkasan Penutupan

```
Lanjut testing debug fungsional otomatis (Anda yang mengeksekusi lewat Playwright) untuk
alur TUTUP SHIFT dari Kasir/Kas Harian di aplikasi POS ini (d:\Develop\Project_pos_mahenz).
Ikuti "Persiapan Lingkungan Teknis" di docs/browser-testing-plan.md. Login sebagai user
kasir yang kas-nya sudah dibuka & dipakai bertransaksi di sub-fase sebelumnya.

Cakupan sub-fase ini:
1. Tutup kas (Close), masukkan jumlah kas fisik yang dihitung manual (cash count) — coba
   masukkan angka yang PAS dengan sistem dan yang BERBEDA (selisih lebih/kurang), verifikasi
   selisih (variance) dihitung dan ditampilkan dengan benar.
2. Verifikasi ringkasan penutupan (total penjualan, total tunai, total non-tunai) sesuai
   dengan transaksi yang sungguh-sungguh terjadi selama shift itu (yang berhasil, bukan yang
   di-void — cross-check dengan Sub-fase 4.6).
3. Setelah kas ditutup, coba buka halaman Kasir lagi — pastikan diminta buka kas baru
   (tidak bisa transaksi pakai kas yang sudah ditutup).
4. Cek riwayat kas (Kas Harian) menampilkan kas yang baru ditutup ini dengan status "Tutup"
   dan semua angka yang benar.

Cek console error dan response API 4xx/5xx tidak wajar. Di akhir, berikan daftar bug dengan
langkah reproduksi. Perbaiki bug kecil/jelas langsung. Untuk selisih perhitungan yang salah,
berhenti dan diskusikan dulu dengan saya.
```

---

# FASE 5 — Pelanggan & Piutang

## Sub-fase 5.1 — Pelanggan (CRUD)

```
Lanjut testing debug fungsional otomatis (Anda yang mengeksekusi lewat Playwright) untuk
halaman Pelanggan (/customers) di aplikasi POS ini (d:\Develop\Project_pos_mahenz). Ikuti
"Persiapan Lingkungan Teknis" di docs/browser-testing-plan.md. Login sebagai admin/admin123.

Cakupan sub-fase ini (fokus HANYA CRUD pelanggan, JANGAN dulu ke transaksi kredit/piutang —
itu di sub-fase berikutnya):
1. Tambah pelanggan baru (nama, kontak, atur limit kredit/credit_limit), edit, toggle
   aktif/nonaktif.
2. Coba input tidak valid: nama kosong, credit_limit negatif.
3. Coba nonaktifkan/hapus pelanggan yang nanti (Sub-fase 5.2/5.3) sudah punya transaksi
   kredit/piutang — kalau sub-fase ini dikerjakan lebih dulu, catat untuk dicoba ulang
   setelahnya.

Cek console error dan response API 4xx/5xx tidak wajar. Di akhir, berikan daftar bug dengan
langkah reproduksi, perbaiki yang jelas/kecil langsung, diskusikan dulu untuk temuan besar.
```

## Sub-fase 5.2 — Transaksi Kredit → Piutang Muncul

```
Lanjut testing debug fungsional otomatis (Anda yang mengeksekusi lewat Playwright) untuk
alur transaksi KREDIT dari Kasir sampai muncul di Piutang, di aplikasi POS ini
(d:\Develop\Project_pos_mahenz). Ikuti "Persiapan Lingkungan Teknis" di
docs/browser-testing-plan.md. Login sebagai user kasir (kas harus sudah dibuka), pakai
pelanggan dari Sub-fase 5.1.

Cakupan sub-fase ini:
1. Buat transaksi di Kasir dengan metode kredit (is_credit=true, pilih pelanggan wajib untuk
   metode ini) — kalau UI tidak menyediakan opsi ini secara eksplisit, cari tahu dari kode
   FE (FE/src/features/sales/cashier) bagaimana cara memicunya, atau laporkan kalau memang
   tidak ada jalur UI untuk membuat transaksi kredit sama sekali (itu sendiri temuan penting).
2. Setelah transaksi kredit dibuat, cek halaman Piutang (/receivables) — pastikan piutang
   baru muncul dengan jumlah dan pelanggan yang benar.
3. Cek transaksi kredit ini TIDAK tercatat sebagai pemasukan tunai/non-tunai di cash-drawer
   (karena belum benar-benar dibayar), tapi tetap tercatat sebagai penjualan di laporan.

Cek console error dan response API 4xx/5xx tidak wajar. Di akhir, berikan daftar bug dengan
langkah reproduksi, perbaiki yang jelas/kecil langsung, diskusikan dulu untuk temuan besar.
```

## Sub-fase 5.3 — Pembayaran Piutang

```
Lanjut testing debug fungsional otomatis (Anda yang mengeksekusi lewat Playwright) untuk
fitur pembayaran piutang (/receivables) di aplikasi POS ini (d:\Develop\Project_pos_mahenz).
Baca docs/audit/02-customers.md dulu sebelum mulai — sudah ada temuan kritis soal fitur
bayar piutang yang diduga rusak, sub-fase ini untuk verifikasi apakah masih terjadi. Ikuti
"Persiapan Lingkungan Teknis" di docs/browser-testing-plan.md. Login sebagai admin/admin123.
Pakai piutang dari Sub-fase 5.2.

Cakupan sub-fase ini:
1. Bayar piutang secara PARTIAL, verifikasi sisa piutang berkurang dengan benar dan status
   berubah jadi "Bayar Sebagian".
2. Bayar sisa piutang hingga LUNAS, verifikasi status berubah jadi "Lunas".
3. Lihat riwayat pembayaran (GetPayments) — pastikan semua pembayaran tercatat lengkap dan
   totalnya cocok dengan total piutang awal.
4. Coba skenario ganjil: bayar melebihi sisa piutang, bayar piutang yang sudah lunas —
   pastikan ada guard yang jelas, bukan data piutang jadi negatif.

Cek console error dan response API 4xx/5xx tidak wajar. Di akhir, berikan daftar bug dengan
langkah reproduksi. Perbaiki bug kecil/jelas langsung. Untuk bug yang cocok dengan temuan
kritis di audit sebelumnya, berhenti dan diskusikan dulu dengan saya.
```

---

# FASE 6 — Keuangan

## Sub-fase 6.1 — Kas Harian (Riwayat & Selisih)

```
Lanjut testing debug fungsional otomatis (Anda yang mengeksekusi lewat Playwright) untuk
halaman Kas Harian (/finance/cash-drawer) di aplikasi POS ini (d:\Develop\Project_pos_mahenz).
Ikuti "Persiapan Lingkungan Teknis" di docs/browser-testing-plan.md. Login sebagai
admin/admin123. Pakai data kas dari Fase 4.

Cakupan sub-fase ini (fokus di luar yang sudah dicek di Sub-fase 4.2/4.7 — di sini dari sudut
pandang admin melihat SEMUA kas, bukan kas sendiri):
1. Lihat riwayat semua kas (dari semua kasir), filter tanggal, filter kasir.
2. Lihat detail satu kas yang sudah ditutup — pastikan semua rincian (modal awal, penjualan
   tunai, pengeluaran, selisih) ditampilkan lengkap dan konsisten dengan yang dicek di
   Sub-fase 4.7.
3. Cek ringkasan (summary) kas harian di halaman ini (kalau ada agregat total semua kasir).

Cek console error dan response API 4xx/5xx tidak wajar. Di akhir, berikan daftar bug dengan
langkah reproduksi, perbaiki yang jelas/kecil langsung, diskusikan dulu untuk temuan besar.
```

## Sub-fase 6.2 — Kas Saya (Isolasi Data per Kasir)

```
Lanjut testing debug fungsional otomatis (Anda yang mengeksekusi lewat Playwright + curl)
untuk halaman Kas Saya (/finance/my-cash) di aplikasi POS ini (d:\Develop\Project_pos_mahenz).
Ikuti "Persiapan Lingkungan Teknis" di docs/browser-testing-plan.md.

Cakupan sub-fase ini (fokus KEAMANAN isolasi data antar kasir):
1. Login sebagai user kasir A (dari Sub-fase 1.2), buka /finance/my-cash, catat data kas
   yang tampil (harus cuma kas milik dia sendiri).
2. Buat user kasir B baru (kalau belum ada) lewat Manajemen User, buat kas untuk kasir B
   juga (buka shift, transaksi sedikit).
3. Login sebagai kasir A lagi, lewat curl langsung panggil endpoint kas dengan mencoba akses
   data kas milik kasir B (misal GET /api/cash-drawer/detail/:id pakai ID kas milik B, sambil
   pakai token milik A) — pastikan DITOLAK (403/404), bukan malah bisa melihat data kas orang
   lain.

Cek console error dan response API 4xx/5xx tidak wajar. Di akhir, berikan daftar bug dengan
langkah reproduksi. JANGAN perbaiki bug kebocoran data antar user secara langsung — berhenti
dan diskusikan dulu dengan saya karena ini isu keamanan (authorization/IDOR).
```

## Sub-fase 6.3 — Pengeluaran

```
Lanjut testing debug fungsional otomatis (Anda yang mengeksekusi lewat Playwright) untuk
halaman Pengeluaran (/finance/expenses) di aplikasi POS ini (d:\Develop\Project_pos_mahenz).
Ikuti "Persiapan Lingkungan Teknis" di docs/browser-testing-plan.md. Login sebagai
admin/admin123.

Cakupan sub-fase ini (fokus HANYA menu ini):
1. Tambah pengeluaran baru dengan berbagai kategori, edit, hapus.
2. Filter berdasarkan tanggal dan kategori, pastikan totalnya benar.
3. Coba input tidak valid: nominal negatif/nol, kategori kosong, tanggal masa depan (kalau
   ada validasi itu).

Cek console error dan response API 4xx/5xx tidak wajar. Di akhir, berikan daftar bug dengan
langkah reproduksi, perbaiki yang jelas/kecil langsung, diskusikan dulu untuk temuan besar.
```

## Sub-fase 6.4 — Dashboard Keuangan (Cross-check)

```
Lanjut testing debug fungsional otomatis (Anda yang mengeksekusi lewat Playwright) untuk
halaman Dashboard Keuangan (/finance) di aplikasi POS ini (d:\Develop\Project_pos_mahenz).
Ikuti "Persiapan Lingkungan Teknis" di docs/browser-testing-plan.md. Login sebagai
admin/admin123.

Cakupan sub-fase ini:
1. Cocokkan angka Pemasukan di dashboard ini dengan total transaksi sukses dari Fase 4
   (jumlahkan manual dari data yang sudah diketahui, bandingkan).
2. Cocokkan angka Pengeluaran dengan data dari Sub-fase 6.3.
3. Cocokkan angka Piutang Terbuka dengan sisa piutang dari Sub-fase 5.3.
4. Coba ganti filter periode (hari ini/minggu ini/bulan ini/custom range), pastikan semua
   angka ikut berubah sesuai periode, TERMASUK transaksi yang baru saja dibuat beberapa menit
   lalu (area yang pernah bermasalah soal filter tanggal date-only vs datetime — pastikan
   tidak muncul lagi di sini).

Cek console error dan response API 4xx/5xx tidak wajar. Di akhir, berikan daftar bug dengan
selisih angka konkret (harusnya X, tapi tampil Y), perbaiki yang jelas/kecil langsung,
diskusikan dulu untuk temuan besar.
```

---

# FASE 7 — Pelaporan & Dashboard Utama

## Sub-fase 7.1 — Dashboard Utama

```
Lanjut testing debug fungsional otomatis (Anda yang mengeksekusi lewat Playwright) untuk
Dashboard utama (/dashboard) di aplikasi POS ini (d:\Develop\Project_pos_mahenz). Ikuti
"Persiapan Lingkungan Teknis" di docs/browser-testing-plan.md. Login sebagai owner/owner123.

Cakupan sub-fase ini:
1. Cocokkan kartu ringkasan (transaksi hari ini, pendapatan, laba kotor, stok menipis,
   piutang terbuka) dengan data riil dari fase-fase sebelumnya.
2. Cek grafik penjualan (sales trend) — hover tiap titik, pastikan tooltip menampilkan
   tanggal terformat dengan benar (bukan ISO string mentah — ini sudah pernah diperbaiki,
   pastikan tidak regresi) dan angka yang akurat.
3. Cek widget Top Produk Terlaris dan Top Kategori — cocokkan dengan transaksi yang pernah
   dibuat.
4. Coba ganti periode (Hari ini/Minggu ini/Bulan ini), pastikan semua widget ikut berubah.

Cek console error dan response API 4xx/5xx tidak wajar. Di akhir, berikan daftar bug dengan
langkah reproduksi, perbaiki yang jelas/kecil langsung, diskusikan dulu untuk temuan besar.
```

## Sub-fase 7.2 — Laporan Penjualan

```
Lanjut testing debug fungsional otomatis (Anda yang mengeksekusi lewat Playwright) untuk
Laporan Penjualan (/reports/sales) di aplikasi POS ini (d:\Develop\Project_pos_mahenz). Ikuti
"Persiapan Lingkungan Teknis" di docs/browser-testing-plan.md. Login sebagai admin/admin123.

Cakupan sub-fase ini:
1. Filter tanggal (termasuk rentang yang mencakup HARI INI dengan transaksi baru dari
   beberapa menit lalu — area yang pernah bermasalah, pastikan sudah benar) dan metode bayar,
   pastikan data & total cocok dengan transaksi riil.
2. Export Excel, buka file hasil export (baca isinya), cocokkan dengan yang tampil di layar.
3. Cek transaksi yang sudah di-void TIDAK ikut terhitung sebagai penjualan di laporan ini.

Cek console error dan response API 4xx/5xx tidak wajar. Di akhir, berikan daftar bug dengan
langkah reproduksi, perbaiki yang jelas/kecil langsung, diskusikan dulu untuk temuan besar.
```

## Sub-fase 7.3 — Laporan Laba Rugi

```
Lanjut testing debug fungsional otomatis (Anda yang mengeksekusi lewat Playwright) untuk
Laporan Laba Rugi (/reports/profit-loss) di aplikasi POS ini (d:\Develop\Project_pos_mahenz).
Ikuti "Persiapan Lingkungan Teknis" di docs/browser-testing-plan.md. Login sebagai
admin/admin123.

Cakupan sub-fase ini:
1. Filter tanggal (termasuk hari ini), pastikan Total Pendapatan cocok dengan Laporan
   Penjualan (Sub-fase 7.2).
2. Verifikasi HPP (harga pokok penjualan) dihitung dari harga beli produk yang benar-benar
   dipakai di transaksi (bukan harga beli produk yang sedang berlaku sekarang kalau harga
   beli sempat berubah — cek apakah sistem mencatat harga beli historis atau ambil harga
   beli terkini, ini beda perilaku yang penting).
3. Cek Total Pengeluaran cocok dengan Sub-fase 6.3, dan Laba Bersih = Laba Kotor - Total
   Pengeluaran benar secara matematis.
4. Export Excel, cocokkan isinya.

Cek console error dan response API 4xx/5xx tidak wajar. Di akhir, berikan daftar bug dengan
langkah reproduksi, perbaiki yang jelas/kecil langsung, diskusikan dulu untuk temuan besar
(terutama soal HPP historis vs current yang bisa jadi keputusan desain, bukan sekadar bug).
```

## Sub-fase 7.4 — Laporan Stok

```
Lanjut testing debug fungsional otomatis (Anda yang mengeksekusi lewat Playwright) untuk
Laporan Stok (/reports/stock) di aplikasi POS ini (d:\Develop\Project_pos_mahenz). Ikuti
"Persiapan Lingkungan Teknis" di docs/browser-testing-plan.md. Login sebagai admin/admin123.

Cakupan sub-fase ini:
1. Filter kategori, cari produk, pastikan stok & nilai stok (stok × harga beli) cocok dengan
   data produk aktual dari Fase 2-3 (setelah pembelian & retur).
2. Cek widget "Produk Stok Rendah" menandai produk yang stoknya di bawah ambang batas dengan
   benar.
3. Export Excel, cocokkan isinya.

Cek console error dan response API 4xx/5xx tidak wajar. Di akhir, berikan daftar bug dengan
langkah reproduksi, perbaiki yang jelas/kecil langsung, diskusikan dulu untuk temuan besar.
```

## Sub-fase 7.5 — Kinerja Kasir

```
Lanjut testing debug fungsional otomatis (Anda yang mengeksekusi lewat Playwright) untuk
Laporan Kinerja Kasir (/reports/cashier) di aplikasi POS ini (d:\Develop\Project_pos_mahenz).
Ikuti "Persiapan Lingkungan Teknis" di docs/browser-testing-plan.md. Login sebagai
admin/admin123.

Cakupan sub-fase ini:
1. Filter tanggal (termasuk hari ini), pastikan data per kasir (jumlah transaksi, total
   penjualan, tunai/non-tunai, void count) sesuai dengan yang sungguh terjadi di Fase 4
   (termasuk kasir A dan kasir B dari Sub-fase 6.2 kalau keduanya sudah pernah bertransaksi).
2. Export Excel, cocokkan isinya.

Cek console error dan response API 4xx/5xx tidak wajar. Di akhir, berikan daftar bug dengan
langkah reproduksi, perbaiki yang jelas/kecil langsung, diskusikan dulu untuk temuan besar.
```

---

# FASE 8 — Sistem Lanjutan & Operasional Sisa

## Sub-fase 8.1 — Printer & Versi Aplikasi

```
Lanjut testing debug fungsional otomatis (Anda yang mengeksekusi lewat Playwright) untuk
halaman Pengaturan Printer (/settings/printer) dan Versi Aplikasi (/settings/versions) di
aplikasi POS ini (d:\Develop\Project_pos_mahenz). Ikuti "Persiapan Lingkungan Teknis" di
docs/browser-testing-plan.md. Login sebagai owner/owner123.

Cakupan sub-fase ini:
1. Pengaturan Printer: ubah pengaturan (lebar kertas, dll kalau ada), simpan, refresh
   halaman (full reload), pastikan tersimpan permanen.
2. Versi Aplikasi: lihat daftar versi, coba form update versi Android kalau ada (meski
   aplikasi Android belum dipakai — pastikan minimal tidak error/crash saat dicoba).

Cek console error dan response API 4xx/5xx tidak wajar. Di akhir, berikan daftar bug dengan
langkah reproduksi, perbaiki yang jelas/kecil langsung, diskusikan dulu untuk temuan besar.
```

## Sub-fase 8.2 — Backup & Restore

```
Lanjut testing debug fungsional otomatis (Anda yang mengeksekusi lewat Playwright) untuk
fitur Backup & Restore di aplikasi POS ini (d:\Develop\Project_pos_mahenz). Ikuti "Persiapan
Lingkungan Teknis" di docs/browser-testing-plan.md. Login sebagai owner/owner123.

PERINGATAN: fitur restore itu DESTRUKTIF (bisa menimpa seluruh data yang sudah dibuat di
fase-fase sebelumnya). Untuk sub-fase ini:
1. Buat backup baru, lihat daftar backup, coba download file backup-nya, verifikasi filenya
   valid (bisa dibuka/ada isinya, bukan file kosong/korup).
2. JANGAN langsung menjalankan proses restore. Jelaskan dulu ke saya: file backup apa yang
   akan dipakai, dan konfirmasi eksplisit dari saya sebelum benar-benar klik restore. Kalau
   saya setujui, baru jalankan dan verifikasi data kembali ke kondisi backup tsb.

Cek console error dan response API 4xx/5xx tidak wajar. Di akhir, berikan daftar bug dengan
langkah reproduksi, perbaiki yang jelas/kecil langsung (di luar aksi restore itu sendiri).
```

## Sub-fase 8.3 — Sync Center

```
Lanjut testing debug fungsional otomatis (Anda yang mengeksekusi lewat Playwright) untuk
halaman Sync Center (/sync) di aplikasi POS ini (d:\Develop\Project_pos_mahenz). Baca
docs/audit/05-operational.md dulu sebelum mulai. Ikuti "Persiapan Lingkungan Teknis" di
docs/browser-testing-plan.md. Login sebagai owner/owner123.

Cakupan sub-fase ini:
1. Coba trigger sync manual, lihat status berubah (Siap/Sinkronisasi/dll), lihat riwayat
   sinkronisasi ter-update.
2. Lihat antrian sync (queue) dan konflik (conflicts) kalau ada data yang sengaja dibuat
   konflik — coba selesaikan satu konflik lewat UI resolve, verifikasi hasilnya masuk akal.
3. Cek jumlah konflik (conflict count) di halaman ini konsisten dengan daftar konflik yang
   sungguh ada.

Cek console error dan response API 4xx/5xx tidak wajar. Di akhir, berikan daftar bug dengan
langkah reproduksi, perbaiki yang jelas/kecil langsung, diskusikan dulu untuk temuan besar.
```

---

# Setelah semua sub-fase selesai

Minta rangkuman akhir: daftar semua bug yang ditemukan di seluruh 24 sub-fase (yang sudah
diperbaiki vs yang masih didiskusikan/tertunda), supaya ada satu dokumen ringkasan final
sebelum lanjut ke tahap berikutnya (misal staging/production).

# Rancangan: Koreksi HPP/Modal Historis (Opsi B) + Menu Rekonsiliasi Modal

> Status: **RANCANGAN KONSEP — belum ada implementasi & belum ada data yang diubah.**
> Dokumen ini fokus pada satu pekerjaan: membetulkan `transaction_items.purchase_price`
> (modal per baris penjualan) yang salah pada data historis, supaya Laporan Laba Rugi &
> kartu "Laba Kotor" di Ringkasan Bisnis kembali akurat.
>
> Konteks & investigasi lengkap ada di `docs/DISKUSI_PERBAIKAN_PELAPORAN.md` §6.1.
> File ini adalah pecahan teknis dari keputusan di sana (bagian G).

---

## 1. Ringkas Masalah (1 paragraf)

Pada transaksi lama (mayoritas Agustus 2026, sebagian awal September), kolom
`transaction_items.purchase_price` (modal per satuan jual) tersimpan SALAH: memakai harga
modal satuan BESAR (mis. per Slop) untuk penjualan satuan KECIL (mis. per Pack/Batang),
tanpa dibagi faktor konversi. Akibatnya HPP membengkak 10-100x pada baris terdampak,
membuat Laba Rugi tampak rugi padahal tidak. Bug pencatatan sudah diperbaiki untuk transaksi
baru (± 8 Sep 2026); yang tersisa adalah membetulkan DATA LAMA. Perbaikan stok terdahulu
TIDAK menyentuh kolom modal ini.

Bukti data (DB dev): 447 dari 1.105 baris "janggal" (modal ≥ harga jual); kontribusi COGS
salah ± Rp 46,9 jt lifetime (Agu Rp 38,9 jt + Sep Rp 7,9 jt + Jul Rp 62 rb).

---

## 2. Keputusan yang Sudah Final

| # | Topik | Keputusan |
|---|-------|-----------|
| 1 | Metode | **Opsi B** — hitung ulang modal dari harga beli NYATA (`purchase_items`), bukan tebakan/harga sekarang. |
| 2 | Cakupan | **Hitung ulang SEMUA baris** `transaction_items` dengan satu metode konsisten (last-cost), bukan hanya baris janggal. Menghindari campur metode. |
| 3 | Aturan harga | **Last-cost**: harga beli terdekat dengan `purchase_date <= tanggal transaksi`. |
| 4 | Fallback | Kalau penjualan terjadi SEBELUM pembelian pertama produk itu → pakai **harga pembelian PERTAMA** (paling awal) sebagai acuan. |
| 5 | Data kosong | Produk yang **sama sekali tak punya pembelian** → TIDAK dikoreksi otomatis, masuk **menu Rekonsiliasi Modal** (input manual admin). Saat ini hanya 1 produk: LA Bold 20 (id 70). |
| 6 | Urutan kerja | **Skrip koreksi dulu**, menu Rekonsiliasi Modal menyusul. |
| 7 | Lingkungan | Kerjakan & uji di **LOKAL/DEV dulu**. Production = urusan user. |

---

## 3. Formula Koreksi (inti Opsi B)

Untuk setiap baris `transaction_items` (ti) pada transaksi `t` (status completed):

```
harga_beli_per_anchor = pi.purchase_price / NULLIF(pi.conversion_qty, 0)
modal_benar_per_satuan_jual = harga_beli_per_anchor * ti.conversion_qty
cogs_baris = ti.quantity * modal_benar_per_satuan_jual
```

Di mana `pi` = baris `purchase_items` acuan, dipilih dengan aturan:
1. **Utama (last-cost):** `purchase_items` dari `purchases` status 'active', produk sama,
   `purchase_date <= DATE(t.transaction_date)`, urut `purchase_date DESC, id DESC`, ambil 1.
2. **Fallback:** kalau tidak ada (penjualan sebelum pembelian pertama), ambil pembelian
   PALING AWAL produk itu: `purchase_date ASC, id ASC`, ambil 1.
3. **Tak ada pembelian sama sekali:** JANGAN koreksi. Tandai baris/produk → Rekonsiliasi Modal.

> Catatan konversi: `purchase_items.purchase_price` adalah harga per SATUAN BELI (paket),
> `conversion_qty` = faktor satuan itu ke anchor. Jadi `purchase_price/conversion_qty`
> = harga per satuan dasar. Dikali `ti.conversion_qty` (faktor satuan JUAL ke anchor)
> = modal per satuan jual. Dimensi konsisten.

---

## 4. Rencana Skrip Koreksi (one-off, pola `BE/cmd/`)

Mengikuti pola skrip backfill yang sudah ada (mis. `BE/cmd/backfill_stock_restore`).

**Nama usulan:** `BE/cmd/backfill_transaction_cogs/` (atau nama lain sesuai preferensi).

**Perilaku wajib:**
- **`--dry-run` (default ON):** tampilkan ringkasan before/after (total COGS lama vs baru,
  jumlah baris berubah, daftar baris yang gagal/masuk manual) TANPA menulis DB.
- **Backup dulu:** skrip menolak jalan mode tulis kalau belum ada backup (atau minimal
  cetak peringatan tegas). Backup manual via mysqldump tetap wajib.
- **Idempotent:** aman dijalankan berulang. Simpan jejak nilai lama sebelum menimpa
  (lihat §5 audit).
- **Baca `MIGRATION_DSN`** (pola sama dgn skrip backfill lain) supaya bisa jalan di prod
  dengan user/password, default root@localhost utk dev.
- **Logging jelas:** per produk/baris — modal lama → modal baru, sumber pembelian yang dipakai
  (last-cost/fallback), dan yang di-skip (alasan).

**Output akhir (mode tulis):**
- Jumlah baris dikoreksi, jumlah di-skip (manual), total COGS lifetime sebelum vs sesudah.
- Daftar produk yang butuh rekonsiliasi manual.

---

## 5. Audit Trail (jejak perubahan)

Sebelum menimpa `transaction_items.purchase_price`, simpan nilai lama supaya bisa ditelusuri
& di-rollback. Opsi (difinalkan saat implementasi):
- Tabel baru mis. `transaction_items_cogs_audit` (transaction_item_id, old_price, new_price,
  source: last_cost/fallback/manual, purchase_item_ref, changed_by, changed_at, note), ATAU
- Kolom cadangan + file log CSV hasil dry-run yang disimpan.

Rekomendasi: tabel audit khusus — konsisten dgn semangat "jejak audit" di Rekonsiliasi Stok.

---

## 6. Menu Rekonsiliasi Modal (role admin) — konsep, dikerjakan SETELAH skrip

Meniru pola menu **Rekonsiliasi Stok** yang sudah ada.

- **Daftar:** baris/produk yang modalnya TIDAK bisa dihitung otomatis (tak ada pembelian,
  atau hasil hitung ditandai meragukan). Analog `needs_stock_review`.
- **Aksi admin:** input modal per satuan yang benar → update `transaction_items.purchase_price`
  baris terkait + tulis audit (siapa, kapan, dari berapa ke berapa, catatan).
- **Cakupan awal:** minimal 1 produk (LA Bold 20). Dibuat generik agar menampung kasus baru
  ke depan (mis. produk baru yang dijual sebelum sempat dicatat pembeliannya).
- **Permission:** khusus role admin (pola RoleGuard menu key, mirip pelaporan lain).

---

## 7. Konsistensi Antar-Menu (dikerjakan bersamaan / setelah koreksi data)

Setelah data modal benar, pertimbangkan menyatukan sumber COGS agar Laba Rugi & Ringkasan
Bisnis tidak drift (lihat DISKUSI §6.1 E):
- COGS: report pakai `profitLossQuery`, business_summary pakai `rangeCOGSQuery`/`monthCOGSQuery`
  (query terpisah, pola sama). Pertimbangkan satukan.
- Definisi "Pendapatan" beda: report = `(total_amount - tax)` diprorata; business_summary
  = `SUM(total_amount)` penuh. Perlu disepakati definisi tunggal (di luar scope koreksi data,
  tapi relevan untuk keakuratan lintas-menu).

---

## 8. Urutan Eksekusi (checklist)

- [ ] 1. Rancang & tulis skrip `backfill_transaction_cogs` (mode dry-run dulu).
- [ ] 2. Jalankan dry-run di DEV → verifikasi before/after (total COGS, baris berubah, skip list).
- [ ] 3. Uji Laba Rugi & Ringkasan Bisnis di DEV setelah koreksi (via browser) → angka wajar.
- [ ] 4. Siapkan audit trail + prosedur backup untuk production.
- [ ] 5. (Menyusul) Desain & implementasi menu Rekonsiliasi Modal (admin).
- [ ] 6. (Opsional) Satukan sumber COGS & definisi Pendapatan antar-menu.

> Production dijalankan oleh user sendiri. Semua langkah di atas divalidasi di dev lebih dulu.

---

## 8b. HASIL ANALISIS MENYELURUH (20 Sep 2026) — WAJIB DIBACA SEBELUM IMPLEMENTASI

Analisis lintas-domain (baca kode 3 sub-agent + query data dev). Mengubah beberapa asumsi.

### 8b.1 Jalur yang menulis `transaction_items.purchase_price` (write-once)
Kolom ini **tidak pernah di-UPDATE ulang** setelah insert (aman ditimpa historis). Yang menulis:
| Jalur | File | Cara isi purchase_price | conversion_qty |
|-------|------|-------------------------|----------------|
| Online (createOnce) | transaction_repo.go ~260-321 | `product_packages.purchase_price` paket terpilih; fallback `products.purchase_price × faktor` | FAKTOR (ResolvePackageFactor, dihitung server) |
| Sync offline (applySyncTransactionOnce) | transaction_repo.go ~550-586 | sama dgn online | dari PAYLOAD CLIENT apa adanya (tak dihitung ulang) |
| Backdate (CreateTransaction) | backdate_repo.go ~167-240 | `products.purchase_price` MENTAH (anchor, TANPA ×faktor) ❌ | `quantity × package.qty` (semantik BEDA) ❌ |

- **VOID** tidak menyentuh transaction_items (purchase_price tetap). Hanya ubah header status + kembalikan stok + void receivable.
- Tidak ada fitur EDIT item transaksi completed. Tidak ada seed/import yang isi purchase_price.
- **Implikasi:** menimpa purchase_price historis AMAN. TAPI jalur BACKDATE (jalur 3) akan tetap memproduksi modal salah ke depan → perlu diperbaiki terpisah (di luar skrip koreksi data, tapi WAJIB kalau fitur backdate dipakai).

### 8b.2 Status & komponen lain (dari data dev + kode)
- Status transaksi nyata: **739 completed, 10 void** ('pending' tak terpakai). Item di transaksi void: 12 → **WAJIB di-skip** (JOIN, filter `t.status='completed'`).
- Laporan COGS (report + business_summary) SEMUA sudah filter `completed` → koreksi transaction_items otomatis merambat konsisten ke Ringkasan Bisnis.
- `is_credit=1` tetap dihitung COGS → **jangan skip**.
- Retur supplier, expired/write-off, expenses: **INDEPENDEN** dari COGS penjualan (tak sentuh transaction_items). Tak perlu dikoreksi, tak ada risiko double-count. (Catatan: kerugian expired TIDAK masuk Laba Rugi sama sekali — gap pelaporan, di luar scope ini.)
- Tidak ada fitur retur penjualan/pelanggan.

### 8b.3 Kondisi tepi data pembelian (dari query dev) — TEMUAN KRITIS
- conversion_qty NULL/0 di transaksi completed: **0** (aman dari bagi-nol).
- product_id NULL / unit_id NULL di transaksi: **0** (semua bisa dipetakan).
- PO void: **0** (tapi tetap harus difilter `p.status='active'` untuk masa depan).
- Penjualan sebelum pembelian pertama: **0** (fallback jarang, tetap perlu ada).
- Multi-pembelian tanggal sama utk produk sama: **6 produk** (mis. Telur 4 baris di 1 tgl) → **butuh tie-break** last-cost (mis. ambil id purchase_items terbesar / harga terbanyak hari itu).

### 8b.4 ⚠️ TEMUAN PALING PENTING: `purchase_items` SENDIRI PUNYA BUG SATUAN
Formula awal Opsi B (`purchase_price / conversion_qty` dari purchase_items) **TIDAK BISA
DIPERCAYA MENTAH**, karena `conversion_qty` di `purchase_items` pun TIDAK KONSISTEN untuk
data lama. Contoh nyata LA Bold Hitam 20 (id 84):
| tgl | unit | conv_qty | harga/beli | harga/anchor | status |
|-----|------|----------|-----------|--------------|--------|
| 6 Agu | Slop | 1.0 | 374.000 | 374.000 | ✅ benar (Slop=anchor) |
| 29 Agu | Pack | **0.1** | 37.500 | **375.000** | ❌ SALAH (harusnya conv 1.0 → 37.500/pack) |
| 31 Agu | Pack | 1.0 | 37.400 | 37.400 | ✅ benar |

Baris 29 Agu: dua baris "Pack" tapi conv_qty beda (0.1 vs 1.0) → harga per-anchor meleset 10x.
Artinya bug level-satuan yang sama TERJADI JUGA di pencatatan pembelian lama, bukan cuma di
penjualan.

**6 produk dengan variasi harga anchor EKSTREM (>3x) — kandidat purchase_items salah:**
Kopi ABC Botol (24x), Djarum Super 12 (10,5x), LA Bold Hitam 20 (10x), Basreng (10x),
Djarum 76 (9,8x), Telur (5x).

**KONSEKUENSI untuk rancangan:** kalau skrip Opsi B asal ambil "last-cost" dari purchase_items,
ia bisa mengambil baris pembelian yang JUSTRU salah → koreksi jadi salah lagi (persis yang
ingin dihindari). Last-cost mentah TIDAK cukup.

### 8b.5 REVISI STRATEGI (menggantikan §3 bila bentrok)
1. **Jangan percaya `purchase_items.conversion_qty` mentah.** Untuk dapat harga beli
   per-anchor yang benar, normalisasi lewat **struktur `product_packages` yang SEKARANG
   sudah benar** (mis. Djarum: Slop anchor 152.500, Pack 15.250; Telur: Krak 22.500).
   Yakni: petakan `purchase_items.package_id`/`unit` → paket yang benar di product_packages,
   ambil faktor konversi RESMI dari situ (ResolvePackageFactor), bukan dari conv_qty tersimpan.
2. **Sumber "harga beli per-anchor yang benar" per pembelian** = `purchase_price` baris beli
   ÷ faktor-resmi-paket-itu (dari product_packages), lalu dipakai sebagai basis last-cost.
3. **Deteksi & tandai baris pembelian janggal** (harga per-anchor menyimpang jauh dari median
   produk itu, atau modal ≥ harga jual eceran) → jangan dipakai sebagai acuan; pakai baris
   pembelian lain yang wajar, atau kirim produk ke **Rekonsiliasi Modal manual**.
4. **Cross-check hasil**: modal per satuan jual hasil hitung TIDAK boleh ≥ harga jual (`ti.price`)
   pada baris itu (indikator masih salah). Baris yang tetap gagal cross-check → manual.
5. Untuk **multi-pembelian tanggal sama** → tie-break: pilih baris dgn faktor paket yang cocok
   & harga wajar; kalau beda-beda, ambil median/terbanyak.
6. **Alternatif lebih sederhana & mungkin lebih andal** (perlu diputuskan): daripada rekonstruksi
   last-cost per tanggal dari purchase_items yang ternyata juga kotor, gunakan
   **`product_packages.purchase_price` SEKARANG** (yang sudah bersih/benar) × faktor resmi paket
   satuan jual di tiap baris transaksi. Kurang "historis per tanggal", tapi jauh lebih bersih
   karena master paket sekarang sudah rapi. Trade-off: tidak menangkap perubahan harga beli
   antar-waktu (mis. Surya 350rb→352,5rb), tapi menghilangkan risiko mengambil data beli kotor.
   → PERLU KEPUTUSAN USER: akurasi historis-per-tanggal (rekonstruksi rumit + bersihkan
     purchase_items dulu) vs kesederhanaan-dan-kebersihan (pakai harga paket sekarang).

### 8b.6 Rekomendasi urutan revisi
1. Bersihkan/normalkan dulu pemahaman harga beli (putuskan §8b.5 poin 6).
2. Skrip koreksi transaction_items (dry-run) dengan cross-check §8b.5 poin 4.
3. Sisa yang gagal cross-check → menu Rekonsiliasi Modal manual.
4. Perbaiki jalur BACKDATE (§8b.1) agar transaksi backdate baru tidak salah lagi.
5. (Opsional) Pertimbangkan membetulkan purchase_items yang janggal juga, kalau dipakai
   laporan pembelian/valuasi.

---

## 8c. METODE FINAL (DISEPAKATI: akurasi utama, data ragu → manual) — 20 Sep 2026

Keputusan user: **akurasi diutamakan**; baris/produk yang meragukan **BUKAN blocker** —
dilempar ke menu Rekonsiliasi Modal (manual admin), pola sama seperti Rekonsiliasi Stok.

### Metode koreksi (TERBUKTI di data dev)
Untuk tiap baris `transaction_items` milik transaksi `status='completed'`:
```
faktor_jual   = ResolvePackageFactor(product_packages[product], ti.unit_id)   // rantai PENUH ke anchor
anchor_cost   = product_packages(product, is_default=1).purchase_price         // harga modal anchor SEKARANG (sudah bersih)
modal_benar   = anchor_cost * faktor_jual
```
Lalu tulis `transaction_items.purchase_price = modal_benar` (dengan audit nilai lama).

Kenapa metode ini benar & akurat:
- `product_packages` SEKARANG sudah bersih/rapi (master satuan sudah diperbaiki) — jadi
  anchor_cost + faktor resmi menghasilkan modal per satuan jual yang konsisten.
- Terbukti: Surya Pack (tersimpan campur 352.500 salah & 35.000) → semua jadi 35.250 konsisten;
  Djarum Pack → 15.250; LA Bold Pack → 37.400; Layangan Pieces → 1.400.
- **WAJIB pakai ResolvePackageFactor rantai PENUH (multi-level)**, bukan 1-level. Bukti:
  Surya "Batang" perlu Batang→Pack→Slop (faktor 0,00625 → modal 2.203), kalau cuma 1-level
  (0,0625) hasil 22.031 (salah). Skrip HARUS meniru model.ResolvePackageFactor persis.

### Kenapa TIDAK pakai last-cost dari purchase_items
Sudah dibuktikan `purchase_items` sendiri KOTOR (conversion_qty tak konsisten, 6 produk
variasi harga anchor ekstrem). Mengambil last-cost dari data kotor = mengoreksi dengan data
salah. `product_packages.purchase_price` sekarang adalah sumber yang sudah dibersihkan →
lebih akurat sebagai basis. (Trade-off: tidak menangkap perubahan harga beli antar-waktu,
mis. Surya 350rb→352,5rb. Kalau user mau presisi per-tanggal itu, harus bersihkan
purchase_items dulu — bisa jadi fase lanjutan, tapi untuk sekarang anchor_cost sekarang
sudah jauh lebih akurat dari kondisi rusak saat ini.)

### Cross-check & kriteria "meragukan" → Rekonsiliasi Modal manual
Baris masuk manual (TIDAK dikoreksi otomatis) bila salah satu:
1. Paket jual (`ti.unit_id`) tidak ditemukan / tidak bisa di-resolve fakotornya. (dev: 2 baris)
2. Produk tidak punya paket anchor / anchor_cost = 0.
3. Hasil `modal_benar >= ti.price` (modal ≥ harga jual — indikator masih janggal).
4. Produk `needs_stock_review=1` atau rantai satuan bercabang (ResolvePackageFactor error).

Proyeksi dev (dgn faktor 1-level di probe; angka manual akan LEBIH KECIL dgn multi-level):
- ok otomatis ~1.258/1.307 (96%), manual ~49 (~4%). Setelah multi-level, manual turun lagi.

### Cakupan & filter (WAJIB)
- HANYA `t.status='completed'` (skip 10 void / 12 item void). Sertakan `is_credit=1`.
- Semua sumber COGS (report + business_summary) otomatis ikut benar (formula sama).

### Dampak (dev, periode 1-20 Sep, metode probe 1-level): COGS 14,79 jt → 7,79 jt.
Angka final skrip (multi-level) sedikit beda tapi arah sama (laba jadi wajar/positif).

### Komponen terpisah yang JUGA perlu ditangani (di luar koreksi data transaksi)
- **Jalur BACKDATE** (backdate_repo.go) menyimpan modal salah (products.purchase_price mentah,
  conversion_qty=qty×pkgqty). Kalau fitur dipakai → perbaiki agar konsisten dgn online
  (pakai product_packages.purchase_price × faktor). Kalau tidak, transaksi backdate baru salah lagi.
- **purchase_items kotor**: opsional dibereskan kalau dipakai laporan pembelian/valuasi.
  Untuk koreksi HPP penjualan ini TIDAK dipakai, jadi tidak wajib sekarang.

### Rencana implementasi (urut)
1. Skrip `BE/cmd/backfill_transaction_cogs` — dry-run dulu: cetak before/after, daftar baris
   OK vs MANUAL, total COGS lama vs baru. Pakai ResolvePackageFactor multi-level (impor dari
   domain product/model, JANGAN hitung ulang faktor manual di SQL).
2. Tabel audit `transaction_items_cogs_audit` (old/new/faktor/anchor/alasan/waktu).
3. Verifikasi Laba Rugi & Ringkasan Bisnis di dev (browser) → angka wajar.
4. Menu **Rekonsiliasi Modal** (admin) untuk baris/produk berstatus MANUAL — input modal benar
   manual, update purchase_price + audit. (dikerjakan setelah skrip)
5. Perbaiki jalur BACKDATE.

### STATUS IMPLEMENTASI (20 Sep 2026)
- ✅ Skrip `BE/cmd/backfill_transaction_cogs/main.go` DIBUAT. Default DRY-RUN; tulis hanya
  dengan flag `--apply`. Pakai `product_model.ResolvePackageFactor` (multi-level, rantai penuh).
  Buat tabel audit `transaction_items_cogs_audit` otomatis saat --apply. `go build` bersih.
- ✅ DRY-RUN dijalankan di DEV. Hasil:
  - Total baris completed: 1.307. Dikoreksi otomatis: 1.294. Manual: 13.
  - Baris OK yang nilainya BERUBAH: 582.
  - Total COGS lifetime: 59.844.235 → 16.447.004 (baris manual masih pakai nilai lama).
  - Manual (5 produk): Beras Ikanku 5kg & Beras Hallo Bro (modal≥jual, margin tipis),
    LA Bold 20 & Dunhill Hitam 16 (faktor tak bisa di-resolve — rantai satuan bercabang),
    Hansaplast (anomali modal 500→50.000 vs jual 40.000).
- ✅ VERIFIKASI koreksi mencurigakan → ternyata BENAR:
  - Roti Winny "Kardus" (1 Kardus=50 Bungkus): ada penjualan per Kardus (jual 150.000) tapi
    modal tersimpan cuma 2.700 (harga per Bungkus) → dikoreksi ke 135.000 (BENAR, lolos
    cross-check < 150.000). Ini kasus "jual satuan besar, modal tercatat satuan kecil" —
    kebalikan bug utama, tetap salah, & berhasil dibetulkan.
  - Toppas Batang: 16.000/12 = 1.333 (benar). Faktor >1 utk satuan > anchor adalah WAJAR.
  - Kesimpulan: metode + cross-check sudah akurat. Faktor >1 bukan bug.
- Hasil dry-run lengkap diarsipkan: `testing/cogs_dryrun.txt` (bisa dibuka user).
- BELUM di-apply. Menunggu keputusan user untuk lanjut --apply di dev (setelah backup),
  lalu verifikasi Laba Rugi via browser.

### CATATAN penyempurnaan skrip yang MUNGKIN diinginkan sebelum --apply
- (opsional) Guard tambahan utk faktor sangat besar (mis. >X) sebagai MANUAL — TIDAK perlu
  karena verifikasi menunjukkan faktor besar (Kardus=50) memang valid & sudah dijaga cross-check.
- Baris manual saat ini dibiarkan pakai nilai lama (belum diubah). Menu Rekonsiliasi Modal
  akan menangani ini. Untuk sementara, COGS baris manual masih mengandung nilai lama (sebagian
  mungkin masih salah, mis. LA Bold modal 355.000) → akan dibereskan via menu manual.

---

## 9. Catatan Teknis Lingkungan (untuk sesi berikutnya)

- Testing/DB probe: mysql WAMP di `C:\wamp64\bin\mysql\mysql8.4.7\bin\mysql.exe`,
  DB dev `pos_retail_db`, user `root` tanpa password (config_dev.json).
- Playwright/Chromium + cache/log HARUS di folder `testing` (bukan disk C).
- Terminal PowerShell kadang tampilkan "Exit Code: -1" palsu & gema karakter — abaikan,
  baca output aktual (tulis ke file lalu baca lebih andal). `npm.ps1` diblokir → pakai
  `npm.cmd`. Env var pakai `$env:` bukan `%VAR%`.
- `docs/` di-.gitignore (tidak dilacak git) — file rancangan ini tidak masuk version control.

# Diskusi Perbaikan Menu Pelaporan

> Status: **CATATAN DISKUSI — belum ada implementasi.**
> Tujuan file ini: jadi papan kerja bersama untuk merancang perbaikan menu Pelaporan,
> dibahas & digarap SATU PER SATU per menu. Setiap keputusan final dicatat di sini
> sebelum masuk ke tahap implementasi.
>
> Cakupan: SEMUA menu Pelaporan KECUALI **Rekonsiliasi Stok** (dikecualikan sesuai
> permintaan — menu itu sudah selesai).

---

## 1. Latar Belakang & Tujuan

Owner toko tidak selalu paham istilah finance/akuntansi. Laporan saat ini sudah benar
secara data (kecuali dugaan bug di Laba Rugi, lihat §3.2), tapi sebagian menyajikan angka
dengan bahasa teknis yang menuntut owner sudah mengerti konsepnya. Tujuan perbaikan:
membuat laporan **menjawab pertanyaan owner**, bukan sekadar **menampilkan data mentah**.

Pergeseran pola pikir yang ingin dicapai:
- DARI (data-driven): "Ini tabel lengkap dengan HPP per produk."
- KE (question-driven): "Bulan ini untung bersih Rp X, naik 8% dari bulan lalu."

---

## 2. Prinsip Desain (acuan untuk semua menu)

1. **Benerin data dulu, baru percantik.** Kalau angkanya salah, tampilan cantik justru
   berbahaya (owner percaya data keliru). Kebenaran angka > penyajian.
2. **Bahasa awam dulu, istilah teknis belakangan.** Ganti/dampingi istilah akuntansi
   dengan bahasa sehari-hari. Boleh tetap tampilkan istilah aslinya sebagai tooltip.
3. **Kesimpulan di atas, detail di bawah.** Owner buru-buru cukup baca 1 kalimat ringkas
   paling atas; yang mau menggali bisa scroll ke tabel detail.
4. **Angka butuh pembanding.** Angka sendirian tak bermakna — sertakan pembanding periode
   sebelumnya / persentase perubahan bila memungkinkan.
5. **Perjelas peran tiap menu.** Kurangi kebingungan "mau tahu X buka menu mana".
6. **Jangan ubah kebenaran data tanpa alasan kuat.** Perbaikan penyajian tidak boleh
   mengubah logika perhitungan yang sudah benar.

---

## 3. Keputusan Arsitektur (hapus / gabung / tambah)

Sudah didiskusikan & disepakati:

- **TIDAK ada menu yang dihapus.** Kelima menu menjawab pertanyaan berbeda.
- **TIDAK ada menu yang digabung.** Ada tumpang tindih (Ringkasan Bisnis menampilkan
  angka dari Penjualan + Laba Rugi + Stok), TAPI perannya beda: Ringkasan Bisnis =
  *lihat cepat*, tiga lainnya = *gali detail*. Solusinya perjelas hubungan antar-menu
  (mis. kartu di Ringkasan Bisnis bisa klik → menu detail), bukan digabung.
- **TIDAK ada menu baru ditambah.** Yang dibutuhkan (pembanding periode, kalimat
  kesimpulan) adalah fitur DI DALAM menu yang sudah ada, bukan menu tersendiri.

Kesimpulan: **fokus benahi ISI tiap menu, bukan mengubah daftar menu.**

---

## 4. Temuan dari Eksplorasi Browser (20 Sep 2026)

Ditinjau langsung via browser (login admin, data dev). Screenshot ada di
`testing/screenshots/pelaporan-*.png`.

| Menu | Kondisi tampilan | Temuan penting |
|------|------------------|----------------|
| Penjualan | 3 KPI (Total Transaksi, Total Pendapatan, Rata-rata/Transaksi) + tabel transaksi | Paling rapi & intuitif. Angka wajar. |
| Laba Rugi | Blok Pendapatan / Pengeluaran (HPP + Expense) / Hasil (Laba Kotor, Laba Bersih) | **DUGAAN BUG:** Pendapatan Rp 6.042.500 tapi HPP Rp 12.556.900 → Laba −Rp 6.514.400. HPP tampak dihitung dari SELURUH pembelian periode, bukan modal barang yang TERJUAL. Perlu diverifikasi ke kode. |
| Stok | 3 KPI (Total Produk 316, Stok Rendah 215, Total Nilai Stok Rp 23.709.804) + tabel + badge "Stok Rendah" | **215/316 produk (68%) berstatus stok rendah** → ambang minimum kemungkinan tidak realistis, badge jadi tidak berguna. |
| Kinerja Kasir | Tabel per kasir (Jml Transaksi, Total Penjualan, Tunai, Non-Tunai, Rata-rata, Void) | Jelas & cukup. Istilah "Void" mungkin perlu dijelaskan. |
| Ringkasan Bisnis | 5 kartu KPI (Transaksi, Pendapatan, Laba Kotor, Stok Menipis, Piutang Terbuka) + Grafik Penjualan + Top Produk. Toggle Hari Ini / Minggu Ini / Bulan Ini | Kandidat "pintu utama" owner. Kartu Laba Kotor & Stok Menipis mengambil angka dari menu Laba Rugi & Stok → tergantung perbaikan #1 & #2. |

Istilah teknis yang saat ini tampil di UI: **HPP**, **Laba Kotor**, **Laba Bersih**,
**Nilai Stok**, **Void**, **Pendapatan**.

---

## 5. Urutan Pembahasan (DISEPAKATI)

Logika urutan: **#1–2 benerin data yang salah/menyesatkan dulu, #3–5 percantik penyajian.**

- [ ] **1. Laba Rugi** — PALING PRIORITAS. Dugaan bug hitung HPP (bukan cuma bahasa).
      Risiko tinggi: owner bisa mengira rugi padahal tidak. Semua menu lain yang
      menampilkan "Laba Kotor" ikut salah kalau ini salah.
      → Langkah pertama: TELUSURI KODE perhitungan HPP untuk memastikan dugaan, sebelum
        ubah tampilan.
- [ ] **2. Stok** — cek ambang "stok rendah" yang bikin 68% produk ketrigger.
- [ ] **3. Ringkasan Bisnis** — dashboard/"wajah" owner. Ditaruh setelah #1 & #2 karena
      kartunya mengambil angka dari kedua menu itu. Fokus: pembanding periode + kesimpulan.
- [ ] **4. Penjualan** — sudah rapi. Hanya penyempurnaan (mis. pembanding periode).
- [ ] **5. Kinerja Kasir** — sudah cukup. Paling ringan (mis. jelaskan istilah "Void").

---

## 6. Catatan per Menu

Format: **Kondisi sekarang** → **Observasi/masalah** → **Ide arah** → **Keputusan final**
(diisi saat sudah sepakat).

> STATUS AKHIR LABA RUGI (21 Sep 2026): SELESAI 4 langkah —
> (1) koreksi data modal historis (skrip backfill_transaction_cogs, 582 baris),
> (2) menu Rekonsiliasi Modal (admin, koreksi 13 baris manual),
> (3) perbaikan jalur backdate (modal & conversion_qty & mutation_type),
> (4) percantik tampilan Laba Rugi (kartu kesimpulan untung/rugi + label bahasa awam +
>     penjelasan rumus, di ProfitLossTab.tsx). Laba Rugi kini akurat & ramah non-finance.

### 6.1 Laba Rugi  ← SEDANG DIBAHAS (investigasi SELESAI 20 Sep 2026)

- **Kondisi sekarang:** Blok Pendapatan (Total Pendapatan) → Pengeluaran (Harga Pokok
  Penjualan/HPP + Total Pengeluaran/Expense) → Hasil (Laba Kotor, Laba Bersih). Filter
  tanggal (Hari ini/Minggu ini/Bulan ini/Reset) + Export Excel.

#### Hasil investigasi kode + data (PENTING)

**Dugaan awal (HPP dari pembelian) TERBUKTI SALAH.** Rumusnya sudah benar secara struktur.
Akar masalah sebenarnya adalah **DATA HISTORIS** — snapshot harga modal di transaksi lama
tercatat dengan level satuan yang salah.

**A. Rumus (sudah benar):**
- COGS dihitung dari `transaction_items` (barang TERJUAL), bukan dari `purchases`.
  Query: `SUM(ti.quantity * ti.purchase_price)` di `report_repo.go` → `profitLossQuery`
  (± baris 58-69). Filter `t.status='completed' AND t.transaction_date BETWEEN ? AND ?`.
- Revenue & COGS satu query, basis tanggal sama (tanggal penjualan) → konsisten.
- `GrossProfit = Revenue − COGS`, `NetProfit = Gross − Expenses`
  (`report_service.go` ± baris 81-82).

**B. Cara modal per item disimpan:**
- `transaction_items.purchase_price` = SNAPSHOT modal per satuan jual saat transaksi.
  Diisi dari `product_packages.purchase_price` paket terpilih; fallback
  `products.purchase_price × conversion_qty` bila kosong
  (`transaction_repo.go` ± baris 282-321 online, 560-584 sync).
- TIDAK ada moving-average. Cost = input manual.

**C. Akar masalah (dari query data dev, periode 1-20 Sep 2026):**
- Total COGS Rp 12.556.900 vs Pendapatan Rp 6.042.500 → cocok dgn UI.
- **65 baris "janggal" (modal ≥ harga jual) menyumbang Rp 7.912.900 = 63% dari COGS.**
- Pola janggal: **modal tercatat = harga satuan BESAR, padahal dijual per satuan KECIL.**
  Contoh nyata:
  - Surya 16: dijual per **Pack** (modal benar 35.000) tapi baris lama catat **352.500**
    (itu harga per **Slop**; 1 Slop = 10 Pack).
  - Layangan: dijual per **Pieces** (benar 1.400) tapi catat **140.000** (harga per Pack;
    1 Pack = 100 Pieces).
- Master data `product_packages` SEKARANG sudah benar (pengecekan modal≥jual = 0 baris).

**D. Ini bug HISTORIS, bukan bug aktif:**
- Sebaran baris janggal (lifetime, 447 dari 1.105 baris = 40%):
  - Jul 2026: 1 baris
  - **Agu 2026: 381 baris, Rp 38.932.085** (mayoritas)
  - Sep 2026: 65 baris, Rp 7.912.900
- Riwayat Surya 16 menunjukkan **titik balik ± 8 Sep 2026**: transaksi ≥ 8 Sep modal
  BENAR (35.000), sebelum itu salah (352.500). Artinya logika pencatatan modal sudah
  diperbaiki di suatu titik; sisa masalah = data lama yang terlanjur tersimpan salah.
- Penyebab teknis lama: modal diambil dari harga ANCHOR (satuan besar) TANPA dikali
  `conversion_qty` (mis. 0,1 utk Pack). Kode sekarang sudah pakai
  `product_packages.purchase_price` per paket → benar.

**E. Hubungan dengan menu lain:**
- **Ringkasan Bisnis** hitung COGS dgn pola sama (`transaction_items`) TAPI query terpisah
  (`business_summary_repo.go` → `rangeCOGSQuery`/`monthCOGSQuery`). Jadi **kartu "Laba
  Kotor" di Ringkasan Bisnis kena bug data yang SAMA.**
- Beda halus: Report pakai `t.transaction_date BETWEEN` (tanpa `DATE()`), business_summary
  pakai `DATE(t.transaction_date) BETWEEN` → bisa beda perlakuan jam di batas periode.
- Beda definisi "Pendapatan": Report = `(total_amount − tax)` diprorata per item;
  business_summary = `SUM(total_amount)` penuh. → Laba Kotor bisa beda tipis antar dua menu.
- Dua query COGS terpisah = risiko drift kalau salah satu diubah tanpa yang lain.

- **Observasi/masalah (ringkас):**
  1. **BUG DATA HISTORIS (utama):** modal transaksi lama salah level satuan → COGS
     membengkak → laba tampak rugi padahal tidak. ~447 baris, mayoritas Agustus 2026.
  2. **Bahasa:** "HPP", "Laba Kotor", "Laba Bersih" istilah akuntansi (perlu didampingi
     bahasa awam).
  3. **Penyajian:** belum ada kalimat kesimpulan besar di atas.
  4. **Konsistensi:** COGS & definisi Pendapatan beda antara Laba Rugi vs Ringkasan Bisnis.

- **Ide arah (kandidat, BELUM diputuskan):**
  - Untuk #1 (data historis): perlu skrip koreksi satu kali yang membetulkan
    `transaction_items.purchase_price` baris lama = modal per satuan jual yang benar
    (mis. `products.purchase_price × conversion_qty`, atau ambil dari
    `product_packages.purchase_price` sesuai satuan). HATI-HATI: ubah data historis =
    operasi sensitif, wajib backup + dry-run + verifikasi per baris. Ini fokus utama.
  - Untuk #2 & #3 (bahasa & penyajian): dibahas setelah angka benar.
  - Untuk #4 (konsistensi): pertimbangkan satukan sumber COGS agar dua menu tak drift.

#### G. KEPUTUSAN: Opsi B (hitung ulang modal dari data pembelian nyata) + menu rekonsiliasi manual

**Dipilih user (akurasi diutamakan; data sedikit karena app baru rilis Agustus 2026).**

**Prinsip koreksi (Opsi B):**
Untuk tiap baris `transaction_items` yang modalnya salah, hitung ulang modal per satuan
jual dari HARGA BELI NYATA yang berlaku pada/sebelum tanggal transaksi:
- Sumber: `purchase_items` (punya `purchase_price` per satuan beli + `conversion_qty` + `unit`).
- Harga beli per SATUAN DASAR (anchor) dari sebuah baris pembelian
  = `purchase_items.purchase_price / purchase_items.conversion_qty`.
- Modal per satuan JUAL di transaksi
  = harga_beli_per_anchor × `transaction_items.conversion_qty`.
- "Berlaku saat transaksi" = ambil pembelian dengan `purchase_date` <= tanggal transaksi,
  yang PALING DEKAT (last cost). Kalau belum ada pembelian sebelum transaksi itu, pakai
  pembelian PERTAMA setelahnya (fallback), atau tandai untuk tinjauan manual.

**Kelayakan data (hasil probe DB dev, lifetime):**
- 1.105 baris transaksi total. **1.104 bisa dikoreksi otomatis** dari data pembelian.
- **Hanya 1 baris** produknya tak punya pembelian sama sekali → LA Bold 20 (id 70) →
  masuk daftar "perlu tinjauan manual".
- Rumus `purchase_price/conversion_qty` terbukti hasilkan modal per-anchor yang akurat
  (mostly sama dengan anchor cost sekarang; beda tipis di produk yg harganya memang
  pernah berubah — inilah nilai tambah Opsi B: pakai harga SAAT ITU, bukan harga sekarang).
- Catatan: keputusan "baris mana yang dikoreksi" — apakah HANYA baris janggal (modal≥jual)
  atau SEMUA baris dihitung ulang demi konsistensi — masih perlu difinalkan. Rekomendasi:
  hitung ulang SEMUA baris agar 1 metode konsisten (last-cost), bukan campur.

**Menu baru: Rekonsiliasi Modal/HPP (role admin) — konsep, mengikuti pola Rekonsiliasi Stok:**
- Menampilkan baris/produk yang modalnya TIDAK bisa dihitung otomatis (tak ada data
  pembelian, atau hasil hitung meragukan), analog produk `needs_stock_review`.
- Admin input modal yang benar secara manual → sistem update `transaction_items.purchase_price`
  baris terkait + catat jejak audit (siapa, kapan, dari berapa ke berapa).
- Tujuan: meminimalkan data rawan; koreksi otomatis untuk yang datanya lengkap, manual
  untuk sisa kecil.

**Aturan eksekusi (WAJIB, semua opsi):**
1. Kerjakan & uji di LOKAL/DEV dulu. Production = urusan user.
2. Backup DB sebelum eksekusi.
3. Skrip DRY-RUN dulu (tampilkan before/after tanpa menyimpan) -> verifikasi angka.
4. Idempotent bila memungkinkan; simpan jejak nilai lama (audit) sebelum menimpa.
5. Setelah koreksi: verifikasi Laba Rugi & kartu Laba Kotor Ringkasan Bisnis jadi wajar.

**Dampak simulasi (periode 1-20 Sep, metode kasar anchor×conv):** COGS 12,5 jt -> ±5,18 jt
(laba jadi positif). Opsi B (last-cost per tanggal) akan lebih presisi dari angka ini.

- **Keputusan final:** **Opsi B** + menu Rekonsiliasi Modal manual (admin). Belum
  diimplementasi — tahap berikutnya: rancang detail skrip koreksi (dry-run dulu) &
  desain menu. Perlu difinalkan: (a) koreksi semua baris vs hanya baris janggal;
  (b) aturan last-cost vs fallback; (c) desain menu rekonsiliasi modal.

- **Catatan verifikasi:** investigasi via baca kode (report_repo, report_service,
  business_summary_repo, transaction_repo, migrations) + query DB dev langsung
  (mysql WAMP). Angka UI (COGS 12.556.900 / Revenue 6.042.500) tereproduksi persis.

#### F. Konfirmasi: migrasi/perbaikan STOK TIDAK menyentuh modal historis

Dikonfirmasi (user + baca kode migrasi 004/005 & skrip backfill di `BE/cmd/`):
- Ada migrasi & perbaikan STOK (RENCANA_PERBAIKAN_STOK_PRESISI): pindah stok ke
  `product_packages.stock`, backfill `quantity`/`conversion_qty`/`package_id`.
- **Semua skrip backfill hanya baca `transaction_items.quantity` & `unit_id`** —
  TIDAK ADA yang menyentuh `transaction_items.purchase_price` (kolom modal).
- Kesimpulan: perbaikan itu membetulkan **stok** + logika pencatatan modal **ke depan**
  (sebab transaksi ≥ ~8 Sep modalnya sudah benar), TAPI **modal historis (transaksi lama)
  tidak ikut dikoreksi**. → Stok sekarang benar, tapi **Laba Rugi periode lama masih salah**.
  Dua hal terpisah; perbaikan stok bukan solusi untuk bug modal historis ini.

### 6.2 Stok  ← A+B SELESAI (30 Sep 2026), C menyusul
- **Kondisi sekarang:** 3 KPI (Total Produk, Stok Rendah, Total Nilai Stok) + tabel
  (kode, nama, kategori, satuan, stok saat ini, nilai stok) + badge "Stok Rendah".
  Filter cari + kategori + Export Excel.
- **Observasi/masalah:** 68% produk berstatus stok rendah (ambang tak realistis?).
  "Total Nilai Stok" berpotensi disalahartikan sebagai "uang/untung".

#### Hasil investigasi data (30 Sep 2026)
- 337 produk aktif. **322 (96%) pakai min_stock=5** (default form, tak pernah disetel).
  Tidak ada yang min_stock=0.
- Distribusi stok riil (satuan anchor): 27 habis, 140 sisa 1-2, 65 sisa 3-5 → banyak yang
  memang tipis. Dengan perbandingan `<=` ambang 5, ~207+ produk masuk "rendah".
- **BUKAN bug hitungan** — angkanya benar. Masalahnya: (1) ambang seragam 5 tak realistis
  per barang, (2) "habis" & "menipis" tak dibedakan (satu badge merah untuk semua).
- Pembengkakan ambang ke satuan terkecil (dugaan awal) BUKAN penyebab utama: uji
  perbandingan di satuan anchor tetap ~233 low.

- **Keputusan final (A+B, tanpa ubah data):**
  - **A. Status 3 tingkat:** Habis (merah, sisa 0) / Menipis (kuning, 0<sisa<=ambang) /
    Aman (tanpa badge). Produk data-diragukan = status "" (tanpa alarm).
    Backend: `StockSummary.Status` di `stock_delta.go`; diteruskan via
    `stock_read_repo.go`; DTO `StockItem.stock_status` + `StockSummary.out_of_stock_count`
    & `low_only_count` (`dto_report.go`); diisi & difilter di `report_repo.go`;
    `report_service.go` + export Excel (Habis/Menipis/Aman). `is_low_stock` &
    `low_stock_count` lama DIPERTAHANKAN (habis+menipis) → menu lain tak berubah.
  - **B. Filter cepat status:** dropdown "Status Stok" (Semua/Habis/Menipis/Aman) di
    `StockReportFilterBar.tsx`; kartu ringkasan jadi 4 (Total Produk, Stok Habis merah,
    Stok Menipis kuning, Total Nilai Stok) di `StockReportSummaryCard.tsx`; badge 3 warna
    di `StockReportTableColumns.tsx`; tipe di `stock.types.ts`.
  - Terverifikasi browser (30 Sep): kartu Habis 27 / Menipis 207, badge kuning/merah,
    filter Habis → 27 data semua badge merah. BE build & FE type-check bersih.
- **Poin C (menyusul, dibahas terpisah):** ambang min_stock seragam 5 belum ideal.
  Opsi: setel per barang lebih mudah, atau hitung otomatis dari rata-rata penjualan.
  BELUM dikerjakan — sesuai permintaan user, dibahas spesifik setelah A+B beres.

### 6.3 Ringkasan Bisnis  ← P1+P2+P3 SELESAI (30 Sep 2026)
- **Kondisi sekarang:** 5 kartu KPI (Transaksi, Pendapatan, Laba Kotor, Stok Menipis,
  Piutang Terbuka) + Grafik Penjualan + Top Produk Terlaris. Toggle periode
  Hari Ini/Minggu Ini/Bulan Ini.

#### Hasil investigasi (30 Sep 2026)
- Laba Kotor & COGS dashboard baca `transaction_items.purchase_price` yang SAMA dgn
  Laba Rugi → sudah ikut terkoreksi backfill HPP. Terverifikasi: Laba Kotor dashboard
  Bulan Ini = Rp 1.035.903, SAMA dengan menu Laba Rugi.
- Stok Menipis pakai BuildStockSummaries/ComputeStockSummary yang sama dgn Laporan Stok,
  tapi hanya baca flag IsLowStock (habis+menipis digabung).
- Temuan: (1) kartu "Stok Menipis" campur habis+menipis; (2) TIDAK ada pembanding periode
  sama sekali (cuma angka telanjang); (3) istilah "Laba Kotor" tanpa penjelasan;
  (4) definisi "Pendapatan" beda dgn Laba Rugi HANYA jika ada pajak (dashboard pakai
  total_amount penuh, Laba Rugi total_amount−tax).

- **Keputusan final (P1+P2+P3 dikerjakan; P4 pajak di-SKIP atas keputusan user):**
  - **P1 Pembanding periode:** tiap kartu Transaksi/Pendapatan/Laba Kotor kini punya
    indikator "▲/▼ x% vs kemarin|minggu lalu|bulan lalu". Backend: `previousPeriodRange`
    (today=kemarin, week=7hr sebelum, month=bulan kalender lalu) + `PrevStats` di
    `dto_business_summary.go`; `GetStats` isi Prev via GetStatsByRange+GetCOGSByRange.
    FE: DeltaLine di `SummaryCards.tsx` (hijau naik/merah turun, "belum ada data" bila
    periode lalu kosong).
  - **P2 Pisah Habis/Menipis:** kartu jadi 6 (Transaksi, Pendapatan, Laba Kotor,
    Stok Habis merah, Stok Menipis kuning, Piutang). Backend repo `GetLowStockCount`
    diganti `GetStockStatusCounts()(out,low)` (pakai s.Status); DTO tambah
    OutOfStockCount+LowOnlyCount. Konsisten dgn Laporan Stok (27 habis / 207 menipis).
  - **P3 Perjelas istilah:** subteks "pendapatan − modal, belum potong biaya" di kartu
    Laba Kotor.
  - Terverifikasi browser (30 Sep, Bulan Ini): 3 indikator % (35%/9%/90%), subteks &
    kartu Habis/Menipis tampil. BE build & FE type-check bersih.
- **P4 (pajak) — DITUNDA:** samakan definisi Pendapatan dgn Laba Rugi HANYA relevan bila
  toko pakai pajak. Belum dicek apakah toko pakai pajak; di-skip atas permintaan user.

### 6.4 Penjualan  ← SELESAI (30 Sep 2026)
- **Kondisi sekarang:** 3 KPI (Total Transaksi, Total Pendapatan, Rata-rata/Transaksi) +
  tabel transaksi (tanggal, kode, kasir, customer, total, metode bayar, status). Filter
  tanggal + metode bayar + Export Excel.
- **Observasi/masalah:** paling intuitif; angka wajar. Kandidat: tambah pembanding periode.

- **Keputusan final (pembanding periode dikerjakan; diskon/pajak di-SKIP):**
  - Tiap kartu (Total Transaksi, Total Pendapatan, Rata-rata/Transaksi) kini punya
    indikator "▲/▼ x% vs periode lalu". Periode pembanding = durasi SAMA dengan rentang
    tanggal terpilih, tepat sebelum periode ini (mengikuti pola previousPeriod Laba Rugi).
    Filter metode bayar/kasir ikut diterapkan ke periode pembanding (apple-to-apple).
  - Backend: `SalesSummary` DTO tambah PrevTransactions/PrevRevenue/PrevAvgPerTx/
    PrevAvailable (dto_report.go); `GetSalesSummaryWithFilters` (report_repo.go) hitung
    prev via helper `prevSalesRange` + query ulang salesSummaryBase dgn kondisi sama.
  - FE: sales.types.ts tambah field prev; SalesReportSummaryCard.tsx pakai DeltaLine
    (pola sama Ringkasan Bisnis).
  - Terverifikasi browser (30 Sep, 1-30 Sep): 3 indikator (34%/7%/40% vs periode lalu).
    BE build & FE type-check bersih.
  - Diskon & Pajak (sudah dihitung backend tapi tak ditampilkan) SENGAJA tidak
    ditambahkan — opsional & tergantung apakah toko pakai pajak (belum dikonfirmasi).

### 6.5 Kinerja Kasir  ← SELESAI (30 Sep 2026)
- **Kondisi sekarang:** tabel per kasir (Jml Transaksi, Total Penjualan, Tunai, Non-Tunai,
  Rata-rata/Transaksi, Void). Filter tanggal + Export Excel.
- **Observasi/masalah:** cukup jelas & angka wajar. Menu paling ringan. Satu-satunya:
  istilah "Void" teknis untuk owner awam. (Toko saat ini 1 kasir aktif.)

- **Keputusan final (perjelas istilah saja):**
  - Kolom "Void" → **"Dibatalkan"** di tabel (CashierPerformanceTableColumns.tsx) +
    tooltip "Jumlah transaksi yang dibatalkan (void)". Header Export Excel juga
    disamakan jadi "Dibatalkan" (report_service.go ExportCashierReport).
  - TIDAK ditambah kartu ringkasan/grafik/pembanding periode — sengaja, karena menu ini
    sifatnya tabel pembanding antar kasir (per-kasir), bukan dashboard angka tunggal.
  - Terverifikasi browser (30 Sep): kolom "Dibatalkan" tampil, "Void" hilang. BE build &
    FE type-check bersih.

---

## 7. Catatan Lintas-Menu (peran & navigasi)

Potensi kebingungan: Penjualan vs Laba Rugi vs Ringkasan Bisnis sama-sama menampilkan
angka penjualan/laba dari sudut berbeda. Pembagian peran yang diinginkan:
- **Ringkasan Bisnis** = pandangan cepat / dashboard (pintu utama owner).
- **Penjualan** = detail transaksi (omzet/pendapatan).
- **Laba Rugi** = untung setelah modal & biaya.
- **Stok** = kondisi barang sekarang.
- **Kinerja Kasir** = performa orang.

---

## 8. Catatan Penting (aturan kerja)

- File ini murni catatan diskusi. **Tidak ada perubahan kode** sampai keputusan final tiap
  menu disepakati.
- Sebelum mengubah sebuah menu, **baca kode terkait dulu** (FE & BE) untuk memastikan apa
  yang SUDAH ada — supaya tidak salah asumsi.
- Testing lewat browser pakai Playwright/Chromium dari folder `testing`. Cache, instalasi,
  session **HARUS di folder testing**, TIDAK boleh di disk C.
- Log lingkungan yang sudah diketahui: terminal PowerShell kadang menampilkan
  "Exit Code: -1" palsu & gema karakter (abaikan, baca output aktual); `npm.ps1` diblokir
  execution policy → pakai `npm.cmd`; env var pakai `$env:` bukan `%VAR%`.

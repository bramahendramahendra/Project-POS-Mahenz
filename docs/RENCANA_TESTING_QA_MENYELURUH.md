# Rencana Testing QA Menyeluruh — Modul Stok & Alur Terkait

> Dokumen prompt eksekusi (bukan hasil testing). Dibuat 16 Agu 2026, menyusul selesainya Fase 0-8 di `docs/RENCANA_PERBAIKAN_STOK_PRESISI.md` (perbaikan presisi stok multi-satuan, semua di DB dev).
>
> Tujuan: uji ULANG seluruh sistem (bukan cuma modul stok) dengan mindset **senior QA / bug hunter** — bukan cuma jalur bahagia (happy path), tapi sengaja mencoba mematahkan sistem: input tidak valid, kondisi batas (boundary), race condition, klik ganda, refresh di tengah proses, dan kombinasi skenario yang belum pernah dicoba sebelumnya. Semua lewat **browser sungguhan** (Playwright), bukan cuma curl API.

## Cara pakai dokumen ini

- Dipecah jadi beberapa fase, tiap fase adalah satu sesi kerja terpisah — **jangan lompat fase**, selesaikan & laporkan hasil fase sebelumnya dulu sebelum lanjut.
- Tiap fase punya blok prompt siap-salin (di bagian "PROMPT — Fase N"). Salin ke sesi AI baru (atau lanjutkan di sesi yang sama).
- **Setiap temuan bug** (baik yang berhasil diperbaiki maupun yang cuma dicatat) WAJIB dilaporkan dengan format: langkah reproduksi, hasil aktual, hasil yang diharapkan, tingkat keparahan (kritis/mayor/minor/kosmetik), dan status (diperbaiki/belum).
- Kalau nemu bug yang butuh keputusan desain/bisnis (bukan bug kode murni), TANYA dulu ke user sebelum "memperbaiki" dengan asumsi sendiri.
- Data uji: gunakan produk/transaksi baru dengan penanda jelas di nama (mis. prefix `QA-TEST-`) supaya gampang dibedakan dari data asli, dan supaya gampang dibersihkan setelah selesai kalau diminta. **Jangan hapus data uji secara sepihak** — biarkan di DB dev sampai user minta dibersihkan (lihat memori proyek soal ini).
- Setiap fase WAJIB: jalankan skenario **berhasil** dan **gagal** untuk tiap alur, bukan cuma satu arah. Termasuk mencoba hal yang SEHARUSNYA ditolak sistem (validasi bekerja), bukan cuma yang seharusnya diterima.
- Cek console browser (devtools) sepanjang SETIAP skenario — bukan cuma di akhir. Screenshot tiap temuan bug.

---

## Fase A — Modul Produk (CRUD, satuan, badge, validasi)

**Cakupan:** Tambah/Edit/Hapus/Nonaktifkan produk, kelola satuan (grosiran), badge `needs_stock_review`, validasi form.

**Skenario wajib dicoba:**
1. Tambah produk baru — satuan tunggal, stok awal 0. Sukses.
2. Tambah produk baru — satuan tunggal, stok awal > 0 (mis. 15). Sukses, cek `stock_mutations` tercatat `reference_type='product_create'`.
3. Tambah produk baru — multi-satuan (anchor + 2 grosiran berjenjang, mis. Kardus→Pack→Pieces) sekaligus dalam satu form. Sukses, cek breakdown di Detail Produk benar.
4. Tambah produk — coba submit dengan field wajib kosong (nama, kategori, satuan dasar, harga jual). Harus ditolak dengan pesan jelas per field, BUKAN silent fail atau 500.
5. Tambah produk — barcode/SKU/nama yang sudah dipakai produk lain. Harus ditolak dengan pesan spesifik.
6. Tambah produk — harga jual < harga beli (margin negatif). Cek: diizinkan tapi ada peringatan visual, atau ditolak? (verifikasi perilaku aktual, laporkan kalau tidak sesuai ekspektasi bisnis).
7. Grosiran: tambah paket baru dengan rasio 1 (mis. 1 Pack = 1 Pack, ref ke diri sendiri secara logis). Harus ditolak.
8. Grosiran: tambah paket yang menciptakan referensi melingkar (A→B→A). Harus ditolak dengan pesan jelas (`wouldCreateCycle`).
9. Grosiran: tambah paket dengan qty/ref_qty desimal aneh (0.001, atau angka sangat besar 999999). Cek breakdown & harga per unit tetap masuk akal, tidak overflow/NaN.
10. Grosiran: hapus paket yang masih dirujuk paket lain (mis. hapus Pack padahal Sachet merujuk ke Pack). Harus ditolak.
11. Grosiran: hapus paket anchor. Harus ditolak (tombol harusnya tidak ada/disabled).
12. Edit produk — ubah stok manual (field "Stok") naik & turun, sebagai admin. Cek `stock_mutations` tercatat `mutation_type='adjustment'`, `user_id` benar.
13. Edit produk — ubah stok manual sebagai role NON-admin (kasir/staff kalau ada akun test). Field harus read-only/terkunci, submit tidak boleh mengubah stok walau dipaksa lewat DevTools mengubah value input.
14. Edit produk — turunkan stok manual sampai MELEBIHI stok yang tersedia (misal stok 5, reserved 3 karena retur pending, coba set ke 0). Cek: ditolak atau diterima? Kalau retur pending menahan sebagian, turun ke bawah reserved harus ditolak.
15. Produk dengan `needs_stock_review=true`: buka Detail & Edit, pastikan badge & catatan tampil. Klik "Tandai Sudah Ditinjau" — cek badge hilang, flag di DB benar-benar `0`. Coba lakukan operasi stok (jual/beli) pada produk ini SEBELUM ditandai ditinjau — harus ditolak eksplisit (`ErrNeedsStockReview`), bukan diproses diam-diam.
16. Nonaktifkan produk yang masih ada stok. Cek: diizinkan (sesuai keputusan desain celah #17), tapi tetap muncul di "produk nonaktif dengan stok" kalau ada laporan itu. Cek juga produk nonaktif tidak muncul di pencarian Kasir.
17. Hapus produk yang belum pernah ada transaksi/pembelian. Sukses.
18. Hapus produk yang sudah pernah dibeli/terjual. Harus ditolak dengan pesan menyarankan nonaktifkan.
19. Hapus produk yang masih ada stok (>0). Harus ditolak.
20. Import produk massal (Excel) — file valid. Sukses, cek breakdown & harga masuk semua.
21. Import produk massal — file dengan baris invalid dicampur baris valid (mis. barcode duplikat, harga negatif). Cek preview menunjukkan mana yang error, baris valid tetap bisa diimpor terpisah dari yang gagal.
22. Cari produk (search & filter kategori/status/stok menipis) dengan kombinasi filter sekaligus. Cek hasil & paginasi konsisten.
23. Sortir tabel produk by Stok (asc/desc) — pastikan urutan benar untuk produk dengan breakdown campuran (bandingkan manual beberapa baris).
24. Cetak label produk — cek harga & barcode yang tercetak sesuai data terbaru (bukan cache basi).

---

## Fase B — Modul Pembelian (Supplier Purchase)

**Cakupan:** Create/Update/Void/Add Items, validasi qty & harga, status pembayaran, expired batch.

**Skenario wajib dicoba:**
1. Buat PO baru — 1 item, 1 satuan, lunas. Sukses, stok naik benar.
2. Buat PO baru — multi item, campuran satuan berbeda per item (termasuk produk dengan >1 satuan, pilih satuan bukan default). Sukses, cek breakdown tiap produk benar.
3. Buat PO — produk dengan HANYA 1 satuan (tidak ada dropdown satuan tampil di FE). Sukses, cek BE fallback ke anchor package benar (celah #15).
4. Buat PO — qty desimal untuk satuan diskrit (mis. beli 2.5 Pack). Cek: ditolak (harus integer untuk diskrit) atau diterima? Verifikasi validasi.
5. Buat PO — qty 0 atau negatif. Harus ditolak validasi form, tidak sampai ke server.
6. Buat PO — harga beli 0 atau negatif. Cek perilaku (harusnya minimal ditolak untuk negatif).
7. Buat PO dengan expired date — isi beberapa rincian tanggal expired yang totalnya TIDAK SAMA dengan qty item (mis. qty 10, expired cuma dialokasikan 7). Cek validasi "Total qty harus sama dengan X" mencegah submit.
8. Buat PO dengan expired date di masa lalu (kadaluarsa sejak sebelum dibeli). Cek: diizinkan (kasus nyata bisa terjadi, barang expired diketahui belakangan) atau ada peringatan.
9. Buat PO — pilih supplier lalu batal pilih (kosongkan), coba submit. Harus ditolak validasi.
10. Status pembayaran "Lunas" — cek paid_amount otomatis = total, tidak bisa diubah manual jadi kurang dari total sambil status tetap Lunas.
11. Status pembayaran "Sebagian" — isi paid_amount < total. Cek sisa hutang terhitung benar, status badge "Bayar Sebagian".
12. Status pembayaran "Hutang" — paid_amount 0. Cek sisa hutang = total.
13. Bayar PO yang berstatus Hutang/Sebagian lewat tombol bayar terpisah (kalau ada) — cek update status & sisa hutang benar setelah dibayar penuh.
14. Edit PO — ubah qty item existing naik. Cek stok bertambah sesuai delta (bukan replace absolut), `stock_mutations` tercatat.
15. Edit PO — ubah qty item existing turun sampai membuat stok produk itu MINUS (karena sebagian sudah terjual di Kasir sejak PO dibuat). Harus DITOLAK dengan pesan jelas, PO tidak boleh tersimpan setengah-setengah.
16. Edit PO — tambah item baru ke PO yang sudah ada (`AddItems`). Sukses, breakdown stok baru benar.
17. Edit PO — hapus salah satu item dari PO existing (kalau fitur ini ada). Cek stok direverse dengan benar.
18. Void PO yang berstatus Lunas. Cek SEMUA stok dari SEMUA item PO itu kembali persis ke sebelum PO dibuat (uji dengan produk rasio non-bulat, cek tidak ada drift +0.001 dst).
19. Void PO yang SUDAH di-void sebelumnya (klik void 2x / lewat 2 tab). Harus ditolak dengan pesan bisnis rapi, BUKAN error 500 atau stok berkurang dobel.
20. Void PO yang sebagian stoknya sudah terjual di Kasir sejak PO dibuat (stok sekarang < qty PO). Cek: void tetap jalan (boleh minus sementara?) atau ditolak — verifikasi & laporkan perilaku aktual, ini kasus tepi penting yang harus jelas aturannya.
21. Void PO yang sebagian item-nya sudah kena write-off expired. Cek stok yang direverse benar (idealnya net dari yang sudah dimusnahkan — sudah dicatat sebagai limitasi minor di dokumen Fase 4, verifikasi ulang perilakunya sekarang).
22. Hapus PO berstatus "active" (belum di-void). Harus ditolak, minta void dulu.
23. Hapus PO yang sudah di-void. Sukses.
24. Generate kode PO otomatis — buka form 2 tab bersamaan, submit keduanya. Cek tidak ada duplikat kode PO (race condition).
25. Filter/sort daftar pembelian (tanggal, status, supplier) kombinasi filter. Cek hasil benar.

---

## Fase C — Modul Kasir / Penjualan

**Cakupan:** Search produk, tambah ke keranjang, checkout, void transaksi, diskon/pajak, metode pembayaran, kredit/piutang.

**Skenario wajib dicoba:**
1. Jual 1 item satuan anchor. Sukses, struk benar, stok berkurang benar.
2. Jual 1 item satuan NON-anchor (mis. Botol dari produk yang anchornya Kardus, rasio non-bulat 1/24 atau 1/3). Sukses, cek `stock_mutations.quantity` presisi penuh (bukan dibulatkan).
3. Jual multi-item, multi-satuan berbeda per item, dalam 1 transaksi. Sukses.
4. Jual qty melebihi stok tersedia (termasuk yang tertahan retur). Harus ditolak dengan pesan ramah SEBELUM atau SAAT checkout, bukan struk keluar dulu baru gagal.
5. Jual produk yang stoknya PAS HABIS sampai 0 (uji boundary, bukan cuma "kurang dari"). Cek berhasil pas ke 0, item berikutnya yang sama harus ditolak (stok 0).
6. Jual produk dengan `needs_stock_review=true`. Harus ditolak eksplisit (`ErrNeedsStockReview`).
7. Jual produk yang barusan dinonaktifkan (Aktif→Nonaktif) — cek tidak lagi muncul di pencarian Kasir sama sekali.
8. Keranjang: ubah qty item via tombol +/- dan lewat input manual. Cek subtotal/total ikut update real-time, tidak nyangkut.
9. Keranjang: set qty ke 0 lewat tombol minus atau input manual. Cek item terhapus dari keranjang atau minimal tidak bisa checkout dengan qty 0.
10. Keranjang: hapus 1 item dari beberapa item. Cek total terhitung ulang benar.
11. Kosongkan keranjang (tombol "Kosongkan"). Cek semua item hilang, tidak ada sisa state nyangkut kalau lanjut tambah produk baru.
12. Diskon: diskon persen (%) dan diskon rupiah (Rp) — coba keduanya, cek kalkulasi benar termasuk pembulatan.
13. Diskon: diskon > 100% atau diskon Rp > subtotal (bikin total negatif). Harus ditolak/di-clamp, tidak boleh total transaksi negatif.
14. Pajak: isi pajak %, cek kalkulasi total benar bersamaan dengan diskon (urutan hitung diskon dulu atau pajak dulu — verifikasi konsisten dengan struk).
15. Pembayaran Tunai — bayar PAS (uang pas). Sukses, kembalian 0.
16. Pembayaran Tunai — bayar KURANG dari total. Tombol proses harus disabled/ditolak.
17. Pembayaran Tunai — bayar LEBIH, cek kembalian dihitung benar (termasuk pembulatan rupiah kalau ada aturan pembulatan kasir).
18. Pembayaran Transfer/QRIS/Kartu — cek field yang relevan muncul/tidak muncul sesuai metode, tidak minta "jumlah bayar tunai" kalau metode non-tunai (atau verifikasi memang tetap minta, sesuai desain).
19. Pembayaran Kredit (kalau tersedia, terhubung ke Piutang) — cek transaksi tercatat sebagai piutang pelanggan, muncul di menu Piutang dengan jumlah benar.
20. Tambah pelanggan ke transaksi (checkbox "Tambah Pelanggan") — cek nama pelanggan tersimpan & tampil di struk/riwayat.
21. Scan barcode produk 1-satuan — langsung masuk keranjang tanpa dialog tambahan.
22. Scan barcode produk multi-satuan — harus tampil pilihan satuan, tidak langsung asal pilih satuan pertama.
23. Scan barcode yang TIDAK terdaftar. Pesan error jelas, tidak crash.
24. Checkout, lalu langsung klik tombol Bayar 2x cepat berturut-turut (double-submit / race condition). Cek tidak menghasilkan 2 transaksi/2x pengurangan stok.
25. Void transaksi yang sudah selesai. Cek stok kembali PERSIS (uji produk rasio non-bulat, pastikan tidak ada drift dibanding sebelum jual — celah #21).
26. Void transaksi yang SUDAH di-void sebelumnya. Harus ditolak rapi.
27. Multi-item checkout dengan item ke-2 sengaja dibuat gagal (stok tidak cukup pas detik terakhir, mis. dijual dari device lain barengan). Cek item pertama TIDAK ikut ter-commit (atomicity, sudah pernah dibuktikan di Fase 4 — verifikasi ulang masih benar).
28. Cetak struk / preview struk — cek semua angka di struk cetak sama persis dengan yang di layar (tidak ada field yang salah mapping).
29. Refresh halaman Kasir di tengah proses checkout (sebelum klik Bayar final). Cek tidak ada transaksi ganda/setengah jalan.

---

## Fase D — Modul Retur ke Supplier

**Cakupan:** Create retur dari PO, approve/reject, reservasi stok.

**Skenario wajib dicoba:**
1. Buat retur dari PO — pilih 1 item, qty penuh sama dengan qty PO. Sukses, `reserved_qty` naik benar di level package yang tepat.
2. Buat retur — qty SEBAGIAN dari qty PO (mis. beli 10, retur 3). Sukses.
3. Buat retur — qty MELEBIHI qty yang dibeli di PO itu (atau melebihi stok yang tersisa kalau sebagian sudah terjual). Harus ditolak.
4. Buat retur — pilih multi item sekaligus dari 1 PO. Sukses, semua ter-reserve benar.
5. Buat retur dari PO yang SUDAH di-void. Cek: field PO itu masih bisa dipilih atau tidak? Kalau bisa, apa yang terjadi — harus ditolak atau ditangani jelas.
6. Buat retur, lalu SEBELUM disetujui, coba jual habis stok produk itu di Kasir sampai ke titik yang membuat qty reserved tidak lagi valid secara fisik. Cek `ApplyStockDelta` menolak jual yang akan "membobol" reservasi (celah #13) — pool bebas = stock - reserved, bukan stock mentah.
7. Approve retur — sukses, cek urutan release-reservasi-lalu-kurangi-stok benar (celah #13/#4), `stock_mutations` tercatat `reference_type='supplier_return'`.
8. Reject retur — cek reservasi dilepas TAPI stok TIDAK berkurang (beda dari approve).
9. Approve retur yang statusnya SUDAH disetujui/ditolak sebelumnya (klik 2x / race). Harus ditolak rapi.
10. Retur untuk produk yang sudah dinonaktifkan setelah retur dibuat tapi sebelum di-approve. Cek approve tetap jalan benar (produk nonaktif ≠ tidak bisa punya mutasi stok).
11. Filter/lihat detail retur — cek breakdown item, alasan, status tampil benar di semua state (Pending/Disetujui/Ditolak).

---

## Fase E — Write-off Kadaluarsa

**Cakupan:** Deteksi batch expired, konfirmasi aman, musnahkan.

**Skenario wajib dicoba:**
1. Batch yang sudah lewat tanggal expired — cek badge "Expired" muncul di tabel produk & modal menampilkan qty+satuan+tanggal yang benar.
2. Batch yang MENDEKATI expired (belum lewat, tapi dalam window peringatan) — cek badge "Mendekati Expired" beda dari "Expired" (kalau ada pembedaan ini).
3. "Sudah Dicek, Aman" (confirm) — cek status batch berubah `cleared`, TIDAK mengubah stok sama sekali (beda dari write-off).
4. "Musnahkan" (write-off) — qty penuh sama dengan batch. Sukses, stok berkurang benar, status `written_off`.
5. Write-off untuk qty yang MELEBIHI stok tersedia saat ini (karena sebagian sudah terjual sejak dibeli). Harus DITOLAK eksplisit (gap #A — tolak, jangan clamp senyap), batch tetap `active`, stok sama sekali tidak tersentuh.
6. Write-off batch yang statusnya SUDAH `written_off`/`cleared` sebelumnya. Harus ditolak, modal seharusnya tidak lagi menampilkan tombol aksi untuk batch itu.
7. Produk dengan BEBERAPA batch expired sekaligus (beda tanggal) — musnahkan salah satu, cek batch lain tidak ikut terpengaruh.
8. Cek nama satuan yang tampil di modal & riwayat batch sesuai satuan asli pembelian (bukan placeholder "unit" — regresi celah lama, pastikan masih benar).

---

## Fase F — Laporan & Dashboard

**Cakupan:** Laporan Stok, Laporan Penjualan, Laba Rugi, Ringkasan Bisnis, konsistensi lintas laporan.

**Skenario wajib dicoba:**
1. Laporan Stok — filter kategori, pencarian nama, sortir tiap kolom (termasuk Stok Saat Ini & Nilai Stok yang sekarang dihitung in-memory sejak Fase 8 — pastikan urutan tetap benar untuk dataset penuh, bukan cuma 1 halaman).
2. Laporan Stok — export Excel, cek angka di file sama dengan yang di layar.
3. Bandingkan "Total Nilai Stok" di Laporan Stok dengan hitung manual (SUM stok anchor × harga beli) untuk beberapa produk sampel.
4. Dashboard/Ringkasan Bisnis — `low_stock_count` dibandingkan dengan jumlah baris "Stok Rendah" di Laporan Stok. Boleh beda (scope filter beda), tapi harus bisa dijelaskan kenapa, bukan angka acak.
5. Laporan Penjualan — filter tanggal termasuk kasus rentang 1 hari, rentang yang tidak ada transaksi sama sekali (cek tidak error, tampil kosong dengan rapi).
6. Laporan Laba Rugi — cek HPP dihitung dari `purchase_price` SNAPSHOT saat transaksi (bukan harga beli sekarang) — uji dengan produk yang harga belinya pernah diubah setelah transaksi lama terjadi.
7. Laporan Kinerja Kasir — cek total transaksi & void count per kasir benar untuk kasir dengan campuran transaksi sukses+void.

---

## Fase G — Lintas Modul: Kasus Tepi & Konkurensi

**Cakupan:** Skenario yang menyentuh lebih dari satu modul sekaligus, race condition, kasus struktural (rantai bercabang, dll).

**Skenario wajib dicoba:**
1. **Concurrency**: buka 2 tab browser, jual produk yang SAMA dengan stok pas-pasan (mis. stok 1) dari kedua tab hampir bersamaan. Cek cuma 1 yang berhasil, yang lain ditolak rapi — TIDAK boleh stok jadi minus.
2. **Concurrency**: 2 tab — 1 tab jual habis stok, tab lain barengan coba retur/write-off produk yang sama. Cek locking `FOR UPDATE` mencegah race (Aturan Operasional #3).
3. Produk dengan struktur `ref_package_id` BERCABANG (celah #10, kalau masih ada produk seperti ini di data dev) — coba jual/beli/retur/write-off. SEMUA harus ditolak eksplisit dengan pesan jelas, tidak ada yang lolos diam-diam.
4. Produk satuan kontinu (Kilogram/Gram/Galon/Gelas) — jual dengan qty desimal (mis. 0.5 Kg). Harus DIIZINKAN (beda dari satuan diskrit), breakdown tetap benar.
5. Produk satuan kontinu — jual qty desimal SANGAT presisi (mis. 0.1234 Kg). Cek tidak ada pembulatan tak terduga di penyimpanan.
6. Alur penuh lintas hari: beli hari ini, jual besok (kalau ada cara ubah tanggal transaksi), retur minggu depan. Cek `stock_mutations` tetap urut berdasarkan `created_at`, bukan tanggal transaksi manual.
7. Produk yang dibuat, langsung diedit paketnya berkali-kali (ubah rasio grosiran) SEBELUM ada transaksi apa pun. Cek breakdown selalu konsisten dengan rasio TERBARU, tidak ada sisa cache rasio lama.
8. Ubah rasio grosiran produk yang SUDAH punya stok & riwayat transaksi. Cek breakdown re-kalkulasi dengan rasio baru — apakah stok existing ikut "terjemahkan ulang" atau tetap di angka lama dengan rasio baru cuma berlaku ke depan? (Verifikasi & laporkan perilaku aktual, ini area abu-abu yang belum eksplisit didokumentasikan.)
9. Login sebagai role berbeda (admin vs kasir vs staff, kalau tersedia) — ulangi beberapa skenario kunci (edit stok, void transaksi, hapus produk) untuk pastikan permission per role konsisten dengan yang didesain.
10. Sesi token kedaluwarsa di tengah proses (kalau bisa disimulasikan) — cek tidak ada data korup, user diarahkan re-login dengan bersih.

---

## Format Laporan Bug (dipakai di semua fase)

Untuk SETIAP temuan, catat dengan format ini di ringkasan akhir tiap fase:

```
### [SEVERITY] Judul singkat bug
- **Lokasi**: halaman/komponen/endpoint
- **Langkah reproduksi**: 1. ... 2. ... 3. ...
- **Hasil aktual**: ...
- **Hasil yang diharapkan**: ...
- **Screenshot**: (path file)
- **Status**: Diperbaiki / Belum diperbaiki (tunggu keputusan user) / Dicatat sebagai limitasi
```

Severity: **Kritis** (data rusak/hilang, stok salah, transaksi ganda) · **Mayor** (fitur gagal total, error 500) · **Minor** (pesan error kurang jelas, validasi longgar tapi tidak merusak data) · **Kosmetik** (tampilan saja).

---

## Prompt Eksekusi per Fase

Salin blok yang sesuai ke sesi kerja. Jangan lompat fase.

---

#### PROMPT — Fase A (Modul Produk)
```
Baca docs/RENCANA_TESTING_QA_MENYELURUH.md di project d:\Develop\Project_pos_mahenz secara lengkap.

TUGAS — Kerjakan FASE A (Modul Produk) sebagai senior QA / bug hunter:
1. Jalankan SEMUA 24 skenario di bagian "Fase A" lewat browser sungguhan (Playwright) — bukan cuma happy path, termasuk skenario yang SEHARUSNYA ditolak sistem.
2. Cek console browser di setiap skenario, screenshot tiap temuan.
3. Untuk tiap bug yang ditemukan: kalau bug kode murni (validasi tidak jalan, error 500, dll) boleh langsung diperbaiki setelah dikonfirmasi lewat reproduksi nyata — tapi kalau butuh keputusan desain/bisnis (mis. "haruskah margin negatif ditolak?"), TANYA ke user dulu, jangan asumsi sendiri.
4. Susun laporan akhir pakai format "Laporan Bug" di dokumen ini, urutkan dari severity tertinggi.
5. JANGAN lanjut ke Fase B sebelum laporan Fase A ditinjau user.
```

---

#### PROMPT — Fase B (Modul Pembelian)
```
Baca docs/RENCANA_TESTING_QA_MENYELURUH.md di project d:\Develop\Project_pos_mahenz secara lengkap. Fase A harus sudah selesai & ditinjau user.

TUGAS — Kerjakan FASE B (Modul Pembelian) sebagai senior QA / bug hunter:
1. Jalankan SEMUA 25 skenario di bagian "Fase B" lewat browser sungguhan, termasuk kasus race condition (poin 24) dan kasus tepi void/edit (poin 15, 18-21) yang butuh verifikasi presisi angka (bandingkan stok sebelum/sesudah manual, bukan cuma percaya tampilan).
2. Perhatikan khusus: poin 20 & 21 (void PO yang sudah kena void/write-off/terjual sebagian) adalah area yang sudah pernah dicatat sebagai limitasi minor sebelumnya — verifikasi ulang perilaku aktualnya sekarang dan laporkan apa adanya, jangan berasumsi sudah beres.
3. Untuk tiap bug: perbaiki kalau jelas bug kode, tanya user kalau butuh keputusan bisnis.
4. Susun laporan akhir, urutkan severity.
5. JANGAN lanjut ke Fase C sebelum laporan Fase B ditinjau user.
```

---

#### PROMPT — Fase C (Modul Kasir/Penjualan)
```
Baca docs/RENCANA_TESTING_QA_MENYELURUH.md di project d:\Develop\Project_pos_mahenz secara lengkap. Fase A & B harus sudah selesai & ditinjau user.

TUGAS — Kerjakan FASE C (Modul Kasir/Penjualan) sebagai senior QA / bug hunter:
1. Jalankan SEMUA 29 skenario di bagian "Fase C" lewat browser sungguhan. Prioritaskan skenario yang menyentuh uang (diskon, pajak, kembalian, poin 12-19) — hitung manual dan bandingkan presisi ke rupiah, dan skenario atomicity/race (poin 24, 27, 29).
2. Untuk poin 27 (atomicity multi-item), gunakan pendekatan yang sama seperti verifikasi Fase 4 sebelumnya: paksa item ke-2 gagal (stok kurang), pastikan item pertama TIDAK ter-commit.
3. Untuk tiap bug: perbaiki kalau jelas bug kode, tanya user kalau butuh keputusan bisnis (terutama soal pembulatan uang/kembalian).
4. Susun laporan akhir, urutkan severity.
5. JANGAN lanjut ke Fase D sebelum laporan Fase C ditinjau user.
```

---

#### PROMPT — Fase D (Modul Retur ke Supplier)
```
Baca docs/RENCANA_TESTING_QA_MENYELURUH.md di project d:\Develop\Project_pos_mahenz secara lengkap. Fase A-C harus sudah selesai & ditinjau user.

TUGAS — Kerjakan FASE D (Modul Retur ke Supplier) sebagai senior QA / bug hunter:
1. Jalankan SEMUA 11 skenario di bagian "Fase D" lewat browser sungguhan. Perhatikan khusus poin 6 (reservasi vs jual bersamaan) — ini celah #13 yang sudah pernah diperbaiki, verifikasi ULANG masih benar dengan skenario browser nyata, bukan cuma percaya dokumentasi lama.
2. Untuk tiap bug: perbaiki kalau jelas bug kode, tanya user kalau butuh keputusan bisnis.
3. Susun laporan akhir, urutkan severity.
4. JANGAN lanjut ke Fase E sebelum laporan Fase D ditinjau user.
```

---

#### PROMPT — Fase E (Write-off Kadaluarsa)
```
Baca docs/RENCANA_TESTING_QA_MENYELURUH.md di project d:\Develop\Project_pos_mahenz secara lengkap. Fase A-D harus sudah selesai & ditinjau user.

TUGAS — Kerjakan FASE E (Write-off Kadaluarsa) sebagai senior QA / bug hunter:
1. Jalankan SEMUA 8 skenario di bagian "Fase E" lewat browser sungguhan. Perhatikan khusus poin 5 (write-off melebihi stok tersedia) — ini gap #A yang sudah diperbaiki (tolak eksplisit, bukan clamp), verifikasi ULANG masih benar.
2. Untuk tiap bug: perbaiki kalau jelas bug kode, tanya user kalau butuh keputusan bisnis.
3. Susun laporan akhir, urutkan severity.
4. JANGAN lanjut ke Fase F sebelum laporan Fase E ditinjau user.
```

---

#### PROMPT — Fase F (Laporan & Dashboard)
```
Baca docs/RENCANA_TESTING_QA_MENYELURUH.md di project d:\Develop\Project_pos_mahenz secara lengkap. Fase A-E harus sudah selesai & ditinjau user.

TUGAS — Kerjakan FASE F (Laporan & Dashboard) sebagai senior QA / bug hunter:
1. Jalankan SEMUA 7 skenario di bagian "Fase F" lewat browser sungguhan. Untuk poin 3 (Nilai Stok), hitung manual minimal 3 produk sampel dan bandingkan presisi ke rupiah dengan angka di layar.
2. Untuk tiap bug: perbaiki kalau jelas bug kode, tanya user kalau butuh keputusan bisnis.
3. Susun laporan akhir, urutkan severity.
4. JANGAN lanjut ke Fase G sebelum laporan Fase F ditinjau user.
```

---

#### PROMPT — Fase G (Lintas Modul: Kasus Tepi & Konkurensi)
```
Baca docs/RENCANA_TESTING_QA_MENYELURUH.md di project d:\Develop\Project_pos_mahenz secara lengkap. Fase A-F harus sudah selesai & ditinjau user. Ini fase paling sulit — fokus ke race condition & kasus struktural yang sengaja dirancang untuk mematahkan sistem.

TUGAS — Kerjakan FASE G (Lintas Modul) sebagai senior QA / bug hunter:
1. Jalankan SEMUA 10 skenario di bagian "Fase G" lewat browser sungguhan — gunakan 2 tab/context browser paralel untuk skenario concurrency (poin 1-2).
2. Poin 3 (produk bercabang) dan poin 8 (ubah rasio produk yang sudah ada transaksi) adalah area abu-abu yang belum pernah eksplisit diuji lewat browser — laporkan perilaku aktual apa adanya, JANGAN asumsi "pasti sudah benar" dari dokumentasi lama.
3. Untuk tiap bug: perbaiki kalau jelas bug kode, tanya user kalau butuh keputusan bisnis/desain (kemungkinan besar beberapa temuan di fase ini butuh keputusan, bukan cuma bug teknis).
4. Susun laporan akhir gabungan SEMUA fase (A-G) sebagai ringkasan eksekutif: total bug ditemukan per severity, berapa yang diperbaiki, apa yang masih terbuka, rekomendasi lanjut/tidak ke tahap berikutnya (mis. rilis production).
```

---

## Status

**Fase A (Modul Produk) — ✅ SELESAI** (18 Agu 2026): semua 24 skenario dijalankan lewat browser sungguhan (Playwright headless Chromium, auth di-inject ke localStorage). Data uji: produk `QA-TEST Produk Satuan Tunggal Stok0` (201, dihapus di skenario 17), `QA-TEST Stok Awal 15` (202, dinonaktifkan skenario 16), `QA-TEST Multi Satuan Berjenjang C` (203), `QA-TEST Import Valid`/`QA-TEST Import Grosir` (204-205, dari skenario import), akun `qatestkasir` (role kasir, dibuat khusus untuk skenario 13 karena sebelumnya tidak ada akun non-admin di data dev). Semua data uji dibiarkan di DB dev (tidak dihapus), kecuali retur pending `RTR-20260818-001` (produk 197, dibuat sebagai fixture untuk skenario 14, sengaja dibiarkan pending — akan relevan lagi untuk Fase D).

### Bug ditemukan & diperbaiki

#### [MAYOR] Edit stok manual produk turun di bawah stok tersedia (reserved) → 500, bukan pesan jelas
- **Lokasi**: `POST /products/update/:id`, dipicu dari form Edit Produk field "Stok"
- **Langkah reproduksi**: 1. Produk 197 (Djarum Super Kretek 12) punya retur supplier pending yang menahan 1 Pack. 2. Buka Edit Produk, set field "Stok" ke 0 (di bawah stok yang sebenarnya tersedia setelah dikurangi reservasi). 3. Klik Simpan → Ya, Simpan.
- **Hasil aktual (sebelum fix)**: Toast merah generik "Internal Server Error" (HTTP 500). Root cause: `ComputeStockDelta` mengembalikan `ErrInsufficientStock` (sentinel error Go biasa), tapi `product_repo.go` `Update()` tidak pernah membungkusnya jadi `*errors.BadRequestError` sebelum diteruskan ke handler — middleware error global menganggapnya error tak dikenal → 500.
- **Hasil yang diharapkan**: Pesan error jelas dalam Bahasa Indonesia menjelaskan stok tidak cukup, HTTP 400.
- **Perbaikan**: Ditambahkan `product_repo.WrapStockError()` (fungsi baru, `stock_delta_repo.go`, di-export supaya bisa dipakai lintas domain) yang menerjemahkan `ErrInsufficientStock`/`ErrNeedsStockReview`/`ErrBranchingChain` jadi `*errors.BadRequestError` dengan pesan spesifik. Dipakai di `product_repo.go` `Update()` dan `Create()` (defensif, untuk stok awal produk baru).
- **Status**: ✅ Diperbaiki & diverifikasi ulang — sekarang muncul "Stok produk ID 197 tidak mencukupi untuk perubahan ini (kemungkinan sebagian sedang ditahan retur/reservasi lain)", HTTP 400.

#### [MAYOR] Jual produk `needs_stock_review=true` → 500, bukan pesan jelas (celah #20 seharusnya menolak eksplisit tapi malah crash)
- **Lokasi**: `POST /transactions/create` (checkout Kasir)
- **Langkah reproduksi**: 1. Cari produk dengan `needs_stock_review=1` (mis. id 53 "Kerupuk", hasil backfill Fase 3 yang gagal dianalisis). 2. Coba checkout via API/Kasir untuk produk itu.
- **Hasil aktual (sebelum fix)**: HTTP 500 "Internal Server Error". Root cause SAMA seperti temuan di atas tapi di jalur berbeda: `transaction_repo.go` `Create()`/`Void()` cuma menangani `ErrInsufficientStock` (lewat konvensi prefix string `"stok_insufficient:"`), TIDAK menangani `ErrNeedsStockReview`/`ErrBranchingChain` — keduanya bocor sebagai error mentah. Ditambah lagi, `transaction_service.go` & `purchase_service.go` (`Void`) MEMAKSA semua error non-cocok-prefix jadi `InternalServerError`, bahkan kalau repo sebenarnya sudah mengembalikan `*errors.BadRequestError` yang benar (pola ini juga berpotensi menyembunyikan pesan jelas dari jalur lain di masa depan).
- **Hasil yang diharapkan**: Pesan jelas menolak transaksi karena produk perlu ditinjau dulu, HTTP 400.
- **Perbaikan**: `transaction_repo.go` `Create()`/`Void()` diganti pakai `product_repo.WrapStockError()` (menggantikan konvensi prefix string lama). `transaction_service.go` & `purchase_service.go` diperbaiki supaya pass-through `*errors.BadRequestError` yang sudah dikenal, bukan dipaksa jadi 500.
- **Status**: ✅ Diperbaiki & diverifikasi ulang — sekarang muncul "Produk ID 53 ditandai perlu ditinjau manual (needs_stock_review), operasi stok diblokir sampai ditinjau admin", HTTP 400.
- **Catatan terkait (belum diperbaiki, prioritas rendah)**: Jalur sync/offline (`ApplySyncTransaction`) punya pola serupa tapi TIDAK crash (error cuma dicocokkan via `strings.Contains(err.Error(), "stok produk")` di `sync_service.go` untuk membedakan status "conflict" vs "failed", tidak lewat middleware HTTP) — kalau error-nya `ErrNeedsStockReview` bukan `ErrInsufficientStock`, akan salah masuk kategori "failed" padahal seharusnya "conflict". Tidak diperbaiki sekarang karena berisiko mengubah teks pesan yang jadi andalan pencocokan string itu tanpa waktu cukup untuk verifikasi menyeluruh — dicatat untuk Fase G atau sesi terpisah.

### Temuan minor (dicatat, tidak diperbaiki — kosmetik/tidak berisiko data)

- **[MINOR] Dropdown "Merujuk ke Satuan" saat edit paket tidak mengecualikan descendant-nya sendiri** — saat edit paket "Pack" (yang sudah dirujuk oleh "Pieces"), dropdown "Merujuk ke Satuan" tetap menampilkan "Pieces" sebagai pilihan, padahal memilihnya akan membuat referensi melingkar. **Tidak berisiko data** — BE (`ResolvePackageFactor`) tetap menolak dengan benar saat disimpan ("Referensi tidak valid: konversi paket 333 melingkar", HTTP 400), jadi ini murni UX (opsi yang ditampilkan seharusnya sudah difilter di FE, bukan mengandalkan penolakan BE). Lokasi: `ProductFormModal.tsx`, fungsi `refOptionsFor` (mode edit) — tidak mengecualikan descendant dari paket yang sedang diedit, cuma mengecualikan diri sendiri.
- **[MINOR] Tidak ada validasi kewajaran (sanity check) untuk rasio grosiran ekstrem** — input qty/ref_qty ekstrem (mis. 0.001 : 999999) diterima tanpa peringatan, menghasilkan estimasi "Harga per satuan" yang absurd (Rp 99+ triliun). Tidak ada crash/overflow/NaN/korupsi data (dikonfirmasi lewat query DB langsung, presisi tersimpan persis). Cuma UX yang bisa ditingkatkan (mis. peringatan lembut kalau rasio > threshold tertentu), bukan bug fungsional.

### Semua 24 skenario — ringkasan hasil

| # | Skenario | Hasil |
|---|---|---|
| 1 | Tambah produk, satuan tunggal, stok 0 | ✅ PASS |
| 2 | Tambah produk, stok awal 15 | ✅ PASS |
| 3 | Tambah produk multi-satuan berjenjang | ✅ PASS (breakdown DB benar: Kardus→10 Pack→20 Pieces) |
| 4 | Submit form kosong | ✅ PASS (12 pesan validasi field, tidak submit) |
| 5 | Barcode/nama duplikat | ✅ PASS ("Barcode sudah digunakan") |
| 6 | Margin negatif | ✅ PASS (ternyata BLOCKED validasi keras, bukan cuma peringatan — "Harga jual tidak boleh lebih rendah dari harga beli") |
| 7 | Grosiran rasio ke diri sendiri | ✅ PASS (self tidak muncul di dropdown ref) |
| 8 | Referensi melingkar A→B→A | ✅ PASS (BE tolak "konversi paket X melingkar") — lihat catatan minor FE di atas |
| 9 | Qty/ref_qty ekstrem (0.001/999999) | ✅ PASS (tersimpan presisi, tidak overflow) — lihat catatan minor |
| 10 | Hapus paket yang masih dirujuk | ✅ PASS (ditolak, row tetap ada) |
| 11 | Hapus paket anchor | ✅ PASS (anchor tidak pernah muncul di tabel grosiran) |
| 12 | Edit stok manual admin (naik/turun) | ✅ PASS (audit trail `stock_mutations` lengkap & benar) |
| 13 | Edit stok sbg non-admin (kasir) | ✅ PASS (diblokir di 3 lapis: routing FE, permission BE, field lock) |
| 14 | Turunkan stok di bawah reserved | ✅ PASS setelah fix (lihat bug MAYOR di atas) |
| 15 | needs_stock_review: badge, tandai ditinjau, blokir operasi | ✅ PASS setelah fix (lihat bug MAYOR di atas) |
| 16 | Nonaktifkan produk berstok | ✅ PASS (stok tidak berubah, hilang dari Kasir) |
| 17 | Hapus produk belum ada transaksi | ✅ PASS |
| 18 | Hapus produk sudah ada transaksi | ✅ PASS (ditolak) |
| 19 | Hapus produk masih ada stok | ✅ PASS (ditolak) |
| 20 | Import Excel valid (termasuk grosiran) | ✅ PASS |
| 21 | Import Excel campur baris invalid | ✅ PASS (per-baris error jelas, baris valid tetap bisa diimport terpisah) |
| 22 | Filter kombinasi (search+kategori+status) | ✅ PASS |
| 23 | Sortir tabel by Stok asc/desc | ✅ PASS (urutan benar termasuk nilai pecahan) |
| 24 | Cetak label | ✅ PASS (barcode & harga sesuai data terbaru) |

**Console browser**: 0 JavaScript error/warning ditemukan di seluruh skenario (hanya network 400/500 yang memang diharapkan sebagai bagian pengujian negatif, bukan JS exception).

**Rekomendasi**: lanjut ke Fase B. Dua bug mayor yang ditemukan (500 seharusnya 400) berpotensi memengaruhi Fase B & C juga (sama-sama lewat `ApplyStockDelta`) — sudah diperbaiki secara terpusat (`WrapStockError`) jadi kemungkinan besar jalur lain (retur, write-off) yang sudah punya pembungkus sendiri tidak terpengaruh, tapi tetap perlu diverifikasi ulang di fase masing-masing.
```

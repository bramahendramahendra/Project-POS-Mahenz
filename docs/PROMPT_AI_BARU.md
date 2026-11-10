# Panduan Pindah ke Claude (untuk user) dan Prompt Awal

Ikuti langkah berurutan. Kamu tinggal menyalin-tempel.

## A. Persiapan (dilakukan sekali, sebelum membuka Claude)

1. **Cek apa yang akan di-commit.** Di terminal, folder root proyek:
   ```
   git status
   ```
   Pastikan TIDAK ada file rahasia/besar yang ikut: `.env`, `.env.local`, `.env.prod`, `backups/*.sql`, `*.log`, `testing/`, `node_modules`. File tersebut seharusnya sudah ter-ignore. Bila muncul di daftar "untracked", jangan di-add; minta Claude/saya memeriksa `.gitignore`.
2. **Commit sebagai titik aman** (kamu yang commit). Disarankan per kelompok supaya riwayat rapi:
   ```
   git add CLAUDE.md docs README.md nginx mobile .gitignore
   git commit -m "docs: rancangan, handoff, nginx portal-toko"
   git add PORTAL
   git commit -m "feat(portal): aplikasi Portal (SSO, admin user dan aplikasi)"
   git add TOKO
   git commit -m "feat(toko): pindah ke TOKO, SSO Portal, toggle login langsung"
   ```
   Bila `git status` menunjukkan rename dari `BE/` dan `FE/` ke `TOKO/`, tidak masalah. Periksa lagi `git status` sampai bersih. Jangan `git push` dulu sebelum tahap 3.6 lolos.
3. **Backup DB dev** (opsional tapi aman): salin `pos_retail_db` dan `pos_portal_db` memakai mysqldump atau phpMyAdmin.
4. Pastikan MySQL WAMP hidup bila ingin lanjut uji.

## B. Membuka Claude

Pilih salah satu:
- **Claude Code (disarankan)**: ia membaca `CLAUDE.md` otomatis dan bisa menjalankan perintah, membaca kode, serta mengedit file seperti Kiro. Buka terminal di `d:\Develop\Project_pos_mahenz`, jalankan `claude`.
- **Claude di chat biasa**: tidak punya akses file. Kamu harus melampirkan `CLAUDE.md`, `docs/HANDOFF_AI.md`, dan `docs/RANCANGAN_PORTAL_TOKO_COFFEE.md` secara manual, dan hasilnya jauh kurang praktis. Hindari untuk implementasi.

## C. Prompt pertama (salin seluruh blok ini ke Claude)

```
Kamu melanjutkan proyek POS (Go/Gin/GORM/MySQL + React/Vite) di folder ini dari AI lain (Kiro)
yang kehabisan token. Balas selalu dalam bahasa Indonesia.

LANGKAH 1 (hanya membaca, JANGAN mengubah apa pun):
Baca berurutan: CLAUDE.md, docs/HANDOFF_AI.md, docs/RANCANGAN_PORTAL_TOKO_COFFEE.md
(terutama bagian 2 tabel keputusan, bagian 5.0 dan 5.0b tahapan, bagian 14 rancangan Coffee, bagian 16 riwayat), lalu README.md.

LANGKAH 2: Ringkas pemahamanmu dalam maksimal 20 baris:
- tiga aplikasi, port, dan DB masing-masing
- status tiap tahap (mana selesai, mana menunggu persetujuan, mana berikutnya)
- 8 keputusan terpenting yang tidak boleh dilanggar
- hal yang belum terverifikasi
- pertanyaan bila ada yang ambigu atau bertentangan antara dokumen dan kode

LANGKAH 3: Berhenti dan tunggu konfirmasi saya. Jangan mengubah kode, jangan commit, jangan push, jangan deploy.

Aturan sepanjang pekerjaan: kerja bertahap dan berhenti melapor setiap akhir tahap; perubahan Toko harus aditif;
struktur dan gaya kode Portal/Coffee wajib sama dengan Toko; setiap keputusan baru dicatat di
docs/RANCANGAN_PORTAL_TOKO_COFFEE.md; verifikasi nyata sebelum menyatakan selesai; bila ragu tanya saya.
```

Tunggu ringkasannya. Koreksi bila ada yang salah, lalu lanjut ke bagian D.

## D. Melanjutkan pekerjaan (urutan prompt berikutnya)

**D1. Tahap 3.6 (uji ujung ke ujung di dev).** Setelah ringkasan Claude benar, kirim:
```
Konfirmasi: ringkasanmu benar. Lanjut tahap 3.6: uji ujung ke ujung di dev (Portal 3000/8080 ke Toko 3001/8081).
Jalankan keempat server (MySQL WAMP sudah hidup). Skenario minimal: login Portal sebagai admin, buat user baru
(harus ganti password), beri akses Toko, buka Toko lewat kartu (SSO), user baru masuk halaman menunggu akses,
beri role di Toko lalu masuk dashboard, logout kembali ke Portal, matikan toggle login langsung lalu cek form login Toko,
kode SSO kedaluwarsa/dipakai ulang ditolak, akun tanpa akses aplikasi tidak bisa launch.
Laporkan hasil tiap skenario, yang gagal dan yang tidak bisa diuji. Jangan commit. Setelah itu berhenti dan tanya saya.
```
Catatan: kunci aplikasi Toko di dev mungkin perlu diisi dulu (lihat HANDOFF bagian 8). Biarkan Claude memandu; ia akan memberi SQL-nya.

**D2. Validasi nginx dan cutover VPS.** Hanya bila 3.6 lolos dan kamu siap deploy:
```
Tinjau nginx/pos-web.conf dan docs/PANDUAN_CUTOVER_PORTAL.md baris demi baris untuk kesalahan.
Jika nginx bisa dipasang/dijalankan lewat Docker atau WSL di mesin ini, validasi dengan nginx -t.
Jelaskan langkah cutover satu per satu dan tunggu saya menjalankannya di VPS. Jangan deploy sendiri.
```

**D3. Fase 4 (Coffee).** Hanya setelah Fase 3 kamu setujui:
```
Fase 3 disetujui. Mulai Fase 4 POS Coffee. Langkah pertama BUKAN kode: baca bagian 14 dokumen rancangan,
susun tahapan Fase 4 (4.1, 4.2, ...) lengkap dengan hasil yang diuji di tiap tahap, tulis ke bagian 5.0c dokumen rancangan,
lalu tampilkan ke saya untuk disetujui. Pertanyaan produk yang masih [TERBUKA] ditanyakan dulu, jangan ditebak.
```

## E. Kebiasaan yang menjaga alur tetap sinkron

- Setiap akhir sesi, minta: `Perbarui docs/HANDOFF_AI.md bagian status tahap dan item terbuka, dan tambahkan riwayat di dokumen rancangan.`
- Awal sesi baru (Claude lupa konteks tiap sesi): `Baca CLAUDE.md dan docs/HANDOFF_AI.md, ringkas status dalam 10 baris, lalu tunggu saya.`
- Jika konteks Claude penuh, gunakan `/compact` atau mulai sesi baru dengan kalimat di atas. Dokumen, bukan chat, yang menyimpan ingatan proyek.
- Commit sendiri setelah tiap tahap disetujui, dengan pesan jelas. Itu titik kembali bila ada yang salah.
- Bila Claude mengusulkan sesuatu yang bertentangan dengan tabel keputusan, minta ia menjelaskan alasannya dan tunggu persetujuanmu sebelum mengubah.
- Jangan menempelkan isi `.env` atau password produksi ke chat.

## F. Pertanyaan yang mungkin muncul
- *Claude bilang butuh izin menjalankan perintah?* Setujui perintah baca/build/test; baca baik-baik sebelum menyetujui perintah yang menghapus data atau `git`.
- *Ada tulisan `.claude/settings.json` berisi izin curl lama?* Itu sisa lama (port 8080 untuk Toko lama). Aman diabaikan atau dihapus.

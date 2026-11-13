# CLAUDE.md - Aturan Kerja Proyek POS Mahenz

File ini dibaca otomatis oleh Claude Code. Baca juga, berurutan:
1. `docs/HANDOFF_AI.md` (status terkini, apa yang sudah dan belum, langkah berikutnya)
2. `docs/RANCANGAN_PORTAL_TOKO_COFFEE.md` (sumber kebenaran: 52 keputusan, rancangan Coffee di bagian 14, riwayat)
3. `README.md`

## Ritme kerja (sama seperti sesi sebelumnya)
Satu tahap sekali jalan, user menyetujui tiap tahap:
- Kerjakan SATU tahap (urutan: 3.6, lalu Fase 4 Coffee sesuai `docs/RANCANGAN_PORTAL_TOKO_COFFEE.md`). Verifikasi nyata, catat ke dokumen rancangan (tabel tahap + riwayat) dan `docs/HANDOFF_AI.md`.
- Setelah tahap selesai: BERHENTI, lapor jujur (apa dilakukan, apa diuji, apa belum diuji), sertakan pesan commit yang disarankan, lalu tanya "Lanjut ke tahap berikutnya?". Jangan mulai tahap berikutnya sebelum user menjawab.
- Fase 4: sebelum kode, susun tahapan 4.1, 4.2, ... di bagian 5.0c dokumen rancangan (dengan hasil uji per tahap), tampilkan, dan tunggu persetujuan user.
- Keputusan produk yang [TERBUKA] atau tidak ada di dokumen: tanyakan ke user (kumpulkan sekaligus), jangan menebak.
- Tindakan di luar mesin dev (VPS, `--apply` DB produksi, uji tablet/APK) atau berisiko (hapus data, ubah history git, ganti rahasia): jelaskan dan tunggu persetujuan.
- Jangan pernah `git commit`/`git push`; user yang commit.
- File `.env*` dan `config_prod.json` sengaja ikut di-commit oleh user (keputusan user). Jangan dipersoalkan atau ditanyakan lagi; tetap jangan menampilkan isinya.

## Aturan wajib
- Balas dan tulis dokumen dalam **bahasa Indonesia**.
- Kerja **bertahap**. Selesai satu tahap: berhenti, lapor jujur (apa yang diuji, apa yang belum), tanya "lanjut?". Jangan lompat tahap.
- Perubahan pada **Toko harus aditif**; jangan merusak production Toko.
- **Jangan `git commit`, `git push`, atau deploy** kecuali user meminta eksplisit. User yang commit.
- Setiap keputusan baru **dicatat** di `docs/RANCANGAN_PORTAL_TOKO_COFFEE.md` (tabel keputusan + riwayat). Jangan menimpa keputusan lama tanpa persetujuan user.
- Bila ragu pada keputusan produk, **tanya user**; jangan menebak.
- Struktur dan gaya kode Portal dan Coffee **wajib sama dengan Toko** (domain/handler/service/repo/dto/model di BE; features/shared di FE). Kode disalin per aplikasi, tidak ada pustaka bersama.
- Jangan baca/tampilkan isi `.env*`, `config_prod.json`, atau rahasia lain di jawaban.
- Operasi berisiko (hapus data, ubah DB produksi, migrasi `--apply`) butuh konfirmasi user dan backup dulu.
- Jangan menyatakan "selesai" tanpa verifikasi nyata (build, lint, type-check, uji API/browser). Nyatakan yang tidak bisa diverifikasi.

## Peta singkat
- `AUTH/BE`: BE User (layanan akun internal, Fase 3B). Hanya `127.0.0.1:8080`, tanpa FE, tidak di nginx. DB `pos_auth_db`. Module Go `auth_api`.
- `PORTAL/{BE,FE}`: login terpusat (SSO), klien BE User. BE 8081, FE 3000. DB `pos_portal_db`. Module Go `portal_api`.
- `TOKO/{BE,FE}`: POS Toko, klien BE User. BE 8082, FE 3001. DB `pos_retail_db`. Module Go `pos_api`.
- `COFFEE/{BE,FE}`: Fase 4 (BE kerangka sementara, 4.1 diulang). BE 8083, FE 3002.
- `nginx/`, `docs/`, `mobile/` (APK Capacitor), `testing/` (di-gitignore, playwright-core).

## Perintah (Windows, PowerShell)
- Gunakan `npm.cmd`, bukan `npm`, bila `npm` terblokir kebijakan eksekusi.
- BE: dari folder `*/BE`: `go build ./...`, `go run main.go` (migrasi DB otomatis saat startup).
- FE: dari folder `*/FE`: `npm.cmd run type-check`, `npm.cmd run lint`, `npm.cmd run build`, `npm.cmd run dev`.
- MySQL WAMP: `C:\wamp64\bin\mysql\mysql8.4.7\bin\mysql.exe -uroot` (tanpa password, dev).

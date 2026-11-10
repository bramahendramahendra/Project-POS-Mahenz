# CLAUDE.md - Aturan Kerja Proyek POS Mahenz

File ini dibaca otomatis oleh Claude Code. Baca juga, berurutan:
1. `docs/HANDOFF_AI.md` (status terkini, apa yang sudah dan belum, langkah berikutnya)
2. `docs/RANCANGAN_PORTAL_TOKO_COFFEE.md` (sumber kebenaran: 52 keputusan, rancangan Coffee di bagian 14, riwayat)
3. `README.md`

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
- `PORTAL/{BE,FE}`: login terpusat (SSO). BE 8080, FE 3000. DB `pos_portal_db`. Module Go `portal_api`.
- `TOKO/{BE,FE}`: POS Toko. BE 8081, FE 3001. DB `pos_retail_db`. Module Go `pos_api`.
- `COFFEE/{BE,FE}`: belum dibuat (Fase 4). BE 8082, FE 3002.
- `nginx/`, `docs/`, `mobile/` (APK Capacitor), `testing/` (di-gitignore, playwright-core).

## Perintah (Windows, PowerShell)
- Gunakan `npm.cmd`, bukan `npm`, bila `npm` terblokir kebijakan eksekusi.
- BE: dari folder `*/BE`: `go build ./...`, `go run main.go` (migrasi DB otomatis saat startup).
- FE: dari folder `*/FE`: `npm.cmd run type-check`, `npm.cmd run lint`, `npm.cmd run build`, `npm.cmd run dev`.
- MySQL WAMP: `C:\wamp64\bin\mysql\mysql8.4.7\bin\mysql.exe -uroot` (tanpa password, dev).

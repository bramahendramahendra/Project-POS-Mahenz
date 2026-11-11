# HANDOFF AI - Portal, Toko, Coffee

Dokumen serah terima dari sesi AI sebelumnya (Kiro) ke AI berikutnya (Claude).
Tanggal: 17 Oktober 2026. Sumber kebenaran keputusan tetap `docs/RANCANGAN_PORTAL_TOKO_COFFEE.md`
(tabel keputusan 1-52 di bagian 2, tahapan di bagian 5.0 dan 5.0b, rancangan Coffee di bagian 14, riwayat #1-#44 di bagian 16).
Dokumen ini hanya merangkum status dan konteks yang tidak ada di sana. Bila bertentangan, dokumen rancangan menang;
bila kode bertentangan dengan dokumen, periksa kode lalu tanyakan ke user.

## 0. Dokumen yang ada di `docs/`
- Inti proyek: `RANCANGAN_PORTAL_TOKO_COFFEE.md` (sumber kebenaran), `HANDOFF_AI.md` (ini), `PROMPT_AI_BARU.md`.
- Deploy aktual: `PANDUAN_CUTOVER_PORTAL.md`, `SOP_DEPLOY.md` (bagian "SOP Deploy per Aplikasi").
- Deploy lama (diberi catatan pemetaan path/service baru, belum ditulis ulang): `PANDUAN_DEPLOY.md`, `DEPLOYMENT_PROD.md`, `DEPLOYMENT_REDEPLOY_FULL.md`, `SETUP_USER_SERVER.md`.
- APK tablet: `PANDUAN_BUILD_APK_TABLET.md`, `RANCANGAN_CETAK_BLUETOOTH_APK_TABLET.md`.
- Dokumen lama yang sudah selesai/usang (Sync Center, migrasi stok, laba rugi, pelaporan, HPP, piutang) sudah dihapus atas keputusan user; riwayatnya ada di git.

## 1. Tujuan besar
Satu repo, tiga aplikasi POS retail (Go/Gin/GORM/MySQL + React/Vite):
- **Portal**: login terpusat (SSO), halaman pilih aplikasi, admin user dan aplikasi.
- **Toko**: aplikasi POS lama (dipindah ke `TOKO/`), kini bisa login lewat Portal.
- **Coffee**: POS kedai kopi baru (BELUM dibuat).

Port: BE Portal 8080, Toko 8081, Coffee 8082. FE Portal 3000, Toko 3001, Coffee 3002.
Service systemd: `pos-portal-backend`, `pos-toko-backend`, `pos-coffee-backend`. DB berbeda per aplikasi.
URL produksi (satu domain, berbasis path): Portal `/`, API Portal `/portal-api/`, Toko `/toko/`, API Toko `/toko/api/`, Coffee kelak `/coffee/` dan `/coffee/api/`.

## 2. Status tahap
| Tahap | Status |
|---|---|
| Fase 1: pindah `BE`/`FE` ke `TOKO/`, port 8081, nginx sumber tunggal | SELESAI |
| Fase 2.1-2.6: Portal BE+FE (auth, ganti password, apps, SSO, admin user/app) | SELESAI, terverifikasi |
| Fase 3.1: SSO di BE Toko (`/api/auth/sso`, migrasi 011) | SELESAI |
| Fase 3.2: toggle login langsung (`direct_login_enabled`) | SELESAI |
| Fase 3.3: FE Toko (base path, `/sso`, `/menunggu-akses`, toggle UI) | SELESAI, e2e browser dev lolos |
| Fase 3.4: skrip `TOKO/BE/cmd/migrate_users_to_portal` | SELESAI, diuji pada salinan DB |
| Fase 3.5: nginx baru, maintenance per aplikasi, panduan cutover, `.env.production` Toko `/toko/` | SELESAI di kode/dokumen, disetujui user |
| Fase 3.6: uji ujung ke ujung di dev | SELESAI (riwayat #46) |
| Fase 3.6b: perbaikan pesan galat SSO Toko (`ServiceUnavailableError` 503) + hapus 4 file sisa | SELESAI, menunggu persetujuan (riwayat #47) |
| Fase 4: POS Coffee | BELUM (rancangan di bagian 14 dokumen rancangan) |

Pekerjaan Fase 1-3.5 sudah di-commit user (commit berikutnya juga oleh user). Belum ada deploy ke VPS.
Standar konsistensi kode Toko/Portal/Coffee: bagian 3.1 dokumen rancangan.

## 3. Hal yang BELUM terverifikasi (jujur)
- `nginx/pos-web.conf` belum divalidasi `nginx -t` (nginx tidak ada di mesin dev).
- `docs/PANDUAN_CUTOVER_PORTAL.md` belum dijalankan di VPS.
- Alur cetak Bluetooth APK di tablet butuh uji manual.
- Build Toko dari kondisi sekarang memakai base `/toko/` sehingga hanya cocok dengan nginx baru.

## 4. Keputusan kunci (ringkas; rincian di dokumen rancangan)
- **SSO**: kode acak 32 byte sekali pakai, TTL 60 detik, hash disimpan. `POST /api/sso/launch` (user login) menghasilkan `redirect_url = <path app>sso?code=...`. `POST /api/sso/exchange` publik dengan header `X-App-Key` dan body `{app_id, code}` (`app_id` = kode teks aplikasi, mis. `"toko"`); penolakan seragam 401. App key `pak_<64 hex>`, Portal hanya menyimpan hash SHA-256, tampil sekali. Kunci aplikasi disimpan di tabel DB aplikasi (Toko: `portal_credentials`), bukan `.env`, diisi lewat SQL manual.
- **Role Portal**: hanya `admin` dan `user`. Migrasi: `admin` Toko menjadi `admin`; semua role lain menjadi `user`. Seed Portal hanya akun `admin`.
- **Admin buat user / reset password**: user wajib ganti password di login pertama (`must_change_password`); selama itu token hanya boleh ke `/auth/me`, `/auth/logout`, `/auth/change-password`.
- **Login Portal**: pembatas 5 gagal per 15 menit, password minimal 8, lupa password lewat admin.
- **Toggle login langsung (Toko)**: hanya role `admin` yang boleh mengubah. Hanya nilai "false" yang mematikan; fail-open. Flag darurat `DirectLoginForceEnabled` di config server mengalahkan DB. Bila Portal belum siap (PortalUrl/credentials kosong), login langsung tetap aktif dan tidak bisa dimatikan. Key dilindungi dari `POST /settings` dan Reset.
- **SSO di Toko**: akun lokal dengan username sama yang belum tertaut DITOLAK (tidak auto-link). User SSO baru mendapat role `menunggu_akses` dan password acak tak dikenal. Role tidak berasal dari Portal; diatur di Toko.
- **Ditolak**: auto-link username, pustaka bersama, endpoint hapus aplikasi (cukup nonaktifkan).
- Hash password Toko adalah Argon2id (package bernama `pkg/bcrypt`, namanya menyesatkan).
- Key localStorage Toko dipindah ke prefix `toko.*` (`bootstrap/migrateStorageKeys.ts`), Portal `portal.*`, Coffee kelak `coffee.*`.
- Mode B (versi lama/baru) dihapus.

## 5. Peta kode penting
- Portal BE: `PORTAL/BE/domain/{auth,user,role,app,sso}`, migrasi `database/migrations/001-004`, config prod: `SecretKey` kosong (dari env `SECRETKEY`).
- Portal FE: routes `/login`, `/apps`, `/account/password`, `/admin/users`, `/admin/apps`; tema indigo/violet (`styles/tokens.css`).
- Toko BE SSO: `pkg/portal/client.go`, `domain/auth/service/auth_service_portal.go`, `domain/auth/repo/auth_repo_portal.go`, migrasi `011_portal_sso.sql`, endpoint `POST /api/auth/sso`, `GET /api/auth/login-options`, `POST /api/auth/direct-login`.
- Toko FE: `src/shared/utils/appUrl.ts`, `src/bootstrap/migrateStorageKeys.ts`, `SsoPage`, `PendingAccessPage`, `features/settings/users/components/DirectLoginCard.tsx`. Env: `.env.development` (Portal `http://localhost:3000/`, API `http://localhost:8081/api`), `.env.production` (`/toko/api`, `/toko/`, `/`).
- Nginx: `nginx/pos-web.conf`, `nginx/maintenance/`. Dokumen deploy: `docs/PANDUAN_CUTOVER_PORTAL.md`, `docs/SOP_DEPLOY.md`.

## 6. Rencana Fase 4 (Coffee), ringkas
Salin `TOKO/BE` dan `TOKO/FE` ke `COFFEE/`, pangkas modul `supplier_return, expiry_batch, customer_balance, capital_reconciliation, backdate, version`; bawa `receivable` dan `customer`. Baseline migrasi baru. Fitur: stok bahan manual + resep catatan, barang jadi (dus/pcs) dipotong saat masuk pesanan, varian/tambahan/meja/antrean harian, open bill, gabung/pindah meja, aturan pesanan lintas hari, tanpa foto/tiket bar, satu kedai, key `coffee.*`, port 8082/3002. Rincian ada di bagian 14 dokumen rancangan. Seperti Portal, kerjakan bertahap dan rancang tahapannya dulu di bagian 5.0 dokumen rancangan sebelum menulis kode, lalu minta persetujuan user.

## 7. Item terbuka kecil
- `set` akses user ke aplikasi dengan id aplikasi yang tidak ada diabaikan diam-diam (pertimbangkan ditolak).
- `TOKO/BE/config/config_prod.json` menyimpan SecretKey/password DB di repo (sarankan pindah ke env).
- Ditunda: opsi "buka aplikasi terakhir otomatis", logout serentak, Husky (cek masih berfungsi setelah pindah folder).
- Di DB dev, Portal masih punya akun `owner` (role admin) dari seed lama; hanya data dev.

## 8. Data dev
User Toko dev: `Banana`, `BramaGanteng` (owner), `admin` (password `admin123`). Sejak tahap 3.6 (11 Nov 2026): skrip migrasi user sudah di-`--apply` pada DB dev (ketiga user Toko tertaut ke Portal), aplikasi `toko` di Portal dev memakai path `http://localhost:3001/`, dan `portal_credentials` Toko dev sudah terisi (kunci hasil rotate-key). Data uji: `kasir72589`, `web12286` (Portal+Toko, password `Baru12345`), `owner72589` (Toko saja, role owner, password `Owner12345`). Untuk produksi path tetap `/toko/` dan kunci dibuat ulang saat cutover.

## 9. Cara menjalankan dev
Jalankan di 4 terminal terpisah (MySQL WAMP harus hidup):
- `PORTAL/BE`: `go run main.go` (8080); `PORTAL/FE`: `npm.cmd run dev` (3000)
- `TOKO/BE`: `go run main.go` (8081; migrasi DB otomatis saat startup); `TOKO/FE`: `npm.cmd run dev` (3001)
Bila perintah berbeda, lihat README di masing-masing folder.

## 10. Catatan lingkungan (Windows)
- Pakai `npm.cmd`. Jangan menjalankan server dev lewat perintah yang memblokir; jalankan di proses latar.
- Jangan menunggu dengan sleep di terminal yang sama; baca output/log sesudahnya.
- Editor sering menambah import otomatis (pernah menggandakan import `validator` di Go/TS); bila build gagal, cek import ganda.
- `gofmt -w` pernah memformat ulang banyak file `PORTAL/BE`; hati-hati membuat diff besar tanpa alasan.
- Skrip uji dan playwright ada di `testing/` (di-gitignore). **Semua file uji/sementara (skrip, log server dev,
  screenshot, dump DB sebelum uji) wajib di `testing/<nama-uji>/`, jangan pernah di disk C** (keputusan user).
  Skrip regresi Portal ke Toko: `testing/portal_toko_e2e/` (lihat README di sana).

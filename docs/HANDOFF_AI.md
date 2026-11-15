# HANDOFF AI - BE User, Portal, Toko, Coffee

Ringkasan status dan konteks untuk sesi AI berikutnya. Diperbarui: 15 November 2026 (setelah cutover produksi).
Sumber kebenaran keputusan: `docs/RANCANGAN_PORTAL_TOKO_COFFEE.md` (tabel keputusan bagian 2, tahapan bagian 5.0-5.0e,
rancangan Coffee bagian 14, riwayat bagian 16). Bila bertentangan, dokumen rancangan menang; bila kode bertentangan
dengan dokumen, periksa kode lalu tanyakan ke user.

## 0. Dokumen yang ada di `docs/`
- Inti: `RANCANGAN_PORTAL_TOKO_COFFEE.md` (sumber kebenaran), `HANDOFF_AI.md` (ini), `PENDING.md` (pekerjaan yang ditunda,
  siap copas), `PROMPT_AI_BARU.md`.
- Deploy (berlaku): `PANDUAN_CUTOVER_PORTAL.md` (cutover sekali jalan, sudah dijalankan 15 Nov 2026), `SOP_DEPLOY.md`
  (bagian "SOP Deploy per Aplikasi" untuk deploy rutin), `SETUP_USER_SERVER.md` (membuat user Linux).
- Deploy lama, **USANG** (struktur sebelum cutover; ditandai di kepala file, disimpan sebagai arsip):
  `PANDUAN_DEPLOY.md`, `DEPLOYMENT_PROD.md`, `DEPLOYMENT_REDEPLOY_FULL.md`.
- APK tablet: `PANDUAN_BUILD_APK_TABLET.md`, `RANCANGAN_CETAK_BLUETOOTH_APK_TABLET.md`.

## 1. Gambaran
Satu repo, empat layanan (Go/Gin/GORM/MySQL + React/Vite), **sudah berjalan di produksi sejak 15 Nov 2026**:
- **BE User** (`AUTH/BE`): layanan akun internal (username, password, akses aplikasi, SSO). Tanpa FE, hanya 127.0.0.1.
- **Portal**: login terpusat, halaman pilih aplikasi, kelola user dan aplikasi (klien BE User).
- **Toko**: POS retail (klien BE User).
- **Coffee**: POS kedai kopi (Fase 4, klien BE User).

| Bagian | Port | Service produksi | Alamat produksi | Database |
|---|---|---|---|---|
| BE User | 127.0.0.1:8080 | `pos-auth-backend` | (tidak di nginx) | `pos_auth_db` |
| Portal | BE 8081, FE dev 3000 | `pos-portal-backend` | `/`, API `/portal-api/` | `pos_portal_db` |
| Toko | BE 8082, FE dev 3001 | `pos-toko-backend` | `/toko/`, API `/toko/api/` (+ `/api/` lama sementara) | `pos_retail_db` |
| Coffee | BE 8083, FE dev 3002 | `pos-coffee-backend` | `/coffee/`, API `/coffee/api/` | `pos_coffee_db` |

Server: kode `/opt/pos-mahenz`, FE di `/var/www/pos-web/{portal,apps/toko,apps/coffee}`, satu akun MySQL `pos_user`
(keputusan 122), rahasia di `config/config_prod.json` tiap BE (keputusan 123). Biner BE wajib di-build di server.

## 2. Status tahap
| Tahap | Status |
|---|---|
| Fase 1: pindah `BE`/`FE` ke `TOKO/`, port 8081, nginx sumber tunggal | SELESAI |
| Fase 2.1-2.6: Portal BE+FE (auth, ganti password, apps, SSO, admin user/app) | SELESAI, terverifikasi |
| Fase 3.1: SSO di BE Toko (`/api/auth/sso`, migrasi 011) | SELESAI |
| Fase 3.2: toggle login langsung (`direct_login_enabled`) | SELESAI |
| Fase 3.3: FE Toko (base path, `/sso`, `/menunggu-akses`, toggle UI) | SELESAI, e2e browser dev lolos |
| Fase 3.4: skrip `TOKO/BE/cmd/migrate_users_to_portal` | SELESAI, diuji pada salinan DB. **Digantikan** `cmd/migrate_users_to_auth` di 3B.5 (tujuan BE User) |
| Fase 3.5: nginx baru, maintenance per aplikasi, panduan cutover, `.env.production` Toko `/toko/` | SELESAI di kode/dokumen, disetujui user |
| Fase 3.6: uji ujung ke ujung di dev | SELESAI (riwayat #46) |
| Fase 3.6b: perbaikan pesan galat SSO Toko (`ServiceUnavailableError` 503) + hapus 4 file sisa | SELESAI, menunggu persetujuan (riwayat #47) |
| Fase 4: POS Coffee | Tahapan 4.1-4.11 di bagian 5.0c (disetujui; baris 4.1 diperbarui setelah Fase 3B dan disetujui). 4.1 (diulang) selesai (riwayat #63). 4.2 selesai (riwayat #64). 4.3 selesai (riwayat #65). 4.4 selesai (riwayat #66). 4.5 selesai (riwayat #67): barang jadi dan pembelian dua jenis (keputusan 91-94). 4.6 selesai (riwayat #68). 4.7a selesai (riwayat #69): pesanan BE (keputusan 99-102). 4.7b selesai (riwayat #70). 4.8 selesai (riwayat #71). 4.9 selesai (riwayat #72). 4.10 selesai (riwayat #73): laporan Coffee (keputusan 114-117; rincian penjualan per menu/kombinasi/kategori/metode bayar/jenis pesanan, Laba Rugi dengan biaya bahan, laporan Bahan Baku dan Barang Terbuang, Ringkasan Bisnis ke menu), uji API 15/15, browser 6/6 laptop/tablet/HP. 4.11 selesai (riwayat #74); **Fase 4 SELESAI** (disetujui user): uji e2e penuh satu hari kerja lewat layar 22/22 laptop/tablet/HP, export Excel Laba Rugi dan label Tertunda terverifikasi, regresi API semua tahap hijau; temuan: simpan pesanan terbuka kini wajib kas terbuka (keputusan 118). Fase 4 selesai bila 4.11 disetujui; deploy VPS seluruh aplikasi dilakukan belakangan (keputusan 61, di luar Fase 4). Riwayat lama: (riwayat #55): `COFFEE/BE` sudah ada (`.env` sudah digeser ke 8083 di 3B.6; isi lain diulang di 4.1), menunggu Fase 3B |
| **Fase 3B: BE User (layanan akun terpisah)** | Keputusan 63-74; rancangan 5.0d disetujui. 3B.1-3B.5 selesai (riwayat #56-#61; skrip `TOKO/BE/cmd/migrate_users_to_auth` sudah di-`--apply` di dev). **3B.6 selesai; Fase 3B SELESAI** (riwayat #62, disetujui user): nginx, README, CLAUDE.md, port Coffee 8083; uji ujung ke ujung Portal 40/40 dan Toko 31/32 (1 karena data kas dev) di laptop, tablet, HP. Dokumen deploy menyusul di fase deploy. Berikutnya: Fase 4 tahap 4.1 Coffee diulang (port 8083). Port baru: User 8080, Portal 8081, Toko 8082, Coffee 8083 |
| Temuan keamanan | Password tersimpan asli di `log_requests` Toko/Portal/Coffee (riwayat #56). Keputusan 74: log lama dibiarkan; penyamaran (`helper/logmask`) sudah dipasang di BE User, Portal (3B.3), dan Toko (3B.4); Coffee juga sudah (`COFFEE/BE/helper/logmask`, 4.1 diulang) |
| Fase Deploy (Portal, Toko, Coffee, BE User) | Tahapan D.1-D.5 di bagian 5.0e (keputusan 119): D.1-D.3 AI di dev, D.4 cutover VPS dan D.5 uji tablet/APK oleh user. D.1 selesai (riwayat #76): nginx blok Coffee, build semua BE (termasuk biner Linux) dan FE produksi lolos, `nginx -t` ditunda ke VPS (keputusan 120); user memeriksa nilai produksi BE. D.2 selesai (riwayat #77): `PANDUAN_CUTOVER_PORTAL.md` ditulis ulang untuk empat aplikasi, SOP per aplikasi diperbarui (keputusan 121-124: produksi struktur lama, satu akun `pos_user`, rahasia di `config_prod.json`, Coffee dipasang bersamaan tanpa akses kasir). D.3 selesai, disetujui (riwayat #78): gladi di salinan produksi 14/14 (port 18080-18083, DB `gladi_*`), panduan diperbaiki (build wajib di server, `Log.Path` relatif). **D.4 selesai 15 Nov 2026** (riwayat #79): cutover VPS dijalankan user, semua berjalan; temuan logout diperbaiki (riwayat #80) dan peringatan data kosong Coffee (riwayat #79, keputusan 125-127), sudah di-deploy. **D.5 ditunda** (lihat `PENDING.md`) |

Pekerjaan Fase 1-3.5 sudah di-commit user (commit berikutnya juga oleh user). Deploy VPS (cutover D.4) sudah dilakukan 15 Nov 2026; produksi kini memakai struktur baru (BE User, Portal, Toko, Coffee).
Standar konsistensi kode Toko/Portal/Coffee: bagian 3.1 dokumen rancangan.

## 3. Hal yang BELUM terverifikasi (jujur)
- D.5: APK di tablet sungguhan dan cetak struk Bluetooth lewat Portal/Toko/Coffee (ditunda, lihat `PENDING.md`).
- Kasir yang sedang login sebelum cutover tetap login sesudahnya (tidak diuji khusus; produksi berjalan).
- Coffee: edit nota pembelian lama yang memuat bahan yang kini nonaktif akan ditolak (bahan harus aktif).

## 4. Keputusan kunci (ringkas; rincian di dokumen rancangan)
- **Akun di BE User** (keputusan 63-77): Portal/Toko/Coffee tidak menyimpan password; login diverifikasi BE User
  (`/api/login/verify`), tiap aplikasi menyimpan profil lokal dengan `auth_user_id`, role lokal, dan sesi sendiri.
  Pembatas login 5 gagal per 15 menit. Satu user satu sesi per aplikasi.
- **Kunci aplikasi**: `pak_<64 hex>`, BE User hanya menyimpan hash. Kunci tiap aplikasi di tabel `auth_credentials`
  database aplikasinya (Portal, Toko, Coffee). Dibuat dengan `AUTH/BE: go run ./cmd/set_app_key --app <kode> --confirm`;
  Toko/Coffee juga bisa diputar dari Portal (Kelola Aplikasi). `portal_credentials` sudah dibuang (migrasi Toko 013).
- **SSO**: kartu di Portal -> `/sso/launch` -> `<app>/sso?code=` (sekali pakai, 60 detik) -> BE aplikasi menukar ke BE User.
  Profil lokal dicari lewat `auth_user_id`; bila belum ada, **ditautkan otomatis** ke profil lokal username sama yang
  belum tertaut (keputusan 67), selain itu dibuat profil baru role `menunggu_akses`.
- **Role**: BE User/Portal hanya `admin` dan `user` (pengelola akun). Role kerja (owner, admin, kasir, menunggu_akses)
  diatur di masing-masing aplikasi. Role Admin aplikasi setara Owner (akses penuh).
- **Login langsung** (form username/password di Toko/Coffee): toggle per aplikasi `settings.direct_login_enabled`, diatur
  admin di aplikasi; fitur sudah ada, tidak perlu dicatat sebagai pekerjaan. Bila terkunci, nyalakan lewat SQL.
  Flag `DirectLoginForceEnabled` sudah dibuang (keputusan 73).
- **Kode**: struktur Portal/Coffee sama dengan Toko (domain/handler/service/repo/dto/model; features/shared), disalin per
  aplikasi tanpa pustaka bersama. Hash password Argon2id (package `pkg/bcrypt`, namanya menyesatkan). Key localStorage
  berawalan `portal.*`, `toko.*`, `coffee.*`.
- **Halaman yang butuh data master kosong** (keputusan 125-127): pakai peringatan syarat (`PrerequisiteGuard`), bukan
  halaman error. Endpoint opsi warisan Toko membalas `data: null` saat tabel kosong; di FE selalu pakai `?? []`.

## 5. Peta kode penting
- BE User: `AUTH/BE/domain/{auth,user,role,app,sso}`, alat `cmd/set_app_key`. Migrasi `001-002`.
- Klien BE User di tiap aplikasi: `pkg/authclient`; alur login/SSO/tautan profil di `domain/auth/service/auth_service_account.go`.
- Portal FE: rute `/login`, `/apps`, `/account/password`, `/admin/users` (tombol Akses aplikasi), `/admin/apps`.
- Toko/Coffee FE: `SsoPage`, `PendingAccessPage`, `shared/utils/appUrl.ts`; `.env.production` memakai base `/toko/` atau `/coffee/`.
- Salin akun Toko ke BE User: `TOKO/BE/cmd/migrate_users_to_auth` (sudah dijalankan di produksi saat cutover).
- Nginx: `nginx/pos-web.conf`, `nginx/maintenance/` (maintenance per aplikasi: portal, toko, coffee, all).
- Coffee: kasir `features/sales/cashier`, open bill `domain/transaction/service/open_bill_service.go`, laporan
  `domain/report` (`coffee_report_*`), migrasi `001-010`.

## 6. Pekerjaan berikutnya
- Ditunda (siap copas, lihat `PENDING.md`): D.5 uji tablet/APK, bersih-bersih server setelah stabil, hapus blok nginx `/api/` lama.
- Perlu diskusi dengan user dulu: perbarui library FE (`npm audit`); perbaiki uji Go lama yang gagal.

## 7. Item terbuka kecil
- Uji Go lama (salinan Toko, ada di Toko dan Coffee): `pkg/binder` `TestBindJSON_LargeNumbers` gagal; paket `helper` panik di
  `go test` karena `config` membaca `./.env` saat init; teardown uji integrasi stok Toko gagal menghapus kemasan anak
  (Coffee sudah diperbaiki, riwayat #69).
- Set akses user ke aplikasi dengan id aplikasi yang tidak ada diabaikan diam-diam (pertimbangkan ditolak).
- Endpoint opsi warisan Toko (`/units/options`, `/categories/options`, `/products/options`, dll.) membalas `null` saat kosong;
  Toko tidak terdampak karena datanya terisi.
- Ditunda: opsi "buka aplikasi terakhir otomatis", logout serentak, Husky.
- Baris `docs`, `nginx`, `.claude`, `.github`, `mobile` di `.gitignore` sengaja di-on/off-kan user saat deploy; abaikan.

## 8. Data dev
- Akun ada di `pos_auth_db` (BE User): `admin`/`admin123`, `Banana`, `BramaGanteng` (owner Toko), plus banyak akun uji
  dari skrip (`tokoNNNNN`, `kaskNNNNN`, dll.). Data uji dibiarkan di DB dev (jangan dihapus).
- Kunci uji BE User di `testing/auth_3b2/keys.txt`, terpasang di `auth_credentials` Portal, Toko, Coffee dev.
- Path aplikasi di BE User dev: Toko `http://localhost:3001/`, Coffee `http://localhost:3002/`.
- Data gladi (salinan produksi) sudah dihapus 15 Nov 2026.

## 9. Cara menjalankan dev
Jalankan di terminal terpisah (MySQL WAMP harus hidup):
- `AUTH/BE` (BE User): `go run main.go` (`127.0.0.1:8080`, wajib hidup sebelum Portal dipakai)
- `PORTAL/BE`: `go run main.go` (8081 sejak 3B.3; kunci Portal di tabel `auth_credentials`); `PORTAL/FE`: `npm.cmd run dev` (3000)
- `TOKO/BE`: `go run main.go` (8082 sejak 3B.4; migrasi DB otomatis saat startup; kunci Toko di tabel `auth_credentials`); `TOKO/FE`: `npm.cmd run dev` (3001)
- `COFFEE/BE`: `go run main.go` (8083; kunci di `auth_credentials`); `COFFEE/FE`: `npm.cmd run dev` (3002). Path aplikasi `coffee` di BE User dev: `http://localhost:3002/`
Bila perintah berbeda, lihat README di masing-masing folder.

## 10. Catatan lingkungan (Windows)
- Pakai `npm.cmd`. Jangan menjalankan server dev lewat perintah yang memblokir; jalankan di proses latar.
- Jangan menunggu dengan sleep di terminal yang sama; baca output/log sesudahnya.
- Editor sering menambah import otomatis (pernah menggandakan import `validator` di Go/TS); bila build gagal, cek import ganda.
- `gofmt -w` pernah memformat ulang banyak file `PORTAL/BE`; hati-hati membuat diff besar tanpa alasan.
- Skrip uji dan playwright ada di `testing/` (di-gitignore). **Semua file uji/sementara (skrip, log server dev,
  screenshot, dump DB sebelum uji) wajib di `testing/<nama-uji>/`, jangan pernah di disk C** (keputusan user).
  Skrip regresi Portal ke Toko: `testing/portal_toko_e2e/` (lihat README di sana).

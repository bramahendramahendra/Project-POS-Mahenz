# HANDOFF AI - Portal, Toko, Coffee

Dokumen serah terima dari sesi AI sebelumnya (Kiro) ke AI berikutnya (Claude).
Tanggal: 17 Oktober 2026. Sumber kebenaran keputusan tetap `docs/RANCANGAN_PORTAL_TOKO_COFFEE.md`
(tabel keputusan 1-52 di bagian 2, tahapan di bagian 5.0 dan 5.0b, rancangan Coffee di bagian 14, riwayat #1-#44 di bagian 16).
Dokumen ini hanya merangkum status dan konteks yang tidak ada di sana. Bila bertentangan, dokumen rancangan menang;
bila kode bertentangan dengan dokumen, periksa kode lalu tanyakan ke user.

## 0. Dokumen yang ada di `docs/`
- Inti proyek: `RANCANGAN_PORTAL_TOKO_COFFEE.md` (sumber kebenaran), `HANDOFF_AI.md` (ini), `PROMPT_AI_BARU.md`.
- Deploy aktual: `PANDUAN_CUTOVER_PORTAL.md` (cutover empat aplikasi, ditulis ulang di D.2), `SOP_DEPLOY.md` (bagian "SOP Deploy per Aplikasi").
- Deploy lama (diberi catatan pemetaan path/service baru, belum ditulis ulang): `PANDUAN_DEPLOY.md`, `DEPLOYMENT_PROD.md`, `DEPLOYMENT_REDEPLOY_FULL.md`, `SETUP_USER_SERVER.md`.
- APK tablet: `PANDUAN_BUILD_APK_TABLET.md`, `RANCANGAN_CETAK_BLUETOOTH_APK_TABLET.md`.
- Dokumen lama yang sudah selesai/usang (Sync Center, migrasi stok, laba rugi, pelaporan, HPP, piutang) sudah dihapus atas keputusan user; riwayatnya ada di git.

## 1. Tujuan besar
Satu repo, tiga aplikasi POS retail (Go/Gin/GORM/MySQL + React/Vite):
- **Portal**: login terpusat (SSO), halaman pilih aplikasi, admin user dan aplikasi.
- **Toko**: aplikasi POS lama (dipindah ke `TOKO/`), kini bisa login lewat Portal.
- **Coffee**: POS kedai kopi baru (BELUM dibuat).

Port (keputusan 64): BE User 8080 (hanya 127.0.0.1, tidak di nginx), Portal 8081, Toko 8082, Coffee 8083 (`.env*`, config, proxy Vite, README, `nginx/pos-web.conf` sudah disesuaikan di 3B.6; dokumen deploy menyusul di fase deploy). FE Portal 3000, Toko 3001, Coffee 3002.
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
| Fase 3.4: skrip `TOKO/BE/cmd/migrate_users_to_portal` | SELESAI, diuji pada salinan DB. **Digantikan** `cmd/migrate_users_to_auth` di 3B.5 (tujuan BE User) |
| Fase 3.5: nginx baru, maintenance per aplikasi, panduan cutover, `.env.production` Toko `/toko/` | SELESAI di kode/dokumen, disetujui user |
| Fase 3.6: uji ujung ke ujung di dev | SELESAI (riwayat #46) |
| Fase 3.6b: perbaikan pesan galat SSO Toko (`ServiceUnavailableError` 503) + hapus 4 file sisa | SELESAI, menunggu persetujuan (riwayat #47) |
| Fase 4: POS Coffee | Tahapan 4.1-4.11 di bagian 5.0c (disetujui; baris 4.1 diperbarui setelah Fase 3B dan disetujui). 4.1 (diulang) selesai (riwayat #63). 4.2 selesai (riwayat #64). 4.3 selesai (riwayat #65). 4.4 selesai (riwayat #66). 4.5 selesai (riwayat #67): barang jadi dan pembelian dua jenis (keputusan 91-94). 4.6 selesai (riwayat #68). 4.7a selesai (riwayat #69): pesanan BE (keputusan 99-102). 4.7b selesai (riwayat #70). 4.8 selesai (riwayat #71). 4.9 selesai (riwayat #72). 4.10 selesai (riwayat #73): laporan Coffee (keputusan 114-117; rincian penjualan per menu/kombinasi/kategori/metode bayar/jenis pesanan, Laba Rugi dengan biaya bahan, laporan Bahan Baku dan Barang Terbuang, Ringkasan Bisnis ke menu), uji API 15/15, browser 6/6 laptop/tablet/HP. 4.11 selesai (riwayat #74); **Fase 4 SELESAI** (disetujui user): uji e2e penuh satu hari kerja lewat layar 22/22 laptop/tablet/HP, export Excel Laba Rugi dan label Tertunda terverifikasi, regresi API semua tahap hijau; temuan: simpan pesanan terbuka kini wajib kas terbuka (keputusan 118). Fase 4 selesai bila 4.11 disetujui; deploy VPS seluruh aplikasi dilakukan belakangan (keputusan 61, di luar Fase 4). Riwayat lama: (riwayat #55): `COFFEE/BE` sudah ada (`.env` sudah digeser ke 8083 di 3B.6; isi lain diulang di 4.1), menunggu Fase 3B |
| **Fase 3B: BE User (layanan akun terpisah)** | Keputusan 63-74; rancangan 5.0d disetujui. 3B.1-3B.5 selesai (riwayat #56-#61; skrip `TOKO/BE/cmd/migrate_users_to_auth` sudah di-`--apply` di dev). **3B.6 selesai; Fase 3B SELESAI** (riwayat #62, disetujui user): nginx, README, CLAUDE.md, port Coffee 8083; uji ujung ke ujung Portal 40/40 dan Toko 31/32 (1 karena data kas dev) di laptop, tablet, HP. Dokumen deploy menyusul di fase deploy. Berikutnya: Fase 4 tahap 4.1 Coffee diulang (port 8083). Port baru: User 8080, Portal 8081, Toko 8082, Coffee 8083 |
| Temuan keamanan | Password tersimpan asli di `log_requests` Toko/Portal/Coffee (riwayat #56). Keputusan 74: log lama dibiarkan; penyamaran (`helper/logmask`) sudah dipasang di BE User, Portal (3B.3), dan Toko (3B.4); Coffee juga sudah (`COFFEE/BE/helper/logmask`, 4.1 diulang) |
| Fase Deploy (Portal, Toko, Coffee, BE User) | Tahapan D.1-D.5 di bagian 5.0e (keputusan 119): D.1-D.3 AI di dev, D.4 cutover VPS dan D.5 uji tablet/APK oleh user. D.1 selesai (riwayat #76): nginx blok Coffee, build semua BE (termasuk biner Linux) dan FE produksi lolos, `nginx -t` ditunda ke VPS (keputusan 120); user memeriksa nilai produksi BE. D.2 selesai (riwayat #77): `PANDUAN_CUTOVER_PORTAL.md` ditulis ulang untuk empat aplikasi, SOP per aplikasi diperbarui (keputusan 121-124: produksi struktur lama, satu akun `pos_user`, rahasia di `config_prod.json`, Coffee dipasang bersamaan tanpa akses kasir). D.3 selesai, disetujui (riwayat #78): gladi di salinan produksi 14/14 (port 18080-18083, DB `gladi_*`), panduan diperbaiki (build wajib di server, `Log.Path` relatif). **D.4 selesai 15 Nov 2026** (riwayat #79): cutover VPS dijalankan user, semua berjalan; temuan logout (bagian 7). Perbaikan temuan logout selesai di dev (riwayat #80), menunggu persetujuan dan deploy FE Toko/Coffee; lalu D.5 uji tablet/APK oleh user |

Pekerjaan Fase 1-3.5 sudah di-commit user (commit berikutnya juga oleh user). Deploy VPS (cutover D.4) sudah dilakukan 15 Nov 2026; produksi kini memakai struktur baru (BE User, Portal, Toko, Coffee).
Standar konsistensi kode Toko/Portal/Coffee: bagian 3.1 dokumen rancangan.

## 3. Hal yang BELUM terverifikasi (jujur)
- `nginx/pos-web.conf` dan `PANDUAN_CUTOVER_PORTAL.md` sudah lolos di VPS (D.4, `nginx -t` OK). D.5 (tablet/APK, cetak Bluetooth) belum.
- Alur cetak Bluetooth APK di tablet butuh uji manual.
- Build Toko dari kondisi sekarang memakai base `/toko/` sehingga hanya cocok dengan nginx baru.
- Coffee: edit nota pembelian lama yang memuat bahan yang kini nonaktif akan ditolak (bahan harus aktif); aktifkan bahan dulu bila perlu.

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

## 7. Daftar pekerjaan pascacutover (dicatat 15 Nov 2026, urut prioritas)
1. **[SELESAI, 15 Nov 2026, laporan user]** Deploy FE Toko dan Coffee ke server (perbaikan riwayat #79 peringatan data kosong Coffee, #80 logout ke Portal). Cukup Frontend, BE tidak berubah.
2. D.5: uji tablet/APK (tampil Portal, masuk Toko/Coffee, cetak struk Bluetooth, alamat printer) oleh user.
3. **[SELESAI DI DEV, menunggu commit + deploy FE Coffee]** Syarat Pembelian Coffee: supplier + produk ATAU bahan baku (keputusan 126); plus syarat Daftar Menu (kategori menu) dan pesan kasir tanpa menu/meja (keputusan 127). Type-check, lint, build lolos; uji `testing/cek_bahan/b_syarat.js` 8/8 di laptop, tablet, HP.
4. Setelah stabil: hapus `/opt/pos-mahenz_lama_<TGL>`, `/var/www/pos-web/dist`, service `pos-backend`.
5. Setelah stabil: beri akses Coffee ke kasir lewat Portal, atur role di Coffee.
6. Setelah semua perangkat memakai alamat baru: hapus blok nginx `location ^~ /api/`.
7. Opsional: matikan login langsung Toko/Coffee bila semua lewat Portal.
8. **[SELESAI, 15 Nov 2026]** Data gladi dihapus atas persetujuan user: database `gladi_*`, dump `pos_retail_db_gladi.sql`, folder `testing/deploy_d3/run/` dan log gladi. Skrip gladi dan `counts_before.json` disimpan.
9. ~~`.gitignore` `docs`/`nginx`~~ **[ABAIKAN]**: baris `docs`, `nginx`, `.claude`, `.github`, `mobile` di `.gitignore` sengaja di-on/off-kan user (dikomentari saat deploy produksi, diaktifkan lagi untuk dev lokal). Bukan masalah; jangan dipersoalkan.
10. Nanti: perbarui library FE yang dilaporkan `npm audit` (tahap tersendiri, uji ulang).
11. Nanti: uji Go lama warisan Toko yang gagal (lihat item kecil di bawah).
12. Nanti: rapikan/tandai usang dokumen deploy lama (`PANDUAN_DEPLOY.md`, `DEPLOYMENT_PROD.md`, dll.).

## 7b. Item terbuka kecil
- **Temuan D.4 (diperbaiki di dev, riwayat #80; sudah di-deploy 15 Nov 2026):** logout dari sesi yang sudah tidak berlaku (akun sama login di perangkat lain) dulu dialihkan interceptor 401 ke `/<app>/login`; kini `/auth/logout` dikecualikan dari refresh di `services/api.client.ts` dan 401 saat logout dianggap berhasil di `features/auth/auth.api.ts`. Uji ulang: `testing/logout_portal/` (`proxy.js`, `repro.js`, `repro_sesi.js`).
- Pasca-cutover (beberapa hari setelah stabil): hapus `/opt/pos-mahenz_lama_<TGL>` dan service `pos-backend`, beri akses Coffee ke kasir, hapus blok `location ^~ /api/` lama di nginx.
- Uji Go lama (salinan Toko, ada di Toko dan Coffee): `pkg/binder` `TestBindJSON_LargeNumbers` gagal; paket `helper` panik di `go test` karena `config` membaca `./.env` saat init; teardown uji integrasi stok Toko masih gagal menghapus kemasan anak (Coffee sudah diperbaiki, riwayat #69).
- `set` akses user ke aplikasi dengan id aplikasi yang tidak ada diabaikan diam-diam (pertimbangkan ditolak).
- `TOKO/BE/config/config_prod.json` menyimpan SecretKey/password DB di repo (sarankan pindah ke env).
- Ditunda: opsi "buka aplikasi terakhir otomatis", logout serentak, Husky (cek masih berfungsi setelah pindah folder).
- Di DB dev, Portal masih punya akun `owner` (role admin) dari seed lama; hanya data dev.

## 8. Data dev
User Toko dev: `Banana`, `BramaGanteng` (owner), `admin` (password `admin123`). Sejak tahap 3.6 (11 Nov 2026): skrip migrasi user sudah di-`--apply` pada DB dev (ketiga user Toko tertaut ke Portal), aplikasi `toko` di Portal dev memakai path `http://localhost:3001/`, dan `portal_credentials` Toko dev sudah terisi (kunci hasil rotate-key). Data uji: `kasir72589`, `web12286` (Portal+Toko, password `Baru12345`), `owner72589` (Toko saja, role owner, password `Owner12345`). Untuk produksi path tetap `/toko/` dan kunci dibuat ulang saat cutover.
Sejak Fase 3B (13-14 Nov 2026): akun ada di `pos_auth_db` (BE User; `admin`/`admin123`, user uji `e2eu39038` password `Baru67890`, akun `tokoNNNNN`/`portalNNNNN` dari skrip uji). Kunci uji BE User di `testing/auth_3b2/keys.txt`; terpasang di `auth_credentials` Portal dev dan Toko dev. Aplikasi `toko` di `pos_auth_db` dev memakai path `http://localhost:3001/`. Tautan Toko lama dikosongkan oleh migrasi 012 (keputusan 77) dan terisi lagi saat login; tabel `portal_credentials` Toko tidak dipakai lagi. Kas dev 12 Nov masih terbuka (dibiarkan).

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

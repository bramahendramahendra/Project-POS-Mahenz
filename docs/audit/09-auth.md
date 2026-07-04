# Audit: Modul Auth (login / token / session)

Scope: `FE/src/features/auth/**`, `FE/src/services/{authClient,api.client}.ts`, `BE/domain/auth/**`, `BE/middleware/auth/**`, `BE/pkg/jwt/jwt.go`, `BE/config/**`

## Temuan

### 1. [KRITIS] Refresh token dead-config — expiry sama dengan access token
- **File**: `BE/domain/auth/service/auth_service.go` (baris 32 & 107)
- **Masalah**: `expiresAt` untuk session dihitung dari `config.Cfg.TokenExpire` (8 jam) di kedua tempat (login & refresh), dan `RefreshToken()` handler (baris 95) mengecek expiry yang sama itu.
- **Dampak**: Refresh token kedaluwarsa di detik yang sama dengan access token yang seharusnya ia perpanjang — silent refresh tidak mungkin terjadi. `RefreshTokenExpire: 604800` di `config_dev.json`/`config_prod.json` adalah config mati.

### 2. [KRITIS] Field mismatch FE/BE membuat refresh selalu gagal
- **File FE**: `FE/src/services/api.client.ts` (baris 97-102) — destructure `access_token` dari response `/auth/refresh`
- **File BE**: `BE/domain/auth/dto/dto_auth.go` (baris 41-45, `RefreshResponse`) — field yang dikirim bernama `"token"`, bukan `"access_token"`
- **Dampak**: `access_token` selalu `undefined`, `setSession` dipanggil dengan token undefined, request yang di-retry mengirim `Authorization: Bearer undefined` dan 401 lagi. Dikombinasi dengan #1, refresh **tidak berfungsi sama sekali** — user dipaksa login ulang setiap kali token kedaluwarsa.

### 3. [KRITIS] Placeholder JWT secret ter-commit ke repo
- **File**: `BE/config/config_prod.json` (baris 2) — `"SecretKey": "pos-retail-secret-key-prod-CHANGE-ME"`
- **Dampak**: Jika secret ini pernah ter-deploy apa adanya (atau polanya bisa ditebak), siapa pun bisa memalsukan JWT valid (HS256, `BE/pkg/jwt/jwt.go` baris 65-73) dan menyamar sebagai user/role mana pun. **Wajib rotasi segera & pindahkan ke env var / secret manager, jangan commit ke git.**

### 4. [Tinggi] Tidak ada rate limiting di `/auth/login`
- **File**: `BE/routes/public_routes.go` (baris 21)
- **Dampak**: Tidak ada throttling/lockout/CAPTCHA di seluruh codebase untuk endpoint ini — brute-force password per username tidak terbatas (hanya mengandalkan lambatnya bcrypt).

### 5. [Tinggi] Logout fail-open di client
- **File**: `FE/src/features/auth/auth.api.ts` (baris 74-76) — `useLogoutMutation.onError` tetap memanggil `finishLogout()`
- **Dampak**: Jika call `/auth/logout` gagal (network error, 500), FE tetap membersihkan session lokal dan tampak logout, padahal row session di BE (`DeleteSessionByToken`) tidak pernah terhapus — token yang dicuri tetap valid & bisa di-replay walau UI menunjukkan sudah logout.

### 6. [Medium] Endpoint publik `/auth/verify-token` jadi oracle informasi akun
- **File**: `BE/domain/auth/handler/auth_handler.go` (baris 103-122)
- **Masalah**: Endpoint ini publik by design (untuk cek token), tapi mengembalikan seluruh claims (`user_id`, `username`, `full_name`, `role`) untuk token apa pun yang secara sintaksis valid.
- **Dampak**: Dikombinasi dengan tidak adanya rate limiting (#4), attacker bisa memakai endpoint ini untuk cek apakah token curian/tebakan masih valid dan mengetahui pemilik akunnya, tanpa perlu header Authorization.

### 7. [Medium] Single-session-per-user — perlu konfirmasi ini disengaja
- **File**: `auth_service.go` — `Login` (baris 52) dan `RefreshToken` (baris 127) memanggil `DeleteSessionByUserID` sebelum membuat session baru.
- **Dampak**: Login dari device/tab lain diam-diam menginvalidasi session sebelumnya tanpa peringatan ke user tersebut. Mungkin memang disengaja (POS single-terminal), tapi perlu dikonfirmasi karena juga memperparah dampak bug #1/#2 (retry silent-refresh bisa ter-invalidate oleh login dari tempat lain).

### 8. [Low] Middleware auth ganda, satu tidak terpakai
- **File**: `BE/middleware/auth/bearer_auth_middleware.go` (tidak terpakai) vs `pos_auth_middleware.go` (`POSBearerAuthMiddleware`, yang benar-benar dipasang di `protected_routes.go` baris 17)
- **Dampak**: Dead code yang berisiko membingungkan pengembangan selanjutnya (perubahan pada middleware yang salah, mengira itu yang aktif).

### 9. [Info] Token disimpan di localStorage
- **File**: `FE/src/features/auth/auth.store.ts` (baris 6-34) — via zustand `persist`, default storage (localStorage)
- **Catatan**: Trade-off standar SPA+Bearer-token, bukan bug baru, tapi menaikkan dampak XSS di bagian app manapun jadi full session takeover.

## Yang Sudah Baik
- Tidak ada username enumeration di login — pesan error sama untuk user tidak ada maupun password salah (`auth_service.go` baris 22-29).
- Password hashing (bcrypt) konsisten.
- Penanganan 401 di FE tersentralisasi dengan single-flight refresh queue (`api.client.ts` baris 74-131) — desainnya bagus, hanya "dikebiri" oleh bug #1/#2 sehingga cabang refresh-nya nyaris selalu gagal.
- Role untuk otorisasi diambil dari session DB (`user_role` di context), bukan dari JWT claims — lebih aman karena role tidak bisa dipalsukan lewat token basi/forged, walau berarti klaim `role` di JWT saat ini tidak terpakai.

## Prioritas Perbaikan
1. Rotasi JWT secret produksi, pindahkan ke env var, jangan commit nilai asli ke repo.
2. Perbaiki field mismatch `access_token`/`token` dan logika expiry refresh token — atau redesain refresh sepenuhnya.
3. Tambah rate limiting/lockout di endpoint login.
4. Perbaiki logout agar tidak fail-open (atau pastikan session tetap terhapus lewat retry/best-effort dengan TTL pendek).
5. Batasi/tutup akses publik `verify-token` atau kurangi informasi yang dikembalikan.

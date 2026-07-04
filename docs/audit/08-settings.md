# Audit: Modul Settings (menus / printer / roles / store / users / versions)

Scope: `FE/src/features/settings/**`, `BE/domain/access/**`, `BE/domain/role/**`, `BE/domain/user/**`, `BE/middleware/**`

Modul ini mengatur akses/permission itu sendiri, jadi bug otorisasi di sini punya dampak sistemik.

## Temuan

### 1. [KRITIS] Privilege escalation via `SetRoleAccess`
- **File**: `BE/domain/access/service/service_access.go` (baris 41-63)
- **Masalah**: `SetRoleAccess` hanya cek role target ada, tidak pernah cek apakah `req.RoleID` adalah role milik pemanggil sendiri, atau apakah itu role sistem yang dilindungi.
- **Dampak**: Role apa pun dengan izin `can_edit` di `sistem.roles` bisa memanggil endpoint ini pada **role_id miliknya sendiri** dan memberi diri sendiri izin penuh ke semua menu (termasuk `sistem.users`, `sistem.roles`) — **privilege escalation penuh**. Catatan: fix commit `b2eaa9d` ("bug protected role dinamis") hanya mengganti `RoleMiddleware("owner")` hardcode jadi `PermissionMiddleware` dinamis; tidak menambah guard self-role/system-role di level service.

### 2. [KRITIS] `UserService.Update` bisa dipakai self-escalation
- **File**: `BE/domain/user/service/user_service.go` (baris 66-76); handler `BE/domain/user/handler/user_handler.go` (baris 98-131)
- **Masalah**: Tidak menerima `currentUserID`, tidak pernah cek `id != currentUserID`.
- **Pembanding**: `Delete` (baris 94-97) secara eksplisit memblokir self-delete.
- **Dampak**: User dengan izin `can_edit` pada `sistem.users` bisa mengedit `role_id` miliknya sendiri (`dto_user.go` baris 30-34) untuk promosi diri ke admin. Tidak ada proteksi "last admin" juga.

### 3. [Tinggi] `UserService.ToggleStatus` tanpa self-check di BE (hanya FE)
- **File**: `BE/domain/user/service/user_service.go` (baris 113-122); handler `user_handler.go` (baris 192-214)
- **Masalah**: FE menonaktifkan toggle untuk diri sendiri (`UserTableColumns.tsx` baris 96), tapi handler tidak meneruskan ID pemanggil ke service.
- **Dampak**: Bisa dilewati lewat API langsung — user menonaktifkan akun sendiri, atau (dikombinasi bug #2) menonaktifkan admin terakhir.

### 4. [Tinggi] `RoleService.Update` tidak ada guard `IsSystem`
- **File**: `BE/domain/role/service/role_service.go` (baris 53-62)
- **Pembanding**: `Delete` (baris 72-74) dan `ToggleStatus` (baris 86-88) sama-sama memblokir role sistem.
- **Dampak**: FE hanya menyembunyikan tombol edit untuk `is_system` (`RolesPage.tsx` baris 138), tapi API langsung bisa rename/redescribe role `owner`/`admin`.

### 5. [Tinggi] Tidak ada proteksi "last admin" sama sekali
- **Dampak**: Dikombinasikan dengan #2, #3, #4 — admin terakhir bisa dihapus/didemosi/dinonaktifkan, berpotensi mengunci seluruh sistem dari akses admin tanpa jalur pemulihan.

### 6. [Medium] Validasi nama role FE vs BE bertentangan dua arah
- **File FE**: `roles.schema.ts` (baris 4-7) — regex `^[a-z0-9_]+$` (izinkan underscore, hanya lowercase)
- **File BE**: `dto_role.go` (baris 17) — `alphanum` (izinkan mixed-case, tidak izinkan underscore)
- **Dampak**: Nama dengan underscore lolos FE tapi ditolak BE; nama mixed-case yang diterima BE ditolak FE — dua arah mismatch, submit form membingungkan.

### 7. [Medium] Validasi username: FE punya max-length, BE tidak
- **File**: `users.schema.ts` (baris 4-9) — `min(3).max(50)` + regex alphanumeric; BE `dto_user.go` (baris 20) — hanya `required,min=3,alphanum`, tanpa batas atas.

### 8. [Low] `RoleGuard` FE tidak konsisten dipakai
- **File**: `UserManagementPage.tsx` (baris 19), `UserTableColumns.tsx` (baris 73, 105), `StoreProfilePage.tsx` (baris 62) pakai `<RoleGuard>`; `RolesPage.tsx` (baris 71-163), `MenusPage.tsx` (baris 71-134), `AppVersionTab.tsx` (baris 34), `PrinterSettingsTab.tsx` (baris 236) tidak sama sekali.
- **Dampak**: Tombol selalu tampil di halaman-halaman ini, mengandalkan BE 403 + toast saja — bukan lubang keamanan (BE tetap menegakkan), tapi UX tidak konsisten.

### 9. [Low] `Update...Payload` di-duplikasi manual
- **File**: `roles.types.ts` (baris 11-20), `users.types.ts` (baris 23-33), `menus.schema.ts` (baris 3-18, schema create/edit juga duplikasi field-by-field)

## Yang Sudah Baik
Password hashing (bcrypt) konsisten antara `Create` dan `ChangePassword` (`user_service.go` baris 42, 87). DTO Update user tidak punya field password, jadi tidak ada risiko menimpa hash dengan string kosong saat edit profil — perubahan password sudah benar diisolasi ke flow terpisah (`ChangePasswordModal`).

## Prioritas Perbaikan
1. Tambah guard self-role/system-role di `SetRoleAccess`.
2. Tambah cek `IsSystem` di `RoleService.Update`.
3. Tambah guard `id != currentUserID` di `UserService.Update` (saat `role_id` berubah) dan di `ToggleStatus`.
4. Tambah proteksi "last active admin" di delete/update/toggle user dan toggle-status role.
5. Selaraskan pemakaian `RoleGuard` di semua halaman settings.

# Audit: Modul Finance (cash-drawer / expenses / my-cash / overview)

Scope: `FE/src/features/finance/**`, `BE/domain/expense/**`, `BE/domain/cash_drawer/**`

## Temuan

### 1. [Tinggi] Expense create/update/delete tanpa cek kepemilikan/role
- **File**: `BE/domain/expense/service/expense_service.go` (`Update` baris 96-139, `Delete`), `BE/domain/expense/handler/expense_handler.go` (baris 100-153)
- **Masalah**: Handler & service tidak menerima/cek `userID`/`role` sama sekali.
- **Pembanding**: `BE/domain/cash_drawer/service/cash_drawer_service.go` (baris 191-193, 220-222, 246-248) konsisten menegakkan `if role != owner && role != admin && cd.UserID != requestingUserID { Unauthorized }`.
- **Dampak**: Cashier mana pun bisa mengedit/menghapus catatan expense milik cashier lain, mengubah diam-diam angka rekonsiliasi cash-drawer orang lain.

### 2. [Tinggi] Update saldo cash-drawer setelah expense tidak transaksional, error di-diamkan
- **File**: `BE/domain/expense/service/expense_service.go` — `Create` (baris 56-71), `Update` delta (baris 109-115), `Delete` (baris 133-136)
- **Masalah**: `s.cashDrawerRepo.UpdateExpenses(...)` dipanggil setelah row expense sudah commit, error-nya dibuang (`_ = s.cashDrawerRepo.UpdateExpenses(...)`). Tidak ada transaction yang membungkus expense-write + drawer-update.
- **Dampak**: Jika update saldo drawer gagal, data expense dan `total_expenses` drawer permanen tidak sinkron, tanpa log error.

### 3. [Medium] Expense list/table tidak ada sorting sama sekali
- **File**: `FE/src/features/finance/expenses/expenses.types.ts` (baris 16-22, `ExpenseListFilter`), `ExpenseTableColumns.tsx` (7 kolom tanpa `sortable: true`), `BE/domain/expense/dto/dto_expense.go` (baris 5-12, `GetAllRequest`)
- **Pembanding**: `CashDrawerListFilter` (baris 66-76) & `CashDrawerTableColumns.tsx` (baris 26, 40, 47, 54, 63) punya sort lengkap end-to-end.
- **Dampak**: Gap fitur end-to-end di satu sibling saja.

### 4. [Medium] `shift_id` wajib di FE, tapi tidak ada validasi di BE
- **File**: `BE/domain/cash_drawer/dto/dto_cash_drawer.go` (baris 24, `OpenRequest.ShiftID *int`) — tidak ada tag `validate`
- **Pembanding**: Field lain di struct yang sama (`OpeningBalance`, `Notes`) punya `validate:"min=0"` / `"max=500"`.
- **Dampak**: Panggilan API langsung bisa membuka cash drawer dengan `shift_id: null`.

### 5. [Medium] Deskripsi wajib di FE, opsional di BE (cash-drawer notes & expense description)
- **File**: `FE/src/features/finance/expenses/expenses.schema.ts` (baris 23, wajib), `BE/domain/expense/dto/dto_expense.go` (baris 21, 35 — hanya `max=255`, tanpa `required`)
- **Dampak**: Deskripsi kosong bisa tersimpan lewat API langsung.

### 6. [Low] `UpdateExpensePayload` di-duplikasi manual
- **File**: `FE/src/features/finance/expenses/expenses.types.ts` (baris 24-40)

### 7. [Low] Inkonsistensi label required (*) antara Open/Close modal
- **File**: `OpenCashDrawerModal.tsx` (baris 161-163, ada `*`) vs `CloseCashDrawerModal.tsx` (baris 88, tidak ada `*`) meski validasi Zod sama-sama wajib `min(0)`.

## Yang Sudah Baik
Pola toast sukses/error + `invalidateQueries` konsisten di semua `api.ts` finance. Semua field uang pakai `RupiahInput` secara konsisten, tidak ada campur float/int.

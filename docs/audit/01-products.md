# Audit: Modul Products (categories / products / units)

Scope: `FE/src/features/products/**`, `BE/domain/product/**`, `BE/domain/product_category/**`, `BE/domain/product_unit/**`

## Temuan

### 1. [Medium] Field `description` produk setengah jadi di seluruh stack
- **File**: `BE/domain/product/model/product.go`, `FE/src/features/products/products/products.schema.ts`, `products.types.ts`, `ProductFormModal.tsx`, `ProductDetailModal.tsx`
- **Masalah**: FE schema/types/`ProductDetailModal` menyimpan & menampilkan `description`, tapi `ProductFormModal` tidak punya input untuk mengisinya, dan BE (`CreateRequest`/`UpdateRequest`/model/repo query) sama sekali tidak punya kolom `description`.
- **Dampak**: Deskripsi produk yang "coba" diisi user tidak pernah tersimpan; `ProductDetailModal` selalu menampilkan kosong. Kemungkinan sisa fitur yang dihapus sebagian atau kontrak FE/BE yang tidak pernah direkonsiliasi.
- **Kategori**: consistency / dead feature

### 2. [Low] Validasi nama produk lebih lemah dari sibling
- **File**: `BE/domain/product/dto/dto_product.go` — `UpdateRequest.Name` hanya `required`
- **Pembanding**: Category & Unit mewajibkan `required,min=2,max=100`
- **Dampak**: Nama produk 1 karakter atau >100 karakter bisa disimpan lewat API langsung meski FE membatasi max=100.

### 3. [Low] Kolom `is_active` di tabel Product tidak `sortable`
- **File**: `FE/src/features/products/products/components/ProductTableColumns.tsx`
- **Pembanding**: Category/Unit table columns identik menandai `is_active` sebagai `sortable: true`.
- **Dampak**: UX tidak konsisten, tanpa alasan fungsional.

### 4. [Low] `UpdateUnitPayload` di-duplikasi manual
- **File**: `FE/src/features/products/units/units.types.ts`
- **Pembanding**: Category/Product menurunkan Update type dari Create type.
- **Dampak**: Field baru di `CreateUnitPayload` tidak otomatis ikut ke `UpdateUnitPayload`, risiko drift diam-diam.

### 5. [Low] Mutation bulk-save packages tanpa toast sukses
- **File**: `FE/src/features/products/products/products.api.ts` — `useSaveProductPackagesBulkMutation`
- **Pembanding**: Semua mutation lain di file yang sama menampilkan toast sukses.
- **Dampak**: User tidak dapat konfirmasi visual setelah bulk-save packages berhasil.

## Yang Sudah Baik
FormModal/Table/TableColumns/FilterBar categories & units nyaris identik (sesuai ekspektasi sibling). Kompleksitas tambahan Product (price tier, package, import, generate barcode/SKU) sudah terisolasi rapi di sub-komponen sendiri, tidak bocor ke kode bersama.

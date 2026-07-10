-- =============================================================
-- Migration 003: Route Registry
-- Daftar path FE yang benar-benar terdaftar sebagai route (lihat
-- FE/src/shared/constants/routes.ts dan FE/src/app/router.tsx).
-- Tabel `menus` hanya boleh menunjuk ke path yang ada di sini,
-- supaya field "Path" di Manajemen Menu tidak bisa diisi sembarangan
-- dan merusak navigasi (path tidak nyambung ke route manapun).
--
-- WAJIB bagi developer: setiap kali menambah route baru di FE
-- (router.tsx), tambahkan migration baru yang meng-INSERT baris path
-- itu ke sini SEBELUM path tersebut bisa dipilih di Manajemen Menu.
-- Tabel ini read-only dari sisi aplikasi (tidak ada CRUD di UI),
-- sengaja hanya bisa diisi lewat migration.
-- =============================================================

CREATE TABLE IF NOT EXISTS route_registry (
    id         INT AUTO_INCREMENT PRIMARY KEY,
    path       VARCHAR(150)  UNIQUE NOT NULL,
    label      VARCHAR(100)  NOT NULL,
    is_active  TINYINT(1)    DEFAULT 1,
    created_at DATETIME      DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME      DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) DEFAULT CHARSET=utf8mb4;

INSERT IGNORE INTO route_registry (path, label) VALUES
    ('/dashboard', 'Dashboard'),
    ('/kasir', 'Kasir'),
    ('/transactions', 'Transaksi'),
    ('/products', 'Produk'),
    ('/products/categories', 'Kategori Produk'),
    ('/products/units', 'Unit Produk'),
    ('/suppliers', 'Supplier'),
    ('/suppliers/purchases', 'Pembelian Supplier'),
    ('/suppliers/returns', 'Retur Supplier'),
    ('/customers', 'Pelanggan'),
    ('/receivables', 'Piutang'),
    ('/finance', 'Dashboard Keuangan'),
    ('/finance/cash-drawer', 'Kas Harian'),
    ('/finance/expenses', 'Pengeluaran'),
    ('/finance/my-cash', 'Kas Saya'),
    ('/reports/sales', 'Laporan Penjualan'),
    ('/reports/profit-loss', 'Laporan Laba Rugi'),
    ('/reports/stock', 'Laporan Stok'),
    ('/reports/cashier', 'Laporan Kinerja Kasir'),
    ('/shifts', 'Shift'),
    ('/sync', 'Sync Center'),
    ('/settings/store', 'Profil Toko'),
    ('/settings/users', 'Manajemen User'),
    ('/settings/roles', 'Manajemen Role'),
    ('/settings/menus', 'Manajemen Menu'),
    ('/settings/printer', 'Printer'),
    ('/settings/versions', 'Versi Aplikasi');

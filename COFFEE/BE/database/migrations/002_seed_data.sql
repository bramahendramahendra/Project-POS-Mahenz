-- =============================================================
-- Migration 002: Seed Data — POS Coffee
-- Data awal yang diperlukan agar aplikasi bisa berjalan.
-- Disalin dari seed Toko (002, 009, 011) tanpa modul yang dipangkas:
-- menu Retur, Versi Aplikasi, Historis (backdate), Rekonsiliasi Modal,
-- dan metode bayar "Saldo" (keputusan 26). Akun awal sama seperti Toko
-- (keputusan 62). Gunakan INSERT IGNORE agar aman jika dijalankan ulang.
-- =============================================================

-- -------------------------------------------------------------
-- Payment Statuses
-- -------------------------------------------------------------
INSERT IGNORE INTO payment_statuses (code, label, is_active, sort_order) VALUES
    ('unpaid',  'Hutang',         1, 1),
    ('partial', 'Bayar Sebagian', 1, 2),
    ('paid',    'Lunas',          1, 3);

-- -------------------------------------------------------------
-- Payment Methods (tanpa "Saldo": saldo pelanggan tidak dipakai di Coffee)
-- -------------------------------------------------------------
INSERT IGNORE INTO payment_methods (code, label, is_active, sort_order) VALUES
    ('cash',     'Tunai',         1, 1),
    ('transfer', 'Transfer Bank', 1, 2),
    ('card',     'Kartu',         1, 3),
    ('qris',     'QRIS',          1, 4),
    ('kredit',   'Hutang',        1, 5);

-- -------------------------------------------------------------
-- Roles (3 role sistem seperti Toko + "menunggu akses" untuk profil SSO baru)
-- -------------------------------------------------------------
INSERT IGNORE INTO roles (name, display_name, description, is_system, is_active) VALUES
    ('owner', 'Owner', 'Akses penuh ke seluruh sistem',             1, 1),
    ('admin', 'Admin', 'Akses manajemen tanpa pengaturan sistem',    1, 1),
    ('kasir', 'Kasir', 'Akses terbatas hanya untuk transaksi kasir', 1, 1),
    ('menunggu_akses', 'Menunggu Akses', 'Akun baru dari Portal, belum diberi izin oleh admin', 1, 1);

-- -------------------------------------------------------------
-- Profil lokal awal (keputusan 62, diperbarui setelah Fase 3B): tanpa password lokal.
-- Akun dan password ada di BE User, profil ini ditautkan otomatis saat login pertama
-- akun BE User dengan username sama (keputusan 67).
-- -------------------------------------------------------------
INSERT IGNORE INTO users (username, password, full_name, role_id)
SELECT 'owner', NULL, 'Owner', id
FROM roles WHERE name = 'owner';

INSERT IGNORE INTO users (username, password, full_name, role_id)
SELECT 'admin', NULL, 'Administrator', id
FROM roles WHERE name = 'admin';

-- -------------------------------------------------------------
-- Pengaturan default kedai
-- direct_login_enabled: login langsung (username + password) aktif sampai dimatikan admin.
-- -------------------------------------------------------------
INSERT IGNORE INTO settings (setting_key, setting_value) VALUES
    ('store_name',                 'POS Coffee'),
    ('store_address',              ''),
    ('store_phone',                ''),
    ('store_email',                ''),
    ('tax_enabled',                '0'),
    ('tax_percent',                '11'),
    ('receipt_footer',             'Terima kasih atas kunjungan Anda'),
    ('stock_notification_enabled', '1'),
    ('pagination_sizes',           '[10,20,50]'),
    ('direct_login_enabled',       'true');

-- -------------------------------------------------------------
-- Menus — navigasi aplikasi (bukan menu jualan kopi, menu jualan memakai domain menu_item)
-- Urutan insert: parent dulu, baru children
-- -------------------------------------------------------------

-- Group: Beranda
INSERT IGNORE INTO menus (key_name, label, icon, path, order_index) VALUES
    ('beranda', 'Beranda', 'Home', NULL, 1);

INSERT IGNORE INTO menus (parent_id, key_name, label, icon, path, order_index)
SELECT m.id, 'beranda.dashboard', 'Dashboard', 'LayoutDashboard', '/dashboard', 1
FROM menus m WHERE m.key_name = 'beranda';

-- Group: Penjualan
INSERT IGNORE INTO menus (key_name, label, icon, path, order_index) VALUES
    ('penjualan', 'Penjualan', 'ShoppingCart', NULL, 2);

INSERT IGNORE INTO menus (parent_id, key_name, label, icon, path, order_index)
SELECT m.id, 'penjualan.kasir', 'Kasir', 'ShoppingCart', '/kasir', 1
FROM menus m WHERE m.key_name = 'penjualan';

INSERT IGNORE INTO menus (parent_id, key_name, label, icon, path, order_index)
SELECT m.id, 'penjualan.transaksi', 'Transaksi', 'Receipt', '/transactions', 2
FROM menus m WHERE m.key_name = 'penjualan';

-- Group: Produk
INSERT IGNORE INTO menus (key_name, label, icon, path, order_index) VALUES
    ('produk', 'Produk', 'Package', NULL, 3);

INSERT IGNORE INTO menus (parent_id, key_name, label, icon, path, order_index)
SELECT m.id, 'produk.produk', 'Produk', 'Package', '/products', 1
FROM menus m WHERE m.key_name = 'produk';

INSERT IGNORE INTO menus (parent_id, key_name, label, icon, path, order_index)
SELECT m.id, 'produk.kategori', 'Kategori', 'Tag', '/products/categories', 2
FROM menus m WHERE m.key_name = 'produk';

INSERT IGNORE INTO menus (parent_id, key_name, label, icon, path, order_index)
SELECT m.id, 'produk.unit', 'Unit', 'PackageSearch', '/products/units', 3
FROM menus m WHERE m.key_name = 'produk';

-- Group: Pengadaan
INSERT IGNORE INTO menus (key_name, label, icon, path, order_index) VALUES
    ('pengadaan', 'Pengadaan', 'Truck', NULL, 4);

INSERT IGNORE INTO menus (parent_id, key_name, label, icon, path, order_index)
SELECT m.id, 'pengadaan.supplier', 'Supplier', 'Truck', '/suppliers', 1
FROM menus m WHERE m.key_name = 'pengadaan';

INSERT IGNORE INTO menus (parent_id, key_name, label, icon, path, order_index)
SELECT m.id, 'pengadaan.pembelian', 'Pembelian', 'ShoppingBag', '/suppliers/purchases', 2
FROM menus m WHERE m.key_name = 'pengadaan';

-- Group: Pelanggan
INSERT IGNORE INTO menus (key_name, label, icon, path, order_index) VALUES
    ('pelanggan', 'Pelanggan', 'Users', NULL, 5);

INSERT IGNORE INTO menus (parent_id, key_name, label, icon, path, order_index)
SELECT m.id, 'pelanggan.pelanggan', 'Pelanggan', 'Users', '/customers', 1
FROM menus m WHERE m.key_name = 'pelanggan';

INSERT IGNORE INTO menus (parent_id, key_name, label, icon, path, order_index)
SELECT m.id, 'pelanggan.piutang', 'Piutang', 'CreditCard', '/receivables', 2
FROM menus m WHERE m.key_name = 'pelanggan';

-- Group: Keuangan
INSERT IGNORE INTO menus (key_name, label, icon, path, order_index) VALUES
    ('keuangan', 'Keuangan', 'Wallet', NULL, 6);

INSERT IGNORE INTO menus (parent_id, key_name, label, icon, path, order_index)
SELECT m.id, 'keuangan.dashboard', 'Dashboard Keuangan', 'Wallet', '/finance', 1
FROM menus m WHERE m.key_name = 'keuangan';

INSERT IGNORE INTO menus (parent_id, key_name, label, icon, path, order_index)
SELECT m.id, 'keuangan.kas_harian', 'Kas Harian', 'Landmark', '/finance/cash-drawer', 2
FROM menus m WHERE m.key_name = 'keuangan';

INSERT IGNORE INTO menus (parent_id, key_name, label, icon, path, order_index)
SELECT m.id, 'keuangan.pengeluaran', 'Pengeluaran', 'TrendingDown', '/finance/expenses', 3
FROM menus m WHERE m.key_name = 'keuangan';

INSERT IGNORE INTO menus (parent_id, key_name, label, icon, path, order_index)
SELECT m.id, 'keuangan.kas_saya', 'Kas Saya', 'Wallet', '/finance/my-cash', 4
FROM menus m WHERE m.key_name = 'keuangan';

-- Group: Pelaporan
INSERT IGNORE INTO menus (key_name, label, icon, path, order_index) VALUES
    ('pelaporan', 'Pelaporan', 'BarChart2', NULL, 7);

INSERT IGNORE INTO menus (parent_id, key_name, label, icon, path, order_index)
SELECT m.id, 'pelaporan.penjualan', 'Penjualan', 'TrendingUp', '/reports/sales', 1
FROM menus m WHERE m.key_name = 'pelaporan';

INSERT IGNORE INTO menus (parent_id, key_name, label, icon, path, order_index)
SELECT m.id, 'pelaporan.laba_rugi', 'Laba Rugi', 'LineChart', '/reports/profit-loss', 2
FROM menus m WHERE m.key_name = 'pelaporan';

INSERT IGNORE INTO menus (parent_id, key_name, label, icon, path, order_index)
SELECT m.id, 'pelaporan.stok', 'Stok', 'PackageSearch', '/reports/stock', 3
FROM menus m WHERE m.key_name = 'pelaporan';

INSERT IGNORE INTO menus (parent_id, key_name, label, icon, path, order_index)
SELECT m.id, 'pelaporan.kinerja_kasir', 'Kinerja Kasir', 'BarChart2', '/reports/cashier', 4
FROM menus m WHERE m.key_name = 'pelaporan';

INSERT IGNORE INTO menus (parent_id, key_name, label, icon, path, order_index)
SELECT m.id, 'pelaporan.ringkasan_bisnis', 'Ringkasan Bisnis', 'BarChart3', '/reports/business-summary', 5
FROM menus m WHERE m.key_name = 'pelaporan';

-- Group: Operasional
INSERT IGNORE INTO menus (key_name, label, icon, path, order_index) VALUES
    ('operasional', 'Operasional', 'Clock', NULL, 8);

INSERT IGNORE INTO menus (parent_id, key_name, label, icon, path, order_index)
SELECT m.id, 'operasional.shift', 'Shift', 'Clock', '/shifts', 1
FROM menus m WHERE m.key_name = 'operasional';

-- Group: Sistem
INSERT IGNORE INTO menus (key_name, label, icon, path, order_index) VALUES
    ('sistem', 'Sistem', 'Settings', NULL, 9);

INSERT IGNORE INTO menus (parent_id, key_name, label, icon, path, order_index)
SELECT m.id, 'sistem.profil_toko', 'Profil Kedai', 'Settings', '/settings/store', 1
FROM menus m WHERE m.key_name = 'sistem';

INSERT IGNORE INTO menus (parent_id, key_name, label, icon, path, order_index)
SELECT m.id, 'sistem.users', 'Manajemen User', 'Users', '/settings/users', 2
FROM menus m WHERE m.key_name = 'sistem';

INSERT IGNORE INTO menus (parent_id, key_name, label, icon, path, order_index)
SELECT m.id, 'sistem.roles', 'Manajemen Role', 'Shield', '/settings/roles', 3
FROM menus m WHERE m.key_name = 'sistem';

INSERT IGNORE INTO menus (parent_id, key_name, label, icon, path, order_index)
SELECT m.id, 'sistem.menus', 'Manajemen Menu', 'LayoutList', '/settings/menus', 4
FROM menus m WHERE m.key_name = 'sistem';

INSERT IGNORE INTO menus (parent_id, key_name, label, icon, path, order_index)
SELECT m.id, 'sistem.printer', 'Printer', 'Printer', '/settings/printer', 5
FROM menus m WHERE m.key_name = 'sistem';

INSERT IGNORE INTO menus (parent_id, key_name, label, icon, path, order_index)
SELECT m.id, 'sistem.backup', 'Backup & Restore', 'DatabaseBackup', '/settings/backup', 6
FROM menus m WHERE m.key_name = 'sistem';

-- -------------------------------------------------------------
-- Route Registry — daftar path FE yang valid, dipakai dropdown "Path"
-- di Manajemen Menu supaya tidak bisa diisi sembarangan.
-- -------------------------------------------------------------
INSERT IGNORE INTO route_registry (path, label) VALUES
    ('/dashboard', 'Dashboard'),
    ('/kasir', 'Kasir'),
    ('/transactions', 'Transaksi'),
    ('/products', 'Produk'),
    ('/products/categories', 'Kategori Produk'),
    ('/products/units', 'Unit Produk'),
    ('/suppliers', 'Supplier'),
    ('/suppliers/purchases', 'Pembelian Supplier'),
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
    ('/reports/business-summary', 'Ringkasan Bisnis'),
    ('/reports/stock-reconciliation', 'Rekonsiliasi Stok'),
    ('/shifts', 'Shift'),
    ('/settings/store', 'Profil Kedai'),
    ('/settings/users', 'Manajemen User'),
    ('/settings/roles', 'Manajemen Role'),
    ('/settings/menus', 'Manajemen Menu'),
    ('/settings/printer', 'Printer'),
    ('/settings/backup', 'Backup & Restore');

-- -------------------------------------------------------------
-- Role Menu Access (sama seperti Toko)
-- Owner : akses penuh semua menu operasional, KECUALI menu teknis
--         (sistem.roles, sistem.menus, sistem.backup) — wilayah Admin IT.
-- Admin : admin IT — akses penuh ke semua menu.
-- Kasir : hanya beranda.dashboard, penjualan.kasir, keuangan.kas_saya.
-- Menunggu akses: tanpa izin menu apa pun (sengaja tidak ada baris).
-- -------------------------------------------------------------

-- OWNER: full access semua menu kecuali menu teknis (wilayah Admin IT)
INSERT IGNORE INTO role_menu_access (role_id, menu_id, can_view, can_create, can_edit, can_delete)
SELECT r.id, m.id, 1, 1, 1, 1
FROM roles r
JOIN menus m ON m.key_name NOT IN ('sistem.roles', 'sistem.menus', 'sistem.backup')
WHERE r.name = 'owner';

-- ADMIN: full access semua menu tanpa kecuali (admin IT)
INSERT IGNORE INTO role_menu_access (role_id, menu_id, can_view, can_create, can_edit, can_delete)
SELECT r.id, m.id, 1, 1, 1, 1
FROM roles r, menus m
WHERE r.name = 'admin';

-- KASIR: dashboard (landing page, view only)
INSERT IGNORE INTO role_menu_access (role_id, menu_id, can_view, can_create, can_edit, can_delete)
SELECT r.id, m.id, 1, 0, 0, 0
FROM roles r
JOIN menus m ON m.key_name = 'beranda.dashboard'
WHERE r.name = 'kasir';

-- KASIR: penjualan.kasir (perlu can_create — transaksi)
INSERT IGNORE INTO role_menu_access (role_id, menu_id, can_view, can_create, can_edit, can_delete)
SELECT r.id, m.id, 1, 1, 0, 0
FROM roles r
JOIN menus m ON m.key_name = 'penjualan.kasir'
WHERE r.name = 'kasir';

-- KASIR: keuangan.kas_saya (buka/tutup kas milik sendiri)
INSERT IGNORE INTO role_menu_access (role_id, menu_id, can_view, can_create, can_edit, can_delete)
SELECT r.id, m.id, 1, 1, 1, 0
FROM roles r
JOIN menus m ON m.key_name = 'keuangan.kas_saya'
WHERE r.name = 'kasir';

-- -------------------------------------------------------------
-- Rekonsiliasi Stok (seperti migrasi 009 Toko: ditambahkan setelah blok di atas,
-- hanya untuk Admin dengan view + edit)
-- -------------------------------------------------------------
INSERT IGNORE INTO menus (parent_id, key_name, label, icon, path, order_index)
SELECT m.id, 'pelaporan.rekonsiliasi_stok', 'Rekonsiliasi Stok', 'ClipboardCheck', '/reports/stock-reconciliation', 6
FROM menus m WHERE m.key_name = 'pelaporan';

INSERT IGNORE INTO role_menu_access (role_id, menu_id, can_view, can_create, can_edit, can_delete)
SELECT r.id, m.id, 1, 0, 1, 0
FROM roles r
JOIN menus m ON m.key_name = 'pelaporan.rekonsiliasi_stok'
WHERE r.name = 'admin';

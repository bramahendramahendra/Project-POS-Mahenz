const ExcelJS = require('exceljs');
const path = require('path');

const HEADER_FILL = { type: 'pattern', pattern: 'solid', fgColor: { argb: 'FF2F5496' } };
const HEADER_FONT = { bold: true, color: { argb: 'FFFFFFFF' } };
const TITLE_FONT = { bold: true, size: 14, color: { argb: 'FF2F5496' } };

function styleHeaderRow(row) {
  row.eachCell((cell) => {
    cell.fill = HEADER_FILL;
    cell.font = HEADER_FONT;
    cell.alignment = { vertical: 'middle', horizontal: 'center', wrapText: true };
    cell.border = {
      top: { style: 'thin' }, left: { style: 'thin' },
      bottom: { style: 'thin' }, right: { style: 'thin' },
    };
  });
  row.height = 28;
}

function addDataBorders(ws, startRow, colCount) {
  for (let r = startRow; r <= ws.rowCount; r++) {
    const row = ws.getRow(r);
    for (let c = 1; c <= colCount; c++) {
      const cell = row.getCell(c);
      cell.border = {
        top: { style: 'thin', color: { argb: 'FFDDDDDD' } },
        left: { style: 'thin', color: { argb: 'FFDDDDDD' } },
        bottom: { style: 'thin', color: { argb: 'FFDDDDDD' } },
        right: { style: 'thin', color: { argb: 'FFDDDDDD' } },
      };
      cell.alignment = { vertical: 'top', wrapText: true };
    }
  }
}

function addTestCaseSheet(wb, name, modul, rows) {
  const ws = wb.addWorksheet(name.slice(0, 31));
  ws.mergeCells('A1:H1');
  ws.getCell('A1').value = `TEST CASE - ${modul}`;
  ws.getCell('A1').font = TITLE_FONT;
  ws.getCell('A1').alignment = { horizontal: 'left' };

  const headers = ['No', 'Role Penguji', 'Skenario Pengujian', 'Langkah Pengujian', 'Data Uji (Input)', 'Expected Result', 'Actual Result', 'Status (Pass/Fail)'];
  const headerRowIdx = 3;
  ws.getRow(headerRowIdx).values = headers;
  styleHeaderRow(ws.getRow(headerRowIdx));

  ws.columns = [
    { width: 5 }, { width: 14 }, { width: 32 }, { width: 42 }, { width: 32 }, { width: 34 }, { width: 20 }, { width: 16 },
  ];

  rows.forEach((r, idx) => {
    ws.addRow([idx + 1, r.role, r.scenario, r.steps, r.data, r.expected, '', '']);
  });

  addDataBorders(ws, headerRowIdx + 1, headers.length);
  ws.views = [{ state: 'frozen', ySplit: headerRowIdx }];
  return ws;
}

function addDummyDataSheet(wb, name, title, headers, rows, widths) {
  const ws = wb.addWorksheet(name.slice(0, 31));
  ws.mergeCells(1, 1, 1, headers.length);
  ws.getCell(1, 1).value = `DUMMY DATA - ${title}`;
  ws.getCell(1, 1).font = TITLE_FONT;

  const headerRowIdx = 3;
  ws.getRow(headerRowIdx).values = headers;
  styleHeaderRow(ws.getRow(headerRowIdx));
  ws.columns = widths.map((w) => ({ width: w }));

  rows.forEach((r) => ws.addRow(r));
  addDataBorders(ws, headerRowIdx + 1, headers.length);
  ws.views = [{ state: 'frozen', ySplit: headerRowIdx }];
  return ws;
}

async function main() {
  const wb = new ExcelJS.Workbook();
  wb.creator = 'QA Documentation Generator';
  wb.created = new Date();

  // ============ SHEET 0: COVER / DAFTAR ISI ============
  const cover = wb.addWorksheet('Daftar Isi');
  cover.mergeCells('A1:D1');
  cover.getCell('A1').value = 'DOKUMEN TESTING APLIKASI POS';
  cover.getCell('A1').font = { bold: true, size: 18, color: { argb: 'FF2F5496' } };
  cover.mergeCells('A2:D2');
  cover.getCell('A2').value = 'Berisi daftar skenario pengujian (test case) per menu dan data dummy untuk kebutuhan testing manual.';
  cover.getCell('A2').font = { italic: true };
  cover.getRow(1).height = 26;

  const menuList = [
    ['1', 'Auth - Login/Logout/Refresh', 'TC - Auth'],
    ['2', 'PIN Kasir', 'TC - PIN'],
    ['3', 'User Management', 'TC - User'],
    ['4', 'Role Management', 'TC - Role'],
    ['5', 'Menu & Access (RBAC)', 'TC - Menu Access'],
    ['6', 'Kategori Produk', 'TC - Kategori'],
    ['7', 'Satuan Produk', 'TC - Satuan'],
    ['8', 'Produk (Master + Import + Paket + Harga)', 'TC - Produk'],
    ['9', 'Supplier', 'TC - Supplier'],
    ['10', 'Pembelian (Purchase)', 'TC - Pembelian'],
    ['11', 'Retur Supplier', 'TC - Retur Supplier'],
    ['12', 'Pelanggan (Customer)', 'TC - Pelanggan'],
    ['13', 'Piutang (Receivable)', 'TC - Piutang'],
    ['14', 'Shift', 'TC - Shift'],
    ['15', 'Kas Harian (Cash Drawer)', 'TC - Kas Harian'],
    ['16', 'Pengeluaran (Expense)', 'TC - Pengeluaran'],
    ['17', 'Transaksi Kasir (POS Sale)', 'TC - Transaksi'],
    ['18', 'Mutasi Stok', 'TC - Mutasi Stok'],
    ['19', 'Laporan (Sales/Laba Rugi/Stok/Kinerja Kasir)', 'TC - Laporan'],
    ['20', 'Dashboard', 'TC - Dashboard'],
    ['21', 'Settings (Toko/Printer/Pajak)', 'TC - Settings'],
    ['22', 'Backup & Restore', 'TC - Backup Restore'],
    ['23', 'Sync Center', 'TC - Sync'],
    ['24', 'App Version', 'TC - Version'],
    ['25', 'Keamanan & Otorisasi Lintas Role', 'TC - Security'],
  ];
  cover.getRow(4).values = ['No', 'Menu / Modul', 'Sheet Test Case', ''];
  styleHeaderRow(cover.getRow(4));
  cover.columns = [{ width: 6 }, { width: 45 }, { width: 22 }, { width: 10 }];
  menuList.forEach((m) => cover.addRow(m));
  addDataBorders(cover, 5, 3);

  cover.addRow([]);
  const dummyHeaderRowIdx = cover.rowCount + 1;
  cover.getRow(dummyHeaderRowIdx).values = ['No', 'Data Dummy', 'Sheet', ''];
  styleHeaderRow(cover.getRow(dummyHeaderRowIdx));
  const dummySheets = [
    ['1', 'Akun Login Siap Pakai', 'Dummy - Akun Login'],
    ['2', 'Data User & Role', 'Dummy - User Role'],
    ['3', 'Data Kategori & Satuan', 'Dummy - Kategori Satuan'],
    ['4', 'Data Produk', 'Dummy - Produk'],
    ['5', 'Data Paket/Varian Satuan Produk', 'Dummy - Produk Paket'],
    ['6', 'Data Harga Grosir Produk', 'Dummy - Produk Grosir'],
    ['7', 'Data Supplier', 'Dummy - Supplier'],
    ['8', 'Data Pembelian', 'Dummy - Pembelian'],
    ['9', 'Data Pelanggan', 'Dummy - Pelanggan'],
    ['10', 'Data Piutang', 'Dummy - Piutang'],
    ['11', 'Data Shift', 'Dummy - Shift'],
    ['12', 'Data Kas Harian', 'Dummy - Kas Harian'],
    ['13', 'Data Pengeluaran', 'Dummy - Pengeluaran'],
    ['14', 'Data Transaksi Kasir', 'Dummy - Transaksi'],
  ];
  dummySheets.forEach((m) => cover.addRow(m));
  addDataBorders(cover, dummyHeaderRowIdx + 1, 3);

  // ============ TEST CASE SHEETS ============

  addTestCaseSheet(wb, 'TC - Auth', 'AUTH (LOGIN/LOGOUT)', [
    { role: 'Owner', scenario: 'Login dengan username & password valid', steps: '1. Buka halaman login\n2. Isi username & password valid\n3. Klik Login', data: 'username: owner\npassword: owner123', expected: 'Login berhasil, redirect ke dashboard, token tersimpan' },
    { role: 'Admin', scenario: 'Login dengan username & password valid', steps: 'Sama seperti di atas', data: 'username: admin\npassword: admin123', expected: 'Login berhasil, redirect ke dashboard' },
    { role: 'Kasir', scenario: 'Login dengan akun kasir', steps: 'Sama seperti di atas', data: 'username: kasir1\npassword: kasir123', expected: 'Login berhasil, menu terbatas sesuai role kasir' },
    { role: 'Semua', scenario: 'Login dengan password salah', steps: '1. Isi username benar, password salah\n2. Klik Login', data: 'username: owner\npassword: salahpassword', expected: 'Muncul error 401 "username atau password salah", tidak redirect' },
    { role: 'Semua', scenario: 'Login dengan username tidak terdaftar', steps: '1. Isi username yang tidak ada\n2. Klik Login', data: 'username: notexist\npassword: apapun123', expected: 'Muncul error 401 Unauthenticated' },
    { role: 'Semua', scenario: 'Login dengan akun nonaktif (is_active=0)', steps: '1. Nonaktifkan user via User Management\n2. Coba login dengan user tersebut', data: 'username: usernonaktif\npassword: password123', expected: 'Login ditolak, pesan akun tidak aktif' },
    { role: 'Semua', scenario: 'Login field kosong', steps: '1. Kosongkan username/password\n2. Klik Login', data: 'username: (kosong)\npassword: (kosong)', expected: 'Validasi required muncul, tidak submit ke server' },
    { role: 'Owner', scenario: 'Login di device kedua menggendorkan sesi device pertama (single session)', steps: '1. Login di Browser A\n2. Login lagi dengan akun sama di Browser B\n3. Kembali ke Browser A, lakukan aksi apapun', data: 'username: owner, password: owner123 (2 device berbeda)', expected: 'Browser A menerima 401 (token invalid) karena sesi lama dihapus' },
    { role: 'Semua', scenario: 'Refresh token', steps: '1. Login\n2. Tunggu/paksa access token dekat expired\n3. Panggil refresh endpoint / reload app', data: 'refresh_token dari hasil login', expected: 'Access token baru diterbitkan, sesi tetap berjalan' },
    { role: 'Semua', scenario: 'Logout', steps: '1. Login\n2. Klik Logout', data: '-', expected: 'Sesi/token dihapus, redirect ke halaman login, token lama tidak bisa dipakai lagi' },
    { role: 'Semua', scenario: 'Akses halaman protected setelah logout dengan token lama', steps: '1. Simpan token sebelum logout\n2. Logout\n3. Panggil API protected dengan token lama', data: 'token lama (expired session)', expected: 'Response 401 Unauthenticated' },
  ]);

  addTestCaseSheet(wb, 'TC - PIN', 'PIN KASIR', [
    { role: 'Kasir/Semua', scenario: 'Set PIN pertama kali', steps: '1. Masuk menu profil/PIN\n2. Set PIN baru', data: 'pin: 1234', expected: 'PIN tersimpan, GET /pin/check mengembalikan sudah ada PIN' },
    { role: 'Kasir/Semua', scenario: 'Set PIN dengan panjang tidak valid (<4 atau >6 digit)', steps: '1. Masukkan PIN 3 digit atau 7 digit', data: 'pin: 123 / pin: 1234567', expected: 'Validasi gagal, error min=4/max=6' },
    { role: 'Kasir/Semua', scenario: 'Set PIN dengan karakter non-numerik', steps: '1. Masukkan PIN mengandung huruf', data: 'pin: 12ab', expected: 'Validasi gagal, error "numeric"' },
    { role: 'Kasir/Semua', scenario: 'Verifikasi PIN benar', steps: '1. Set PIN 1234\n2. Verifikasi dengan PIN 1234', data: 'pin: 1234', expected: 'Verifikasi sukses' },
    { role: 'Kasir/Semua', scenario: 'Verifikasi PIN salah', steps: '1. Verifikasi dengan PIN berbeda dari yang tersimpan', data: 'pin: 9999', expected: 'Verifikasi gagal' },
    { role: 'Kasir/Semua', scenario: 'Ubah PIN (change PIN)', steps: '1. Masukkan PIN lama benar\n2. Masukkan PIN baru', data: 'pin lama: 1234, pin baru: 5678', expected: 'PIN berhasil diubah, verifikasi dengan PIN lama gagal setelah itu' },
  ]);

  addTestCaseSheet(wb, 'TC - User', 'USER MANAGEMENT (Owner Only)', [
    { role: 'Owner', scenario: 'Lihat daftar user', steps: '1. Buka menu User Management', data: '-', expected: 'List user tampil lengkap dengan role & status' },
    { role: 'Admin/Kasir', scenario: 'Akses menu User Management dibatasi', steps: '1. Login sebagai admin/kasir\n2. Akses menu/API /users/list', data: '-', expected: 'Admin/kasir mendapat 403 Forbidden (hanya Owner yang boleh)' },
    { role: 'Owner', scenario: 'Tambah user baru dengan data valid', steps: '1. Klik Tambah User\n2. Isi form\n3. Simpan', data: 'username: kasir2, password: kasir123, full_name: Budi Kasir, role: kasir, is_active: true', expected: 'User baru tersimpan dan muncul di list' },
    { role: 'Owner', scenario: 'Tambah user dengan username sudah dipakai', steps: '1. Isi username yang sudah ada', data: 'username: owner (duplikat)', expected: 'Error validasi unique username' },
    { role: 'Owner', scenario: 'Tambah user dengan password < 6 karakter', steps: '1. Isi password pendek', data: 'password: 123', expected: 'Validasi gagal min=6' },
    { role: 'Owner', scenario: 'Tambah user dengan username < 3 karakter / non-alfanumerik', steps: '1. Isi username tidak valid', data: 'username: a! ', expected: 'Validasi gagal' },
    { role: 'Owner', scenario: 'Edit user (ubah full_name & role) tanpa ganti password', steps: '1. Pilih user\n2. Ubah nama/role\n3. Kosongkan password\n4. Simpan', data: 'full_name: Budi Santoso, role_id: admin', expected: 'Data terupdate, password lama tetap berlaku' },
    { role: 'Owner', scenario: 'Nonaktifkan user (toggle status)', steps: '1. Klik toggle status pada salah satu user', data: 'user: kasir2', expected: 'Status user berubah menjadi nonaktif, user tidak bisa login' },
    { role: 'Owner', scenario: 'Hapus user', steps: '1. Klik hapus pada user test', data: 'user: kasir2', expected: 'User terhapus dari list' },
    { role: 'Owner', scenario: 'Hapus/nonaktifkan diri sendiri sedang login', steps: '1. Coba hapus/nonaktifkan akun owner yang sedang digunakan', data: 'user: owner (diri sendiri)', expected: 'Idealnya ditolak/diberi warning (cek apakah backend mengizinkan - potensi bug)' },
  ]);

  addTestCaseSheet(wb, 'TC - Role', 'ROLE MANAGEMENT', [
    { role: 'Owner/Admin', scenario: 'Lihat daftar role', steps: '1. Buka menu Role Management', data: '-', expected: 'List role (owner, admin, kasir, dst) tampil' },
    { role: 'Kasir', scenario: 'Akses menu Role Management dibatasi', steps: '1. Login sebagai kasir, akses /roles/list', data: '-', expected: '403 Forbidden' },
    { role: 'Owner', scenario: 'Tambah role custom baru', steps: '1. Klik Tambah Role\n2. Isi nama & display name\n3. Simpan', data: 'name: supervisor, display_name: Supervisor, description: Role supervisor toko', expected: 'Role baru tersimpan, muncul di list' },
    { role: 'Admin', scenario: 'Admin mencoba create/update/delete role', steps: '1. Login sebagai admin\n2. Panggil POST /roles/create', data: '-', expected: '403 Forbidden (hanya owner)' },
    { role: 'Owner', scenario: 'Tambah role dengan nama duplikat', steps: '1. Isi nama role yang sudah ada', data: 'name: kasir (duplikat)', expected: 'Validasi gagal unique' },
    { role: 'Owner', scenario: 'Hapus role sistem (owner/admin/kasir)', steps: '1. Coba hapus role "kasir" yang is_system=1', data: 'role: kasir', expected: 'Ditolak, role sistem tidak boleh dihapus' },
    { role: 'Owner', scenario: 'Hapus role custom yang masih dipakai user', steps: '1. Buat role "supervisor"\n2. Assign ke salah satu user\n3. Hapus role tersebut', data: 'role: supervisor', expected: 'Idealnya ditolak/diberi peringatan karena masih ada user terkait (cek behaviour aktual)' },
    { role: 'Owner', scenario: 'Nonaktifkan role (toggle status)', steps: '1. Toggle status role custom', data: 'role: supervisor', expected: 'Status berubah nonaktif' },
  ]);

  addTestCaseSheet(wb, 'TC - Menu Access', 'MENU & ACCESS (RBAC MATRIX)', [
    { role: 'Semua', scenario: 'Ambil menu sesuai role login (my menu)', steps: '1. Login sebagai masing-masing role\n2. Panggil /menus/my', data: 'role: owner / admin / kasir', expected: 'Struktur menu tampil berbeda sesuai hak akses role' },
    { role: 'Owner/Admin', scenario: 'Lihat daftar menu & struktur parent-child', steps: '1. Buka menu Management Menu', data: '-', expected: 'Menu tree tampil dengan grouping benar' },
    { role: 'Owner', scenario: 'Tambah menu baru & reorder', steps: '1. Tambah menu baru\n2. Ubah urutan (reorder)', data: 'key_name: test.menu, label: Test Menu', expected: 'Menu baru tersimpan, urutan berubah sesuai reorder' },
    { role: 'Owner', scenario: 'Atur permission role terhadap menu (set access)', steps: '1. Buka Role Access\n2. Pilih role "supervisor"\n3. Set can_view/create/edit/delete per menu\n4. Simpan', data: 'role: supervisor, menu: produk.produk -> view=1, create=1, edit=0, delete=0', expected: 'Permission tersimpan; user dengan role supervisor hanya bisa lihat & tambah produk, tidak bisa edit/hapus' },
    { role: 'Admin', scenario: 'Admin mencoba set permission (should be owner only)', steps: '1. Login admin\n2. Panggil POST /roles/:id/menus/set', data: '-', expected: '403 Forbidden' },
    { role: 'Supervisor (custom)', scenario: 'Verifikasi permission granular benar-benar diterapkan di setiap endpoint', steps: '1. Login sebagai role custom dengan permission terbatas\n2. Coba create/edit/delete pada modul yang diberi permission terbatas', data: 'lihat skenario di atas', expected: 'Aksi yang tidak diizinkan mengembalikan 403 "Anda tidak memiliki akses ke fitur ini"' },
  ]);

  addTestCaseSheet(wb, 'TC - Kategori', 'KATEGORI PRODUK', [
    { role: 'Owner/Admin', scenario: 'Lihat daftar kategori', steps: '1. Buka menu Kategori Produk', data: '-', expected: 'List kategori tampil' },
    { role: 'Owner/Admin', scenario: 'Tambah kategori baru dengan data valid', steps: '1. Klik Tambah\n2. Isi nama & kode\n3. Simpan', data: 'name: Minuman, code: MIN, description: Kategori minuman', expected: 'Kategori tersimpan dan tampil di list' },
    { role: 'Owner/Admin', scenario: 'Tambah kategori dengan kode duplikat', steps: '1. Isi kode yang sudah ada', data: 'code: MIN (duplikat)', expected: 'Validasi gagal unique code' },
    { role: 'Owner/Admin', scenario: 'Tambah kategori dengan kode > 10 karakter', steps: '1. Isi kode terlalu panjang', data: 'code: KATEGORIPANJANGSEKALI', expected: 'Validasi gagal max length' },
    { role: 'Owner/Admin', scenario: 'Edit kategori', steps: '1. Pilih kategori\n2. Ubah nama\n3. Simpan', data: 'name: Minuman Segar', expected: 'Data terupdate' },
    { role: 'Owner/Admin', scenario: 'Nonaktifkan kategori (toggle status)', steps: '1. Toggle status kategori', data: 'kategori: Minuman', expected: 'Status nonaktif, tidak muncul di dropdown produk aktif' },
    { role: 'Owner/Admin', scenario: 'Hapus kategori yang masih dipakai produk', steps: '1. Hapus kategori yang punya relasi produk', data: 'kategori: Makanan (punya produk)', expected: 'category_id produk terkait menjadi NULL (ON DELETE SET NULL) - verifikasi produk tidak error' },
    { role: 'Kasir', scenario: 'Kasir mengakses list kategori (tanpa permission middleware)', steps: '1. Login kasir\n2. Panggil GET/POST /categories/list', data: '-', expected: 'CATATAN: endpoint list tidak dijaga permission, kemungkinan tetap bisa diakses - verifikasi apakah ini sesuai ekspektasi bisnis' },
  ]);

  addTestCaseSheet(wb, 'TC - Satuan', 'SATUAN PRODUK (UNIT)', [
    { role: 'Owner/Admin', scenario: 'Lihat daftar satuan', steps: '1. Buka menu Satuan Produk', data: '-', expected: 'List satuan tampil (Pcs, Box, Karton, dst)' },
    { role: 'Owner/Admin', scenario: 'Tambah satuan baru', steps: '1. Klik Tambah\n2. Isi nama & abbreviation\n3. Simpan', data: 'name: Karton, abbreviation: KRT', expected: 'Satuan tersimpan' },
    { role: 'Owner/Admin', scenario: 'Tambah satuan tanpa nama (required)', steps: '1. Kosongkan nama\n2. Simpan', data: 'name: (kosong)', expected: 'Validasi gagal required' },
    { role: 'Owner/Admin', scenario: 'Edit satuan', steps: '1. Pilih satuan\n2. Ubah abbreviation\n3. Simpan', data: 'abbreviation: CRT', expected: 'Data terupdate' },
    { role: 'Owner/Admin', scenario: 'Nonaktifkan / hapus satuan yang dipakai produk', steps: '1. Hapus satuan yang direferensikan produk', data: 'unit: Pcs (dipakai produk)', expected: 'unit_id produk terkait menjadi NULL (SET NULL) - verifikasi tampilan produk' },
  ]);

  addTestCaseSheet(wb, 'TC - Produk', 'PRODUK (MASTER, IMPORT, PAKET, HARGA GROSIR)', [
    { role: 'Owner/Admin', scenario: 'Lihat daftar produk (list, search, options)', steps: '1. Buka menu Produk\n2. Coba pencarian nama/barcode', data: 'search: "Indomie"', expected: 'List produk sesuai filter tampil' },
    { role: 'Owner/Admin', scenario: 'Tambah produk baru dengan data lengkap & valid', steps: '1. Klik Tambah Produk\n2. Isi semua field\n3. Simpan', data: 'barcode: 8991002100016, sku: SKU-0001, name: Indomie Goreng, category_id: Makanan, purchase_price: 2500, selling_price: 3000, stock: 100, min_stock: 10, unit_id: Pcs', expected: 'Produk tersimpan dan tampil di list' },
    { role: 'Owner/Admin', scenario: 'Tambah produk dengan barcode/SKU duplikat', steps: '1. Isi barcode yang sudah dipakai produk lain', data: 'barcode: 8991002100016 (duplikat)', expected: 'Validasi gagal unique' },
    { role: 'Owner/Admin', scenario: 'Tambah produk tanpa category_id / unit_id (required)', steps: '1. Kosongkan category_id atau unit_id\n2. Simpan', data: 'category_id: 0, unit_id: 0', expected: 'Validasi gagal required/min>0' },
    { role: 'Owner/Admin', scenario: 'Tambah produk dengan harga/stok negatif', steps: '1. Isi selling_price atau stock negatif', data: 'selling_price: -1000, stock: -5', expected: 'Validasi gagal min=0' },
    { role: 'Owner/Admin', scenario: 'Generate barcode otomatis', steps: '1. Klik tombol Generate Barcode saat tambah produk', data: '-', expected: 'Barcode unik terisi otomatis di field' },
    { role: 'Owner/Admin', scenario: 'Generate SKU otomatis', steps: '1. Klik tombol Generate SKU', data: '-', expected: 'SKU unik terisi otomatis' },
    { role: 'Owner/Admin', scenario: 'Cari produk berdasarkan barcode (kasir scan)', steps: '1. Panggil GET /products/by-barcode/:barcode dengan barcode valid', data: 'barcode: 8991002100016', expected: 'Detail produk ditemukan' },
    { role: 'Owner/Admin', scenario: 'Cari produk dengan barcode tidak terdaftar', steps: '1. Panggil by-barcode dengan barcode asal', data: 'barcode: 0000000000000', expected: 'Response 404 not found' },
    { role: 'Owner/Admin', scenario: 'Tambah paket/varian satuan produk (product package)', steps: '1. Buka detail produk\n2. Tambah paket baru: nama paket, satuan, konversi qty, harga', data: 'package_name: Dus (isi 40), unit_id: Dus, conversion_qty: 40, selling_price: 110000, is_default: false', expected: 'Paket tersimpan dan tampil di detail produk' },
    { role: 'Owner/Admin', scenario: 'Tambah paket dengan conversion_qty <= 0', steps: '1. Isi conversion_qty 0', data: 'conversion_qty: 0', expected: 'Validasi gagal min=0.001' },
    { role: 'Owner/Admin', scenario: 'Set harga grosir/tier price produk', steps: '1. Buka detail produk\n2. Tambah tier harga: nama tier, min_qty, harga', data: 'tier_name: Grosir 12+, min_qty: 12, price: 2800', expected: 'Tier harga tersimpan, dipakai saat transaksi qty >= 12' },
    { role: 'Owner/Admin', scenario: 'Import produk via Excel - preview valid', steps: '1. Download template import\n2. Isi data sesuai template\n3. Upload & preview', data: 'nama: Teh Botol, barcode: 8992388123456, kategori: Minuman, harga_beli: 4000, harga_jual: 5000, stok: 50, stok_minimum: 5, satuan: Pcs', expected: 'Preview menampilkan baris valid, siap diimport' },
    { role: 'Owner/Admin', scenario: 'Import produk via Excel - ada baris invalid (duplikat barcode/kategori tidak ada)', steps: '1. Upload file dengan salah satu baris barcode duplikat atau kategori tidak ditemukan\n2. Lihat hasil preview', data: 'baris 2: barcode sudah ada; baris 3: kategori "XYZ" tidak ada', expected: 'Preview menandai baris tersebut sebagai Error, baris valid tetap bisa diimport (bulk)' },
    { role: 'Owner/Admin', scenario: 'Hapus produk yang sudah pernah bertransaksi', steps: '1. Hapus produk yang sudah punya riwayat transaction_items', data: 'produk: Indomie Goreng (sudah terjual)', expected: 'Verifikasi apakah dihapus permanen atau soft-delete, dan dampaknya ke riwayat transaksi (product_id SET NULL, product_name snapshot tetap tampil)' },
    { role: 'Owner/Admin', scenario: 'Nonaktifkan produk (toggle status)', steps: '1. Toggle status produk', data: 'produk: Teh Botol', expected: 'Produk nonaktif tidak muncul di pencarian kasir' },
    { role: 'Kasir', scenario: 'Kasir mengakses list produk tanpa permission (potensi gap keamanan)', steps: '1. Login kasir\n2. Panggil /products/list secara langsung', data: '-', expected: 'CATATAN: endpoint tidak dijaga permission - cek apakah ini disengaja' },
  ]);

  addTestCaseSheet(wb, 'TC - Supplier', 'SUPPLIER', [
    { role: 'Owner/Admin', scenario: 'Lihat daftar supplier', steps: '1. Buka menu Supplier', data: '-', expected: 'List supplier tampil' },
    { role: 'Owner/Admin', scenario: 'Tambah supplier baru', steps: '1. Klik Tambah\n2. Isi data\n3. Simpan', data: 'supplier_code: SUP-001, name: PT Sumber Makmur, phone: 081234567890, email: sumber@makmur.co.id, contact_person: Andi', expected: 'Supplier tersimpan' },
    { role: 'Owner/Admin', scenario: 'Tambah supplier dengan kode duplikat', steps: '1. Isi kode yang sudah dipakai', data: 'supplier_code: SUP-001 (duplikat)', expected: 'Validasi gagal unique' },
    { role: 'Owner/Admin', scenario: 'Edit data supplier', steps: '1. Pilih supplier\n2. Ubah alamat/telepon\n3. Simpan', data: 'phone: 081298765432', expected: 'Data terupdate' },
    { role: 'Owner/Admin', scenario: 'Nonaktifkan supplier', steps: '1. Toggle status supplier', data: 'supplier: PT Sumber Makmur', expected: 'Status nonaktif, tidak muncul di dropdown pembelian baru' },
    { role: 'Owner/Admin', scenario: 'Hapus supplier yang masih punya riwayat pembelian', steps: '1. Hapus supplier yang punya purchases terkait', data: 'supplier: PT Sumber Makmur (sudah ada transaksi pembelian)', expected: 'supplier_id di purchases menjadi NULL (SET NULL), riwayat pembelian tetap ada' },
  ]);

  addTestCaseSheet(wb, 'TC - Pembelian', 'PEMBELIAN (SUPPLIER PURCHASE)', [
    { role: 'Owner/Admin', scenario: 'Generate kode pembelian otomatis', steps: '1. Klik Tambah Pembelian', data: '-', expected: 'purchase_code otomatis terisi, format unik' },
    { role: 'Owner/Admin', scenario: 'Buat pembelian baru dengan beberapa item', steps: '1. Pilih supplier\n2. Tambah item produk + qty + harga beli\n3. Isi diskon\n4. Simpan', data: 'supplier: PT Sumber Makmur, tanggal: 2026-07-01, item1: Indomie Goreng qty=100 harga=2500, item2: Teh Botol qty=50 harga=4000, discount_amount: 10000', expected: 'Pembelian tersimpan, total_amount terhitung benar, stok produk bertambah, payment_status=unpaid' },
    { role: 'Owner/Admin', scenario: 'Buat pembelian tanpa item (items kosong)', steps: '1. Simpan pembelian tanpa item', data: 'items: []', expected: 'Validasi gagal, minimal 1 item' },
    { role: 'Owner/Admin', scenario: 'Bayar pembelian secara penuh (lunas)', steps: '1. Pilih pembelian unpaid\n2. Bayar sejumlah total_amount', data: 'amount: sama dengan total_amount, payment_method: cash', expected: 'payment_status berubah menjadi paid, remaining_amount = 0' },
    { role: 'Owner/Admin', scenario: 'Bayar pembelian secara sebagian (partial/installment)', steps: '1. Bayar sebagian dari total', data: 'amount: 50% dari total_amount', expected: 'payment_status menjadi partial, remaining_amount berkurang sesuai pembayaran' },
    { role: 'Owner/Admin', scenario: 'Lihat riwayat pembayaran pembelian', steps: '1. Buka detail pembelian\n2. Lihat tab/list payments', data: '-', expected: 'Semua pembayaran (purchase_payments) tampil sesuai urutan tanggal' },
    { role: 'Owner/Admin', scenario: 'Bayar melebihi remaining_amount', steps: '1. Bayar dengan jumlah lebih besar dari sisa hutang', data: 'amount: lebih besar dari remaining_amount', expected: 'Idealnya ditolak/divalidasi - verifikasi behaviour aktual (potensi bug jika tidak divalidasi)' },
    { role: 'Owner/Admin', scenario: 'Edit pembelian yang sudah dibayar sebagian', steps: '1. Ubah item/jumlah pada pembelian dengan status partial', data: '-', expected: 'Verifikasi apakah sistem mengizinkan edit dan bagaimana dampaknya ke stok & remaining_amount' },
    { role: 'Owner/Admin', scenario: 'Hapus pembelian', steps: '1. Hapus data pembelian test', data: '-', expected: 'Pembelian terhapus; verifikasi apakah stok yang sudah ditambahkan ikut dikembalikan/disesuaikan' },
  ]);

  addTestCaseSheet(wb, 'TC - Retur Supplier', 'RETUR SUPPLIER', [
    { role: 'Owner/Admin', scenario: 'Buat retur baru dari pembelian yang sudah ada', steps: '1. Pilih pembelian\n2. Pilih item yang diretur + qty\n3. Isi alasan retur\n4. Simpan', data: 'purchase: PB-0001, item: Teh Botol qty=5, reason: Barang rusak saat pengiriman', expected: 'Retur tersimpan dengan status "pending"' },
    { role: 'Owner/Admin', scenario: 'Buat retur tanpa alasan (reason required)', steps: '1. Kosongkan alasan\n2. Simpan', data: 'reason: (kosong)', expected: 'Validasi gagal required' },
    { role: 'Owner/Admin', scenario: 'Ubah status retur menjadi approved', steps: '1. Pilih retur status pending\n2. Update status ke approved', data: 'status: approved', expected: 'Status berubah approved; verifikasi efek ke stok (apakah stok berkurang lagi/disesuaikan) dan ke saldo hutang supplier' },
    { role: 'Owner/Admin', scenario: 'Ubah status retur menjadi rejected', steps: '1. Update status retur ke rejected', data: 'status: rejected', expected: 'Status berubah rejected, tidak ada efek ke stok/keuangan' },
    { role: 'Owner/Admin', scenario: 'Hapus retur', steps: '1. Hapus data retur test', data: '-', expected: 'Retur terhapus dari list' },
  ]);

  addTestCaseSheet(wb, 'TC - Pelanggan', 'PELANGGAN (CUSTOMER)', [
    { role: 'Owner/Admin', scenario: 'Lihat daftar pelanggan', steps: '1. Buka menu Pelanggan', data: '-', expected: 'List pelanggan tampil' },
    { role: 'Owner/Admin', scenario: 'Tambah pelanggan baru', steps: '1. Klik Tambah\n2. Isi data\n3. Simpan', data: 'customer_code: CUST-001, name: Ibu Sari, phone: 082211223344, address: Jl. Melati No.5, credit_limit: 500000', expected: 'Pelanggan tersimpan' },
    { role: 'Owner/Admin', scenario: 'Tambah pelanggan dengan kode duplikat', steps: '1. Isi kode yang sudah ada', data: 'customer_code: CUST-001 (duplikat)', expected: 'Validasi gagal unique' },
    { role: 'Owner/Admin', scenario: 'Edit data pelanggan & limit kredit', steps: '1. Ubah credit_limit\n2. Simpan', data: 'credit_limit: 1000000', expected: 'Data terupdate' },
    { role: 'Owner/Admin', scenario: 'Nonaktifkan pelanggan', steps: '1. Toggle status pelanggan', data: 'pelanggan: Ibu Sari', expected: 'Status nonaktif, tidak muncul di dropdown transaksi kredit baru' },
    { role: 'Owner/Admin', scenario: 'Hapus pelanggan yang masih memiliki piutang aktif', steps: '1. Hapus pelanggan yang punya receivable belum lunas', data: 'pelanggan: Ibu Sari (masih punya piutang)', expected: 'Verifikasi apakah dihapus dan bagaimana dampaknya ke data piutang (customer_id SET NULL)' },
  ]);

  addTestCaseSheet(wb, 'TC - Piutang', 'PIUTANG (RECEIVABLE)', [
    { role: 'Owner/Admin', scenario: 'Piutang otomatis terbentuk dari transaksi kredit', steps: '1. Buat transaksi kasir dengan is_credit=true & pelanggan\n2. Cek menu Piutang', data: 'transaksi kredit total 150000, customer: Ibu Sari', expected: 'Muncul record piutang baru dengan total_amount=150000, status=unpaid' },
    { role: 'Owner/Admin', scenario: 'Lihat ringkasan piutang (summary)', steps: '1. Buka menu Piutang > Summary', data: '-', expected: 'Total piutang keseluruhan & per pelanggan tampil benar' },
    { role: 'Owner/Admin', scenario: 'Bayar piutang secara penuh', steps: '1. Pilih piutang unpaid\n2. Bayar sesuai total_amount', data: 'amount: sama dengan total_amount', expected: 'Status menjadi paid, remaining_amount=0' },
    { role: 'Owner/Admin', scenario: 'Bayar piutang sebagian (partial)', steps: '1. Bayar sebagian dari total_amount', data: 'amount: 50% dari total_amount', expected: 'Status menjadi partial' },
    { role: 'Owner/Admin', scenario: 'Lihat riwayat pembayaran piutang tertentu', steps: '1. Buka detail piutang\n2. Lihat payments', data: '-', expected: 'Semua pembayaran tampil urut tanggal' },
    { role: 'Owner/Admin', scenario: 'Piutang jatuh tempo (due_date terlewat)', steps: '1. Set due_date piutang di masa lampau\n2. Lihat list/summary piutang', data: 'due_date: 2026-06-01 (lewat)', expected: 'Idealnya ada indikator piutang jatuh tempo/overdue' },
  ]);

  addTestCaseSheet(wb, 'TC - Shift', 'SHIFT', [
    { role: 'Owner/Admin', scenario: 'Lihat daftar shift', steps: '1. Buka menu Shift', data: '-', expected: 'List shift (Pagi, Siang, Malam) tampil' },
    { role: 'Owner/Admin', scenario: 'Tambah shift baru', steps: '1. Klik Tambah\n2. Isi nama & jam mulai/selesai\n3. Simpan', data: 'name: Shift Malam, start_time: 20:00, end_time: 04:00', expected: 'Shift tersimpan' },
    { role: 'Owner/Admin', scenario: 'Tambah shift dengan jam mulai = jam selesai / tidak logis', steps: '1. Isi start_time = end_time', data: 'start_time: 08:00, end_time: 08:00', expected: 'Verifikasi apakah ada validasi jam (potensi gap validasi)' },
    { role: 'Owner/Admin', scenario: 'Edit shift', steps: '1. Ubah jam shift\n2. Simpan', data: 'end_time: 16:00', expected: 'Data terupdate' },
    { role: 'Owner/Admin', scenario: 'Lihat summary shift', steps: '1. Buka Summary Shift', data: '-', expected: 'Ringkasan performa per shift tampil' },
    { role: 'Owner/Admin', scenario: 'Nonaktifkan shift yang sedang dipakai', steps: '1. Toggle status shift yang sedang aktif dipakai kasir', data: 'shift: Shift Pagi', expected: 'Verifikasi dampak ke kasir yang sedang menggunakan shift tersebut' },
  ]);

  addTestCaseSheet(wb, 'TC - Kas Harian', 'KAS HARIAN (CASH DRAWER)', [
    { role: 'Kasir', scenario: 'Buka kas (open drawer) di awal shift', steps: '1. Login kasir\n2. Buka menu Kas Saya\n3. Input saldo awal\n4. Buka Kas', data: 'opening_balance: 500000, shift: Shift Pagi, open_notes: Kas awal shift pagi', expected: 'Kas berhasil dibuka, status=open, muncul di /cash-drawer/my-cash' },
    { role: 'Kasir', scenario: 'Buka kas kedua kali padahal kas masih terbuka', steps: '1. Kas sudah dalam status open\n2. Coba buka kas baru lagi', data: '-', expected: 'Idealnya ditolak (satu kasir hanya boleh 1 kas terbuka) - verifikasi behaviour' },
    { role: 'Kasir', scenario: 'Kas otomatis terupdate saat ada transaksi cash', steps: '1. Buka kas\n2. Lakukan transaksi penjualan pembayaran cash', data: 'transaksi: total 50000, payment_method: cash', expected: 'total_sales & total_cash_sales pada cash_drawer bertambah' },
    { role: 'Kasir', scenario: 'Tutup kas (close drawer) di akhir shift', steps: '1. Input saldo akhir aktual\n2. Tutup kas', data: 'closing_balance: 1200000', expected: 'Kas tertutup, status=closed, expected_balance & difference terhitung otomatis' },
    { role: 'Kasir', scenario: 'Tutup kas dengan selisih (selisih kas)', steps: '1. Input closing_balance berbeda dari expected_balance', data: 'expected_balance: 1200000, closing_balance: 1150000', expected: 'difference tercatat -50000, tampil sebagai selisih kurang' },
    { role: 'Owner/Admin', scenario: 'Lihat daftar & ringkasan semua kas harian', steps: '1. Buka menu Kas Harian (bukan "kas saya")', data: '-', expected: 'List seluruh kas dari semua kasir tampil dengan summary' },
    { role: 'Owner/Admin', scenario: 'Update manual total_sales/total_expenses pada kas', steps: '1. Buka detail kas\n2. Update total sales/expenses secara manual', data: '-', expected: 'Data terupdate, verifikasi hanya owner/admin yang bisa melakukan ini' },
  ]);

  addTestCaseSheet(wb, 'TC - Pengeluaran', 'PENGELUARAN (EXPENSE)', [
    { role: 'Kasir/Owner/Admin', scenario: 'Tambah pengeluaran baru', steps: '1. Buka menu Pengeluaran\n2. Isi kategori, deskripsi, jumlah\n3. Simpan', data: 'expense_date: 2026-07-02, category: Operasional, description: Beli alat tulis kantor, amount: 75000, payment_method: cash', expected: 'Pengeluaran tersimpan, total_expenses di kas harian bertambah (jika cash)' },
    { role: 'Semua', scenario: 'Tambah pengeluaran tanpa amount (required)', steps: '1. Kosongkan amount\n2. Simpan', data: 'amount: (kosong)', expected: 'Validasi gagal required' },
    { role: 'Owner/Admin', scenario: 'Edit pengeluaran', steps: '1. Ubah jumlah/kategori\n2. Simpan', data: 'amount: 100000', expected: 'Data terupdate, kas harian ikut disesuaikan' },
    { role: 'Owner/Admin', scenario: 'Hapus pengeluaran', steps: '1. Hapus data pengeluaran test', data: '-', expected: 'Pengeluaran terhapus, total_expenses kas harian dikembalikan/disesuaikan' },
    { role: 'Owner/Admin', scenario: 'Lihat daftar pengeluaran dengan filter tanggal/kategori', steps: '1. Filter berdasarkan tanggal & kategori', data: 'tanggal: 2026-07-01 s/d 2026-07-02', expected: 'List sesuai filter' },
  ]);

  addTestCaseSheet(wb, 'TC - Transaksi', 'TRANSAKSI KASIR (POS SALE)', [
    { role: 'Kasir', scenario: 'Transaksi penjualan tunai (cash) sederhana', steps: '1. Login kasir, buka kas\n2. Scan/pilih produk\n3. Input qty\n4. Pilih payment_method cash, input uang bayar\n5. Simpan transaksi', data: 'item: Indomie Goreng qty=3 @3000, payment_method: cash, payment_amount: 10000', expected: 'Transaksi tersimpan, change_amount=1000, stok produk berkurang 3, kas harian bertambah 9000' },
    { role: 'Kasir', scenario: 'Transaksi dengan banyak item & diskon per item', steps: '1. Tambah beberapa item berbeda\n2. Beri discount_item pada salah satu item', data: 'item1: Indomie qty=5, item2: Teh Botol qty=2 discount_item=500/item', expected: 'subtotal & total_amount terhitung benar setelah diskon' },
    { role: 'Kasir', scenario: 'Transaksi dengan metode pembayaran non-cash (transfer/qris/card)', steps: '1. Pilih payment_method transfer/qris/card\n2. Simpan', data: 'payment_method: qris', expected: 'Transaksi tersimpan; total_cash_sales pada kas TIDAK bertambah (hanya total_sales)' },
    { role: 'Kasir', scenario: 'Transaksi kredit (is_credit=true) dengan pelanggan', steps: '1. Pilih pelanggan\n2. Set is_credit=true\n3. Simpan', data: 'customer: Ibu Sari, payment_method: kredit, total: 150000', expected: 'Transaksi tersimpan dan otomatis membuat record di Piutang' },
    { role: 'Kasir', scenario: 'Transaksi dengan payment_method tidak valid', steps: '1. Kirim payment_method di luar enum', data: 'payment_method: ewallet', expected: 'Validasi gagal oneof=cash transfer qris card kredit' },
    { role: 'Kasir', scenario: 'Transaksi dengan device_source tidak valid', steps: '1. Kirim device_source di luar enum', data: 'device_source: ios', expected: 'Validasi gagal oneof=desktop web android' },
    { role: 'Kasir', scenario: 'Transaksi tanpa item (items kosong)', steps: '1. Simpan transaksi tanpa item', data: 'items: []', expected: 'Validasi gagal, minimal 1 item' },
    { role: 'Kasir', scenario: 'Transaksi dengan qty item 0 atau negatif', steps: '1. Isi qty 0/negatif pada salah satu item', data: 'qty: 0', expected: 'Validasi gagal min=0.001' },
    { role: 'Kasir', scenario: 'Transaksi dengan stok produk tidak cukup', steps: '1. Pilih produk dengan stok tersisa 2\n2. Input qty beli 10', data: 'produk stok=2, qty beli=10', expected: 'Verifikasi apakah sistem menolak transaksi/menampilkan warning stok tidak cukup' },
    { role: 'Kasir', scenario: 'Transaksi menggunakan harga grosir (tier price) sesuai qty', steps: '1. Beli produk dengan qty >= min_qty tier grosir', data: 'produk tier grosir 12+: price=2800, qty beli=15', expected: 'Harga yang dipakai adalah harga tier grosir, bukan harga normal' },
    { role: 'Kasir', scenario: 'Transaksi menggunakan paket/varian satuan (bukan satuan dasar)', steps: '1. Pilih produk dengan varian paket "Dus isi 40"\n2. Beli 2 dus', data: 'package: Dus (isi 40), qty=2', expected: 'Stok berkurang sebesar qty x conversion_qty (2x40=80), harga sesuai harga dus' },
    { role: 'Owner/Admin', scenario: 'Void transaksi (hanya owner/admin)', steps: '1. Pilih transaksi completed\n2. Void transaksi', data: 'transaction_code: TRX-0001', expected: 'Status berubah menjadi void; verifikasi stok & kas dikembalikan sesuai (stock_mutation type=void)' },
    { role: 'Kasir', scenario: 'Kasir mencoba void transaksi (harus ditolak)', steps: '1. Login kasir\n2. Panggil /transactions/void/:id', data: '-', expected: '403 Forbidden, hanya owner/admin yang boleh void' },
    { role: 'Owner/Admin', scenario: 'Lihat detail & daftar transaksi dengan filter tanggal/kasir', steps: '1. Filter transaksi berdasarkan tanggal, kasir, status', data: 'tanggal: 2026-07-01 s/d 2026-07-02', expected: 'List transaksi sesuai filter tampil dengan benar' },
  ]);

  addTestCaseSheet(wb, 'TC - Mutasi Stok', 'MUTASI STOK (STOCK MUTATION - Owner/Admin)', [
    { role: 'Owner/Admin', scenario: 'Lihat riwayat mutasi stok semua produk', steps: '1. Buka menu Mutasi Stok', data: '-', expected: 'List mutasi (in/out/adjustment/void) tampil lengkap dengan stock_before & stock_after' },
    { role: 'Owner/Admin', scenario: 'Lihat riwayat mutasi stok per produk tertentu', steps: '1. Buka /stock-mutations/product/:product_id', data: 'product: Indomie Goreng', expected: 'Hanya mutasi produk tersebut yang tampil, urut kronologis' },
    { role: 'Owner/Admin', scenario: 'Verifikasi mutasi tercatat otomatis saat pembelian (in)', steps: '1. Buat pembelian baru\n2. Cek mutasi stok produk terkait', data: '-', expected: 'Muncul record mutation_type=in sesuai qty pembelian' },
    { role: 'Owner/Admin', scenario: 'Verifikasi mutasi tercatat otomatis saat transaksi kasir (out)', steps: '1. Buat transaksi penjualan\n2. Cek mutasi stok', data: '-', expected: 'Muncul record mutation_type=out sesuai qty terjual' },
    { role: 'Owner/Admin', scenario: 'Verifikasi mutasi tercatat saat void transaksi', steps: '1. Void sebuah transaksi\n2. Cek mutasi stok', data: '-', expected: 'Muncul record mutation_type=void mengembalikan stok' },
    { role: 'Kasir', scenario: 'Kasir mengakses menu mutasi stok (harus ditolak)', steps: '1. Login kasir\n2. Akses /stock-mutations/list', data: '-', expected: '403 Forbidden (RoleMiddleware owner/admin only)' },
  ]);

  addTestCaseSheet(wb, 'TC - Laporan', 'LAPORAN (SALES / LABA RUGI / STOK / KINERJA KASIR)', [
    { role: 'Owner/Admin', scenario: 'Lihat laporan penjualan (sales) dengan filter tanggal', steps: '1. Buka menu Laporan Penjualan\n2. Filter rentang tanggal', data: 'tanggal: 2026-06-01 s/d 2026-07-02', expected: 'Data penjualan sesuai rentang tampil, total sesuai transaksi completed' },
    { role: 'Owner/Admin', scenario: 'Lihat grafik tren penjualan (sales chart)', steps: '1. Buka chart penjualan', data: '-', expected: 'Grafik tampil sesuai data transaksi' },
    { role: 'Owner/Admin', scenario: 'Export laporan penjualan ke Excel/PDF', steps: '1. Klik Export pada Laporan Penjualan', data: '-', expected: 'File terunduh dengan data sesuai filter' },
    { role: 'Owner/Admin', scenario: 'Lihat laporan laba rugi (profit & loss)', steps: '1. Buka menu Laba Rugi\n2. Filter periode', data: 'periode: Juni 2026', expected: 'Pendapatan, HPP, dan laba/rugi bersih tampil sesuai perhitungan' },
    { role: 'Owner/Admin', scenario: 'Export laporan laba rugi', steps: '1. Klik Export pada Laba Rugi', data: '-', expected: 'File terunduh sesuai data' },
    { role: 'Owner/Admin', scenario: 'Lihat laporan stok (stock report) & ringkasan', steps: '1. Buka menu Laporan Stok', data: '-', expected: 'Daftar stok saat ini, termasuk produk di bawah min_stock ditandai' },
    { role: 'Owner/Admin', scenario: 'Export laporan stok', steps: '1. Klik Export pada Laporan Stok', data: '-', expected: 'File terunduh sesuai data stok saat ini' },
    { role: 'Owner/Admin', scenario: 'Lihat laporan kinerja kasir per periode', steps: '1. Buka menu Kinerja Kasir\n2. Filter tanggal/kasir', data: 'kasir: kasir1, periode: Juli 2026', expected: 'Total transaksi & omzet per kasir tampil benar' },
    { role: 'Owner/Admin', scenario: 'Export laporan kinerja kasir', steps: '1. Klik Export pada Kinerja Kasir', data: '-', expected: 'File terunduh sesuai data' },
    { role: 'Kasir', scenario: 'Kasir mengakses endpoint list/summary laporan (tanpa permission)', steps: '1. Login kasir\n2. Panggil POST /reports/sales/list atau /sales/summary', data: '-', expected: 'CATATAN: endpoint list/summary tidak dijaga permission (hanya export yang dijaga) - verifikasi apakah ini gap keamanan yang perlu diperbaiki' },
  ]);

  addTestCaseSheet(wb, 'TC - Dashboard', 'DASHBOARD', [
    { role: 'Owner/Admin', scenario: 'Lihat statistik ringkas dashboard (stats)', steps: '1. Login owner/admin\n2. Buka Dashboard', data: '-', expected: 'Total penjualan hari ini, jumlah transaksi, dll tampil benar' },
    { role: 'Owner/Admin', scenario: 'Lihat tren penjualan (sales-trend)', steps: '1. Lihat grafik tren di dashboard', data: '-', expected: 'Grafik tren sesuai data transaksi terbaru' },
    { role: 'Owner/Admin', scenario: 'Lihat produk & kategori terlaris (top-products/top-categories)', steps: '1. Lihat widget produk/kategori terlaris', data: '-', expected: 'Data terurut sesuai jumlah/omzet penjualan' },
    { role: 'Owner/Admin', scenario: 'Lihat breakdown metode pembayaran (payment-methods)', steps: '1. Lihat widget metode pembayaran', data: '-', expected: 'Persentase/nilai per metode pembayaran sesuai transaksi' },
    { role: 'Kasir', scenario: 'Kasir memanggil API dashboard langsung (tanpa role guard di BE)', steps: '1. Login kasir\n2. Panggil GET /dashboard/stats langsung via API', data: '-', expected: 'CATATAN: API dashboard tidak memiliki RoleMiddleware di BE (hanya dibatasi di FE) - verifikasi apakah kasir bisa mengakses data ini, ini adalah gap keamanan yang perlu dikonfirmasi ke tim dev' },
  ]);

  addTestCaseSheet(wb, 'TC - Settings', 'SETTINGS (TOKO / PRINTER / PAJAK)', [
    { role: 'Semua', scenario: 'Lihat semua settings', steps: '1. Panggil GET /settings', data: '-', expected: 'Semua key-value settings tampil (store_name, tax_percent, dst)' },
    { role: 'Owner/Admin', scenario: 'Update settings umum', steps: '1. Ubah beberapa setting\n2. Simpan', data: 'tax_enabled: true, tax_percent: 11', expected: 'Setting terupdate' },
    { role: 'Admin', scenario: 'Reset settings ke default (admin only)', steps: '1. Login admin\n2. Klik Reset Settings', data: '-', expected: 'Semua setting kembali ke nilai default awal' },
    { role: 'Owner', scenario: 'Owner mencoba reset settings', steps: '1. Login owner\n2. Panggil /settings/reset', data: '-', expected: 'Verifikasi apakah owner diizinkan (endpoint didefinisikan admin only) - kemungkinan owner mendapat 403, cek apakah ini sesuai ekspektasi bisnis' },
    { role: 'Owner/Admin', scenario: 'Update profil toko (store settings)', steps: '1. Buka Settings > Profil Toko\n2. Ubah nama, alamat, telepon, email\n3. Simpan', data: 'store_name: Toko Mahenz, store_address: Jl. Contoh No.1, store_phone: 021-555123, store_email: info@mahenz.co.id', expected: 'Profil toko terupdate, tampil di struk transaksi' },
    { role: 'Owner/Admin', scenario: 'Update setting printer', steps: '1. Buka Settings > Printer\n2. Ubah konfigurasi printer\n3. Simpan', data: 'printer_name: EPSON-TM88, paper_size: 58mm', expected: 'Setting printer terupdate' },
    { role: 'Kasir', scenario: 'Lihat profil toko (view only, tidak bisa ubah)', steps: '1. Login kasir\n2. Buka Settings > Profil Toko', data: '-', expected: 'Data tampil read-only, tombol simpan tidak tersedia/403 jika dipaksa via API' },
  ]);

  addTestCaseSheet(wb, 'TC - Backup Restore', 'BACKUP & RESTORE', [
    { role: 'Owner/Admin', scenario: 'Buat backup database baru', steps: '1. Buka menu Backup\n2. Klik Backup Sekarang', data: '-', expected: 'File backup baru terbuat dan muncul di list' },
    { role: 'Owner/Admin', scenario: 'Lihat daftar file backup', steps: '1. Buka list backup', data: '-', expected: 'Semua file backup tampil dengan tanggal/ukuran' },
    { role: 'Owner/Admin', scenario: 'Download file backup', steps: '1. Klik download pada salah satu file backup', data: '-', expected: 'File terunduh dengan benar dan valid (tidak corrupt)' },
    { role: 'Admin', scenario: 'Restore database dari file backup (admin only)', steps: '1. Login admin\n2. Pilih file backup\n3. Klik Restore', data: 'file: backup_2026-07-01.sql', expected: 'Database berhasil dipulihkan sesuai isi backup - LAKUKAN DI ENVIRONMENT TEST, JANGAN DI PRODUKSI' },
    { role: 'Owner', scenario: 'Owner mencoba restore (harus ditolak jika endpoint admin only)', steps: '1. Login owner\n2. Panggil POST /restore', data: '-', expected: 'Verifikasi apakah owner mendapat 403 - cek konsistensi kebijakan akses' },
  ]);

  addTestCaseSheet(wb, 'TC - Sync', 'SYNC CENTER (Offline-first Desktop/Android)', [
    { role: 'Owner/Admin', scenario: 'Lihat daftar konflik sinkronisasi (conflicts)', steps: '1. Buka menu Sync Center > Conflicts', data: '-', expected: 'List konflik pending tampil dengan data desktop vs online' },
    { role: 'Owner/Admin', scenario: 'Lihat jumlah konflik (conflicts count)', steps: '1. Lihat badge/counter jumlah konflik', data: '-', expected: 'Angka sesuai jumlah data status=pending' },
    { role: 'Owner/Admin', scenario: 'Resolve konflik - pilih data desktop', steps: '1. Pilih salah satu konflik\n2. Pilih resolusi "desktop"', data: 'resolution: desktop, resolved_action: approve', expected: 'Data desktop dipakai sebagai final, status konflik menjadi resolved' },
    { role: 'Owner/Admin', scenario: 'Resolve konflik - pilih data online & reject', steps: '1. Pilih konflik lain\n2. Pilih resolusi "online", resolved_action: reject', data: 'resolution: online, resolved_action: reject', expected: 'Status konflik resolved sesuai pilihan' },
    { role: 'Owner/Admin', scenario: 'Lihat sync queue (item menunggu sinkron)', steps: '1. Buka menu Queue', data: '-', expected: 'List item dengan status pending/syncing/synced/failed tampil' },
    { role: 'Owner/Admin', scenario: 'Lihat riwayat sinkronisasi per device (history)', steps: '1. Buka menu History', data: '-', expected: 'Riwayat total/synced/conflict/failed items per device tampil' },
    { role: 'Semua device', scenario: 'Push data dari device desktop/android ke server', steps: '1. Simulasikan push payload dari device offline', data: 'device_id: DSK-001, entity_type: transaction, action: create', expected: 'Data masuk ke sync_queue dan diproses; jika ada konflik data, masuk ke sync_conflicts' },
  ]);

  addTestCaseSheet(wb, 'TC - Version', 'APP VERSION', [
    { role: 'Publik/App Android', scenario: 'Cek update versi Android (public)', steps: '1. Panggil GET /version/android tanpa login', data: '-', expected: 'Info versi terbaru (is_latest, is_mandatory, download_url) tampil tanpa perlu auth' },
    { role: 'Owner/Admin', scenario: 'Lihat daftar semua versi aplikasi', steps: '1. Buka menu App Version', data: '-', expected: 'List versi desktop & android tampil' },
    { role: 'Owner/Admin', scenario: 'Tambah/update versi Android baru dengan status mandatory', steps: '1. Tambah versi baru\n2. Set is_mandatory=true, is_latest=true', data: 'platform: android, version: 2.5.0, download_url: https://.../app-2.5.0.apk, is_mandatory: true', expected: 'Versi baru tersimpan, versi sebelumnya is_latest diset false otomatis (verifikasi behaviour)' },
    { role: 'Kasir', scenario: 'Kasir mencoba akses menu App Version (harus ditolak)', steps: '1. Login kasir\n2. Akses /version/list', data: '-', expected: '403 Forbidden (owner/admin only)' },
  ]);

  addTestCaseSheet(wb, 'TC - Security', 'KEAMANAN & OTORISASI LINTAS ROLE (CROSS-CHECK)', [
    { role: 'Kasir', scenario: 'Kasir mengakses endpoint /users/* (User Management)', steps: '1. Login kasir\n2. Panggil semua endpoint /users/*', data: '-', expected: '403 Forbidden untuk semua aksi (hanya owner)' },
    { role: 'Kasir', scenario: 'Kasir mengakses endpoint /stock-mutations/*', steps: '1. Login kasir\n2. Panggil /stock-mutations/list', data: '-', expected: '403 Forbidden (owner/admin only)' },
    { role: 'Kasir', scenario: 'Kasir mengakses endpoint /finance/*', steps: '1. Login kasir\n2. Panggil /finance/summary', data: '-', expected: '403 Forbidden (owner/admin only)' },
    { role: 'Kasir', scenario: 'Kasir mengakses /sync/* selain push', steps: '1. Login kasir\n2. Panggil /sync/conflicts', data: '-', expected: '403 Forbidden (owner/admin only)' },
    { role: 'Kasir', scenario: 'Kasir mengakses /categories/list, /products/list, /customers/list, /suppliers/list (GAP: tidak ada permission middleware)', steps: '1. Login kasir\n2. Panggil masing-masing endpoint langsung via API/Postman', data: '-', expected: 'Dokumentasikan hasil aktual: apakah data tetap bisa diakses walau menu tidak tampil di FE - laporkan sebagai temuan jika dianggap gap keamanan' },
    { role: 'Kasir', scenario: 'Kasir mengakses /dashboard/* (GAP: tidak ada RoleMiddleware di BE)', steps: '1. Login kasir\n2. Panggil /dashboard/stats via API langsung', data: '-', expected: 'Dokumentasikan hasil aktual - laporkan sebagai temuan jika data sensitif bisa diakses kasir' },
    { role: 'Admin', scenario: 'Admin mengakses endpoint khusus owner (users, roles CUD, menus CUD)', steps: '1. Login admin\n2. Coba create/update/delete pada roles, menus, dan semua aksi users', data: '-', expected: '403 Forbidden untuk semua aksi tersebut, hanya view roles/menus yang diizinkan' },
    { role: 'Semua', scenario: 'Akses API tanpa token / dengan token invalid', steps: '1. Panggil endpoint protected tanpa header Authorization / dengan token acak', data: 'Authorization: Bearer invalidtoken123', expected: '401 Unauthenticated untuk semua endpoint protected' },
    { role: 'Semua', scenario: 'Akses API dengan token expired', steps: '1. Gunakan token yang sudah melewati waktu expired', data: '-', expected: '401 Unauthenticated, harus refresh token / login ulang' },
    { role: 'Semua', scenario: 'Verifikasi mismatch algoritma hash password (bcrypt vs argon2id) pada akun default', steps: '1. Coba login menggunakan akun seed default owner/owner123 dan admin/admin123', data: 'username: owner, password: owner123', expected: 'Pastikan login sukses; jika gagal karena mismatch hash (seed argon2id vs verifikasi bcrypt), laporkan sebagai bug ke tim dev sebelum lanjut testing modul lain' },
  ]);

  // ============ DUMMY DATA SHEETS ============

  addDummyDataSheet(wb, 'Dummy - Akun Login', 'AKUN LOGIN SIAP PAKAI', ['Role', 'Username', 'Password', 'Full Name', 'Status', 'Catatan'], [
    ['Owner', 'owner', 'owner123', 'Owner Toko', 'Aktif', 'Akun default hasil seed - verifikasi dulu bisa login (lihat catatan hash password di TC-Security)'],
    ['Admin', 'admin', 'admin123', 'Admin Toko', 'Aktif', 'Akun default hasil seed'],
    ['Kasir', 'kasir1', 'kasir123', 'Budi Kasir', 'Aktif', 'Harus dibuat manual dulu via User Management (tidak ada di seed default)'],
    ['Kasir', 'kasir2', 'kasir123', 'Siti Kasir', 'Aktif', 'Untuk uji multi-kasir / shift berbeda'],
    ['Kasir (nonaktif)', 'kasir3', 'kasir123', 'Andi Nonaktif', 'Nonaktif', 'Untuk uji skenario login akun nonaktif'],
    ['Supervisor (custom role)', 'supervisor1', 'super123', 'Rina Supervisor', 'Aktif', 'Untuk uji role custom dengan permission granular via Role Access'],
  ], [22, 16, 16, 20, 12, 46]);

  addDummyDataSheet(wb, 'Dummy - User Role', 'DATA USER & ROLE', ['username', 'password', 'full_name', 'role', 'is_active', 'pin'], [
    ['owner', 'owner123', 'Owner Toko', 'owner', 'true', '1234'],
    ['admin', 'admin123', 'Admin Toko', 'admin', 'true', '1234'],
    ['kasir1', 'kasir123', 'Budi Kasir', 'kasir', 'true', '5678'],
    ['kasir2', 'kasir123', 'Siti Kasir', 'kasir', 'true', '5678'],
    ['kasir3', 'kasir123', 'Andi Nonaktif', 'kasir', 'false', '9999'],
    ['supervisor1', 'super123', 'Rina Supervisor', 'supervisor (custom)', 'true', '1111'],
  ], [16, 16, 20, 22, 12, 10]);

  addDummyDataSheet(wb, 'Dummy - Kategori Satuan', 'DATA KATEGORI & SATUAN', ['Tipe', 'name', 'code/abbreviation', 'description', 'is_active'], [
    ['Kategori', 'Makanan', 'MKN', 'Produk makanan ringan & instan', 'true'],
    ['Kategori', 'Minuman', 'MIN', 'Produk minuman kemasan', 'true'],
    ['Kategori', 'Kebutuhan Rumah Tangga', 'KRT', 'Sabun, deterjen, dll', 'true'],
    ['Kategori', 'Rokok', 'RKK', 'Produk rokok', 'true'],
    ['Kategori', 'ATK', 'ATK', 'Alat tulis kantor', 'false'],
    ['Satuan', 'Pcs', 'Pcs', 'Satuan per buah/pieces', 'true'],
    ['Satuan', 'Box', 'Box', 'Satuan per box/kotak', 'true'],
    ['Satuan', 'Karton', 'Krt', 'Satuan per karton (isi banyak box)', 'true'],
    ['Satuan', 'Lusin', 'Lsn', 'Satuan 12 pcs', 'true'],
    ['Satuan', 'Kg', 'Kg', 'Satuan kilogram', 'true'],
  ], [12, 26, 20, 34, 10]);

  addDummyDataSheet(wb, 'Dummy - Produk', 'DATA PRODUK', ['barcode', 'sku', 'name', 'category', 'unit', 'purchase_price', 'selling_price', 'stock', 'min_stock', 'is_active'], [
    ['8991002100016', 'SKU-0001', 'Indomie Goreng', 'Makanan', 'Pcs', '2500', '3000', '200', '20', 'true'],
    ['8992388123456', 'SKU-0002', 'Teh Botol Sosro 450ml', 'Minuman', 'Pcs', '4000', '5000', '150', '15', 'true'],
    ['8996001600028', 'SKU-0003', 'Aqua Botol 600ml', 'Minuman', 'Pcs', '2800', '3500', '300', '30', 'true'],
    ['8999999000017', 'SKU-0004', 'Rinso Anti Noda 800g', 'Kebutuhan Rumah Tangga', 'Pcs', '15000', '18500', '50', '5', 'true'],
    ['8991234567890', 'SKU-0005', 'Sampoerna A Mild 16', 'Rokok', 'Pcs', '22000', '24000', '80', '10', 'true'],
    ['8990001112223', 'SKU-0006', 'Chitato Sapi Panggang 68g', 'Makanan', 'Pcs', '8500', '10000', '5', '10', 'true - stok di bawah min_stock (untuk uji notifikasi stok)'],
    ['8990004445556', 'SKU-0007', 'Kopi Kapal Api Sachet', 'Minuman', 'Pcs', '1200', '1500', '0', '20', 'true - stok habis (untuk uji stok 0)'],
    ['8990007778889', 'SKU-0008', 'Produk Nonaktif Test', 'Makanan', 'Pcs', '1000', '1500', '10', '5', 'false - untuk uji produk nonaktif tidak muncul di kasir'],
  ], [16, 12, 28, 22, 10, 14, 14, 10, 10, 40]);

  const wsPackages = addDummyDataSheet(wb, 'Dummy - Produk Paket', 'DATA PAKET/VARIAN SATUAN PRODUK (lanjutan sheet Produk)', ['product', 'package_name', 'unit', 'conversion_qty', 'selling_price', 'purchase_price', 'is_default'], [
    ['Indomie Goreng', '(satuan dasar)', 'Pcs', '1', '3000', '2500', 'true'],
    ['Indomie Goreng', 'Dus (isi 40)', 'Karton', '40', '110000', '95000', 'false'],
    ['Aqua Botol 600ml', '(satuan dasar)', 'Pcs', '1', '3500', '2800', 'true'],
    ['Aqua Botol 600ml', 'Dus (isi 24)', 'Karton', '24', '78000', '65000', 'false'],
  ], [24, 24, 10, 16, 14, 14, 12]);
  wsPackages.getRow(2).font = { italic: true, size: 10 };

  addDummyDataSheet(wb, 'Dummy - Produk Grosir', 'DATA HARGA GROSIR (TIER PRICE) - lanjutan sheet Produk', ['product', 'tier_name', 'min_qty', 'price'], [
    ['Indomie Goreng', 'Grosir 12+', '12', '2800'],
    ['Indomie Goreng', 'Grosir 40+ (1 dus)', '40', '2700'],
    ['Aqua Botol 600ml', 'Grosir 24+', '24', '3200'],
  ], [24, 24, 12, 12]);

  addDummyDataSheet(wb, 'Dummy - Supplier', 'DATA SUPPLIER', ['supplier_code', 'name', 'address', 'phone', 'email', 'contact_person', 'is_active'], [
    ['SUP-001', 'PT Sumber Makmur', 'Jl. Industri No.10, Jakarta', '081234567890', 'sumber@makmur.co.id', 'Andi Wijaya', 'true'],
    ['SUP-002', 'CV Berkat Jaya', 'Jl. Raya Bogor No.25', '082198765432', 'berkat@jaya.co.id', 'Dewi Lestari', 'true'],
    ['SUP-003', 'Distributor Rokok Nusantara', 'Jl. Gatot Subroto No.7', '081211122233', 'info@rokoknusantara.co.id', 'Hendra', 'true'],
    ['SUP-004', 'Supplier Nonaktif Test', 'Jl. Test No.1', '080000000000', 'test@test.com', 'Test', 'false'],
  ], [14, 26, 30, 16, 26, 18, 10]);

  addDummyDataSheet(wb, 'Dummy - Pembelian', 'DATA PEMBELIAN (PURCHASE) - lanjutan sheet Supplier', ['purchase_code', 'supplier', 'purchase_date', 'items (product x qty @price)', 'discount_amount', 'total_amount', 'payment_status', 'paid_amount'], [
    ['PB-0001', 'PT Sumber Makmur', '2026-06-25', 'Indomie Goreng x100 @2500; Teh Botol x50 @4000', '10000', '440000', 'paid', '440000'],
    ['PB-0002', 'CV Berkat Jaya', '2026-06-28', 'Aqua Botol x150 @2800', '0', '420000', 'partial', '200000'],
    ['PB-0003', 'Distributor Rokok Nusantara', '2026-07-01', 'Sampoerna A Mild x80 @22000', '20000', '1740000', 'unpaid', '0'],
  ], [14, 24, 14, 44, 14, 14, 14, 14]);

  addDummyDataSheet(wb, 'Dummy - Pelanggan', 'DATA PELANGGAN', ['customer_code', 'name', 'phone', 'address', 'credit_limit', 'is_active'], [
    ['CUST-001', 'Ibu Sari', '082211223344', 'Jl. Melati No.5', '500000', 'true'],
    ['CUST-002', 'Bapak Joko', '081199887766', 'Jl. Anggrek No.12', '1000000', 'true'],
    ['CUST-003', 'Toko Kelontong Makmur', '081277889900', 'Jl. Pasar Baru No.3', '2000000', 'true'],
    ['CUST-004', 'Pelanggan Nonaktif Test', '080000000001', 'Jl. Test No.2', '0', 'false'],
  ], [14, 26, 16, 26, 14, 10]);

  addDummyDataSheet(wb, 'Dummy - Piutang', 'DATA PIUTANG (RECEIVABLE) - lanjutan sheet Pelanggan', ['customer', 'transaction_ref', 'total_amount', 'paid_amount', 'remaining_amount', 'status', 'due_date'], [
    ['Ibu Sari', 'TRX-0010', '150000', '0', '150000', 'unpaid', '2026-07-15'],
    ['Bapak Joko', 'TRX-0015', '500000', '250000', '250000', 'partial', '2026-07-10'],
    ['Toko Kelontong Makmur', 'TRX-0020', '1200000', '1200000', '0', 'paid', '2026-06-30'],
  ], [24, 16, 14, 14, 16, 12, 14]);

  addDummyDataSheet(wb, 'Dummy - Shift', 'DATA SHIFT', ['name', 'start_time', 'end_time', 'is_active'], [
    ['Shift Pagi', '07:00', '15:00', 'true'],
    ['Shift Siang', '15:00', '22:00', 'true'],
    ['Shift Malam', '22:00', '07:00', 'true'],
  ], [16, 14, 14, 10]);

  addDummyDataSheet(wb, 'Dummy - Kas Harian', 'DATA KAS HARIAN (CASH DRAWER) - lanjutan sheet Shift', ['kasir', 'shift', 'opening_balance', 'closing_balance', 'expected_balance', 'difference', 'status', 'open_notes'], [
    ['kasir1', 'Shift Pagi', '500000', '1250000', '1250000', '0', 'closed', 'Kas awal shift pagi normal'],
    ['kasir2', 'Shift Siang', '300000', '980000', '1000000', '-20000', 'closed', 'Ada selisih kurang 20.000'],
    ['kasir1', 'Shift Malam', '400000', '', '', '', 'open', 'Kas sedang berjalan (belum ditutup)'],
  ], [12, 14, 16, 16, 16, 14, 12, 34]);

  addDummyDataSheet(wb, 'Dummy - Pengeluaran', 'DATA PENGELUARAN (EXPENSE)', ['expense_date', 'category', 'description', 'amount', 'payment_method', 'user'], [
    ['2026-07-01', 'Operasional', 'Beli alat tulis kantor', '75000', 'cash', 'kasir1'],
    ['2026-07-01', 'Listrik & Air', 'Bayar listrik bulanan', '450000', 'transfer', 'admin'],
    ['2026-07-02', 'Transportasi', 'Ongkos kirim barang retur', '50000', 'cash', 'kasir2'],
    ['2026-07-02', 'Konsumsi', 'Makan siang karyawan', '120000', 'cash', 'owner'],
  ], [14, 16, 32, 14, 16, 12]);

  addDummyDataSheet(wb, 'Dummy - Transaksi', 'DATA TRANSAKSI KASIR (POS SALE)', ['transaction_code', 'kasir', 'items (product x qty @price)', 'payment_method', 'payment_amount', 'total_amount', 'customer', 'is_credit', 'status', 'device_source'], [
    ['TRX-0001', 'kasir1', 'Indomie Goreng x3 @3000', 'cash', '10000', '9000', '-', 'false', 'completed', 'web'],
    ['TRX-0002', 'kasir1', 'Aqua Botol x2 @3500; Teh Botol x1 @5000', 'qris', '12000', '12000', '-', 'false', 'completed', 'android'],
    ['TRX-0003', 'kasir2', 'Rinso Anti Noda x1 @18500', 'transfer', '18500', '18500', '-', 'false', 'completed', 'desktop'],
    ['TRX-0010', 'kasir1', 'Sampoerna A Mild x5 @24000; barang diambil kredit', 'kredit', '0', '150000', 'Ibu Sari', 'true', 'completed', 'web'],
    ['TRX-0020', 'kasir2', 'Indomie Goreng x40 (1 dus) @110000', 'cash', '110000', '110000', '-', 'false', 'completed', 'web'],
    ['TRX-0099', 'kasir1', 'Indomie Goreng x2 @3000 (contoh transaksi untuk di-void)', 'cash', '6000', '6000', '-', 'false', 'void (setelah di-void owner)', 'web'],
  ], [16, 10, 46, 14, 14, 14, 14, 10, 24, 12]);

  const outPath = path.join(__dirname, 'Dokumen_Testing_POS.xlsx');
  await wb.xlsx.writeFile(outPath);
  console.log('Excel file created at:', outPath);
}

main().catch((err) => {
  console.error(err);
  process.exit(1);
});

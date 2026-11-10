// Migrasi sekali jalan nama key localStorage ke format `<aplikasi>.<nama>`.
//
// WAJIB di-import paling awal di main.tsx, sebelum modul store (zustand persist) dievaluasi,
// karena persist membaca localStorage saat store dibuat. Dengan begitu kasir yang sedang login
// tidak ter-logout, keranjang yang belum selesai tidak hilang, dan printer Bluetooth tablet
// tetap diingat.

const RENAMES: Record<string, string> = {
  'auth-session': 'toko.auth-session',
  'cashier-draft': 'toko.cashier-draft',
  'backdate-cashier-draft': 'toko.backdate-cashier-draft',
  'pos.native_bt_printer_address': 'toko.native_bt_printer_address',
}

function migrate(): void {
  try {
    for (const [oldKey, newKey] of Object.entries(RENAMES)) {
      const oldValue = localStorage.getItem(oldKey)
      if (oldValue === null) continue
      // Nilai baru yang sudah ada tidak ditimpa; key lama tetap dibersihkan.
      if (localStorage.getItem(newKey) === null) {
        localStorage.setItem(newKey, oldValue)
      }
      localStorage.removeItem(oldKey)
    }
  } catch {
    // localStorage tidak tersedia (mode privat/kuota): lewati, aplikasi tetap jalan.
  }
}

migrate()

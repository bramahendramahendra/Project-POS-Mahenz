import type { CapacitorConfig } from '@capacitor/cli'

/**
 * Konfigurasi Capacitor untuk APK POS Mahenz (Opsi A).
 *
 * WebView membuka web produksi langsung (server.url), jadi APK hanya "cangkang"
 * + jembatan Bluetooth native. Update tampilan/fitur web tidak perlu rebuild APK.
 *
 * Catatan:
 * - server.url WAJIB HTTPS (https://mahenz.duckdns.org sudah aktif) agar WebView
 *   mau membukanya tanpa mengaktifkan cleartext.
 * - androidScheme 'https' agar origin WebView konsisten dengan domain.
 * - webDir menunjuk 'www' (placeholder) — pada Opsi A isinya tidak dipakai WebView
 *   karena memuat dari server.url, tetapi `cap sync` tetap butuh folder ini ada.
 */
const config: CapacitorConfig = {
  appId: 'id.mahenz.pos',
  appName: 'POS Mahenz',
  webDir: 'www',
  server: {
    url: 'https://mahenz.duckdns.org',
    cleartext: false,
    androidScheme: 'https',
  },
}

export default config

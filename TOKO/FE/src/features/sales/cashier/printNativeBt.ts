/**
 * Jembatan cetak Bluetooth NATIVE untuk APK (Capacitor).
 *
 * Dipakai HANYA saat aplikasi berjalan di dalam APK (WebView Capacitor). Di
 * browser biasa, modul ini aman: semua akses ke Capacitor dibungkus pengecekan,
 * dan `isNativeApp()` mengembalikan false sehingga metode ini tidak ditampilkan.
 *
 * Desain (Cara A): FE TIDAK menambah dependency npm Capacitor. Plugin diakses
 * lewat objek global `window.Capacitor` yang disuntikkan WebView saat runtime.
 * Semua "barang Android" (plugin native, Gradle, build) ada di folder `mobile/`.
 *
 * Plugin yang dipakai: @kduma-autoid/capacitor-bluetooth-printer
 *   - list(): daftar device Bluetooth ter-pair
 *   - connectAndPrint({ address, data }): konek + cetak (data = base64 ESC/POS)
 */
import { buildReceiptEscposBytes, type ReceiptData } from './blePrinter'

// ---- Tipe minimal untuk objek global Capacitor (runtime di WebView) ----

type BluetoothDeviceType = 'unknown' | 'classic' | 'le' | 'dual'

interface NativeBluetoothDevice {
  name: string
  address: string
  type: BluetoothDeviceType
}

interface BluetoothPrinterPlugin {
  list(): Promise<{ devices: NativeBluetoothDevice[] }>
  connect(options: { address: string }): Promise<void>
  disconnect(): Promise<void>
  print(options: { data: string }): Promise<void>
  connectAndPrint(options: { address: string; data: string }): Promise<void>
}

interface CapacitorGlobal {
  isNativePlatform?: () => boolean
  Plugins?: {
    BluetoothPrinter?: BluetoothPrinterPlugin
  }
}

function getCapacitor(): CapacitorGlobal | null {
  if (typeof window === 'undefined') return null
  const cap = (window as unknown as { Capacitor?: CapacitorGlobal }).Capacitor
  return cap ?? null
}

/** True hanya bila aplikasi berjalan di dalam APK Capacitor (bukan browser). */
export function isNativeApp(): boolean {
  const cap = getCapacitor()
  return !!cap && typeof cap.isNativePlatform === 'function' && cap.isNativePlatform()
}

function getPlugin(): BluetoothPrinterPlugin | null {
  return getCapacitor()?.Plugins?.BluetoothPrinter ?? null
}

/** Encode byte ke base64 standar (plugin menerima string base64). */
function bytesToBase64(bytes: Uint8Array): string {
  let binary = ''
  for (let i = 0; i < bytes.length; i++) binary += String.fromCharCode(bytes[i])
  return btoa(binary)
}

// Mengingat printer terakhir yang dipilih user, agar cetak berikutnya tidak perlu
// memilih ulang. Disimpan di localStorage (per perangkat).
const LAST_PRINTER_KEY = 'toko.native_bt_printer_address'

function getRememberedAddress(): string | null {
  try {
    return localStorage.getItem(LAST_PRINTER_KEY)
  } catch {
    return null
  }
}

function rememberAddress(address: string): void {
  try {
    localStorage.setItem(LAST_PRINTER_KEY, address)
  } catch {
    // abaikan bila localStorage tidak tersedia
  }
}

/**
 * Pilih printer Bluetooth: pakai alamat yang diberikan, atau yang terakhir
 * diingat, atau otomatis pilih SATU-satunya printer bila hanya ada satu device
 * ter-pair. Mengembalikan alamat printer, atau melempar error dengan pesan jelas.
 */
async function resolvePrinterAddress(
  plugin: BluetoothPrinterPlugin,
  preferred?: string,
): Promise<string> {
  if (preferred) return preferred

  const remembered = getRememberedAddress()
  const { devices } = await plugin.list()

  if (!devices || devices.length === 0) {
    throw new Error('Tidak ada printer Bluetooth yang dipasangkan. Pair printer EP58M di Pengaturan Bluetooth tablet dulu.')
  }

  // Bila alamat yang diingat masih ada di daftar, pakai itu.
  if (remembered && devices.some((d) => d.address === remembered)) {
    return remembered
  }

  // Hanya satu device → pakai langsung.
  if (devices.length === 1) return devices[0].address

  // Banyak device & belum ada yang diingat: pakai yang pertama bertipe classic
  // (printer SPP seperti EP58M), jika tidak ada ambil device pertama.
  const classic = devices.find((d) => d.type === 'classic' || d.type === 'dual')
  return (classic ?? devices[0]).address
}

/**
 * Cetak struk langsung ke printer Bluetooth lewat plugin native (APK).
 * @param data  data struk siap tampil.
 * @param address  alamat printer opsional; bila kosong dipilih otomatis.
 */
export async function printViaNativeBt(data: ReceiptData, address?: string): Promise<void> {
  const plugin = getPlugin()
  if (!plugin) {
    throw new Error('Jembatan Bluetooth tidak tersedia. Fitur ini hanya jalan di aplikasi (APK), bukan browser.')
  }

  const target = await resolvePrinterAddress(plugin, address)
  const payload = bytesToBase64(buildReceiptEscposBytes(data))

  await plugin.connectAndPrint({ address: target, data: payload })
  rememberAddress(target)
}

import { formatRupiah } from '@/shared/utils'

import type { CartItem, CartSummary, CheckoutResponse, Discount, PaymentMethod, Tax } from './cashier.types'

const PAYMENT_LABELS: Record<PaymentMethod, string> = {
  cash: 'Tunai',
  transfer: 'Transfer',
  qris: 'QRIS',
  card: 'Kartu',
  kredit: 'Kredit',
  balance: 'Saldo',
}

const ESC_INIT = new Uint8Array([0x1b, 0x40]) // ESC @ — reset printer
const CUT_PAPER = new Uint8Array([0x1d, 0x56, 0x00]) // GS V 0 — full cut
const BLE_CHUNK_SIZE = 20 // MTU default BLE, aman untuk semua printer

// Perintah ESC/POS untuk format teks native (hasil tajam, sama seperti e-nota).
const ESC_ALIGN_LEFT = [0x1b, 0x61, 0x00]
const ESC_ALIGN_CENTER = [0x1b, 0x61, 0x01]
const ESC_BOLD_ON = [0x1b, 0x45, 0x01]
const ESC_BOLD_OFF = [0x1b, 0x45, 0x00]
const ESC_FEED_LINES = (n: number) => [0x1b, 0x64, n] // ESC d n — feed n baris

export function isBlePrintSupported(): boolean {
  return typeof navigator !== 'undefined' && 'bluetooth' in navigator && !!navigator.bluetooth
}

/** Printer BLE yang browser sudah punya izin persist (pernah dipilih user sebelumnya). */
export async function getRememberedPrinter(): Promise<BluetoothDevice | null> {
  if (!isBlePrintSupported()) return null
  try {
    const devices = await navigator.bluetooth!.getDevices()
    return devices[0] ?? null
  } catch {
    return null
  }
}

/** Memicu picker native browser untuk memilih printer BLE — dipanggil dari klik tombol eksplisit. */
export async function pickAndRememberPrinter(): Promise<BluetoothDevice> {
  if (!isBlePrintSupported()) {
    throw new Error('Browser ini tidak mendukung Web Bluetooth.')
  }
  return navigator.bluetooth!.requestDevice({
    acceptAllDevices: true,
    optionalServices: [
      '000018f0-0000-1000-8000-00805f9b34fb', // service umum dipakai printer thermal BLE generik
      '49535343-fe7d-4ae5-8fa9-9fafd205e455', // varian umum lain (ISSC transparent UART)
    ],
  })
}

function findWritableCharacteristic(
  services: BluetoothRemoteGATTService[]
): Promise<BluetoothRemoteGATTCharacteristic | null> {
  return (async () => {
    for (const service of services) {
      const chars = await service.getCharacteristics().catch(() => [])
      const writable = chars.find((c) => c.properties.write || c.properties.writeWithoutResponse)
      if (writable) return writable
    }
    return null
  })()
}

async function writeInChunks(char: BluetoothRemoteGATTCharacteristic, bytes: Uint8Array): Promise<void> {
  const write = char.properties.writeWithoutResponse
    ? char.writeValueWithoutResponse.bind(char)
    : char.writeValue.bind(char)
  for (let i = 0; i < bytes.length; i += BLE_CHUNK_SIZE) {
    await write(bytes.slice(i, i + BLE_CHUNK_SIZE))
  }
}

/** Kirim teks struk sebagai perintah ESC/POS mentah ke printer BLE yang sudah dipilih. */
export async function printViaBle(device: BluetoothDevice, text: string): Promise<void> {
  if (!device.gatt) throw new Error('Device tidak punya GATT server.')
  const server = device.gatt.connected ? device.gatt : await device.gatt.connect()
  try {
    const services = await server.getPrimaryServices()
    const char = await findWritableCharacteristic(services)
    if (!char) throw new Error('Tidak ditemukan characteristic yang bisa ditulis pada printer ini.')

    const body = new TextEncoder().encode(text)
    const payload = new Uint8Array(ESC_INIT.length + body.length + CUT_PAPER.length)
    payload.set(ESC_INIT, 0)
    payload.set(body, ESC_INIT.length)
    payload.set(CUT_PAPER, ESC_INIT.length + body.length)

    await writeInChunks(char, payload)
  } finally {
    server.disconnect()
  }
}

/**
 * Kirim byte ESC/POS yang SUDAH jadi (dari buildReceiptEscposBytes) ke printer
 * BLE. Berbeda dari printViaBle yang menerima teks lalu membungkusnya — di sini
 * byte sudah lengkap termasuk init/align/bold/feed/cut, jadi hasilnya identik
 * dengan jalur RawBT & Web Serial.
 */
export async function printViaBleBytes(device: BluetoothDevice, bytes: Uint8Array): Promise<void> {
  if (!device.gatt) throw new Error('Device tidak punya GATT server.')
  const server = device.gatt.connected ? device.gatt : await device.gatt.connect()
  try {
    const services = await server.getPrimaryServices()
    const char = await findWritableCharacteristic(services)
    if (!char) throw new Error('Tidak ditemukan characteristic yang bisa ditulis pada printer ini.')
    await writeInChunks(char, bytes)
  } finally {
    server.disconnect()
  }
}

function padLine(left: string, right: string, width: number): string {
  const space = Math.max(1, width - left.length - right.length)
  return left + ' '.repeat(space) + right + '\n'
}

function center(text: string, width: number): string {
  const space = Math.max(0, Math.floor((width - text.length) / 2))
  return ' '.repeat(space) + text + '\n'
}

/** Susun ulang data struk jadi teks polos rata-kolom monospace untuk dikirim ke printer ESC/POS. */
export function receiptToPlainText(params: {
  storeName: string
  storeSub?: string
  footer?: string
  paperSize: '58mm' | '80mm'
  checkoutData: CheckoutResponse
  cart: CartItem[]
  summary: CartSummary
  discount: Discount
  tax: Tax
  paymentMethod: PaymentMethod
  amountPaid: number
  customerName?: string
}): string {
  const { storeName, storeSub, footer, paperSize, checkoutData, cart, summary, discount, tax, paymentMethod, amountPaid, customerName } = params
  const width = paperSize === '58mm' ? 32 : 48
  const line = '-'.repeat(width) + '\n'
  const change = Math.max(0, amountPaid - summary.grandTotal)

  let out = center(storeName, width)
  if (storeSub) out += center(storeSub, width)
  out += line
  out += padLine('No. Transaksi', checkoutData.transaction_code, width)
  out += padLine('Tanggal', new Date(checkoutData.transaction_date).toLocaleString('id-ID'), width)
  if (customerName) out += padLine('Pelanggan', customerName, width)
  out += padLine('Pembayaran', PAYMENT_LABELS[paymentMethod], width)
  out += line

  for (const item of cart) {
    const price = item.effective_price ?? item.price
    out += padLine(item.product_name, formatRupiah(item.subtotal), width)
    out += `  ${item.unit_name} x${item.qty} @ ${formatRupiah(price)}\n`
  }
  out += line

  out += padLine('Subtotal', formatRupiah(summary.subtotal), width)
  if (summary.discountAmount > 0) {
    out += padLine(`Diskon${discount.type === 'percent' ? ` (${discount.value}%)` : ''}`, `-${formatRupiah(summary.discountAmount)}`, width)
  }
  if (summary.taxAmount > 0) {
    out += padLine(`Pajak (${tax.percent}%)`, `+${formatRupiah(summary.taxAmount)}`, width)
  }
  out += padLine('TOTAL', formatRupiah(summary.grandTotal), width)
  out += padLine(`Dibayar (${PAYMENT_LABELS[paymentMethod]})`, formatRupiah(amountPaid), width)
  out += padLine('Kembalian', formatRupiah(change), width)
  out += line

  if (footer) out += center(footer, width)
  out += '\n\n\n'
  return out
}

// ===========================================================================
// Cetak ESC/POS native lintas-jalur (Web Bluetooth & Web Serial) — hasil setajam
// e-nota karena printer memakai font internalnya, bukan render grafis browser.
// ===========================================================================

export type ReceiptData = {
  storeName: string
  storeSub?: string
  footer?: string
  paperSize: '58mm' | '80mm'
  checkoutData: CheckoutResponse
  cart: CartItem[]
  summary: CartSummary
  discount: Discount
  tax: Tax
  paymentMethod: PaymentMethod
  amountPaid: number
  customerName?: string
}

export type DevicePlatform = 'android' | 'windows' | 'other'

/** Deteksi platform perangkat dari user agent untuk memilih jalur cetak yang tepat. */
export function getPlatform(): DevicePlatform {
  if (typeof navigator === 'undefined') return 'other'
  const ua = navigator.userAgent.toLowerCase()
  if (ua.includes('android')) return 'android'
  if (ua.includes('windows')) return 'windows'
  return 'other'
}

/**
 * Bangun byte ESC/POS lengkap dari data struk, dengan perintah format native
 * (align center untuk header, bold untuk TOTAL, feed + cut di akhir). Teks tubuh
 * memakai layout rata-kolom monospace yang sama dengan receiptToPlainText.
 */
export function buildReceiptEscposBytes(data: ReceiptData): Uint8Array {
  const width = data.paperSize === '58mm' ? 32 : 48
  const lineSep = '-'.repeat(width) + '\n'
  const change = Math.max(0, data.amountPaid - data.summary.grandTotal)
  const enc = new TextEncoder()
  const parts: number[] = []

  const push = (bytes: number[]) => parts.push(...bytes)
  const pushText = (text: string) => parts.push(...enc.encode(text))

  push([...ESC_INIT])

  // Header toko — center, nama tebal.
  push(ESC_ALIGN_CENTER)
  push(ESC_BOLD_ON)
  pushText(data.storeName + '\n')
  push(ESC_BOLD_OFF)
  if (data.storeSub) pushText(data.storeSub + '\n')
  push(ESC_ALIGN_LEFT)

  pushText(lineSep)
  pushText(padLine('No. Transaksi', data.checkoutData.transaction_code, width))
  pushText(padLine('Tanggal', new Date(data.checkoutData.transaction_date).toLocaleString('id-ID'), width))
  if (data.customerName) pushText(padLine('Pelanggan', data.customerName, width))
  pushText(padLine('Pembayaran', PAYMENT_LABELS[data.paymentMethod], width))
  pushText(lineSep)

  for (const item of data.cart) {
    const price = item.effective_price ?? item.price
    pushText(padLine(item.product_name, formatRupiah(item.subtotal), width))
    pushText(`  ${item.unit_name} x${item.qty} @ ${formatRupiah(price)}\n`)
  }
  pushText(lineSep)

  pushText(padLine('Subtotal', formatRupiah(data.summary.subtotal), width))
  if (data.summary.discountAmount > 0) {
    pushText(padLine(`Diskon${data.discount.type === 'percent' ? ` (${data.discount.value}%)` : ''}`, `-${formatRupiah(data.summary.discountAmount)}`, width))
  }
  if (data.summary.taxAmount > 0) {
    pushText(padLine(`Pajak (${data.tax.percent}%)`, `+${formatRupiah(data.summary.taxAmount)}`, width))
  }
  // TOTAL — tebal agar menonjol.
  push(ESC_BOLD_ON)
  pushText(padLine('TOTAL', formatRupiah(data.summary.grandTotal), width))
  push(ESC_BOLD_OFF)
  pushText(padLine(`Dibayar (${PAYMENT_LABELS[data.paymentMethod]})`, formatRupiah(data.amountPaid), width))
  pushText(padLine('Kembalian', formatRupiah(change), width))
  pushText(lineSep)

  if (data.footer) {
    push(ESC_ALIGN_CENTER)
    pushText(data.footer + '\n')
    push(ESC_ALIGN_LEFT)
  }

  // Feed beberapa baris lalu potong (jika printer mendukung auto-cut).
  push(ESC_FEED_LINES(3))
  push([...CUT_PAPER])

  return new Uint8Array(parts)
}

// ===========================================================================
// Web Serial API — cetak langsung ke COM port dari BROWSER di mesin kasir
// (laptop/PC Windows). Ini pengganti jalur lama yang menitipkan cetak ke BE:
// di production BE ada di server cloud yang tidak punya akses ke printer, jadi
// COM port HARUS dibuka dari sisi klien. Web Serial hanya tersedia di konteks
// aman (HTTPS atau localhost) pada Chrome/Edge desktop.
// ===========================================================================

// Minimal typing agar tidak bergantung pada lib DOM "serial" yang mungkin belum
// aktif di tsconfig. Hanya method yang kita pakai yang dideklarasikan.
interface MinimalSerialPort {
  open(options: { baudRate: number }): Promise<void>
  close(): Promise<void>
  readonly writable: WritableStream<Uint8Array> | null
}
interface MinimalSerial {
  getPorts(): Promise<MinimalSerialPort[]>
  requestPort(): Promise<MinimalSerialPort>
}

function getSerial(): MinimalSerial | null {
  const nav = navigator as unknown as { serial?: MinimalSerial }
  return nav.serial ?? null
}

/** True bila browser mendukung Web Serial (Chrome/Edge desktop di konteks aman). */
export function isWebSerialSupported(): boolean {
  return typeof navigator !== 'undefined' && getSerial() !== null
}

/**
 * Minta user memilih COM port printer (memicu dialog native browser). Harus
 * dipanggil dari gesture klik. Browser mengingat izinnya sehingga pemanggilan
 * berikutnya bisa pakai getRememberedSerialPort tanpa dialog lagi.
 */
export async function pickSerialPort(): Promise<MinimalSerialPort> {
  const serial = getSerial()
  if (!serial) throw new Error('Browser ini tidak mendukung Web Serial (butuh Chrome/Edge + HTTPS).')
  return serial.requestPort()
}

/** Port serial yang izinnya sudah pernah diberikan user sebelumnya (bila ada). */
export async function getRememberedSerialPort(): Promise<MinimalSerialPort | null> {
  const serial = getSerial()
  if (!serial) return null
  try {
    const ports = await serial.getPorts()
    return ports[0] ?? null
  } catch {
    return null
  }
}

/** Tulis byte ESC/POS ke COM port lewat Web Serial. baudRate umum printer: 9600. */
export async function printViaWebSerial(port: MinimalSerialPort, bytes: Uint8Array): Promise<void> {
  await port.open({ baudRate: 9600 })
  try {
    if (!port.writable) throw new Error('Port serial tidak bisa ditulis.')
    const writer = port.writable.getWriter()
    try {
      await writer.write(bytes)
      // Jeda singkat agar buffer terkirim sebelum port ditutup.
      await new Promise((r) => setTimeout(r, 400))
    } finally {
      writer.releaseLock()
    }
  } finally {
    await port.close()
  }
}

/**
 * Bangun HTML isi struk (untuk dicetak via iframe, lihat printFrame.ts).
 * Memakai kelas .center/.bold/.row/.sep/.total yang didefinisikan di iframe.
 */
export function buildReceiptHtml(data: ReceiptData): string {
  const esc = (s: string) =>
    s.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;')
  const change = Math.max(0, data.amountPaid - data.summary.grandTotal)
  const tanggal = new Date(data.checkoutData.transaction_date).toLocaleString('id-ID')

  const rows: string[] = []

  rows.push(`<div class="center store-name">${esc(data.storeName)}</div>`)
  if (data.storeSub) rows.push(`<div class="center">${esc(data.storeSub)}</div>`)
  rows.push('<hr />')

  rows.push(`<div class="info"><span class="lbl">No:</span> <span class="val">${esc(data.checkoutData.transaction_code)}</span></div>`)
  rows.push(`<div class="info"><span class="lbl">Tgl:</span> <span class="val">${esc(tanggal)}</span></div>`)
  if (data.customerName) rows.push(`<div class="info"><span class="lbl">Plg:</span> <span class="val">${esc(data.customerName)}</span></div>`)
  rows.push(`<div class="info"><span class="lbl">Bayar:</span> <span class="val">${esc(PAYMENT_LABELS[data.paymentMethod])}</span></div>`)
  rows.push('<hr />')

  for (const item of data.cart) {
    const price = item.effective_price ?? item.price
    rows.push(`<div class="item-name">${esc(item.product_name)}</div>`)
    rows.push(`<div class="row"><span class="item-sub">${esc(item.unit_name)} x${item.qty} @ ${esc(formatRupiah(price))}</span><span>${esc(formatRupiah(item.subtotal))}</span></div>`)
  }
  rows.push('<hr />')

  rows.push(`<div class="row"><span>Subtotal</span><span>${esc(formatRupiah(data.summary.subtotal))}</span></div>`)
  if (data.summary.discountAmount > 0) {
    const label = `Diskon${data.discount.type === 'percent' ? ` (${data.discount.value}%)` : ''}`
    rows.push(`<div class="row"><span>${esc(label)}</span><span>-${esc(formatRupiah(data.summary.discountAmount))}</span></div>`)
  }
  if (data.summary.taxAmount > 0) {
    rows.push(`<div class="row"><span>Pajak (${data.tax.percent}%)</span><span>+${esc(formatRupiah(data.summary.taxAmount))}</span></div>`)
  }
  rows.push('<div class="sep"></div>')
  rows.push(`<div class="row total"><span>TOTAL</span><span>${esc(formatRupiah(data.summary.grandTotal))}</span></div>`)
  rows.push(`<div class="row"><span>Dibayar</span><span>${esc(formatRupiah(data.amountPaid))}</span></div>`)
  rows.push(`<div class="row"><span>Kembalian</span><span>${esc(formatRupiah(change))}</span></div>`)
  rows.push('<hr />')

  if (data.footer) rows.push(`<div class="center">${esc(data.footer)}</div>`)

  return rows.join('\n')
}

// ===========================================================================
// RawBT — cetak lewat aplikasi RawBT di Android. Aplikasi web membentuk URL
// intent berisi data ESC/POS (base64) lalu mengarahkan browser ke sana; Android
// melempar ke RawBT yang mencetak via Bluetooth Classic. Hanya untuk Android.
// ===========================================================================

/** Deteksi perangkat Android (untuk menandai ketersediaan jalur RawBT). */
export function isAndroid(): boolean {
  return typeof navigator !== 'undefined' && /android/i.test(navigator.userAgent)
}

/** Encode byte ke base64 standar untuk disematkan di URL intent RawBT. */
function bytesToBase64(bytes: Uint8Array): string {
  let binary = ''
  for (let i = 0; i < bytes.length; i++) {
    binary += String.fromCharCode(bytes[i])
  }
  return btoa(binary)
}

/**
 * Cetak lewat aplikasi RawBT di Android. Membentuk URL intent sesuai skema resmi
 * RawBT lalu mengarahkan browser ke sana. Melempar error bila bukan Android
 * (tidak ada handler intent rawbt di platform lain).
 */
export function printViaRawBt(bytes: Uint8Array): void {
  if (!isAndroid()) {
    throw new Error('Cetak via RawBT hanya tersedia di perangkat Android.')
  }
  const b64 = bytesToBase64(bytes)
  const intentUrl = `intent:base64,${b64}#Intent;scheme=rawbt;package=ru.a402d.rawbtprinter;end;`
  window.location.href = intentUrl
}

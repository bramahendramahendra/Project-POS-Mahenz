import { formatRupiah } from '@/shared/utils'

import type { CartItem, CartSummary, CheckoutResponse, Discount, PaymentMethod, Tax } from './cashier.types'

const PAYMENT_LABELS: Record<PaymentMethod, string> = {
  cash: 'Tunai',
  transfer: 'Transfer',
  qris: 'QRIS',
  card: 'Kartu',
  kredit: 'Kredit',
}

const ESC_INIT = new Uint8Array([0x1b, 0x40]) // ESC @ — reset printer
const CUT_PAPER = new Uint8Array([0x1d, 0x56, 0x00]) // GS V 0 — full cut
const BLE_CHUNK_SIZE = 20 // MTU default BLE, aman untuk semua printer

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

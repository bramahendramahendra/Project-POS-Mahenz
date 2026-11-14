import { useCallback, useState } from 'react'
import { Printer, ShoppingCart, X, Bluetooth, Monitor, FileText, Usb, Server, Smartphone, Wifi, MessageCircle } from 'lucide-react'

import { ActionModal } from '@/shared/components'
import { Button } from '@/shared/components/ui/button'
import { ScrollArea } from '@/shared/components/ui/scroll-area'
import { DialogFooter } from '@/shared/components/ui/dialog'
import { formatRupiah } from '@/shared/utils'
import { toast } from 'sonner'

import { useStoreProfileQuery } from '@/features/settings/store'
import { usePrinterSettingsQuery } from '@/features/settings/printer'
import type { PrinterSettings } from '@/features/settings/printer'

import {
  getRememberedPrinter,
  printViaBleBytes,
  isBlePrintSupported,
  pickAndRememberPrinter,
  buildReceiptEscposBytes,
  buildReceiptHtml,
  isWebSerialSupported,
  getRememberedSerialPort,
  pickSerialPort,
  printViaWebSerial,
  printViaRawBt,
  type ReceiptData,
} from '../blePrinter'
import { printReceiptHtml } from '../printFrame'
import { saveReceiptPdf } from '../receiptPdf'
import { printViaServer, printViaLan } from '../printApi'
import { shareReceiptViaWhatsApp } from '../receiptWhatsApp'
import { isNativeApp, printViaNativeBt } from '../printNativeBt'
import type {
  CartItem,
  CartSummary,
  CheckoutResponse,
  Discount,
  PaymentMethod,
  Tax,
} from '../cashier.types'
import { cartItemDetails, cartItemName, cartItemQtyLine, orderInfoRows } from '../cashier.utils'

interface ReceiptPrintProps {
  open: boolean
  onClose: () => void
  checkoutData: CheckoutResponse
  cart: CartItem[]
  summary: CartSummary
  discount: Discount
  tax: Tax
  paymentMethod: PaymentMethod
  amountPaid: number
  customerName?: string
  /** bill = tagihan sementara pesanan terbuka (keputusan 108): tanpa baris pembayaran. */
  mode?: 'checkout' | 'reprint' | 'bill'
}

// Pilihan cetak yang ditampilkan ke user. Tiap pilihan melakukan SATU jalur
// persis sesuai namanya — tidak ada fallback otomatis ke jalur lain, dan TIDAK
// ada metode yang dipilih otomatis oleh sistem. User selalu memilih sendiri tiap
// kali mencetak. Bila jalur tidak didukung perangkat, user diberi pesan jelas.
//   server    : BE mengirim ESC/POS ke COM port di komputer server (paling andal)
//   lan       : BE mengirim ESC/POS ke printer jaringan (IP) via TCP
//   serial    : Web Serial langsung ke COM port (butuh HTTPS/localhost, desktop)
//   bluetooth : Web Bluetooth langsung ke printer (butuh HTTPS/localhost + BLE)
//   native-bt : cetak Bluetooth langsung lewat APK (Capacitor) — hanya muncul di APK
//   rawbt     : lewat aplikasi RawBT di Android (intent) — untuk printer Bluetooth
//   pdf       : buat file PDF struk (selalu jalan, termasuk HTTP)
//   whatsapp  : bagikan ringkasan struk via WhatsApp (bukan cetak kertas)
//   browser   : dialog print bawaan browser (selalu jalan, termasuk HTTP)
type PrintChoice = 'native-bt' | 'server' | 'lan' | 'serial' | 'bluetooth' | 'rawbt' | 'pdf' | 'whatsapp' | 'browser'

// Daftar opsi cetak yang ditampilkan ke user, apa adanya. User memilih sendiri
// jalur mana yang dipakai setiap kali mencetak — tidak ada yang otomatis.
// `devices` menandai perangkat yang cocok, sebagai petunjuk visual untuk kasir.
// `nativeOnly` = hanya ditampilkan bila berjalan di dalam APK (Capacitor).
const PRINT_OPTIONS: Array<{
  key: PrintChoice
  title: string
  desc: string
  devices: string
  icon: typeof Bluetooth
  nativeOnly?: boolean
}> = [
  {
    key: 'native-bt',
    title: 'Printer Bluetooth (Aplikasi)',
    desc: 'Cetak langsung ke printer Bluetooth (mis. EP58M) lewat aplikasi. Hasil tajam & cepat. Hanya tersedia bila membuka lewat aplikasi (APK), bukan browser.',
    devices: 'Tablet / HP (Aplikasi)',
    icon: Bluetooth,
    nativeOnly: true,
  },
  {
    key: 'server',
    title: 'Cetak via Server (Printer di Komputer)',
    desc: 'Kirim ke printer thermal yang terhubung di komputer server (COM/Bluetooth). Paling stabil, bisa dipakai laptop maupun tablet selama satu jaringan dengan server.',
    devices: 'Desktop & Tablet',
    icon: Server,
  },
  {
    key: 'lan',
    title: 'Printer Jaringan / WiFi (LAN)',
    desc: 'Kirim ke printer yang punya alamat IP di jaringan (port 9100). Atur IP printer di pengaturan server. Bisa dipakai banyak perangkat sekaligus.',
    devices: 'Desktop & Tablet',
    icon: Wifi,
  },
  {
    key: 'serial',
    title: 'Printer USB / Kabel (Komputer ini)',
    desc: 'Cetak langsung ke printer yang tercolok di komputer ini lewat port USB/COM. Hanya untuk desktop/laptop (Chrome/Edge).',
    devices: 'Desktop',
    icon: Usb,
  },
  {
    key: 'bluetooth',
    title: 'Printer Bluetooth (Langsung)',
    desc: 'Cetak langsung ke printer thermal Bluetooth dari perangkat ini. Hanya untuk printer Bluetooth LE; printer Bluetooth lama mungkin tidak terdeteksi.',
    devices: 'Desktop & Tablet (Chrome)',
    icon: Bluetooth,
  },
  {
    key: 'rawbt',
    title: 'Printer via RawBT (Android)',
    desc: 'Cetak ke printer Bluetooth lewat aplikasi RawBT. Perlu aplikasi RawBT terpasang di perangkat Android dan printer sudah dipilih di dalamnya.',
    devices: 'Tablet / HP Android',
    icon: Smartphone,
  },
  {
    key: 'pdf',
    title: 'Simpan / Bagikan PDF',
    desc: 'Buat file PDF struk. Bisa diarsip, dikirim ke pelanggan, atau dicetak dari aplikasi perangkat. Jalan di semua perangkat.',
    devices: 'Desktop & Tablet',
    icon: FileText,
  },
  {
    key: 'whatsapp',
    title: 'Kirim via WhatsApp',
    desc: 'Bagikan ringkasan struk ke pelanggan lewat WhatsApp (bukan cetak kertas). WhatsApp akan terbuka dengan teks struk siap kirim.',
    devices: 'Desktop & Tablet',
    icon: MessageCircle,
  },
  {
    key: 'browser',
    title: 'Browser / Dialog Print',
    desc: 'Buka dialog print bawaan browser. Cadangan yang jalan di mana saja.',
    devices: 'Desktop & Tablet',
    icon: Monitor,
  },
]

const DEFAULT_PRINTER_SETTINGS: PrinterSettings = {
  paper_size: '80mm',
  receipt_header: '',
  receipt_footer: 'Terima kasih telah berbelanja',
  show_logo: false,
  auto_print: false,
  print_method: 'browser',
}

const PAYMENT_LABELS: Record<PaymentMethod, string> = {
  cash: 'Tunai',
  transfer: 'Transfer',
  qris: 'QRIS',
  card: 'Kartu',
  kredit: 'Hutang',
  balance: 'Saldo',
}

function formatDate(dateStr: string): string {
  const d = new Date(dateStr)
  return d.toLocaleString('id-ID', {
    day: '2-digit',
    month: 'short',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}

export function ReceiptPrint({
  open,
  onClose,
  checkoutData,
  cart,
  summary,
  discount,
  tax,
  paymentMethod,
  amountPaid,
  customerName,
  mode = 'checkout',
}: ReceiptPrintProps) {
  const change = amountPaid - summary.grandTotal
  const { data: storeProfile } = useStoreProfileQuery()
  const { data: printerSettingsData } = usePrinterSettingsQuery()
  const printerSettings = printerSettingsData ?? DEFAULT_PRINTER_SETTINGS

  // Pilihan metode cetak yang muncul saat tombol Cetak ditekan.
  const [methodPickerOpen, setMethodPickerOpen] = useState(false)

  const buildReceiptData = useCallback(
    (): ReceiptData => ({
      storeName: printerSettings.receipt_header || storeProfile?.name || 'POS Coffee',
      storeSub: [storeProfile?.address, storeProfile?.phone].filter(Boolean).join(' • '),
      footer: printerSettings.receipt_footer,
      paperSize: printerSettings.paper_size,
      checkoutData, cart, summary, discount, tax, paymentMethod, amountPaid, customerName,
      unpaid: mode === 'bill',
    }),
    [printerSettings, storeProfile, checkoutData, cart, summary, discount, tax, paymentMethod, amountPaid, customerName, mode],
  )

  // Jalur NATIVE BLUETOOTH (APK): cetak langsung ke printer Bluetooth lewat
  // plugin native Capacitor. Hanya aktif di dalam APK (kartunya pun hanya muncul
  // di APK). Tanpa fallback otomatis.
  const printNative = useCallback(async () => {
    const data = buildReceiptData()
    try {
      await printViaNativeBt(data)
      toast.success('Struk dikirim ke printer')
    } catch (err) {
      toast.error((err as Error)?.message || 'Gagal cetak via Bluetooth aplikasi.')
    }
  }, [buildReceiptData])

  // Jalur SERVER: BE yang mengirim ESC/POS ke printer (COM port) di komputer
  // server. Paling andal & tajam; dipakai laptop maupun tablet selama satu
  // jaringan dengan server. Tanpa fallback otomatis.
  const printServer = useCallback(async () => {
    const data = buildReceiptData()
    try {
      await printViaServer(data)
      toast.success('Struk dikirim ke printer')
    } catch (err) {
      toast.error((err as Error)?.message || 'Gagal cetak via server. Pastikan printer di komputer server menyala.')
    }
  }, [buildReceiptData])

  // Jalur LAN: BE meneruskan ESC/POS ke printer jaringan (IP) via TCP. Host/port
  // diambil dari konfigurasi server. Tanpa fallback otomatis.
  const printLan = useCallback(async () => {
    const data = buildReceiptData()
    try {
      await printViaLan(data)
      toast.success('Struk dikirim ke printer jaringan')
    } catch (err) {
      toast.error((err as Error)?.message || 'Gagal cetak ke printer jaringan. Pastikan IP printer benar & printer menyala.')
    }
  }, [buildReceiptData])

  // Jalur WHATSAPP: bagikan ringkasan struk ke pelanggan (bukan cetak kertas).
  const shareWhatsApp = useCallback(() => {
    const data = buildReceiptData()
    try {
      shareReceiptViaWhatsApp(data)
    } catch (err) {
      toast.error((err as Error)?.message || 'Gagal membuka WhatsApp')
    }
  }, [buildReceiptData])

  // Jalur RAWBT: cetak lewat aplikasi RawBT di Android (intent). Hanya tersedia
  // di Android; bila bukan Android, user diberi pesan. Tanpa fallback.
  const printRawBt = useCallback(() => {
    const data = buildReceiptData()
    try {
      printViaRawBt(buildReceiptEscposBytes(data))
    } catch (err) {
      toast.error((err as Error)?.message || 'Gagal cetak via RawBT')
    }
  }, [buildReceiptData])

  // Jalur BLUETOOTH (Web Bluetooth): kirim ESC/POS langsung ke printer. Hanya
  // tersedia di konteks aman (HTTPS/localhost) & untuk printer Bluetooth LE.
  // Tidak ada fallback — bila gagal, user diberi pesan dan bisa pilih metode
  // lain sendiri.
  const printBluetooth = useCallback(async () => {
    const data = buildReceiptData()
    if (!isBlePrintSupported()) {
      toast.error('Cetak Bluetooth tidak tersedia di perangkat/alamat ini. Pilih metode lain (mis. Simpan PDF).')
      return
    }
    try {
      const ble = (await getRememberedPrinter()) ?? (await pickAndRememberPrinter())
      await printViaBleBytes(ble, buildReceiptEscposBytes(data))
      toast.success('Struk dikirim ke printer')
    } catch (err) {
      toast.error((err as Error)?.message || 'Gagal cetak via Bluetooth')
    }
  }, [buildReceiptData])

  // Jalur SERIAL (Web Serial): kirim ESC/POS langsung ke COM port. Hanya
  // tersedia di konteks aman (HTTPS/localhost), umumnya desktop. Tanpa fallback.
  const printSerial = useCallback(async () => {
    const data = buildReceiptData()
    if (!isWebSerialSupported()) {
      toast.error('Cetak USB/Serial tidak tersedia di perangkat/alamat ini. Pilih metode lain (mis. Simpan PDF).')
      return
    }
    try {
      const port = (await getRememberedSerialPort()) ?? (await pickSerialPort())
      await printViaWebSerial(port, buildReceiptEscposBytes(data))
      toast.success('Struk dikirim ke printer')
    } catch (err) {
      toast.error((err as Error)?.message || 'Gagal cetak via USB/Serial')
    }
  }, [buildReceiptData])

  // Jalur BROWSER: cetak lewat iframe tersembunyi (deterministik, tidak
  // terpengaruh layout aplikasi — menghindari halaman kosong/dobel).
  const printBrowser = useCallback(() => {
    const data = buildReceiptData()
    printReceiptHtml(buildReceiptHtml(data), data.paperSize)
  }, [buildReceiptData])

  // Jalur PDF: buat file PDF struk di browser (tanpa server/Bluetooth).
  // Jalan di HTTP maupun HTTPS, di tablet maupun desktop. File bisa diarsip,
  // dibagikan ke pelanggan, atau dibuka lalu dicetak ke printer thermal lewat
  // aplikasi bawaan perangkat.
  const printPdf = useCallback(() => {
    const data = buildReceiptData()
    try {
      saveReceiptPdf(data)
      toast.success('Struk PDF berhasil dibuat')
    } catch (err) {
      toast.error((err as Error)?.message || 'Gagal membuat PDF struk')
    }
  }, [buildReceiptData])

  const runPrintMethod = useCallback(
    (method: PrintChoice) => {
      switch (method) {
        case 'native-bt':
          void printNative()
          break
        case 'server':
          void printServer()
          break
        case 'lan':
          void printLan()
          break
        case 'serial':
          void printSerial()
          break
        case 'bluetooth':
          void printBluetooth()
          break
        case 'rawbt':
          printRawBt()
          break
        case 'pdf':
          printPdf()
          break
        case 'whatsapp':
          shareWhatsApp()
          break
        case 'browser':
          printBrowser()
          break
      }
    },
    [printNative, printServer, printLan, printSerial, printBluetooth, printRawBt, printPdf, shareWhatsApp, printBrowser],
  )

  // Tombol Cetak: buka pemilih metode agar user tetap bisa memilih tiap cetak.
  const handlePrint = useCallback(() => {
    setMethodPickerOpen(true)
  }, [])

  const handlePickMethod = useCallback(
    (method: PrintChoice) => {
      setMethodPickerOpen(false)
      runPrintMethod(method)
    },
    [runPrintMethod],
  )

  // Catatan: TIDAK ada auto-print. Sesuai permintaan, setiap cetak selalu lewat
  // modal pemilih metode (handlePrint membuka pemilih), tidak pernah otomatis.

  const footer = (
    <div className="no-print">
      <DialogFooter className="border-t px-6 py-4">
        <Button variant="outline" onClick={handlePrint} className="gap-1.5">
          <Printer size={14} />
          Cetak
        </Button>
        {mode === 'checkout' ? (
          <Button onClick={onClose} className="gap-1.5">
            <ShoppingCart size={14} />
            Transaksi Baru
          </Button>
        ) : (
          <Button onClick={onClose} className="gap-1.5">
            <X size={14} />
            Tutup
          </Button>
        )}
      </DialogFooter>
    </div>
  )

  // Isi struk (tanpa wrapper). Dipakai dua kali: sebagai preview di dalam modal
  // (tampil di layar) dan sebagai konten cetak di dalam PrintPortal (tercetak).
  const receiptInner = (
    <>
      {/* Header kedai */}
      <div className="text-center space-y-0.5">
        {printerSettings.show_logo && storeProfile?.logo_url && (
          <img src={storeProfile.logo_url} alt="Logo" className="mx-auto mb-1 max-h-16 max-w-20 object-contain" />
        )}
        <p className="text-base font-bold text-gray-900 tracking-wide">
          {printerSettings.receipt_header || storeProfile?.name || 'POS Coffee'}
        </p>
        {[storeProfile?.address, storeProfile?.phone].filter(Boolean).length > 0 && (
          <p className="text-sm text-gray-500">
            {[storeProfile?.address, storeProfile?.phone].filter(Boolean).join(' • ')}
          </p>
        )}
      </div>

      <hr className="border-dashed border-gray-300" />

      {/* Nomor antrean besar (keputusan 56, 95) */}
      {checkoutData.queue_label && (
        <>
          <div className="text-center">
            <p className="text-xs uppercase tracking-wide text-gray-500">Antrean</p>
            <p className="text-4xl font-bold text-gray-900" aria-label="Nomor antrean">{checkoutData.queue_label}</p>
          </div>
          <hr className="border-dashed border-gray-300" />
        </>
      )}

      {/* Info transaksi */}
      <div className="space-y-1 text-xs text-gray-600">
        <div className="flex justify-between">
          <span className="text-gray-400">No. Transaksi</span>
          <span className="font-medium text-gray-800">{checkoutData.transaction_code}</span>
        </div>
        <div className="flex justify-between">
          <span className="text-gray-400">Tanggal</span>
          <span className="text-gray-700">{formatDate(checkoutData.transaction_date)}</span>
        </div>
        {customerName && (
          <div className="flex justify-between">
            <span className="text-gray-400">Pelanggan</span>
            <span className="text-gray-700">{customerName}</span>
          </div>
        )}
        {orderInfoRows(checkoutData).map(([label, value]) => (
          <div key={label} className="flex justify-between">
            <span className="text-gray-400">{label}</span>
            <span className="text-gray-700">{value}</span>
          </div>
        ))}
        {mode === 'bill' ? (
          <p className="pt-1 text-center text-sm font-bold text-red-600">TAGIHAN - BELUM DIBAYAR</p>
        ) : (
          <div className="flex justify-between">
            <span className="text-gray-400">Pembayaran</span>
            <span className="text-gray-700">{PAYMENT_LABELS[paymentMethod]}</span>
          </div>
        )}
      </div>

      <hr className="border-dashed border-gray-300" />

      {/* Item list */}
      <div className="space-y-2.5">
        {cart.map((item) => (
          <div key={item.key} className="space-y-0.5">
            <div className="flex justify-between items-start gap-2">
              <div className="flex-1 min-w-0">
                <p className="text-sm font-medium text-gray-800">{cartItemName(item)}</p>
                {cartItemDetails(item).map((d) => (
                  <p key={d} className="text-xs text-gray-500">{d}</p>
                ))}
                <p className="text-xs text-gray-400">{cartItemQtyLine(item)}</p>
              </div>
              <span className="text-sm font-semibold text-gray-800 shrink-0">
                {formatRupiah(item.subtotal)}
              </span>
            </div>
            {item.discount_amount && item.discount_amount > 0 && (
              <div className="flex justify-between text-xs text-red-500 pl-2">
                <span>
                  Diskon{' '}
                  {item.discount_type === 'percent'
                    ? `${item.discount_value}%`
                    : formatRupiah(item.discount_value ?? 0)}
                </span>
                <span>-{formatRupiah(item.discount_amount)}</span>
              </div>
            )}
          </div>
        ))}
      </div>

      <hr className="border-dashed border-gray-300" />

      {/* Summary */}
      <div className="space-y-1.5 text-sm">
        <div className="flex justify-between text-gray-600">
          <span>Subtotal</span>
          <span>{formatRupiah(summary.subtotal)}</span>
        </div>
        {summary.discountAmount > 0 && (
          <div className="flex justify-between text-green-600">
            <span>Diskon{discount.type === 'percent' ? ` (${discount.value}%)` : ''}</span>
            <span>-{formatRupiah(summary.discountAmount)}</span>
          </div>
        )}
        {summary.taxAmount > 0 && (
          <div className="flex justify-between text-gray-600">
            <span>Pajak ({tax.percent}%)</span>
            <span>+{formatRupiah(summary.taxAmount)}</span>
          </div>
        )}
        <hr className="border-gray-200" />
        <div className="flex justify-between font-bold text-base text-gray-900">
          <span>TOTAL</span>
          <span>{formatRupiah(summary.grandTotal)}</span>
        </div>
        {mode !== 'bill' && (
          <>
            <div className="flex justify-between text-gray-600">
              <span>Dibayar</span>
              <span>{formatRupiah(amountPaid)}</span>
            </div>
            <div className="flex justify-between font-semibold">
              <span>Kembalian</span>
              <span className="text-green-600">{formatRupiah(Math.max(0, change))}</span>
            </div>
          </>
        )}
      </div>

      <div className="text-center pt-1">
        <p className="text-xs text-gray-400">{printerSettings.receipt_footer || ''}</p>
      </div>
    </>
  )

  const receiptMaxWidth = printerSettings.paper_size === '58mm' ? '210px' : '300px'

  return (
    <>
    <ActionModal
      open={open}
      onOpenChange={(val) => { if (!val) onClose() }}
      title={mode === 'bill' ? 'Tagihan Sementara' : 'Struk Transaksi'}
      description="Detail struk transaksi"
      contentClassName="max-w-sm"
      footer={footer}
    >
        {/* Preview struk di layar (TIDAK ikut tercetak — cetak dilakukan dari
            PrintPortal di bawah). */}
        <ScrollArea style={{ maxHeight: '65vh' }}>
        <div
          className="no-print px-6 py-5 space-y-4"
          style={{ maxWidth: receiptMaxWidth, margin: '0 auto' }}
        >
          {receiptInner}
        </div>
        </ScrollArea>
    </ActionModal>



    {/* Pemilih metode cetak — muncul saat tombol Cetak ditekan. */}
    <ActionModal
      open={methodPickerOpen}
      onOpenChange={setMethodPickerOpen}
      title="Pilih Metode Cetak"
      description="Pilih cara mencetak struk"
      contentClassName="max-w-xs"
      hideFooter
    >
      <div className="px-5 py-4 space-y-3">
        {PRINT_OPTIONS.filter((opt) => !opt.nativeOnly || isNativeApp()).map((opt) => {
          const Icon = opt.icon
          return (
            <button
              key={opt.key}
              type="button"
              onClick={() => handlePickMethod(opt.key)}
              className="w-full rounded-lg border-2 border-gray-200 p-3 text-left transition-colors hover:bg-gray-50 hover:border-[#2c3e50]"
            >
              <div className="flex items-center gap-2">
                <Icon size={18} className="text-[#2c3e50] shrink-0" />
                <span className="text-sm font-semibold text-gray-800">{opt.title}</span>
                <span className="ml-auto text-[10px] font-medium text-gray-600 bg-gray-100 px-1.5 py-0.5 rounded shrink-0">
                  {opt.devices}
                </span>
              </div>
              <p className="text-xs text-gray-500 mt-1">{opt.desc}</p>
            </button>
          )
        })}
      </div>
    </ActionModal>
    </>
  )
}

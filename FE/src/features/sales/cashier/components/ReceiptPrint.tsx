import { useCallback, useEffect, useState } from 'react'
import { Printer, ShoppingCart, X, Wallet, Bluetooth, Monitor, FileText, Usb, Server, Smartphone } from 'lucide-react'
import { useMutation } from '@tanstack/react-query'

import { ActionModal } from '@/shared/components'
import { Button } from '@/shared/components/ui/button'
import { ScrollArea } from '@/shared/components/ui/scroll-area'
import { DialogFooter } from '@/shared/components/ui/dialog'
import { formatRupiah } from '@/shared/utils'
import { RupiahInput } from '@/shared/components/ui/rupiah-input'
import { api } from '@/services'
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
import { printViaServer } from '../printApi'
import type {
  CartItem,
  CartSummary,
  CheckoutResponse,
  Discount,
  PaymentMethod,
  Tax,
} from '../cashier.types'

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
  mode?: 'checkout' | 'reprint'
}

// Pilihan cetak yang ditampilkan ke user. Tiap pilihan melakukan SATU jalur
// persis sesuai namanya — tidak ada fallback otomatis ke jalur lain, dan TIDAK
// ada metode yang dipilih otomatis oleh sistem. User selalu memilih sendiri tiap
// kali mencetak. Bila jalur tidak didukung perangkat, user diberi pesan jelas.
//   server    : BE mengirim ESC/POS ke COM port di komputer server (paling andal)
//   serial    : Web Serial langsung ke COM port (butuh HTTPS/localhost, desktop)
//   bluetooth : Web Bluetooth langsung ke printer (butuh HTTPS/localhost + BLE)
//   rawbt     : lewat aplikasi RawBT di Android (intent) — untuk printer Bluetooth
//   pdf       : buat file PDF struk (selalu jalan, termasuk HTTP)
//   browser   : dialog print bawaan browser (selalu jalan, termasuk HTTP)
type PrintChoice = 'server' | 'serial' | 'bluetooth' | 'rawbt' | 'pdf' | 'browser'

// Daftar opsi cetak yang ditampilkan ke user, apa adanya. User memilih sendiri
// jalur mana yang dipakai setiap kali mencetak — tidak ada yang otomatis.
// `devices` menandai perangkat yang cocok, sebagai petunjuk visual untuk kasir.
const PRINT_OPTIONS: Array<{
  key: PrintChoice
  title: string
  desc: string
  devices: string
  icon: typeof Bluetooth
}> = [
  {
    key: 'server',
    title: 'Cetak via Server (Printer di Komputer)',
    desc: 'Kirim ke printer thermal yang terhubung di komputer server (COM/Bluetooth). Paling stabil, bisa dipakai laptop maupun tablet selama satu jaringan dengan server.',
    devices: 'Desktop & Tablet',
    icon: Server,
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
      storeName: printerSettings.receipt_header || storeProfile?.name || 'POS System',
      storeSub: [storeProfile?.address, storeProfile?.phone].filter(Boolean).join(' • '),
      footer: printerSettings.receipt_footer,
      paperSize: printerSettings.paper_size,
      checkoutData, cart, summary, discount, tax, paymentMethod, amountPaid, customerName,
    }),
    [printerSettings, storeProfile, checkoutData, cart, summary, discount, tax, paymentMethod, amountPaid, customerName],
  )

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
        case 'server':
          void printServer()
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
        case 'browser':
          printBrowser()
          break
      }
    },
    [printServer, printSerial, printBluetooth, printRawBt, printPdf, printBrowser],
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

  // Save-to-balance feature
  const showSaveToBalance = mode === 'checkout' && !!customerName && paymentMethod === 'cash' && change > 0
  const [saveAmount, setSaveAmount] = useState(0)
  const [savedToBalance, setSavedToBalance] = useState(false)

  // Inisialisasi nilai simpan-ke-saldo saat struk dibuka. set-state di effect
  // memang disengaja: menyiapkan default form dari nilai kembalian.
  useEffect(() => {
    if (open && showSaveToBalance) {
      // eslint-disable-next-line react-hooks/set-state-in-effect
      setSaveAmount(Math.max(0, change))
      setSavedToBalance(false)
    }
  }, [open, showSaveToBalance, change])

  const saveToBalanceMutation = useMutation({
    mutationFn: (payload: { transaction_id: number; amount: number }) =>
      api.post('/transactions/save-to-balance', payload),
    onSuccess: () => {
      setSavedToBalance(true)
      toast.success(`Rp ${saveAmount.toLocaleString('id-ID')} berhasil disimpan ke saldo pelanggan`)
    },
    onError: (err: Error) => {
      toast.error(err.message || 'Gagal menyimpan ke saldo')
    },
  })

  function handleSaveToBalance() {
    if (saveAmount <= 0 || saveAmount > change) return
    saveToBalanceMutation.mutate({
      transaction_id: checkoutData.id,
      amount: saveAmount,
    })
  }

  const footer = (
    <div className="no-print">
      {/* Simpan kembalian ke saldo — only checkout mode, cash, has customer, has change */}
      {showSaveToBalance && !savedToBalance && (
        <div className="border-t px-6 py-3 bg-blue-50 space-y-2">
          <div className="flex items-center gap-2">
            <Wallet size={16} className="text-blue-600" />
            <span className="text-sm font-semibold text-blue-800">Simpan kembalian ke Saldo Pelanggan?</span>
          </div>
          <div className="flex items-center gap-2">
            <RupiahInput
              value={saveAmount}
              onChange={(v) => setSaveAmount(Math.min(v, change))}
              className="flex-1 h-9 text-sm"
            />
            <Button
              size="sm"
              variant="outline"
              onClick={onClose}
              className="shrink-0"
            >
              Tidak
            </Button>
            <Button
              size="sm"
              onClick={handleSaveToBalance}
              disabled={saveAmount <= 0 || saveToBalanceMutation.isPending}
              className="shrink-0 bg-blue-600 hover:bg-blue-700"
            >
              {saveToBalanceMutation.isPending ? 'Menyimpan...' : 'Ya, Simpan'}
            </Button>
          </div>
        </div>
      )}
      {savedToBalance && (
        <div className="border-t px-6 py-3 bg-green-50">
          <p className="text-sm text-green-700 font-medium">
            ✓ {formatRupiah(saveAmount)} berhasil disimpan ke saldo pelanggan
          </p>
        </div>
      )}
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
      {/* Toko header */}
      <div className="text-center space-y-0.5">
        {printerSettings.show_logo && storeProfile?.logo_url && (
          <img src={storeProfile.logo_url} alt="Logo" className="mx-auto mb-1 max-h-16 max-w-20 object-contain" />
        )}
        <p className="text-base font-bold text-gray-900 tracking-wide">
          {printerSettings.receipt_header || storeProfile?.name || 'POS System'}
        </p>
        {[storeProfile?.address, storeProfile?.phone].filter(Boolean).length > 0 && (
          <p className="text-sm text-gray-500">
            {[storeProfile?.address, storeProfile?.phone].filter(Boolean).join(' • ')}
          </p>
        )}
      </div>

      <hr className="border-dashed border-gray-300" />

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
        <div className="flex justify-between">
          <span className="text-gray-400">Pembayaran</span>
          <span className="text-gray-700">{PAYMENT_LABELS[paymentMethod]}</span>
        </div>
      </div>

      <hr className="border-dashed border-gray-300" />

      {/* Item list */}
      <div className="space-y-2.5">
        {cart.map((item) => (
          <div key={`${item.product_id}-${item.unit_id}`} className="space-y-0.5">
            <div className="flex justify-between items-start gap-2">
              <div className="flex-1 min-w-0">
                <p className="text-sm font-medium text-gray-800 truncate">{item.product_name}</p>
                <p className="text-xs text-gray-400">
                  {item.unit_name} &times; {item.qty} &nbsp;@&nbsp;
                  {formatRupiah(item.effective_price ?? item.price)}
                </p>
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
        <div className="flex justify-between text-gray-600">
          <span>Dibayar</span>
          <span>{formatRupiah(amountPaid)}</span>
        </div>
        <div className="flex justify-between font-semibold">
          <span>Kembalian</span>
          <span className="text-green-600">{formatRupiah(Math.max(0, change))}</span>
        </div>
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
      title="Struk Transaksi"
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
        {PRINT_OPTIONS.map((opt) => {
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

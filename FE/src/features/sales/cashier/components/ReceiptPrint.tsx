import { useCallback, useEffect, useRef, useState } from 'react'
import { Printer, ShoppingCart, X, Wallet } from 'lucide-react'
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

import { getRememberedPrinter, printViaBle, receiptToPlainText } from '../blePrinter'
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

const DEFAULT_PRINTER_SETTINGS: PrinterSettings = {
  paper_size: '80mm',
  receipt_header: '',
  receipt_footer: 'Terima kasih telah berbelanja',
  show_logo: false,
  auto_print: false,
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
  const autoPrinted = useRef(false)

  const handlePrint = useCallback(async () => {
    const bleDevice = await getRememberedPrinter()
    if (bleDevice) {
      const text = receiptToPlainText({
        storeName: printerSettings.receipt_header || storeProfile?.name || 'POS System',
        storeSub: [storeProfile?.address, storeProfile?.phone].filter(Boolean).join(' • '),
        footer: printerSettings.receipt_footer,
        paperSize: printerSettings.paper_size,
        checkoutData, cart, summary, discount, tax, paymentMethod, amountPaid, customerName,
      })
      try {
        await printViaBle(bleDevice, text)
        return
      } catch {
        // Printer BLE gagal (di luar jangkauan, dst) — fallback ke dialog print biasa.
      }
    }
    window.print()
  }, [printerSettings, storeProfile, checkoutData, cart, summary, discount, tax, paymentMethod, amountPaid, customerName])

  useEffect(() => {
    if (open && mode === 'checkout' && printerSettings.auto_print && !autoPrinted.current) {
      autoPrinted.current = true
      handlePrint()
    }
    if (!open) autoPrinted.current = false
  }, [open, mode, printerSettings.auto_print, handlePrint])

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

  return (
    <ActionModal
      open={open}
      onOpenChange={(val) => { if (!val) onClose() }}
      title="Struk Transaksi"
      description="Detail struk transaksi"
      contentClassName="max-w-sm"
      footer={footer}
    >
        {/* Receipt preview */}
        <ScrollArea style={{ maxHeight: '65vh' }}>
        <div
          className="print-root px-6 py-5 space-y-4"
          style={{ maxWidth: printerSettings.paper_size === '58mm' ? '210px' : '300px', margin: '0 auto' }}
        >

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
        </div>
        </ScrollArea>
    </ActionModal>
  )
}

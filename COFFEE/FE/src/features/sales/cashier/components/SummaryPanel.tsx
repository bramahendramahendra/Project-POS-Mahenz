import { useMemo, useState } from 'react'
import { AlertTriangle, ClipboardList, Receipt, Trash2 } from 'lucide-react'

import { ConfirmDialog } from '@/shared/components'
import { Button } from '@/shared/components/ui/button'
import { ScrollArea } from '@/shared/components/ui/scroll-area'
import { formatRupiah } from '@/shared/utils'
import { useCashDrawerCurrentQuery } from '@/features/finance/cash-drawer'

import { useCashierStore } from '../cashier.store'
import { calcSubtotal } from '../cashier.utils'
import { useAddBillItemsMutation, useOpenBillsQuery, useSaveOpenBillMutation } from '../openBill.api'
import { makeBillStore, useActiveBill } from '../useActiveBill'
import { ActiveBillPanel } from './ActiveBillPanel'
import { CartItemRow } from './CartItemRow'
import { CustomerSelector } from './CustomerSelector'
import { DiscountInput } from './DiscountInput'
import { OpenBillsDialog } from './OpenBillsDialog'
import { OrderInfoPanel } from './OrderInfoPanel'
import { TaxInput } from './TaxInput'

/**
 * Panel pesanan kasir. Pesanan baru: info pesanan, keranjang, Simpan (bayar nanti) atau Bayar.
 * Pesanan terbuka (tahap 4.8): item tersimpan, item baru ditambahkan ke pesanan, lalu Bayar/utang.
 */
export function SummaryPanel() {
  const store = useCashierStore()
  const { cart, discount, tax, orderType, table, orderName, openPaymentModal, clearCart, setActiveBill } = store
  const { activeBillId, bill, billOpen, billLines, summary } = useActiveBill()
  const [clearConfirmOpen, setClearConfirmOpen] = useState(false)
  const [openListOpen, setOpenListOpen] = useState(false)
  const { data: openBills = [] } = useOpenBillsQuery()
  const { mutate: saveBill, isPending: isSaving } = useSaveOpenBillMutation()
  const { mutate: addItems, isPending: isAdding } = useAddBillItemsMutation()
  const billStore = useMemo(() => makeBillStore(billLines), [billLines])

  const { data: currentDrawer } = useCashDrawerCurrentQuery()
  const kasOpen = currentDrawer?.status === 'open'
  // Tahap 4.9: kas tertunda (dari hari sebelumnya) tidak bisa dipakai berjualan baru; hanya menyelesaikan pesanan tertunda.
  const kasPending = kasOpen && !!currentDrawer?.is_pending
  const overdueBills = openBills.filter((b) => b.is_overdue)
  const billOverdue = billOpen && !!bill?.is_overdue
  const needsTable = !activeBillId && orderType === 'dine_in' && !table
  const cartTotal = calcSubtotal(cart)
  const showDiscountTax = billOpen ? cart.length === 0 && billLines.length > 0 : cart.length > 0

  const payloadItems = () =>
    cart.map((i) => ({
      variant_id: i.variant_id,
      quantity: i.qty,
      modifier_option_ids: i.modifiers.map((m) => m.option_id),
      notes: i.notes,
      discount_item: i.discount_amount ?? 0,
    }))

  const handleSave = () =>
    saveBill(
      {
        order_type: orderType,
        table_id: orderType === 'dine_in' ? table?.id : undefined,
        order_name: orderName.trim() || undefined,
        device_source: 'web',
        items: payloadItems(),
      },
      { onSuccess: () => clearCart() },
    )

  const handleAddToBill = () => {
    if (!activeBillId) return
    addItems(
      { id: activeBillId, items: payloadItems() },
      {
        onSuccess: () => {
          clearCart()
          setActiveBill(activeBillId)
        },
      },
    )
  }

  return (
    <div className="flex h-full flex-col border-l border-gray-200 bg-white">
      <div className="flex shrink-0 items-center gap-2 border-b px-4 py-3">
        <Receipt size={16} className="text-gray-500" />
        <span className="text-sm font-semibold text-gray-800">{billOpen ? 'Pesanan Terbuka' : 'Pesanan'}</span>
        <Button variant="outline" size="sm" className="ml-auto h-7 gap-1 text-xs" onClick={() => setOpenListOpen(true)} aria-label="Daftar pesanan terbuka">
          <ClipboardList size={12} />
          Terbuka ({openBills.length})
        </Button>
        {cart.length > 0 && (
          <Button variant="ghost" size="sm" className="h-7 gap-1 px-2 text-xs text-red-500 hover:bg-red-50 hover:text-red-600" onClick={() => setClearConfirmOpen(true)}>
            <Trash2 size={12} />
            Kosongkan
          </Button>
        )}
      </div>

      {activeBillId ? bill ? <ActiveBillPanel bill={bill} /> : <div className="h-24 shrink-0 animate-pulse border-b bg-gray-50" /> : <OrderInfoPanel />}
      <CustomerSelector />

      <ScrollArea className="min-h-0 flex-1">
        {billOpen && cart.length > 0 && <p className="px-4 pt-3 text-xs font-semibold uppercase tracking-wide text-blue-700">Item baru (belum disimpan)</p>}
        {cart.length === 0 ? (
          <div className="flex items-center justify-center py-8">
            <p className="text-center text-sm text-gray-400">{billOpen ? 'Ketuk menu untuk menambah item ke pesanan ini.' : 'Belum ada item. Ketuk menu untuk menambah.'}</p>
          </div>
        ) : (
          <ul className="divide-y divide-gray-100">
            {cart.map((item) => (
              <CartItemRow key={item.key} item={item} />
            ))}
          </ul>
        )}
        {/* Diskon dan pajak ikut tergulir; pada pesanan terbuka diisi saat bayar (keputusan 104). */}
        {showDiscountTax && (
          <div className="border-t">
            {billOpen ? (
              <>
                <DiscountInput store={billStore} />
                <TaxInput store={billStore} />
              </>
            ) : (
              <>
                <DiscountInput />
                <TaxInput />
              </>
            )}
          </div>
        )}
      </ScrollArea>

      <div className="shrink-0 space-y-1.5 border-t bg-gray-50 px-4 py-3 text-sm">
        {billOpen && cart.length > 0 ? (
          <div className="flex justify-between font-semibold text-blue-800">
            <span>Item baru</span>
            <span>{formatRupiah(cartTotal)}</span>
          </div>
        ) : (
          <>
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
            <div className="flex items-baseline justify-between border-t pt-2 font-bold text-gray-900">
              <span className="text-sm">TOTAL</span>
              <span className="text-xl" aria-label="Total pesanan">
                {formatRupiah(summary.grandTotal)}
              </span>
            </div>
          </>
        )}
      </div>

      {!kasOpen && (
        <div className="mx-4 mb-2 flex shrink-0 items-start gap-2 rounded-md border border-amber-200 bg-amber-50 px-3 py-2 text-xs text-amber-700">
          <AlertTriangle size={13} className="mt-0.5 shrink-0" />
          <span>Kas belum dibuka. Buka kas terlebih dahulu untuk memproses pembayaran.</span>
        </div>
      )}
      {kasPending && (
        <div className="mx-4 mb-2 flex shrink-0 items-start gap-2 rounded-md border border-amber-300 bg-amber-50 px-3 py-2 text-xs text-amber-800" role="alert">
          <AlertTriangle size={13} className="mt-0.5 shrink-0" />
          <span>Kas Anda masih tertunda dari hari sebelumnya. Selesaikan pesanan tertunda, lalu tutup kas itu sebelum berjualan lagi.</span>
        </div>
      )}
      {overdueBills.length > 0 && !billOpen && (
        <button type="button" onClick={() => setOpenListOpen(true)} className="mx-4 mb-2 shrink-0 rounded-md border border-amber-300 bg-amber-50 px-3 py-2 text-left text-xs text-amber-800 hover:bg-amber-100">
          Ada {overdueBills.length} pesanan tertunda dari hari sebelumnya. Ketuk untuk menyelesaikan.
        </button>
      )}
      {billOverdue && cart.length > 0 && (
        <p className="mx-4 mb-2 shrink-0 text-xs text-amber-700">Pesanan tertunda tidak bisa ditambah item. Kosongkan item baru.</p>
      )}
      {needsTable && cart.length > 0 && <p className="mx-4 mb-2 shrink-0 text-xs text-amber-700">Pilih meja untuk pesanan makan di tempat.</p>}

      <div className="flex shrink-0 gap-2 border-t px-4 py-3">
        {billOpen ? (
          cart.length > 0 ? (
            <Button className="h-11 w-full text-base" onClick={handleAddToBill} disabled={isAdding || billOverdue || kasPending}>
              Tambahkan ke pesanan · {formatRupiah(cartTotal)}
            </Button>
          ) : (
            <Button className="h-11 w-full gap-2 text-base" onClick={openPaymentModal} disabled={!kasOpen || billLines.length === 0 || (kasPending && !billOverdue)}>
              Bayar {formatRupiah(summary.grandTotal)}
            </Button>
          )
        ) : (
          <>
            <Button variant="outline" className="h-11 flex-1 text-sm" onClick={handleSave} disabled={cart.length === 0 || needsTable || isSaving || kasPending}>
              Simpan (bayar nanti)
            </Button>
            <Button className="h-11 flex-1 gap-2 text-base" onClick={openPaymentModal} disabled={cart.length === 0 || !kasOpen || needsTable || kasPending}>
              Bayar {cart.length > 0 && formatRupiah(summary.grandTotal)}
            </Button>
          </>
        )}
      </div>

      {openListOpen && (
        <OpenBillsDialog
          title="Pesanan Terbuka"
          description="Pesanan yang belum dibayar dari semua kasir."
          actionLabel="Buka"
          excludeId={activeBillId}
          onClose={() => setOpenListOpen(false)}
          onPick={(b) => {
            setActiveBill(b.id)
            setOpenListOpen(false)
          }}
        />
      )}

      <ConfirmDialog
        open={clearConfirmOpen}
        onOpenChange={setClearConfirmOpen}
        title={billOpen ? 'Kosongkan Item Baru' : 'Kosongkan Pesanan'}
        description={billOpen ? 'Item baru yang belum disimpan akan dihapus. Pesanan terbuka tidak berubah.' : 'Semua item, meja, dan nama pemesan akan dihapus. Yakin?'}
        confirmLabel="Ya, Kosongkan"
        variant="destructive"
        onConfirm={() => {
          clearCart()
          if (billOpen) setActiveBill(activeBillId)
          setClearConfirmOpen(false)
        }}
      />
    </div>
  )
}

import { useEffect } from 'react'
import { toast } from 'sonner'

import { useCashDrawerCurrentQuery } from '@/features/finance/cash-drawer'

import { useCashierStore } from './cashier.store'
import { CartEditableList } from './components/CartEditableList'
import { PaymentModal } from './components/PaymentModal'
import { ProductSearch } from './components/ProductSearch'
import { SummaryPanel } from './components/SummaryPanel'

export function CashierPage() {
  const { data: currentDrawer, isLoading: isLoadingDrawer } = useCashDrawerCurrentQuery()
  const { paymentModalOpen, closePaymentModal } = useCashierStore()

  useEffect(() => {
    if (!isLoadingDrawer && !currentDrawer) {
      toast.warning('Belum ada kas yang dibuka. Buka kas/shift terlebih dahulu.', {
        duration: 5000,
        id: 'no-active-shift',
      })
    }
  }, [currentDrawer, isLoadingDrawer])

  return (
    <div
      style={{
        display: 'flex',
        height: 'calc(100vh - var(--navbar-height))',
        overflow: 'hidden',
      }}
    >
      {/* Panel Kiri — Search + Keranjang Editable */}
      <div
        style={{
          flex: 1,
          display: 'flex',
          flexDirection: 'column',
          borderRight: '1px solid var(--color-border)',
          overflow: 'hidden',
        }}
        className="bg-gray-50"
      >
        {/* Search */}
        <div className="px-4 pt-4 pb-3 shrink-0">
          <ProductSearch />
        </div>

        {/* Keranjang editable — flex-1, scrollable */}
        <CartEditableList />
      </div>

      {/* Panel Kanan — Ringkasan (360px fixed) */}
      <div style={{ width: '360px', flexShrink: 0, overflow: 'hidden', display: 'flex', flexDirection: 'column' }}>
        <SummaryPanel />
      </div>

      <PaymentModal
        open={paymentModalOpen}
        onOpenChange={(open) => {
          if (!open) closePaymentModal()
        }}
      />
    </div>
  )
}

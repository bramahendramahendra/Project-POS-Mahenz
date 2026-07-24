import { useEffect, useState } from 'react'
import { toast } from 'sonner'

import { useCashDrawerCurrentQuery } from '@/features/finance/cash-drawer'
import { useBreakpoint } from '@/shared/hooks'

import { useCashierStore } from './cashier.store'
import { CartEditableList } from './components/CartEditableList'
import { PaymentModal } from './components/PaymentModal'
import { ProductSearch } from './components/ProductSearch'
import { SummaryPanel } from './components/SummaryPanel'

export function CashierPage() {
  const { data: currentDrawer, isLoading: isLoadingDrawer } = useCashDrawerCurrentQuery()
  const { paymentModalOpen, closePaymentModal, cart } = useCashierStore()
  const isDesktop = useBreakpoint('lg')
  const [mobileTab, setMobileTab] = useState<'produk' | 'keranjang'>('produk')
  const itemCount = cart.reduce((sum, i) => sum + i.qty, 0)

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
        flexDirection: isDesktop ? 'row' : 'column',
        height: 'calc(100vh - var(--navbar-height))',
        overflow: 'hidden',
      }}
    >
      {/* Switch tab Produk/Keranjang — hanya tablet & mobile */}
      {!isDesktop && (
        <div className="flex shrink-0 border-b bg-white">
          <button
            type="button"
            onClick={() => setMobileTab('produk')}
            className={`flex-1 py-2.5 text-sm font-medium border-b-2 transition-colors ${
              mobileTab === 'produk'
                ? 'border-blue-600 text-blue-600'
                : 'border-transparent text-gray-500'
            }`}
          >
            Produk
          </button>
          <button
            type="button"
            onClick={() => setMobileTab('keranjang')}
            className={`flex-1 py-2.5 text-sm font-medium border-b-2 transition-colors relative ${
              mobileTab === 'keranjang'
                ? 'border-blue-600 text-blue-600'
                : 'border-transparent text-gray-500'
            }`}
          >
            Keranjang
            {itemCount > 0 && (
              <span className="ml-1.5 inline-flex items-center justify-center rounded-full bg-blue-600 text-white text-[10px] min-w-[16px] h-4 px-1">
                {itemCount}
              </span>
            )}
          </button>
        </div>
      )}

      {/* Panel Kiri — Search + Keranjang Editable */}
      <div
        style={{
          flex: 1,
          display: isDesktop || mobileTab === 'produk' ? 'flex' : 'none',
          flexDirection: 'column',
          borderRight: isDesktop ? '1px solid var(--color-border)' : undefined,
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

      {/* Panel Kanan — Ringkasan (360px fixed di desktop, full-width tab di mobile/tablet) */}
      <div
        style={{
          width: isDesktop ? '360px' : '100%',
          flexShrink: 0,
          overflow: 'hidden',
          display: isDesktop || mobileTab === 'keranjang' ? 'flex' : 'none',
          flexDirection: 'column',
        }}
      >
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

import { useEffect, useState } from 'react'
import { toast } from 'sonner'
import { CalendarClock, User } from 'lucide-react'

import { useBreakpoint } from '@/shared/hooks'

import { useBackdateCashDrawerCurrentQuery } from '../backdate-cash/backdate-cash.api'
import { useBackdateCashierStore } from './backdate-cashier.store'
import { BackdateCartEditableList } from './components/BackdateCartEditableList'
import { BackdatePaymentModal } from './components/BackdatePaymentModal'
import { BackdateProductSearch } from './components/BackdateProductSearch'
import { BackdateSummaryPanel } from './components/BackdateSummaryPanel'

export function BackdateCashierPage() {
  const { data: currentDrawer, isLoading: isLoadingDrawer } = useBackdateCashDrawerCurrentQuery()
  const { paymentModalOpen, closePaymentModal, cart } = useBackdateCashierStore()
  const isDesktop = useBreakpoint('lg')
  const [mobileTab, setMobileTab] = useState<'produk' | 'keranjang'>('produk')
  const itemCount = cart.reduce((sum, i) => sum + i.qty, 0)

  useEffect(() => {
    if (!isLoadingDrawer && !currentDrawer) {
      toast.warning('Belum ada kas historis yang dibuka. Buka kas historis terlebih dahulu.', {
        duration: 5000,
        id: 'no-active-backdate-shift',
      })
    }
  }, [currentDrawer, isLoadingDrawer])

  const drawerDate = currentDrawer
    ? new Date(currentDrawer.open_time).toLocaleDateString('id-ID', { weekday: 'short', year: 'numeric', month: 'short', day: 'numeric' })
    : ''

  return (
    <div
      style={{
        display: 'flex',
        flexDirection: 'column',
        height: 'calc(100vh - var(--navbar-height))',
        overflow: 'hidden',
      }}
    >
      {/* Info bar — tanggal + kasir dari kas historis */}
      {currentDrawer && (
        <div className="flex items-center gap-4 px-4 py-2 bg-blue-50 border-b border-blue-200 shrink-0 text-sm">
          <span className="flex items-center gap-1.5 text-blue-700 font-medium">
            <CalendarClock size={14} />
            {drawerDate}
          </span>
          <span className="flex items-center gap-1.5 text-blue-600">
            <User size={14} />
            Kasir: {currentDrawer.user_name}
          </span>
        </div>
      )}

      <div
        style={{
          display: 'flex',
          flexDirection: isDesktop ? 'row' : 'column',
          flex: 1,
          overflow: 'hidden',
        }}
      >
        {/* Mobile tabs */}
        {!isDesktop && (
          <div className="flex shrink-0 border-b bg-white">
            <button
              type="button"
              onClick={() => setMobileTab('produk')}
              className={`flex-1 py-2.5 text-sm font-medium border-b-2 transition-colors ${
                mobileTab === 'produk' ? 'border-blue-600 text-blue-600' : 'border-transparent text-gray-500'
              }`}
            >
              Produk
            </button>
            <button
              type="button"
              onClick={() => setMobileTab('keranjang')}
              className={`flex-1 py-2.5 text-sm font-medium border-b-2 transition-colors relative ${
                mobileTab === 'keranjang' ? 'border-blue-600 text-blue-600' : 'border-transparent text-gray-500'
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

        {/* Left panel */}
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
          <div className="px-4 pt-4 pb-3 shrink-0">
            <BackdateProductSearch />
          </div>
          <BackdateCartEditableList />
        </div>

        {/* Right panel */}
        <div
          style={{
            width: isDesktop ? '360px' : '100%',
            flexShrink: 0,
            overflow: 'hidden',
            display: isDesktop || mobileTab === 'keranjang' ? 'flex' : 'none',
            flexDirection: 'column',
          }}
        >
          <BackdateSummaryPanel />
        </div>
      </div>

      <BackdatePaymentModal
        open={paymentModalOpen}
        onOpenChange={(open) => { if (!open) closePaymentModal() }}
        backdateDrawer={currentDrawer}
      />
    </div>
  )
}

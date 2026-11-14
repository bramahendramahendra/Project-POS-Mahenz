import { useEffect, useState } from 'react'
import { toast } from 'sonner'

import { useCashDrawerCurrentQuery } from '@/features/finance/cash-drawer'
import { useBreakpoint } from '@/shared/hooks'

import { useCashierStore } from './cashier.store'
import { MenuCatalog } from './components/MenuCatalog'
import { PaymentModal } from './components/PaymentModal'
import { SummaryPanel } from './components/SummaryPanel'

/**
 * Kasir Coffee (tahap 4.7b): kiri katalog menu (tombol besar per kategori), kanan panel pesanan.
 * Tablet/HP: dua tab "Menu" dan "Pesanan".
 */
export function CashierPage() {
  const { data: currentDrawer, isLoading: isLoadingDrawer } = useCashDrawerCurrentQuery()
  const { paymentModalOpen, closePaymentModal, cart } = useCashierStore()
  const isDesktop = useBreakpoint('lg')
  const [mobileTab, setMobileTab] = useState<'menu' | 'pesanan'>('menu')
  const itemCount = cart.reduce((sum, i) => sum + i.qty, 0)

  useEffect(() => {
    if (!isLoadingDrawer && !currentDrawer) {
      toast.warning('Belum ada kas yang dibuka. Buka kas/shift terlebih dahulu.', {
        duration: 5000,
        id: 'no-active-shift',
      })
    }
  }, [currentDrawer, isLoadingDrawer])

  const tabClass = (active: boolean) =>
    `relative flex-1 border-b-2 py-2.5 text-sm font-medium transition-colors ${
      active ? 'border-blue-600 text-blue-600' : 'border-transparent text-gray-500'
    }`

  return (
    <div
      // Perangkat sentuh (tablet/HP/APK): tombol dan input kasir minimal 40px agar mudah ditekan.
      className="touch:[&_button:not([role=checkbox])]:min-h-10 touch:[&_button:not([role=checkbox])]:min-w-10 touch:[&_input]:min-h-10"
      style={{
        display: 'flex',
        flexDirection: isDesktop ? 'row' : 'column',
        // Dikurangi padding atas-bawah <main> (2 x 24px) agar tombol Bayar tidak terdorong ke luar layar.
        height: 'calc(100vh - var(--navbar-height) - 48px)',
        overflow: 'hidden',
      }}
    >
      {!isDesktop && (
        <div className="flex shrink-0 border-b bg-white" role="tablist">
          <button type="button" role="tab" aria-selected={mobileTab === 'menu'} onClick={() => setMobileTab('menu')} className={tabClass(mobileTab === 'menu')}>
            Menu
          </button>
          <button type="button" role="tab" aria-selected={mobileTab === 'pesanan'} onClick={() => setMobileTab('pesanan')} className={tabClass(mobileTab === 'pesanan')}>
            Pesanan
            {itemCount > 0 && (
              <span className="ml-1.5 inline-flex h-4 min-w-[16px] items-center justify-center rounded-full bg-blue-600 px-1 text-[10px] text-white">
                {itemCount}
              </span>
            )}
          </button>
        </div>
      )}

      <div
        style={{
          flex: 1,
          display: isDesktop || mobileTab === 'menu' ? 'flex' : 'none',
          flexDirection: 'column',
          overflow: 'hidden',
        }}
        className="bg-gray-50"
      >
        <MenuCatalog />
      </div>

      <div
        style={{
          width: isDesktop ? '380px' : '100%',
          flexShrink: 0,
          overflow: 'hidden',
          display: isDesktop || mobileTab === 'pesanan' ? 'flex' : 'none',
          flexDirection: 'column',
          flex: isDesktop ? undefined : 1,
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

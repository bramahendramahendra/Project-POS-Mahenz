import { Link } from 'react-router-dom'
import { BarChart3, Package, ShoppingBag } from 'lucide-react'

import { Card, CardContent } from '@/shared/components/ui/card'
import { ROUTES } from '@/shared/constants/routes'
import { useAuth } from '@/features/auth/hooks/useAuth'
import { useMenuStore } from '@/features/menu/menu.store'

import { GreetingHeader } from './components/GreetingHeader'
import { CashDrawerStatusCard } from './components/CashDrawerStatusCard'
import { MyRecentTransactions } from './components/MyRecentTransactions'
import { MyTodayPerformance } from './components/MyTodayPerformance'

const SHORTCUTS = [
  { label: 'Ringkasan Bisnis', path: ROUTES.REPORTS_BUSINESS_SUMMARY, icon: BarChart3 },
  { label: 'Produk', path: ROUTES.PRODUCTS, icon: Package },
  { label: 'Pembelian', path: ROUTES.SUPPLIER_PURCHASES, icon: ShoppingBag },
]

function ShortcutGrid() {
  const hasAccess = useMenuStore((s) => s.hasAccess)
  const visible = SHORTCUTS.filter((s) => hasAccess(shortcutMenuKey(s.path)))

  if (visible.length === 0) return null

  return (
    <div className="grid grid-cols-1 gap-3 sm:grid-cols-3">
      {visible.map((s) => (
        <Link key={s.path} to={s.path}>
          <Card className="hover:border-gray-300 transition-colors">
            <CardContent className="pt-4 pb-4 flex items-center gap-3">
              <div className="rounded-lg bg-gray-100 p-2 text-gray-600">
                <s.icon size={18} />
              </div>
              <span className="text-sm font-medium text-gray-800">{s.label}</span>
            </CardContent>
          </Card>
        </Link>
      ))}
    </div>
  )
}

function shortcutMenuKey(path: string): string {
  switch (path) {
    case ROUTES.REPORTS_BUSINESS_SUMMARY: return 'pelaporan.ringkasan_bisnis'
    case ROUTES.PRODUCTS: return 'produk.produk'
    case ROUTES.SUPPLIER_PURCHASES: return 'pengadaan.pembelian'
    default: return ''
  }
}

export function DashboardPage() {
  const { user } = useAuth()
  const hasAccess = useMenuStore((s) => s.hasAccess)
  const isKasirAccessible = hasAccess('penjualan.kasir')

  return (
    <div className="space-y-4">
      <GreetingHeader fullName={user?.fullName ?? ''} />

      {isKasirAccessible ? (
        <>
          <CashDrawerStatusCard />
          <MyTodayPerformance />
          <MyRecentTransactions />
        </>
      ) : (
        <ShortcutGrid />
      )}
    </div>
  )
}

import { Fragment, useState } from 'react'
import { ChevronDown, ChevronRight } from 'lucide-react'
import { useQuery } from '@tanstack/react-query'

import { api } from '@/services'
import { formatRupiah } from '@/shared/utils'
import { ORDER_TYPE_LABEL, type OrderType } from '@/features/sales/cashier/cashier.types'

import type { SalesFilter } from '../sales.types'

// Rincian penjualan Coffee (tahap 4.10, keputusan 115): per menu (bisa dibuka per kombinasi), per kategori
// menu, per metode bayar, per jenis pesanan.

interface SalesBreakdown {
  total_revenue: number
  total_qty: number
  menus: { menu_item_id: number; menu_name: string; category_name: string; qty: number; revenue: number; share: number; variants: { variant_name: string; qty: number; revenue: number }[] }[]
  categories: { category_name: string; qty: number; revenue: number; share: number }[]
  payment_methods: { key: string; count: number; total: number }[]
  order_types: { key: string; count: number; total: number }[]
}

const PAYMENT_LABEL: Record<string, string> = { cash: 'Tunai', transfer: 'Transfer', qris: 'QRIS', card: 'Kartu', kredit: 'Hutang' }
type Tab = 'menu' | 'category' | 'payment' | 'order'
const TABS: { key: Tab; label: string }[] = [
  { key: 'menu', label: 'Per Menu' },
  { key: 'category', label: 'Per Kategori' },
  { key: 'payment', label: 'Metode Bayar' },
  { key: 'order', label: 'Jenis Pesanan' },
]

const pct = (n: number) => `${n.toLocaleString('id-ID', { maximumFractionDigits: 1 })}%`
const th = 'px-3 py-2 text-left text-xs font-medium text-gray-500'
const td = 'px-3 py-2'

export function SalesBreakdownSection({ filter }: { filter: SalesFilter }) {
  const [tab, setTab] = useState<Tab>('menu')
  const [open, setOpen] = useState<Record<number, boolean>>({})
  const range = { date_from: filter.date_from, date_to: filter.date_to }
  const { data, isLoading } = useQuery({
    queryKey: ['reports', 'sales-breakdown', range],
    queryFn: () => api.post<SalesBreakdown>('/reports/sales/breakdown', range),
  })

  return (
    <div className="rounded-lg border bg-white">
      <div className="flex flex-wrap items-center gap-2 border-b px-3 py-2" role="tablist" aria-label="Rincian penjualan">
        {TABS.map((t) => (
          <button
            key={t.key}
            type="button"
            role="tab"
            aria-selected={tab === t.key}
            onClick={() => setTab(t.key)}
            className={`rounded-full px-3 py-1 text-sm font-medium ${tab === t.key ? 'bg-[#2c3e50] text-white' : 'text-gray-600 hover:bg-gray-100'}`}
          >
            {t.label}
          </button>
        ))}
        {data && <span className="ml-auto text-xs text-gray-500">Pendapatan {formatRupiah(data.total_revenue)} (setelah diskon, tanpa pajak)</span>}
      </div>
      <div className="overflow-x-auto">
        {isLoading || !data ? (
          <div className="h-32 animate-pulse bg-gray-50" />
        ) : tab === 'menu' ? (
          <table className="w-full min-w-[520px] text-sm">
            <thead className="border-b bg-gray-50">
              <tr>
                <th className={th}>Menu</th>
                <th className={th}>Kategori</th>
                <th className={`${th} text-right`}>Terjual</th>
                <th className={`${th} text-right`}>Pendapatan</th>
                <th className={`${th} text-right`}>%</th>
              </tr>
            </thead>
            <tbody className="divide-y">
              {data.menus.length === 0 && (
                <tr>
                  <td colSpan={5} className="px-3 py-6 text-center text-gray-400">Belum ada penjualan di periode ini</td>
                </tr>
              )}
              {data.menus.map((m) => (
                <Fragment key={m.menu_item_id}>
                  <tr aria-label={`Penjualan menu ${m.menu_name}`}>
                    <td className={td}>
                      {m.variants.length > 0 ? (
                        <button type="button" className="flex items-center gap-1 font-medium text-gray-800" aria-expanded={!!open[m.menu_item_id]} onClick={() => setOpen((p) => ({ ...p, [m.menu_item_id]: !p[m.menu_item_id] }))}>
                          {open[m.menu_item_id] ? <ChevronDown size={14} /> : <ChevronRight size={14} />}
                          {m.menu_name}
                        </button>
                      ) : (
                        <span className="pl-5 font-medium text-gray-800">{m.menu_name}</span>
                      )}
                    </td>
                    <td className={`${td} text-gray-500`}>{m.category_name || '—'}</td>
                    <td className={`${td} text-right`}>{m.qty}</td>
                    <td className={`${td} text-right font-medium`}>{formatRupiah(m.revenue)}</td>
                    <td className={`${td} text-right text-gray-500`}>{pct(m.share)}</td>
                  </tr>
                  {open[m.menu_item_id] &&
                    m.variants.map((v) => (
                      <tr key={v.variant_name} className="bg-gray-50/60 text-xs text-gray-600">
                        <td className={`${td} pl-10`}>{v.variant_name}</td>
                        <td className={td} />
                        <td className={`${td} text-right`}>{v.qty}</td>
                        <td className={`${td} text-right`}>{formatRupiah(v.revenue)}</td>
                        <td className={td} />
                      </tr>
                    ))}
                </Fragment>
              ))}
            </tbody>
          </table>
        ) : tab === 'category' ? (
          <table className="w-full min-w-[420px] text-sm">
            <thead className="border-b bg-gray-50">
              <tr>
                <th className={th}>Kategori menu</th>
                <th className={`${th} text-right`}>Terjual</th>
                <th className={`${th} text-right`}>Pendapatan</th>
                <th className={`${th} text-right`}>%</th>
              </tr>
            </thead>
            <tbody className="divide-y">
              {data.categories.map((c) => (
                <tr key={c.category_name}>
                  <td className={td}>{c.category_name}</td>
                  <td className={`${td} text-right`}>{c.qty}</td>
                  <td className={`${td} text-right font-medium`}>{formatRupiah(c.revenue)}</td>
                  <td className={`${td} text-right text-gray-500`}>{pct(c.share)}</td>
                </tr>
              ))}
            </tbody>
          </table>
        ) : (
          <table className="w-full min-w-[420px] text-sm">
            <thead className="border-b bg-gray-50">
              <tr>
                <th className={th}>{tab === 'payment' ? 'Metode bayar' : 'Jenis pesanan'}</th>
                <th className={`${th} text-right`}>Transaksi</th>
                <th className={`${th} text-right`}>Total (termasuk pajak)</th>
              </tr>
            </thead>
            <tbody className="divide-y">
              {(tab === 'payment' ? data.payment_methods : data.order_types).map((g) => (
                <tr key={g.key}>
                  <td className={td}>{tab === 'payment' ? (PAYMENT_LABEL[g.key] ?? g.key) : (ORDER_TYPE_LABEL[g.key as OrderType] ?? 'Lainnya')}</td>
                  <td className={`${td} text-right`}>{g.count}</td>
                  <td className={`${td} text-right font-medium`}>{formatRupiah(g.total)}</td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>
    </div>
  )
}

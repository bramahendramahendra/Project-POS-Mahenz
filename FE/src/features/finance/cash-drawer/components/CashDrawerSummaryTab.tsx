import { useState } from 'react'

import { Card, CardContent, CardHeader, CardTitle } from '@/shared/components/ui/card'
import { Badge } from '@/shared/components/ui/badge'
import { ScrollArea, ScrollBar } from '@/shared/components/ui/scroll-area'
import { formatRupiah } from '@/shared/utils'
import { useAuthStore } from '@/features/auth'
import { ROLES } from '@/shared/constants/roles'

import { useCashDrawerSummaryQuery } from '../cash-drawer.api'
import type { CashDrawer, CashDrawerListFilter } from '../cash-drawer.types'
import { CashDrawerFilterBar } from './CashDrawerFilterBar'

function todayString(): string {
  return new Date().toISOString().split('T')[0]
}

function monthStartString(): string {
  const d = new Date()
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-01`
}

function formatDateTime(dateStr: string): string {
  return new Date(dateStr).toLocaleString('id-ID', {
    day: '2-digit',
    month: 'short',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}

const defaultFilter: CashDrawerListFilter = {
  page: 1,
  limit: 1000,
  start_date: monthStartString(),
  end_date: todayString(),
}

interface SummaryCardProps {
  title: string
  value: number
  valueClass?: string
}

function SummaryCard({ title, value, valueClass = '' }: SummaryCardProps) {
  return (
    <Card>
      <CardHeader className="pb-2">
        <CardTitle className="text-sm font-medium text-gray-500">{title}</CardTitle>
      </CardHeader>
      <CardContent>
        <p className={`text-lg font-semibold ${valueClass}`}>{formatRupiah(value)}</p>
      </CardContent>
    </Card>
  )
}

export function CashDrawerSummaryTab() {
  const [filter, setFilter] = useState<CashDrawerListFilter>(defaultFilter)

  const { user } = useAuthStore()
  const isAdminOrOwner = user?.roleName === ROLES.OWNER || user?.roleName === ROLES.ADMIN

  const { data, isLoading } = useCashDrawerSummaryQuery(filter)
  const summary = data

  const handleFilterChange = (newFilter: CashDrawerListFilter) => {
    setFilter(newFilter)
  }

  const handleReset = () => {
    setFilter(defaultFilter)
  }

  return (
    <div className="space-y-6">
      <CashDrawerFilterBar
        filter={filter}
        onChange={handleFilterChange}
        onReset={handleReset}
        showKasirFilter={isAdminOrOwner}
      />

      {isLoading && (
        <div className="py-8 text-center text-sm text-gray-400">Memuat rekap...</div>
      )}

      {summary && !isLoading && (
        <>
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
            <SummaryCard title="Total Saldo Awal Tunai" value={summary.total_opening} />
            <SummaryCard title="Total Saldo Akhir Tunai" value={summary.total_closing} />
            <SummaryCard
              title="Total Pengeluaran"
              value={summary.total_expenses}
              valueClass="text-red-600"
            />
            <SummaryCard
              title="Selisih Bersih"
              value={summary.net}
              valueClass={summary.net >= 0 ? 'text-green-600' : 'text-red-600'}
            />
          </div>

          <ScrollArea className="rounded-lg border border-gray-200">
            <table className="w-full text-sm">
              <thead className="bg-gray-50 text-gray-500 text-xs uppercase">
                <tr>
                  <th className="px-4 py-3 text-left">Waktu Buka</th>
                  <th className="px-4 py-3 text-left">Shift</th>
                  {isAdminOrOwner && <th className="px-4 py-3 text-left">Kasir</th>}
                  <th className="px-4 py-3 text-right">Saldo Awal Tunai</th>
                  <th className="px-4 py-3 text-right">Saldo Akhir Tunai</th>
                  <th className="px-4 py-3 text-right">Pengeluaran</th>
                  <th className="px-4 py-3 text-right">Selisih</th>
                  <th className="px-4 py-3 text-center">Status</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-100">
                {summary.records.length === 0 && (
                  <tr>
                    <td
                      colSpan={isAdminOrOwner ? 8 : 7}
                      className="px-4 py-8 text-center text-gray-400"
                    >
                      Tidak ada data pada periode ini
                    </td>
                  </tr>
                )}
                {summary.records.map((row: CashDrawer) => {
                  const diff = row.difference ?? 0
                  return (
                    <tr key={row.id} className="hover:bg-gray-50">
                      <td className="px-4 py-3 text-gray-600">{formatDateTime(row.open_time)}</td>
                      <td className="px-4 py-3 text-gray-600">{row.shift_name ?? '—'}</td>
                      {isAdminOrOwner && (
                        <td className="px-4 py-3 font-medium">{row.user_name}</td>
                      )}
                      <td className="px-4 py-3 text-right">{formatRupiah(row.opening_balance)}</td>
                      <td className="px-4 py-3 text-right">
                        {row.status === 'closed' && row.closing_balance != null
                          ? formatRupiah(row.closing_balance)
                          : '—'}
                      </td>
                      <td className="px-4 py-3 text-right text-red-600">
                        {formatRupiah(row.total_expenses)}
                      </td>
                      <td
                        className={`px-4 py-3 text-right font-medium ${
                          row.status !== 'closed'
                            ? 'text-gray-400'
                            : diff === 0
                              ? 'text-gray-500'
                              : diff > 0
                                ? 'text-green-600'
                                : 'text-red-600'
                        }`}
                      >
                        {row.status === 'closed'
                          ? `${diff > 0 ? '+' : ''}${formatRupiah(diff)}`
                          : '—'}
                      </td>
                      <td className="px-4 py-3 text-center">
                        {row.status === 'closed' ? (
                          <Badge variant="secondary">Tutup</Badge>
                        ) : (
                          <Badge variant="default" className="bg-green-600">Buka</Badge>
                        )}
                      </td>
                    </tr>
                  )
                })}
              </tbody>
            </table>
            <ScrollBar orientation="horizontal" />
          </ScrollArea>
        </>
      )}
    </div>
  )
}

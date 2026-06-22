import { useState } from 'react'

import { formatRupiah, monthStart, todayStr } from '@/shared/utils'

import { useProfitLossReportQuery } from '../profit-loss.api'
import { ProfitLossFilterBar } from './ProfitLossFilterBar'

interface PLRowProps {
  label: string
  value: number
  indent?: boolean
  bold?: boolean
  colored?: boolean
}

function PLRow({ label, value, indent = false, bold = false, colored = false }: PLRowProps) {
  const colorClass = colored
    ? value >= 0
      ? 'text-green-600'
      : 'text-red-600'
    : 'text-gray-800'

  return (
    <div
      className={`flex justify-between items-center py-2.5 border-b border-gray-100 last:border-0 ${
        indent ? 'pl-4' : ''
      }`}
    >
      <span className={`text-sm ${bold ? 'font-semibold text-gray-900' : 'text-gray-600'}`}>
        {label}
      </span>
      <span className={`text-sm font-medium tabular-nums ${bold ? 'font-bold text-base' : ''} ${colorClass}`}>
        {value < 0 ? '-' : ''}
        {formatRupiah(Math.abs(value))}
      </span>
    </div>
  )
}

function SectionHeader({ children }: { children: React.ReactNode }) {
  return (
    <div className="px-4 py-2 bg-gray-50 rounded-t-lg border-b">
      <p className="text-xs font-semibold text-gray-500 uppercase tracking-wide">{children}</p>
    </div>
  )
}

function Skeleton() {
  return (
    <div className="space-y-3">
      {Array.from({ length: 8 }).map((_, i) => (
        <div key={i} className="h-9 animate-pulse rounded bg-gray-100" />
      ))}
    </div>
  )
}

interface DateFilter {
  date_from?: string
  date_to?: string
}

export function ProfitLossTab() {
  const [filter, setFilter] = useState<DateFilter>({
    date_from: monthStart(),
    date_to: todayStr(),
  })

  const { data, isLoading } = useProfitLossReportQuery({
    date_from: filter.date_from || undefined,
    date_to: filter.date_to || undefined,
  })

  const report = data

  return (
    <div className="space-y-4">
      <ProfitLossFilterBar filter={filter} onChange={setFilter} />

      {isLoading && <Skeleton />}

      {!isLoading && !report && (
        <div className="py-12 text-center text-sm text-gray-400">
          Belum ada data untuk periode yang dipilih
        </div>
      )}

      {!isLoading && report && (
        <div className="max-w-lg space-y-4">
          <div className="rounded-lg border bg-white overflow-hidden">
            <SectionHeader>Pendapatan</SectionHeader>
            <div className="px-4">
              <PLRow label="Total Penjualan" value={report.total_sales} />
              <PLRow label="Retur Penjualan" value={-report.total_returns} indent />
              <PLRow
                label="Net Penjualan"
                value={report.total_sales - report.total_returns}
                bold
              />
            </div>
          </div>

          <div className="rounded-lg border bg-white overflow-hidden">
            <SectionHeader>Pengeluaran</SectionHeader>
            <div className="px-4">
              <PLRow label="Harga Pokok Penjualan (HPP)" value={report.total_hpp} />
              <PLRow label="Total Pengeluaran (Expense)" value={report.total_expense} />
              <PLRow
                label="Total Pengeluaran"
                value={report.total_hpp + report.total_expense}
                bold
              />
            </div>
          </div>

          <div className="rounded-lg border bg-white overflow-hidden">
            <SectionHeader>Hasil</SectionHeader>
            <div className="px-4">
              <PLRow label="Laba Kotor" value={report.gross_profit} bold colored />
              <PLRow label="Laba Bersih" value={report.net_profit} bold colored />
            </div>
          </div>
        </div>
      )}
    </div>
  )
}

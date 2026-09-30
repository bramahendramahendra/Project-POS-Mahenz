import { useState } from 'react'
import type { ReactNode } from 'react'
import { TrendingUp, TrendingDown, AlertTriangle, CheckCircle2 } from 'lucide-react'

import { formatRupiah, monthStart, todayStr } from '@/shared/utils'

import { useProfitLossReportQuery } from '../profit-loss.api'
import type { ProfitLossDateFilter, ProfitLossItem, ProfitLossReport } from '../profit-loss.types'
import { ProfitLossFilterBar } from './ProfitLossFilterBar'

// ── helper format ──────────────────────────────────────────────
function rp(v: number): string {
  return (v < 0 ? '-' : '') + formatRupiah(Math.abs(v))
}

// ── Kartu HERO: untung / rugi ──────────────────────────────────
function HeroCard({ report }: { report: ProfitLossReport }) {
  const net = report.net_profit
  const isProfit = net >= 0

  // margin: dari tiap Rp 100 penjualan, berapa jadi untung bersih
  const marginPct = report.total_revenue > 0 ? (net / report.total_revenue) * 100 : 0
  const marginPer100 = Math.round((marginPct / 100) * 100) // = marginPct dibulatkan
  const marginLabel =
    marginPct >= 15 ? 'margin sehat' : marginPct >= 5 ? 'margin tipis' : marginPct > 0 ? 'margin sangat tipis' : 'rugi'

  // pembanding periode sebelumnya
  let deltaNode: ReactNode = null
  if (report.prev_available) {
    const prev = report.prev_net_profit
    const diff = net - prev
    const naik = diff >= 0
    const arrow = naik ? '▲' : '▼'
    const pctText =
      prev !== 0 ? `${arrow} ${Math.abs((diff / Math.abs(prev)) * 100).toFixed(0)}% ` : `${arrow} `
    deltaNode = (
      <div className="mt-2 flex items-center gap-2">
        <span
          className={`inline-flex items-center gap-1 rounded-full px-2.5 py-1 text-xs font-semibold ${
            naik ? 'bg-green-100 text-green-700' : 'bg-red-100 text-red-700'
          }`}
        >
          {naik ? <TrendingUp size={12} /> : <TrendingDown size={12} />}
          {pctText}vs periode sebelumnya
        </span>
        <span className="text-[11px] text-gray-400">(sebelumnya {rp(prev)})</span>
      </div>
    )
  }

  return (
    <div
      className={`rounded-xl border p-5 ${
        isProfit ? 'border-green-200 bg-green-50' : 'border-red-200 bg-red-50'
      }`}
    >
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        {/* kiri: untung/rugi */}
        <div className="flex items-start gap-3">
          <div
            className={`flex h-12 w-12 shrink-0 items-center justify-center rounded-full ${
              isProfit ? 'bg-green-600' : 'bg-red-600'
            }`}
          >
            {isProfit ? (
              <CheckCircle2 size={26} className="text-white" />
            ) : (
              <AlertTriangle size={24} className="text-white" />
            )}
          </div>
          <div>
            <p className={`text-xs font-medium ${isProfit ? 'text-green-700' : 'text-red-700'}`}>
              Periode ini toko {isProfit ? 'UNTUNG' : 'RUGI'}
            </p>
            <p className={`text-3xl font-extrabold tabular-nums ${isProfit ? 'text-green-600' : 'text-red-600'}`}>
              {rp(net)}
            </p>
            {deltaNode}
          </div>
        </div>

        {/* kanan: margin (hanya kalau untung & ada pendapatan) */}
        {report.total_revenue > 0 && (
          <div className="border-t pt-3 sm:border-l sm:border-t-0 sm:pl-6 sm:pt-0">
            {isProfit ? (
              <>
                <p className="text-[11px] text-gray-500">Dari tiap Rp 100 penjualan</p>
                <p className="text-2xl font-extrabold text-teal-700">Rp {marginPer100}</p>
                <span className="mt-1 inline-block rounded-full bg-teal-100 px-2 py-0.5 text-[11px] font-medium text-teal-700">
                  jadi untung — {marginLabel}
                </span>
              </>
            ) : (
              <div className="max-w-[220px]">
                <p className="text-[11px] font-medium text-gray-500">Penyebab utama</p>
                <p className="text-xs text-gray-700">
                  Modal barang &amp; biaya lebih besar dari uang masuk penjualan. Cek produk yang
                  dijual rugi di bawah.
                </p>
              </div>
            )}
          </div>
        )}
      </div>
    </div>
  )
}

// ── Bar proporsi: kemana uang penjualan pergi ─────────────────
function ProportionBar({ report }: { report: ProfitLossReport }) {
  const rev = report.total_revenue
  if (rev <= 0) return null

  const cogs = Math.max(0, report.total_cogs)
  const exp = Math.max(0, report.total_expenses)
  const profit = Math.max(0, report.net_profit)
  const denom = cogs + exp + profit || 1

  const pct = (v: number) => Math.round((v / denom) * 100)
  const cogsPct = pct(cogs)
  const expPct = pct(exp)
  const profitPct = pct(profit)

  return (
    <div className="rounded-xl border bg-white p-4">
      <p className="text-sm font-semibold text-gray-700">Kemana uang penjualan pergi?</p>
      <p className="mb-3 text-[11px] text-gray-400">Total penjualan {rp(rev)}</p>

      <div className="flex h-7 w-full overflow-hidden rounded-md bg-gray-100">
        {cogs > 0 && (
          <div className="flex items-center justify-center bg-amber-500" style={{ width: `${cogsPct}%` }}>
            {cogsPct >= 12 && <span className="text-[11px] font-semibold text-white">Modal {cogsPct}%</span>}
          </div>
        )}
        {exp > 0 && (
          <div className="flex items-center justify-center bg-indigo-500" style={{ width: `${expPct}%` }}>
            {expPct >= 12 && <span className="text-[11px] font-semibold text-white">Biaya {expPct}%</span>}
          </div>
        )}
        {profit > 0 && (
          <div className="flex items-center justify-center bg-green-600" style={{ width: `${profitPct}%` }}>
            {profitPct >= 8 && <span className="text-[11px] font-semibold text-white">{profitPct}%</span>}
          </div>
        )}
      </div>

      <div className="mt-2 flex flex-wrap gap-x-4 gap-y-1 text-[11px] text-gray-500">
        <span className="flex items-center gap-1"><i className="h-2.5 w-2.5 rounded-sm bg-amber-500" /> Modal barang terjual</span>
        <span className="flex items-center gap-1"><i className="h-2.5 w-2.5 rounded-sm bg-indigo-500" /> Biaya operasional</span>
        <span className="flex items-center gap-1"><i className="h-2.5 w-2.5 rounded-sm bg-green-600" /> Untung bersih</span>
      </div>
      {exp === 0 && (
        <p className="mt-2 text-[11px] text-amber-600">
          Biaya operasional Rp 0 — kalau ada listrik/gaji/sewa yang belum dicatat, untung bersih
          sebenarnya bisa lebih kecil.
        </p>
      )}
    </div>
  )
}

// ── Rincian angka ─────────────────────────────────────────────
function DetailRow({ label, hint, value, tone, bold }: { label: string; hint?: string; value: number; tone?: 'neg' | 'pos'; bold?: boolean }) {
  const color = tone === 'neg' ? 'text-red-600' : tone === 'pos' ? 'text-green-600' : 'text-gray-800'
  return (
    <div className="flex items-center justify-between border-b border-gray-100 py-2.5 last:border-0">
      <span className={bold ? 'font-semibold text-gray-900' : 'text-gray-600'}>
        <span className="text-sm">{label}</span>
        {hint && <span className="block text-[11px] font-normal text-gray-400">{hint}</span>}
      </span>
      <span className={`text-sm tabular-nums ${bold ? 'text-base font-bold' : 'font-medium'} ${color}`}>{rp(value)}</span>
    </div>
  )
}

function DetailCard({ report }: { report: ProfitLossReport }) {
  return (
    <div className="rounded-xl border bg-white p-4">
      <p className="mb-1 text-sm font-semibold text-gray-700">Rincian</p>
      <DetailRow label="Uang masuk penjualan" hint="setelah diskon, tanpa pajak" value={report.total_revenue} />
      <DetailRow label="− Modal barang terjual (HPP)" hint="harga beli barang yang laku" value={report.total_cogs} tone="neg" />
      <DetailRow label="− Biaya operasional" hint="listrik, gaji, sewa, dll" value={report.total_expenses} tone="neg" />
      <DetailRow label="= Untung Kotor" hint="pendapatan − modal barang" value={report.gross_profit} tone={report.gross_profit >= 0 ? 'pos' : 'neg'} bold />
      <div className={`mt-2 flex items-center justify-between rounded-lg border px-3 py-2.5 ${report.net_profit >= 0 ? 'border-green-200 bg-green-50' : 'border-red-200 bg-red-50'}`}>
        <span className="text-sm font-bold text-gray-900">= Untung Bersih</span>
        <span className={`text-base font-extrabold tabular-nums ${report.net_profit >= 0 ? 'text-green-600' : 'text-red-600'}`}>{rp(report.net_profit)}</span>
      </div>
    </div>
  )
}

// ── Produk: paling untung + jual rugi ─────────────────────────
function ProductInsight({ items }: { items: ProfitLossItem[] }) {
  const sorted = [...items].sort((a, b) => b.gross_profit - a.gross_profit)
  const topProfit = sorted.filter((i) => i.gross_profit > 0).slice(0, 5)
  const lossProducts = sorted.filter((i) => i.gross_profit < 0).sort((a, b) => a.gross_profit - b.gross_profit)
  const maxProfit = topProfit.length > 0 ? topProfit[0].gross_profit : 1

  return (
    <div className="rounded-xl border bg-white p-4">
      <p className="mb-3 text-sm font-semibold text-gray-700">Produk paling menguntungkan</p>
      {topProfit.length === 0 && <p className="text-xs text-gray-400">Belum ada produk yang menghasilkan untung di periode ini.</p>}
      <div className="space-y-2">
        {topProfit.map((it, idx) => (
          <div key={it.product_id} className="flex items-center gap-3">
            <span className="w-4 text-right text-xs text-gray-400">{idx + 1}</span>
            <span className="w-40 shrink-0 truncate text-sm text-gray-600" title={it.product_name}>{it.product_name}</span>
            <div className="h-2.5 flex-1 overflow-hidden rounded-full bg-green-100">
              <div className="h-full rounded-full bg-green-500" style={{ width: `${Math.max(4, (it.gross_profit / maxProfit) * 100)}%` }} />
            </div>
            <span className="w-28 text-right text-sm font-semibold tabular-nums text-green-600">+ {formatRupiah(it.gross_profit)}</span>
          </div>
        ))}
      </div>

      {lossProducts.length > 0 && (
        <div className="mt-4 rounded-lg border border-red-200 bg-red-50 p-3">
          <p className="flex items-center gap-1.5 text-sm font-bold text-red-700">
            <AlertTriangle size={14} /> Produk yang dijual RUGI ({lossProducts.length})
          </p>
          <div className="mt-2 space-y-1.5">
            {lossProducts.slice(0, 5).map((it) => (
              <div key={it.product_id} className="flex items-center justify-between text-xs">
                <span className="truncate text-red-900" title={it.product_name}>
                  {it.product_name}
                  <span className="text-red-400"> — jual {formatRupiah(it.total_revenue / (it.qty_sold || 1))}/unit, modal {formatRupiah(it.purchase_price)}</span>
                </span>
                <span className="ml-2 shrink-0 font-semibold text-red-600">{rp(it.gross_profit)}</span>
              </div>
            ))}
          </div>
          <p className="mt-2 text-[11px] text-red-500">Pertimbangkan naikkan harga jual, atau cek harga beli produk ini di menu Rekonsiliasi Modal.</p>
        </div>
      )}
    </div>
  )
}

function Skeleton() {
  return (
    <div className="space-y-3">
      {Array.from({ length: 6 }).map((_, i) => (
        <div key={i} className="h-16 animate-pulse rounded-xl bg-gray-100" />
      ))}
    </div>
  )
}

export function ProfitLossTab() {
  const [filter, setFilter] = useState<ProfitLossDateFilter>({
    date_from: monthStart(),
    date_to: todayStr(),
  })

  const { data: report, isLoading } = useProfitLossReportQuery(filter)

  const handleFilterChange = (newFilter: ProfitLossDateFilter) => setFilter(newFilter)

  return (
    <div className="space-y-4">
      <ProfitLossFilterBar filter={filter} onChange={handleFilterChange} />

      {isLoading && <Skeleton />}

      {!isLoading && !report && (
        <div className="py-12 text-center text-sm text-gray-400">Belum ada data untuk periode yang dipilih</div>
      )}

      {!isLoading && report && (
        <div className="space-y-4">
          <HeroCard report={report} />
          <div className="grid grid-cols-1 gap-4 lg:grid-cols-2">
            <div className="space-y-4">
              <ProportionBar report={report} />
              <DetailCard report={report} />
            </div>
            <ProductInsight items={report.items ?? []} />
          </div>
          <p className="text-[11px] leading-relaxed text-gray-400">
            Catatan: &quot;Modal barang terjual (HPP)&quot; adalah harga beli barang yang benar-benar laku,
            bukan total belanja. Angka merah berarti rugi.
          </p>
        </div>
      )}
    </div>
  )
}

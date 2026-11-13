import { FormModal } from '@/shared/components'
import { Badge } from '@/shared/components/ui/badge'
import { formatDateTime, formatRupiah } from '@/shared/utils'

import { useCashDrawerDetailQuery } from '../cash-drawer.api'
import type { CashDrawerDetail } from '../cash-drawer.types'

interface CashDrawerDetailModalProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  cashDrawerId?: number
}

interface SummaryRowProps {
  label: string
  value: string
  valueClass?: string
}

function SummaryRow({ label, value, valueClass = '' }: SummaryRowProps) {
  return (
    <div className="flex justify-between items-center py-2 border-b border-gray-100 last:border-0">
      <span className="text-sm text-gray-500">{label}</span>
      <span className={`text-sm font-medium ${valueClass}`}>{value}</span>
    </div>
  )
}

export function CashDrawerDetailModal({ open, onOpenChange, cashDrawerId }: CashDrawerDetailModalProps) {
  const enabled = open && (cashDrawerId ?? 0) > 0
  const { data, isLoading } = useCashDrawerDetailQuery(enabled ? (cashDrawerId as number) : null)
  const detail: CashDrawerDetail | undefined = data ?? undefined

  const diff = detail?.difference ?? 0

  return (
    <FormModal
      open={open}
      onOpenChange={onOpenChange}
      title="Detail Kas Harian"
      size="sm"
      hideFooter
    >
      {isLoading || !detail ? (
        <div className="space-y-4">
          {[1, 2, 3, 4].map((i) => (
            <div key={i} className="h-8 animate-pulse rounded-md bg-gray-100" />
          ))}
        </div>
      ) : (
        <div className="space-y-4">
          <div className="flex items-center justify-between">
            <div>
              <p className="font-semibold text-gray-700">{formatDateTime(detail.open_time)}</p>
              {detail.shift_name && (
                <p className="text-xs text-gray-500 mt-0.5">
                  {detail.shift_name}
                  {detail.shift_start && detail.shift_end
                    ? ` (${detail.shift_start} – ${detail.shift_end})`
                    : ''}
                </p>
              )}
            </div>
            {detail.status === 'closed' ? (
              <Badge variant="secondary">Tutup</Badge>
            ) : (
              <Badge variant="default" className="bg-green-600">Buka</Badge>
            )}
          </div>

          <div className="bg-gray-50 rounded-lg p-4 space-y-1">
            <SummaryRow label="Kasir" value={detail.cashier_name} />
            <SummaryRow label="Saldo Awal Tunai" value={formatRupiah(detail.opening_balance)} />
            <SummaryRow
              label="Total Penjualan Tunai"
              value={formatRupiah(detail.total_cash_sales)}
              valueClass="text-green-600"
            />
            <SummaryRow
              label="Total Pengeluaran"
              value={formatRupiah(detail.total_expenses)}
              valueClass="text-red-600"
            />
            <SummaryRow
              label="Saldo Ekspektasi"
              value={formatRupiah(detail.expected_balance)}
            />
            {detail.status === 'closed' && (
              <>
                <SummaryRow
                  label="Saldo Akhir Tunai"
                  value={formatRupiah(detail.closing_balance ?? 0)}
                  valueClass="font-semibold"
                />
                <SummaryRow
                  label="Selisih"
                  value={`${diff >= 0 ? '+' : ''}${formatRupiah(diff)}`}
                  valueClass={
                    diff === 0 ? 'text-gray-500' : diff > 0 ? 'text-green-600' : 'text-red-600'
                  }
                />
              </>
            )}
          </div>

          {detail.non_cash_sales?.length > 0 && (
            <div className="rounded-lg border border-gray-100 px-4 py-3 space-y-1">
              <p className="text-xs font-medium text-gray-500 uppercase tracking-wide pb-1">Non-Tunai (Informasi)</p>
              {detail.non_cash_sales.map((item) => (
                <SummaryRow key={item.payment_method} label={item.label} value={formatRupiah(item.total)} />
              ))}
              <SummaryRow
                label="Total Non-Tunai"
                value={formatRupiah(detail.non_cash_sales.reduce((sum, i) => sum + i.total, 0))}
                valueClass="font-semibold"
              />
            </div>
          )}

          {detail.status === 'closed' && (
            <div className="space-y-1 text-sm">
              {detail.close_time && (
                <p className="text-gray-500">
                  Ditutup: <span className="text-gray-700">{formatDateTime(detail.close_time)}</span>
                </p>
              )}
              {detail.notes && (
                <p className="text-gray-500">
                  Catatan tutup: <span className="text-gray-700">{detail.notes}</span>
                </p>
              )}
            </div>
          )}
          {detail.open_notes && (
            <p className="text-sm text-gray-500">
              Catatan buka: <span className="text-gray-700">{detail.open_notes}</span>
            </p>
          )}
        </div>
      )}
    </FormModal>
  )
}

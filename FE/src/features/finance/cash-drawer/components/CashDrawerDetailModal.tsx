import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/shared/components/ui/dialog'
import { Badge } from '@/shared/components/ui/badge'
import { formatRupiah } from '@/shared/utils'

import { useCashDrawerDetailQuery } from '../cash-drawer.api'
import type { CashDrawerDetail } from '../cash-drawer.types'

interface CashDrawerDetailModalProps {
  cashDrawerId: number | null
  onClose: () => void
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

export function CashDrawerDetailModal({ cashDrawerId, onClose }: CashDrawerDetailModalProps) {
  const { data, isLoading } = useCashDrawerDetailQuery(cashDrawerId)
  const detail: CashDrawerDetail | undefined = data ?? undefined

  const diff = detail?.difference ?? 0

  return (
    <Dialog open={cashDrawerId !== null} onOpenChange={(open) => !open && onClose()}>
      <DialogContent className="max-w-md">
        <DialogHeader>
          <DialogTitle>Detail Kas Harian</DialogTitle>
        </DialogHeader>

        {isLoading && (
          <div className="py-8 text-center text-sm text-gray-400">Memuat data...</div>
        )}

        {detail && !isLoading && (
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
              <SummaryRow label="Saldo Awal" value={formatRupiah(detail.opening_balance)} />
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
                    label="Saldo Akhir (Aktual)"
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
      </DialogContent>
    </Dialog>
  )
}

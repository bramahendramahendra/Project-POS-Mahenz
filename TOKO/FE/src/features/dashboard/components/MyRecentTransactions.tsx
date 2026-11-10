import { useState } from 'react'
import { Receipt } from 'lucide-react'

import { Card, CardContent } from '@/shared/components/ui/card'
import { StatusBadge } from '@/shared/components'
import { formatRupiah } from '@/shared/utils'
import { TransactionDetailModal } from '@/features/sales/transactions/components/TransactionDetailModal'
import { PAYMENT_LABELS } from '@/features/sales/transactions/transactions.utils'
import type { PaymentMethod } from '@/features/sales/transactions/transactions.types'

import { useRecentTransactionsQuery } from '../dashboard.api'

function formatDateTimeShort(dateStr: string): string {
  return new Date(dateStr).toLocaleString('id-ID', {
    day: '2-digit', month: '2-digit', year: 'numeric', hour: '2-digit', minute: '2-digit',
  })
}

export function MyRecentTransactions() {
  const { data, isLoading } = useRecentTransactionsQuery(5)
  const [selectedId, setSelectedId] = useState<number | null>(null)

  const items = data ?? []

  return (
    <>
      <Card>
        <CardContent className="pt-4 pb-4">
          <h3 className="font-semibold text-gray-700 text-sm mb-3">Transaksi Terakhir Saya</h3>

          {isLoading ? (
            <div className="space-y-2">
              {Array.from({ length: 3 }).map((_, i) => (
                <div key={i} className="h-12 animate-pulse rounded-lg bg-gray-100" />
              ))}
            </div>
          ) : items.length === 0 ? (
            <div className="flex flex-col items-center justify-center py-8 text-gray-400">
              <Receipt size={28} className="mb-2" />
              <p className="text-sm">Belum ada transaksi</p>
            </div>
          ) : (
            <ul className="divide-y divide-gray-100">
              {items.map((item) => (
                <li key={item.id}>
                  <button
                    type="button"
                    onClick={() => setSelectedId(item.id)}
                    className="w-full flex items-center justify-between gap-3 py-2.5 text-left hover:bg-gray-50 -mx-2 px-2 rounded-md transition-colors"
                  >
                    <div className="min-w-0">
                      <p className="text-sm font-medium text-gray-800 truncate">{item.transaction_code}</p>
                      <p className="text-xs text-gray-400">
                        {formatDateTimeShort(item.transaction_date)} · {PAYMENT_LABELS[item.payment_method as PaymentMethod] ?? item.payment_method}
                      </p>
                    </div>
                    <div className="flex items-center gap-2 shrink-0">
                      <span className="text-sm font-semibold text-gray-800">{formatRupiah(item.total_amount)}</span>
                      <StatusBadge status={item.status === 'completed' ? 'success' : 'void'} />
                    </div>
                  </button>
                </li>
              ))}
            </ul>
          )}
        </CardContent>
      </Card>

      {selectedId !== null && (
        <TransactionDetailModal transactionId={selectedId} onClose={() => setSelectedId(null)} />
      )}
    </>
  )
}

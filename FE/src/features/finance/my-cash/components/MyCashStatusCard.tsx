import { useState } from 'react'

import { Badge } from '@/shared/components/ui/badge'
import { Button } from '@/shared/components/ui/button'
import { Card, CardContent } from '@/shared/components/ui/card'
import { formatRupiah } from '@/shared/utils'
import { OpenCashDrawerModal } from '@/features/finance/cash-drawer/components/OpenCashDrawerModal'
import { CloseCashDrawerModal } from '@/features/finance/cash-drawer/components/CloseCashDrawerModal'

import type { MyCashData } from '../my-cash.types'

interface MyCashStatusCardProps {
  data?: MyCashData
  isLoading: boolean
}

function StatBox({
  label,
  value,
  valueClass,
}: {
  label: string
  value: string
  valueClass?: string
}) {
  return (
    <div className="rounded-lg bg-gray-50 px-4 py-3 flex flex-col gap-1">
      <p className="text-xs text-gray-500">{label}</p>
      <p className={`text-sm font-semibold ${valueClass ?? ''}`}>{value}</p>
    </div>
  )
}

export function MyCashStatusCard({ data, isLoading }: MyCashStatusCardProps) {
  const [openModalOpen, setOpenModalOpen] = useState(false)
  const [closeModalOpen, setCloseModalOpen] = useState(false)

  const isOpen = data?.status === 'open'

  if (isLoading) {
    return (
      <Card>
        <CardContent className="pt-4 pb-4 space-y-3">
          <div className="flex items-center justify-between">
            <div className="h-5 w-24 animate-pulse rounded bg-gray-100" />
            <div className="h-8 w-24 animate-pulse rounded bg-gray-100" />
          </div>
          <div className="grid grid-cols-2 gap-3 sm:grid-cols-4">
            {[...Array(4)].map((_, i) => (
              <div key={i} className="h-16 animate-pulse rounded-lg bg-gray-100" />
            ))}
          </div>
        </CardContent>
      </Card>
    )
  }

  return (
    <>
      <Card>
        <CardContent className="pt-4 pb-4 space-y-3">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-2">
              {isOpen ? (
                <Badge variant="default" className="bg-green-600">● Buka</Badge>
              ) : (
                <Badge variant="secondary">● Tutup</Badge>
              )}
              {isOpen && data?.shift_name && (
                <span className="text-sm text-gray-500">
                  {data.shift_name}
                  {data.shift_start && data.shift_end
                    ? ` (${data.shift_start} – ${data.shift_end})`
                    : ''}
                </span>
              )}
            </div>

            {isOpen ? (
              <Button
                variant="outline"
                size="sm"
                onClick={() => setCloseModalOpen(true)}
                className="text-red-600 border-red-200 hover:bg-red-50 hover:text-red-700"
              >
                Tutup Kas
              </Button>
            ) : (
              <Button size="sm" onClick={() => setOpenModalOpen(true)}>
                Buka Kas
              </Button>
            )}
          </div>

          {isOpen && (
            <>
              <div className="grid grid-cols-2 gap-3 sm:grid-cols-4">
                <StatBox
                  label="Saldo Awal Tunai"
                  value={formatRupiah(data?.opening_balance ?? 0)}
                />
                <StatBox
                  label="Total Masuk Tunai"
                  value={formatRupiah(data?.total_cash_sales ?? 0)}
                  valueClass="text-green-600"
                />
                <StatBox
                  label="Total Keluar"
                  value={formatRupiah(data?.total_expenses ?? 0)}
                  valueClass="text-red-600"
                />
                <StatBox
                  label="Saldo Ekspektasi"
                  value={formatRupiah(data?.expected_balance ?? 0)}
                  valueClass="text-blue-600"
                />
              </div>

              {(data?.non_cash_sales?.length ?? 0) > 0 && (
                <div className="rounded-lg border border-gray-100 px-4 py-3 space-y-2">
                  <p className="text-xs font-medium text-gray-500 uppercase tracking-wide">Non-Tunai (Informasi)</p>
                  <div className="grid grid-cols-2 gap-2 sm:grid-cols-3">
                    {data!.non_cash_sales.map((item) => (
                      <StatBox
                        key={item.payment_method}
                        label={item.label}
                        value={formatRupiah(item.total)}
                      />
                    ))}
                  </div>
                  <div className="flex justify-between items-center pt-1 border-t border-gray-100">
                    <p className="text-xs text-gray-500">Total Non-Tunai</p>
                    <p className="text-sm font-semibold text-gray-700">
                      {formatRupiah(data!.non_cash_sales.reduce((sum, i) => sum + i.total, 0))}
                    </p>
                  </div>
                </div>
              )}
            </>
          )}
        </CardContent>
      </Card>

      <OpenCashDrawerModal
        open={openModalOpen}
        onOpenChange={setOpenModalOpen}
      />

      <CloseCashDrawerModal
        open={closeModalOpen}
        onOpenChange={setCloseModalOpen}
        cashDrawerId={data?.id ?? null}
      />
    </>
  )
}

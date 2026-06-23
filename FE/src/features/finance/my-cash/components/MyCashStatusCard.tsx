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

export function MyCashStatusCard({ data, isLoading }: MyCashStatusCardProps) {
  const [openModalOpen, setOpenModalOpen] = useState(false)
  const [closeModalOpen, setCloseModalOpen] = useState(false)

  const isOpen = data?.status === 'open'

  if (isLoading) {
    return (
      <Card>
        <CardContent className="pt-4 pb-4 space-y-3">
          <div className="h-5 w-24 animate-pulse rounded bg-gray-100" />
          <div className="h-8 w-48 animate-pulse rounded bg-gray-100" />
          <div className="grid grid-cols-3 gap-4 pt-2">
            <div className="h-12 animate-pulse rounded bg-gray-100" />
            <div className="h-12 animate-pulse rounded bg-gray-100" />
            <div className="h-12 animate-pulse rounded bg-gray-100" />
          </div>
        </CardContent>
      </Card>
    )
  }

  return (
    <>
      <Card>
        <CardContent className="pt-4 pb-4">
          <div className="flex items-start justify-between">
            <div className="space-y-3 flex-1">
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

              {isOpen && (
                <div className="grid grid-cols-2 gap-x-8 gap-y-2 sm:grid-cols-4">
                  <div>
                    <p className="text-xs text-gray-500">Saldo Buka</p>
                    <p className="text-sm font-semibold">{formatRupiah(data?.opening_balance ?? 0)}</p>
                  </div>
                  <div>
                    <p className="text-xs text-gray-500">Total Masuk</p>
                    <p className="text-sm font-semibold text-green-600">{formatRupiah(data?.total_cash_sales ?? 0)}</p>
                  </div>
                  <div>
                    <p className="text-xs text-gray-500">Total Keluar</p>
                    <p className="text-sm font-semibold text-red-600">{formatRupiah(data?.total_expenses ?? 0)}</p>
                  </div>
                  <div>
                    <p className="text-xs text-gray-500">Saldo Ekspektasi</p>
                    <p className="text-sm font-semibold">{formatRupiah(data?.expected_balance ?? 0)}</p>
                  </div>
                </div>
              )}
            </div>

            <div className="ml-4">
              {isOpen ? (
                <Button
                  variant="outline"
                  onClick={() => setCloseModalOpen(true)}
                  className="text-red-600 border-red-200 hover:bg-red-50 hover:text-red-700"
                >
                  Tutup Kas
                </Button>
              ) : (
                <Button onClick={() => setOpenModalOpen(true)}>
                  Buka Kas
                </Button>
              )}
            </div>
          </div>
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

import { useState } from 'react'

import { PageHeader, DataTable } from '@/shared/components'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/shared/components/ui/tabs'

import { useMyCashQuery } from './my-cash.api'
import type { CashDrawerTransaction, CashDrawerExpenseItem, NonCashTransaction } from './my-cash.types'
import { MyCashStatusCard } from './components/MyCashStatusCard'
import { buildMyCashTransactionColumns } from './components/MyCashTransactionColumns'
import { buildMyCashNonCashTransactionColumns } from './components/MyCashNonCashTransactionColumns'
import { buildMyCashExpenseColumns } from './components/MyCashExpenseColumns'

export function MyCashPage() {
  const { data, isLoading } = useMyCashQuery()
  const [activeTab, setActiveTab] = useState('transaksi-tunai')

  const transactions = (data?.transactions ?? []) as CashDrawerTransaction[]
  const nonCashTransactions = (data?.non_cash_transactions ?? []) as NonCashTransaction[]
  const expenses = (data?.expenses ?? []) as CashDrawerExpenseItem[]

  const transactionColumns = buildMyCashTransactionColumns()
  const nonCashTransactionColumns = buildMyCashNonCashTransactionColumns()
  const expenseColumns = buildMyCashExpenseColumns()

  return (
    <div className="space-y-4">
      <PageHeader
        title="Kas Saya"
        breadcrumbs={[{ label: 'Keuangan' }, { label: 'Kas Saya' }]}
      />

      <MyCashStatusCard data={data ?? undefined} isLoading={isLoading} />

      <Tabs value={activeTab} onValueChange={setActiveTab}>
        <TabsList>
          <TabsTrigger value="transaksi-tunai">
            Transaksi Tunai ({transactions.length})
          </TabsTrigger>
          <TabsTrigger value="transaksi-non-tunai">
            Transaksi Non-Tunai ({nonCashTransactions.length})
          </TabsTrigger>
          <TabsTrigger value="pengeluaran">
            Pengeluaran ({expenses.length})
          </TabsTrigger>
        </TabsList>

        <TabsContent value="transaksi-tunai" className="mt-3">
          <DataTable<CashDrawerTransaction & Record<string, unknown>>
            columns={transactionColumns}
            data={transactions as (CashDrawerTransaction & Record<string, unknown>)[]}
            isLoading={isLoading}
            emptyMessage="Belum ada transaksi tunai hari ini"
          />
        </TabsContent>

        <TabsContent value="transaksi-non-tunai" className="mt-3">
          <DataTable<NonCashTransaction & Record<string, unknown>>
            columns={nonCashTransactionColumns}
            data={nonCashTransactions as (NonCashTransaction & Record<string, unknown>)[]}
            isLoading={isLoading}
            emptyMessage="Belum ada transaksi non-tunai hari ini"
          />
        </TabsContent>

        <TabsContent value="pengeluaran" className="mt-3">
          <DataTable<CashDrawerExpenseItem & Record<string, unknown>>
            columns={expenseColumns}
            data={expenses as (CashDrawerExpenseItem & Record<string, unknown>)[]}
            isLoading={isLoading}
            emptyMessage="Belum ada pengeluaran hari ini"
          />
        </TabsContent>
      </Tabs>
    </div>
  )
}

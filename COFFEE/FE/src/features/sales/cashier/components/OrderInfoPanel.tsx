import { Input } from '@/shared/components/ui/input'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/shared/components/ui/select'
import { diningTableApi } from '@/features/menu-items/menu-items.api'

import { useCashierStore } from '../cashier.store'
import { ORDER_TYPE_LABEL, type OrderType } from '../cashier.types'

const ORDER_TYPES: OrderType[] = ['take_away', 'dine_in']

/** Jenis pesanan, meja (makan di tempat), dan nama pemesan opsional (keputusan 25, 97, 99). */
export function OrderInfoPanel() {
  const { orderType, table, orderName, setOrderType, setTable, setOrderName } = useCashierStore()
  const { data: tables = [] } = diningTableApi.useOptions()
  const isDineIn = orderType === 'dine_in'

  return (
    <div className="shrink-0 space-y-2.5 border-b px-4 py-3">
      <div className="grid grid-cols-2 gap-2" role="radiogroup" aria-label="Jenis pesanan">
        {ORDER_TYPES.map((t) => (
          <button
            key={t}
            type="button"
            role="radio"
            aria-checked={orderType === t}
            onClick={() => setOrderType(t)}
            className={`min-h-10 rounded-lg border-2 px-2 text-sm font-medium transition-colors ${
              orderType === t ? 'border-[#2c3e50] bg-[#2c3e50] text-white' : 'border-gray-200 text-gray-600 hover:border-gray-300'
            }`}
          >
            {ORDER_TYPE_LABEL[t]}
          </button>
        ))}
      </div>
      <div className={`grid gap-2 ${isDineIn ? 'grid-cols-2' : 'grid-cols-1'}`}>
        {isDineIn && (
          <Select
            value={table ? String(table.id) : ''}
            onValueChange={(v) => {
              const found = tables.find((t) => String(t.id) === v)
              setTable(found ? { id: found.id, name: found.name } : null)
            }}
          >
            <SelectTrigger className={`w-full ${table ? '' : 'border-amber-400'}`} aria-label="Meja">
              <SelectValue placeholder="Pilih meja" />
            </SelectTrigger>
            <SelectContent>
              {tables.length === 0 && <div className="px-3 py-2 text-xs text-gray-500">Belum ada meja aktif</div>}
              {tables.map((t) => (
                <SelectItem key={t.id} value={String(t.id)}>
                  {t.name}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        )}
        <Input
          value={orderName}
          maxLength={50}
          onChange={(e) => setOrderName(e.target.value)}
          placeholder="Nama pemesan (opsional)"
          aria-label="Nama pemesan"
        />
      </div>
    </div>
  )
}

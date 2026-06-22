import { useState } from 'react'

import { Checkbox } from '@/shared/components/ui/checkbox'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/shared/components/ui/select'

import { useCustomerListQuery } from '../cashier.api'
import { useCashierStore } from '../cashier.store'

export function CustomerSelector() {
  const [showCustomer, setShowCustomer] = useState(false)
  const { selectedCustomer, setCustomer } = useCashierStore()
  const { data: customerData } = useCustomerListQuery({ page: 1, limit: 200, search: '' })
  const customers = customerData?.data ?? []

  return (
    <div className="px-4 py-2.5 border-b shrink-0">
      <label className="flex items-center gap-2 cursor-pointer select-none w-fit">
        <Checkbox
          checked={showCustomer}
          onCheckedChange={(v) => {
            setShowCustomer(!!v)
            if (!v) setCustomer(null)
          }}
        />
        <span className="text-xs text-gray-500">Tambah Pelanggan</span>
        {selectedCustomer && (
          <span className="text-xs font-medium text-blue-600">— {selectedCustomer.name}</span>
        )}
      </label>
      {showCustomer && (
        <div className="mt-2">
          <Select
            value={selectedCustomer ? String(selectedCustomer.id) : 'none'}
            onValueChange={(v) => {
              if (v === 'none') {
                setCustomer(null)
              } else {
                const c = customers.find((c) => String(c.id) === v)
                if (c) setCustomer({ id: c.id, name: c.name })
              }
            }}
          >
            <SelectTrigger className="h-8 text-sm border-dashed">
              <SelectValue placeholder="Pilih pelanggan..." />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="none">— Tanpa Pelanggan —</SelectItem>
              {customers.map((c) => (
                <SelectItem key={c.id} value={String(c.id)}>
                  {c.name}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>
      )}
    </div>
  )
}

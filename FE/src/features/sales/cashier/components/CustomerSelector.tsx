import { useState } from 'react'
import { X } from 'lucide-react'

import { Checkbox } from '@/shared/components/ui/checkbox'
import { Button } from '@/shared/components/ui/button'
import { AsyncCombobox } from '@/shared/components/ui/async-combobox'

import { useCustomerListQuery } from '../cashier.api'
import { useCashierStore } from '../cashier.store'
import type { Customer } from '@/features/customers'

export function CustomerSelector() {
  const [showCustomer, setShowCustomer] = useState(false)
  const [keyword, setKeyword] = useState('')
  const { selectedCustomer, setCustomer } = useCashierStore()

  const { data: customerData, isFetching } = useCustomerListQuery({
    page: 1,
    limit: 20,
    search: keyword,
    is_active: true,
  })
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
        <div className="mt-2 flex items-center gap-1.5">
          <AsyncCombobox<Customer>
            className="h-8 text-sm border-dashed"
            value={selectedCustomer?.id}
            selectedLabel={selectedCustomer?.name}
            onSearch={setKeyword}
            onValueChange={(_v, item) => {
              if (item) setCustomer({ id: item.id, name: item.name })
            }}
            options={customers}
            getOptionValue={(c) => c.id}
            getOptionLabel={(c) => c.name}
            isLoading={isFetching}
            placeholder="Pilih pelanggan..."
            searchPlaceholder="Cari nama pelanggan..."
            emptyText="Pelanggan tidak ditemukan."
          />
          {selectedCustomer && (
            <Button
              type="button"
              variant="ghost"
              size="icon"
              className="h-8 w-8 shrink-0 text-gray-400 hover:text-gray-600"
              onClick={() => setCustomer(null)}
            >
              <X size={14} />
            </Button>
          )}
        </div>
      )}
    </div>
  )
}

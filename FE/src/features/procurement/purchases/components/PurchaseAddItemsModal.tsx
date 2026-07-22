import { useState } from 'react'
import { useForm, useFieldArray, useWatch, Controller } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { Plus, Trash2 } from 'lucide-react'

import { ConfirmDialog, FormModal } from '@/shared/components'
import { Button } from '@/shared/components/ui/button'
import { Input } from '@/shared/components/ui/input'
import { RupiahInput } from '@/shared/components/ui/rupiah-input'
import { AsyncCombobox } from '@/shared/components/ui/async-combobox'
import { formatRupiah } from '@/shared/utils'
import { useQueryClient } from '@tanstack/react-query'
import { queryKeys } from '@/shared/constants'
import { api } from '@/services'
import { useProductSearchQuery } from '@/features/products/products'
import type { ProductPackage, ProductSearchOption } from '@/features/products/products/products.types'

import { purchaseItemSchema } from '../purchases.schema'
import { useSupplierPurchaseDetailQuery, useAddPurchaseItemsMutation } from '../purchases.api'
import type { SupplierPurchase } from '../purchases.types'

const addItemsSchema = z.object({
  items: z.array(purchaseItemSchema).min(1, 'Minimal tambah 1 item'),
})

type AddItemsFormValues = z.infer<typeof addItemsSchema>

const emptyItem = { product_id: 0, product_name: '', quantity: 1, price: 0, unit: '', conversion_qty: 1 }

interface PurchaseAddItemsModalProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  purchase: SupplierPurchase | null
}

function ProductCell({
  value,
  label,
  onChange,
  error,
}: {
  value: number
  label?: string
  onChange: (id: number, name: string) => void
  error?: boolean
}) {
  const [keyword, setKeyword] = useState('')
  const { data: results, isFetching } = useProductSearchQuery(keyword)

  return (
    <AsyncCombobox<ProductSearchOption>
      className={`h-8 text-xs ${error ? 'border-red-500' : ''}`}
      value={value > 0 ? value : undefined}
      selectedLabel={label}
      onSearch={setKeyword}
      onValueChange={(_v, item) => {
        if (item) onChange(item.id, item.name)
      }}
      options={results}
      getOptionValue={(p) => p.id}
      getOptionLabel={(p) => p.name}
      isLoading={isFetching}
      placeholder="Pilih produk"
      searchPlaceholder="Cari produk..."
      emptyText="Produk tidak ditemukan."
    />
  )
}

interface ConfirmState {
  isConfirming: boolean
  pendingValues: AddItemsFormValues | null
}

const initialConfirmState: ConfirmState = { isConfirming: false, pendingValues: null }

export function PurchaseAddItemsModal({ open, onOpenChange, purchase }: PurchaseAddItemsModalProps) {
  // isConfirming & pendingValues digabung jadi 1 state supaya efek "reset saat modal
  // ditutup" cukup 1 kali setState, bukan berurutan (react-hooks/set-state-in-effect).
  const [confirmState, setConfirmState] = useState<ConfirmState>(initialConfirmState)

  const purchaseId = purchase?.id ?? null
  const { data: detail } = useSupplierPurchaseDetailQuery(open ? purchaseId : null)
  const { mutate: addItems, isPending } = useAddPurchaseItemsMutation()

  const {
    control,
    handleSubmit,
    setValue,
    reset,
    formState: { errors },
  } = useForm<AddItemsFormValues>({
    resolver: zodResolver(addItemsSchema),
    defaultValues: { items: [emptyItem] },
  })

  const { fields, append, remove } = useFieldArray({ control, name: 'items' })
  const queryClient = useQueryClient()
  const watchItems = useWatch({ control, name: 'items' })
  const [itemPackageId, setItemPackageId] = useState<Record<number, number>>({})

  async function handleProductChange(index: number, id: number, productName: string) {
    setValue(`items.${index}.product_id`, id, { shouldValidate: true })
    setValue(`items.${index}.product_name`, productName)

    const packages = await queryClient
      .fetchQuery<ProductPackage[]>({
        queryKey: queryKeys.products.productUnits(id),
        queryFn: () => api.post<ProductPackage[]>(`/products/${id}/packages/list`, {}),
      })
      .catch(() => [] as ProductPackage[])

    const validPackages = Array.isArray(packages) ? packages : []
    const defaultPkg = validPackages.find((pkg) => pkg.is_default) ?? validPackages[0]

    setItemPackageId((prev) => ({ ...prev, [index]: defaultPkg?.id ?? 0 }))
    setValue(`items.${index}.unit`, defaultPkg?.unit_name ?? 'pcs')
    setValue(`items.${index}.price`, defaultPkg?.purchase_price ?? 0)
    setValue(`items.${index}.conversion_qty`, defaultPkg?.resolved_factor ?? 1)
  }

  const newItemsTotal = watchItems.reduce((sum, item) => sum + (item.quantity || 0) * (item.price || 0), 0)
  const currentTotal = detail?.total_amount ?? 0
  const paidAmount = detail?.paid_amount ?? 0
  const projectedTotal = currentTotal + newItemsTotal
  const projectedRemaining = Math.max(0, projectedTotal - paidAmount)
  const projectedStatus = projectedRemaining <= 0 ? 'Lunas' : 'Sebagian'

  function handleClose() {
    reset({ items: [emptyItem] })
    setConfirmState(initialConfirmState)
    onOpenChange(false)
  }

  function onSubmit(values: AddItemsFormValues) {
    setConfirmState({ isConfirming: true, pendingValues: values })
  }

  function handleConfirmedSave() {
    const { pendingValues } = confirmState
    if (!pendingValues || !purchaseId) return
    addItems(
      {
        id: purchaseId,
        items: pendingValues.items.map((item, index) => ({
          product_id: item.product_id,
          package_id: itemPackageId[index] || undefined,
          quantity: item.quantity,
          purchase_price: item.price,
          unit: item.unit,
          conversion_qty: item.conversion_qty ?? 1,
        })),
      },
      {
        onSuccess: () => handleClose(),
        onError: () => setConfirmState((s) => ({ ...s, isConfirming: false })),
      },
    )
  }

  return (
    <>
      <FormModal
        open={open && !confirmState.isConfirming}
        onOpenChange={(v) => { if (!v && !confirmState.isConfirming) handleClose() }}
        title={`Tambah Item — ${purchase?.purchase_code ?? ''}`}
        size="full"
        isLoading={isPending}
        onSubmit={handleSubmit(onSubmit)}
        submitLabel="Simpan Item Baru"
      >
        {!detail ? (
          <div className="space-y-3 py-4">
            {[1, 2, 3].map((i) => (
              <div key={i} className="h-8 animate-pulse rounded-md bg-gray-100" />
            ))}
          </div>
        ) : (
          <div className="space-y-5">
            <div className="rounded-lg border overflow-hidden">
              <table className="w-full text-xs">
                <thead className="bg-gray-50 border-b">
                  <tr>
                    <th className="px-3 py-2 text-left font-medium text-gray-500">Produk (sudah ada)</th>
                    <th className="px-3 py-2 text-right font-medium text-gray-500">Qty</th>
                    <th className="px-3 py-2 text-left font-medium text-gray-500">Satuan</th>
                    <th className="px-3 py-2 text-right font-medium text-gray-500">Harga</th>
                    <th className="px-3 py-2 text-right font-medium text-gray-500">Subtotal</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-gray-100">
                  {detail.items.map((item) => (
                    <tr key={item.id} className="text-gray-500 bg-gray-50/50">
                      <td className="px-3 py-2">{item.product_name}</td>
                      <td className="px-3 py-2 text-right">{item.quantity}</td>
                      <td className="px-3 py-2">{item.unit}</td>
                      <td className="px-3 py-2 text-right">{formatRupiah(item.purchase_price)}</td>
                      <td className="px-3 py-2 text-right">{formatRupiah(item.subtotal)}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
              <p className="px-3 py-1.5 text-[11px] text-gray-400 bg-gray-50 border-t">
                Item di atas sudah tercatat &amp; tidak bisa diubah di sini. Tambahkan item baru di bawah.
              </p>
            </div>

            <div className="space-y-2">
              <div className="flex items-center justify-between">
                <p className="text-sm font-medium">Item Baru</p>
                <Button
                  type="button"
                  variant="outline"
                  size="sm"
                  onClick={() => append(emptyItem)}
                  className="h-7 gap-1 text-xs"
                >
                  <Plus className="h-3 w-3" />
                  Tambah Item
                </Button>
              </div>

              <div className="rounded-lg border overflow-x-auto">
                <table className="w-full min-w-[700px] text-sm">
                  <thead className="bg-gray-50 border-b">
                    <tr>
                      <th className="px-3 py-2 text-left font-medium text-gray-500 w-[30%]">Produk</th>
                      <th className="px-3 py-2 text-right font-medium text-gray-500 w-[10%]">Qty</th>
                      <th className="px-3 py-2 text-left font-medium text-gray-500 w-[12%]">Satuan</th>
                      <th className="px-3 py-2 text-right font-medium text-gray-500 w-[20%]">Harga</th>
                      <th className="px-3 py-2 text-right font-medium text-gray-500 w-[18%]">Subtotal</th>
                      <th className="px-3 py-2 w-8" />
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-gray-100">
                    {fields.map((field, index) => {
                      const qty = watchItems[index]?.quantity || 0
                      const price = watchItems[index]?.price || 0
                      const itemErrors = errors.items?.[index]
                      return (
                        <tr key={field.id}>
                          <td className="px-3 py-2">
                            <ProductCell
                              value={watchItems[index]?.product_id ?? 0}
                              label={watchItems[index]?.product_name}
                              onChange={(id, name) => handleProductChange(index, id, name)}
                              error={!!itemErrors?.product_id}
                            />
                            {itemErrors?.product_id && (
                              <p className="text-xs text-red-500 mt-0.5">{itemErrors.product_id.message}</p>
                            )}
                          </td>
                          <td className="px-3 py-2">
                            <Controller
                              control={control}
                              name={`items.${index}.quantity`}
                              render={({ field: f }) => (
                                <Input
                                  type="number"
                                  min={1}
                                  step={1}
                                  value={f.value}
                                  onChange={(e) => f.onChange(Number(e.target.value))}
                                  className={`h-8 text-xs text-right ${itemErrors?.quantity ? 'border-red-500' : ''}`}
                                />
                              )}
                            />
                          </td>
                          <td className="px-3 py-2 text-xs text-gray-700">
                            {watchItems[index]?.unit || '-'}
                          </td>
                          <td className="px-3 py-2">
                            <Controller
                              control={control}
                              name={`items.${index}.price`}
                              render={({ field: f }) => (
                                <RupiahInput
                                  value={f.value}
                                  onChange={f.onChange}
                                  className={`h-8 text-xs ${itemErrors?.price ? 'border-red-500' : ''}`}
                                />
                              )}
                            />
                          </td>
                          <td className="px-3 py-2 text-right text-xs font-medium">
                            {formatRupiah(qty * price)}
                          </td>
                          <td className="px-3 py-2">
                            {fields.length > 1 && (
                              <Button
                                type="button"
                                variant="ghost"
                                size="sm"
                                className="h-7 w-7 p-0"
                                onClick={() => remove(index)}
                              >
                                <Trash2 className="h-3.5 w-3.5 text-red-400" />
                              </Button>
                            )}
                          </td>
                        </tr>
                      )
                    })}
                  </tbody>
                </table>
              </div>
            </div>

            <div className="rounded-lg bg-gray-50 p-3 space-y-1 text-sm">
              <div className="flex justify-between">
                <span className="text-gray-500">Total Sekarang</span>
                <span>{formatRupiah(currentTotal)}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-gray-500">Tambahan Item Baru</span>
                <span>+{formatRupiah(newItemsTotal)}</span>
              </div>
              <div className="flex justify-between border-t pt-1 font-semibold">
                <span>Total Baru</span>
                <span>{formatRupiah(projectedTotal)}</span>
              </div>
              <div className="flex justify-between text-xs text-gray-500 pt-1">
                <span>Status pembayaran setelah disimpan</span>
                <span className="font-medium">{projectedStatus}</span>
              </div>
            </div>
          </div>
        )}
      </FormModal>

      <ConfirmDialog
        open={confirmState.isConfirming}
        onOpenChange={(v) => { if (!v) handleClose() }}
        title="Tambah Item Pembelian"
        description={`Total PO akan berubah dari ${formatRupiah(currentTotal)} menjadi ${formatRupiah(projectedTotal)}. Status pembayaran akan menjadi "${projectedStatus}". Lanjutkan?`}
        confirmLabel="Ya, Simpan"
        isLoading={isPending}
        onConfirm={handleConfirmedSave}
      />
    </>
  )
}

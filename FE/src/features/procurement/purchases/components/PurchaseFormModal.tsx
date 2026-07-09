import { useEffect, useState } from 'react'
import { useForm, useFieldArray, Controller } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { Plus, Trash2 } from 'lucide-react'

import { ConfirmDialog, FormModal } from '@/shared/components'
import { Input } from '@/shared/components/ui/input'
import { Label } from '@/shared/components/ui/label'
import { Textarea } from '@/shared/components/ui/textarea'
import { Button } from '@/shared/components/ui/button'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/shared/components/ui/select'
import { formatRupiah, todayStr } from '@/shared/utils'
import { RupiahInput } from '@/shared/components/ui/rupiah-input'
import { api } from '@/services'
import { useSupplierOptionsQuery } from '@/features/procurement/suppliers'
import type { SupplierOption } from '@/features/procurement/suppliers'
import { useProductSearchQuery } from '@/features/products/products'
import { AsyncCombobox } from '@/shared/components/ui/async-combobox'
import { useQueryClient } from '@tanstack/react-query'
import { queryKeys } from '@/shared/constants'
import {
  useCreateSupplierPurchaseMutation,
  useUpdateSupplierPurchaseMutation,
  useGeneratePurchaseCodeQuery,
} from '../purchases.api'
import { usePaymentStatusesQuery } from '../payment-statuses.api'
import { usePaymentMethodsQuery } from '../payment-methods.api'
import type { PaymentStatus, SupplierPurchase } from '../purchases.types'
import type { ProductPackage, ProductSearchOption } from '@/features/products/products/products.types'
import { purchaseSchema, type PurchaseFormValues } from '../purchases.schema'

interface PurchaseFormModalProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  initialData?: SupplierPurchase | null
}

function PurchaseProductCell({
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

const emptyValues: PurchaseFormValues = {
  purchase_date: todayStr(),
  invoice_number: '',
  supplier_id: 0,
  items: [{ product_id: 0, product_name: '', quantity: 1, price: 0, unit: '', conversion_qty: 1 }],
  discount_amount: 0,
  notes: '',
  payment_status: 'paid',
  paid_amount: 0,
  payment_method: 'cash',
}

function buildDefaultValues(data: SupplierPurchase): PurchaseFormValues {
  return {
    purchase_date: data.purchase_date,
    invoice_number: data.invoice_number,
    supplier_id: data.supplier_id,
    items: data.items.map((item) => ({
      product_id: item.product_id,
      product_name: item.product_name,
      quantity: item.quantity,
      price: item.purchase_price,
      unit: item.unit,
      conversion_qty: item.conversion_qty,
    })),
    discount_amount: data.discount_amount,
    notes: data.notes ?? '',
    payment_status: data.payment_status,
    paid_amount: data.paid_amount,
    payment_method: 'cash',
  }
}

export function PurchaseFormModal({ open, onOpenChange, initialData }: PurchaseFormModalProps) {
  const isEditMode = !!initialData

  const [isConfirming, setIsConfirming] = useState(false)
  const [pendingValues, setPendingValues] = useState<PurchaseFormValues | null>(null)

  const [supplierKeyword, setSupplierKeyword] = useState('')
  const [selectedSupplierLabel, setSelectedSupplierLabel] = useState('')
  const { data: supplierOptions = [], isFetching: isSuppliersFetching } = useSupplierOptionsQuery(supplierKeyword)

  const queryClient = useQueryClient()
  const { mutate: create, isPending: isCreating } = useCreateSupplierPurchaseMutation()
  const { mutate: update, isPending: isUpdating } = useUpdateSupplierPurchaseMutation()
  const isPending = isCreating || isUpdating

  const { data: codeData, isFetching: isGeneratingCode } = useGeneratePurchaseCodeQuery(
    open && !isEditMode
  )
  const { data: paymentStatuses = [] } = usePaymentStatusesQuery()
  const { data: paymentMethods = [] } = usePaymentMethodsQuery()

  const {
    register,
    handleSubmit,
    control,
    watch,
    setValue,
    reset,
    formState: { errors },
  } = useForm<PurchaseFormValues>({
    resolver: zodResolver(purchaseSchema),
    defaultValues: emptyValues,
  })

  const { fields, append, remove } = useFieldArray({ control, name: 'items' })

  const [itemUnitOptions, setItemUnitOptions] = useState<Record<number, ProductPackage[]>>({})
  const [itemSelectedPackageId, setItemSelectedPackageId] = useState<Record<number, number>>({})
  const [itemRefPurchasePrice, setItemRefPurchasePrice] = useState<Record<number, number>>({})

  const watchItems = watch('items')
  const watchDiscount = watch('discount_amount') ?? 0
  const watchPaymentStatus = watch('payment_status')
  const watchPaymentMethod = watch('payment_method')
  const watchSupplierId = watch('supplier_id')

  const subtotal = watchItems.reduce((sum, item) => sum + (item.quantity || 0) * (item.price || 0), 0)
  const total = Math.max(0, subtotal - (watchDiscount || 0))

  useEffect(() => {
    if (!open) {
      reset({ ...emptyValues, purchase_date: todayStr() })
      setItemUnitOptions({})
      setItemSelectedPackageId({})
      setItemRefPurchasePrice({})
      setIsConfirming(false)
      setPendingValues(null)
      setSupplierKeyword('')
      setSelectedSupplierLabel('')
      return
    }
    if (isEditMode && initialData) {
      reset(buildDefaultValues(initialData))
      setItemUnitOptions({})
      setItemSelectedPackageId({})
      setSelectedSupplierLabel(initialData.supplier_name)
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [open, initialData])

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

    setItemRefPurchasePrice((prev) => ({ ...prev, [index]: defaultPkg?.purchase_price ?? 0 }))

    if (validPackages.length > 1) {
      setItemUnitOptions((prev) => ({ ...prev, [index]: validPackages }))
      setItemSelectedPackageId((prev) => ({ ...prev, [index]: defaultPkg?.id ?? 0 }))
      setValue(`items.${index}.unit`, defaultPkg?.unit_name ?? 'pcs')
      setValue(`items.${index}.price`, defaultPkg?.purchase_price ?? 0)
      setValue(`items.${index}.conversion_qty`, defaultPkg?.conversion_qty ?? 1)
    } else {
      setItemUnitOptions((prev) => {
        const next = { ...prev }
        delete next[index]
        return next
      })
      setItemSelectedPackageId((prev) => {
        const next = { ...prev }
        delete next[index]
        return next
      })
      setValue(`items.${index}.unit`, defaultPkg?.unit_name ?? 'pcs')
      setValue(`items.${index}.price`, defaultPkg?.purchase_price ?? 0)
      setValue(`items.${index}.conversion_qty`, 1)
    }
  }

  function handleUnitChange(index: number, packageId: string) {
    const id = Number(packageId)
    const pkg = itemUnitOptions[index]?.find((p) => p.id === id)
    if (!pkg) return
    setItemSelectedPackageId((prev) => ({ ...prev, [index]: id }))
    setItemRefPurchasePrice((prev) => ({ ...prev, [index]: pkg.purchase_price }))
    setValue(`items.${index}.unit`, pkg.unit_name)
    setValue(`items.${index}.price`, pkg.purchase_price)
    setValue(`items.${index}.conversion_qty`, pkg.conversion_qty ?? 1)
  }

  function unitOptionLabel(pkg: ProductPackage) {
    return pkg.conversion_qty > 1
      ? `${pkg.package_name || pkg.unit_name} (x${pkg.conversion_qty})`
      : pkg.package_name || pkg.unit_name
  }

  const handleClose = () => {
    setIsConfirming(false)
    setPendingValues(null)
    onOpenChange(false)
  }

  function onSubmit(values: PurchaseFormValues) {
    setPendingValues(values)
    setIsConfirming(true)
  }

  function handleConfirmedSave() {
    if (!pendingValues) return

    const payload = {
      ...pendingValues,
      items: pendingValues.items.map((item) => ({
        product_id: item.product_id,
        quantity: item.quantity,
        purchase_price: item.price,
        unit: item.unit,
        conversion_qty: item.conversion_qty ?? 1,
      })),
      paid_amount: pendingValues.payment_status === 'paid' ? total : pendingValues.paid_amount,
    }

    if (isEditMode && initialData) {
      update(
        { id: initialData.id, ...payload },
        {
          onSuccess: () => handleClose(),
          onError: () => setIsConfirming(false),
        },
      )
    } else {
      create(payload, {
        onSuccess: () => handleClose(),
        onError: () => setIsConfirming(false),
      })
    }
  }

  return (
    <>
    <FormModal
      open={open && !isConfirming}
      onOpenChange={(open) => { if (!open && !isConfirming) handleClose() }}
      title={isEditMode ? 'Edit Pembelian' : 'Tambah Pembelian'}
      size="full"
      isLoading={isPending}
      onSubmit={handleSubmit(onSubmit)}
      submitLabel={isEditMode ? 'Simpan Perubahan' : 'Simpan Pembelian'}
    >
      <div className="space-y-5">
        <div className="grid grid-cols-4 gap-3">
          <div className="space-y-1.5">
            <Label htmlFor="pur-code">Kode PO</Label>
            <Input
              id="pur-code"
              value={
                isEditMode
                  ? (initialData?.purchase_code ?? '')
                  : isGeneratingCode
                    ? '...'
                    : (codeData?.purchase_code ?? '')
              }
              readOnly
              className="bg-gray-50 font-mono text-blue-700 font-medium"
            />
          </div>
          <div className="space-y-1.5">
            <Label htmlFor="pur-date">
              Tanggal <span className="text-red-500">*</span>
            </Label>
            <Input
              id="pur-date"
              type="date"
              max={todayStr()}
              {...register('purchase_date')}
              className={errors.purchase_date ? 'border-red-500' : ''}
            />
            {errors.purchase_date && (
              <p className="text-xs text-red-500">{errors.purchase_date.message}</p>
            )}
          </div>
          <div className="space-y-1.5">
            <Label htmlFor="pur-inv">
              No. Faktur <span className="text-red-500">*</span>
            </Label>
            <Input
              id="pur-inv"
              placeholder="INV-001"
              {...register('invoice_number')}
              className={errors.invoice_number ? 'border-red-500' : ''}
            />
            {errors.invoice_number && (
              <p className="text-xs text-red-500">{errors.invoice_number.message}</p>
            )}
          </div>
          <div className="space-y-1.5">
            <Label>
              Supplier <span className="text-red-500">*</span>
            </Label>
            <AsyncCombobox<SupplierOption>
              className={errors.supplier_id ? 'border-red-500' : ''}
              value={watchSupplierId > 0 ? watchSupplierId : undefined}
              selectedLabel={selectedSupplierLabel}
              onSearch={setSupplierKeyword}
              onValueChange={(v, item) => {
                setValue('supplier_id', Number(v ?? 0), { shouldValidate: true })
                if (item) setSelectedSupplierLabel(item.name)
              }}
              options={supplierOptions}
              getOptionValue={(s) => s.id}
              getOptionLabel={(s) => s.name}
              isLoading={isSuppliersFetching}
              placeholder="Pilih supplier"
              searchPlaceholder="Cari supplier..."
              emptyText="Supplier tidak ditemukan."
            />
            {errors.supplier_id && (
              <p className="text-xs text-red-500">{errors.supplier_id.message}</p>
            )}
          </div>
        </div>

        <div className="space-y-2">
          <div className="flex items-center justify-between">
            <Label>
              Item Produk <span className="text-red-500">*</span>
            </Label>
            <Button
              type="button"
              variant="outline"
              size="sm"
              onClick={() => append({ product_id: 0, product_name: '', quantity: 1, price: 0, unit: '', conversion_qty: 1 })}
              className="h-7 gap-1 text-xs"
            >
              <Plus className="h-3 w-3" />
              Tambah Item
            </Button>
          </div>

          <div className="rounded-lg border overflow-x-auto">
            <table className="w-full min-w-[820px] text-sm">
              <thead className="bg-gray-50 border-b">
                <tr>
                  <th className="px-3 py-2 text-left font-medium text-gray-500 min-w-[220px] w-[28%]">Produk</th>
                  <th className="px-3 py-2 text-right font-medium text-gray-500 min-w-[70px] w-[8%]">Qty</th>
                  <th className="px-3 py-2 text-left font-medium text-gray-500 min-w-[110px] w-[9%]">Satuan</th>
                  <th className="px-3 py-2 text-right font-medium text-gray-500 min-w-[120px] w-[14%]">Ref. Harga Beli</th>
                  <th className="px-3 py-2 text-right font-medium text-gray-500 min-w-[140px] w-[19%]">Harga</th>
                  <th className="px-3 py-2 text-right font-medium text-gray-500 min-w-[110px] w-[16%]">Subtotal</th>
                  <th className="px-3 py-2 w-8" />
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-100">
                {fields.map((field, index) => {
                  const qty = watchItems[index]?.quantity || 0
                  const price = watchItems[index]?.price || 0
                  const currentProductId = watchItems[index]?.product_id
                  const itemErrors = errors.items?.[index]
                  return (
                    <tr key={field.id}>
                      <td className="px-3 py-2">
                        <PurchaseProductCell
                          value={currentProductId ?? 0}
                          label={watchItems[index]?.product_name}
                          onChange={(id, name) => handleProductChange(index, id, name)}
                          error={!!itemErrors?.product_id}
                        />
                        {itemErrors?.product_id && (
                          <p className="text-xs text-red-500 mt-0.5">{itemErrors.product_id.message}</p>
                        )}
                      </td>
                      <td className="px-3 py-2">
                        <Input
                          type="number"
                          min={1}
                          step={1}
                          {...register(`items.${index}.quantity`, { valueAsNumber: true })}
                          className={`h-8 text-xs text-right ${itemErrors?.quantity ? 'border-red-500' : ''}`}
                        />
                        {itemErrors?.quantity && (
                          <p className="text-xs text-red-500 mt-0.5">{itemErrors.quantity.message}</p>
                        )}
                      </td>
                      <td className="px-3 py-2 text-xs text-gray-700">
                        {itemUnitOptions[index]?.length > 1 ? (
                          <Select
                            value={String(itemSelectedPackageId[index] ?? '')}
                            onValueChange={(v) => handleUnitChange(index, v)}
                          >
                            <SelectTrigger className="h-8 text-xs min-w-[120px]">
                              <SelectValue placeholder="-" />
                            </SelectTrigger>
                            <SelectContent>
                              {itemUnitOptions[index].map((pkg) => (
                                <SelectItem key={pkg.id} value={String(pkg.id)}>
                                  {unitOptionLabel(pkg)}
                                </SelectItem>
                              ))}
                            </SelectContent>
                          </Select>
                        ) : (
                          watchItems[index]?.unit || '-'
                        )}
                      </td>
                      <td className="px-3 py-2 text-right text-xs text-gray-400">
                        {itemRefPurchasePrice[index] != null
                          ? formatRupiah(itemRefPurchasePrice[index])
                          : '-'}
                      </td>
                      <td className="px-3 py-2">
                        <Controller
                          control={control}
                          name={`items.${index}.price`}
                          render={({ field }) => (
                            <RupiahInput
                              value={field.value}
                              onChange={field.onChange}
                              className={`h-8 text-xs ${itemErrors?.price ? 'border-red-500' : ''}`}
                            />
                          )}
                        />
                        {itemErrors?.price && (
                          <p className="text-xs text-red-500 mt-0.5">{itemErrors.price.message}</p>
                        )}
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
          {errors.items?.root && (
            <p className="text-xs text-red-500">{errors.items.root.message}</p>
          )}
          {errors.items?.message && (
            <p className="text-xs text-red-500">{errors.items.message}</p>
          )}
        </div>

        <div className="grid grid-cols-2 gap-4">
          <div className="space-y-3">
            <div className="space-y-1.5">
              <Label htmlFor="pur-discount">Diskon (Rp)</Label>
              <Controller
                control={control}
                name="discount_amount"
                render={({ field }) => (
                  <RupiahInput
                    id="pur-discount"
                    value={field.value ?? 0}
                    onChange={field.onChange}
                    className={errors.discount_amount ? 'border-red-500' : ''}
                  />
                )}
              />
              {errors.discount_amount && (
                <p className="text-xs text-red-500">{errors.discount_amount.message}</p>
              )}
            </div>

            <div className="space-y-1.5">
              <Label>Status Pembayaran</Label>
              <Select
                value={watchPaymentStatus}
                onValueChange={(v) => setValue('payment_status', v as PaymentStatus, { shouldValidate: true })}
              >
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  {paymentStatuses.map((s) => (
                    <SelectItem key={s.code} value={s.code}>
                      {s.label}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>

            {watchPaymentStatus === 'partial' && (
              <div className="space-y-1.5">
                <Label htmlFor="pur-paid">Jumlah Dibayar (Rp)</Label>
                <Controller
                  control={control}
                  name="paid_amount"
                  render={({ field }) => (
                    <RupiahInput
                      id="pur-paid"
                      value={field.value}
                      onChange={(v) => field.onChange(v)}
                      className={errors.paid_amount ? 'border-red-500' : ''}
                    />
                  )}
                />
                {errors.paid_amount && (
                  <p className="text-xs text-red-500">{errors.paid_amount.message}</p>
                )}
              </div>
            )}

            {watchPaymentStatus !== 'unpaid' && (
              <div className="space-y-1.5">
                <Label>Metode Pembayaran</Label>
                <Select
                  value={watchPaymentMethod ?? 'cash'}
                  onValueChange={(v) => setValue('payment_method', v, { shouldValidate: true })}
                >
                  <SelectTrigger className={errors.payment_method ? 'border-red-500' : ''}>
                    <SelectValue placeholder="Pilih metode" />
                  </SelectTrigger>
                  <SelectContent>
                    {paymentMethods.map((m) => (
                      <SelectItem key={m.code} value={m.code}>
                        {m.label}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
                {errors.payment_method && (
                  <p className="text-xs text-red-500">{errors.payment_method.message}</p>
                )}
              </div>
            )}
          </div>

          <div className="space-y-3">
            <div className="rounded-lg bg-gray-50 p-3 space-y-1 text-sm">
              <div className="flex justify-between">
                <span className="text-gray-500">Subtotal</span>
                <span>{formatRupiah(subtotal)}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-gray-500">Diskon</span>
                <span className="text-red-500">-{formatRupiah(watchDiscount || 0)}</span>
              </div>
              <div className="flex justify-between border-t pt-1 font-semibold">
                <span>Total</span>
                <span>{formatRupiah(total)}</span>
              </div>
            </div>

            <div className="space-y-1.5">
              <Label htmlFor="pur-notes">Catatan</Label>
              <Textarea
                id="pur-notes"
                {...register('notes')}
                placeholder="Catatan pembelian (opsional)"
                className="resize-none"
                rows={3}
              />
            </div>
          </div>
        </div>
      </div>
    </FormModal>

    <ConfirmDialog
      open={isConfirming}
      onOpenChange={(val) => { if (!val) handleClose() }}
      title={isEditMode ? 'Update Pembelian' : 'Tambah Pembelian'}
      description={`Yakin ingin ${isEditMode ? 'memperbarui' : 'menyimpan'} data pembelian ini?`}
      confirmLabel="Ya, Simpan"
      isLoading={isPending}
      onConfirm={handleConfirmedSave}
    />
    </>
  )
}

import { Fragment, useEffect, useState } from 'react'
import { useForm, useFieldArray, useWatch, Controller, type Control } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { Info, Lock, Plus, Trash2, Unlock } from 'lucide-react'

import { ConfirmDialog, FormModal } from '@/shared/components'
import { Checkbox } from '@/shared/components/ui/checkbox'
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
import { formatRupiah, todayStr, toISODate } from '@/shared/utils'
import { RupiahInput } from '@/shared/components/ui/rupiah-input'
import { api } from '@/services'
import { useSupplierOptionsQuery } from '@/features/procurement/suppliers'
import type { SupplierOption } from '@/features/procurement/suppliers'
import { useProductSearchQuery, fetchProductDetail } from '@/features/products/products'
import { formatResolvedFactor } from '@/features/products/products/products.utils'
import { AsyncCombobox } from '@/shared/components/ui/async-combobox'
import { useQueryClient } from '@tanstack/react-query'
import { queryKeys } from '@/shared/constants'
import {
  useCreateSupplierPurchaseMutation,
  useUpdateSupplierPurchaseMutation,
  useGeneratePurchaseCodeQuery,
  useSupplierPurchaseDetailQuery,
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

// ExpiryBatchSection: rincian tanggal expired opsional untuk 1 baris item PO. Checkbox
// dan repeater-nya sengaja disatukan di sini (bukan di komponen utama) karena useFieldArray
// per baris cuma bisa dipanggil sebagai hook terpisah per index — tidak bisa dipanggil
// dinamis di dalam .map() komponen induk.
function ExpiryBatchSection({
  control,
  index,
  itemQuantity,
}: {
  control: Control<PurchaseFormValues>
  index: number
  itemQuantity: number
}) {
  const { fields, append, remove, replace } = useFieldArray({
    control,
    name: `items.${index}.expiry_batches`,
  })
  const watchBatches = useWatch({ control, name: `items.${index}.expiry_batches` }) ?? []
  const enabled = fields.length > 0
  const totalAllocated = watchBatches.reduce((sum, b) => sum + (b?.qty || 0), 0)
  const isBalanced = enabled && Math.abs(totalAllocated - itemQuantity) < 0.0001

  function handleToggle(checked: boolean) {
    if (checked) {
      append({ qty: itemQuantity || 0, expired_date: '' })
    } else {
      replace([])
    }
  }

  return (
    <div className="space-y-2">
      <div className="flex items-center gap-2">
        <Checkbox
          id={`expiry-toggle-${index}`}
          checked={enabled}
          onCheckedChange={(v) => handleToggle(v === true)}
        />
        <label htmlFor={`expiry-toggle-${index}`} className="text-xs font-medium text-gray-700 cursor-pointer">
          Produk ini ada tanggal expired
        </label>
      </div>

      {enabled && (
        <div className="rounded-lg border border-dashed border-indigo-200 bg-indigo-50/60 p-3">
          <div className="flex items-center justify-between mb-2">
            <span className="text-xs font-semibold text-gray-700">Rincian Tanggal Expired</span>
            <span className="text-[11px] text-gray-500">Total qty harus sama dengan {itemQuantity}</span>
          </div>

          <div className="space-y-1.5">
            {fields.map((f, bi) => (
              <div key={f.id} className="grid grid-cols-[110px_1fr_auto] gap-2 items-center rounded-md border bg-white p-1.5">
                <Controller
                  control={control}
                  name={`items.${index}.expiry_batches.${bi}.qty`}
                  render={({ field }) => (
                    <Input
                      type="number"
                      min={0.001}
                      step="any"
                      value={field.value ?? ''}
                      onChange={(e) => field.onChange(e.target.valueAsNumber || 0)}
                      className="h-7 text-xs"
                      placeholder="Qty"
                    />
                  )}
                />
                <Controller
                  control={control}
                  name={`items.${index}.expiry_batches.${bi}.expired_date`}
                  render={({ field }) => (
                    <Input
                      type="date"
                      value={field.value ?? ''}
                      onChange={field.onChange}
                      className="h-7 text-xs"
                    />
                  )}
                />
                <Button
                  type="button"
                  variant="ghost"
                  size="sm"
                  className="h-7 w-7 p-0"
                  onClick={() => remove(bi)}
                >
                  <Trash2 className="h-3.5 w-3.5 text-red-400" />
                </Button>
              </div>
            ))}
          </div>

          <Button
            type="button"
            variant="outline"
            size="sm"
            className="mt-2 h-7 gap-1 text-xs"
            onClick={() => append({ qty: 0, expired_date: '' })}
          >
            <Plus className="h-3 w-3" />
            Tambah Tanggal Expired
          </Button>

          <div
            className={`mt-2 flex justify-between rounded-md border px-2.5 py-1.5 text-xs font-medium ${
              isBalanced
                ? 'border-green-200 bg-green-50 text-green-700'
                : 'border-red-200 bg-red-50 text-red-600'
            }`}
          >
            <span>Qty teralokasi</span>
            <span>
              {totalAllocated} / {itemQuantity} {isBalanced ? '✓' : ''}
            </span>
          </div>
        </div>
      )}
    </div>
  )
}

const emptyValues: PurchaseFormValues = {
  purchase_date: todayStr(),
  invoice_number: '',
  no_invoice: false,
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
    // API mengembalikan datetime ISO penuh (mis. "2026-07-10T00:00:00+07:00"), tapi
    // <input type="date"> butuh format yyyy-MM-dd persis, kalau tidak inputnya kosong.
    purchase_date: toISODate(new Date(data.purchase_date)),
    invoice_number: data.invoice_number,
    no_invoice: data.invoice_number.startsWith('AUTO-'),
    supplier_id: data.supplier_id,
    items: data.items.map((item) => ({
      product_id: item.product_id,
      product_name: item.product_name,
      quantity: item.quantity,
      price: item.purchase_price,
      unit: item.unit,
      conversion_qty: item.conversion_qty,
      // Wajib dibawa balik ke form — Update() menghapus semua purchase_items lama
      // (CASCADE ke product_expiry_batches) lalu re-insert cuma dari payload ini. Kalau
      // expiry_batches tidak di-hydrate di sini, submit edit apapun (bahkan tanpa
      // perubahan) akan menghapus permanen rincian tanggal expired yang sudah ada.
      expiry_batches: item.expiry_batches?.map((b) => ({ qty: b.qty, expired_date: b.expired_date })),
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
  // Baris di tabel/list pembelian tidak menyertakan `items` (hanya field ringkasan),
  // jadi saat edit kita perlu fetch detail lengkap dulu — memakai initialData langsung
  // sebagai defaultValues akan crash karena item.items undefined.
  const { data: fullPurchaseDetail } = useSupplierPurchaseDetailQuery(
    open && isEditMode ? (initialData?.id ?? null) : null
  )
  const { data: paymentStatuses = [] } = usePaymentStatusesQuery()
  const { data: paymentMethods = [] } = usePaymentMethodsQuery()

  const {
    register,
    handleSubmit,
    control,
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

  // Snapshot item ASLI (sebelum diedit) + stok produk saat ini — diambil sekali saat
  // modal Edit dibuka, dipakai buat validasi "produk yang dihapus/diganti dari daftar
  // tidak boleh bikin stoknya minus" (lihat purchase_repo.go Update() — aturan yang
  // sama persis, cuma di sini dicek proaktif sebelum user sempat klik Simpan).
  const [originalItemsSnapshot, setOriginalItemsSnapshot] = useState<
    { product_id: number; product_name: string; quantity: number; conversion_qty: number }[]
  >([])
  const [originalStockMap, setOriginalStockMap] = useState<Record<number, number>>({})

  // Menandai PO mana yang UI-state-nya (itemUnitOptions dkk) sudah disinkronkan, supaya
  // reset di bawah hanya terjadi sekali per data baru — pola "adjust state saat render"
  // yang direkomendasikan React untuk sinkron ke data eksternal, bukan react ke event user.
  const [loadedDetailId, setLoadedDetailId] = useState<number | null>(null)

  const watchItems = useWatch({ control, name: 'items' })
  const watchDiscount = useWatch({ control, name: 'discount_amount' }) ?? 0
  const watchPaymentStatus = useWatch({ control, name: 'payment_status' })
  const watchPaymentMethod = useWatch({ control, name: 'payment_method' })
  const watchSupplierId = useWatch({ control, name: 'supplier_id' })
  const watchNoInvoice = useWatch({ control, name: 'no_invoice' }) ?? false
  const watchPaidAmount = useWatch({ control, name: 'paid_amount' }) ?? 0

  const subtotal = watchItems.reduce((sum, item) => sum + (item.quantity || 0) * (item.price || 0), 0)
  const total = Math.max(0, subtotal - (watchDiscount || 0))

  // originalPaymentStatus: status pembayaran PO SEBELUM diedit (bukan pilihan di form,
  // karena field Status Pembayaran dikunci begitu status aslinya bukan "unpaid" — lihat
  // rancangan tier di bawah). undefined saat mode Tambah baru (semua bebas).
  const originalPaymentStatus = isEditMode ? initialData?.payment_status : undefined
  const originalPaidAmount = initialData?.paid_amount ?? 0
  // Tier kunci field: "unpaid" bebas semua, "partial" kunci Diskon+Status Pembayaran,
  // "paid" kunci itu + Metode Pembayaran juga (lihat diskusi rancangan form Edit).
  const lockDiscountAndStatus = isEditMode && originalPaymentStatus !== undefined && originalPaymentStatus !== 'unpaid'
  const lockPaymentMethod = isEditMode && originalPaymentStatus === 'paid'

  // Validasi live #1: PO yang ASALNYA Lunas — paid_amount aslinya tetap jadi acuan
  // (bukan field yang bisa diedit user, lihat handleConfirmedSave), jadi total baru
  // tidak boleh turun sampai di bawah itu. Untuk tier "partial", pengecekan setara
  // sudah ditangani oleh superRefine di purchaseSchema (field Jumlah Dibayar-nya
  // sendiri, yang masih terbuka).
  const totalBelowOriginalPaid = originalPaymentStatus === 'paid' && total < originalPaidAmount

  // Validasi live #2: produk yang ada di item ASLI tapi sudah tidak ada lagi di item
  // yang sedang diedit sekarang (dihapus atau diganti ke produk lain) — cek apakah
  // rollback qty pembelian aslinya bikin stok produk itu minus. Kalau minus, berarti
  // sebagian stok dari pembelian ini sudah terjual/keluar lewat jalur lain (mis. Kasir),
  // jadi tidak bisa "ditarik balik" lagi. Aturan yang sama persis dicek ulang di
  // purchase_repo.go Update() — di sini cuma supaya user tahu sebelum klik Simpan.
  const currentProductIds = new Set(watchItems.map((item) => item.product_id).filter((id) => id > 0))
  const removedProductWarnings = originalItemsSnapshot
    .filter((old) => !currentProductIds.has(old.product_id))
    .map((old) => {
      const currentStock = originalStockMap[old.product_id] ?? Infinity
      const neededRollback = old.quantity * old.conversion_qty
      return { ...old, currentStock, neededRollback, blocked: currentStock - neededRollback < 0 }
    })
    .filter((w) => w.blocked)

  const hasBlockingValidation = totalBelowOriginalPaid || removedProductWarnings.length > 0

  // effectivePaidAmount: paid_amount yang SUNGGUH akan dikirim ke backend — logika ini
  // dipakai bareng di payload submit (handleConfirmedSave) dan di ringkasan dialog
  // konfirmasi, supaya keduanya selalu konsisten (satu sumber kebenaran).
  const effectivePaidAmount =
    originalPaymentStatus === 'paid' ? originalPaidAmount : watchPaymentStatus === 'paid' ? total : watchPaidAmount
  const effectiveRemaining = Math.max(0, total - effectivePaidAmount)

  // Kode PO dipakai sebagai basis nomor faktur otomatis saat toggle "faktur tidak
  // ada" aktif — di mode tambah dari hasil generate-code, di mode edit dari PO yang
  // sudah tersimpan.
  const currentPurchaseCode = isEditMode
    ? (initialData?.purchase_code ?? '')
    : (codeData?.purchase_code ?? '')

  // Sinkronkan invoice_number ke placeholder otomatis begitu kode PO tersedia — perlu
  // effect terpisah karena saat toggle pertama kali diklik, kode PO (async) mungkin
  // belum sempat ter-fetch.
  useEffect(() => {
    if (watchNoInvoice && currentPurchaseCode) {
      setValue('invoice_number', `AUTO-${currentPurchaseCode}`, { shouldValidate: true })
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [watchNoInvoice, currentPurchaseCode])

  function handleToggleNoInvoice() {
    const next = !watchNoInvoice
    setValue('no_invoice', next)
    setValue('invoice_number', next && currentPurchaseCode ? `AUTO-${currentPurchaseCode}` : '', {
      shouldValidate: true,
    })
  }

  if (open && isEditMode && fullPurchaseDetail && fullPurchaseDetail.id !== loadedDetailId) {
    setLoadedDetailId(fullPurchaseDetail.id)
    setItemUnitOptions({})
    setItemSelectedPackageId({})
    setSelectedSupplierLabel(fullPurchaseDetail.supplier_name)

    const snapshot = fullPurchaseDetail.items.map((item) => ({
      product_id: item.product_id,
      product_name: item.product_name,
      quantity: item.quantity,
      conversion_qty: item.conversion_qty || 1,
    }))
    setOriginalItemsSnapshot(snapshot)
    setOriginalStockMap({})
    Promise.all(
      snapshot.map((item) =>
        fetchProductDetail(item.product_id)
          .then((p) => [item.product_id, p.stock] as const)
          // Gagal fetch (mis. produk sudah dihapus) — anggap stok "aman tak terbatas" di
          // sisi FE; backend tetap jadi penjaga akhir yang sebenarnya kalau ini salah.
          .catch(() => [item.product_id, Infinity] as const)
      )
    ).then((pairs) => setOriginalStockMap(Object.fromEntries(pairs)))
  }

  // reset() react-hook-form murni imperatif (bukan React state setter), jadi tetap aman
  // dipanggil di effect — hanya dipicu sekali per kedatangan data fetch yang baru.
  useEffect(() => {
    if (open && isEditMode && fullPurchaseDetail) {
      reset(buildDefaultValues(fullPurchaseDetail))
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [open, isEditMode, fullPurchaseDetail])

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
      setValue(`items.${index}.conversion_qty`, defaultPkg?.resolved_factor ?? 1)
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
    setValue(`items.${index}.conversion_qty`, pkg.resolved_factor ?? 1)
  }

  function unitOptionLabel(pkg: ProductPackage) {
    const name = pkg.package_name ? `${pkg.unit_name} (${pkg.package_name})` : pkg.unit_name
    return pkg.resolved_factor > 1 ? `${name} (x${formatResolvedFactor(pkg.resolved_factor)})` : name
  }

  const handleClose = () => {
    reset({ ...emptyValues, purchase_date: todayStr() })
    setItemUnitOptions({})
    setItemSelectedPackageId({})
    setItemRefPurchasePrice({})
    setOriginalItemsSnapshot([])
    setOriginalStockMap({})
    setIsConfirming(false)
    setPendingValues(null)
    setSupplierKeyword('')
    setSelectedSupplierLabel('')
    setLoadedDetailId(null)
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
      // no_invoice cuma state UI form (toggle "faktur tidak ada"), bukan field yang
      // dikenal backend — di-set undefined supaya tidak ikut terserialisasi di payload.
      no_invoice: undefined,
      items: pendingValues.items.map((item, index) => ({
        product_id: item.product_id,
        package_id: itemSelectedPackageId[index] || undefined,
        quantity: item.quantity,
        purchase_price: item.price,
        unit: item.unit,
        conversion_qty: item.conversion_qty ?? 1,
        expiry_batches: item.expiry_batches && item.expiry_batches.length > 0 ? item.expiry_batches : undefined,
      })),
      // paid_amount: kalau PO ini ASALNYA sudah Lunas, jangan ikut menimpanya jadi
      // total baru (Status Pembayaran & field ini terkunci di tier itu) — pertahankan
      // angka aslinya apa adanya, biar backend yang hitung ulang remaining_amount &
      // otomatis turunkan status jadi "Bayar Sebagian" kalau total naik. Untuk tier
      // "partial", field Jumlah Dibayar masih terbuka jadi tetap dari input user. Untuk
      // PO yang asalnya "unpaid" (atau mode Tambah baru), Status Pembayaran sendiri
      // masih bebas dipilih user, jadi perilaku lama (paid_amount = total kalau pilih
      // Lunas) tetap berlaku. effectivePaidAmount sudah menghitung ketiga kasus ini
      // (dipakai juga di ringkasan dialog konfirmasi, satu sumber kebenaran).
      paid_amount: effectivePaidAmount,
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
      submitDisabled={hasBlockingValidation}
      onSubmit={handleSubmit(onSubmit)}
      submitLabel={isEditMode ? 'Simpan Perubahan' : 'Simpan Pembelian'}
    >
      {isEditMode && !fullPurchaseDetail ? (
        <div className="space-y-3 py-4">
          {[1, 2, 3, 4].map((i) => (
            <div key={i} className="h-8 animate-pulse rounded-md bg-gray-100" />
          ))}
        </div>
      ) : (
      <div className="space-y-5">
        {lockDiscountAndStatus && (
          <div className="flex items-start gap-2 rounded-lg border border-blue-200 bg-blue-50 px-3 py-2 text-xs text-blue-700">
            <Info size={14} className="mt-0.5 shrink-0" />
            <span>
              {originalPaymentStatus === 'paid'
                ? 'PO ini sudah lunas — Diskon, Status Pembayaran, dan Metode Pembayaran tidak bisa diubah dari sini.'
                : `PO ini sudah dibayar sebagian (${formatRupiah(originalPaidAmount)} dari ${formatRupiah(initialData?.total_amount ?? 0)}) — Diskon dan Status Pembayaran tidak bisa diubah dari sini.`}
            </span>
          </div>
        )}

        {removedProductWarnings.map((w) => (
          <div
            key={w.product_id}
            className="flex items-start gap-2 rounded-lg border border-red-200 bg-red-50 px-3 py-2 text-xs text-red-700"
          >
            <Info size={14} className="mt-0.5 shrink-0" />
            <span>
              Produk &quot;{w.product_name}&quot; tidak bisa dihapus/diganti dari pembelian ini — sebagian stoknya
              sudah terjual (sisa stok {w.currentStock}, butuh rollback {w.neededRollback}). Kembalikan produk ini
              ke daftar item untuk melanjutkan.
            </span>
          </div>
        ))}

        {totalBelowOriginalPaid && (
          <div className="flex items-start gap-2 rounded-lg border border-red-200 bg-red-50 px-3 py-2 text-xs text-red-700">
            <Info size={14} className="mt-0.5 shrink-0" />
            <span>
              Total baru ({formatRupiah(total)}) tidak boleh lebih kecil dari yang sudah dibayar (
              {formatRupiah(originalPaidAmount)}). Tambah kembali item yang dikurangi, atau batalkan perubahan ini.
            </span>
          </div>
        )}

        <div className="grid grid-cols-1 gap-3 sm:grid-cols-2 lg:grid-cols-4">
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
              No. Faktur {!watchNoInvoice && <span className="text-red-500">*</span>}
            </Label>
            <div className="relative">
              <Input
                id="pur-inv"
                placeholder={watchNoInvoice ? '' : 'INV-001'}
                readOnly={watchNoInvoice}
                {...register('invoice_number')}
                className={`pr-9 ${watchNoInvoice ? 'bg-gray-50 text-gray-700' : ''} ${errors.invoice_number ? 'border-red-500' : ''}`}
              />
              <button
                type="button"
                onClick={handleToggleNoInvoice}
                className="absolute right-2 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600"
                title={watchNoInvoice ? 'Klik untuk isi manual' : 'Faktur tidak ada'}
              >
                {watchNoInvoice ? <Lock size={14} /> : <Unlock size={14} />}
              </button>
            </div>
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
                    <Fragment key={field.id}>
                    <tr>
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
                    <tr>
                      <td colSpan={7} className="bg-gray-50/60 px-3 pb-3 pt-0">
                        <ExpiryBatchSection control={control} index={index} itemQuantity={qty} />
                        {itemErrors?.expiry_batches && (
                          <p className="text-xs text-red-500 mt-1">
                            {(itemErrors.expiry_batches as { message?: string }).message}
                          </p>
                        )}
                      </td>
                    </tr>
                    </Fragment>
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

        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
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
                    disabled={lockDiscountAndStatus}
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
                disabled={lockDiscountAndStatus}
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
                  disabled={lockPaymentMethod}
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
      )}
    </FormModal>

    <ConfirmDialog
      open={isConfirming}
      onOpenChange={(val) => {
        if (!val) {
          setIsConfirming(false)
          setPendingValues(null)
        }
      }}
      title={isEditMode ? 'Update Pembelian' : 'Tambah Pembelian'}
      description={
        isEditMode && initialData && total !== initialData.total_amount
          ? `Total pembelian berubah dari ${formatRupiah(initialData.total_amount)} menjadi ${formatRupiah(total)}. Sisa tagihan: ${formatRupiah(effectiveRemaining)}. Yakin ingin memperbarui?`
          : `Yakin ingin ${isEditMode ? 'memperbarui' : 'menyimpan'} data pembelian ini?`
      }
      confirmLabel="Ya, Simpan"
      isLoading={isPending}
      onConfirm={handleConfirmedSave}
    />
    </>
  )
}

import { useEffect, useRef, useState } from 'react'
import { useForm, useWatch, Controller } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { Lock, Pencil, Plus, Trash2, Unlock } from 'lucide-react'
import { toast } from 'sonner'

import { ConfirmDialog, FormModal } from '@/shared/components'
import { Button } from '@/shared/components/ui/button'
import { Input } from '@/shared/components/ui/input'
import { Label } from '@/shared/components/ui/label'
import { RupiahInput } from '@/shared/components/ui/rupiah-input'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/shared/components/ui/select'
import { useUnitOptionsQuery } from '@/features/products/units'
import { useCategoryOptionsQuery } from '@/features/products/categories'
import { useAuth } from '@/features/auth'

import {
  useCreateProductMutation,
  useCreateProductPackageMutation,
  useDeleteProductPackageMutation,
  useGenerateBarcodeQuery,
  useGenerateSkuQuery,
  useProductDetailQuery,
  useProductPackagesQuery,
  useUpdateProductMutation,
  useUpdateProductPackageMutation,
} from '../products.api'
import { calcMargin } from '../products.utils'
import { productSchema } from '../products.schema'
import type { ProductFormValues, GrosirFormValues } from '../products.schema'
import { GrosirRowForm } from './GrosirRowForm'
import type { PackageRefOption } from './GrosirRowForm'
import type { PackageDraftPayload, Product, ProductPackage } from '../products.types'

interface ProductFormModalProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  product?: Product | null
}

function mapProductToForm(product: Product): ProductFormValues {
  return {
    name: product.name,
    sku: product.sku ?? '',
    barcode: product.barcode ?? '',
    category_id: product.category_id ?? 0,
    purchase_price: product.purchase_price,
    selling_price: product.selling_price,
    stock: product.stock,
    min_stock: product.min_stock,
    unit_id: product.unit_id ?? 0,
    is_active: product.is_active,
  }
}

function packageLabel(pkg: ProductPackage): string {
  const name = pkg.package_name ? `${pkg.unit_name} (${pkg.package_name})` : pkg.unit_name
  return pkg.is_default ? `${name} (Satuan Dasar)` : name
}

// ─── Draft paket (mode Tambah Produk, belum tersimpan ke server) ─────────────

interface GrosirDraft extends GrosirFormValues {
  tempId: number
}

function draftLabel(d: GrosirDraft, unitName: string): string {
  return d.package_name ? `${unitName} (${d.package_name})` : unitName
}

/** Cek apakah menjadikan `newRef` sebagai acuan draft `editingTempId` akan membuat rantai melingkar. */
function wouldCreateCycle(drafts: GrosirDraft[], editingTempId: number | null, newRef: number): boolean {
  if (newRef === 0) return false
  if (newRef === editingTempId) return true
  let current = newRef
  const visited = new Set<number>()
  while (current !== 0) {
    if (current === editingTempId) return true
    if (visited.has(current)) return true
    visited.add(current)
    const d = drafts.find((x) => x.tempId === current)
    if (!d) return false
    current = d.ref_package_id
  }
  return false
}

/** Urutkan draft supaya tiap paket selalu muncul SETELAH paket yang dirujuknya (syarat BE). */
function topoSortDrafts(drafts: GrosirDraft[]): GrosirDraft[] {
  const byId = new Map(drafts.map((d) => [d.tempId, d]))
  const result: GrosirDraft[] = []
  const visited = new Set<number>()

  function visit(d: GrosirDraft) {
    if (visited.has(d.tempId)) return
    visited.add(d.tempId)
    if (d.ref_package_id !== 0) {
      const parent = byId.get(d.ref_package_id)
      if (parent) visit(parent)
    }
    result.push(d)
  }

  drafts.forEach(visit)
  return result
}

// ─── Main modal ───────────────────────────────────────────────────────────────

const defaultValues: ProductFormValues = {
  name: '', sku: '', barcode: '', category_id: 0,
  purchase_price: 0, selling_price: 0, stock: 0, min_stock: 5, unit_id: 0, is_active: true,
}

export function ProductFormModal({ open, onOpenChange, product }: ProductFormModalProps) {
  const isEdit = product != null
  const productId = product?.id
  const { isAdmin } = useAuth()

  const [generateSkuEnabled, setGenerateSkuEnabled] = useState(false)
  const [showGrosirForm, setShowGrosirForm] = useState(false)
  const [editingPackage, setEditingPackage] = useState<ProductPackage | null>(null)
  const [grosirDrafts, setGrosirDrafts] = useState<GrosirDraft[]>([])
  const [editingDraftTempId, setEditingDraftTempId] = useState<number | null>(null)
  const nextTempIdRef = useRef(1)
  const [isConfirming, setIsConfirming] = useState(false)
  const [pendingValues, setPendingValues] = useState<ProductFormValues | null>(null)
  const [barcodeLocked, setBarcodeLocked] = useState(true)

  const { data: detailData, isLoading: isLoadingDetail } = useProductDetailQuery(
    isEdit && open ? (productId as number) : 0
  )
  const { data: existingPackages = [] } = useProductPackagesQuery(
    isEdit && open ? (productId as number) : 0
  )
  const { data: categories = [] } = useCategoryOptionsQuery()
  const { data: units = [] } = useUnitOptionsQuery()

  const { mutate: createProduct, isPending: isCreating } = useCreateProductMutation()
  const { mutate: updateProduct, isPending: isUpdating } = useUpdateProductMutation()
  const { mutate: createPackage, isPending: isCreatingPackage } = useCreateProductPackageMutation(productId ?? 0)
  const { mutate: updatePackage, isPending: isUpdatingPackage } = useUpdateProductPackageMutation(productId ?? 0)
  const { mutate: deletePackage } = useDeleteProductPackageMutation(productId ?? 0)
  const isPending = isCreating || isUpdating
  const isSavingPackage = isCreatingPackage || isUpdatingPackage

  const { register, handleSubmit, reset, setValue, control, formState: { errors } } = useForm<ProductFormValues>({
    resolver: zodResolver(productSchema),
    defaultValues,
  })

  const categoryIdValue = useWatch({ control, name: 'category_id' })
  const unitIdValue = useWatch({ control, name: 'unit_id' })
  const purchasePriceValue = useWatch({ control, name: 'purchase_price' })
  const sellingPriceValue = useWatch({ control, name: 'selling_price' })
  const margin = calcMargin(purchasePriceValue, sellingPriceValue)

  const { data: barcodeData, isFetching: isFetchingBarcode, refetch: refetchBarcode } = useGenerateBarcodeQuery()
  const { data: skuData, isFetching: isFetchingSku, refetch: refetchSku } = useGenerateSkuQuery(
    categoryIdValue ?? 0,
    !isEdit && generateSkuEnabled
  )

  useEffect(() => {
    const barcode = barcodeData?.barcode
    if (barcode) setValue('barcode', barcode, { shouldValidate: true })
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [barcodeData])

  useEffect(() => {
    const sku = skuData?.sku
    if (sku && generateSkuEnabled) setValue('sku', sku, { shouldValidate: true })
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [skuData, generateSkuEnabled])

  useEffect(() => {
    if (!open) return
    if (isEdit && detailData) {
      reset(mapProductToForm(detailData))
    } else if (!isEdit) {
      reset(defaultValues)
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [open, detailData, isEdit])

  const handleClose = () => {
    setGenerateSkuEnabled(false)
    setShowGrosirForm(false)
    setEditingPackage(null)
    setGrosirDrafts([])
    setEditingDraftTempId(null)
    nextTempIdRef.current = 1
    setPendingValues(null)
    setIsConfirming(false)
    setBarcodeLocked(true)
    onOpenChange(false)
  }

  const onSubmit = (values: ProductFormValues) => {
    setPendingValues(values)
    setIsConfirming(true)
  }

  const handleConfirmedSave = () => {
    if (!pendingValues) return

    if (isEdit && productId) {
      // unit_id (satuan anchor) sengaja tidak dikirim — permanen sejak dibuat,
      // cuma bisa diganti lewat mekanisme khusus, bukan form ini.
      updateProduct(
        {
          id: productId,
          name: pendingValues.name,
          sku: pendingValues.sku,
          barcode: pendingValues.barcode,
          category_id: pendingValues.category_id,
          purchase_price: pendingValues.purchase_price,
          selling_price: pendingValues.selling_price,
          stock: pendingValues.stock,
          min_stock: pendingValues.min_stock,
          is_active: pendingValues.is_active,
        },
        {
          onSuccess: () => handleClose(),
          onError: () => setIsConfirming(false),
        }
      )
    } else {
      // Satuan lain yang diisi di form ini (belum tersimpan) dikirim sekaligus bersama
      // produk, diurutkan dulu supaya tiap paket selalu muncul setelah paket yang
      // dirujuknya — server memprosesnya dalam satu transaksi (lihat BE product_repo.go Create).
      const packages: PackageDraftPayload[] = topoSortDrafts(grosirDrafts).map((d) => ({
        temp_id: d.tempId,
        unit_id: d.unit_id,
        package_name: d.package_name || undefined,
        ref_temp_id: d.ref_package_id,
        qty: d.qty,
        ref_qty: d.ref_qty,
        purchase_price: d.purchase_price,
        selling_price: d.selling_price,
      }))
      createProduct(
        { ...pendingValues, packages: packages.length > 0 ? packages : undefined },
        {
          onSuccess: () => handleClose(),
          onError: () => setIsConfirming(false),
        }
      )
    }
  }

  // ─── Grosir/paket: dua mode — draft (lokal, belum simpan) & live (API langsung) ──

  const anchorLabel = isEdit
    ? packageLabel(existingPackages.find((p) => p.is_default) ?? ({} as ProductPackage))
    : `${units.find((u) => u.id === unitIdValue)?.name ?? 'Satuan Dasar'} (Satuan Dasar)`

  const refOptionsFor = (excludeTempId: number | null): PackageRefOption[] => {
    if (isEdit) {
      // Anchor selalu ditaruh pertama supaya jadi default pilihan (opsi paling wajar
      // buat paket baru), tidak tergantung urutan hasil query dari server.
      return [...existingPackages]
        .sort((a, b) => Number(b.is_default) - Number(a.is_default))
        .filter((p) => excludeTempId == null || p.id !== excludeTempId)
        .map((p) => ({ id: p.id, label: packageLabel(p) }))
    }
    return [
      { id: 0, label: anchorLabel },
      ...grosirDrafts
        .filter((d) => d.tempId !== excludeTempId)
        .map((d) => ({ id: d.tempId, label: draftLabel(d, units.find((u) => u.id === d.unit_id)?.name ?? '') })),
    ]
  }

  const handleGrosirSave = (values: GrosirFormValues) => {
    if (isEdit && productId) {
      const payload = {
        unit_id: values.unit_id,
        package_name: values.package_name || undefined,
        ref_package_id: values.ref_package_id,
        qty: values.qty,
        ref_qty: values.ref_qty,
        purchase_price: values.purchase_price,
        selling_price: values.selling_price,
      }
      const onSuccess = () => { setShowGrosirForm(false); setEditingPackage(null) }
      if (editingPackage) {
        updatePackage({ packageId: editingPackage.id, ...payload }, { onSuccess })
      } else {
        createPackage(payload, { onSuccess })
      }
      return
    }

    // Mode draft: validasi tidak melingkar, lalu simpan ke state lokal saja.
    if (wouldCreateCycle(grosirDrafts, editingDraftTempId, values.ref_package_id)) {
      toast.error('Satuan ini tidak bisa merujuk ke rantai yang berujung pada dirinya sendiri')
      return
    }
    if (editingDraftTempId != null) {
      setGrosirDrafts((prev) => prev.map((d) => (d.tempId === editingDraftTempId ? { ...values, tempId: d.tempId } : d)))
    } else {
      const tempId = nextTempIdRef.current++
      setGrosirDrafts((prev) => [...prev, { ...values, tempId }])
    }
    setShowGrosirForm(false)
    setEditingDraftTempId(null)
  }

  const handleDeleteDraft = (tempId: number) => {
    const stillReferenced = grosirDrafts.some((d) => d.ref_package_id === tempId)
    if (stillReferenced) {
      toast.error('Satuan ini masih dirujuk satuan lain, hapus/pindahkan itu dulu')
      return
    }
    setGrosirDrafts((prev) => prev.filter((d) => d.tempId !== tempId))
  }

  const isLoadingContent = isEdit && (isLoadingDetail || !detailData)
  const nonDefaultPackages = existingPackages.filter((p) => !p.is_default)
  const canAddGrosir = isEdit || unitIdValue > 0

  return (
    <>
      <FormModal
        open={open && !isConfirming}
        onOpenChange={(val) => {
          if (!val && !isConfirming) handleClose()
        }}
        title={isEdit ? 'Edit Produk' : 'Tambah Produk'}
        size="lg"
        isLoading={isPending}
        onSubmit={handleSubmit(onSubmit)}
        submitLabel="Simpan"
      >
        {isLoadingContent ? (
          <div className="space-y-4">
            {[1, 2, 3, 4].map((i) => (
              <div key={i} className="h-10 animate-pulse rounded-md bg-gray-100" />
            ))}
          </div>
        ) : (
          <div className="space-y-4">
            {/* Nama Produk */}
            <div className="space-y-1.5">
              <Label htmlFor="name">
                Nama Produk <span className="text-red-500">*</span>
              </Label>
              <Input
                id="name"
                {...register('name')}
                placeholder="Nama produk"
                className={errors.name ? 'border-red-500' : ''}
              />
              {errors.name && <p className="text-xs text-red-500">{errors.name.message}</p>}
            </div>

            {/* Kategori + Satuan */}
            <div className="grid grid-cols-2 gap-3">
              <div className="space-y-1.5">
                <Label>
                  Kategori <span className="text-red-500">*</span>
                </Label>
                <Select
                  value={categoryIdValue !== undefined && categoryIdValue > 0 ? String(categoryIdValue) : ''}
                  onValueChange={(v) => {
                    setValue('category_id', Number(v), { shouldValidate: true })
                    if (!generateSkuEnabled) return
                    setGenerateSkuEnabled(false)
                    setValue('sku', '')
                  }}
                >
                  <SelectTrigger className={errors.category_id ? 'border-red-500' : ''}>
                    <SelectValue placeholder="Pilih Kategori" />
                  </SelectTrigger>
                  <SelectContent>
                    {categories.map((c) => (
                      <SelectItem key={c.id} value={String(c.id)}>
                        {c.name}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
                {errors.category_id && (
                  <p className="text-xs text-red-500">{errors.category_id.message}</p>
                )}
              </div>
              <div className="space-y-1.5">
                <Label>
                  Satuan Dasar <span className="text-red-500">*</span>
                </Label>
                {isEdit ? (
                  <>
                    <div className="flex h-9 items-center rounded-md border bg-gray-50 px-3 text-sm text-gray-700">
                      {units.find((u) => u.id === unitIdValue)?.name ?? '—'}
                    </div>
                    <p className="text-[11px] text-gray-400">
                      Satuan pencatatan stok, permanen sejak produk dibuat — tidak bisa diganti di sini.
                    </p>
                  </>
                ) : (
                  <>
                    <Select
                      value={unitIdValue > 0 ? String(unitIdValue) : ''}
                      onValueChange={(v) => setValue('unit_id', Number(v), { shouldValidate: true })}
                    >
                      <SelectTrigger className={errors.unit_id ? 'border-red-500' : ''}>
                        <SelectValue placeholder="Pilih Satuan" />
                      </SelectTrigger>
                      <SelectContent>
                        {units.map((u) => (
                          <SelectItem key={u.id} value={String(u.id)}>
                            {u.name} ({u.abbreviation})
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                    {errors.unit_id && <p className="text-xs text-red-500">{errors.unit_id.message}</p>}
                  </>
                )}
              </div>
            </div>

            {/* Barcode + SKU */}
            <div className="grid grid-cols-2 gap-3">
              <div className="space-y-1.5">
                <Label htmlFor="barcode">
                  Barcode <span className="text-red-500">*</span>
                </Label>
                <div className="flex gap-1.5">
                  <div className="relative flex-1">
                    <Input
                      id="barcode"
                      {...register('barcode')}
                      readOnly={barcodeLocked}
                      placeholder={isEdit ? '' : (barcodeLocked ? 'Klik 🔓 untuk edit manual' : 'Ketik barcode manual')}
                      className={`pr-9 ${barcodeLocked ? 'bg-gray-50 text-gray-700' : ''} ${errors.barcode ? 'border-red-500' : ''}`}
                    />
                    {!isEdit && (
                      <button
                        type="button"
                        onClick={() => setBarcodeLocked((v) => !v)}
                        className="absolute right-2 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600"
                        title={barcodeLocked ? 'Klik untuk edit manual' : 'Kunci barcode'}
                      >
                        {barcodeLocked ? <Lock size={14} /> : <Unlock size={14} />}
                      </button>
                    )}
                  </div>
                  {!isEdit && barcodeLocked && (
                    <button
                      type="button"
                      disabled={isFetchingBarcode}
                      onClick={() => refetchBarcode()}
                      className="shrink-0 rounded-md border border-gray-300 px-2.5 text-xs text-gray-600 hover:bg-gray-50 disabled:cursor-not-allowed disabled:opacity-40"
                    >
                      {isFetchingBarcode ? '...' : 'Generate'}
                    </button>
                  )}
                </div>
                {errors.barcode && <p className="text-xs text-red-500">{errors.barcode.message}</p>}
              </div>
              <div className="space-y-1.5">
                <Label htmlFor="sku">
                  SKU / Kode <span className="text-red-500">*</span>
                </Label>
                <div className="flex gap-1.5">
                  <Input
                    id="sku"
                    {...register('sku')}
                    readOnly
                    placeholder={isEdit ? '' : (categoryIdValue ? 'Klik Generate' : 'Pilih kategori dulu')}
                    className={`bg-gray-50 text-gray-700 ${errors.sku ? 'border-red-500' : ''}`}
                  />
                  {!isEdit && (
                    <button
                      type="button"
                      disabled={!categoryIdValue || isFetchingSku}
                      onClick={() => {
                        if (generateSkuEnabled) {
                          refetchSku()
                        } else {
                          setGenerateSkuEnabled(true)
                        }
                      }}
                      className="shrink-0 rounded-md border border-gray-300 px-2.5 text-xs text-gray-600 hover:bg-gray-50 disabled:cursor-not-allowed disabled:opacity-40"
                    >
                      {isFetchingSku ? '...' : 'Generate'}
                    </button>
                  )}
                </div>
                {errors.sku && <p className="text-xs text-red-500">{errors.sku.message}</p>}
              </div>
            </div>

            {/* Harga Beli + Harga Jual + Margin */}
            <div className="grid grid-cols-3 gap-3">
              <div className="space-y-1.5">
                <Label htmlFor="purchase_price">Harga Beli</Label>
                <Controller
                  name="purchase_price"
                  control={control}
                  render={({ field }) => (
                    <RupiahInput
                      id="purchase_price"
                      value={field.value}
                      onChange={field.onChange}
                      className={errors.purchase_price ? 'border-red-500' : ''}
                    />
                  )}
                />
                {errors.purchase_price && (
                  <p className="text-xs text-red-500">{errors.purchase_price.message}</p>
                )}
              </div>
              <div className="space-y-1.5">
                <Label htmlFor="selling_price">
                  Harga Jual <span className="text-red-500">*</span>
                </Label>
                <Controller
                  name="selling_price"
                  control={control}
                  render={({ field }) => (
                    <RupiahInput
                      id="selling_price"
                      value={field.value}
                      onChange={field.onChange}
                      className={errors.selling_price ? 'border-red-500' : ''}
                    />
                  )}
                />
                {errors.selling_price && (
                  <p className="text-xs text-red-500">{errors.selling_price.message}</p>
                )}
              </div>
              <div className="space-y-1.5">
                <Label>Margin</Label>
                <div
                  className={`flex h-9 items-center rounded-md border px-3 text-sm font-medium ${
                    margin < 0
                      ? 'border-red-300 bg-red-100 text-red-700'
                      : margin >= 30
                        ? 'border-green-200 bg-green-50 text-green-700'
                        : margin >= 15
                          ? 'border-amber-200 bg-amber-50 text-amber-700'
                          : margin > 0
                            ? 'border-red-200 bg-red-50 text-red-600'
                            : 'border-gray-200 bg-gray-50 text-gray-400'
                  }`}
                >
                  {margin}%
                </div>
                {margin < 0 && (
                  <p className="text-xs text-red-600">
                    Harga jual lebih murah dari harga beli — produk ini akan dijual rugi.
                  </p>
                )}
              </div>
            </div>

            {/* Stok + Stok Minimum */}
            <div className="grid grid-cols-2 gap-3">
              <div className="space-y-1.5">
                <Label htmlFor="stock">Stok</Label>
                <Input
                  id="stock"
                  type="number"
                  min={0}
                  readOnly={!isAdmin}
                  {...register('stock', { valueAsNumber: true })}
                  className={`${!isAdmin ? 'bg-gray-50 text-gray-500 cursor-not-allowed' : ''} ${errors.stock ? 'border-red-500' : ''}`}
                />
                {!isAdmin && (
                  <p className="text-xs text-gray-500">Hanya Admin yang bisa mengubah stok.</p>
                )}
                {errors.stock && <p className="text-xs text-red-500">{errors.stock.message}</p>}
              </div>
              <div className="space-y-1.5">
                <Label htmlFor="min_stock">Stok Minimum</Label>
                <Input
                  id="min_stock"
                  type="number"
                  min={0}
                  {...register('min_stock', { valueAsNumber: true })}
                  className={errors.min_stock ? 'border-red-500' : ''}
                />
                <p className="text-xs text-gray-500">
                  Batas stok terendah sebelum muncul peringatan stok hampir habis.
                </p>
                {errors.min_stock && (
                  <p className="text-xs text-red-500">{errors.min_stock.message}</p>
                )}
              </div>
            </div>

            {/* Grosiran / Satuan Lain */}
            <div className="space-y-2 border-t pt-4">
              <div className="flex items-center justify-between">
                <div>
                  <Label className="text-sm font-semibold">Grosiran / Satuan Lain</Label>
                  <p className="text-xs text-gray-500 mt-0.5">
                    {canAddGrosir
                      ? 'Tambah satuan lain (lebih besar maupun lebih kecil dari satuan yang sudah ada), merujuk ke satuan manapun yang sudah diisi.'
                      : 'Pilih Satuan Dasar dulu, baru bisa tambah satuan lain.'}
                  </p>
                </div>
                {canAddGrosir && !showGrosirForm && (
                  <Button
                    type="button"
                    variant="outline"
                    size="sm"
                    className="gap-1 h-7 text-xs"
                    onClick={() => { setEditingPackage(null); setEditingDraftTempId(null); setShowGrosirForm(true) }}
                  >
                    <Plus size={12} /> Tambah Paket
                  </Button>
                )}
              </div>

              {isEdit && nonDefaultPackages.length > 0 && (
                <div className="rounded-md border text-xs overflow-hidden">
                  <table className="w-full">
                    <thead className="bg-gray-50">
                      <tr>
                        {['Nama Paket', 'Rasio', 'H. Beli', 'H. Jual', ''].map((h) => (
                          <th key={h} className="px-2 py-1.5 text-left font-medium text-gray-600">{h}</th>
                        ))}
                      </tr>
                    </thead>
                    <tbody>
                      {nonDefaultPackages.map((pkg) => (
                        <tr key={pkg.id} className="border-t">
                          <td className="px-2 py-1.5 font-medium">{packageLabel(pkg)}</td>
                          <td className="px-2 py-1.5 text-gray-600">
                            {pkg.qty} {pkg.unit_name} = {pkg.ref_qty}{' '}
                            {existingPackages.find((p) => p.id === pkg.ref_package_id)?.unit_name ?? ''}
                          </td>
                          <td className="px-2 py-1.5">Rp {pkg.purchase_price.toLocaleString('id-ID')}</td>
                          <td className="px-2 py-1.5">Rp {pkg.selling_price.toLocaleString('id-ID')}</td>
                          <td className="px-2 py-1.5">
                            <div className="flex gap-1 justify-end">
                              <Button
                                type="button"
                                variant="ghost"
                                size="icon"
                                className="h-6 w-6 text-gray-500 hover:text-blue-600"
                                onClick={() => { setEditingPackage(pkg); setShowGrosirForm(true) }}
                              >
                                <Pencil size={12} />
                              </Button>
                              <Button
                                type="button"
                                variant="ghost"
                                size="icon"
                                className="h-6 w-6 text-gray-500 hover:text-red-600"
                                onClick={() => deletePackage(pkg.id)}
                              >
                                <Trash2 size={12} />
                              </Button>
                            </div>
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              )}

              {!isEdit && grosirDrafts.length > 0 && (
                <div className="rounded-md border text-xs overflow-hidden">
                  <table className="w-full">
                    <thead className="bg-gray-50">
                      <tr>
                        {['Nama Paket', 'Rasio', 'H. Beli', 'H. Jual', ''].map((h) => (
                          <th key={h} className="px-2 py-1.5 text-left font-medium text-gray-600">{h}</th>
                        ))}
                      </tr>
                    </thead>
                    <tbody>
                      {grosirDrafts.map((d) => {
                        const unitName = units.find((u) => u.id === d.unit_id)?.name ?? ''
                        const refName = d.ref_package_id === 0
                          ? anchorLabel
                          : draftLabel(
                              grosirDrafts.find((x) => x.tempId === d.ref_package_id) ?? d,
                              units.find((u) => u.id === grosirDrafts.find((x) => x.tempId === d.ref_package_id)?.unit_id)?.name ?? ''
                            )
                        return (
                          <tr key={d.tempId} className="border-t">
                            <td className="px-2 py-1.5 font-medium">{draftLabel(d, unitName)}</td>
                            <td className="px-2 py-1.5 text-gray-600">
                              {d.qty} {unitName} = {d.ref_qty} {refName}
                            </td>
                            <td className="px-2 py-1.5">Rp {d.purchase_price.toLocaleString('id-ID')}</td>
                            <td className="px-2 py-1.5">Rp {d.selling_price.toLocaleString('id-ID')}</td>
                            <td className="px-2 py-1.5">
                              <div className="flex gap-1 justify-end">
                                <Button
                                  type="button"
                                  variant="ghost"
                                  size="icon"
                                  className="h-6 w-6 text-gray-500 hover:text-blue-600"
                                  onClick={() => { setEditingDraftTempId(d.tempId); setShowGrosirForm(true) }}
                                >
                                  <Pencil size={12} />
                                </Button>
                                <Button
                                  type="button"
                                  variant="ghost"
                                  size="icon"
                                  className="h-6 w-6 text-gray-500 hover:text-red-600"
                                  onClick={() => handleDeleteDraft(d.tempId)}
                                >
                                  <Trash2 size={12} />
                                </Button>
                              </div>
                            </td>
                          </tr>
                        )
                      })}
                    </tbody>
                  </table>
                </div>
              )}

              {canAddGrosir && showGrosirForm && (
                <GrosirRowForm
                  refOptions={
                    isEdit
                      ? refOptionsFor(editingPackage?.id ?? null)
                      : refOptionsFor(editingDraftTempId)
                  }
                  availableUnits={units}
                  isSaving={isSavingPackage}
                  initialValues={
                    isEdit && editingPackage
                      ? {
                          unit_id: editingPackage.unit_id,
                          package_name: editingPackage.package_name ?? '',
                          ref_package_id: editingPackage.ref_package_id ?? 0,
                          qty: editingPackage.qty ?? 1,
                          ref_qty: editingPackage.ref_qty ?? 1,
                          purchase_price: editingPackage.purchase_price,
                          selling_price: editingPackage.selling_price,
                        }
                      : !isEdit && editingDraftTempId != null
                        ? (() => {
                            const d = grosirDrafts.find((x) => x.tempId === editingDraftTempId)
                            return d ? { ...d } : undefined
                          })()
                        : undefined
                  }
                  onSave={handleGrosirSave}
                  onCancel={() => { setShowGrosirForm(false); setEditingPackage(null); setEditingDraftTempId(null) }}
                />
              )}
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
        title={isEdit ? 'Update Produk' : 'Tambah Produk'}
        description={`Yakin ingin ${isEdit ? 'mengupdate' : 'menambahkan'} produk "${pendingValues?.name}"?`}
        confirmLabel="Ya, Simpan"
        isLoading={isPending}
        onConfirm={handleConfirmedSave}
      />
    </>
  )
}

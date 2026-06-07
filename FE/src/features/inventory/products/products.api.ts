import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'

import { api } from '@/services'
import { queryKeys } from '@/shared/constants'

import type {
  CreatePriceTierPayload,
  CreateProductPackagePayload,
  CreateProductPayload,
  PriceTier,
  Product,
  ProductFilter,
  ProductPackage,
  UpdatePriceTierPayload,
  UpdateProductPayload,
} from './products.types'
// ─── Import ───────────────────────────────────────────────────────────────────

export interface ImportPreviewRow {
  no: number
  nama: string
  barcode: string
  kategori: string
  harga_beli: number
  harga_jual: number
  margin: number
  stok: number
  stok_minimum: number
  satuan: string
  satuan_id: number
  valid: boolean
  errors: string[]
  warnings: string[]
}

export interface ImportPreviewGrosirRow {
  no_produk: number
  nama_paket: string
  satuan: string
  satuan_id: number
  konversi: number
  harga_beli: number
  harga_jual: number
  valid: boolean
  errors: string[]
}

export interface ImportPreviewResponse {
  rows: ImportPreviewRow[]
  grosir: ImportPreviewGrosirRow[]
}

export interface ImportBulkRow {
  no: number
  nama: string
  barcode: string
  kategori: string
  harga_beli: number
  harga_jual: number
  stok: number
  stok_minimum: number
  satuan: string
  satuan_id: number
}

export interface GrosirImportRow {
  no_produk: number
  nama_paket: string
  satuan: string
  satuan_id: number
  konversi: number
  harga_beli: number
  harga_jual: number
}

interface ImportBulkResult {
  success: number
  failed: { baris: number; data: ImportBulkRow; alasan: string }[]
}

export interface ImportBulkPayload {
  rows: ImportBulkRow[]
  grosir: GrosirImportRow[]
}

export function useImportPreviewMutation() {
  return useMutation({
    mutationFn: (file: File) => {
      const formData = new FormData()
      formData.append('file', file)
      return api.post<ImportPreviewResponse>('/products/import-preview', formData)
    },
    onError: (e: Error) => toast.error(e.message),
  })
}

export function useImportProductsBulkMutation() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (payload: ImportBulkPayload) =>
      api.post<ImportBulkResult>('/products/import-bulk', payload),
    onSuccess: (data) => {
      qc.invalidateQueries({ queryKey: queryKeys.products.all() })
      const result = data as unknown as ImportBulkResult
      const failedCount = result.failed?.length ?? 0
      toast.success(
        `${result.success} produk berhasil diimport${failedCount > 0 ? `, ${failedCount} gagal` : ''}`
      )
    },
    onError: (e: Error) => toast.error(e.message),
  })
}

// ─── Generate ─────────────────────────────────────────────────────────────────

export function useGenerateBarcodeQuery() {
  return useQuery({
    queryKey: ['generate-barcode'],
    queryFn: () => api.get<{ barcode: string }>('/products/generate-barcode'),
    enabled: false,
  })
}

export function useGenerateSkuQuery(categoryId: number, enabled: boolean) {
  return useQuery({
    queryKey: ['generate-sku', categoryId],
    queryFn: () => api.get<{ sku: string }>('/products/generate-sku', { category_id: categoryId }),
    enabled: enabled && categoryId > 0,
    staleTime: Infinity,
    gcTime: 0,
  })
}

// ─── Queries ──────────────────────────────────────────────────────────────────

export interface ProductListData {
  items: Product[]
  total: number
  page: number
  limit: number
}

export function useProductListQuery(filter?: ProductFilter) {
  return useQuery({
    queryKey: queryKeys.products.list(filter as Record<string, unknown>),
    queryFn: () => {
      const { page_size, ...rest } = filter ?? {}
      return api.get<ProductListData>('/products', { ...rest, limit: page_size })
    },
  })
}

export function useProductDetailQuery(id: number) {
  return useQuery({
    queryKey: queryKeys.products.detail(id),
    queryFn: () => api.get<Product>(`/products/${id}`),
    enabled: id > 0,
  })
}

export function useProductBarcodeQuery(code: string, enabled: boolean) {
  return useQuery({
    queryKey: queryKeys.products.barcode(code),
    queryFn: () => api.get<{ product: Product }>(`/products/barcode/${code}`),
    enabled: enabled && code.length > 0,
  })
}

export function useProductPackagesQuery(productId: number) {
  return useQuery({
    queryKey: queryKeys.products.productUnits(productId),
    queryFn: () => api.get<ProductPackage[]>(`/products/${productId}/packages`),
    enabled: productId > 0,
  })
}

export function useProductPricesQuery(productId: number) {
  return useQuery({
    queryKey: queryKeys.products.priceTiers(productId),
    queryFn: () => api.get<PriceTier[]>(`/products/${productId}/prices`),
    enabled: productId > 0,
  })
}

// ─── Product Mutations ────────────────────────────────────────────────────────

export function useCreateProductMutation() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (payload: CreateProductPayload) => api.post<Product>('/products', payload),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: queryKeys.products.all() })
    },
    onError: (e: Error) => toast.error(e.message),
  })
}

export function useUpdateProductMutation() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id, ...payload }: UpdateProductPayload & { id: number }) =>
      api.put<Product>(`/products/${id}`, payload),
    onSuccess: (_data, { id }) => {
      qc.invalidateQueries({ queryKey: queryKeys.products.all() })
      qc.invalidateQueries({ queryKey: queryKeys.products.detail(id) })
    },
    onError: (e: Error) => toast.error(e.message),
  })
}

export function useDeleteProductMutation() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id: number) => api.delete<void>(`/products/${id}`),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: queryKeys.products.all() })
    },
    onError: (e: Error) => toast.error(e.message),
  })
}

export function useToggleProductStatusMutation() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id: number) => api.patch<void>(`/products/${id}/toggle-status`),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: queryKeys.products.all() })
    },
    onError: (e: Error) => toast.error(e.message),
  })
}

// ─── Bulk Toggle Status ───────────────────────────────────────────────────────

export function useBulkToggleProductStatusMutation() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (ids: number[]) =>
      Promise.all(ids.map((id) => api.patch<void>(`/products/${id}/toggle-status`))),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: queryKeys.products.all() })
    },
    onError: (e: Error) => toast.error(e.message),
  })
}

// ─── Product Package Mutations ────────────────────────────────────────────────

export function useSaveProductPackagesBulkMutation() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ productId, packages }: { productId: number; packages: CreateProductPackagePayload[] }) =>
      api.post<void>(`/products/${productId}/packages`, { packages }),
    onSuccess: (_data, { productId }) => {
      qc.invalidateQueries({ queryKey: queryKeys.products.productUnits(productId) })
    },
    onError: (e: Error) => toast.error(e.message),
  })
}

export function useAddProductPackageMutation(productId: number) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (payload: CreateProductPackagePayload) =>
      api.post<ProductPackage>(`/products/${productId}/packages`, { packages: [payload] }),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: queryKeys.products.productUnits(productId) })
    },
    onError: (e: Error) => toast.error(e.message),
  })
}

export function useDeleteProductPackageMutation(productId: number) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (packageId: number) => api.delete<void>(`/products/${productId}/packages/${packageId}`),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: queryKeys.products.productUnits(productId) })
    },
    onError: (e: Error) => toast.error(e.message),
  })
}

// ─── Price Tier Mutations ─────────────────────────────────────────────────────

export function useAddPriceTierMutation(productId: number) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (payload: CreatePriceTierPayload) =>
      api.post<PriceTier>(`/products/${productId}/prices`, payload),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: queryKeys.products.priceTiers(productId) })
    },
    onError: (e: Error) => toast.error(e.message),
  })
}

export function useUpdatePriceTierMutation(productId: number) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ priceId, ...payload }: UpdatePriceTierPayload) =>
      api.put<PriceTier>(`/products/${productId}/prices/${priceId}`, payload),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: queryKeys.products.priceTiers(productId) })
    },
    onError: (e: Error) => toast.error(e.message),
  })
}

export function useDeletePriceTierMutation(productId: number) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (priceId: number) => api.delete<void>(`/products/${productId}/prices/${priceId}`),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: queryKeys.products.priceTiers(productId) })
    },
    onError: (e: Error) => toast.error(e.message),
  })
}

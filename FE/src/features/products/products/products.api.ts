import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'

import { api, apiClient } from '@/services'
import { queryKeys } from '@/shared/constants'
import type { PaginatedData } from '@/shared/types'

import type {
  CreatePriceTierPayload,
  CreateProductPackagePayload,
  CreateProductPayload,
  GrosirImportRow,
  ImportBulkPayload,
  ImportBulkResult,
  ImportBulkRow,
  ImportPreviewGrosirRow,
  ImportPreviewResponse,
  ImportPreviewRow,
  PriceTier,
  Product,
  ProductListFilter,
  ProductOption,
  ProductPackage,
  ProductSearchOption,
  UpdatePriceTierPayload,
  UpdateProductPayload,
} from './products.types'

export type {
  ImportPreviewRow,
  ImportPreviewGrosirRow,
  ImportPreviewResponse,
  ImportBulkRow,
  GrosirImportRow,
  ImportBulkResult,
  ImportBulkPayload,
}

// ─── Import Helpers ───────────────────────────────────────────────────────────

export async function downloadImportTemplate(): Promise<void> {
  const response = await apiClient.post('/products/import-template', {}, { responseType: 'blob' })
  const url = URL.createObjectURL(response.data as Blob)
  const a = document.createElement('a')
  a.href = url
  a.download = 'template_import_produk.xlsx'
  a.click()
  URL.revokeObjectURL(url)
}

// ─── Import Mutations ─────────────────────────────────────────────────────────

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
    queryFn: () => api.post<{ barcode: string }>('/products/generate-barcode', {}),
    enabled: false,
  })
}

export function useGenerateSkuQuery(categoryId: number, enabled: boolean) {
  return useQuery({
    queryKey: ['generate-sku', categoryId],
    queryFn: () => api.post<{ sku: string }>('/products/generate-sku', { category_id: categoryId }),
    enabled: enabled && categoryId > 0,
    staleTime: Infinity,
    gcTime: 0,
  })
}

// ─── Queries ──────────────────────────────────────────────────────────────────

export function useProductListQuery(filter: ProductListFilter) {
  return useQuery({
    queryKey: queryKeys.products.list(filter as unknown as Record<string, unknown>),
    queryFn: () => api.post<PaginatedData<Product>>('/products/list', filter),
  })
}

export function useProductOptionsQuery() {
  return useQuery({
    queryKey: queryKeys.products.options(),
    queryFn: () => api.post<ProductOption[]>('/products/options', {}),
  })
}

export function useProductSearchQuery(keyword: string, enabled = true) {
  return useQuery({
    queryKey: ['products', 'search', keyword],
    queryFn: () => api.post<ProductSearchOption[]>('/products/search', { q: keyword, limit: 20 }),
    enabled: enabled && keyword.length >= 2,
  })
}

export function fetchProductDetail(id: number) {
  return api.post<Product>(`/products/detail/${id}`, {})
}

export function fetchProductPackages(productId: number) {
  return api.post<ProductPackage[]>(`/products/${productId}/packages/list`, {})
}

export function fetchProductPrices(productId: number) {
  return api.post<PriceTier[]>(`/products/${productId}/prices/list`, {})
}

export function useProductDetailQuery(id: number) {
  return useQuery({
    queryKey: queryKeys.products.detail(id),
    queryFn: () => fetchProductDetail(id),
    enabled: id > 0,
  })
}

export function useProductBarcodeQuery(code: string, enabled: boolean) {
  return useQuery({
    queryKey: queryKeys.products.barcode(code),
    queryFn: () => api.post<{ product: Product }>(`/products/by-barcode/${code}`, {}),
    enabled: enabled && code.length > 0,
  })
}

export function useProductPackagesQuery(productId: number) {
  return useQuery({
    queryKey: queryKeys.products.productUnits(productId),
    queryFn: () => fetchProductPackages(productId),
    enabled: productId > 0,
  })
}

export function useProductPricesQuery(productId: number) {
  return useQuery({
    queryKey: queryKeys.products.priceTiers(productId),
    queryFn: () => fetchProductPrices(productId),
    enabled: productId > 0,
  })
}

// ─── Product Mutations ────────────────────────────────────────────────────────

export function useCreateProductMutation() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (payload: CreateProductPayload) => api.post<Product>('/products/create', payload),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: queryKeys.products.all() })
      toast.success('Produk berhasil ditambahkan')
    },
    onError: (e: Error) => toast.error(e.message),
  })
}

export function useUpdateProductMutation() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id, ...payload }: UpdateProductPayload & { id: number }) =>
      api.post<Product>(`/products/update/${id}`, payload),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: queryKeys.products.all() })
      toast.success('Produk berhasil diperbarui')
    },
    onError: (e: Error) => toast.error(e.message),
  })
}

export function useDeleteProductMutation() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id: number) => api.post<void>(`/products/delete/${id}`, {}),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: queryKeys.products.all() })
      toast.success('Produk berhasil dihapus')
    },
    onError: (e: Error) => toast.error(e.message),
  })
}

export function useToggleProductStatusMutation() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id }: { id: number; isActive: boolean }) =>
      api.post<void>(`/products/toggle-status/${id}`, {}),
    onSuccess: (_data, { isActive }) => {
      qc.invalidateQueries({ queryKey: queryKeys.products.all() })
      toast.success(`Produk berhasil ${isActive ? 'dinonaktifkan' : 'diaktifkan'}`)
    },
    onError: (e: Error) => toast.error(e.message),
  })
}

// ─── Bulk Toggle Status ───────────────────────────────────────────────────────

export function useBulkToggleProductStatusMutation() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ ids }: { ids: number[]; label: string }) =>
      Promise.all(ids.map((id) => api.post<void>(`/products/toggle-status/${id}`, {}))),
    onSuccess: (_data, { ids, label }) => {
      qc.invalidateQueries({ queryKey: queryKeys.products.all() })
      toast.success(`${ids.length} produk berhasil ${label}`)
    },
    onError: (e: Error) => toast.error(e.message),
  })
}

// ─── Product Package Mutations ────────────────────────────────────────────────

export function useSaveProductPackagesBulkMutation() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ productId, packages }: { productId: number; packages: CreateProductPackagePayload[] }) =>
      api.post<void>(`/products/${productId}/packages/save`, { packages }),
    onSuccess: (_data, { productId }) => {
      qc.invalidateQueries({ queryKey: queryKeys.products.productUnits(productId) })
      toast.success('Kemasan produk berhasil disimpan')
    },
    onError: (e: Error) => toast.error(e.message),
  })
}

export function useDeleteProductPackageMutation(productId: number) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (packageId: number) =>
      api.post<void>(`/products/${productId}/packages/delete/${packageId}`, {}),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: queryKeys.products.productUnits(productId) })
    },
    onError: (e: Error) => toast.error(e.message),
  })
}

// ─── Price Tier Mutations ─────────────────────────────────────────────────────

export function useSavePriceTiersMutation(productId: number) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (prices: CreatePriceTierPayload[]) =>
      api.post<void>(`/products/${productId}/prices/save`, { prices }),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: queryKeys.products.priceTiers(productId) })
    },
    onError: (e: Error) => toast.error(e.message),
  })
}

export function useAddPriceTierMutation(productId: number) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (payload: CreatePriceTierPayload) =>
      api.post<void>(`/products/${productId}/prices/create`, payload),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: queryKeys.products.priceTiers(productId) })
      toast.success('Harga tier berhasil ditambahkan')
    },
    onError: (e: Error) => toast.error(e.message),
  })
}

export function useUpdatePriceTierMutation(productId: number) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ priceId, ...payload }: UpdatePriceTierPayload & { priceId: number }) =>
      api.post<void>(`/products/${productId}/prices/update/${priceId}`, payload),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: queryKeys.products.priceTiers(productId) })
      toast.success('Harga tier berhasil diperbarui')
    },
    onError: (e: Error) => toast.error(e.message),
  })
}

export function useDeletePriceTierMutation(productId: number) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (priceId: number) =>
      api.post<void>(`/products/${productId}/prices/delete/${priceId}`, {}),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: queryKeys.products.priceTiers(productId) })
      toast.success('Harga tier berhasil dihapus')
    },
    onError: (e: Error) => toast.error(e.message),
  })
}

import { TriangleAlert } from 'lucide-react'

import { DetailField, FormModal, StatusBadge } from '@/shared/components'
import { Button } from '@/shared/components/ui/button'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/shared/components/ui/tabs'
import { formatDate, formatRupiah } from '@/shared/utils'

import { useProductDetailQuery, useProductPackagesQuery } from '../products.api'
import { useExpiryBatchHistoryQuery } from '../expiry-batches.api'
import type { ExpiryBatchStatus } from '../expiry-batches.types'
import { calcMargin, formatPackageBreakdown } from '../products.utils'
import { ProductPurchaseHistory } from './ProductPurchaseHistory'
import { ProductSaleHistory } from './ProductSaleHistory'

const EXPIRY_STATUS_LABEL: Record<ExpiryBatchStatus, { label: string; className: string }> = {
  active: { label: 'Perlu Dicek', className: 'bg-amber-100 text-amber-700' },
  cleared: { label: 'Sudah Dicek, Aman', className: 'bg-green-100 text-green-700' },
  written_off: { label: 'Dimusnahkan', className: 'bg-red-100 text-red-600' },
}

interface ProductDetailModalProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  productId?: number
}

export function ProductDetailModal({ open, onOpenChange, productId }: ProductDetailModalProps) {
  const enabled = open && (productId ?? 0) > 0
  const { data: product, isLoading } = useProductDetailQuery(enabled ? (productId as number) : 0)
  const { data: units = [] } = useProductPackagesQuery(enabled ? (productId as number) : 0)
  const { data: expiryHistory = [] } = useExpiryBatchHistoryQuery(enabled ? productId : undefined)

  const margin = product ? calcMargin(product.purchase_price, product.selling_price) : 0
  const grosirUnits = (units ?? []).filter((u) => !u.is_default)
  const defaultUnit = (units ?? []).find((u) => u.is_default)
  const refUnitName = (refPackageId: number | null) =>
    units.find((p) => p.id === refPackageId)?.unit_name ?? product?.unit_name ?? ''

  return (
    <FormModal
      open={open}
      onOpenChange={onOpenChange}
      title="Detail Produk"
      size="lg"
      hideFooter
    >
      {isLoading || !product ? (
        <div className="space-y-4">
          {[1, 2, 3, 4].map((i) => (
            <div key={i} className="h-8 animate-pulse rounded-md bg-gray-100" />
          ))}
        </div>
      ) : (
        <Tabs defaultValue="detail" className="w-full">
          <TabsList className="mb-4">
            <TabsTrigger value="detail">Detail</TabsTrigger>
            <TabsTrigger value="packages">Paket Satuan</TabsTrigger>
            <TabsTrigger value="purchase-history">Riwayat Pembelian</TabsTrigger>
            <TabsTrigger value="sale-history">Riwayat Penjualan</TabsTrigger>
          </TabsList>

          {/* ═══ TAB: Detail ═══ */}
          <TabsContent value="detail" className="space-y-4 text-sm">
            {/* Identitas */}
            <div className="grid grid-cols-1 gap-3 sm:grid-cols-2">
              <DetailField label="Nama Produk" value={product.name} />
              <DetailField label="Status">
                <StatusBadge status={product.is_active ? 'active' : 'inactive'} />
              </DetailField>
            </div>

            {product.needs_stock_review && (
              <div className="flex items-start gap-2 rounded-md border border-amber-200 bg-amber-50 px-3 py-2 text-xs text-amber-800">
                <TriangleAlert size={14} className="mt-0.5 shrink-0" />
                <span>
                  <span className="font-semibold">Perlu ditinjau:</span> data stok produk ini belum
                  terverifikasi penuh dari migrasi sistem.
                  {product.stock_review_note ? ` ${product.stock_review_note}` : ''}
                </span>
              </div>
            )}

            <div className="grid grid-cols-1 gap-3 sm:grid-cols-2">
              <DetailField label="Barcode">
                <code className="text-xs text-gray-700">{product.barcode || '—'}</code>
              </DetailField>
              <DetailField label="SKU / Kode">
                <code className="text-xs text-gray-700">{product.sku || '—'}</code>
              </DetailField>
            </div>

            <div className="grid grid-cols-1 gap-3 sm:grid-cols-2">
              <DetailField label="Kategori" value={product.category_name || '—'} />
              <DetailField label="Satuan" value={product.unit_name || '—'} />
            </div>

            {/* Harga */}
            <div className="rounded-md border bg-gray-50 p-3 grid grid-cols-1 gap-3 sm:grid-cols-3">
              <DetailField label="Harga Beli" value={formatRupiah(product.purchase_price)} />
              <DetailField label="Harga Jual" value={formatRupiah(product.selling_price)} />
              <DetailField label="Margin">
                <span
                  className={`inline-flex items-center rounded-full px-2 py-0.5 text-xs font-medium ${
                    margin >= 30
                      ? 'bg-green-100 text-green-700'
                      : margin >= 15
                        ? 'bg-amber-100 text-amber-700'
                        : 'bg-red-100 text-red-600'
                  }`}
                >
                  {margin}%
                </span>
              </DetailField>
            </div>

            {/* Stok */}
            <div className="grid grid-cols-1 gap-3 sm:grid-cols-2">
              <DetailField label="Stok">
                <span
                  className={`font-medium ${
                    product.stock === 0
                      ? 'text-red-600'
                      : product.is_low_stock
                        ? 'text-amber-600'
                        : 'text-gray-800'
                  }`}
                >
                  {formatPackageBreakdown(units)}
                </span>
              </DetailField>
              <DetailField label="Stok Minimum" value={String(product.min_stock)} />
            </div>

            {product.reserved_qty > 0 && (
              <div className="grid grid-cols-1 gap-3 sm:grid-cols-2">
                <DetailField label="Stok Direservasi">
                  <span className="font-medium text-amber-600">
                    {formatPackageBreakdown(units, (p) => p.reserved_qty)}
                  </span>
                </DetailField>
                <DetailField label="Stok Tersedia">
                  <span className="font-medium text-gray-800">
                    {formatPackageBreakdown(units, (p) => p.stock - p.reserved_qty)}
                  </span>
                </DetailField>
              </div>
            )}

            {/* Batch Expired */}
            {expiryHistory.length > 0 && (
              <div className="space-y-2 border-t pt-3">
                <p className="text-xs font-semibold text-gray-600 uppercase tracking-wide">Batch Expired</p>
                <div className="rounded-md border overflow-x-auto">
                  <table className="w-full text-xs">
                    <thead className="bg-gray-50">
                      <tr>
                        {['Qty', 'Tanggal Expired', 'Status'].map((h) => (
                          <th key={h} className="px-2 py-1.5 text-left font-medium text-gray-600">{h}</th>
                        ))}
                      </tr>
                    </thead>
                    <tbody>
                      {expiryHistory.map((b) => (
                        <tr key={b.id} className="border-t">
                          <td className="px-2 py-1.5 font-medium">{b.qty}</td>
                          <td className="px-2 py-1.5 text-gray-600">{formatDate(b.expired_date)}</td>
                          <td className="px-2 py-1.5">
                            <span
                              className={`inline-flex items-center rounded-full px-2 py-0.5 text-[11px] font-medium ${EXPIRY_STATUS_LABEL[b.status].className}`}
                            >
                              {EXPIRY_STATUS_LABEL[b.status].label}
                            </span>
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              </div>
            )}
          </TabsContent>

          {/* ═══ TAB: Paket Satuan ═══ */}
          <TabsContent value="packages" className="space-y-3 text-sm">
            <p className="text-xs text-gray-500">Daftar satuan/kemasan yang tersedia untuk produk ini.</p>

            {/* Default package */}
            {defaultUnit && (
              <div className="flex items-center justify-between rounded-lg border border-blue-200 bg-blue-50 p-3">
                <div>
                  <div className="flex items-center gap-2">
                    <span className="font-semibold text-gray-800">{defaultUnit.unit_name}</span>
                    <span className="inline-flex items-center rounded-full bg-blue-200 px-2 py-0.5 text-[10px] font-semibold text-blue-800">
                      Satuan Dasar
                    </span>
                  </div>
                  <p className="text-xs text-gray-500 mt-0.5">1 {defaultUnit.unit_name} = 1 unit terkecil</p>
                </div>
                <div className="text-right text-xs">
                  <p>Beli: {formatRupiah(defaultUnit.purchase_price)} | Jual: <strong>{formatRupiah(defaultUnit.selling_price)}</strong></p>
                  <p className="text-gray-500 mt-0.5">Stok: {defaultUnit.stock} {defaultUnit.unit_name}</p>
                </div>
              </div>
            )}

            {/* Other packages */}
            {grosirUnits.map((u) => (
              <div key={u.id} className="flex items-center justify-between rounded-lg border p-3">
                <div>
                  <div className="flex items-center gap-2">
                    <span className="font-semibold text-gray-800">
                      {u.package_name ? `${u.unit_name} (${u.package_name})` : u.unit_name}
                    </span>
                    <span className="inline-flex items-center rounded-full bg-gray-100 px-2 py-0.5 text-[10px] font-medium text-gray-600">
                      Turunan
                    </span>
                  </div>
                  <p className="text-xs text-gray-500 mt-0.5">
                    {u.qty} {u.unit_name} = {u.ref_qty} {refUnitName(u.ref_package_id)}
                  </p>
                </div>
                <div className="text-right text-xs">
                  <p>Beli: {formatRupiah(u.purchase_price)} | Jual: <strong>{formatRupiah(u.selling_price)}</strong></p>
                  <p className="text-gray-500 mt-0.5">Stok: {u.stock} {u.unit_name}</p>
                </div>
              </div>
            ))}

            {grosirUnits.length === 0 && !defaultUnit && (
              <p className="text-center text-xs text-gray-400 py-6">Belum ada paket satuan</p>
            )}
          </TabsContent>

          {/* ═══ TAB: Riwayat Pembelian ═══ */}
          <TabsContent value="purchase-history">
            <ProductPurchaseHistory productId={product.id} />
          </TabsContent>

          {/* ═══ TAB: Riwayat Penjualan ═══ */}
          <TabsContent value="sale-history">
            <ProductSaleHistory productId={product.id} />
          </TabsContent>
        </Tabs>
      )}

      <div className="flex justify-end border-t pt-3 mt-4">
        <Button variant="outline" onClick={() => onOpenChange(false)}>
          Tutup
        </Button>
      </div>
    </FormModal>
  )
}

import { FileDown } from 'lucide-react'

import { Button } from '@/shared/components/ui/button'

interface ProductBulkActionBarProps {
  count: number
  allActive: boolean
  allInactive: boolean
  isBulkToggling: boolean
  onToggleStatus: () => void
  onExport: () => void
  onPrintLabel: () => void
  onClear: () => void
}

export function ProductBulkActionBar({
  count,
  allActive,
  allInactive,
  isBulkToggling,
  onToggleStatus,
  onExport,
  onPrintLabel,
  onClear,
}: ProductBulkActionBarProps) {
  const showBulkStatus = allActive || allInactive

  return (
    <div className="flex items-center gap-3 rounded-lg border bg-blue-50 px-4 py-2 text-sm">
      <span className="font-medium text-blue-700">{count} produk dipilih</span>
      <div className="ml-auto flex gap-2">
        {showBulkStatus && (
          <Button
            variant="outline"
            size="sm"
            disabled={isBulkToggling}
            onClick={onToggleStatus}
            className={allActive ? 'text-amber-600 hover:text-amber-700' : 'text-green-600 hover:text-green-700'}
          >
            {isBulkToggling ? 'Memproses...' : allActive ? 'Nonaktifkan' : 'Aktifkan'}
          </Button>
        )}
        <Button variant="outline" size="sm" onClick={onExport} className="gap-1">
          <FileDown size={14} />
          Export Excel
        </Button>
        <Button variant="outline" size="sm" onClick={onPrintLabel}>
          Cetak Label
        </Button>
        <Button variant="outline" size="sm" onClick={onClear}>
          Batalkan Pilihan
        </Button>
      </div>
    </div>
  )
}

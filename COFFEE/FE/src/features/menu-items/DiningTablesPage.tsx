import { diningTableApi } from './menu-items.api'
import type { DiningTable } from './menu-items.types'
import { SimpleMasterPage } from './components/SimpleMasterPage'

// Daftar meja kedai, dikelola admin/owner; kasir hanya melihat (keputusan 56, 97, 98).
// Meja dipilih di kasir untuk pesanan dine-in (tahap 4.7).
export function DiningTablesPage() {
  return (
    <SimpleMasterPage<DiningTable>
      title="Meja"
      description="Daftar meja yang bisa dipilih kasir untuk pesanan makan di tempat. Urutan menentukan urutan tampil di kasir."
      itemLabel="Meja"
      menuKey="menu_jual.meja"
      masterApi={diningTableApi}
      usage={(row) => row.usage_count}
      usageLabel={(n) => (n > 0 ? `${n} pesanan` : 'Belum dipakai')}
      namePlaceholder="Contoh: Meja 1, Teras 2"
    />
  )
}

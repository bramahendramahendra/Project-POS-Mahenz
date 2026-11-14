import { menuCategoryApi } from './menu-items.api'
import type { MenuCategory } from './menu-items.types'
import { SimpleMasterPage } from './components/SimpleMasterPage'

// Kategori menu jualan (mis. Kopi, Non-Kopi, Makanan).
export function MenuCategoriesPage() {
  return (
    <SimpleMasterPage<MenuCategory>
      title="Kategori Menu"
      description="Pengelompokan menu di kasir dan laporan."
      itemLabel="Kategori"
      menuKey="menu_jual.kategori"
      masterApi={menuCategoryApi}
      usage={(row) => row.item_count}
      usageLabel={(n) => (n > 0 ? `${n} menu` : 'Belum dipakai')}
      namePlaceholder="Contoh: Kopi, Non-Kopi, Makanan"
    />
  )
}

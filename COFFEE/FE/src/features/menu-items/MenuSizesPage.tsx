import { menuSizeApi } from './menu-items.api'
import type { MenuSize } from './menu-items.types'
import { SimpleMasterPage } from './components/SimpleMasterPage'

// Daftar ukuran varian menu, dikelola admin (keputusan 85).
export function MenuSizesPage() {
  return (
    <SimpleMasterPage<MenuSize>
      title="Ukuran"
      description="Ukuran yang bisa dipilih untuk varian menu (mis. Regular, Large)."
      itemLabel="Ukuran"
      menuKey="menu_jual.ukuran"
      masterApi={menuSizeApi}
      usage={(row) => row.usage_count}
      usageLabel={(n) => (n > 0 ? `${n} kombinasi varian` : 'Belum dipakai')}
      namePlaceholder="Contoh: Regular, Large"
    />
  )
}

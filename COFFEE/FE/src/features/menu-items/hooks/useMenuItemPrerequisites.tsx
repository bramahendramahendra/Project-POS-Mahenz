import { Tag } from 'lucide-react'

import { menuCategoryApi } from '../menu-items.api'

// Syarat halaman Daftar Menu: kategori menu wajib dipilih di form menu (pola sama dengan useProductPrerequisites).
export function useMenuItemPrerequisites() {
  const { data: categoriesData, isLoading } = menuCategoryApi.useOptions()
  const hasCategories = (categoriesData ?? []).length > 0

  return {
    isLoading,
    items: [
      {
        label: 'Belum ada kategori menu — tambahkan di menu Kategori Menu',
        metLabel: 'Kategori menu sudah tersedia',
        met: hasCategories,
        icon: <Tag size={24} />,
      },
    ],
  }
}

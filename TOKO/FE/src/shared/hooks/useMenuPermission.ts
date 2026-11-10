import { useMenuStore } from '@/features/menu/menu.store'

/**
 * Hook untuk mengambil permission user pada menu tertentu.
 * Gunakan key_name menu sebagai argument, contoh: 'inventory.products'
 */
export function useMenuPermission(menuKey: string) {
  const getPermission = useMenuStore((s) => s.getPermission)
  const perm = getPermission(menuKey)

  return {
    can_view:   perm?.can_view   ?? false,
    can_create: perm?.can_create ?? false,
    can_edit:   perm?.can_edit   ?? false,
    can_delete: perm?.can_delete ?? false,
  }
}

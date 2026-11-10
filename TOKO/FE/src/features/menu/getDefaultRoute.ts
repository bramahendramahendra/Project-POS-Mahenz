import { ROUTES } from '@/shared/constants/routes'

import type { MenuItem } from './menu.types'

/**
 * Cari route pertama yang boleh diakses user, urut order_index, depth-first.
 * Dipakai sebagai tujuan redirect setelah login dan fallback saat akses ditolak,
 * menggantikan logic hardcode "kasir -> /kasir, selain itu -> /dashboard".
 */
export function getDefaultRoute(menus: MenuItem[]): string {
  const sorted = [...menus].sort((a, b) => a.order_index - b.order_index)

  for (const item of sorted) {
    if (item.path && item.permission.can_view) {
      return item.path
    }
    if (item.children.length > 0) {
      const childRoute = getDefaultRoute(item.children)
      if (childRoute !== ROUTES.LOGIN) {
        return childRoute
      }
    }
  }

  return ROUTES.LOGIN
}

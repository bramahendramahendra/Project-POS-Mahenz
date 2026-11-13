import { ROUTES } from '@/shared/constants/routes'

import type { MenuItem } from './menu.types'

/** Route pertama yang boleh diakses (urut order_index, depth-first), atau null bila tidak ada. */
function findFirstRoute(menus: MenuItem[]): string | null {
  const sorted = [...menus].sort((a, b) => a.order_index - b.order_index)

  for (const item of sorted) {
    if (item.path && item.permission.can_view) {
      return item.path
    }
    if (item.children.length > 0) {
      const childRoute = findFirstRoute(item.children)
      if (childRoute !== null) {
        return childRoute
      }
    }
  }

  return null
}

/**
 * Cari route pertama yang boleh diakses user, urut order_index, depth-first.
 * Dipakai sebagai tujuan redirect setelah login dan fallback saat akses ditolak,
 * menggantikan logic hardcode "kasir -> /kasir, selain itu -> /dashboard".
 *
 * Bila user tidak punya akses menu apa pun (mis. akun baru dari Portal dengan role
 * "menunggu akses"), tujuannya halaman "menunggu akses", bukan halaman login
 * (halaman login akan mengembalikan user yang sudah login ke sini lagi dan membuat loop).
 */
export function getDefaultRoute(menus: MenuItem[]): string {
  return findFirstRoute(menus) ?? ROUTES.PENDING_ACCESS
}

import { useEffect } from 'react'
import { Navigate } from 'react-router-dom'

import { ROUTES } from '@/shared/constants'
import { useMyMenusQuery } from '@/features/menu/menu.api'
import { useMenuStore } from '@/features/menu/menu.store'
import { getDefaultRoute } from '@/features/menu/getDefaultRoute'

import { useAuthStore } from '../auth.store'

export function RootRedirect() {
  const isAuthenticated = useAuthStore((s) => s.isAuthenticated)

  const menus = useMenuStore((s) => s.menus)
  const isLoaded = useMenuStore((s) => s.isLoaded)
  const setMenus = useMenuStore((s) => s.setMenus)

  const { data } = useMyMenusQuery(isAuthenticated && !isLoaded)

  useEffect(() => {
    if (data && !isLoaded) setMenus(data)
  }, [data, isLoaded, setMenus])

  if (!isAuthenticated) {
    return <Navigate to={ROUTES.LOGIN} replace />
  }

  if (!isLoaded) {
    return (
      <div className="flex min-h-screen items-center justify-center text-sm text-gray-400">
        Memuat...
      </div>
    )
  }

  return <Navigate to={getDefaultRoute(menus)} replace />
}

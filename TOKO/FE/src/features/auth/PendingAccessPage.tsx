import { useEffect } from 'react'
import { Navigate } from 'react-router-dom'
import { Hourglass, Loader2, LogOut } from 'lucide-react'

import { config, ROUTES } from '@/shared/constants'
import { getDefaultRoute } from '@/features/menu/getDefaultRoute'
import { useMyMenusQuery } from '@/features/menu/menu.api'
import { useMenuStore } from '@/features/menu/menu.store'
import { appUrl, portalUrl } from '@/shared/utils/appUrl'

import { useLogoutMutation } from './auth.api'
import { useAuthStore } from './auth.store'

// Tujuan akun yang sudah login tetapi belum punya izin menu apa pun
// (mis. akun baru dari Portal dengan role "menunggu akses").
export function PendingAccessPage() {
  const isAuthenticated = useAuthStore((s) => s.isAuthenticated)
  const user = useAuthStore((s) => s.user)
  const menus = useMenuStore((s) => s.menus)
  const isLoaded = useMenuStore((s) => s.isLoaded)
  const setMenus = useMenuStore((s) => s.setMenus)
  const { mutate: logout, isPending } = useLogoutMutation()

  const { data } = useMyMenusQuery(isAuthenticated && !isLoaded)

  useEffect(() => {
    if (data && !isLoaded) setMenus(data)
  }, [data, isLoaded, setMenus])

  if (!isAuthenticated || !user) {
    return <Navigate to={ROUTES.LOGIN} replace />
  }

  if (!isLoaded) {
    return (
      <div className="flex min-h-screen items-center justify-center text-sm text-gray-400">
        <Loader2 size={16} className="mr-2 animate-spin" /> Memuat...
      </div>
    )
  }

  // Sudah punya akses menu: tidak perlu di halaman ini.
  const target = getDefaultRoute(menus)
  if (target !== ROUTES.PENDING_ACCESS) {
    return <Navigate to={target} replace />
  }

  return (
    <div className="min-h-screen flex flex-col items-center justify-center bg-[#2c3e50] px-4">
      <div className="mb-6 text-center">
        <h1 className="text-3xl font-bold text-white tracking-wide">{config.appName}</h1>
      </div>

      <div className="bg-white rounded-xl shadow-lg p-8 w-full max-w-md text-center space-y-4">
        <Hourglass size={36} className="mx-auto text-amber-500" />
        <h2 className="text-lg font-semibold text-gray-800">Menunggu pengaturan akses</h2>
        <p className="text-sm text-gray-600">
          Halo, {user.fullName}. Akun Anda sudah terdaftar, tetapi admin belum mengatur izin akses
          di aplikasi ini. Hubungi admin, lalu buka kembali aplikasi ini
          {portalUrl ? ' dari Portal' : ''} setelah akses diatur.
        </p>

        <div className="space-y-2 pt-2">
          {portalUrl && (
            <a
              href={portalUrl}
              className="block w-full rounded-md bg-[#2c3e50] px-4 py-2 text-sm font-medium text-white hover:bg-[#34495e]"
            >
              Kembali ke Portal
            </a>
          )}
          <button
            type="button"
            onClick={() => logout()}
            disabled={isPending}
            className="flex w-full items-center justify-center gap-2 rounded-md border border-gray-300 px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50 disabled:opacity-60"
          >
            <LogOut size={16} /> Keluar
          </button>
          {/* Cadangan bila Portal belum terpasang: kembali ke halaman login */}
          {!portalUrl && (
            <a href={appUrl(ROUTES.LOGIN)} className="block text-xs text-gray-400 hover:text-gray-600">
              Ke halaman login
            </a>
          )}
        </div>
      </div>
    </div>
  )
}

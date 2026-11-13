import { useEffect } from 'react'
import { ExternalLink, Loader2 } from 'lucide-react'

import { config } from '@/shared/constants'
import { useMyMenusQuery } from '@/features/menu/menu.api'
import { useMenuStore } from '@/features/menu/menu.store'
import { getDefaultRoute } from '@/features/menu/getDefaultRoute'
import { appUrl, portalUrl } from '@/shared/utils/appUrl'

import { useLoginOptionsQuery } from './auth.api'
import { useAuth } from './hooks/useAuth'
import { LoginForm } from './components/LoginForm'

export function LoginPage() {
  const { isAuthenticated } = useAuth()
  const menus = useMenuStore((s) => s.menus)
  const isLoaded = useMenuStore((s) => s.isLoaded)
  const setMenus = useMenuStore((s) => s.setMenus)

  const { data } = useMyMenusQuery(isAuthenticated && !isLoaded)
  const { data: options, isLoading: optionsLoading } = useLoginOptionsQuery()

  useEffect(() => {
    if (data && !isLoaded) setMenus(data)
  }, [data, isLoaded, setMenus])

  useEffect(() => {
    if (isAuthenticated && isLoaded) {
      window.location.href = appUrl(getDefaultRoute(menus))
    }
  }, [isAuthenticated, isLoaded, menus])

  if (isAuthenticated) return null

  // Bila opsi login gagal dimuat, form login langsung tetap ditampilkan agar tidak ada yang terkunci.
  const directEnabled = options ? options.direct_login_enabled : true
  const portalAvailable = !!portalUrl && (options ? options.portal_login_enabled : true)

  const portalButton = (
    <a
      href={portalUrl}
      className="flex w-full items-center justify-center gap-2 rounded-md border border-[#2c3e50] px-4 py-2 text-sm font-medium text-[#2c3e50] transition-colors hover:bg-[#2c3e50]/5 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[#2c3e50] focus-visible:ring-offset-2"
    >
      <ExternalLink size={16} />
      Masuk lewat Portal
    </a>
  )

  return (
    <div className="min-h-screen flex flex-col items-center justify-center bg-[#2c3e50] px-4">
      {/* App name */}
      <div className="mb-6 text-center">
        <h1 className="text-3xl font-bold text-white tracking-wide">{config.appName}</h1>
        <p className="text-white/60 mt-1 text-sm">Masuk ke akun Anda</p>
      </div>

      {/* Login card */}
      <div className="bg-white rounded-xl shadow-lg p-8 w-full max-w-sm">
        {optionsLoading ? (
          <div className="flex items-center justify-center gap-2 py-6 text-sm text-gray-500">
            <Loader2 size={16} className="animate-spin" /> Memuat...
          </div>
        ) : directEnabled ? (
          <div className="space-y-4">
            <LoginForm />
            {portalAvailable && (
              <>
                <div className="flex items-center gap-3 text-xs text-gray-400">
                  <span className="h-px flex-1 bg-gray-200" />
                  atau
                  <span className="h-px flex-1 bg-gray-200" />
                </div>
                {portalButton}
              </>
            )}
          </div>
        ) : (
          <div className="space-y-4 text-center">
            <p className="text-sm text-gray-600">
              Login dengan username dan password dinonaktifkan. Masuk lewat Portal untuk membuka
              aplikasi ini.
            </p>
            {portalAvailable && portalButton}
          </div>
        )}
      </div>

      {/* Footer */}
      <p className="mt-6 text-white/40 text-xs text-center">v{config.version} &middot; {config.appName}</p>
    </div>
  )
}

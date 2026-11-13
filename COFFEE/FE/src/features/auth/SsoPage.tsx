import { useEffect, useState } from 'react'
import { useSearchParams } from 'react-router-dom'
import { AlertCircle, Loader2 } from 'lucide-react'

import { config, ROUTES } from '@/shared/constants'
import { appUrl, portalUrl } from '@/shared/utils/appUrl'

import { useLoginOptionsQuery, useSsoLoginMutation } from './auth.api'

// Kode SSO hanya berlaku sekali. React StrictMode (dev) menjalankan effect dua kali, jadi
// kode yang sudah dikirim dicatat di tingkat modul agar tidak ditukar dua kali.
const submittedCodes = new Set<string>()

// Halaman tujuan dari Portal: /sso?code=XXXX. Menukar kode menjadi sesi Coffee lalu masuk.
export function SsoPage() {
  const [params] = useSearchParams()
  const code = params.get('code')?.trim() ?? ''
  const { mutateAsync } = useSsoLoginMutation()
  const { data: options } = useLoginOptionsQuery()
  // Galat disimpan di state lokal (bukan state mutation): di React StrictMode (dev) komponen
  // di-unmount lalu mount ulang tepat setelah mutate dipanggil, sehingga observer mutation
  // kehilangan hasilnya. Promise dari mutateAsync tetap sampai.
  const [errorMessage, setErrorMessage] = useState<string | null>(null)

  useEffect(() => {
    if (!code || submittedCodes.has(code)) return
    submittedCodes.add(code)
    // Kode dihapus dari address bar supaya tidak tersimpan di riwayat / ikut tersalin.
    window.history.replaceState(null, '', appUrl(ROUTES.SSO))
    mutateAsync({ code, device_info: 'web' }).catch((e: unknown) => {
      setErrorMessage(e instanceof Error ? e.message : 'Login lewat Portal gagal')
    })
  }, [code, mutateAsync])

  const failed = errorMessage !== null || (!code && submittedCodes.size === 0)
  const message = errorMessage ?? 'Kode login tidak ditemukan. Silakan buka aplikasi ini dari Portal.'
  const directEnabled = options ? options.direct_login_enabled : true

  return (
    <div className="min-h-screen flex flex-col items-center justify-center bg-[#2c3e50] px-4">
      <div className="mb-6 text-center">
        <h1 className="text-3xl font-bold text-white tracking-wide">{config.appName}</h1>
      </div>

      <div className="bg-white rounded-xl shadow-lg p-8 w-full max-w-sm text-center">
        {failed ? (
          <div className="space-y-4" role="alert">
            <AlertCircle size={32} className="mx-auto text-red-500" />
            <p className="text-sm text-gray-700">{message}</p>
            <div className="space-y-2">
              {portalUrl && (
                <a
                  href={portalUrl}
                  className="block w-full rounded-md bg-[#2c3e50] px-4 py-2 text-sm font-medium text-white hover:bg-[#34495e]"
                >
                  Kembali ke Portal
                </a>
              )}
              {directEnabled && (
                <a
                  href={appUrl(ROUTES.LOGIN)}
                  className="block w-full rounded-md border border-gray-300 px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50"
                >
                  Login dengan username
                </a>
              )}
            </div>
          </div>
        ) : (
          <div className="flex flex-col items-center gap-3 py-4 text-gray-600" role="status">
            <Loader2 size={28} className="animate-spin text-[#2c3e50]" />
            <p className="text-sm">Memasuki {config.appName}...</p>
          </div>
        )}
      </div>
    </div>
  )
}

import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'

import { api, authApi } from '@/services'
import { queryKeys } from '@/shared/constants'
import { ApiError } from '@/shared/types'
import { appUrl, portalUrl } from '@/shared/utils/appUrl'
import { useMenuStore } from '@/features/menu/menu.store'
import { getDefaultRoute } from '@/features/menu/getDefaultRoute'

import { useAuthStore } from './auth.store'
import type {
  AuthUser,
  LoginOptions,
  LoginRequest,
  LoginResponse,
  SsoRequest,
} from './auth.types'
import type { MenuItem } from '@/features/menu/menu.types'

// Dipakai bersama login langsung dan login lewat Portal: simpan sesi, muat menu, lalu
// pindah (reload penuh) ke halaman tujuan agar state aplikasi bersih untuk sesi baru.
function useCompleteLogin() {
  const setSession = useAuthStore((s) => s.setSession)
  const setMenus = useMenuStore((s) => s.setMenus)

  return async (data: LoginResponse) => {
    const user: AuthUser = {
      id: data.user.id,
      username: data.user.username,
      fullName: data.user.full_name,
      roleId: data.user.role_id,
      roleName: data.user.role_name,
      apps: 'web',
    }

    setSession({
      accessToken: data.token,
      refreshToken: data.refresh_token,
      expiresAt: data.expires_at,
      user,
    })

    // Fetch menu akses langsung setelah login — dipakai juga untuk menentukan halaman tujuan
    let menus: MenuItem[] = []
    try {
      menus = await api.post<MenuItem[]>('/menus/my', {})
      setMenus(menus)
    } catch {
      // Menu gagal diambil tidak menghentikan login
    }

    window.location.href = appUrl(getDefaultRoute(menus))
  }
}

export function useLoginMutation() {
  const completeLogin = useCompleteLogin()

  return useMutation({
    mutationFn: (payload: LoginRequest) => authApi.post<LoginResponse>('/auth/login', payload),
    onSuccess: completeLogin,
    onError: (error: Error) => {
      toast.error(error.message)
    },
  })
}

// Login lewat Portal: menukar kode sekali pakai (dari URL /sso?code=...) menjadi sesi Coffee.
export function useSsoLoginMutation() {
  const completeLogin = useCompleteLogin()

  return useMutation({
    mutationFn: (payload: SsoRequest) => authApi.post<LoginResponse>('/auth/sso', payload),
    onSuccess: completeLogin,
  })
}

// Cara login yang tersedia (publik). Bila gagal diambil, halaman login tetap menampilkan form
// login langsung agar tidak ada yang terkunci.
export function useLoginOptionsQuery() {
  return useQuery({
    queryKey: queryKeys.auth.loginOptions(),
    queryFn: () => authApi.get<LoginOptions>('/auth/login-options'),
    staleTime: 0,
    retry: false,
  })
}

// 401 saat logout = sesi di BE sudah berakhir; hasil akhirnya sama dengan logout berhasil.
function isSessionGone(error: unknown): boolean {
  return error instanceof ApiError && error.statusCode === 401
}

export function useLogoutMutation() {
  const clearSession = useAuthStore((s) => s.clearSession)
  const clearMenus = useMenuStore((s) => s.clearMenus)
  const queryClient = useQueryClient()

  const finishLogout = () => {
    clearSession()
    clearMenus()
    queryClient.clear()
    // Bila Portal terpasang, keluar dari aplikasi mengembalikan user ke Portal.
    if (portalUrl) window.location.href = portalUrl
  }

  return useMutation({
    mutationFn: () => api.post<void>('/auth/logout'),
    // Best-effort retry dulu sebelum menyerah — supaya session di BE (DeleteSessionByToken)
    // tetap kemungkinan besar terhapus walau ada error jaringan sesaat. 401 tidak diulang:
    // sesi di BE sudah tidak ada (mis. tergantikan login di perangkat lain), jadi sudah keluar.
    retry: (failureCount, error) => failureCount < 2 && !isSessionGone(error),
    retryDelay: (attempt) => Math.min(1000 * 2 ** attempt, 5000),

    onSuccess: () => {
      finishLogout()
      toast.success('Logout berhasil')
    },

    onError: (error: Error) => {
      // Setelah retry di atas tetap gagal, sesi lokal tetap dibersihkan agar user tidak
      // terjebak di UI (mis. benar-benar offline). Session di BE mungkin belum terhapus
      // sampai token tersebut kedaluwarsa secara alami.
      finishLogout()
      if (isSessionGone(error)) {
        toast.success('Logout berhasil')
      } else {
        toast.error('Logout dari server gagal, sesi lokal tetap dibersihkan')
      }
    },
  })
}

export function useGetMeQuery() {
  const isAuthenticated = useAuthStore((s) => s.isAuthenticated)

  return useQuery({
    queryKey: queryKeys.auth.profile(),
    queryFn: () => api.post<AuthUser>('/auth/me', {}),
    enabled: isAuthenticated,
  })
}

// Menyalakan/mematikan login langsung (hanya admin; BE menolak role lain).
export function useSetDirectLoginMutation() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (enabled: boolean) => api.post<void>('/auth/direct-login', { enabled }),
    onSuccess: (_data, enabled) => {
      void queryClient.invalidateQueries({ queryKey: queryKeys.auth.loginOptions() })
      toast.success(enabled ? 'Login langsung dinyalakan' : 'Login langsung dimatikan')
    },
    onError: (error: Error) => toast.error(error.message),
  })
}

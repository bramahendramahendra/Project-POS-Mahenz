import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'

import { api, authApi } from '@/services'
import { queryKeys } from '@/shared/constants'
import { useMenuStore } from '@/features/menu/menu.store'
import { getDefaultRoute } from '@/features/menu/getDefaultRoute'

import { useAuthStore } from './auth.store'
import type { AuthUser, LoginRequest, LoginResponse } from './auth.types'
import type { MenuItem } from '@/features/menu/menu.types'

export function useLoginMutation() {
  const setSession = useAuthStore((s) => s.setSession)
  const setMenus = useMenuStore((s) => s.setMenus)

  return useMutation({
    mutationFn: (payload: LoginRequest) => authApi.post<LoginResponse>('/auth/login', payload),

    onSuccess: async (data) => {
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

      window.location.href = getDefaultRoute(menus)
    },

    onError: (error: Error) => {
      toast.error(error.message)
    },
  })
}

export function useLogoutMutation() {
  const clearSession = useAuthStore((s) => s.clearSession)
  const clearMenus = useMenuStore((s) => s.clearMenus)
  const queryClient = useQueryClient()

  const finishLogout = () => {
    clearSession()
    clearMenus()
    queryClient.clear()
  }

  return useMutation({
    mutationFn: () => api.post<void>('/auth/logout'),
    // Best-effort retry dulu sebelum menyerah — supaya session di BE (DeleteSessionByToken)
    // tetap kemungkinan besar terhapus walau ada error jaringan sesaat.
    retry: 2,
    retryDelay: (attempt) => Math.min(1000 * 2 ** attempt, 5000),

    onSuccess: () => {
      finishLogout()
      toast.success('Logout berhasil')
    },

    onError: () => {
      // Setelah retry di atas tetap gagal, sesi lokal tetap dibersihkan agar user tidak
      // terjebak di UI (mis. benar-benar offline). Session di BE mungkin belum terhapus
      // sampai token tersebut kedaluwarsa secara alami.
      finishLogout()
      toast.error('Logout dari server gagal, sesi lokal tetap dibersihkan')
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

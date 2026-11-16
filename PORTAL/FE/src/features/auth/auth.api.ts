import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'

import { api, authApi } from '@/services'
import { queryKeys } from '@/shared/constants'
import { ROUTES } from '@/shared/constants/routes'

import { useAuthStore } from './auth.store'
import type {
  AuthUser,
  ChangePasswordRequest,
  LoginRequest,
  LoginResponse,
  MeResponse,
} from './auth.types'

export function useLoginMutation() {
  const setSession = useAuthStore((s) => s.setSession)

  return useMutation({
    mutationFn: (payload: LoginRequest) => authApi.post<LoginResponse>('/auth/login', payload),

    onSuccess: (data) => {
      const user: AuthUser = {
        id: data.user.id,
        username: data.user.username,
        fullName: data.user.full_name,
        roleId: data.user.role_id,
        roleName: data.user.role_name,
        mustChangePassword: data.user.must_change_password,
        apps: 'web',
      }

      setSession({
        accessToken: data.token,
        refreshToken: data.refresh_token,
        expiresAt: data.expires_at,
        user,
      })

      // Muat ulang halaman tujuan agar seluruh state aplikasi bersih untuk sesi baru.
      window.location.href = user.mustChangePassword ? ROUTES.CHANGE_PASSWORD : ROUTES.HOME
    },
  })
}

export function useLogoutMutation() {
  const clearSession = useAuthStore((s) => s.clearSession)
  const queryClient = useQueryClient()

  const finishLogout = () => {
    clearSession()
    queryClient.clear()
  }

  return useMutation({
    mutationFn: () => api.post<void>('/auth/logout'),
    // Best-effort retry dulu sebelum menyerah — supaya sesi di BE tetap
    // kemungkinan besar terhapus walau ada error jaringan sesaat.
    retry: 2,
    retryDelay: (attempt) => Math.min(1000 * 2 ** attempt, 5000),

    onSuccess: () => {
      finishLogout()
      toast.success('Anda telah keluar')
    },

    onError: () => {
      // Sesi lokal tetap dibersihkan agar user tidak terjebak di UI.
      finishLogout()
      toast.error('Keluar dari server gagal, sesi lokal tetap dibersihkan')
    },
  })
}

// Profil terbaru dari BE User (keputusan 128): dipanggil sekali saat Portal dibuka (ProtectedRoute), lalu nama, role,
// dan status wajib ganti password di sesi tersimpan disegarkan, sehingga menu langsung sesuai bila admin mengubah role.
// Bila gagal, data sesi tersimpan tetap dipakai (401 tetap diurus api client: refresh atau keluar).
export function useGetMeQuery() {
  const isAuthenticated = useAuthStore((s) => s.isAuthenticated)
  const updateUser = useAuthStore((s) => s.updateUser)

  return useQuery({
    queryKey: queryKeys.auth.profile(),
    queryFn: async () => {
      const me = await api.post<MeResponse>('/auth/me', {})
      updateUser({
        fullName: me.full_name,
        roleId: me.role_id,
        roleName: me.role_name,
        mustChangePassword: me.must_change_password,
      })
      return me
    },
    enabled: isAuthenticated,
  })
}

export function useChangePasswordMutation() {
  const updateUser = useAuthStore((s) => s.updateUser)

  return useMutation({
    mutationFn: (payload: ChangePasswordRequest) => api.post<void>('/auth/change-password', payload),
    onSuccess: () => {
      updateUser({ mustChangePassword: false })
      toast.success('Password berhasil diubah')
    },
  })
}

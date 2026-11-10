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

export function useGetMeQuery() {
  const isAuthenticated = useAuthStore((s) => s.isAuthenticated)

  return useQuery({
    queryKey: queryKeys.auth.profile(),
    queryFn: () => api.post<MeResponse>('/auth/me', {}),
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

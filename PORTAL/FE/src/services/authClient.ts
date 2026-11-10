import axios from 'axios'

import { ApiError } from '@/shared/types'

// Instance khusus endpoint auth (login, refresh, register, dst).
// Sengaja tanpa interceptor refresh-token: endpoint-endpoint ini
// yang MENGHASILKAN token, jadi tidak boleh masuk ke alur auto-refresh
// milik apiClient (lihat api.client.ts) untuk menghindari redirect/reload
// tak terduga saat kredensial salah atau refresh token invalid.
const authClient = axios.create({
  baseURL: import.meta.env.VITE_API_URL,
  headers: { 'Content-Type': 'application/json' },
  timeout: 30_000,
})

authClient.interceptors.response.use(
  (response) => {
    const body = response.data as { status: boolean; message: string; data: unknown }
    if (body.status === false) {
      throw new ApiError(body.message, response.status)
    }
    response.data = body.data
    return response
  },
  (error) => {
    const message: string = error.response?.data?.message ?? error.message ?? 'Terjadi kesalahan'
    const statusCode: number = error.response?.status ?? 0
    return Promise.reject(new ApiError(message, statusCode))
  }
)

export { authClient }

export const authApi = {
  post: <T>(url: string, data?: unknown) => authClient.post<T>(url, data).then((r) => r.data),
}

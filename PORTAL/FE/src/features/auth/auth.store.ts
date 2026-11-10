import { create } from 'zustand'
import { persist } from 'zustand/middleware'

import type { AuthState, SetSessionPayload } from './auth.types'

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      accessToken: null,
      refreshToken: null,
      expiresAt: null,
      user: null,
      isAuthenticated: false,

      setSession: (payload: SetSessionPayload) =>
        set({
          accessToken: payload.accessToken,
          refreshToken: payload.refreshToken,
          expiresAt: payload.expiresAt,
          user: payload.user,
          isAuthenticated: true,
        }),

      updateUser: (patch) =>
        set((state) => (state.user ? { user: { ...state.user, ...patch } } : state)),

      clearSession: () =>
        set({
          accessToken: null,
          refreshToken: null,
          expiresAt: null,
          user: null,
          isAuthenticated: false,
        }),
    }),
    // Format nama key browser: <app>.<nama>
    { name: 'portal.auth-session' }
  )
)

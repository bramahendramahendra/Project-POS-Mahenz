/* eslint-disable react-refresh/only-export-components */
import { lazy } from 'react'
import { createBrowserRouter, Navigate } from 'react-router-dom'

import { ChangePasswordPage, LoginPage, ProtectedRoute, RootRedirect } from '@/features/auth'
import { ROLES } from '@/shared/constants/roles'
import { ROUTES } from '@/shared/constants/routes'

import { LazyRoute } from './LazyRoute'

const HomePage = lazy(() => import('@/features/home/HomePage').then((m) => ({ default: m.HomePage })))
const AdminUsersPage = lazy(() =>
  import('@/features/admin/users/AdminUsersPage').then((m) => ({ default: m.AdminUsersPage }))
)
const AdminAppsPage = lazy(() =>
  import('@/features/admin/apps/AdminAppsPage').then((m) => ({ default: m.AdminAppsPage }))
)

export const router = createBrowserRouter([
  { path: '/', element: <RootRedirect /> },
  { path: ROUTES.LOGIN, element: <LoginPage /> },

  // Halaman setelah login (semua role).
  {
    element: <ProtectedRoute />,
    children: [
      { path: ROUTES.HOME, element: <LazyRoute><HomePage /></LazyRoute> },
      { path: ROUTES.CHANGE_PASSWORD, element: <ChangePasswordPage /> },
    ],
  },

  // Khusus admin Portal.
  {
    element: <ProtectedRoute roles={[ROLES.ADMIN]} />,
    children: [
      { path: ROUTES.ADMIN_USERS, element: <LazyRoute><AdminUsersPage /></LazyRoute> },
      { path: ROUTES.ADMIN_APPS, element: <LazyRoute><AdminAppsPage /></LazyRoute> },
    ],
  },

  // 404
  { path: '*', element: <Navigate to="/" replace /> },
])

import type { Role } from '@/shared/types'

export interface AppUser {
  id: number
  username: string
  full_name: string
  role_id: number
  role_name: Role
  is_active: boolean
  created_at: string
}

export interface UserListFilter {
  page: number
  limit: number
  search: string
  role_id?: number
  is_active?: boolean
  sort_by?: string
  sort_order?: 'asc' | 'desc'
}

export interface CreateUserPayload {
  username: string
  password: string
  full_name: string
  role_id: number
}

export type UpdateUserPayload = Omit<CreateUserPayload, 'username' | 'password'>

export interface ChangePasswordPayload {
  password: string
}

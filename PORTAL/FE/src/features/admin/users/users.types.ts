export type AdminUser = {
  id: number
  username: string
  full_name: string
  role_id: number
  role_name: string
  is_active: boolean
  must_change_password: boolean
  created_at: string
}

export interface UsersFilter {
  page: number
  limit: number
  search?: string
}

export interface CreateUserPayload {
  username: string
  password: string
  full_name: string
  role_id: number
}

export interface UpdateUserPayload {
  full_name: string
  role_id: number
}

export interface RoleOption {
  id: number
  display_name: string
}

export interface UserAccessItem {
  app_id: number
  code: string
  name: string
  is_active: boolean
  granted: boolean
}

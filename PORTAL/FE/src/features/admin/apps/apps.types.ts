export type AdminApp = {
  id: number
  code: string
  name: string
  description: string
  path: string
  icon: string
  sort_order: number
  is_active: boolean
  has_key: boolean
  user_count: number
  created_at: string
}

export interface CreateAppPayload {
  code: string
  name: string
  description: string
  path: string
  icon: string
  sort_order: number
}

export type UpdateAppPayload = Omit<CreateAppPayload, 'code'>

export interface CreateAppResult {
  app: AdminApp
  app_key: string
}

export interface RotateKeyResult {
  app_key: string
}

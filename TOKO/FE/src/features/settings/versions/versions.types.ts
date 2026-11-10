import type { Platform } from '@/shared/types'

export interface AppVersion {
  id: number
  platform: Platform
  version: string
  download_url: string
  release_notes: string
  is_mandatory: boolean
  is_latest: boolean
  created_at: string
}

export interface CreateAppVersionPayload {
  version: string
  download_url: string
  release_notes?: string
  is_mandatory: boolean
}

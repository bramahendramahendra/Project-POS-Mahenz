import type { Platform } from '@/shared/types'

export interface StoreProfile {
  name: string
  address?: string
  phone?: string
  email?: string
  logo_url?: string
  tax_default?: number
}

export interface PrinterSettings {
  paper_size: '58mm' | '80mm'
  receipt_header: string
  receipt_footer: string
  show_logo: boolean
  auto_print: boolean
}

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

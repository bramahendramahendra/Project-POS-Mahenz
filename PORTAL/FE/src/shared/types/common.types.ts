export interface SelectOption {
  label: string
  value: string | number
}

export interface DateRangeFilter {
  start_date: string
  end_date: string
}

export type SortOrder = 'asc' | 'desc'

export type Platform = 'web' | 'desktop' | 'android'

// Role Portal hanya dua: admin dan user (lihat ROLES di shared/constants).
export type Role = string

import type { Dayjs } from 'dayjs'

import { getWIBNow } from '@/shared/utils/date'

export function getGreeting(date: Dayjs = getWIBNow()): string {
  const hour = date.hour()
  if (hour < 10) return 'Selamat pagi'
  if (hour < 15) return 'Selamat siang'
  if (hour < 18) return 'Selamat sore'
  return 'Selamat malam'
}

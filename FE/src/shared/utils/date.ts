import dayjs from 'dayjs'
import utc from 'dayjs/plugin/utc'
import timezone from 'dayjs/plugin/timezone'

dayjs.extend(utc)
dayjs.extend(timezone)

export const APP_TIMEZONE = 'Asia/Jakarta'

/** Mengembalikan waktu sekarang sebagai objek dayjs yang sudah di-set ke WIB (Asia/Jakarta). */
export function getWIBNow() {
  return dayjs().tz(APP_TIMEZONE)
}

/** Mem-parse value ke objek dayjs yang di-anchor ke WIB, dipakai untuk parsing DAN formatting. */
function toWIB(date: string | Date) {
  return dayjs(date).tz(APP_TIMEZONE)
}

const MONTHS_SHORT = [
  'Jan',
  'Feb',
  'Mar',
  'Apr',
  'Mei',
  'Jun',
  'Jul',
  'Ags',
  'Sep',
  'Okt',
  'Nov',
  'Des',
]

export function formatDate(date: string | Date): string {
  const d = toWIB(date)
  return `${d.date()} ${MONTHS_SHORT[d.month()]} ${d.year()}`
}

export function formatDateShort(date: string | Date): string {
  const d = toWIB(date)
  return `${d.date()} ${MONTHS_SHORT[d.month()]}`
}

export function formatDateTime(date: string | Date): string {
  const d = toWIB(date)
  const hours = String(d.hour()).padStart(2, '0')
  const minutes = String(d.minute()).padStart(2, '0')
  return `${formatDate(date)}, ${hours}:${minutes}`
}

export function formatRelative(date: string | Date): string {
  const d = toWIB(date)
  const now = getWIBNow()
  const diffMs = now.valueOf() - d.valueOf()
  const diffSeconds = Math.floor(diffMs / 1000)
  const diffMinutes = Math.floor(diffSeconds / 60)
  const diffHours = Math.floor(diffMinutes / 60)
  const diffDays = Math.floor(diffHours / 24)

  if (diffSeconds < 60) return 'baru saja'
  if (diffMinutes < 60) return `${diffMinutes} menit yang lalu`
  if (diffHours < 24) return `${diffHours} jam yang lalu`
  if (diffDays === 1) return 'kemarin'
  return formatDate(date)
}

export function toISODate(date: Date): string {
  return toWIB(date).format('YYYY-MM-DD')
}

export function todayStr(): string {
  return getWIBNow().format('YYYY-MM-DD')
}

export function monthStart(): string {
  return getWIBNow().startOf('month').format('YYYY-MM-DD')
}

export function weekStart(): string {
  const now = getWIBNow()
  const dayOfWeek = now.day() // 0 = Minggu ... 6 = Sabtu
  return now.subtract(dayOfWeek - 1, 'day').format('YYYY-MM-DD')
}

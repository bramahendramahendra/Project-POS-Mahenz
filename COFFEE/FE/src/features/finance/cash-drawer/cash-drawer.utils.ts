/** Kas tertunda: dibuka sebelum hari ini (waktu perangkat) dan belum ditutup (tahap 4.9, keputusan 111). */
export function isPendingOpenTime(openTime?: string): boolean {
  if (!openTime) return false
  const d = new Date(openTime)
  const today = new Date()
  const key = (x: Date) => x.getFullYear() * 10000 + (x.getMonth() + 1) * 100 + x.getDate()
  return key(d) < key(today)
}

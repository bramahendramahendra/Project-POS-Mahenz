// Persentase perubahan current vs previous (laporan: Ringkasan Bisnis, Laporan Penjualan).
// null = tidak bisa dibandingkan (periode sebelumnya nol atau negatif).
export function pctChange(current: number, previous: number): number | null {
  if (previous <= 0) return null
  return ((current - previous) / previous) * 100
}

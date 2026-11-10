// Alamat aplikasi dengan memperhitungkan base path (mis. '/toko/' saat disajikan di bawah Portal).
// Dipakai untuk navigasi penuh (window.location) yang tidak melewati React Router.

const BASE = import.meta.env.BASE_URL // selalu berakhiran '/', mis. '/' atau '/toko/'

/** Basename untuk React Router; undefined bila aplikasi ada di root. */
export const routerBasename: string | undefined = BASE === '/' ? undefined : BASE.replace(/\/$/, '')

/** Mengubah path internal ('/kasir') menjadi alamat penuh dengan base path ('/toko/kasir'). */
export function appUrl(path: string): string {
  const clean = path.startsWith('/') ? path : `/${path}`
  return BASE.replace(/\/$/, '') + clean
}

/**
 * Alamat Portal (untuk tombol "Masuk lewat Portal" / "Kembali ke Portal").
 * Kosong = Portal belum dipasang untuk build ini, sehingga fitur terkait disembunyikan.
 */
export const portalUrl: string = (import.meta.env.VITE_PORTAL_URL as string | undefined) ?? ''

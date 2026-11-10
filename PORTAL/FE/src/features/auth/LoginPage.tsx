import { Navigate } from 'react-router-dom'
import { BarChart3, Coffee, ShieldCheck, Store } from 'lucide-react'

import { BrandMark } from '@/shared/components/layouts'
import { config, ROUTES } from '@/shared/constants'

import { LoginForm } from './components/LoginForm'
import { useAuthStore } from './auth.store'

const SHOWCASE = [
  {
    icon: Store,
    title: 'POS Toko',
    desc: 'Kasir, stok, pembelian, dan laporan toko.',
    status: 'Tersedia',
    delay: '[animation-delay:150ms]',
  },
  {
    icon: Coffee,
    title: 'POS Coffee',
    desc: 'Pesanan, meja, dan stok bahan kedai kopi.',
    status: 'Segera',
    delay: '[animation-delay:300ms]',
  },
  {
    icon: BarChart3,
    title: 'Satu akun',
    desc: 'Login sekali, pilih aplikasi sesuai izin Anda.',
    status: 'SSO',
    delay: '[animation-delay:450ms]',
  },
] as const

export function LoginPage() {
  const isAuthenticated = useAuthStore((s) => s.isAuthenticated)
  const mustChangePassword = useAuthStore((s) => s.user?.mustChangePassword ?? false)

  if (isAuthenticated) {
    return <Navigate to={mustChangePassword ? ROUTES.CHANGE_PASSWORD : ROUTES.HOME} replace />
  }

  return (
    <div className="grid min-h-screen lg:grid-cols-[1.1fr_1fr]">
      {/* Panel kiri: branding */}
      <aside className="portal-hero relative isolate hidden overflow-hidden lg:flex lg:flex-col lg:justify-between lg:p-12">
        <div className="portal-grid-pattern absolute inset-0 -z-10" aria-hidden="true" />
        <div
          className="absolute -left-24 top-1/3 -z-10 h-72 w-72 animate-float-slow rounded-full bg-indigo-500/30 blur-3xl"
          aria-hidden="true"
        />
        <div
          className="absolute -right-16 bottom-10 -z-10 h-80 w-80 animate-float rounded-full bg-fuchsia-500/25 blur-3xl"
          aria-hidden="true"
        />

        <BrandMark tone="light" className="animate-fade-up" />

        <div className="max-w-xl">
          <h1 className="animate-fade-up text-5xl font-semibold leading-[1.1] tracking-tight text-white [animation-delay:80ms]">
            Satu pintu untuk <span className="text-gradient">semua aplikasi</span> Anda.
          </h1>
          <p className="mt-5 max-w-md animate-fade-up text-base leading-relaxed text-white/65 [animation-delay:160ms]">
            Masuk sekali, lalu buka POS Toko, POS Coffee, dan aplikasi lain tanpa login berulang.
          </p>

          <ul className="mt-10 space-y-3">
            {SHOWCASE.map(({ icon: Icon, title, desc, status, delay }) => (
              <li
                key={title}
                className={`glass flex animate-fade-up items-center gap-4 rounded-2xl px-5 py-4 ${delay}`}
              >
                <span className="flex h-11 w-11 shrink-0 items-center justify-center rounded-xl bg-white/12 text-white">
                  <Icon size={20} />
                </span>
                <span className="min-w-0 flex-1">
                  <span className="block text-sm font-semibold text-white">{title}</span>
                  <span className="block text-sm text-white/60">{desc}</span>
                </span>
                <span className="rounded-full bg-white/12 px-2.5 py-1 text-[11px] font-medium uppercase tracking-wider text-white/80">
                  {status}
                </span>
              </li>
            ))}
          </ul>
        </div>

        <p className="flex items-center gap-2 text-xs text-white/45">
          <ShieldCheck size={14} /> Akses diatur oleh admin Portal &middot; v{config.version}
        </p>
      </aside>

      {/* Panel kanan: form */}
      <main className="relative flex items-center justify-center bg-background px-5 py-12 sm:px-8">
        <div
          className="pointer-events-none absolute inset-x-0 top-0 -z-0 h-64 bg-gradient-to-b from-indigo-100/70 to-transparent lg:hidden"
          aria-hidden="true"
        />
        <div className="relative z-10 w-full max-w-md animate-fade-up">
          <BrandMark className="mb-10 lg:hidden" />

          <div className="rounded-3xl border border-border/70 bg-white p-8 shadow-soft sm:p-10">
            <h2 className="text-2xl font-semibold tracking-tight text-foreground">
              Selamat datang kembali
            </h2>
            <p className="mt-1.5 text-sm text-muted-foreground">
              Masuk untuk melanjutkan ke aplikasi Anda.
            </p>

            <div className="mt-8">
              <LoginForm />
            </div>
          </div>

          <p className="mt-6 text-center text-xs text-muted-foreground">
            Lupa password? Hubungi admin Portal untuk mengatur ulang.
          </p>
        </div>
      </main>
    </div>
  )
}

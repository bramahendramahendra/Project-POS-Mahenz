import { AppWindow, ArrowUpRight, Coffee, Loader2, Store, type LucideIcon } from 'lucide-react'

import { cn } from '@/shared/utils'

import type { PortalApp } from '../home.types'

// Ikon dipilih dari kolom `icon` aplikasi; nilai yang tidak dikenal memakai ikon umum.
const ICONS: Record<string, LucideIcon> = {
  store: Store,
  coffee: Coffee,
}

// Warna aksen per kartu (bergilir) agar daftar tidak monoton.
const ACCENTS = [
  'from-indigo-500 to-violet-500',
  'from-amber-500 to-orange-500',
  'from-emerald-500 to-teal-500',
  'from-sky-500 to-blue-500',
  'from-fuchsia-500 to-pink-500',
] as const

interface AppCardProps {
  app: PortalApp
  index: number
  isLaunching: boolean
  disabled: boolean
  onOpen: (app: PortalApp) => void
}

export function AppCard({ app, index, isLaunching, disabled, onOpen }: AppCardProps) {
  const Icon = ICONS[app.icon] ?? AppWindow
  const accent = ACCENTS[index % ACCENTS.length]

  return (
    <button
      type="button"
      onClick={() => onOpen(app)}
      disabled={disabled}
      aria-label={`Buka ${app.name}`}
      style={{ animationDelay: `${index * 70}ms` }}
      className={cn(
        'group relative flex animate-fade-up flex-col items-start gap-5 overflow-hidden rounded-3xl border border-border/70 bg-white p-6 text-left shadow-sm transition-all duration-200',
        'hover:-translate-y-1 hover:border-brand/30 hover:shadow-glow',
        'focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2',
        'disabled:cursor-not-allowed disabled:opacity-60 disabled:hover:translate-y-0 disabled:hover:shadow-sm'
      )}
    >
      <div
        className={cn(
          'pointer-events-none absolute -right-10 -top-10 h-36 w-36 rounded-full bg-gradient-to-br opacity-10 blur-2xl transition-opacity group-hover:opacity-25',
          accent
        )}
        aria-hidden="true"
      />

      <span
        className={cn(
          'flex h-14 w-14 items-center justify-center rounded-2xl bg-gradient-to-br text-white shadow-lg',
          accent
        )}
      >
        <Icon size={26} />
      </span>

      <span className="block w-full">
        <span className="block text-lg font-semibold tracking-tight text-foreground">{app.name}</span>
        {app.description && (
          <span className="mt-1 block text-sm leading-relaxed text-muted-foreground">
            {app.description}
          </span>
        )}
      </span>

      <span className="mt-auto flex items-center gap-1.5 text-sm font-medium text-brand">
        {isLaunching ? (
          <>
            <Loader2 size={16} className="animate-spin" /> Membuka...
          </>
        ) : (
          <>
            Buka aplikasi
            <ArrowUpRight
              size={16}
              className="transition-transform group-hover:-translate-y-0.5 group-hover:translate-x-0.5"
            />
          </>
        )}
      </span>
    </button>
  )
}

export function AppCardSkeleton() {
  return (
    <div
      className="flex animate-pulse flex-col gap-5 rounded-3xl border border-border/70 bg-white p-6"
      aria-hidden="true"
    >
      <div className="h-14 w-14 rounded-2xl bg-secondary" />
      <div className="space-y-2">
        <div className="h-5 w-32 rounded bg-secondary" />
        <div className="h-4 w-full rounded bg-secondary" />
        <div className="h-4 w-2/3 rounded bg-secondary" />
      </div>
      <div className="h-4 w-24 rounded bg-secondary" />
    </div>
  )
}

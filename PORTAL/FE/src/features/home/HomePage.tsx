import { AlertCircle, LayoutGrid, RefreshCw } from 'lucide-react'

import { useAuth } from '@/features/auth/hooks/useAuth'

import { AppCard, AppCardSkeleton } from './components/AppCard'
import { useLaunchAppMutation, useMyAppsQuery } from './home.api'

function firstName(fullName: string): string {
  return fullName.trim().split(/\s+/)[0] ?? ''
}

export function HomePage() {
  const { user } = useAuth()
  const { data: apps, isLoading, isError, error, refetch, isFetching } = useMyAppsQuery()
  const { mutate: launch, isPending, variables } = useLaunchAppMutation()

  return (
    <div>
      <div className="animate-fade-up">
        <p className="text-sm font-medium text-brand">Portal</p>
        <h1 className="mt-1 text-3xl font-semibold tracking-tight">
          Halo, {firstName(user?.fullName ?? '')}
        </h1>
        <p className="mt-1.5 text-muted-foreground">Pilih aplikasi yang ingin Anda buka.</p>
      </div>

      <section className="mt-8" aria-live="polite">
        {isLoading && (
          <div className="grid gap-5 sm:grid-cols-2 lg:grid-cols-3">
            <AppCardSkeleton />
            <AppCardSkeleton />
          </div>
        )}

        {isError && (
          <div
            role="alert"
            className="flex flex-col items-center rounded-3xl border border-red-200 bg-red-50/70 px-6 py-14 text-center"
          >
            <AlertCircle size={28} className="text-red-500" />
            <h2 className="mt-3 text-lg font-semibold text-red-800">Daftar aplikasi gagal dimuat</h2>
            <p className="mt-1 max-w-sm text-sm text-red-700/80">
              {error instanceof Error ? error.message : 'Terjadi kesalahan.'}
            </p>
            <button
              type="button"
              onClick={() => refetch()}
              disabled={isFetching}
              className="mt-5 inline-flex h-10 items-center gap-2 rounded-xl bg-white px-4 text-sm font-medium text-red-700 shadow-sm ring-1 ring-red-200 transition-colors hover:bg-red-50 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-red-400 disabled:opacity-60"
            >
              <RefreshCw size={16} className={isFetching ? 'animate-spin' : undefined} />
              Coba lagi
            </button>
          </div>
        )}

        {apps && apps.length > 0 && (
          <div className="grid gap-5 sm:grid-cols-2 lg:grid-cols-3">
            {apps.map((app, index) => (
              <AppCard
                key={app.id}
                app={app}
                index={index}
                isLaunching={isPending && variables === app.code}
                disabled={isPending}
                onOpen={(a) => launch(a.code)}
              />
            ))}
          </div>
        )}

        {apps && apps.length === 0 && (
          <div className="flex flex-col items-center justify-center rounded-3xl border border-dashed border-border bg-white/60 px-6 py-16 text-center">
            <span className="flex h-14 w-14 items-center justify-center rounded-2xl bg-accent text-accent-foreground">
              <LayoutGrid size={26} />
            </span>
            <h2 className="mt-4 text-lg font-semibold">Belum ada aplikasi untuk Anda</h2>
            <p className="mt-1 max-w-sm text-sm text-muted-foreground">
              Akun Anda belum diberi akses ke aplikasi apa pun. Hubungi admin Portal untuk meminta akses.
            </p>
          </div>
        )}
      </section>
    </div>
  )
}

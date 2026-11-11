import { useState, type ReactNode } from 'react'
import { NavLink, useNavigate } from 'react-router-dom'
import {
  AppWindow,
  ChevronDown,
  KeyRound,
  LayoutGrid,
  LogOut,
  ShieldCheck,
  Users,
  type LucideIcon,
} from 'lucide-react'

import { useLogoutMutation } from '@/features/auth/auth.api'
import { useAuth } from '@/features/auth/hooks/useAuth'
import { ConfirmDialog } from '@/shared/components/ConfirmDialog'
import { Avatar, AvatarFallback } from '@/shared/components/ui/avatar'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/shared/components/ui/dropdown-menu'
import { ROLES, ROUTES } from '@/shared/constants'
import { cn } from '@/shared/utils'

import { BrandMark } from './BrandMark'

const NAV_ITEMS: { to: string; label: string; icon: LucideIcon }[] = [
  { to: ROUTES.HOME, label: 'Aplikasi', icon: LayoutGrid },
  { to: ROUTES.ADMIN_USERS, label: 'Kelola User', icon: Users },
  { to: ROUTES.ADMIN_APPS, label: 'Kelola Aplikasi', icon: AppWindow },
]

function initialsOf(name: string): string {
  const parts = name.trim().split(/\s+/).filter(Boolean)
  if (parts.length === 0) return '?'
  const first = parts[0]?.[0] ?? ''
  const last = parts.length > 1 ? (parts[parts.length - 1]?.[0] ?? '') : ''
  return (first + last).toUpperCase()
}

export function PortalLayout({ children }: { children: ReactNode }) {
  const navigate = useNavigate()
  const { user, mustChangePassword, isAdmin } = useAuth()
  const { mutate: logout, isPending } = useLogoutMutation()
  const [confirmOpen, setConfirmOpen] = useState(false)

  const roleLabel = user?.roleName === ROLES.ADMIN ? 'Admin' : 'User'

  return (
    <div className="min-h-screen bg-background">
      <header
        className="sticky top-0 z-30 border-b border-border/70 bg-white/80 backdrop-blur-md"
        style={{ height: 'var(--navbar-height)' }}
      >
        <div className="mx-auto flex h-full max-w-6xl items-center justify-between px-4 sm:px-6">
          <div className="flex items-center gap-8">
            <BrandMark />

            {isAdmin && !mustChangePassword && (
              // Mulai lg (1024px): di tablet tegak menu header terlalu sempit, jadi masuk ke menu akun.
              <nav className="hidden items-center gap-1 lg:flex" aria-label="Navigasi utama">
                {NAV_ITEMS.map(({ to, label, icon: Icon }) => (
                  <NavLink
                    key={to}
                    to={to}
                    className={({ isActive }) =>
                      cn(
                        'flex items-center gap-2 rounded-lg px-3 py-2 text-sm font-medium transition-colors',
                        isActive
                          ? 'bg-accent text-accent-foreground'
                          : 'text-muted-foreground hover:bg-secondary hover:text-foreground'
                      )
                    }
                  >
                    <Icon size={16} />
                    {label}
                  </NavLink>
                ))}
              </nav>
            )}
          </div>

          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <button
                type="button"
                className="group flex items-center gap-3 rounded-full border border-transparent py-1 pl-1 pr-3 transition-colors hover:border-border hover:bg-white focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
                aria-label="Menu akun"
              >
                <Avatar className="h-9 w-9">
                  <AvatarFallback className="bg-gradient-to-br from-brand-light to-brand-violet text-xs font-semibold text-white">
                    {initialsOf(user?.fullName ?? '')}
                  </AvatarFallback>
                </Avatar>
                <span className="hidden text-left leading-tight sm:block">
                  <span className="block text-sm font-medium text-foreground">{user?.fullName}</span>
                  <span className="block text-xs text-muted-foreground">{roleLabel}</span>
                </span>
                <ChevronDown
                  size={16}
                  className="text-muted-foreground transition-transform group-data-[state=open]:rotate-180"
                />
              </button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end" className="w-60">
              <DropdownMenuLabel className="font-normal">
                <p className="text-sm font-medium text-foreground">{user?.fullName}</p>
                <p className="flex items-center gap-1 text-xs text-muted-foreground">
                  <ShieldCheck size={12} /> @{user?.username} &middot; {roleLabel}
                </p>
              </DropdownMenuLabel>
              <DropdownMenuSeparator />
              {isAdmin &&
                !mustChangePassword &&
                NAV_ITEMS.map(({ to, label, icon: Icon }) => (
                  <DropdownMenuItem key={to} onSelect={() => navigate(to)} className="lg:hidden">
                    <Icon size={16} className="mr-2" /> {label}
                  </DropdownMenuItem>
                ))}
              {!mustChangePassword && (
                <DropdownMenuItem onSelect={() => navigate(ROUTES.CHANGE_PASSWORD)}>
                  <KeyRound size={16} className="mr-2" /> Ganti password
                </DropdownMenuItem>
              )}
              <DropdownMenuItem
                onSelect={() => setConfirmOpen(true)}
                className="text-destructive focus:text-destructive"
              >
                <LogOut size={16} className="mr-2" /> Keluar
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      </header>

      <main className="mx-auto w-full max-w-6xl px-4 py-8 sm:px-6">{children}</main>

      <ConfirmDialog
        open={confirmOpen}
        onOpenChange={setConfirmOpen}
        title="Keluar dari Portal?"
        description="Anda perlu login lagi untuk membuka aplikasi."
        confirmLabel="Keluar"
        variant="destructive"
        isLoading={isPending}
        onConfirm={() => logout()}
      />
    </div>
  )
}

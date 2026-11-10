import { useId } from 'react'

import { cn } from '@/shared/utils'

interface BrandMarkProps {
  /** light: untuk latar gelap, dark: untuk latar terang */
  tone?: 'light' | 'dark'
  showText?: boolean
  className?: string
}

export function BrandMark({ tone = 'dark', showText = true, className }: BrandMarkProps) {
  // id unik per instance: gradient di dalam subtree display:none (mis. panel kiri login
  // di layar kecil) tidak bisa dirujuk instance lain yang terlihat.
  const gradId = `brand-grad-${useId().replace(/:/g, '')}`

  return (
    <div className={cn('flex items-center gap-3', className)}>
      <svg
        viewBox="0 0 64 64"
        className="h-9 w-9 shrink-0 drop-shadow-sm"
        aria-hidden="true"
        focusable="false"
      >
        <defs>
          <linearGradient id={gradId} x1="8" y1="6" x2="58" y2="60" gradientUnits="userSpaceOnUse">
            <stop stopColor="#6366f1" />
            <stop offset="1" stopColor="#7c3aed" />
          </linearGradient>
        </defs>
        <rect width="64" height="64" rx="16" fill={`url(#${gradId})`} />
        <rect x="14" y="14" width="15" height="15" rx="4" fill="#fff" />
        <rect x="35" y="14" width="15" height="15" rx="4" fill="#fff" fillOpacity=".7" />
        <rect x="14" y="35" width="15" height="15" rx="4" fill="#fff" fillOpacity=".7" />
        <rect x="35" y="35" width="15" height="15" rx="4" fill="#fff" fillOpacity=".4" />
      </svg>
      {showText && (
        <div className="leading-tight">
          <p
            className={cn(
              'text-base font-semibold tracking-tight',
              tone === 'light' ? 'text-white' : 'text-foreground'
            )}
          >
            Mahenz Portal
          </p>
          <p
            className={cn(
              'text-[11px] uppercase tracking-[0.18em]',
              tone === 'light' ? 'text-white/55' : 'text-muted-foreground'
            )}
          >
            Satu pintu, semua aplikasi
          </p>
        </div>
      )}
    </div>
  )
}

import { useEffect, useRef, useState } from 'react'
import { Clock } from 'lucide-react'

import { cn } from '@/shared/utils'

interface TimePickerInputProps {
  value?: string
  onChange?: (value: string) => void
  placeholder?: string
  className?: string
  hasError?: boolean
}

const HOURS = Array.from({ length: 24 }, (_, i) => i)
const MINUTES = Array.from({ length: 60 }, (_, i) => i)

const ITEM_HEIGHT = 36

export function TimePickerInput({
  value = '',
  onChange,
  placeholder = '--:--',
  className,
  hasError,
}: TimePickerInputProps) {
  const [open, setOpen] = useState(false)
  const containerRef = useRef<HTMLDivElement>(null)
  const hourRef = useRef<HTMLDivElement>(null)
  const minuteRef = useRef<HTMLDivElement>(null)

  const selectedHour = value ? parseInt(value.split(':')[0], 10) : null
  const selectedMinute = value ? parseInt(value.split(':')[1], 10) : null

  useEffect(() => {
    if (!open) return
    setTimeout(() => {
      if (hourRef.current && selectedHour !== null) {
        hourRef.current.scrollTop = selectedHour * ITEM_HEIGHT - ITEM_HEIGHT * 2
      }
      if (minuteRef.current && selectedMinute !== null) {
        minuteRef.current.scrollTop = selectedMinute * ITEM_HEIGHT - ITEM_HEIGHT * 2
      }
    }, 80)
  }, [open])

  useEffect(() => {
    const handleClickOutside = (e: MouseEvent) => {
      if (containerRef.current && !containerRef.current.contains(e.target as Node)) {
        setOpen(false)
      }
    }
    document.addEventListener('mousedown', handleClickOutside)
    return () => document.removeEventListener('mousedown', handleClickOutside)
  }, [])

  const handleSelect = (hour: number | null, minute: number | null) => {
    const h = hour !== null ? hour : (selectedHour ?? 0)
    const m = minute !== null ? minute : (selectedMinute ?? 0)
    onChange?.(`${String(h).padStart(2, '0')}:${String(m).padStart(2, '0')}`)
  }

  return (
    <div ref={containerRef} className="relative">
      <button
        type="button"
        onClick={() => setOpen((v) => !v)}
        className={cn(
          'flex h-10 w-full items-center justify-between rounded-md border border-input bg-background px-3 py-2 text-sm',
          'ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2',
          'hover:bg-accent/40 transition-colors',
          !value && 'text-muted-foreground',
          hasError && 'border-red-500',
          open && 'ring-2 ring-ring ring-offset-2',
          className
        )}
      >
        <span className="font-mono tracking-widest">{value || placeholder}</span>
        <Clock className="h-4 w-4 text-muted-foreground" />
      </button>

      {open && (
        <div className="mt-1.5 w-48 rounded-xl border border-border bg-popover shadow-md">
          <div className="px-3 pt-3 pb-1 text-xs font-medium text-muted-foreground text-center tracking-wide">
            Pilih Waktu
          </div>

          <div className="flex divide-x divide-border">
            {/* Kolom Jam */}
            <div className="flex-1 flex flex-col items-center">
              <span className="text-[10px] text-muted-foreground py-1">Jam</span>
              <div
                ref={hourRef}
                className="h-44 overflow-y-auto w-full"
                style={{ scrollbarWidth: 'none' }}
              >
                {HOURS.map((h) => (
                  <button
                    key={h}
                    type="button"
                    onClick={() => handleSelect(h, null)}
                    style={{ height: ITEM_HEIGHT }}
                    className={cn(
                      'w-full text-sm font-mono text-center transition-colors hover:bg-accent rounded-md',
                      selectedHour === h
                        ? 'bg-primary text-primary-foreground font-semibold hover:bg-primary'
                        : 'text-foreground'
                    )}
                  >
                    {String(h).padStart(2, '0')}
                  </button>
                ))}
              </div>
            </div>

            {/* Kolom Menit */}
            <div className="flex-1 flex flex-col items-center">
              <span className="text-[10px] text-muted-foreground py-1">Menit</span>
              <div
                ref={minuteRef}
                className="h-44 overflow-y-auto w-full"
                style={{ scrollbarWidth: 'none' }}
              >
                {MINUTES.map((m) => (
                  <button
                    key={m}
                    type="button"
                    onClick={() => handleSelect(null, m)}
                    style={{ height: ITEM_HEIGHT }}
                    className={cn(
                      'w-full text-sm font-mono text-center transition-colors hover:bg-accent rounded-md',
                      selectedMinute === m
                        ? 'bg-primary text-primary-foreground font-semibold hover:bg-primary'
                        : 'text-foreground'
                    )}
                  >
                    {String(m).padStart(2, '0')}
                  </button>
                ))}
              </div>
            </div>
          </div>

          <div className="p-2 border-t border-border">
            <button
              type="button"
              onClick={() => setOpen(false)}
              className="w-full text-xs text-center py-1.5 rounded-md bg-primary text-primary-foreground hover:bg-primary/90 transition-colors font-medium"
            >
              Selesai
            </button>
          </div>
        </div>
      )}
    </div>
  )
}

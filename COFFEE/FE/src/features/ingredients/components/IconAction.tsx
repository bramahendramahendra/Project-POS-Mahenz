import type { ReactNode } from 'react'

import { Button } from '@/shared/components/ui/button'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/shared/components/ui/tooltip'

interface IconActionProps {
  label: string
  onClick: () => void
  className?: string
  children: ReactNode
}

// Tombol ikon kolom aksi dengan tooltip dan aria-label.
export function IconAction({ label, onClick, className, children }: IconActionProps) {
  return (
    <Tooltip>
      <TooltipTrigger asChild>
        <Button variant="ghost" size="icon" aria-label={label} className={`h-7 w-7 text-gray-500 ${className ?? ''}`} onClick={onClick}>
          {children}
        </Button>
      </TooltipTrigger>
      <TooltipContent>{label}</TooltipContent>
    </Tooltip>
  )
}

import { type ReactNode } from 'react'
import { AlertTriangle, Loader2 } from 'lucide-react'

import {
  AlertDialog,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from '@/shared/components/ui/alert-dialog'
import { Button } from '@/shared/components/ui/button'

interface ExtendedConfirmDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  title: string
  description: ReactNode
  /** Konten tambahan di bawah description (mis. input konfirmasi ketik ulang teks) — tidak ada di ConfirmDialog biasa. */
  children?: ReactNode
  confirmLabel?: string
  cancelLabel?: string
  variant?: 'default' | 'destructive'
  isLoading?: boolean
  confirmDisabled?: boolean
  onConfirm: () => void
}

export function ExtendedConfirmDialog({
  open,
  onOpenChange,
  title,
  description,
  children,
  confirmLabel = 'Ya, Lanjutkan',
  cancelLabel = 'Batal',
  variant = 'default',
  isLoading,
  confirmDisabled,
  onConfirm,
}: ExtendedConfirmDialogProps) {
  const handleOpenChange = (val: boolean) => {
    if (isLoading) return
    onOpenChange(val)
  }

  return (
    <AlertDialog open={open} onOpenChange={handleOpenChange}>
      <AlertDialogContent>
        <AlertDialogHeader>
          <div className="flex items-center gap-3">
            <AlertTriangle
              size={22}
              className={variant === 'destructive' ? 'text-red-500' : 'text-yellow-500'}
            />
            <AlertDialogTitle>{title}</AlertDialogTitle>
          </div>
          <AlertDialogDescription className="pt-1">{description}</AlertDialogDescription>
        </AlertDialogHeader>

        {children}

        <AlertDialogFooter>
          <Button
            type="button"
            variant="outline"
            onClick={() => handleOpenChange(false)}
            disabled={isLoading}
          >
            {cancelLabel}
          </Button>
          <Button
            type="button"
            variant={variant === 'destructive' ? 'destructive' : 'default'}
            onClick={onConfirm}
            disabled={isLoading || confirmDisabled}
          >
            {isLoading && <Loader2 size={14} className="animate-spin" />}
            {confirmLabel}
          </Button>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  )
}

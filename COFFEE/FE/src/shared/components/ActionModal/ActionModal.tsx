import { type FormEventHandler, type ReactNode } from 'react'
import { Loader2 } from 'lucide-react'

import { Button } from '@/shared/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/shared/components/ui/dialog'
import { cn } from '@/shared/utils'

interface ActionModalProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  title: string
  description?: string
  /** Kelas untuk DialogContent (ukuran/tinggi) — pemanggil kontrol penuh, tidak dipaksa ke skala FormModal. */
  contentClassName?: string
  /** Override kelas DialogHeader (default border-b px-6 py-4) — beberapa modal pakai padding lebih rapat. */
  headerClassName?: string
  /** Override kelas DialogDescription — default: visible kalau `description` diisi, sr-only kalau tidak. */
  descriptionClassName?: string
  isLoading?: boolean
  children: ReactNode

  /** Bungkus children+footer di dalam <form> asli — dipakai saat butuh submit-on-Enter (mis. form pembayaran kasir). */
  asForm?: boolean
  onFormSubmit?: FormEventHandler<HTMLFormElement>

  /** Footer sepenuhnya custom (menggantikan footer default Batal/Submit) — dipakai kalau jumlah/label tombol tidak mengikuti pola submit-cancel biasa. */
  footer?: ReactNode
  hideFooter?: boolean

  /** Footer default (dipakai kalau `footer` & `hideFooter` tidak diisi) — dipanggil lewat onClick kalau bukan asForm, atau lewat submit native form kalau asForm. */
  onSubmit?: () => void
  submitLabel?: string
  cancelLabel?: string
  submitDisabled?: boolean
}

export function ActionModal({
  open,
  onOpenChange,
  title,
  description,
  contentClassName,
  headerClassName,
  descriptionClassName,
  isLoading,
  children,
  asForm,
  onFormSubmit,
  footer,
  hideFooter,
  onSubmit,
  submitLabel = 'Simpan',
  cancelLabel = 'Batal',
  submitDisabled,
}: ActionModalProps) {
  function handleOpenChange(val: boolean) {
    if (isLoading) return
    onOpenChange(val)
  }

  const defaultFooter = (
    <DialogFooter className="border-t px-6 py-4">
      <Button
        type="button"
        variant="outline"
        onClick={() => handleOpenChange(false)}
        disabled={isLoading}
      >
        {cancelLabel}
      </Button>
      <Button type={asForm ? 'submit' : 'button'} onClick={asForm ? undefined : onSubmit} disabled={isLoading || submitDisabled}>
        {isLoading && <Loader2 size={14} className="animate-spin" />}
        {submitLabel}
      </Button>
    </DialogFooter>
  )

  const resolvedFooter = footer ?? (hideFooter ? null : defaultFooter)

  const body = (
    <>
      {children}
      {resolvedFooter}
    </>
  )

  return (
    <Dialog open={open} onOpenChange={handleOpenChange}>
      <DialogContent
        className={cn('flex flex-col gap-0 p-0', contentClassName)}
        onInteractOutside={(e) => e.preventDefault()}
        onEscapeKeyDown={(e) => {
          if (isLoading) e.preventDefault()
        }}
      >
        <DialogHeader className={headerClassName ?? 'border-b px-6 py-4'}>
          <DialogTitle>{title}</DialogTitle>
          <DialogDescription className={descriptionClassName ?? (description ? '' : 'sr-only')}>
            {description ?? title}
          </DialogDescription>
        </DialogHeader>

        {asForm ? <form onSubmit={onFormSubmit}>{body}</form> : body}
      </DialogContent>
    </Dialog>
  )
}

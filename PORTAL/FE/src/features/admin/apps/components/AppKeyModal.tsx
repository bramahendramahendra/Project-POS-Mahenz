import { useState } from 'react'
import { Check, Copy, KeyRound, TriangleAlert } from 'lucide-react'

import { Button } from '@/shared/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/shared/components/ui/dialog'

interface AppKeyModalProps {
  /** Kunci asli. null/undefined = modal tertutup. Hanya hidup di memori selama modal terbuka. */
  appKey: string | null
  appName: string
  /** true bila kunci hasil rotasi (kunci lama sudah tidak berlaku) */
  isRotation?: boolean
  onClose: () => void
}

async function copyText(text: string): Promise<boolean> {
  try {
    await navigator.clipboard.writeText(text)
    return true
  } catch {
    // Fallback untuk konteks non-HTTPS / izin clipboard ditolak.
    const ta = document.createElement('textarea')
    ta.value = text
    ta.style.position = 'fixed'
    ta.style.opacity = '0'
    document.body.appendChild(ta)
    ta.select()
    const ok = document.execCommand('copy')
    document.body.removeChild(ta)
    return ok
  }
}

export function AppKeyModal({ appKey, appName, isRotation, onClose }: AppKeyModalProps) {
  // Menyimpan kunci yang sudah disalin; status "tersalin" otomatis hilang saat kunci berganti.
  const [copiedKey, setCopiedKey] = useState<string | null>(null)
  const copied = !!appKey && copiedKey === appKey

  const handleCopy = async () => {
    if (!appKey) return
    if (await copyText(appKey)) setCopiedKey(appKey)
  }

  return (
    <Dialog open={!!appKey} onOpenChange={(v) => !v && onClose()}>
      <DialogContent
        className="max-w-[560px]"
        onInteractOutside={(e) => e.preventDefault()}
        onEscapeKeyDown={(e) => e.preventDefault()}
      >
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2">
            <KeyRound size={18} className="text-brand" />
            {isRotation ? 'Kunci baru' : 'Kunci aplikasi'} — {appName}
          </DialogTitle>
          <DialogDescription>
            Salin dan simpan kunci ini sekarang. Kunci tidak akan ditampilkan lagi.
          </DialogDescription>
        </DialogHeader>

        <div className="space-y-3">
          <div className="flex items-start gap-2 rounded-lg border border-amber-200 bg-amber-50 px-3 py-2.5 text-sm text-amber-800">
            <TriangleAlert size={16} className="mt-0.5 shrink-0" />
            <span>
              {isRotation
                ? 'Kunci lama sudah tidak berlaku. Perbarui kunci di database aplikasi tujuan agar login SSO tetap berjalan.'
                : 'Masukkan kunci ini ke database aplikasi tujuan (tabel kunci) bersama kode aplikasinya.'}
            </span>
          </div>

          <div className="flex items-stretch gap-2">
            <code
              data-testid="app-key"
              className="block min-w-0 flex-1 select-all break-all rounded-lg border border-border bg-secondary/60 px-3 py-2.5 font-mono text-xs leading-relaxed"
            >
              {appKey}
            </code>
            <Button type="button" variant="outline" onClick={() => void handleCopy()} className="h-auto shrink-0">
              {copied ? <Check size={16} className="text-emerald-600" /> : <Copy size={16} />}
              {copied ? 'Tersalin' : 'Salin'}
            </Button>
          </div>
        </div>

        <DialogFooter>
          <Button type="button" onClick={onClose}>
            Saya sudah menyimpan kunci
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}

import { useRef, useState } from 'react'
import { Database, Download, RotateCcw, Upload } from 'lucide-react'
import { toast } from 'sonner'

import { RoleGuard } from '@/shared/components'
import { Button } from '@/shared/components/ui/button'

import {
  downloadBackupFile,
  useBackupListQuery,
  useCreateBackupMutation,
  useRestoreBackupMutation,
} from '../backup.api'
import type { BackupInfo } from '../backup.types'
import { RestoreConfirmDialog } from './RestoreConfirmDialog'

function formatDate(str: string): string {
  return new Date(str).toLocaleString('id-ID', {
    day: '2-digit',
    month: 'short',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}

export function BackupTab() {
  const { data, isLoading } = useBackupListQuery()
  const { mutate: createBackup, isPending: isCreating } = useCreateBackupMutation()
  const { mutate: restoreBackup, isPending: isRestoring } = useRestoreBackupMutation()
  const fileInputRef = useRef<HTMLInputElement>(null)

  const [restoreTarget, setRestoreTarget] = useState<BackupInfo | { filename: string; file: File } | null>(null)

  const files = data?.files ?? []

  const handleDownload = (filename: string) => {
    downloadBackupFile(filename).catch(() => toast.error('Gagal mengunduh file backup'))
  }

  const handleUploadClick = () => fileInputRef.current?.click()

  const handleFileSelected = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
    e.target.value = ''
    if (!file) return
    if (!file.name.toLowerCase().endsWith('.sql')) {
      toast.error('Hanya file .sql yang diizinkan')
      return
    }
    setRestoreTarget({ filename: file.name, file })
  }

  const handleConfirmRestore = () => {
    if (!restoreTarget) return
    const source = 'file' in restoreTarget ? { file: restoreTarget.file } : { filename: restoreTarget.filename }
    restoreBackup(source, {
      onSuccess: () => {
        toast.success('Restore berhasil dilakukan. Silakan login ulang.')
        setRestoreTarget(null)
      },
    })
  }

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between flex-wrap gap-3">
        <p className="text-sm text-gray-500 max-w-xl">
          Backup menyimpan seluruh data database saat ini ke file .sql. Restore akan{' '}
          <strong className="text-red-600">menimpa seluruh data yang ada sekarang</strong> — gunakan
          dengan sangat hati-hati.
        </p>
        <div className="flex gap-2">
          <RoleGuard menuKey="sistem.backup" action="can_delete">
            <input
              ref={fileInputRef}
              type="file"
              accept=".sql"
              className="hidden"
              onChange={handleFileSelected}
            />
            <Button variant="outline" className="gap-1.5" onClick={handleUploadClick}>
              <Upload size={14} />
              Restore dari File Lain
            </Button>
          </RoleGuard>
          <RoleGuard menuKey="sistem.backup" action="can_create">
            <Button className="gap-1.5" onClick={() => createBackup()} disabled={isCreating}>
              <Database size={14} />
              {isCreating ? 'Membuat Backup...' : 'Buat Backup Baru'}
            </Button>
          </RoleGuard>
        </div>
      </div>

      {isLoading ? (
        <div className="space-y-3">
          {Array.from({ length: 3 }).map((_, i) => (
            <div key={i} className="h-12 animate-pulse rounded-lg bg-gray-100" />
          ))}
        </div>
      ) : files.length === 0 ? (
        <p className="py-8 text-center text-sm text-gray-400">Belum ada backup yang dibuat</p>
      ) : (
        <div className="rounded-lg border bg-white overflow-hidden">
          <table className="w-full text-sm">
            <thead className="bg-gray-50 border-b">
              <tr>
                <th className="px-4 py-3 text-left font-medium text-gray-600">Nama File</th>
                <th className="px-4 py-3 text-left font-medium text-gray-600">Ukuran</th>
                <th className="px-4 py-3 text-left font-medium text-gray-600">Dibuat</th>
                <th className="px-4 py-3 text-right font-medium text-gray-600">Aksi</th>
              </tr>
            </thead>
            <tbody>
              {files.map((f) => (
                <tr key={f.filename} className="border-b last:border-0 hover:bg-gray-50">
                  <td className="px-4 py-3 font-mono text-xs text-gray-700">{f.filename}</td>
                  <td className="px-4 py-3 text-gray-500">{f.size}</td>
                  <td className="px-4 py-3 text-gray-500">{formatDate(f.created_at)}</td>
                  <td className="px-4 py-3">
                    <div className="flex items-center justify-end gap-1.5">
                      <RoleGuard menuKey="sistem.backup" action="can_view">
                        <Button
                          variant="outline"
                          size="sm"
                          className="gap-1.5 h-8"
                          onClick={() => handleDownload(f.filename)}
                        >
                          <Download size={13} />
                          Download
                        </Button>
                      </RoleGuard>
                      <RoleGuard menuKey="sistem.backup" action="can_delete">
                        <Button
                          variant="outline"
                          size="sm"
                          className="gap-1.5 h-8 text-red-600 border-red-200 hover:bg-red-50 hover:text-red-700"
                          onClick={() => setRestoreTarget(f)}
                        >
                          <RotateCcw size={13} />
                          Restore
                        </Button>
                      </RoleGuard>
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      <RestoreConfirmDialog
        open={!!restoreTarget}
        onOpenChange={(open) => { if (!open) setRestoreTarget(null) }}
        filename={restoreTarget?.filename ?? null}
        createdAt={restoreTarget && 'created_at' in restoreTarget ? restoreTarget.created_at : null}
        isLoading={isRestoring}
        onConfirm={handleConfirmRestore}
      />
    </div>
  )
}

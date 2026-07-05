import { Plus } from 'lucide-react'

import { RoleGuard } from '@/shared/components'
import { Badge } from '@/shared/components/ui/badge'
import { Button } from '@/shared/components/ui/button'
import { useDisclosure } from '@/shared/hooks'

import { useAppVersionListQuery } from '../../settings.api'
import { AppVersionFormModal } from './AppVersionFormModal'

const PLATFORM_LABEL: Record<string, string> = {
  web:     'Web',
  desktop: 'Desktop',
  android: 'Android',
}

function formatDate(str: string): string {
  return new Date(str).toLocaleDateString('id-ID', {
    day:   '2-digit',
    month: 'short',
    year:  'numeric',
  })
}

export function AppVersionTab() {
  const { data: versions = [], isLoading } = useAppVersionListQuery()
  const { isOpen: formOpen, open: openForm, close: closeForm } = useDisclosure()

  const handleOpenAdd = () => openForm()
  const handleCloseForm = () => closeForm()

  return (
    <div className="space-y-4">
      <div className="flex justify-end">
        <RoleGuard menuKey="sistem.versi" action="can_create">
          <Button onClick={handleOpenAdd}>
            <Plus size={14} className="mr-2" />
            Tambah Versi
          </Button>
        </RoleGuard>
      </div>

      {isLoading ? (
        <div className="space-y-3">
          {Array.from({ length: 3 }).map((_, i) => (
            <div key={i} className="h-12 animate-pulse rounded-lg bg-gray-100" />
          ))}
        </div>
      ) : versions.length === 0 ? (
        <p className="py-8 text-center text-sm text-gray-400">Belum ada data versi aplikasi</p>
      ) : (
        <div className="rounded-lg border bg-white overflow-hidden">
          <table className="w-full text-sm">
            <thead className="bg-gray-50 border-b">
              <tr>
                <th className="px-4 py-3 text-left font-medium text-gray-600">Platform</th>
                <th className="px-4 py-3 text-left font-medium text-gray-600">Versi</th>
                <th className="px-4 py-3 text-left font-medium text-gray-600">Catatan Rilis</th>
                <th className="px-4 py-3 text-left font-medium text-gray-600">Tanggal</th>
                <th className="px-4 py-3 text-center font-medium text-gray-600">Status</th>
                <th className="px-4 py-3 text-left font-medium text-gray-600">Download</th>
              </tr>
            </thead>
            <tbody>
              {versions.map((v) => (
                <tr key={v.id} className="border-b last:border-0 hover:bg-gray-50">
                  <td className="px-4 py-3 font-medium">
                    {PLATFORM_LABEL[v.platform] ?? v.platform}
                  </td>
                  <td className="px-4 py-3 font-mono text-blue-600">{v.version}</td>
                  <td className="px-4 py-3 text-gray-500 text-xs max-w-xs truncate">
                    {v.release_notes || '-'}
                  </td>
                  <td className="px-4 py-3 text-gray-500">{formatDate(v.created_at)}</td>
                  <td className="px-4 py-3 text-center">
                    <div className="flex items-center justify-center gap-1.5">
                      {v.is_latest && <Badge variant="default">Terbaru</Badge>}
                      {v.is_mandatory && <Badge variant="destructive">Wajib</Badge>}
                      {!v.is_latest && !v.is_mandatory && <Badge variant="secondary">Lama</Badge>}
                    </div>
                  </td>
                  <td className="px-4 py-3">
                    <a
                      href={v.download_url}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="text-xs text-blue-600 hover:underline"
                    >
                      Download
                    </a>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      <AppVersionFormModal
        open={formOpen}
        onOpenChange={(open) => { if (!open) handleCloseForm() }}
      />
    </div>
  )
}

import { useState } from 'react'

import { ConfirmDialog } from '@/shared/components'
import { Button } from '@/shared/components/ui/button'
import { formatDateTime } from '@/shared/utils'

import {
  useResolveConflictMutation,
  useSyncConflictsQuery,
} from '../sync.api'
import type { SyncConflict } from '../sync.types'

const ENTITY_TYPE_LABEL: Record<string, string> = {
  product: 'PRODUK',
  transaction: 'TRANSAKSI',
  expense: 'PENGELUARAN',
}

function parseJsonSafe(raw: string): Record<string, unknown> {
  try {
    const parsed = JSON.parse(raw)
    return typeof parsed === 'object' && parsed !== null ? parsed : {}
  } catch {
    return {}
  }
}

function DataDiff({
  serverData,
  localData,
}: {
  serverData: Record<string, unknown>
  localData: Record<string, unknown>
}) {
  const allKeys = Array.from(new Set([...Object.keys(serverData), ...Object.keys(localData)]))

  return (
    <div className="grid grid-cols-2 gap-3 mt-3 text-xs">
      <div className="space-y-1">
        <p className="font-semibold text-gray-500 uppercase tracking-wide mb-1.5">Data Server</p>
        {allKeys.map((key) => {
          const isDiff = String(serverData[key]) !== String(localData[key])
          return (
            <div
              key={key}
              className={`flex justify-between px-2 py-1 rounded ${isDiff ? 'bg-yellow-100' : 'bg-gray-50'}`}
            >
              <span className="text-gray-500">{key}</span>
              <span className="font-medium text-gray-800">{String(serverData[key] ?? '—')}</span>
            </div>
          )
        })}
      </div>
      <div className="space-y-1">
        <p className="font-semibold text-gray-500 uppercase tracking-wide mb-1.5">Data Lokal (Offline)</p>
        {allKeys.map((key) => {
          const isDiff = String(serverData[key]) !== String(localData[key])
          return (
            <div
              key={key}
              className={`flex justify-between px-2 py-1 rounded ${isDiff ? 'bg-yellow-100' : 'bg-gray-50'}`}
            >
              <span className="text-gray-500">{key}</span>
              <span className="font-medium text-gray-800">{String(localData[key] ?? '—')}</span>
            </div>
          )
        })}
      </div>
    </div>
  )
}

function ConflictCard({ conflict }: { conflict: SyncConflict }) {
  const [keepServerOpen, setKeepServerOpen] = useState(false)
  const [useLocalOpen, setUseLocalOpen] = useState(false)

  const { mutate: resolve, isPending: isResolving } = useResolveConflictMutation()

  const serverData = parseJsonSafe(conflict.online_data)
  const localData = parseJsonSafe(conflict.desktop_data)
  const entityLabel = ENTITY_TYPE_LABEL[conflict.entity_type] ?? conflict.entity_type.toUpperCase()

  return (
    <div className="rounded-xl border border-orange-200 bg-white p-4 space-y-2">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-2">
          <span>⚠️</span>
          <span className="font-semibold text-gray-800">
            KONFLIK {entityLabel} — #{conflict.entity_id}
          </span>
        </div>
      </div>

      <p className="text-xs text-gray-500">
        Perangkat: <span className="font-medium">{conflict.device_id}</span> ·{' '}
        {formatDateTime(conflict.created_at)}
      </p>

      <DataDiff serverData={serverData} localData={localData} />

      <div className="flex gap-2 justify-end pt-2 border-t">
        <Button
          size="sm"
          variant="outline"
          className="text-red-600 border-red-200 hover:bg-red-50"
          onClick={() => setUseLocalOpen(true)}
          disabled={isResolving}
        >
          ✗ Pakai Data Lokal
        </Button>
        <Button
          size="sm"
          onClick={() => setKeepServerOpen(true)}
          disabled={isResolving}
        >
          ✓ Terima Server
        </Button>
      </div>

      {/* Semantik BE: action='reject' = buang data offline, pertahankan data server. */}
      <ConfirmDialog
        open={keepServerOpen}
        onOpenChange={setKeepServerOpen}
        title="Terima Data Server"
        description={`Data server akan dipertahankan untuk "${entityLabel} #${conflict.entity_id}". Data lokal (offline) akan dibuang. Lanjutkan?`}
        confirmLabel="Terima Server"
        isLoading={isResolving}
        onConfirm={() => resolve({ id: conflict.id, action: 'reject' }, { onSuccess: () => setKeepServerOpen(false) })}
      />
      {/* Semantik BE: action='approve' = terapkan data offline ke server. */}
      <ConfirmDialog
        open={useLocalOpen}
        onOpenChange={setUseLocalOpen}
        title="Pakai Data Lokal"
        description={`Data lokal (offline) akan diterapkan untuk "${entityLabel} #${conflict.entity_id}", menimpa data server. Lanjutkan?`}
        confirmLabel="Pakai Data Lokal"
        variant="destructive"
        isLoading={isResolving}
        onConfirm={() => resolve({ id: conflict.id, action: 'approve' }, { onSuccess: () => setUseLocalOpen(false) })}
      />
    </div>
  )
}

export function ConflictList() {
  const { data, isLoading } = useSyncConflictsQuery()
  const list = data?.data ?? []

  if (isLoading) {
    return (
      <div className="space-y-3">
        {Array.from({ length: 2 }).map((_, i) => (
          <div key={i} className="h-40 animate-pulse rounded-xl bg-gray-100" />
        ))}
      </div>
    )
  }

  if (list.length === 0) {
    return (
      <div className="rounded-xl border bg-green-50 border-green-200 py-10 text-center">
        <p className="text-green-700 font-medium">Tidak ada konflik</p>
        <p className="text-sm text-green-600 mt-1">Semua data sudah tersinkronisasi dengan baik</p>
      </div>
    )
  }

  return (
    <div className="space-y-4">
      {list.map((c) => (
        <ConflictCard key={c.id} conflict={c} />
      ))}
    </div>
  )
}

import { RefreshCw } from 'lucide-react'

import { Button } from '@/shared/components/ui/button'

import { useTriggerSyncMutation } from '../sync.api'
import { useSyncStatus } from '../hooks/useSyncStatus'

export function SyncStatusCard() {
  const { conflictCount } = useSyncStatus()
  const { mutate: triggerSync, isPending } = useTriggerSyncMutation()

  const current = conflictCount > 0
    ? { icon: '⚠️', text: 'Ada Konflik', color: 'text-orange-600', bg: 'bg-orange-50 border-orange-200' }
    : { icon: '🔵', text: 'Siap', color: 'text-gray-600', bg: 'bg-gray-50' }

  return (
    <div className={`rounded-xl border p-4 ${current.bg}`}>
      <div className="flex items-center justify-between mb-3">
        <h3 className="font-semibold text-gray-700">Status Sinkronisasi</h3>
        <Button
          size="sm"
          variant="outline"
          onClick={() => triggerSync()}
          disabled={isPending}
          className="gap-1.5"
        >
          <RefreshCw size={14} className={isPending ? 'animate-spin' : ''} />
          Sync Manual
        </Button>
      </div>

      <div className="flex items-center gap-2 mb-2">
        <span className="text-xl">{current.icon}</span>
        <span className={`font-semibold ${current.color}`}>{current.text}</span>
      </div>

      <div className="flex gap-6 text-sm text-gray-600">
        <span>
          Konflik:{' '}
          <span className={`font-medium ${conflictCount > 0 ? 'text-orange-600' : ''}`}>
            {conflictCount} item
          </span>
        </span>
      </div>
    </div>
  )
}

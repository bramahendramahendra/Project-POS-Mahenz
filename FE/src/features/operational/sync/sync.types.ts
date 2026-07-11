export type EntityType = 'product' | 'transaction' | 'expense'

export interface SyncConflict {
  id: number
  entity_type: EntityType | string
  entity_id: number
  local_id: string
  device_id: string
  desktop_data: string
  online_data: string
  desktop_time: string
  online_time: string
  status: 'pending' | 'resolved'
  created_at: string
}

export interface ConflictListResponse {
  data: SyncConflict[]
  total: number
  page: number
  limit: number
}

export interface SyncQueueItem {
  id: number
  device_id: string
  entity_type: EntityType | string
  entity_id: number
  action: 'create' | 'update' | 'delete'
  status: 'pending' | 'syncing' | 'synced' | 'failed'
  retry_count: number
  created_at: string
}

export interface QueueListResponse {
  data: SyncQueueItem[]
  total: number
}

export interface SyncHistoryItem {
  id: number
  device_id: string
  device_type: 'desktop' | 'web' | 'android'
  total_items: number
  synced_items: number
  conflict_items: number
  failed_items: number
  pending_items: number
  duration_ms: number
  status: 'success' | 'partial' | 'failed'
  started_at: string
  finished_at: string | null
}

export interface SyncHistoryListResponse {
  data: SyncHistoryItem[]
  total: number
  page: number
  limit: number
}

export interface SyncFilter {
  page?: number
  limit?: number
}

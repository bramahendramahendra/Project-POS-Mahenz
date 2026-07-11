import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'

import { api, apiClient } from '@/services'
import { queryKeys } from '@/shared/constants'

import type { BackupInfo } from './backup.types'

export function useBackupListQuery() {
  return useQuery({
    queryKey: queryKeys.backup.list(),
    queryFn: () => api.get<{ files: BackupInfo[] }>('/backup/list'),
  })
}

export function useCreateBackupMutation() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: () => api.post<BackupInfo>('/backup'),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: queryKeys.backup.list() })
      toast.success('Backup berhasil dibuat')
    },
    onError: (e: Error) => toast.error(e.message),
  })
}

export async function downloadBackupFile(filename: string): Promise<void> {
  const response = await apiClient.get(`/backup/download/${encodeURIComponent(filename)}`, {
    responseType: 'blob',
  })
  const url = URL.createObjectURL(response.data as Blob)
  const a = document.createElement('a')
  a.href = url
  a.download = filename
  a.click()
  URL.revokeObjectURL(url)
}

async function fetchBackupAsFile(filename: string): Promise<File> {
  const response = await apiClient.get(`/backup/download/${encodeURIComponent(filename)}`, {
    responseType: 'blob',
  })
  return new File([response.data as Blob], filename, { type: 'application/sql' })
}

export function useRestoreBackupMutation() {
  return useMutation({
    mutationFn: async (source: { filename: string } | { file: File }) => {
      const file = 'file' in source ? source.file : await fetchBackupAsFile(source.filename)
      const formData = new FormData()
      formData.append('file', file)
      return api.post('/restore', formData)
    },
    onError: (e: Error) => toast.error(e.message),
  })
}

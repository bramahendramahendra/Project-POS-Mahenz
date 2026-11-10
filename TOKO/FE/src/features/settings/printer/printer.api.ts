import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'

import { api } from '@/services'
import { queryKeys } from '@/shared/constants'

import type { PrinterSettings } from './printer.types'

export function usePrinterSettingsQuery() {
  return useQuery({
    queryKey: queryKeys.settings.printer(),
    queryFn: () => api.get<PrinterSettings>('/settings/printer'),
  })
}

export function useUpdatePrinterSettingsMutation() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (payload: PrinterSettings) => api.post<PrinterSettings>('/settings/printer', payload),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: queryKeys.settings.printer() })
      toast.success('Pengaturan printer berhasil disimpan')
    },
    onError: (e: Error) => toast.error(e.message),
  })
}

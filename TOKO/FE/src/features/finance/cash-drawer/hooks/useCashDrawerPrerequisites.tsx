import { Clock } from 'lucide-react'

import { useShiftOptionsQuery } from '@/features/operational/shifts'

export function useCashDrawerPrerequisites() {
  const { data: shiftsData, isLoading } = useShiftOptionsQuery()

  const hasShifts = (shiftsData ?? []).length > 0

  return {
    isLoading,
    items: [
      {
        label: 'Belum ada shift aktif — tambahkan di menu Operasional › Manajemen Shift',
        metLabel: 'Shift aktif tersedia',
        met: hasShifts,
        icon: <Clock size={24} />,
      },
    ],
  }
}

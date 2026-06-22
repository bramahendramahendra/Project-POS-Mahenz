import type { ReactNode } from 'react'
import { Clock } from 'lucide-react'

import { useShiftOptionsQuery } from '@/features/operational/shifts'

interface ShiftPrerequisiteGuardProps {
  children: ReactNode
}

export function ShiftPrerequisiteGuard({ children }: ShiftPrerequisiteGuardProps) {
  const { data: shiftsRaw, isLoading } = useShiftOptionsQuery()
  const shifts = shiftsRaw ?? []

  if (isLoading) {
    return (
      <div className="space-y-3">
        {[1, 2, 3, 4, 5].map((i) => (
          <div key={i} className="h-12 animate-pulse rounded-md bg-gray-100" />
        ))}
      </div>
    )
  }

  const hasShifts = shifts.length > 0

  if (!hasShifts) {
    return (
      <div className="flex flex-col items-center justify-center rounded-lg border border-dashed bg-white px-6 py-16 text-center">
        <div className="mb-4">
          <div className="rounded-full bg-amber-50 p-4">
            <Clock size={28} className="text-amber-500" />
          </div>
        </div>
        <h3 className="mb-1 text-base font-semibold text-gray-800">Belum bisa membuka Kas Harian</h3>
        <p className="mb-1 text-sm text-gray-500">
          Sebelum menggunakan Kas Harian, pastikan data berikut sudah tersedia:
        </p>
        <ul className="mb-6 text-sm">
          <li className="flex items-center gap-2 text-amber-600">
            <span>!</span>
            Belum ada shift aktif — tambahkan di menu Operasional &rsaquo; Manajemen Shift
          </li>
        </ul>
      </div>
    )
  }

  return <>{children}</>
}

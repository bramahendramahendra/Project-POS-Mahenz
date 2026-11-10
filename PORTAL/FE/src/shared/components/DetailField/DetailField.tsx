import type { ReactNode } from 'react'

interface DetailFieldProps {
  label: string
  value?: string
  children?: ReactNode
}

export function DetailField({ label, value, children }: DetailFieldProps) {
  return (
    <div className="space-y-0.5">
      <p className="text-xs text-gray-500">{label}</p>
      {children ?? <p className="font-medium text-gray-800">{value ?? '—'}</p>}
    </div>
  )
}

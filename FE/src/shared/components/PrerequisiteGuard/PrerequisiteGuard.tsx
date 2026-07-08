import type { ReactNode } from 'react'

interface PrerequisiteItem {
  label: string
  metLabel?: string
  met: boolean
  icon?: ReactNode
}

interface PrerequisiteGuardProps {
  isLoading: boolean
  title: string
  description: string
  items: PrerequisiteItem[]
  children: ReactNode
}

export function PrerequisiteGuard({ isLoading, title, description, items, children }: PrerequisiteGuardProps) {
  if (isLoading) {
    return (
      <div className="space-y-3">
        {[1, 2, 3, 4, 5].map((i) => (
          <div key={i} className="h-12 animate-pulse rounded-md bg-gray-100" />
        ))}
      </div>
    )
  }

  const allMet = items.every((item) => item.met)

  if (!allMet) {
    return (
      <div className="flex flex-col items-center justify-center rounded-lg border border-dashed bg-white px-6 py-16 text-center">
        <div className="mb-4 flex gap-3">
          {items.map((item, i) => (
            <div key={i} className={`rounded-full p-3 ${item.met ? 'bg-green-50' : 'bg-amber-50'}`}>
              <span className={item.met ? 'text-green-500' : 'text-amber-500'}>{item.icon}</span>
            </div>
          ))}
        </div>
        <h3 className="mb-1 text-base font-semibold text-gray-800">{title}</h3>
        <p className="mb-1 text-sm text-gray-500">{description}</p>
        <ul className="mb-6 text-sm">
          {items.map((item, i) => (
            <li key={i} className={`flex items-center gap-2 ${item.met ? 'text-green-600' : 'text-amber-600'}`}>
              <span>{item.met ? '✓' : '!'}</span>
              {item.met ? (item.metLabel ?? item.label) : item.label}
            </li>
          ))}
        </ul>
      </div>
    )
  }

  return <>{children}</>
}

interface SummaryCardProps {
  label: string
  value: string
  sub?: string
  color: 'gray' | 'blue' | 'red' | 'green' | 'orange'
}

const colorMap: Record<SummaryCardProps['color'], string> = {
  gray:   'bg-gray-50 border-gray-200 text-gray-700',
  blue:   'bg-blue-50 border-blue-100 text-blue-700',
  red:    'bg-red-50 border-red-100 text-red-700',
  green:  'bg-green-50 border-green-100 text-green-700',
  orange: 'bg-orange-50 border-orange-100 text-orange-700',
}

const subColorMap: Record<SummaryCardProps['color'], string> = {
  gray:   'text-gray-400',
  blue:   'text-blue-400',
  red:    'text-red-400',
  green:  'text-green-500',
  orange: 'text-orange-400',
}

export function SummaryCard({ label, value, sub, color }: SummaryCardProps) {
  return (
    <div className={`rounded-lg border p-2.5 ${colorMap[color]}`}>
      <p className="text-xs opacity-70 mb-0.5">{label}</p>
      <p className="font-semibold text-sm leading-tight">{value}</p>
      {sub && <p className={`text-xs mt-0.5 ${subColorMap[color]}`}>{sub}</p>}
    </div>
  )
}

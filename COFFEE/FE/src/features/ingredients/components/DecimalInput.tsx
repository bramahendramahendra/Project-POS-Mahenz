import { useState } from 'react'

import { Input } from '@/shared/components/ui/input'

import { parseDecimalInput } from '../ingredients.utils'

interface DecimalInputProps {
  id?: string
  value: number | null | undefined
  onChange: (value: number | null) => void
  placeholder?: string
  className?: string
  'aria-invalid'?: boolean
}

// Input angka desimal (koma atau titik) untuk jumlah stok bahan. Kosong = null.
export function DecimalInput({ id, value, onChange, placeholder, className, ...rest }: DecimalInputProps) {
  const [text, setText] = useState(value == null || Number.isNaN(value) ? '' : String(value).replace('.', ','))

  // Sinkron bila nilai diubah dari luar (mis. form di-reset): disesuaikan saat render, bukan di effect.
  const [prevValue, setPrevValue] = useState(value)
  if (value !== prevValue) {
    setPrevValue(value)
    const current = parseDecimalInput(text)
    if (value == null && text !== '') setText('')
    else if (value != null && !Number.isNaN(value) && current !== value) setText(String(value).replace('.', ','))
  }

  return (
    <Input
      id={id}
      type="text"
      inputMode="decimal"
      autoComplete="off"
      value={text}
      placeholder={placeholder}
      className={className}
      aria-invalid={rest['aria-invalid']}
      onChange={(e) => {
        const next = e.target.value
        if (!/^[0-9.,\s]*$/.test(next)) return
        setText(next)
        const parsed = parseDecimalInput(next)
        onChange(Number.isNaN(parsed) ? null : parsed)
      }}
    />
  )
}

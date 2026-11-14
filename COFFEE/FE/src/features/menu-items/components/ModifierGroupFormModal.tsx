import { useState } from 'react'
import { Plus, Trash2 } from 'lucide-react'

import { FormModal } from '@/shared/components'
import { Button } from '@/shared/components/ui/button'
import { Checkbox } from '@/shared/components/ui/checkbox'
import { Input } from '@/shared/components/ui/input'
import { Label } from '@/shared/components/ui/label'
import { RupiahInput } from '@/shared/components/ui/rupiah-input'
import { Switch } from '@/shared/components/ui/switch'
import { ToggleGroup, ToggleGroupItem } from '@/shared/components/ui/toggle-group'

import { useSaveModifierGroupMutation } from '../menu-items.api'
import type { ModifierGroup, ModifierOptionInput, SelectionType } from '../menu-items.types'

interface ModifierGroupFormModalProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  group: ModifierGroup | null
}

const emptyOption = (): ModifierOptionInput => ({ id: 0, name: '', price: 0, is_default: false, is_active: true })

// Form kelompok tambahan (keputusan 54, 87): pilih satu / boleh banyak, wajib atau opsional,
// pilihan dengan harga (0 = gratis), paling banyak satu pilihan bawaan.
// Komponen dibuat ulang (key) setiap kali dibuka, jadi nilai awal cukup dari props.
export function ModifierGroupFormModal({ open, onOpenChange, group }: ModifierGroupFormModalProps) {
  const [name, setName] = useState(group?.name ?? '')
  const [selectionType, setSelectionType] = useState<SelectionType>(group?.selection_type ?? 'single')
  const [isRequired, setIsRequired] = useState(group?.is_required ?? false)
  const [options, setOptions] = useState<ModifierOptionInput[]>(
    group ? group.options.map((o) => ({ ...o })) : [emptyOption(), emptyOption()]
  )
  const [error, setError] = useState('')
  const { mutate: save, isPending } = useSaveModifierGroupMutation()

  const updateOption = (index: number, patch: Partial<ModifierOptionInput>) =>
    setOptions((prev) =>
      prev.map((o, i) => {
        if (i === index) return { ...o, ...patch }
        // Pilihan bawaan hanya satu: memilih satu menghapus tanda di pilihan lain.
        return patch.is_default ? { ...o, is_default: false } : o
      })
    )

  const submit = () => {
    const cleaned = options.map((o) => ({ ...o, name: o.name.trim() })).filter((o) => o.name !== '' || o.id > 0)
    if (!name.trim()) return setError('Nama kelompok wajib diisi')
    if (cleaned.length === 0) return setError('Minimal satu pilihan')
    if (cleaned.some((o) => o.name === '')) return setError('Nama pilihan tidak boleh kosong')
    if (!cleaned.some((o) => o.is_active)) return setError('Minimal satu pilihan aktif')
    if (cleaned.some((o) => o.is_default && !o.is_active)) return setError('Pilihan bawaan harus aktif')
    const names = cleaned.map((o) => o.name.toLowerCase())
    if (new Set(names).size !== names.length) return setError('Nama pilihan tidak boleh ganda')
    setError('')
    save(
      { id: group?.id, name: name.trim(), selection_type: selectionType, is_required: isRequired, options: cleaned },
      { onSuccess: () => onOpenChange(false) }
    )
  }

  return (
    <FormModal
      open={open}
      onOpenChange={onOpenChange}
      title={group ? 'Edit Kelompok Tambahan' : 'Tambah Kelompok Tambahan'}
      size="lg"
      isLoading={isPending}
      onSubmit={submit}
      submitLabel="Simpan"
    >
      <div className="space-y-4">
        <div className="space-y-1.5">
          <Label htmlFor="group-name">
            Nama Kelompok <span className="text-red-500">*</span>
          </Label>
          <Input id="group-name" value={name} onChange={(e) => setName(e.target.value)} placeholder="Contoh: Level Gula, Topping" />
        </div>

        <div className="flex flex-wrap items-center gap-6">
          <div className="space-y-1.5">
            <Label>Cara Memilih</Label>
            <ToggleGroup type="single" value={selectionType} onValueChange={(v) => v && setSelectionType(v as SelectionType)}>
              <ToggleGroupItem value="single">Pilih satu</ToggleGroupItem>
              <ToggleGroupItem value="multiple">Boleh banyak</ToggleGroupItem>
            </ToggleGroup>
          </div>
          <label className="flex items-center gap-2 pt-5 text-sm text-gray-700">
            <Switch checked={isRequired} onCheckedChange={setIsRequired} aria-label="Wajib dipilih" />
            Wajib dipilih kasir
          </label>
        </div>

        <div className="space-y-2">
          <Label>Pilihan</Label>
          <div className="hidden grid-cols-[1fr_150px_80px_70px_36px] gap-2 px-1 text-xs font-medium text-gray-500 sm:grid">
            <span>Nama</span>
            <span>Harga (0 = gratis)</span>
            <span className="text-center">Bawaan</span>
            <span className="text-center">Aktif</span>
            <span />
          </div>
          {options.map((o, i) => (
            <div key={o.id || `new-${i}`} className="grid grid-cols-1 gap-2 rounded-md border p-2 sm:grid-cols-[1fr_150px_80px_70px_36px] sm:items-center sm:border-0 sm:p-0 sm:px-1">
              <Input
                aria-label={`Nama pilihan ${i + 1}`}
                value={o.name}
                onChange={(e) => updateOption(i, { name: e.target.value })}
                placeholder="Contoh: Normal, Extra shot"
              />
              <RupiahInput value={o.price} onChange={(v) => updateOption(i, { price: v })} placeholder="0" />
              <label className="flex items-center gap-2 text-sm sm:justify-center">
                <Checkbox
                  checked={o.is_default}
                  onCheckedChange={(v) => updateOption(i, { is_default: v === true })}
                  aria-label={`Pilihan bawaan ${i + 1}`}
                />
                <span className="sm:hidden">Bawaan</span>
              </label>
              <label className="flex items-center gap-2 text-sm sm:justify-center">
                <Checkbox
                  checked={o.is_active}
                  onCheckedChange={(v) => updateOption(i, { is_active: v === true })}
                  aria-label={`Pilihan aktif ${i + 1}`}
                />
                <span className="sm:hidden">Aktif</span>
              </label>
              <Button
                type="button"
                variant="ghost"
                size="icon"
                aria-label={`Hapus pilihan ${i + 1}`}
                className="h-8 w-8 text-gray-500 hover:text-red-600"
                onClick={() => setOptions((prev) => prev.filter((_, j) => j !== i))}
              >
                <Trash2 size={14} />
              </Button>
            </div>
          ))}
          <Button type="button" variant="outline" size="sm" className="gap-1" onClick={() => setOptions((prev) => [...prev, emptyOption()])}>
            <Plus size={14} />
            Tambah Pilihan
          </Button>
        </div>

        {error && <p className="text-sm text-red-500">{error}</p>}
      </div>
    </FormModal>
  )
}

import { Plus, Trash2 } from 'lucide-react'

import { Button } from '@/shared/components/ui/button'
import { Input } from '@/shared/components/ui/input'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/shared/components/ui/select'
import type { IngredientOption } from '@/features/ingredients/ingredients.types'

import type { RecipeLine } from '../menu-items.types'

interface RecipeEditorProps {
  lines: RecipeLine[]
  onChange: (lines: RecipeLine[]) => void
  ingredients: IngredientOption[]
  label: string // untuk aria-label, mis. "Regular · Panas"
}

// Resep = baris bahan + takaran teks (keputusan 86). Catatan saja, tidak mengurangi stok.
export function RecipeEditor({ lines, onChange, ingredients, label }: RecipeEditorProps) {
  const update = (index: number, patch: Partial<RecipeLine>) =>
    onChange(lines.map((l, i) => (i === index ? { ...l, ...patch } : l)))

  return (
    <div className="space-y-2">
      {lines.length === 0 && <p className="text-xs text-gray-500">Belum ada resep (opsional).</p>}
      {lines.map((line, i) => (
        <div key={i} className="grid grid-cols-[1fr_140px_36px] items-center gap-2">
          <Select
            value={line.ingredient_id ? String(line.ingredient_id) : ''}
            onValueChange={(v) => update(i, { ingredient_id: Number(v) })}
          >
            <SelectTrigger className="w-full" aria-label={`Bahan resep ${label} baris ${i + 1}`}>
              <SelectValue placeholder="Pilih bahan..." />
            </SelectTrigger>
            <SelectContent>
              {ingredients.map((ing) => (
                <SelectItem
                  key={ing.id}
                  value={String(ing.id)}
                  disabled={ing.id !== line.ingredient_id && lines.some((l) => l.ingredient_id === ing.id)}
                >
                  {ing.name}
                </SelectItem>
              ))}
              {line.ingredient_id > 0 && !ingredients.some((ing) => ing.id === line.ingredient_id) && (
                <SelectItem value={String(line.ingredient_id)}>{line.ingredient_name ?? 'Bahan nonaktif'}</SelectItem>
              )}
            </SelectContent>
          </Select>
          <Input
            aria-label={`Takaran ${label} baris ${i + 1}`}
            value={line.amount}
            maxLength={50}
            onChange={(e) => update(i, { amount: e.target.value })}
            placeholder="Mis. 18 gram"
          />
          <Button
            type="button"
            variant="ghost"
            size="icon"
            aria-label={`Hapus baris resep ${label} ${i + 1}`}
            className="h-8 w-8 text-gray-500 hover:text-red-600"
            onClick={() => onChange(lines.filter((_, j) => j !== i))}
          >
            <Trash2 size={14} />
          </Button>
        </div>
      ))}
      <Button
        type="button"
        variant="outline"
        size="sm"
        className="gap-1"
        disabled={ingredients.length === 0}
        onClick={() => onChange([...lines, { ingredient_id: 0, amount: '' }])}
      >
        <Plus size={14} />
        Tambah Bahan
      </Button>
      {ingredients.length === 0 && <p className="text-xs text-amber-600">Belum ada bahan aktif. Tambahkan di menu Bahan Baku.</p>}
    </div>
  )
}

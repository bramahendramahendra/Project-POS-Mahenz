import { useEffect, useState } from 'react'
import { Check, ChevronsUpDown, Loader2 } from 'lucide-react'

import { cn } from '@/shared/utils'
import { useDebounce } from '@/shared/hooks'
import { Button } from '@/shared/components/ui/button'
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from '@/shared/components/ui/command'
import { Popover, PopoverContent, PopoverTrigger } from '@/shared/components/ui/popover'

interface AsyncComboboxProps<T> {
  value?: string | number
  onValueChange: (value: string | number | undefined, item?: T) => void
  onSearch: (keyword: string) => void
  options: T[] | null | undefined
  getOptionValue: (item: T) => string | number
  getOptionLabel: (item: T) => string
  selectedLabel?: string
  isLoading?: boolean
  disabled?: boolean
  placeholder?: string
  searchPlaceholder?: string
  emptyText?: string
  minChars?: number
  debounceMs?: number
  className?: string
}

export function AsyncCombobox<T>({
  value,
  onValueChange,
  onSearch,
  options,
  getOptionValue,
  getOptionLabel,
  selectedLabel,
  isLoading,
  disabled,
  placeholder = 'Pilih...',
  searchPlaceholder = 'Ketik untuk mencari...',
  emptyText = 'Tidak ada hasil.',
  minChars = 2,
  debounceMs = 300,
  className,
}: AsyncComboboxProps<T>) {
  const [open, setOpen] = useState(false)
  const [keyword, setKeyword] = useState('')
  const debouncedKeyword = useDebounce(keyword, debounceMs)

  useEffect(() => {
    if (debouncedKeyword.trim().length >= minChars) {
      onSearch(debouncedKeyword.trim())
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [debouncedKeyword])

  const safeOptions = options ?? []
  const selected = safeOptions.find((o) => String(getOptionValue(o)) === String(value))
  const label = selected ? getOptionLabel(selected) : selectedLabel

  return (
    <Popover
      open={open}
      onOpenChange={(next) => {
        setOpen(next)
        if (!next) setKeyword('')
      }}
    >
      <PopoverTrigger asChild>
        <Button
          type="button"
          variant="outline"
          role="combobox"
          aria-expanded={open}
          disabled={disabled}
          className={cn('h-9 w-full justify-between text-sm font-normal', !label && 'text-muted-foreground', className)}
        >
          <span className="truncate">{label || placeholder}</span>
          <ChevronsUpDown className="ml-2 h-4 w-4 shrink-0 opacity-50" />
        </Button>
      </PopoverTrigger>
      <PopoverContent
        className="z-[1200] w-[var(--radix-popover-trigger-width)] p-0 pointer-events-auto"
        align="start"
      >
        <Command shouldFilter={false}>
          <CommandInput
            value={keyword}
            onValueChange={setKeyword}
            placeholder={searchPlaceholder}
            onKeyDown={(e) => e.stopPropagation()}
          />
          <CommandList>
            {keyword.trim().length < minChars ? (
              <div className="py-6 text-center text-sm text-muted-foreground">
                Ketik minimal {minChars} huruf untuk mencari
              </div>
            ) : isLoading ? (
              <div className="flex items-center justify-center gap-2 py-6 text-sm text-muted-foreground">
                <Loader2 className="h-4 w-4 animate-spin" /> Mencari...
              </div>
            ) : (
              <>
                <CommandEmpty>{emptyText}</CommandEmpty>
                <CommandGroup>
                  {safeOptions.map((option) => {
                    const optValue = getOptionValue(option)
                    return (
                      <CommandItem
                        key={optValue}
                        value={String(optValue)}
                        onSelect={() => {
                          onValueChange(optValue, option)
                          setOpen(false)
                        }}
                      >
                        <Check
                          className={cn(
                            'mr-2 h-4 w-4',
                            String(value) === String(optValue) ? 'opacity-100' : 'opacity-0'
                          )}
                        />
                        {getOptionLabel(option)}
                      </CommandItem>
                    )
                  })}
                </CommandGroup>
              </>
            )}
          </CommandList>
        </Command>
      </PopoverContent>
    </Popover>
  )
}

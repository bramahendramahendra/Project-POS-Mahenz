import { usePageSizeOptionsQuery } from '@/features/settings/settings.api'

const FALLBACK: number[] = [10, 20, 50]

export function usePageSizeOptions(): number[] {
  const { data } = usePageSizeOptionsQuery()
  return data ?? FALLBACK
}

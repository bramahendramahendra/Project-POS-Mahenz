import { Suspense } from 'react'

import { ErrorBoundary, PageLoader } from '@/shared/components'

export function LazyRoute({ children }: { children: React.ReactNode }) {
  return (
    <ErrorBoundary>
      <Suspense fallback={<PageLoader />}>{children}</Suspense>
    </ErrorBoundary>
  )
}

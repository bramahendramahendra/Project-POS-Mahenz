import * as CollapsiblePrimitive from '@radix-ui/react-collapsible'

import { cn } from '@/shared/utils'

const Collapsible = CollapsiblePrimitive.Root

const CollapsibleTrigger = CollapsiblePrimitive.Trigger

function CollapsibleContent({
  className,
  ...props
}: React.ComponentPropsWithoutRef<typeof CollapsiblePrimitive.Content>) {
  return (
    <CollapsiblePrimitive.Content
      className={cn(
        'overflow-hidden',
        'data-[state=open]:animate-collapsible-down',
        'data-[state=closed]:animate-collapsible-up',
        className,
      )}
      {...props}
    />
  )
}

export { Collapsible, CollapsibleTrigger, CollapsibleContent }

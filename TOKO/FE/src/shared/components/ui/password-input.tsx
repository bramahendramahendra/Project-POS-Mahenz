import * as React from 'react'
import { Eye, EyeOff } from 'lucide-react'

import { cn } from '@/shared/utils'
import { Input } from './input'

const PasswordInput = React.forwardRef<HTMLInputElement, React.ComponentProps<'input'>>(
  ({ className, ...props }, ref) => {
    const [showPassword, setShowPassword] = React.useState(false)

    return (
      <div className="relative">
        <Input
          {...props}
          ref={ref}
          type={showPassword ? 'text' : 'password'}
          className={cn('pr-10', className)}
        />
        <button
          type="button"
          tabIndex={-1}
          onClick={() => setShowPassword((v) => !v)}
          className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600"
        >
          {showPassword ? <EyeOff size={16} /> : <Eye size={16} />}
        </button>
      </div>
    )
  }
)
PasswordInput.displayName = 'PasswordInput'

export { PasswordInput }

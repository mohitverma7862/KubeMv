import { Slot } from '@radix-ui/react-slot'
import { cva, type VariantProps } from 'class-variance-authority'
import type { ButtonHTMLAttributes } from 'react'
import { cn } from '../lib/cn'

const buttonVariants = cva(
  'inline-flex items-center justify-center rounded-md text-sm font-medium transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[var(--color-kp-accent)] disabled:pointer-events-none disabled:opacity-50',
  {
    variants: {
      variant: {
        primary: 'bg-[var(--color-kp-accent)] text-white hover:brightness-110',
        secondary:
          'bg-[var(--color-kp-surface-2)] text-[var(--color-kp-text)] border border-[var(--color-kp-border)] hover:bg-[var(--color-kp-surface)]',
        ghost: 'text-[var(--color-kp-muted)] hover:bg-[var(--color-kp-surface-2)] hover:text-[var(--color-kp-text)]',
        danger: 'bg-[var(--color-kp-critical)] text-white hover:brightness-110',
      },
      size: {
        sm: 'h-8 px-3',
        md: 'h-9 px-4',
        lg: 'h-10 px-5',
      },
    },
    defaultVariants: {
      variant: 'primary',
      size: 'md',
    },
  },
)

export interface KPButtonProps
  extends ButtonHTMLAttributes<HTMLButtonElement>,
    VariantProps<typeof buttonVariants> {
  asChild?: boolean
}

export function KPButton({
  className,
  variant,
  size,
  asChild = false,
  ...props
}: KPButtonProps) {
  const Comp = asChild ? Slot : 'button'
  return <Comp className={cn(buttonVariants({ variant, size }), className)} {...props} />
}

import { cn } from '../lib/cn'

type Status = 'healthy' | 'warning' | 'critical' | 'unknown'

const styles: Record<Status, string> = {
  healthy: 'text-[var(--color-kp-healthy)] border-[var(--color-kp-healthy)]/40 bg-[var(--color-kp-healthy)]/10',
  warning: 'text-[var(--color-kp-warning)] border-[var(--color-kp-warning)]/40 bg-[var(--color-kp-warning)]/10',
  critical: 'text-[var(--color-kp-critical)] border-[var(--color-kp-critical)]/40 bg-[var(--color-kp-critical)]/10',
  unknown: 'text-[var(--color-kp-muted)] border-[var(--color-kp-border)] bg-[var(--color-kp-surface-2)]',
}

export function KPStatusBadge({
  status,
  label,
}: {
  status: Status
  label: string
}) {
  return (
    <span
      className={cn(
        'inline-flex items-center gap-2 rounded-full border px-2.5 py-0.5 text-xs font-medium uppercase tracking-wide',
        styles[status],
      )}
    >
      <span aria-hidden className="h-1.5 w-1.5 rounded-full bg-current" />
      {label}
    </span>
  )
}

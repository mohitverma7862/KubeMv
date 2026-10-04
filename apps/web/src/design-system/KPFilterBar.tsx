import { KPButton } from './KPButton'

export function KPFilterBar({
  value,
  onChange,
  onRefresh,
  placeholder = 'Filter with / or :pods /api',
}: {
  value: string
  onChange: (value: string) => void
  onRefresh: () => void
  placeholder?: string
}) {
  return (
    <div className="flex items-center gap-2">
      <input
        value={value}
        onChange={(e) => onChange(e.target.value)}
        placeholder={placeholder}
        className="flex-1 rounded-md border border-[var(--color-kp-border)] bg-[var(--color-kp-surface-2)] px-3 py-2 text-sm"
      />
      <KPButton type="button" variant="secondary" size="sm" onClick={onRefresh}>Refresh</KPButton>
    </div>
  )
}

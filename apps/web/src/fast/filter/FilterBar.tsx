import { CommandAutocomplete } from '../command/CommandAutocomplete'

export function FilterBar({
  value,
  onChange,
}: {
  value: string
  onChange?: (value: string) => void
}) {
  return (
    <label className="km-fast-muted" style={{ display: 'block', padding: '4px 8px' }}>
      Filter
      <input
        aria-label="Filter"
        value={value}
        onChange={(event) => onChange?.(event.target.value)}
        style={{
          marginLeft: 8,
          width: 'min(420px, 60%)',
          background: 'transparent',
          border: '1px solid var(--km-border)',
          borderRadius: 'var(--km-radius)',
        }}
      />
      <CommandAutocomplete value={value} />
    </label>
  )
}

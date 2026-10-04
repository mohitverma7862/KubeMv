export function KPNamespaceSelector({
  value,
  onChange,
  namespaces,
}: {
  value: string
  onChange: (ns: string) => void
  namespaces: string[]
}) {
  const options = ['all', ...namespaces]
  return (
    <label className="flex items-center gap-2 text-sm text-[var(--color-kp-muted)]">
      Namespace
      <select
        className="rounded-md border border-[var(--color-kp-border)] bg-[var(--color-kp-surface-2)] px-2 py-1 text-[var(--color-kp-text)]"
        value={value}
        onChange={(e) => onChange(e.target.value)}
      >
        {options.map((ns) => (
          <option key={ns} value={ns}>{ns}</option>
        ))}
      </select>
    </label>
  )
}

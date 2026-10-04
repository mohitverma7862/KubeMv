import type { ClusterRef } from '../services/types'

export function KPClusterSelector({
  clusters,
  value,
  onChange,
}: {
  clusters: ClusterRef[]
  value: string
  onChange: (id: string) => void
}) {
  return (
    <label className="flex items-center gap-2 text-sm text-[var(--color-kp-muted)]">
      Cluster
      <select
        className="rounded-md border border-[var(--color-kp-border)] bg-[var(--color-kp-surface-2)] px-2 py-1 text-[var(--color-kp-text)]"
        value={value}
        onChange={(e) => onChange(e.target.value)}
      >
        {clusters.map((c) => (
          <option key={c.id} value={c.id}>{c.name}</option>
        ))}
      </select>
    </label>
  )
}

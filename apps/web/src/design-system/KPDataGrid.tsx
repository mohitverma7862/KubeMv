import type { ResourceRow } from '../services/types'
import { cn } from '../lib/cn'

export function KPDataGrid({
  rows,
  selectedName,
  onSelect,
}: {
  rows: ResourceRow[]
  selectedName?: string
  onSelect: (row: ResourceRow) => void
}) {
  if (rows.length === 0) {
    return <div className="p-6 text-sm text-[var(--color-kp-muted)]">No resources match the current filter.</div>
  }
  return (
    <div className="overflow-auto">
      <table className="min-w-full text-left text-sm">
        <thead className="sticky top-0 bg-[var(--color-kp-surface)] text-xs uppercase tracking-wide text-[var(--color-kp-muted)]">
          <tr>
            <th className="px-3 py-2">Namespace</th>
            <th className="px-3 py-2">Name</th>
            <th className="px-3 py-2">Status</th>
            <th className="px-3 py-2">Age</th>
            <th className="px-3 py-2">Extra</th>
          </tr>
        </thead>
        <tbody>
          {rows.map((row) => (
            <tr
              key={`${row.namespace}/${row.name}`}
              className={cn(
                'cursor-pointer border-t border-[var(--color-kp-border)] hover:bg-[var(--color-kp-surface-2)]',
                selectedName === row.name && 'bg-[var(--color-kp-surface-2)]',
              )}
              onClick={() => onSelect(row)}
            >
              <td className="px-3 py-2 text-[var(--color-kp-muted)]">{row.namespace || '—'}</td>
              <td className="px-3 py-2 font-medium">{row.name}</td>
              <td className="px-3 py-2">{row.status}</td>
              <td className="px-3 py-2 tabular-nums">{row.age}</td>
              <td className="px-3 py-2 text-xs text-[var(--color-kp-muted)]">
                {row.extra ? Object.entries(row.extra).map(([k, v]) => `${k}=${v}`).join(' ') : '—'}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}

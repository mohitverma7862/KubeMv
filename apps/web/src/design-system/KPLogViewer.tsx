import { useMemo, useState } from 'react'
import { KPButton } from './KPButton'

export function KPLogViewer({
  logs,
  loading,
  onRefresh,
  onSearch,
}: {
  logs: string
  loading?: boolean
  onRefresh: () => void
  onSearch: (q: string) => void
}) {
  const [search, setSearch] = useState('')
  const lines = useMemo(() => logs.split('\n').filter(Boolean), [logs])
  const filtered = useMemo(() => {
    if (!search) return lines
    const q = search.toLowerCase()
    return lines.filter((l) => l.toLowerCase().includes(q))
  }, [lines, search])

  return (
    <div className="flex h-full min-h-[240px] flex-col gap-2">
      <div className="flex flex-wrap items-center gap-2">
        <input
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          onKeyDown={(e) => e.key === 'Enter' && onSearch(search)}
          placeholder="Search logs (client filter)"
          className="min-w-[200px] flex-1 rounded-md border border-[var(--color-kp-border)] bg-[var(--color-kp-surface-2)] px-2 py-1 text-xs"
        />
        <KPButton size="sm" variant="secondary" onClick={() => onSearch(search)}>Apply</KPButton>
        <KPButton size="sm" variant="secondary" onClick={onRefresh} disabled={loading}>Refresh</KPButton>
      </div>
      <pre className="flex-1 overflow-auto rounded-md border border-[var(--color-kp-border)] bg-[#0a0d12] p-2 text-xs leading-relaxed">
        {loading ? 'Loading logs…' : filtered.join('\n') || 'No log lines.'}
      </pre>
    </div>
  )
}

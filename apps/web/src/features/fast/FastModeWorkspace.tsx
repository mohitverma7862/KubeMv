import { useQuery } from '@tanstack/react-query'
import { useEffect, useRef } from 'react'
import { KPDataGrid } from '../../design-system/KPDataGrid'
import { KPFilterBar } from '../../design-system/KPFilterBar'
import { KPYamlPanel } from '../../design-system/KPYamlPanel'
import { KPStatusBadge } from '../../design-system/KPStatusBadge'
import { parseColonCommand } from '../../lib/commands'
import { api } from '../../services/api'
import type { ResourceRow } from '../../services/types'
import { useResourceStore } from '../../stores/resourceStore'
import { useSessionStore } from '../../stores/sessionStore'
import { useUIStore } from '../../stores/uiStore'

export function FastModeWorkspace() {
  const token = useSessionStore((s) => s.token)
  const clusterID = useUIStore((s) => s.clusterID)
  const namespace = useUIStore((s) => s.namespace)
  const { kind, filter, setFilter, setKind, selected, select, detail, setDetail } = useResourceStore()
  const filterRef = useRef<HTMLInputElement>(null)

  const { data: rows = [], refetch, isLoading } = useQuery({
    queryKey: ['resources', clusterID, kind, namespace, filter],
    enabled: Boolean(token),
    queryFn: () => {
      const colon = parseColonCommand(filter.startsWith(':') ? filter : '')
      const q = colon.filter ?? (filter.startsWith(':') ? '' : filter)
      const ns = colon.namespace ?? namespace
      const labels = q.includes('=') && !q.startsWith('/') ? q : undefined
      const search = q.startsWith('/') ? q.slice(1) : labels ? '' : q
      return api.resources(token!, clusterID, kind, {
        namespace: ns === 'all' ? '*' : ns,
        q: search,
        labels,
      })
    },
  })

  useEffect(() => {
    if (!selected || !token) return
    const ns = selected.namespace || '_'
    api.resourceDetail(token, clusterID, kind, ns, selected.name).then(setDetail).catch(() => setDetail(null))
  }, [selected, token, clusterID, kind, setDetail])

  function onFilterChange(value: string) {
    setFilter(value)
    if (value.startsWith(':')) {
      const parsed = parseColonCommand(value)
      if (parsed.kind) setKind(parsed.kind)
      if (parsed.namespace) useUIStore.getState().setNamespace(parsed.namespace)
    }
  }

  function onSelectRow(row: ResourceRow) {
    select(row)
  }

  return (
    <div className="flex h-[calc(100vh-10rem)] min-h-0 flex-col gap-3">
      <div className="flex items-center justify-between gap-3">
        <div>
          <h1 className="text-xl font-semibold capitalize">{kind}</h1>
          <p className="text-xs text-[var(--color-kp-muted)]">
            Press <kbd className="rounded border px-1">/</kbd> to filter, <kbd className="rounded border px-1">:</kbd>{' '}
            for commands, <kbd className="rounded border px-1">Ctrl+K</kbd> palette
          </p>
        </div>
        <KPStatusBadge status="healthy" label={`ns: ${namespace}`} />
      </div>

      <KPFilterBar
        value={filter}
        onChange={onFilterChange}
        onRefresh={() => refetch()}
        placeholder="Filter (/payment) or :pods -n payments /api"
      />

      <div className="grid min-h-0 flex-1 grid-cols-1 gap-3 lg:grid-cols-2 xl:grid-cols-[1.2fr_1fr_1fr]">
        <section className="min-h-0 overflow-hidden rounded-lg border border-[var(--color-kp-border)] bg-[var(--color-kp-surface)]">
          {isLoading ? <div className="p-4 text-sm text-[var(--color-kp-muted)]">Loading…</div> : null}
          <KPDataGrid rows={rows} selectedName={selected?.name} onSelect={onSelectRow} />
        </section>

        <section className="min-h-0 overflow-auto rounded-lg border border-[var(--color-kp-border)] bg-[var(--color-kp-surface)] p-4">
          <h2 className="text-sm font-medium text-[var(--color-kp-muted)]">Details</h2>
          {selected ? (
            <div className="mt-3 space-y-2 text-sm">
              <div className="text-lg font-semibold">{selected.name}</div>
              <div>Status: {selected.status}</div>
              <div>Age: {selected.age}</div>
              {selected.extra ? (
                <ul className="text-xs text-[var(--color-kp-muted)]">
                  {Object.entries(selected.extra).map(([k, v]) => (
                    <li key={k}>{k}: {v}</li>
                  ))}
                </ul>
              ) : null}
            </div>
          ) : (
            <p className="mt-3 text-sm text-[var(--color-kp-muted)]">Select a resource row.</p>
          )}
        </section>

        <section className="min-h-0 overflow-auto rounded-lg border border-[var(--color-kp-border)] bg-[var(--color-kp-surface)] p-4 lg:col-span-2 xl:col-span-1">
          <h2 className="text-sm font-medium text-[var(--color-kp-muted)]">YAML</h2>
          {detail ? <KPYamlPanel yaml={detail.yaml} /> : <p className="mt-3 text-sm text-[var(--color-kp-muted)]">YAML loads on selection.</p>}
          {detail?.events?.length ? (
            <div className="mt-4">
              <h3 className="text-xs uppercase tracking-wide text-[var(--color-kp-muted)]">Events</h3>
              <ul className="mt-2 space-y-1 text-xs">
                {detail.events.map((ev, i) => (
                  <li key={i} className="rounded border border-[var(--color-kp-border)] p-2">
                    <strong>{ev.type}</strong> {ev.reason} — {ev.message}
                  </li>
                ))}
              </ul>
            </div>
          ) : null}
        </section>
      </div>
      <input ref={filterRef} className="sr-only" aria-hidden tabIndex={-1} />
    </div>
  )
}

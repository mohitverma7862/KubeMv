import { useQuery } from '@tanstack/react-query'
import { useState } from 'react'
import { KPButton } from '../../design-system/KPButton'
import { KPFilterBar } from '../../design-system/KPFilterBar'
import { KPResourceGraph } from '../../design-system/KPResourceGraph'
import { api } from '../../services/api'
import { useSessionStore } from '../../stores/sessionStore'
import { useUIStore } from '../../stores/uiStore'

type TopologyMode = 'workload' | 'network' | 'dependency'

export function TopologyExplorer({
  initialRoot,
  initialMode = 'workload',
}: {
  initialRoot?: string
  initialMode?: TopologyMode
}) {
  const token = useSessionStore((s) => s.token)
  const clusterID = useUIStore((s) => s.clusterID)
  const namespace = useUIStore((s) => s.namespace)
  const [mode, setMode] = useState<TopologyMode>(initialMode)
  const [search, setSearch] = useState('')
  const [rootName, setRootName] = useState(initialRoot ?? 'payment-api')

  const { data, refetch, isLoading } = useQuery({
    queryKey: ['topology', clusterID, namespace, mode, rootName, search],
    enabled: Boolean(token),
    queryFn: () =>
      api.topology(token!, clusterID, {
        mode,
        namespace: namespace === 'all' ? 'payments' : namespace,
        rootName,
        q: search,
      }),
  })

  return (
    <div className="space-y-4">
      <div className="flex flex-wrap items-center justify-between gap-3">
        <div>
          <h1 className="text-2xl font-semibold">Resource topology</h1>
          <p className="text-sm text-[var(--color-kp-muted)]">XRay-style graph with health overlay (border color)</p>
        </div>
        <div className="flex gap-2">
          {(['workload', 'network', 'dependency'] as TopologyMode[]).map((m) => (
            <KPButton key={m} size="sm" variant={mode === m ? 'primary' : 'secondary'} onClick={() => setMode(m)}>
              {m}
            </KPButton>
          ))}
        </div>
      </div>

      <label className="text-sm text-[var(--color-kp-muted)]">
        Root workload
        <input
          className="ml-2 rounded border border-[var(--color-kp-border)] bg-[var(--color-kp-surface-2)] px-2 py-1"
          value={rootName}
          onChange={(e) => setRootName(e.target.value)}
        />
      </label>

      <KPFilterBar value={search} onChange={setSearch} onRefresh={() => refetch()} placeholder="Filter graph nodes" />

      {isLoading ? <p className="text-sm text-[var(--color-kp-muted)]">Building graph…</p> : null}
      {data ? <KPResourceGraph graph={data} /> : null}
    </div>
  )
}

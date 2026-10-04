import { useQuery } from '@tanstack/react-query'
import { KPHealthScore } from '../design-system/KPHealthScore'
import { KPStatusBadge } from '../design-system/KPStatusBadge'
import { api } from '../services/api'
import { useSessionStore } from '../stores/sessionStore'
import { useUIStore } from '../stores/uiStore'

export function OverviewPage() {
  const token = useSessionStore((s) => s.token)
  const clusterID = useUIStore((s) => s.clusterID)

  const { data, isLoading, error } = useQuery({
    queryKey: ['overview', clusterID],
    enabled: Boolean(token),
    queryFn: () => api.overview(token!, clusterID),
  })

  if (isLoading) return <p className="text-[var(--color-kp-muted)]">Loading cluster overview…</p>
  if (error) return <p className="text-[var(--color-kp-critical)]">{(error as Error).message}</p>
  if (!data) return null

  const healthTone = data.health.score >= 85 ? 'healthy' : data.health.score >= 70 ? 'warning' : 'critical'

  return (
    <div className="space-y-4">
      <div className="flex flex-wrap items-center justify-between gap-3">
        <div>
          <h1 className="text-2xl font-semibold">{data.cluster.name}</h1>
          <p className="text-sm text-[var(--color-kp-muted)]">Cluster overview (stub connector)</p>
        </div>
        <KPStatusBadge status={healthTone} label={healthTone} />
      </div>

      <div className="grid gap-4 md:grid-cols-4">
        <KPHealthScore score={data.health.score} />
        <MetricCard label="Nodes" value={String(data.health.nodeCount)} />
        <MetricCard label="Pods" value={String(data.health.podCount)} />
        <MetricCard label="CPU / Memory" value={`${data.health.cpuPercent}% / ${data.health.memPercent}%`} />
      </div>

      <section className="rounded-lg border border-[var(--color-kp-border)] bg-[var(--color-kp-surface)] p-4">
        <h2 className="text-sm font-medium uppercase tracking-wide text-[var(--color-kp-muted)]">Namespaces</h2>
        <ul className="mt-3 grid gap-2 md:grid-cols-2">
          {data.namespaces.map((ns) => (
            <li
              key={ns.name}
              className="flex items-center justify-between rounded-md border border-[var(--color-kp-border)] bg-[var(--color-kp-surface-2)] px-3 py-2 text-sm"
            >
              <span>{ns.name}</span>
              <KPStatusBadge status="healthy" label={ns.status} />
            </li>
          ))}
        </ul>
      </section>
    </div>
  )
}

function MetricCard({ label, value }: { label: string; value: string }) {
  return (
    <div className="rounded-lg border border-[var(--color-kp-border)] bg-[var(--color-kp-surface)] p-4">
      <div className="text-xs uppercase tracking-wide text-[var(--color-kp-muted)]">{label}</div>
      <div className="mt-1 text-2xl font-semibold tabular-nums">{value}</div>
    </div>
  )
}

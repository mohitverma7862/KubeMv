import { useQuery } from '@tanstack/react-query'
import { useMemo, useState } from 'react'
import { Link } from 'react-router-dom'
import { KPMetricChart } from '../../design-system/KPMetricChart'
import { KPButton } from '../../design-system/KPButton'
import { api } from '../../services/api'
import { useSessionStore } from '../../stores/sessionStore'
import { useUIStore } from '../../stores/uiStore'

export function ObservabilityExplorer() {
  const token = useSessionStore((s) => s.token)!
  const { clusterID, namespace } = useUIStore()
  const [workloadName, setWorkloadName] = useState('payment-api')
  const [workloadKind, setWorkloadKind] = useState('Deployment')

  const { data: dashboard, isLoading, refetch } = useQuery({
    queryKey: ['obs-dashboard', clusterID, namespace, workloadKind, workloadName],
    queryFn: () =>
      api.observabilityDashboard(token, clusterID, {
        namespace,
        kind: workloadKind,
        name: workloadName,
      }),
  })

  const downTargets = useMemo(
    () => dashboard?.targets.filter((t) => !t.up) ?? [],
    [dashboard],
  )

  return (
    <div className="flex h-full flex-col gap-4 overflow-auto p-4">
      <div className="flex flex-wrap items-end gap-3">
        <div>
          <label className="text-xs text-[var(--color-kp-muted)]">Workload kind</label>
          <select
            className="mt-1 block rounded border border-[var(--color-kp-border)] bg-[var(--color-kp-surface)] px-2 py-1"
            value={workloadKind}
            onChange={(e) => setWorkloadKind(e.target.value)}
          >
            <option>Deployment</option>
            <option>StatefulSet</option>
            <option>DaemonSet</option>
          </select>
        </div>
        <div>
          <label className="text-xs text-[var(--color-kp-muted)]">Name</label>
          <input
            className="mt-1 block rounded border border-[var(--color-kp-border)] bg-[var(--color-kp-surface)] px-2 py-1 font-mono"
            value={workloadName}
            onChange={(e) => setWorkloadName(e.target.value)}
          />
        </div>
        <KPButton size="sm" variant="secondary" onClick={() => refetch()}>
          Refresh
        </KPButton>
        {dashboard?.logDeepLink && (
          <Link to={dashboard.logDeepLink} className="text-sm text-[var(--color-kp-accent)] underline">
            Open pod logs in Fast Mode
          </Link>
        )}
      </div>

      {isLoading && <p className="text-sm text-[var(--color-kp-muted)]">Loading observability dashboard…</p>}

      {dashboard && (
        <>
          <div className="flex flex-wrap gap-2 text-sm">
            {dashboard.grafanaUrl && (
              <a className="text-[var(--color-kp-accent)] underline" href={dashboard.grafanaUrl} target="_blank" rel="noreferrer">
                Grafana
              </a>
            )}
            {dashboard.lokiUrl && (
              <a className="text-[var(--color-kp-accent)] underline" href={dashboard.lokiUrl} target="_blank" rel="noreferrer">
                Loki
              </a>
            )}
            {dashboard.prometheusUrl && (
              <span className="text-[var(--color-kp-muted)]">Prometheus: {dashboard.prometheusUrl}</span>
            )}
          </div>

          {downTargets.length > 0 && (
            <div className="rounded border border-amber-500/40 bg-amber-500/10 px-3 py-2 text-sm">
              {downTargets.length} scrape target(s) down in namespace {namespace}
            </div>
          )}

          <div className="grid gap-4 md:grid-cols-2">
            {dashboard.presets.map((preset) => (
              <div key={preset.id} className="space-y-2">
                <div className="flex items-center justify-between">
                  <h3 className="font-medium">{preset.title}</h3>
                  <span className="font-mono text-xs text-[var(--color-kp-muted)]">{preset.unit}</span>
                </div>
                <KPMetricChart points={preset.series.points} unit={preset.unit} />
                <code className="block truncate text-xs text-[var(--color-kp-muted)]">{preset.query}</code>
              </div>
            ))}
          </div>

          <section>
            <h3 className="mb-2 text-sm font-semibold uppercase tracking-wide text-[var(--color-kp-muted)]">
              Scrape targets
            </h3>
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-[var(--color-kp-border)] text-left text-[var(--color-kp-muted)]">
                  <th className="py-1">Target</th>
                  <th>Job</th>
                  <th>Instance</th>
                  <th>Status</th>
                </tr>
              </thead>
              <tbody>
                {dashboard.targets.map((t) => (
                  <tr key={`${t.namespace}/${t.name}`} className="border-b border-[var(--color-kp-border)]/50">
                    <td className="py-1 font-mono">{t.namespace}/{t.kind}/{t.name}</td>
                    <td>{t.job}</td>
                    <td className="font-mono text-xs">{t.instance}</td>
                    <td className={t.up ? 'text-emerald-500' : 'text-red-400'}>{t.up ? 'UP' : 'DOWN'}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </section>
        </>
      )}
    </div>
  )
}

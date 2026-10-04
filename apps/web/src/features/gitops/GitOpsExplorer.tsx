import { useQuery } from '@tanstack/react-query'
import { api } from '../../services/api'
import { useSessionStore } from '../../stores/sessionStore'
import { useUIStore } from '../../stores/uiStore'

const statusClass: Record<string, string> = {
  Synced: 'text-emerald-400',
  OutOfSync: 'text-amber-300',
  Unknown: 'text-[var(--color-kp-muted)]',
}

export function GitOpsExplorer() {
  const token = useSessionStore((s) => s.token)!
  const { clusterID, namespace } = useUIStore()

  const { data, isLoading, refetch } = useQuery({
    queryKey: ['gitops-overview', clusterID, namespace],
    queryFn: () => api.gitopsOverview(token, clusterID, namespace),
  })

  return (
    <div className="flex flex-col gap-4">
      <div className="flex items-center gap-3">
        <button
          type="button"
          className="rounded border border-[var(--color-kp-border)] px-3 py-1 text-sm"
          onClick={() => refetch()}
        >
          Refresh
        </button>
        <span className="text-sm text-[var(--color-kp-muted)]">
          Namespace: <span className="font-mono">{namespace}</span>
        </span>
      </div>

      {isLoading && <p className="text-sm text-[var(--color-kp-muted)]">Loading GitOps overview…</p>}

      {data && (
        <>
          <div className="grid grid-cols-2 gap-3 text-sm md:grid-cols-4">
            <Stat label="Applications" value={data.stats.applications} />
            <Stat label="Synced" value={data.stats.synced} />
            <Stat label="Out of sync" value={data.stats.outOfSync} />
            <Stat label="Drift items" value={data.stats.driftItems} />
          </div>

          <section>
            <h2 className="mb-2 text-sm font-semibold uppercase tracking-wide text-[var(--color-kp-muted)]">
              Applications
            </h2>
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-[var(--color-kp-border)] text-left text-[var(--color-kp-muted)]">
                  <th className="py-1">App</th>
                  <th>Provider</th>
                  <th>Revision</th>
                  <th>Sync</th>
                  <th>Health</th>
                </tr>
              </thead>
              <tbody>
                {data.applications.map((app) => (
                  <tr key={app.name} className="border-b border-[var(--color-kp-border)]/50">
                    <td className="py-1">
                      <div className="font-mono">{app.namespace}/{app.name}</div>
                      <div className="text-xs text-[var(--color-kp-muted)]">{app.path}</div>
                    </td>
                    <td>{app.provider}</td>
                    <td className="font-mono text-xs">{app.revision}</td>
                    <td className={statusClass[app.syncStatus] ?? statusClass.Unknown}>{app.syncStatus}</td>
                    <td>{app.health}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </section>

          <section>
            <h2 className="mb-2 text-sm font-semibold uppercase tracking-wide text-[var(--color-kp-muted)]">
              Manifest drift
            </h2>
            <div className="space-y-2">
              {data.drift.map((d) => (
                <article key={d.id} className="rounded-lg border border-[var(--color-kp-border)] bg-[var(--color-kp-surface)] px-3 py-2">
                  <div className="flex flex-wrap gap-2 text-sm">
                    <span className="font-medium">{d.resource}</span>
                    <span className="text-[var(--color-kp-muted)]">{d.field}</span>
                    <span className="text-xs uppercase text-amber-300">{d.severity}</span>
                  </div>
                  <p className="mt-1 font-mono text-xs">
                    git: {d.gitValue} → live: {d.liveValue}
                  </p>
                  <p className="mt-1 text-xs text-[var(--color-kp-muted)]">{d.suggestion}</p>
                </article>
              ))}
            </div>
          </section>

          <section>
            <h2 className="mb-2 text-sm font-semibold uppercase tracking-wide text-[var(--color-kp-muted)]">
              Pipeline hooks
            </h2>
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-[var(--color-kp-border)] text-left text-[var(--color-kp-muted)]">
                  <th className="py-1">Run</th>
                  <th>Trigger</th>
                  <th>Status</th>
                  <th>Commit</th>
                </tr>
              </thead>
              <tbody>
                {data.pipelines.map((p) => (
                  <tr key={p.id} className="border-b border-[var(--color-kp-border)]/50">
                    <td className="py-1">
                      <a className="text-[var(--color-kp-accent)] underline" href={p.url} target="_blank" rel="noreferrer">
                        {p.name} ({p.id})
                      </a>
                    </td>
                    <td>{p.trigger}</td>
                    <td>{p.status}</td>
                    <td className="font-mono text-xs">{p.commit}</td>
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

function Stat({ label, value }: { label: string; value: number }) {
  return (
    <div className="rounded-lg border border-[var(--color-kp-border)] bg-[var(--color-kp-surface)] px-3 py-2">
      <div className="text-xs text-[var(--color-kp-muted)]">{label}</div>
      <div className="text-lg font-semibold">{value}</div>
    </div>
  )
}

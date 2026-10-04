import { useQuery } from '@tanstack/react-query'
import { api } from '../../services/api'
import { useSessionStore } from '../../stores/sessionStore'
import { useUIStore } from '../../stores/uiStore'

const severityClass: Record<string, string> = {
  high: 'text-red-400 border-red-500/40 bg-red-500/10',
  medium: 'text-amber-300 border-amber-500/40 bg-amber-500/10',
  low: 'text-sky-300 border-sky-500/40 bg-sky-500/10',
}

const riskClass: Record<string, string> = {
  high: 'text-red-400',
  medium: 'text-amber-300',
  low: 'text-emerald-400',
}

export function SecurityExplorer() {
  const token = useSessionStore((s) => s.token)!
  const { clusterID, namespace } = useUIStore()

  const { data, isLoading, refetch } = useQuery({
    queryKey: ['security-summary', clusterID, namespace],
    queryFn: () => api.securitySummary(token, clusterID, namespace),
  })

  return (
    <div className="flex flex-col gap-4">
      <div className="flex flex-wrap items-center gap-3">
        <button
          type="button"
          className="rounded border border-[var(--color-kp-border)] px-3 py-1 text-sm"
          onClick={() => refetch()}
        >
          Refresh
        </button>
        <span className="text-sm text-[var(--color-kp-muted)]">
          Namespace scope: <span className="font-mono text-[var(--color-kp-text)]">{namespace}</span>
        </span>
      </div>

      {isLoading && <p className="text-sm text-[var(--color-kp-muted)]">Loading security summary…</p>}

      {data && (
        <>
          <div className="grid gap-4 md:grid-cols-[auto_1fr]">
            <div className="flex h-28 w-28 flex-col items-center justify-center rounded-xl border border-[var(--color-kp-border)] bg-[var(--color-kp-surface-elevated)]">
              <span className="text-3xl font-bold">{data.grade}</span>
              <span className="text-sm text-[var(--color-kp-muted)]">{data.score}/100</span>
            </div>
            <div className="grid grid-cols-2 gap-3 text-sm md:grid-cols-4">
              <Stat label="RoleBindings" value={data.stats.roleBindings} />
              <Stat label="ClusterRoleBindings" value={data.stats.clusterRoleBindings} />
              <Stat label="High findings" value={data.stats.highFindings} />
              <Stat label="Medium / Low" value={`${data.stats.mediumFindings} / ${data.stats.lowFindings}`} />
            </div>
          </div>

          <section>
            <h2 className="mb-2 text-sm font-semibold uppercase tracking-wide text-[var(--color-kp-muted)]">
              Policy hints
            </h2>
            <div className="space-y-2">
              {data.findings.map((f) => (
                <article
                  key={f.id}
                  className={`rounded-lg border px-3 py-2 ${severityClass[f.severity] ?? severityClass.low}`}
                >
                  <div className="flex flex-wrap items-center gap-2">
                    <span className="text-xs font-semibold uppercase">{f.severity}</span>
                    <span className="text-xs opacity-80">{f.category}</span>
                    <span className="font-medium">{f.title}</span>
                  </div>
                  <p className="mt-1 text-sm">{f.message}</p>
                  <p className="mt-1 font-mono text-xs opacity-90">{f.resource}</p>
                  <p className="mt-2 text-xs opacity-90">Remediation: {f.remediation}</p>
                </article>
              ))}
            </div>
          </section>

          <section>
            <h2 className="mb-2 text-sm font-semibold uppercase tracking-wide text-[var(--color-kp-muted)]">
              RBAC explorer
            </h2>
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-[var(--color-kp-border)] text-left text-[var(--color-kp-muted)]">
                  <th className="py-1">Binding</th>
                  <th>Role</th>
                  <th>Subjects</th>
                  <th>Risk</th>
                </tr>
              </thead>
              <tbody>
                {data.rbac.map((b) => (
                  <tr key={`${b.kind}/${b.name}`} className="border-b border-[var(--color-kp-border)]/50">
                    <td className="py-1 font-mono text-xs">
                      {b.namespace ? `${b.namespace}/` : ''}
                      {b.kind}/{b.name}
                    </td>
                    <td className="font-mono text-xs">{b.roleRef}</td>
                    <td className="text-xs">{b.subjects.join(', ')}</td>
                    <td className={riskClass[b.risk] ?? ''}>{b.risk}</td>
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

function Stat({ label, value }: { label: string; value: string | number }) {
  return (
    <div className="rounded-lg border border-[var(--color-kp-border)] bg-[var(--color-kp-surface)] px-3 py-2">
      <div className="text-xs text-[var(--color-kp-muted)]">{label}</div>
      <div className="text-lg font-semibold">{value}</div>
    </div>
  )
}

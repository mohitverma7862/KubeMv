import { useQuery, useQueryClient } from '@tanstack/react-query'
import { useState } from 'react'
import { KPButton } from '../../design-system/KPButton'
import { KPConfirmationDialog } from '../../design-system/KPConfirmationDialog'
import { KPRolloutBar } from '../../design-system/KPRolloutBar'
import { api } from '../../services/api'
import type { ResourceKindPath, ResourceRow } from '../../services/types'
import { useSessionStore } from '../../stores/sessionStore'
import { useUIStore } from '../../stores/uiStore'

const scalableKinds: ResourceKindPath[] = ['deployments', 'statefulsets']

export function WorkloadOpsPanel({ kind, workload }: { kind: ResourceKindPath; workload: ResourceRow }) {
  const token = useSessionStore((s) => s.token)
  const clusterID = useUIStore((s) => s.clusterID)
  const qc = useQueryClient()
  const [replicas, setReplicas] = useState('3')
  const [confirm, setConfirm] = useState<null | 'restart' | 'rollback'>(null)
  const [message, setMessage] = useState<string | null>(null)

  const { data: rollout } = useQuery({
    queryKey: ['rollout', clusterID, kind, workload.namespace, workload.name],
    enabled: Boolean(token) && kind === 'deployments',
    queryFn: () => api.rolloutStatus(token!, clusterID, kind, workload.namespace, workload.name),
  })

  async function runScale() {
    if (!token) return
    const result = await api.scaleWorkload(token, clusterID, kind, workload.namespace, workload.name, Number(replicas))
    setMessage(result.message)
    qc.invalidateQueries({ queryKey: ['resources'] })
  }

  async function runRestart() {
    if (!token) return
    const result = await api.restartWorkload(token, clusterID, kind, workload.namespace, workload.name)
    setMessage(result.message)
    setConfirm(null)
  }

  async function runRollback() {
    if (!token || !rollout?.revisions?.length) return
    const rev = rollout.revisions[0].revision
    const result = await api.rollbackDeployment(token, clusterID, workload.namespace, workload.name, rev)
    setMessage(result.message)
    setConfirm(null)
  }

  return (
    <section className="rounded-lg border border-[var(--color-kp-border)] bg-[var(--color-kp-surface)] p-3">
      <h2 className="text-sm font-semibold">Workload actions</h2>
      <p className="text-xs text-[var(--color-kp-muted)]">{workload.namespace}/{workload.name}</p>

      {kind === 'deployments' && rollout ? (
        <div className="mt-3">
          <h3 className="text-xs uppercase tracking-wide text-[var(--color-kp-muted)]">Rollout</h3>
          <KPRolloutBar revisions={rollout.revisions} />
        </div>
      ) : null}

      {scalableKinds.includes(kind) ? (
        <div className="mt-3 flex flex-wrap items-end gap-2">
          <label className="text-xs">Replicas
            <input className="ml-2 w-20 rounded border px-2 py-1" value={replicas} onChange={(e) => setReplicas(e.target.value)} />
          </label>
          <KPButton size="sm" onClick={runScale}>Scale</KPButton>
        </div>
      ) : null}

      {kind === 'deployments' ? (
        <div className="mt-3 flex gap-2">
          <KPButton size="sm" variant="secondary" onClick={() => setConfirm('restart')}>Restart rollout</KPButton>
          <KPButton size="sm" variant="secondary" onClick={() => setConfirm('rollback')}>Rollback</KPButton>
        </div>
      ) : null}

      {message ? <p className="mt-3 text-xs text-[var(--color-kp-healthy)]">{message}</p> : null}

      <KPConfirmationDialog
        open={confirm === 'restart'}
        title="Restart deployment"
        message={`Restart rollout for ${workload.name}?`}
        risk="LOW"
        onCancel={() => setConfirm(null)}
        onConfirm={runRestart}
      />
      <KPConfirmationDialog
        open={confirm === 'rollback'}
        title="Rollback deployment"
        message={`Rollback ${workload.name} to previous revision?`}
        risk="MEDIUM"
        onCancel={() => setConfirm(null)}
        onConfirm={runRollback}
      />
    </section>
  )
}

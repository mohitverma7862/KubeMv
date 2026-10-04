import type { RolloutRevision } from '../services/types'

export function KPRolloutBar({ revisions }: { revisions: RolloutRevision[] }) {
  return (
    <div className="space-y-3">
      {revisions.map((rev) => (
        <div key={rev.revision}>
          <div className="flex justify-between text-xs text-[var(--color-kp-muted)]">
            <span>Revision {rev.revision}</span>
            <span>{rev.status}</span>
          </div>
          <div className="mt-1 h-2 overflow-hidden rounded bg-[var(--color-kp-surface-2)]">
            <div className="h-full bg-[var(--color-kp-accent)]" style={{ width: `${rev.progress}%` }} />
          </div>
          <div className="mt-1 text-xs tabular-nums">{rev.ready}/{rev.replicas} ready</div>
        </div>
      ))}
    </div>
  )
}

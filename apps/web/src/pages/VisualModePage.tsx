import { useUIStore } from '../stores/uiStore'

export function VisualModePage() {
  const { clusterID, namespace, mode } = useUIStore()

  return (
    <div className="space-y-4">
      <h1 className="text-2xl font-semibold">Visual Mode</h1>
      <p className="text-sm text-[var(--color-kp-muted)]">
        Dashboards, topology, and charts will expand in later phases. Active UI mode: <strong>{mode}</strong>.
      </p>
      <div className="rounded-lg border border-[var(--color-kp-border)] bg-[var(--color-kp-surface)] p-4 text-sm">
        Preserved context — cluster <code>{clusterID}</code>, namespace <code>{namespace}</code>.
      </div>
    </div>
  )
}

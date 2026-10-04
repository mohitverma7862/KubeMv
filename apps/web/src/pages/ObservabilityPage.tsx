import { ObservabilityExplorer } from '../features/observability/ObservabilityExplorer'

export function ObservabilityPage() {
  return (
    <div className="flex h-full flex-col">
      <header className="border-b border-[var(--color-kp-border)] px-4 py-3">
        <h1 className="text-lg font-semibold">Observability</h1>
        <p className="text-sm text-[var(--color-kp-muted)]">
          Unified metrics, scrape health, and links to logs — Prometheus-ready when configured.
        </p>
      </header>
      <ObservabilityExplorer />
    </div>
  )
}

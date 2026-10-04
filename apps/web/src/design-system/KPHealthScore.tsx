export function KPHealthScore({ score, label = 'Cluster Health' }: { score: number; label?: string }) {
  const tone =
    score >= 85 ? 'text-[var(--color-kp-healthy)]' : score >= 70 ? 'text-[var(--color-kp-warning)]' : 'text-[var(--color-kp-critical)]'

  return (
    <div className="rounded-lg border border-[var(--color-kp-border)] bg-[var(--color-kp-surface)] p-4">
      <div className="text-xs uppercase tracking-wide text-[var(--color-kp-muted)]">{label}</div>
      <div className={['mt-1 text-3xl font-semibold tabular-nums', tone].join(' ')}>{score}/100</div>
    </div>
  )
}

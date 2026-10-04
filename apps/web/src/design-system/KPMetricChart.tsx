import type { MetricSample } from '../services/types'

function normalize(points: MetricSample[]) {
  if (points.length === 0) {
    return { path: '', min: 0, max: 1, latest: 0 }
  }
  const values = points.map((p) => p.value)
  const min = Math.min(...values)
  const max = Math.max(...values)
  const span = max - min || 1
  const width = 320
  const height = 72
  const coords = points.map((p, i) => {
    const x = (i / Math.max(points.length - 1, 1)) * width
    const y = height - ((p.value - min) / span) * (height - 8) - 4
    return `${x},${y}`
  })
  return {
    path: `M ${coords.join(' L ')}`,
    min,
    max,
    latest: values[values.length - 1],
    width,
    height,
  }
}

export function KPMetricChart({
  points,
  unit,
  accent = 'var(--color-kp-accent)',
}: {
  points: MetricSample[]
  unit?: string
  accent?: string
}) {
  const { path, min, max, latest, width, height } = normalize(points)
  return (
    <div className="rounded-md border border-[var(--color-kp-border)] bg-[var(--color-kp-surface-elevated)] p-3">
      <div className="mb-2 flex items-baseline justify-between text-sm">
        <span className="text-[var(--color-kp-muted)]">
          {min.toFixed(1)} – {max.toFixed(1)} {unit}
        </span>
        <span className="font-mono font-semibold">
          {latest.toFixed(2)} {unit}
        </span>
      </div>
      <svg viewBox={`0 0 ${width} ${height}`} className="h-20 w-full" role="img" aria-hidden>
        <path d={path} fill="none" stroke={accent} strokeWidth="2" vectorEffect="non-scaling-stroke" />
      </svg>
    </div>
  )
}

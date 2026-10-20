export function ResourceHeader({
  filter,
  liveLabel,
}: {
  filter: string
  liveLabel: string
}) {
  return (
    <div className="km-fast-header" style={{ height: 32, borderBottom: '1px solid var(--km-border)' }} data-testid="fast-table-header">
      <span>{filter ? `⚑ /filter: ${filter}` : '⚑ /filter'}</span>
      <span className="km-fast-muted" style={{ marginLeft: 'auto' }}>
        {liveLabel}
      </span>
    </div>
  )
}

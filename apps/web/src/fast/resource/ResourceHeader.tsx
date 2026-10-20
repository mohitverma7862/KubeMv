export function ResourceHeader({
  filter,
  shown,
  total,
  liveLabel,
  watchAge,
  columnsOpen,
  onToggleColumns,
}: {
  filter: string
  shown: number
  total: number
  liveLabel: string
  watchAge?: string
  columnsOpen?: boolean
  onToggleColumns?: () => void
}) {
  const count = filter ? `${shown}/${total}` : `${total}`
  return (
    <div className="km-fast-header km-fast-table-toolbar" data-testid="fast-table-header" data-filter-active={filter ? 'true' : 'false'}>
      <span>{filter ? `⚑ /filter: ${filter}` : '⚑ /filter'}</span>
      <span data-testid="fast-count">{count}</span>
      <button type="button" className="km-fast-tab" aria-expanded={columnsOpen} onClick={onToggleColumns}>
        cols
      </button>
      <span className="km-fast-muted" style={{ marginLeft: 'auto' }}>
        {watchAge ? `${watchAge}  ` : ''}
        {liveLabel}
      </span>
    </div>
  )
}

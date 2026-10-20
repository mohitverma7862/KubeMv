import { statusShortcuts } from './registry'

export function CommandBar({
  readOnly,
  portForwardCount,
  version,
}: {
  readOnly: boolean
  portForwardCount: number
  version?: string
}) {
  return (
    <footer className="km-fast-status" data-testid="fast-status">
      {statusShortcuts().map((item) => (
        <span key={item.id}>
          <kbd>{item.keys}</kbd> {item.label}
        </span>
      ))}
      <span className="km-fast-muted">ctx context</span>
      <span className="km-fast-muted">ns namespace</span>
      <span>{readOnly ? 'R readonly' : 'R writable'}</span>
      <span>N forwards {portForwardCount}</span>
      <span className="km-fast-muted" style={{ marginLeft: 'auto' }}>
        {version ?? 'KubeMv'}
      </span>
    </footer>
  )
}

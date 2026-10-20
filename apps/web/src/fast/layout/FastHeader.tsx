import type { FastHeaderModel } from '../types'

const connectionLabel: Record<FastHeaderModel['connection'], string> = {
  live: '● LIVE',
  reconnecting: '◌ RECONNECTING',
  disconnected: '✕ DISCONNECTED',
}

export function FastHeader({ model }: { model: FastHeaderModel }) {
  const connection =
    model.connection === 'reconnecting' && model.reconnectAttempt
      ? `${connectionLabel.reconnecting} ${model.reconnectAttempt}`
      : connectionLabel[model.connection]

  return (
    <header className="km-fast-header" data-testid="fast-header">
      <span className="km-fast-brand">KubeMv</span>
      <span>Fast Mode (K9s-Parity)</span>
      <span>{model.clusterName || 'no cluster'}</span>
      <span className="km-fast-muted">user: {model.userLabel || '—'}</span>
      <span className="km-fast-muted">context: {model.contextName || '—'}</span>
      <span className="km-fast-muted">namespace: {model.namespace || '—'}</span>
      {model.production ? <span className="km-fast-badge">PRODUCTION</span> : null}
      {model.readOnly ? (
        <span className="km-fast-badge" data-tone="readonly">
          READ ONLY
        </span>
      ) : null}
      <span className="km-fast-badge" data-state={model.connection} data-testid="fast-connection">
        {connection}
      </span>
      <span className="km-fast-muted" style={{ marginLeft: 'auto' }}>
        {model.clock}
      </span>
    </header>
  )
}

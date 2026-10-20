import type { Ref } from 'react'
import { statusShortcuts } from './registry'

export function CommandBar({
  readOnly,
  portForwardCount,
  version,
  contextName,
  namespace,
  commandActive,
  command,
  commandError,
  commandRef,
  onCommandChange,
  onCommandSubmit,
}: {
  readOnly: boolean
  portForwardCount: number
  version?: string
  contextName?: string
  namespace?: string
  commandActive?: boolean
  command?: string
  commandError?: string
  commandRef?: Ref<HTMLInputElement>
  onCommandChange?: (value: string) => void
  onCommandSubmit?: () => void
}) {
  return (
    <footer className="km-fast-status" data-testid="fast-status">
      {commandActive ? (
        <label>
          :
          <input
            ref={commandRef}
            aria-label="Command"
            value={command ?? ''}
            onChange={(event) => onCommandChange?.(event.target.value)}
            onKeyDown={(event) => {
              if (event.key === 'Enter') {
                event.preventDefault()
                onCommandSubmit?.()
              }
            }}
            style={{
              marginLeft: 6,
              width: 180,
              background: 'transparent',
              border: '1px solid var(--km-border)',
              borderRadius: 'var(--km-radius)',
            }}
          />
        </label>
      ) : null}
      {commandError ? <span className="km-fast-muted">{commandError}</span> : null}
      {statusShortcuts().map((item) => (
        <span key={item.id}>
          <kbd>{item.keys}</kbd> {item.label}
        </span>
      ))}
      <span className="km-fast-muted">ctx {contextName || '—'}</span>
      <span className="km-fast-muted">ns {namespace || '—'}</span>
      <span>{readOnly ? 'R readonly' : 'R writable'}</span>
      <span>N forwards {portForwardCount}</span>
      <span className="km-fast-muted" style={{ marginLeft: 'auto' }}>
        {version ?? 'KubeMv'}
      </span>
    </footer>
  )
}

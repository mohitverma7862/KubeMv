import { FastConsole } from '../fast/FastConsole'
import { useSessionStore } from '../stores/sessionStore'

export function FastModePage() {
  const clear = useSessionStore((s) => s.clear)
  return (
    <div style={{ height: '100%' }}>
      <FastConsole onSignOut={clear} />
    </div>
  )
}

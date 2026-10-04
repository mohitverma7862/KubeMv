import { FitAddon } from '@xterm/addon-fit'
import { Terminal } from '@xterm/xterm'
import { useEffect, useRef } from 'react'
import '@xterm/xterm/css/xterm.css'

export function KPTerminal({ wsUrl }: { wsUrl: string | null }) {
  const containerRef = useRef<HTMLDivElement>(null)
  const termRef = useRef<Terminal | null>(null)

  useEffect(() => {
    if (!containerRef.current || !wsUrl) return
    const term = new Terminal({
      theme: { background: '#0a0d12', foreground: '#e8edf4', cursor: '#4f8cff' },
      fontSize: 12,
      fontFamily: 'ui-monospace, SFMono-Regular, Menlo, monospace',
    })
    const fit = new FitAddon()
    term.loadAddon(fit)
    term.open(containerRef.current)
    fit.fit()
    termRef.current = term

    const ws = new WebSocket(wsUrl)
    ws.onmessage = (ev) => term.write(String(ev.data))
    ws.onopen = () => term.writeln('connected')
    ws.onclose = () => term.writeln('\r\n[session closed]')
    term.onData((data) => {
      if (ws.readyState === WebSocket.OPEN) ws.send(data)
    })

    return () => {
      ws.close()
      term.dispose()
      termRef.current = null
    }
  }, [wsUrl])

  if (!wsUrl) {
    return <p className="text-sm text-[var(--color-kp-muted)]">Select a pod to open exec.</p>
  }

  return <div ref={containerRef} className="h-64 w-full overflow-hidden rounded-md border border-[var(--color-kp-border)] p-1" />
}

import { useQuery } from '@tanstack/react-query'
import { useMemo, useState } from 'react'
import { KPButton } from '../../design-system/KPButton'
import { KPLogViewer } from '../../design-system/KPLogViewer'
import { KPTerminal } from '../../design-system/KPTerminal'
import { api } from '../../services/api'
import type { ResourceRow } from '../../services/types'
import { useSessionStore } from '../../stores/sessionStore'
import { useUIStore } from '../../stores/uiStore'

type OpsTab = 'logs' | 'exec' | 'portforward'

export function PodOpsPanel({ pod }: { pod: ResourceRow }) {
  const token = useSessionStore((s) => s.token)
  const clusterID = useUIStore((s) => s.clusterID)
  const [tab, setTab] = useState<OpsTab>('logs')
  const [container, setContainer] = useState('')
  const [previous, setPrevious] = useState(false)
  const [logSearch, setLogSearch] = useState('')
  const [localPort, setLocalPort] = useState('8080')
  const [remotePort, setRemotePort] = useState('8080')

  const { data: containers = [] } = useQuery({
    queryKey: ['containers', clusterID, pod.namespace, pod.name],
    enabled: Boolean(token),
    queryFn: () => api.podContainers(token!, clusterID, pod.namespace, pod.name),
  })

  const activeContainer = container || containers[0]?.name || 'app'

  const { data: logs = '', isLoading: logsLoading, refetch: refetchLogs } = useQuery({
    queryKey: ['logs', clusterID, pod.namespace, pod.name, activeContainer, previous, logSearch],
    enabled: Boolean(token) && tab === 'logs',
    queryFn: () =>
      api.podLogs(token!, clusterID, pod.namespace, pod.name, {
        container: activeContainer,
        previous,
        q: logSearch,
      }),
    refetchInterval: tab === 'logs' ? 3000 : false,
  })

  const { data: portForwards = [], refetch: refetchPF } = useQuery({
    queryKey: ['portforwards', clusterID],
    enabled: Boolean(token) && tab === 'portforward',
    queryFn: () => api.listPortForwards(token!, clusterID),
  })

  const execUrl = useMemo(() => {
    if (!token || tab !== 'exec') return null
    const params = new URLSearchParams({
      token,
      namespace: pod.namespace,
      pod: pod.name,
      container: activeContainer,
    })
    const proto = window.location.protocol === 'https:' ? 'wss' : 'ws'
    const host = window.location.host
    return `${proto}://${host}/api/v1/ws/clusters/${clusterID}/exec?${params}`
  }, [token, tab, clusterID, pod, activeContainer])

  async function startPortForward() {
    if (!token) return
    await api.createPortForward(token, clusterID, {
      namespace: pod.namespace,
      pod: pod.name,
      localPort: Number(localPort),
      remotePort: Number(remotePort),
    })
    refetchPF()
  }

  return (
    <section className="rounded-lg border border-[var(--color-kp-border)] bg-[var(--color-kp-surface)] p-3">
      <div className="flex flex-wrap items-center justify-between gap-2">
        <div>
          <h2 className="text-sm font-semibold">Pod operations</h2>
          <p className="text-xs text-[var(--color-kp-muted)]">{pod.namespace}/{pod.name}</p>
        </div>
        <label className="text-xs text-[var(--color-kp-muted)]">
          Container
          <select
            className="ml-2 rounded border border-[var(--color-kp-border)] bg-[var(--color-kp-surface-2)] px-2 py-1 text-[var(--color-kp-text)]"
            value={activeContainer}
            onChange={(e) => setContainer(e.target.value)}
          >
            {containers.map((c) => (
              <option key={c.name} value={c.name}>{c.name} ({c.restarts} restarts)</option>
            ))}
          </select>
        </label>
      </div>

      <div className="mt-3 flex gap-2">
        {(['logs', 'exec', 'portforward'] as OpsTab[]).map((t) => (
          <KPButton key={t} size="sm" variant={tab === t ? 'primary' : 'secondary'} onClick={() => setTab(t)}>
            {t === 'portforward' ? 'Port forward' : t.charAt(0).toUpperCase() + t.slice(1)}
          </KPButton>
        ))}
        {tab === 'logs' ? (
          <label className="ml-auto flex items-center gap-2 text-xs text-[var(--color-kp-muted)]">
            <input type="checkbox" checked={previous} onChange={(e) => setPrevious(e.target.checked)} />
            Previous container logs
          </label>
        ) : null}
      </div>

      <div className="mt-3">
        {tab === 'logs' ? (
          <KPLogViewer logs={logs} loading={logsLoading} onRefresh={() => refetchLogs()} onSearch={setLogSearch} />
        ) : null}
        {tab === 'exec' ? <KPTerminal wsUrl={execUrl} /> : null}
        {tab === 'portforward' ? (
          <div className="space-y-3 text-sm">
            <div className="flex flex-wrap items-end gap-2">
              <label className="text-xs">Local<input className="ml-1 rounded border px-2 py-1" value={localPort} onChange={(e) => setLocalPort(e.target.value)} /></label>
              <label className="text-xs">Remote<input className="ml-1 rounded border px-2 py-1" value={remotePort} onChange={(e) => setRemotePort(e.target.value)} /></label>
              <KPButton size="sm" onClick={startPortForward}>Start</KPButton>
            </div>
            <ul className="space-y-1 text-xs">
              {portForwards.map((pf) => (
                <li key={pf.id} className="rounded border border-[var(--color-kp-border)] p-2">
                  {pf.pod} {pf.localPort}→{pf.remotePort} — <a className="text-[var(--color-kp-accent)]" href={pf.url}>{pf.url}</a>
                </li>
              ))}
            </ul>
          </div>
        ) : null}
      </div>
    </section>
  )
}

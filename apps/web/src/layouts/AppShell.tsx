import { useQuery } from '@tanstack/react-query'
import { useMemo, useState } from 'react'
import { NavLink, Outlet, useNavigate } from 'react-router-dom'
import { KPClusterSelector } from '../design-system/KPClusterSelector'
import { KPCommandPalette } from '../design-system/KPCommandPalette'
import { KPCommandPaletteModal } from '../design-system/KPCommandPaletteModal'
import { KPNamespaceSelector } from '../design-system/KPNamespaceSelector'
import { KPButton } from '../design-system/KPButton'
import { useHotkeys } from '../hooks/useHotkeys'
import type { CommandDefinition } from '../lib/commands'
import { api } from '../services/api'
import { useResourceStore } from '../stores/resourceStore'
import { useSessionStore } from '../stores/sessionStore'
import { useUIStore } from '../stores/uiStore'
import type { ClusterRef } from '../services/types'

const navItems = [
  { to: '/', label: 'Overview' },
  { to: '/fast', label: 'Fast Mode' },
  { to: '/visual', label: 'Visual Mode' },
  { to: '/workloads', label: 'Workloads' },
  { to: '/network', label: 'Network' },
  { to: '/security', label: 'Security' },
]

export function AppShell({ clusters }: { clusters: ClusterRef[] }) {
  const navigate = useNavigate()
  const { principal, clear, token } = useSessionStore()
  const { mode, clusterID, namespace, setClusterID, setMode, setNamespace } = useUIStore()
  const setKind = useResourceStore((s) => s.setKind)
  const [paletteOpen, setPaletteOpen] = useState(false)

  const { data: overview } = useQuery({
    queryKey: ['overview', clusterID],
    enabled: Boolean(token),
    queryFn: () => api.overview(token!, clusterID),
  })

  const namespaceNames = useMemo(
    () => overview?.namespaces.map((n) => n.name) ?? ['default', 'payments', 'platform'],
    [overview],
  )

  useHotkeys({
    'Ctrl+K': () => setPaletteOpen(true),
    'Ctrl+R': () => window.location.reload(),
  })

  function onCommand(cmd: CommandDefinition) {
    if (cmd.kind) {
      setKind(cmd.kind)
      navigate('/fast')
      setMode('fast')
    }
  }

  return (
    <div className="grid h-full grid-rows-[auto_1fr_auto]">
      <KPCommandPaletteModal open={paletteOpen} onClose={() => setPaletteOpen(false)} onSelect={onCommand} />
      <header className="flex flex-wrap items-center gap-3 border-b border-[var(--color-kp-border)] bg-[var(--color-kp-surface)] px-4 py-3">
        <div className="text-lg font-semibold tracking-tight">KubeMv</div>
        <KPClusterSelector clusters={clusters} value={clusterID} onChange={setClusterID} />
        <KPNamespaceSelector value={namespace} onChange={setNamespace} namespaces={namespaceNames} />
        <div className="flex-1" />
        <KPCommandPalette onOpen={() => setPaletteOpen(true)} />
        <div className="flex items-center gap-2">
          <KPButton
            size="sm"
            variant={mode === 'fast' ? 'primary' : 'secondary'}
            onClick={() => {
              setMode('fast')
              navigate('/fast')
            }}
          >
            Fast
          </KPButton>
          <KPButton
            size="sm"
            variant={mode === 'visual' ? 'primary' : 'secondary'}
            onClick={() => {
              setMode('visual')
              navigate('/visual')
            }}
          >
            Visual
          </KPButton>
        </div>
        <div className="text-sm text-[var(--color-kp-muted)]">{principal?.displayName}</div>
        <KPButton size="sm" variant="ghost" onClick={clear}>Sign out</KPButton>
      </header>

      <div className="grid min-h-0 grid-cols-[220px_1fr]">
        <aside className="border-r border-[var(--color-kp-border)] bg-[var(--color-kp-surface)] p-3">
          <nav className="flex flex-col gap-1">
            {navItems.map((item) => (
              <NavLink
                key={item.to}
                to={item.to}
                className={({ isActive }) =>
                  [
                    'rounded-md px-3 py-2 text-sm',
                    isActive
                      ? 'bg-[var(--color-kp-surface-2)] text-[var(--color-kp-text)]'
                      : 'text-[var(--color-kp-muted)] hover:bg-[var(--color-kp-surface-2)] hover:text-[var(--color-kp-text)]',
                  ].join(' ')
                }
              >
                {item.label}
              </NavLink>
            ))}
          </nav>
        </aside>
        <main className="min-h-0 overflow-auto p-4">
          <Outlet />
        </main>
      </div>

      <footer className="flex items-center justify-between border-t border-[var(--color-kp-border)] bg-[var(--color-kp-surface)] px-4 py-2 text-xs text-[var(--color-kp-muted)]">
        <span>Phase 1 — K9s core navigation (stub or live cluster via kubeconfig)</span>
        <span>Fast / Visual context: {clusterID} / {namespace}</span>
      </footer>
    </div>
  )
}

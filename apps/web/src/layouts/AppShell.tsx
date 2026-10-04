import { NavLink, Outlet } from 'react-router-dom'
import { KPClusterSelector } from '../design-system/KPClusterSelector'
import { KPCommandPalette } from '../design-system/KPCommandPalette'
import { KPButton } from '../design-system/KPButton'
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
  const { principal, clear } = useSessionStore()
  const { mode, clusterID, setClusterID, setMode } = useUIStore()

  return (
    <div className="grid h-full grid-rows-[auto_1fr_auto]">
      <header className="flex flex-wrap items-center gap-3 border-b border-[var(--color-kp-border)] bg-[var(--color-kp-surface)] px-4 py-3">
        <div className="text-lg font-semibold tracking-tight">KubeMv</div>
        <KPClusterSelector clusters={clusters} value={clusterID} onChange={setClusterID} />
        <div className="flex-1" />
        <KPCommandPalette />
        <div className="flex items-center gap-2">
          <KPButton
            size="sm"
            variant={mode === 'fast' ? 'primary' : 'secondary'}
            onClick={() => setMode('fast')}
          >
            Fast
          </KPButton>
          <KPButton
            size="sm"
            variant={mode === 'visual' ? 'primary' : 'secondary'}
            onClick={() => setMode('visual')}
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
        <span>Phase 0 foundation — stub cluster data</span>
        <span>Context preserved across Fast / Visual mode switches</span>
      </footer>
    </div>
  )
}

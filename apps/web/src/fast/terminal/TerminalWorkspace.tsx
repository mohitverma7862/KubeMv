import type { WorkspaceTab } from '../types'
import { DescribeViewer } from './DescribeViewer'
import { EventViewer } from './EventViewer'
import { ExecTerminal } from './ExecTerminal'
import { LogViewer } from './LogViewer'
import { YamlViewer } from './YamlViewer'

const tabs: WorkspaceTab[] = ['logs', 'events', 'yaml', 'describe', 'terminal', 'rollout', 'history']

export function TerminalWorkspace({
  tab,
  onTab,
  logs,
  yaml,
  events,
  describeText,
}: {
  tab: WorkspaceTab
  onTab?: (tab: WorkspaceTab) => void
  logs?: string
  yaml?: string
  events?: string
  describeText?: string
}) {
  const body =
    tab === 'logs' ? (
      <LogViewer text={logs} />
    ) : tab === 'yaml' ? (
      <YamlViewer text={yaml} />
    ) : tab === 'events' ? (
      <EventViewer text={events} />
    ) : tab === 'describe' ? (
      <DescribeViewer text={describeText} />
    ) : tab === 'terminal' ? (
      <ExecTerminal />
    ) : (
      <p className="km-fast-muted">{tab} is not connected in this phase.</p>
    )

  return (
    <section data-testid="fast-workspace" style={{ display: 'flex', flexDirection: 'column', minHeight: 0, height: '100%' }}>
      <div className="km-fast-nav" style={{ height: 36 }}>
        {tabs.map((item) => (
          <button
            key={item}
            type="button"
            className="km-fast-tab"
            aria-current={item === tab ? 'page' : undefined}
            onClick={() => onTab?.(item)}
          >
            {item}
          </button>
        ))}
      </div>
      <div style={{ margin: 0, padding: 10 }}>{body}</div>
    </section>
  )
}

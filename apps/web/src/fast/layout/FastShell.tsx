import { AIActionOverlay } from '../ai/AIActionOverlay'
import { CommandPalette } from '../command/CommandPalette'
import { HelpOverlay } from '../command/HelpOverlay'
import { FilterBar } from '../filter/FilterBar'
import { ResourceInspector } from '../resource/ResourceInspector'
import { ResourceTable } from '../resource/ResourceTable'
import { ResourceTree } from '../resource/ResourceTree'
import { TerminalWorkspace } from '../terminal/TerminalWorkspace'
import type { FastShellProps } from '../types'
import { FastHeader } from './FastHeader'
import { ResourceNav } from './ResourceNav'
import { SplitLayout } from './SplitLayout'
import { StatusBar } from './StatusBar'
import '../theme/kubemv-terminal.css'

const liveLabel: Record<FastShellProps['header']['connection'], string> = {
  live: '● Live',
  reconnecting: '◌ Reconnecting',
  disconnected: '✕ Disconnected',
}

export function FastShell(props: FastShellProps) {
  const events = props.detail?.events?.map((event) => `${event.type} ${event.reason} ${event.message}`).join('\n')
  return (
    <div className="km-fast" data-testid="fast-shell">
      <FastHeader model={props.header} />
      <ResourceNav items={props.nav} activeKind={props.activeKind} onSelectKind={props.onSelectKind} />
      <div className="km-fast-workspace">
        <SplitLayout
          testId="fast-upper"
          main={
            <>
              <FilterBar value={props.filter} onChange={props.onFilterChange} />
              <ResourceTable
                columns={props.columns}
                rows={props.rows}
                selectedKey={props.selectedKey}
                onSelectRow={props.onSelectRow}
                filter={props.filter}
                liveLabel={liveLabel[props.header.connection]}
              />
            </>
          }
          side={<ResourceInspector detail={props.detail} />}
        />
        <SplitLayout
          testId="fast-lower"
          main={
            <TerminalWorkspace
              tab={props.workspaceTab}
              onTab={props.onWorkspaceTab}
              logs={props.logs}
              yaml={props.detail?.yaml}
              events={events}
              describeText={props.describeText}
            />
          }
          side={
            <>
              <ResourceTree nodes={props.tree ?? []} />
              <p className="km-fast-muted" style={{ padding: '0 10px 10px' }}>
                j/k move · Enter open · / filter · : command · ? help
              </p>
            </>
          }
        />
      </div>
      <StatusBar readOnly={props.header.readOnly} portForwardCount={props.portForwardCount ?? 0} />
      <HelpOverlay open={Boolean(props.helpOpen)} onClose={() => props.onHelpOpenChange?.(false)} />
      <CommandPalette
        open={Boolean(props.paletteOpen)}
        items={props.nav}
        onClose={() => props.onPaletteOpenChange?.(false)}
        onSelectKind={props.onSelectKind}
      />
      <AIActionOverlay
        open={Boolean(props.aiOpen)}
        finding={props.aiFinding}
        onClose={() => props.onAiOpenChange?.(false)}
      />
    </div>
  )
}

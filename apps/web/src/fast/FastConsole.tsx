import { useEffect, useRef, useState } from 'react'
import { parseColonCommand } from '../lib/commands'
import type { ResourceKindPath, ResourceRow } from '../services/types'
import { dispatchKey, focusRegion } from './command/dispatcher'
import { sampleCluster } from './data/sampleCluster'
import type { FastDataSource } from './data/source'
import { applyFilter } from './filter/FilterEngine'
import { FastShell } from './layout/FastShell'
import { moreNav, referenceNav } from './nav'
import { columnsFor } from './resource/columns'
import { rowKey } from './resource/rowKey'
import type { WorkspaceTab } from './types'

function formatClock(date: Date): string {
  return date.toLocaleTimeString(undefined, { hour12: false })
}

export function FastConsole({
  source = sampleCluster,
  clock,
  onVisual,
}: {
  source?: FastDataSource
  clock?: string
  onVisual?: () => void
}) {
  const [kind, setKind] = useState<ResourceKindPath>('pods')
  const [selectedKey, setSelectedKey] = useState<string>()
  const [filter, setFilter] = useState('')
  const [command, setCommand] = useState('')
  const [commandActive, setCommandActive] = useState(false)
  const [commandError, setCommandError] = useState('')
  const [workspaceTab, setWorkspaceTab] = useState<WorkspaceTab>('logs')
  const [helpOpen, setHelpOpen] = useState(false)
  const [paletteOpen, setPaletteOpen] = useState(false)
  const [aiOpen, setAiOpen] = useState(false)
  const [now, setNow] = useState('')
  const filterRef = useRef<HTMLInputElement>(null)
  const commandRef = useRef<HTMLInputElement>(null)
  const overlayOpen = helpOpen || paletteOpen || aiOpen

  useEffect(() => {
    if (clock) return
    const tick = () => setNow(formatClock(new Date()))
    const start = window.setTimeout(tick, 0)
    const timer = window.setInterval(tick, 1000)
    return () => {
      window.clearTimeout(start)
      window.clearInterval(timer)
    }
  }, [clock])

  useEffect(() => {
    if (commandActive) commandRef.current?.focus()
  }, [commandActive])

  useEffect(() => {
    const onKey = (event: KeyboardEvent) => {
      const action = dispatchKey(event, { region: focusRegion(event.target), overlayOpen })
      if (!action) return
      event.preventDefault()
      if (action.type === 'open-help') {
        setPaletteOpen(false)
        setAiOpen(false)
        setHelpOpen(true)
      } else if (action.type === 'open-palette') {
        setHelpOpen(false)
        setAiOpen(false)
        setPaletteOpen(true)
      } else if (action.type === 'open-ai') {
        setHelpOpen(false)
        setPaletteOpen(false)
        setAiOpen(true)
      } else if (action.type === 'enter-filter') {
        filterRef.current?.focus()
      } else if (action.type === 'enter-command') {
        setCommandError('')
        setCommandActive(true)
      } else if (action.type === 'close') {
        if (overlayOpen) {
          setHelpOpen(false)
          setPaletteOpen(false)
          setAiOpen(false)
        } else if (commandActive) {
          setCommandActive(false)
          setCommand('')
          setCommandError('')
        } else if (document.activeElement instanceof HTMLElement) {
          document.activeElement.blur()
        }
      }
    }
    window.addEventListener('keydown', onKey)
    return () => window.removeEventListener('keydown', onKey)
  }, [commandActive, overlayOpen])

  const rows = applyFilter(source.rows[kind] ?? [], filter)
  const selected = rows.find((row) => rowKey(row) === selectedKey) ?? rows[0]
  const activeKey = selected ? rowKey(selected) : undefined
  const detail = activeKey
    ? (source.details[activeKey] ?? { row: selected as ResourceRow, yaml: '# sample manifest omitted\n' })
    : null

  function selectKind(next: ResourceKindPath) {
    setKind(next)
    setSelectedKey(undefined)
    setCommandError('')
  }

  function submitCommand() {
    const raw = command.trim()
    const parsed = parseColonCommand(raw.startsWith(':') ? raw : `:${raw}`)
    if (!parsed.kind) {
      setCommandError('unknown command')
      return
    }
    selectKind(parsed.kind)
    if (parsed.filter) setFilter(parsed.filter)
    setCommand('')
    setCommandActive(false)
    setCommandError('')
  }

  const watchLabel = source.id === 'sample' ? '● Live · sample' : undefined

  return (
    <FastShell
      header={{ ...source.header, clock: clock ?? now }}
      nav={referenceNav}
      moreNav={moreNav}
      activeKind={kind}
      onSelectKind={selectKind}
      onVisual={onVisual}
      columns={columnsFor(kind)}
      rows={rows}
      selectedKey={activeKey}
      onSelectRow={(row) => setSelectedKey(rowKey(row))}
      filter={filter}
      onFilterChange={setFilter}
      filterInputRef={filterRef}
      watchLabel={watchLabel}
      detail={detail}
      tree={activeKey ? source.trees[activeKey] : []}
      workspaceTab={workspaceTab}
      onWorkspaceTab={setWorkspaceTab}
      logs={activeKey ? source.logs[activeKey] : undefined}
      describeText={activeKey ? source.describe[activeKey] : undefined}
      helpOpen={helpOpen}
      onHelpOpenChange={setHelpOpen}
      paletteOpen={paletteOpen}
      onPaletteOpenChange={setPaletteOpen}
      aiOpen={aiOpen}
      onAiOpenChange={setAiOpen}
      command={command}
      commandActive={commandActive}
      commandError={commandError}
      commandInputRef={commandRef}
      onCommandChange={setCommand}
      onCommandSubmit={submitCommand}
      portForwardCount={0}
    />
  )
}

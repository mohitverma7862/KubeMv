import { useEffect, useMemo, useRef, useState } from 'react'
import { parseColonCommand } from '../lib/commands'
import type { ResourceKindPath, ResourceRow } from '../services/types'
import { dispatchKey, focusRegion } from './command/dispatcher'
import { sampleCluster } from './data/sampleCluster'
import type { FastDataSource } from './data/source'
import { applyFilter } from './filter/FilterEngine'
import { FastShell } from './layout/FastShell'
import { moreNav, referenceNav } from './nav'
import { columnsFor } from './resource/columns'
import { freshView, goBack, visitKind, type ViewSnapshot } from './resource/navigation'
import { cycleSort, formatWatchAge, moveSelection, reconcileRows, sortRows } from './resource/tableModel'
import { rowKey } from './resource/rowKey'
import type { WorkspaceTab } from './types'

function formatClock(date: Date): string {
  return date.toLocaleTimeString(undefined, { hour12: false })
}

export function FastConsole({
  source = sampleCluster,
  clock,
  onSignOut,
}: {
  source?: FastDataSource
  clock?: string
  onSignOut?: () => void
}) {
  const [view, setView] = useState<ViewSnapshot>(() => freshView('pods'))
  const [command, setCommand] = useState('')
  const [commandActive, setCommandActive] = useState(false)
  const [commandError, setCommandError] = useState('')
  const [helpOpen, setHelpOpen] = useState(false)
  const [paletteOpen, setPaletteOpen] = useState(false)
  const [aiOpen, setAiOpen] = useState(false)
  const [now, setNow] = useState('')
  const [watchAge, setWatchAge] = useState('◷ —')
  const [scrollToken, setScrollToken] = useState(0)
  const filterRef = useRef<HTMLInputElement>(null)
  const commandRef = useRef<HTMLInputElement>(null)
  const stack = useRef<ViewSnapshot[]>([])
  const memory = useRef(new Map<ResourceKindPath, ViewSnapshot>())
  const scrollRef = useRef(0)
  const refreshedAt = useRef<number | null>(null)
  const [revision, setRevision] = useState(0)
  const overlayOpen = helpOpen || paletteOpen || aiOpen
  const kind = view.kind
  const columns = useMemo(() => columnsFor(kind), [kind])

  useEffect(() => {
    const tick = () => {
      const at = refreshedAt.current ?? Date.now()
      if (refreshedAt.current == null) refreshedAt.current = at
      setWatchAge(formatWatchAge(Date.now() - at))
      if (!clock) setNow(formatClock(new Date()))
    }
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

  const sourceRows = source.rows[kind] ?? []
  const [stableRows, setStableRows] = useState(sourceRows)
  const [trackedKind, setTrackedKind] = useState(kind)
  const [trackedRows, setTrackedRows] = useState(sourceRows)
  const [trackedRevision, setTrackedRevision] = useState(revision)
  if (trackedKind !== kind || trackedRows !== sourceRows || trackedRevision !== revision) {
    const reset = trackedKind !== kind || trackedRevision !== revision
    setTrackedKind(kind)
    setTrackedRows(sourceRows)
    setTrackedRevision(revision)
    setStableRows(reset ? sourceRows : reconcileRows(stableRows, sourceRows))
  }
  const filtered = applyFilter(stableRows, view.filter)
  const rows = sortRows(filtered, columns, view.sort)
  const selected = rows.find((row) => rowKey(row) === view.selectedKey) ?? rows[0]
  const activeKey = selected ? rowKey(selected) : undefined
  const detail = activeKey
    ? (source.details[activeKey] ?? { row: selected as ResourceRow, yaml: '# sample manifest omitted\n' })
    : null

  function capture(next: Partial<ViewSnapshot> = {}): ViewSnapshot {
    return { ...view, scrollTop: scrollRef.current, ...next }
  }

  function applyView(next: ViewSnapshot) {
    scrollRef.current = next.scrollTop
    setScrollToken((token) => token + 1)
    setView(next)
  }

  function selectKind(nextKind: ResourceKindPath, filter?: string) {
    if (nextKind === kind) {
      if (filter) setView((current) => ({ ...current, filter }))
      return
    }
    const current = capture()
    memory.current.set(kind, current)
    const step = visitKind(stack.current, current, nextKind, memory.current.get(nextKind))
    stack.current = step.stack
    setCommandError('')
    applyView(filter ? { ...step.current, filter } : step.current)
  }

  function move(action: 'next' | 'prev' | 'first' | 'last') {
    const index = rows.findIndex((row) => rowKey(row) === activeKey)
    const next = moveSelection(index, rows.length, action)
    const row = rows[next]
    if (!row) return
    setView((current) => ({ ...current, selectedKey: rowKey(row) }))
  }

  function openRow() {
    if (!activeKey || view.workspaceTab === 'describe') return
    stack.current.push(capture())
    setView((current) => ({ ...current, workspaceTab: 'describe' }))
  }

  function refresh() {
    refreshedAt.current = Date.now()
    setWatchAge('◷ 0s')
    setRevision((value) => value + 1)
  }

  function close() {
    if (overlayOpen) {
      setHelpOpen(false)
      setPaletteOpen(false)
      setAiOpen(false)
      return
    }
    if (commandActive) {
      setCommandActive(false)
      setCommand('')
      setCommandError('')
      return
    }
    if (document.activeElement instanceof HTMLInputElement) {
      document.activeElement.blur()
      return
    }
    const back = goBack(stack.current)
    if (!back) return
    memory.current.set(kind, capture())
    stack.current = back.stack
    applyView(back.current)
  }

  const actions = useRef({ selectKind, move, openRow, refresh, close })
  useEffect(() => {
    actions.current = { selectKind, move, openRow, refresh, close }
  })

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
        actions.current.close()
      } else if (action.type === 'move-next') {
        actions.current.move('next')
      } else if (action.type === 'move-prev') {
        actions.current.move('prev')
      } else if (action.type === 'first') {
        actions.current.move('first')
      } else if (action.type === 'last') {
        actions.current.move('last')
      } else if (action.type === 'open') {
        actions.current.openRow()
      } else if (action.type === 'refresh') {
        actions.current.refresh()
      }
    }
    window.addEventListener('keydown', onKey)
    return () => window.removeEventListener('keydown', onKey)
  }, [overlayOpen])

  function submitCommand() {
    const raw = command.trim()
    const parsed = parseColonCommand(raw.startsWith(':') ? raw : `:${raw}`)
    if (!parsed.kind) {
      setCommandError('unknown command')
      return
    }
    selectKind(parsed.kind, parsed.filter || undefined)
    setCommand('')
    setCommandActive(false)
    setCommandError('')
  }

  function toggleColumn(columnId: string) {
    if (columnId === 'name') return
    setView((current) => {
      const hidden = current.hiddenColumnIds.includes(columnId)
        ? current.hiddenColumnIds.filter((id) => id !== columnId)
        : [...current.hiddenColumnIds, columnId]
      return { ...current, hiddenColumnIds: hidden }
    })
  }

  const watchLabel = source.id === 'sample' ? '● Live · sample' : undefined

  return (
    <FastShell
      header={{ ...source.header, clock: clock ?? now }}
      nav={referenceNav}
      moreNav={moreNav}
      activeKind={kind}
      onSelectKind={selectKind}
      onSignOut={onSignOut}
      columns={columns}
      rows={rows}
      totalCount={stableRows.length}
      selectedKey={activeKey}
      onSelectRow={(row) => setView((current) => ({ ...current, selectedKey: rowKey(row) }))}
      filter={view.filter}
      onFilterChange={(filter) => setView((current) => ({ ...current, filter }))}
      filterInputRef={filterRef}
      watchLabel={watchLabel}
      watchAge={watchAge}
      sort={view.sort}
      onSort={(columnId) => setView((current) => ({ ...current, sort: cycleSort(current.sort, columnId) }))}
      hiddenColumnIds={view.hiddenColumnIds}
      onToggleColumn={toggleColumn}
      columnWidths={view.columnWidths}
      onColumnWidth={(columnId, width) =>
        setView((current) => ({ ...current, columnWidths: { ...current.columnWidths, [columnId]: width } }))
      }
      scrollTop={view.scrollTop}
      scrollToken={scrollToken}
      onTableScroll={(scrollTop) => {
        scrollRef.current = scrollTop
      }}
      detail={detail}
      tree={activeKey ? source.trees[activeKey] : []}
      workspaceTab={view.workspaceTab}
      onWorkspaceTab={(workspaceTab: WorkspaceTab) => setView((current) => ({ ...current, workspaceTab }))}
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

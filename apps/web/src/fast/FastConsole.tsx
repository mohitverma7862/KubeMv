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
import { kindForDigit, kindForNavLetter } from './command/navKeys'
import { nextWorkspaceTab } from './command/actions'
import { freshView, goBack, goForward, visitKind, type ViewSnapshot } from './resource/navigation'
import { cycleSort, formatWatchAge, moveSelection, reconcileRows, sortRows } from './resource/tableModel'
import { rowKey } from './resource/rowKey'
import type { WorkspaceTab } from './types'

const TABLE_PAGE_SIZE = 12

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
  const forwardStack = useRef<ViewSnapshot[]>([])
  const lastCommand = useRef('')
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
    forwardStack.current = []
    const step = visitKind(stack.current, current, nextKind, memory.current.get(nextKind))
    stack.current = step.stack
    setCommandError('')
    applyView(filter ? { ...step.current, filter } : step.current)
  }

  function move(action: 'next' | 'prev' | 'first' | 'last' | 'page-up' | 'page-down') {
    const index = rows.findIndex((row) => rowKey(row) === activeKey)
    const next = moveSelection(index, rows.length, action, TABLE_PAGE_SIZE)
    const row = rows[next]
    if (!row) return
    setView((current) => ({ ...current, selectedKey: rowKey(row) }))
  }

  function openRow(row?: ResourceRow) {
    const key = row ? rowKey(row) : activeKey
    if (!key) return
    stack.current.push(capture({ selectedKey: key }))
    forwardStack.current = []
    setView((current) => ({ ...current, selectedKey: key, workspaceTab: 'describe' }))
  }

  function setWorkspaceTab(tab: WorkspaceTab) {
    if (!activeKey) return
    setView((current) => ({ ...current, workspaceTab: tab }))
  }

  function applyResourceAction(action: string) {
    switch (action) {
      case 'logs':
        setWorkspaceTab('logs')
        break
      case 'shell':
      case 'attach':
        setWorkspaceTab('terminal')
        break
      case 'describe':
      case 'container':
        openRow()
        break
      case 'yaml':
      case 'edit':
        setWorkspaceTab('yaml')
        break
      case 'port-forward':
        setWorkspaceTab('history')
        setCommandError('port-forward: use pod ops API (Phase 8)')
        break
      case 'delete':
        setCommandError('delete: confirmation required (Phase 8)')
        break
      default:
        break
    }
  }

  function historyBack() {
    const current = capture()
    const back = goBack(stack.current)
    if (!back) return
    forwardStack.current.push(current)
    memory.current.set(kind, current)
    stack.current = back.stack
    applyView(back.current)
  }

  function historyForward() {
    const current = capture()
    const fwd = goForward(forwardStack.current)
    if (!fwd) return
    stack.current.push(current)
    forwardStack.current = fwd.forward
    applyView(fwd.current)
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
    historyBack()
  }

  const navItems = useMemo(() => [...referenceNav, ...moreNav], [])

  const actions = useRef({
    selectKind,
    move,
    openRow,
    refresh,
    close,
    historyBack,
    historyForward,
    applyResourceAction,
    setWorkspaceTab,
  })
  useEffect(() => {
    actions.current = {
      selectKind,
      move,
      openRow,
      refresh,
      close,
      historyBack,
      historyForward,
      applyResourceAction,
      setWorkspaceTab,
    }
  })

  useEffect(() => {
    const onKey = (event: KeyboardEvent) => {
      const action = dispatchKey(event, {
        region: focusRegion(event.target),
        overlayOpen,
        commandActive,
      })
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
      } else if (action.type === 'close' || action.type === 'back') {
        actions.current.close()
      } else if (action.type === 'quit') {
        onSignOut?.()
      } else if (action.type === 'move-next') {
        actions.current.move('next')
      } else if (action.type === 'move-prev') {
        actions.current.move('prev')
      } else if (action.type === 'page-up') {
        actions.current.move('page-up')
      } else if (action.type === 'page-down') {
        actions.current.move('page-down')
      } else if (action.type === 'first') {
        actions.current.move('first')
      } else if (action.type === 'last') {
        actions.current.move('last')
      } else if (action.type === 'open') {
        actions.current.openRow()
      } else if (action.type === 'refresh') {
        actions.current.refresh()
      } else if (action.type === 'history-back') {
        actions.current.historyBack()
      } else if (action.type === 'history-forward') {
        actions.current.historyForward()
      } else if (action.type === 'repeat-command') {
        if (lastCommand.current) {
          setCommand(lastCommand.current)
          setCommandActive(true)
        }
      } else if (action.type === 'command-clear') {
        setCommand('')
        setCommandError('')
      } else if (action.type === 'workspace-tab-next') {
        setView((current) => ({ ...current, workspaceTab: nextWorkspaceTab(current.workspaceTab, 1) }))
      } else if (action.type === 'workspace-tab-prev') {
        setView((current) => ({ ...current, workspaceTab: nextWorkspaceTab(current.workspaceTab, -1) }))
      } else if (action.type === 'resource-action') {
        actions.current.applyResourceAction(action.action)
      } else if (action.type === 'nav-digit') {
        const next = kindForDigit(navItems, action.digit)
        if (next) actions.current.selectKind(next)
      } else if (action.type === 'nav-letter') {
        const next = kindForNavLetter(navItems, action.letter)
        if (next) actions.current.selectKind(next)
      }
    }
    window.addEventListener('keydown', onKey)
    return () => window.removeEventListener('keydown', onKey)
  }, [overlayOpen, commandActive, navItems, onSignOut])

  function submitCommand() {
    const raw = command.trim()
    const line = raw.startsWith(':') ? raw : `:${raw}`
    lastCommand.current = line
    const parsed = parseColonCommand(line)
    if (parsed.type === 'unknown') {
      setCommandError('unknown command')
      return
    }
    if (parsed.type === 'quit') {
      onSignOut?.()
      return
    }
    if (parsed.type === 'help') {
      setHelpOpen(true)
      setCommandActive(false)
      setCommand('')
      return
    }
    if (parsed.type === 'ctx') {
      setCommandError('use cluster selector for context (:ctx)')
      return
    }
    if (parsed.type === 'clear-filter') {
      setView((current) => ({ ...current, filter: '' }))
      setCommandActive(false)
      setCommand('')
      setCommandError('')
      return
    }
    if (parsed.type === 'resource') {
      selectKind(parsed.kind, parsed.filter || undefined)
      setCommand('')
      setCommandActive(false)
      setCommandError('')
    }
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
      onActivateRow={(row) => openRow(row)}
      onPaletteCommand={(cmd) => {
        if (cmd.id === 'help') {
          setHelpOpen(true)
          return
        }
        if (cmd.id === 'quit') {
          onSignOut?.()
          return
        }
        if (cmd.id === 'clear') {
          setView((current) => ({ ...current, filter: '' }))
          return
        }
        if (cmd.id === 'ctx') {
          setCommandError('use cluster selector for context (:ctx)')
        }
      }}
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

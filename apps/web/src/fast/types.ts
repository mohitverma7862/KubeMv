import type { Ref } from 'react'
import type { EventRow, ResourceDetail, ResourceKindPath, ResourceRow } from '../services/types'

export type ConnectionState = 'live' | 'reconnecting' | 'disconnected'

export type WorkspaceTab = 'logs' | 'events' | 'yaml' | 'describe' | 'terminal' | 'rollout' | 'history'

export type FastHeaderModel = {
  clusterName: string
  userLabel: string
  contextName: string
  namespace: string
  production: boolean
  readOnly: boolean
  connection: ConnectionState
  reconnectAttempt?: number
  clock: string
}

export type ResourceNavItem = {
  id: string
  label: string
  shortcut: string
  kind?: ResourceKindPath
  available: boolean
}

export type TableColumn = {
  id: string
  header: string
  value: (row: ResourceRow) => string
}

export type SortState = { columnId: string; direction: 'asc' | 'desc' } | null

export type OwnerNode = {
  id: string
  kind: string
  name: string
  children?: OwnerNode[]
}

export type AIFinding = {
  rootCause: string
  evidence: string[]
  recommendedFix: string
  commands: string[]
}

export type FastShellProps = {
  header: FastHeaderModel
  nav: ResourceNavItem[]
  moreNav?: ResourceNavItem[]
  activeKind?: ResourceKindPath
  onSelectKind?: (kind: ResourceKindPath) => void
  onVisual?: () => void
  columns: TableColumn[]
  rows: ResourceRow[]
  totalCount?: number
  selectedKey?: string
  onSelectRow?: (row: ResourceRow) => void
  filter: string
  watchAge?: string
  sort?: SortState
  onSort?: (columnId: string) => void
  hiddenColumnIds?: string[]
  onToggleColumn?: (columnId: string) => void
  columnWidths?: Record<string, number>
  onColumnWidth?: (columnId: string, width: number) => void
  scrollTop?: number
  scrollToken?: number
  onTableScroll?: (scrollTop: number) => void
  onFilterChange?: (value: string) => void
  filterInputRef?: Ref<HTMLInputElement>
  watchLabel?: string
  command?: string
  commandActive?: boolean
  commandError?: string
  commandInputRef?: Ref<HTMLInputElement>
  onCommandChange?: (value: string) => void
  onCommandSubmit?: () => void
  detail?: ResourceDetail | null
  tree?: OwnerNode[]
  workspaceTab: WorkspaceTab
  onWorkspaceTab?: (tab: WorkspaceTab) => void
  logs?: string
  describeText?: string
  helpOpen?: boolean
  onHelpOpenChange?: (open: boolean) => void
  paletteOpen?: boolean
  onPaletteOpenChange?: (open: boolean) => void
  aiOpen?: boolean
  onAiOpenChange?: (open: boolean) => void
  aiFinding?: AIFinding | null
  portForwardCount?: number
}

export type { EventRow, ResourceDetail, ResourceKindPath, ResourceRow }

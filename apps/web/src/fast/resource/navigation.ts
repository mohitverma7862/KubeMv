import type { ResourceKindPath } from '../../services/types'
import type { SortState, WorkspaceTab } from '../types'

export type ViewSnapshot = {
  kind: ResourceKindPath
  selectedKey?: string
  filter: string
  sort: SortState
  scrollTop: number
  hiddenColumnIds: string[]
  columnWidths: Record<string, number>
  workspaceTab: WorkspaceTab
}

export function freshView(kind: ResourceKindPath): ViewSnapshot {
  return {
    kind,
    filter: '',
    sort: null,
    scrollTop: 0,
    hiddenColumnIds: [],
    columnWidths: {},
    workspaceTab: 'logs',
  }
}

export function visitKind(
  stack: ViewSnapshot[],
  current: ViewSnapshot,
  nextKind: ResourceKindPath,
  remembered?: ViewSnapshot,
): { stack: ViewSnapshot[]; current: ViewSnapshot } {
  if (nextKind === current.kind) return { stack, current }
  return {
    stack: [...stack, current],
    current: remembered ?? freshView(nextKind),
  }
}

export function goBack(stack: ViewSnapshot[]): { stack: ViewSnapshot[]; current: ViewSnapshot } | null {
  if (stack.length === 0) return null
  return { stack: stack.slice(0, -1), current: stack[stack.length - 1] }
}

export function goForward(
  forward: ViewSnapshot[],
): { forward: ViewSnapshot[]; current: ViewSnapshot } | null {
  if (forward.length === 0) return null
  return { forward: forward.slice(0, -1), current: forward[forward.length - 1] }
}

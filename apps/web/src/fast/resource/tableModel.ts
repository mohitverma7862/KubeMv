import type { ResourceRow } from '../../services/types'
import type { SortState, TableColumn } from '../types'
import { rowKey } from './rowKey'

export const ROW_HEIGHT = 32
export const OVERSCAN = 4

export function defaultColumnWidth(id: string): number {
  if (id === 'name') return 280
  if (id === 'node' || id === 'ip') return 160
  return 112
}

export function nextColumnWidth(current: number, delta: number, min = 72): number {
  return Math.max(min, Math.round(current + delta))
}

export function visibleWindow(
  count: number,
  scrollTop: number,
  viewport: number,
  rowHeight = ROW_HEIGHT,
  overscan = OVERSCAN,
): { start: number; end: number } {
  const height = viewport > 0 ? viewport : rowHeight * 12
  const start = Math.max(0, Math.floor(scrollTop / rowHeight) - overscan)
  const end = Math.min(count, Math.ceil((scrollTop + height) / rowHeight) + overscan)
  return { start, end }
}

export function moveSelection(index: number, count: number, action: 'next' | 'prev' | 'first' | 'last'): number {
  if (count <= 0) return -1
  if (action === 'first') return 0
  if (action === 'last') return count - 1
  if (index < 0) return action === 'prev' ? count - 1 : 0
  if (action === 'next') return Math.min(count - 1, index + 1)
  return Math.max(0, index - 1)
}

function ageSeconds(value: string): number | null {
  const match = /^(\d+)(s|m|h|d)$/.exec(value.trim())
  if (!match) return null
  const amount = Number(match[1])
  const unit = match[2]
  if (unit === 's') return amount
  if (unit === 'm') return amount * 60
  if (unit === 'h') return amount * 3600
  return amount * 86400
}

export function compareValues(left: string, right: string): number {
  const leftAge = ageSeconds(left)
  const rightAge = ageSeconds(right)
  if (leftAge != null && rightAge != null) return leftAge - rightAge
  const leftNumber = Number(left)
  const rightNumber = Number(right)
  if (left.trim() !== '' && right.trim() !== '' && Number.isFinite(leftNumber) && Number.isFinite(rightNumber)) {
    return leftNumber - rightNumber
  }
  return left.localeCompare(right, undefined, { numeric: true, sensitivity: 'base' })
}

export function sortRows(rows: ResourceRow[], columns: TableColumn[], sort: SortState): ResourceRow[] {
  if (!sort) return rows
  const column = columns.find((item) => item.id === sort.columnId)
  if (!column) return rows
  const direction = sort.direction === 'asc' ? 1 : -1
  return [...rows].sort((left, right) => direction * compareValues(column.value(left), column.value(right)))
}

export function cycleSort(current: SortState, columnId: string): SortState {
  if (current?.columnId !== columnId) return { columnId, direction: 'asc' }
  if (current.direction === 'asc') return { columnId, direction: 'desc' }
  return null
}

function sameRow(left: ResourceRow, right: ResourceRow): boolean {
  return (
    left.kind === right.kind &&
    left.namespace === right.namespace &&
    left.name === right.name &&
    left.status === right.status &&
    left.age === right.age &&
    JSON.stringify(left.labels ?? {}) === JSON.stringify(right.labels ?? {}) &&
    JSON.stringify(left.extra ?? {}) === JSON.stringify(right.extra ?? {})
  )
}

/** Keep object identity for rows whose contents did not change. */
export function reconcileRows(previous: ResourceRow[], next: ResourceRow[]): ResourceRow[] {
  const prior = new Map(previous.map((row) => [rowKey(row), row]))
  return next.map((row) => {
    const existing = prior.get(rowKey(row))
    if (existing && sameRow(existing, row)) return existing
    return row
  })
}

export function formatWatchAge(elapsedMs: number): string {
  const seconds = Math.max(0, Math.floor(elapsedMs / 1000))
  if (seconds < 60) return `◷ ${seconds}s`
  const minutes = Math.floor(seconds / 60)
  return `◷ ${minutes}m`
}

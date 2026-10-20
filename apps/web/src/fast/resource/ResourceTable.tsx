import { useEffect, useRef, useState } from 'react'
import type { ResourceRow } from '../../services/types'
import type { SortState, TableColumn } from '../types'
import { ResourceHeader } from './ResourceHeader'
import { ResourceRowView } from './ResourceRowView'
import { rowKey } from './rowKey'
import { ROW_HEIGHT, defaultColumnWidth, nextColumnWidth, visibleWindow } from './tableModel'

export function ResourceTable({
  columns,
  rows,
  totalCount,
  selectedKey,
  onSelectRow,
  filter,
  liveLabel,
  watchAge,
  sort = null,
  onSort,
  hiddenColumnIds = [],
  onToggleColumn,
  columnWidths = {},
  onColumnWidth,
  scrollTop = 0,
  scrollToken = 0,
  onScroll,
  viewportHeight,
}: {
  columns: TableColumn[]
  rows: ResourceRow[]
  totalCount?: number
  selectedKey?: string
  onSelectRow?: (row: ResourceRow) => void
  filter: string
  liveLabel: string
  watchAge?: string
  sort?: SortState
  onSort?: (columnId: string) => void
  hiddenColumnIds?: string[]
  onToggleColumn?: (columnId: string) => void
  columnWidths?: Record<string, number>
  onColumnWidth?: (columnId: string, width: number) => void
  scrollTop?: number
  scrollToken?: number
  onScroll?: (scrollTop: number) => void
  viewportHeight?: number
}) {
  const scroller = useRef<HTMLDivElement>(null)
  const seenToken = useRef(scrollToken)
  const [offset, setOffset] = useState(scrollTop)
  const [tokenSeen, setTokenSeen] = useState(scrollToken)
  if (tokenSeen !== scrollToken) {
    setTokenSeen(scrollToken)
    setOffset(scrollTop)
  }
  const [measured, setMeasured] = useState(0)
  const [columnsOpen, setColumnsOpen] = useState(false)
  const visibleColumns = columns.filter((column) => !hiddenColumnIds.includes(column.id))
  const widthOf = (id: string) => columnWidths[id] ?? defaultColumnWidth(id)
  const tableWidth = visibleColumns.reduce((sum, column) => sum + widthOf(column.id), 0)
  const windowed = visibleWindow(rows.length, offset, viewportHeight ?? measured)
  const selectedIndex = rows.findIndex((row) => rowKey(row) === selectedKey)

  useEffect(() => {
    const element = scroller.current
    if (!element || selectedIndex < 0) return
    const top = selectedIndex * ROW_HEIGHT
    const bottom = top + ROW_HEIGHT
    const view = element.clientHeight
    if (view <= 0) return
    if (top < element.scrollTop) element.scrollTop = top
    else if (bottom > element.scrollTop + view) element.scrollTop = bottom - view
  }, [selectedIndex])

  useEffect(() => {
    const element = scroller.current
    if (!element || seenToken.current === scrollToken) return
    element.scrollTop = offset
    seenToken.current = scrollToken
  }, [scrollToken, offset])

  useEffect(() => {
    const element = scroller.current
    if (!element || viewportHeight != null) return
    const measure = () => setMeasured(element.clientHeight)
    measure()
    if (typeof ResizeObserver === 'undefined') return
    const observer = new ResizeObserver(measure)
    observer.observe(element)
    return () => observer.disconnect()
  }, [viewportHeight])

  function resize(columnId: string, event: React.PointerEvent<HTMLSpanElement>) {
    event.preventDefault()
    event.stopPropagation()
    const startX = event.clientX
    const startWidth = widthOf(columnId)
    const move = (pointer: PointerEvent) => onColumnWidth?.(columnId, nextColumnWidth(startWidth, pointer.clientX - startX))
    const stop = () => {
      window.removeEventListener('pointermove', move)
      window.removeEventListener('pointerup', stop)
    }
    window.addEventListener('pointermove', move)
    window.addEventListener('pointerup', stop)
  }

  return (
    <div className="km-fast-table-wrap" data-testid="fast-table">
      <ResourceHeader
        filter={filter}
        shown={rows.length}
        total={totalCount ?? rows.length}
        liveLabel={liveLabel}
        watchAge={watchAge}
        columnsOpen={columnsOpen}
        onToggleColumns={() => setColumnsOpen((open) => !open)}
      />
      {columnsOpen ? (
        <ul className="km-fast-columns" data-testid="fast-columns">
          {columns.map((column) => (
            <li key={column.id}>
              <label>
                <input
                  type="checkbox"
                  checked={!hiddenColumnIds.includes(column.id)}
                  disabled={column.id === 'name'}
                  onChange={() => onToggleColumn?.(column.id)}
                />
                {column.header}
              </label>
            </li>
          ))}
        </ul>
      ) : null}
      <div
        ref={scroller}
        className="km-fast-table-scroll"
        onScroll={(event) => {
          const top = event.currentTarget.scrollTop
          setOffset(top)
          onScroll?.(top)
        }}
      >
        <div role="table" style={{ minWidth: tableWidth }}>
          <div role="row" className="km-fast-head" data-testid="fast-columns-head">
            {visibleColumns.map((column) => (
              <button
                key={column.id}
                type="button"
                role="columnheader"
                className="km-fast-head-cell"
                style={{ width: widthOf(column.id) }}
                data-sort={sort?.columnId === column.id ? sort.direction : undefined}
                onClick={() => onSort?.(column.id)}
              >
                {column.header}
                {sort?.columnId === column.id ? (sort.direction === 'asc' ? ' ▲' : ' ▼') : ''}
                <span className="km-fast-resize" onPointerDown={(event) => resize(column.id, event)} />
              </button>
            ))}
          </div>
          {rows.length === 0 ? (
            <p className="km-fast-muted" style={{ padding: 12 }}>
              No resources loaded.
            </p>
          ) : (
            <div style={{ position: 'relative', height: rows.length * ROW_HEIGHT }}>
              {rows.slice(windowed.start, windowed.end).map((row, offset) => (
                <ResourceRowView
                  key={rowKey(row)}
                  row={row}
                  columns={visibleColumns}
                  widths={columnWidths}
                  selected={rowKey(row) === selectedKey}
                  top={(windowed.start + offset) * ROW_HEIGHT}
                  onSelect={onSelectRow}
                />
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  )
}

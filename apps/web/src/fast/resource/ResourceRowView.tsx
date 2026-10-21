import { memo } from 'react'
import type { ResourceRow } from '../../services/types'
import type { TableColumn } from '../types'
import { ROW_HEIGHT, defaultColumnWidth } from './tableModel'

export const ResourceRowView = memo(function ResourceRowView({
  row,
  columns,
  selected,
  top,
  widths,
  onSelect,
  onActivate,
}: {
  row: ResourceRow
  columns: TableColumn[]
  selected: boolean
  top: number
  widths: Record<string, number>
  onSelect?: (row: ResourceRow) => void
  onActivate?: (row: ResourceRow) => void
}) {
  return (
    <div
      role="row"
      aria-selected={selected}
      data-testid="fast-row"
      className="km-fast-row"
      style={{ top, height: ROW_HEIGHT }}
      onClick={() => onSelect?.(row)}
      onDoubleClick={() => onActivate?.(row)}
    >
      {columns.map((column) => (
        <span key={column.id} role="cell" style={{ width: widths[column.id] ?? defaultColumnWidth(column.id) }}>
          {column.value(row)}
        </span>
      ))}
    </div>
  )
})

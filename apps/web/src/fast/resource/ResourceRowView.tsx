import type { ResourceRow } from '../../services/types'
import type { TableColumn } from '../types'

export function ResourceRowView({
  row,
  columns,
  selected,
  onSelect,
}: {
  row: ResourceRow
  columns: TableColumn[]
  selected: boolean
  onSelect?: (row: ResourceRow) => void
}) {
  return (
    <tr aria-selected={selected} onClick={() => onSelect?.(row)}>
      {columns.map((column) => (
        <td key={column.id}>{column.value(row)}</td>
      ))}
    </tr>
  )
}

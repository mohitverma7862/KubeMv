import type { ResourceRow } from '../../services/types'
import type { TableColumn } from '../types'
import { ResourceHeader } from './ResourceHeader'
import { ResourceRowView } from './ResourceRowView'
import { rowKey } from './rowKey'

export function ResourceTable({
  columns,
  rows,
  selectedKey,
  onSelectRow,
  filter,
  liveLabel,
}: {
  columns: TableColumn[]
  rows: ResourceRow[]
  selectedKey?: string
  onSelectRow?: (row: ResourceRow) => void
  filter: string
  liveLabel: string
}) {
  return (
    <div data-testid="fast-table">
      <ResourceHeader filter={filter} liveLabel={liveLabel} />
      {rows.length === 0 ? (
        <p className="km-fast-muted" style={{ padding: 12 }}>
          No resources loaded.
        </p>
      ) : (
        <table className="km-fast-table">
          <thead>
            <tr>
              {columns.map((column) => (
                <th key={column.id}>{column.header}</th>
              ))}
            </tr>
          </thead>
          <tbody>
            {rows.map((row) => (
              <ResourceRowView
                key={rowKey(row)}
                row={row}
                columns={columns}
                selected={rowKey(row) === selectedKey}
                onSelect={onSelectRow}
              />
            ))}
          </tbody>
        </table>
      )}
    </div>
  )
}

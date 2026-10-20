import type { ResourceRow } from '../../services/types'

export function rowKey(row: ResourceRow): string {
  return `${row.namespace}/${row.kind}/${row.name}`
}

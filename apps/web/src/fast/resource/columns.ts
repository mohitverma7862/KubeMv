import type { ResourceKindPath, ResourceRow } from '../../services/types'
import type { TableColumn } from '../types'

function field(row: ResourceRow, key: string): string {
  return row.extra?.[key] ?? '—'
}

const name: TableColumn = { id: 'name', header: 'NAME', value: (row) => row.name }
const namespace: TableColumn = { id: 'namespace', header: 'NAMESPACE', value: (row) => row.namespace || '—' }
const status: TableColumn = { id: 'status', header: 'STATUS', value: (row) => row.status }
const age: TableColumn = { id: 'age', header: 'AGE', value: (row) => row.age }

const podColumns: TableColumn[] = [
  name,
  { id: 'ready', header: 'READY', value: (row) => field(row, 'ready') },
  status,
  { id: 'restarts', header: 'RESTARTS', value: (row) => field(row, 'restarts') },
  age,
  { id: 'ip', header: 'IP', value: (row) => field(row, 'ip') },
  { id: 'node', header: 'NODE', value: (row) => field(row, 'node') },
  { id: 'cpu', header: 'CPU', value: (row) => field(row, 'cpu') },
  { id: 'mem', header: 'MEM', value: (row) => field(row, 'mem') },
]

const workloadColumns: TableColumn[] = [
  name,
  namespace,
  { id: 'ready', header: 'READY', value: (row) => field(row, 'ready') },
  status,
  age,
]

const fallbackColumns: TableColumn[] = [name, namespace, status, age]

const columnsByKind: Partial<Record<ResourceKindPath, TableColumn[]>> = {
  pods: podColumns,
  deployments: workloadColumns,
  statefulsets: workloadColumns,
  daemonsets: workloadColumns,
  jobs: workloadColumns,
}

export function columnsFor(kind: ResourceKindPath | undefined): TableColumn[] {
  if (!kind) return fallbackColumns
  return columnsByKind[kind] ?? fallbackColumns
}

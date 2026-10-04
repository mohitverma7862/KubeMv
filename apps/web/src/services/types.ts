export interface ClusterRef {
  id: string
  name: string
}

export interface ClusterHealth {
  score: number
  nodeCount: number
  podCount: number
  healthy: number
  warning: number
  failed: number
  cpuPercent: number
  memPercent: number
}

export interface Namespace {
  name: string
  status: string
  createdAt: string
}

export interface Principal {
  id: string
  displayName: string
  email: string
  roles: string[]
}

export interface ClusterOverview {
  cluster: ClusterRef
  health: ClusterHealth
  namespaces: Namespace[]
}

export type ResourceKindPath =
  | 'pods'
  | 'deployments'
  | 'services'
  | 'nodes'
  | 'events'
  | 'configmaps'
  | 'secrets'
  | 'crds'
  | 'namespaces'

export interface ResourceRow {
  kind: string
  namespace: string
  name: string
  status: string
  age: string
  labels?: Record<string, string>
  extra?: Record<string, string>
}

export interface EventRow {
  type: string
  reason: string
  message: string
  object: string
  age: string
  namespace: string
}

export interface ResourceDetail {
  row: ResourceRow
  yaml: string
  events?: EventRow[]
  related?: ResourceRow[]
}

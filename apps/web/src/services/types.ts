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

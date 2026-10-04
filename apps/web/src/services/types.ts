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
  | 'statefulsets'
  | 'daemonsets'
  | 'jobs'
  | 'cronjobs'
  | 'hpa'
  | 'pdb'

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

export interface PodContainer {
  name: string
  image: string
  ready: boolean
  restarts: number
}

export interface RolloutRevision {
  revision: number
  replicas: number
  ready: number
  progress: number
  status: string
}

export interface RolloutStatus {
  kind: string
  name: string
  namespace: string
  replicas: number
  ready: number
  updated: number
  available: number
  strategy: string
  revisions: RolloutRevision[]
}

export interface GraphNode {
  id: string
  kind: string
  name: string
  namespace: string
  status: string
  healthScore: number
}

export interface GraphEdge {
  source: string
  target: string
  relation: string
}

export interface TopologyGraph {
  mode: string
  namespace: string
  root: string
  nodes: GraphNode[]
  edges: GraphEdge[]
}

export interface MutationResult {
  action: string
  risk: string
  status: string
  message: string
  auditId: string
}

export interface PortForwardSession {
  id: string
  namespace: string
  pod: string
  localPort: number
  remotePort: number
  status: string
  url: string
}

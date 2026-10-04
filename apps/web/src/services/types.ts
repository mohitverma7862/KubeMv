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

export interface MetricSample {
  timestamp: string
  value: number
}

export interface MetricSeries {
  labels: Record<string, string>
  points: MetricSample[]
}

export interface MetricsQueryResult {
  resultType: string
  series: MetricSeries[]
}

export interface ScrapeTarget {
  namespace: string
  kind: string
  name: string
  job: string
  instance: string
  up: boolean
  lastScrape: string
  labels: Record<string, string>
}

export interface ObservabilityPreset {
  id: string
  title: string
  unit: string
  query: string
  series: MetricSeries
}

export interface SecurityStats {
  roleBindings: number
  clusterRoleBindings: number
  highFindings: number
  mediumFindings: number
  lowFindings: number
}

export interface RBACBinding {
  kind: string
  namespace?: string
  name: string
  roleRef: string
  subjects: string[]
  risk: string
}

export interface PolicyFinding {
  id: string
  severity: string
  category: string
  title: string
  message: string
  resource: string
  remediation: string
}

export interface GitOpsStats {
  applications: number
  synced: number
  outOfSync: number
  driftItems: number
}

export interface GitOpsApplication {
  name: string
  namespace: string
  provider: string
  repository: string
  revision: string
  path: string
  syncStatus: string
  health: string
  lastSyncedAt: string
}

export interface GitOpsDrift {
  id: string
  severity: string
  resource: string
  field: string
  gitValue: string
  liveValue: string
  suggestion: string
}

export interface PipelineRun {
  id: string
  name: string
  trigger: string
  status: string
  commit: string
  startedAt: string
  url: string
}

export interface AssistResource {
  namespace: string
  kind: string
  name: string
}

export interface Hypothesis {
  title: string
  likelihood: string
  evidence: string
}

export interface TriageResult {
  summary: string
  confidence: number
  hypotheses: Hypothesis[]
  signals: string[]
  disclaimer: string
}

export interface RunbookStep {
  order: number
  title: string
  command?: string
  caution?: string
}

export interface Runbook {
  id: string
  title: string
  steps: RunbookStep[]
}

export interface AutomationHook {
  id: string
  title: string
  description: string
  risk: string
  requiresApproval: boolean
  dryRunSupported: boolean
}

export interface AssistBundle {
  resource: AssistResource
  triage: TriageResult
  runbook: Runbook
  hooks: AutomationHook[]
}

export interface HookDryRunResult {
  hookId: string
  status: string
  message: string
  plannedActions: string[]
  auditId: string
}

export interface GitOpsOverview {
  namespace: string
  stats: GitOpsStats
  applications: GitOpsApplication[]
  drift: GitOpsDrift[]
  pipelines: PipelineRun[]
}

export interface SecuritySummary {
  namespace: string
  score: number
  grade: string
  stats: SecurityStats
  rbac: RBACBinding[]
  findings: PolicyFinding[]
}

export interface ObservabilityDashboard {
  namespace: string
  kind: string
  name: string
  presets: ObservabilityPreset[]
  targets: ScrapeTarget[]
  prometheusUrl?: string
  grafanaUrl?: string
  lokiUrl?: string
  logDeepLink?: string
}

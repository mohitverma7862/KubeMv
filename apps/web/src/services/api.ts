import type {
  ClusterOverview,
  ClusterRef,
  Principal,
  ResourceDetail,
  PodContainer,
  MutationResult,
  PortForwardSession,
  ResourceKindPath,
  ResourceRow,
  RolloutStatus,
  TopologyGraph,
  ObservabilityDashboard,
  MetricsQueryResult,
  ScrapeTarget,
  SecuritySummary,
  GitOpsOverview,
} from './types'

interface ApiEnvelope<T> {
  data?: T
  error?: { code: string; message: string }
}

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const response = await fetch(path, {
    ...init,
    headers: {
      'Content-Type': 'application/json',
      ...(init?.headers ?? {}),
    },
  })
  const body = (await response.json()) as ApiEnvelope<T>
  if (!response.ok || body.error) {
    throw new Error(body.error?.message ?? `request failed: ${response.status}`)
  }
  return body.data as T
}

export const api = {
  meta: () => request<Record<string, unknown>>('/api/v1/meta'),
  login: (username: string, password: string) =>
    request<{ token: string; principal: Principal }>('/api/v1/auth/login', {
      method: 'POST',
      body: JSON.stringify({ username, password }),
    }),
  me: (token: string) =>
    request<Principal>('/api/v1/auth/me', {
      headers: { Authorization: `Bearer ${token}` },
    }),
  clusters: (token: string) =>
    request<ClusterRef[]>('/api/v1/clusters', {
      headers: { Authorization: `Bearer ${token}` },
    }),
  overview: (token: string, clusterID: string) =>
    request<ClusterOverview>(`/api/v1/clusters/${clusterID}/overview`, {
      headers: { Authorization: `Bearer ${token}` },
    }),
  resources: (
    token: string,
    clusterID: string,
    kind: ResourceKindPath,
    params: { namespace?: string; q?: string; labels?: string },
  ) => {
    const search = new URLSearchParams()
    if (params.namespace) search.set('namespace', params.namespace)
    if (params.q) search.set('q', params.q)
    if (params.labels) search.set('labels', params.labels)
    const qs = search.toString()
    return request<ResourceRow[]>(
      `/api/v1/clusters/${clusterID}/resources/${kind}${qs ? `?${qs}` : ''}`,
      { headers: { Authorization: `Bearer ${token}` } },
    )
  },
  resourceDetail: (
    token: string,
    clusterID: string,
    kind: ResourceKindPath,
    namespace: string,
    name: string,
  ) =>
    request<ResourceDetail>(
      `/api/v1/clusters/${clusterID}/resources/${kind}/${namespace || '_'}/${name}`,
      { headers: { Authorization: `Bearer ${token}` } },
    ),
  podContainers: (token: string, clusterID: string, namespace: string, name: string) =>
    request<PodContainer[]>(`/api/v1/clusters/${clusterID}/pods/${namespace}/${name}/containers`, {
      headers: { Authorization: `Bearer ${token}` },
    }),
  podLogs: (
    token: string,
    clusterID: string,
    namespace: string,
    name: string,
    params: { container?: string; previous?: boolean; q?: string; tail?: number },
  ) => {
    const search = new URLSearchParams()
    if (params.container) search.set('container', params.container)
    if (params.previous) search.set('previous', 'true')
    if (params.q) search.set('q', params.q)
    if (params.tail) search.set('tail', String(params.tail))
    const qs = search.toString()
    return request<{ logs: string }>(
      `/api/v1/clusters/${clusterID}/pods/${namespace}/${name}/logs${qs ? `?${qs}` : ''}`,
      { headers: { Authorization: `Bearer ${token}` } },
    ).then((r) => r.logs)
  },
  listPortForwards: (token: string, clusterID: string) =>
    request<PortForwardSession[]>(`/api/v1/clusters/${clusterID}/portforwards`, {
      headers: { Authorization: `Bearer ${token}` },
    }),
  createPortForward: (token: string, clusterID: string, body: Record<string, unknown>) =>
    request<PortForwardSession>(`/api/v1/clusters/${clusterID}/portforwards`, {
      method: 'POST',
      headers: { Authorization: `Bearer ${token}` },
      body: JSON.stringify(body),
    }),
  rolloutStatus: (token: string, clusterID: string, kind: ResourceKindPath, namespace: string, name: string) =>
    request<RolloutStatus>(`/api/v1/clusters/${clusterID}/workloads/${kind}/${namespace}/${name}/rollout`, {
      headers: { Authorization: `Bearer ${token}` },
    }),
  scaleWorkload: (token: string, clusterID: string, kind: ResourceKindPath, namespace: string, name: string, replicas: number) =>
    request<MutationResult>(`/api/v1/clusters/${clusterID}/workloads/${kind}/${namespace}/${name}/scale`, {
      method: 'POST',
      headers: { Authorization: `Bearer ${token}` },
      body: JSON.stringify({ replicas }),
    }),
  restartWorkload: (token: string, clusterID: string, kind: ResourceKindPath, namespace: string, name: string) =>
    request<MutationResult>(`/api/v1/clusters/${clusterID}/workloads/${kind}/${namespace}/${name}/restart`, {
      method: 'POST',
      headers: { Authorization: `Bearer ${token}` },
    }),
  rollbackDeployment: (token: string, clusterID: string, namespace: string, name: string, revision: number) =>
    request<MutationResult>(
      `/api/v1/clusters/${clusterID}/workloads/deployments/${namespace}/${name}/rollback?revision=${revision}`,
      { method: 'POST', headers: { Authorization: `Bearer ${token}` } },
    ),
  topology: (
    token: string,
    clusterID: string,
    params: { mode?: string; namespace?: string; rootName?: string; q?: string },
  ) => {
    const search = new URLSearchParams()
    if (params.mode) search.set('mode', params.mode)
    if (params.namespace) search.set('namespace', params.namespace)
    if (params.rootName) search.set('rootName', params.rootName)
    if (params.q) search.set('q', params.q)
    const qs = search.toString()
    return request<TopologyGraph>(`/api/v1/clusters/${clusterID}/topology${qs ? `?${qs}` : ''}`, {
      headers: { Authorization: `Bearer ${token}` },
    })
  },
  observabilityDashboard: (
    token: string,
    clusterID: string,
    params: { namespace?: string; kind?: string; name?: string },
  ) => {
    const search = new URLSearchParams()
    if (params.namespace) search.set('namespace', params.namespace)
    if (params.kind) search.set('kind', params.kind)
    if (params.name) search.set('name', params.name)
    const qs = search.toString()
    return request<ObservabilityDashboard>(
      `/api/v1/clusters/${clusterID}/observability/dashboard${qs ? `?${qs}` : ''}`,
      { headers: { Authorization: `Bearer ${token}` } },
    )
  },
  observabilityTargets: (token: string, clusterID: string, namespace?: string) => {
    const qs = namespace ? `?namespace=${encodeURIComponent(namespace)}` : ''
    return request<ScrapeTarget[]>(`/api/v1/clusters/${clusterID}/observability/targets${qs}`, {
      headers: { Authorization: `Bearer ${token}` },
    })
  },
  gitopsOverview: (token: string, clusterID: string, namespace?: string) => {
    const qs = namespace ? `?namespace=${encodeURIComponent(namespace)}` : ''
    return request<GitOpsOverview>(`/api/v1/clusters/${clusterID}/gitops/overview${qs}`, {
      headers: { Authorization: `Bearer ${token}` },
    })
  },
  securitySummary: (token: string, clusterID: string, namespace?: string) => {
    const qs = namespace ? `?namespace=${encodeURIComponent(namespace)}` : ''
    return request<SecuritySummary>(`/api/v1/clusters/${clusterID}/security/summary${qs}`, {
      headers: { Authorization: `Bearer ${token}` },
    })
  },
  observabilityMetrics: (
    token: string,
    clusterID: string,
    params: { query: string; start?: number; end?: number; step?: number },
  ) => {
    const search = new URLSearchParams({ query: params.query })
    if (params.start) search.set('start', String(params.start))
    if (params.end) search.set('end', String(params.end))
    if (params.step) search.set('step', String(params.step))
    return request<MetricsQueryResult>(
      `/api/v1/clusters/${clusterID}/observability/metrics?${search.toString()}`,
      { headers: { Authorization: `Bearer ${token}` } },
    )
  },
}

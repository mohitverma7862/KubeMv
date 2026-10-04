import type { ClusterOverview, ClusterRef, Principal } from './types'

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
}

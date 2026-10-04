import { apiRequest } from "./client";
import type {
  Cluster,
  ClusterRegistration,
  Health,
  PluginManifest,
  Session,
} from "./types";

export async function fetchHealth(): Promise<Health> {
  return apiRequest<Health>("/api/v1/health");
}

export async function fetchSession(): Promise<Session | null> {
  try {
    return await apiRequest<Session>("/api/v1/auth/session");
  } catch (error) {
    if (error instanceof Error && "status" in error && error.status === 401) {
      return null;
    }
    throw error;
  }
}

export async function login(username: string, password: string): Promise<Session> {
  return apiRequest<Session>("/api/v1/auth/login", {
    body: { username, password },
  });
}

export async function logout(csrfToken: string): Promise<void> {
  await apiRequest<void>("/api/v1/auth/logout", {
    method: "POST",
    body: {},
    csrfToken,
  });
}

export async function fetchClusters(): Promise<Cluster[]> {
  const body = await apiRequest<{ clusters: Cluster[] }>("/api/v1/clusters");
  return body.clusters;
}

export async function registerCluster(
  registration: ClusterRegistration,
  csrfToken: string,
): Promise<Cluster> {
  return apiRequest<Cluster>("/api/v1/clusters", {
    body: registration,
    csrfToken,
  });
}

export async function removeCluster(id: string, csrfToken: string): Promise<void> {
  await apiRequest<void>(`/api/v1/clusters/${id}`, {
    method: "DELETE",
    csrfToken,
  });
}

export async function fetchPlugins(): Promise<PluginManifest[]> {
  const body = await apiRequest<{ plugins: PluginManifest[] }>("/api/v1/plugins");
  return body.plugins;
}

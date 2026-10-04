export type Health = {
  status: "ok";
  service: string;
  version: string;
  phase: number;
  capabilities: string[];
};

export type Principal = {
  username: string;
  displayName: string;
  roles: string[];
};

export type Session = {
  principal: Principal;
  expiresAt: string;
  csrfToken: string;
};

export type ClusterProvider = "generic" | "eks" | "gke" | "aks";

export type Cluster = {
  id: string;
  name: string;
  provider: ClusterProvider;
  context: string;
  connectionState: "not_connected";
  connectionDetail: string;
  createdAt: string;
};

export type ClusterRegistration = {
  name: string;
  provider: ClusterProvider;
  context: string;
  kubeconfigRef?: string;
};

export type PluginManifest = {
  id: string;
  name: string;
  version: string;
  category: string;
  description: string;
  permissions: string[];
};

export type ApiErrorBody = {
  error: {
    code: string;
    message: string;
  };
};

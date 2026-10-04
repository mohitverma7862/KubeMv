import type { ClusterProvider } from "../api/types";

export const providers: { value: ClusterProvider; label: string }[] = [
  { value: "generic", label: "Kubernetes" },
  { value: "eks", label: "Amazon EKS" },
  { value: "gke", label: "Google GKE" },
  { value: "aks", label: "Azure AKS" },
];

const namePattern = /^[A-Za-z0-9][A-Za-z0-9._-]{0,62}$/;

export type ClusterDraft = {
  name: string;
  provider: ClusterProvider;
  context: string;
  kubeconfigRef: string;
};

export function emptyClusterDraft(): ClusterDraft {
  return { name: "", provider: "generic", context: "", kubeconfigRef: "" };
}

export function validateClusterDraft(draft: ClusterDraft): string | null {
  if (!namePattern.test(draft.name.trim())) {
    return "Name must start with a letter or digit and use only letters, digits, '.', '_' or '-'.";
  }
  if (!providers.some((provider) => provider.value === draft.provider)) {
    return "Choose a supported provider.";
  }
  if (!namePattern.test(draft.context.trim())) {
    return "Context must start with a letter or digit and use only letters, digits, '.', '_' or '-'.";
  }
  const ref = draft.kubeconfigRef.trim();
  if (ref.length > 512 || /\s/.test(ref) || ref.includes("..")) {
    return "Kubeconfig reference must be a single opaque server-side token.";
  }
  const lowered = ref.toLowerCase();
  const banned = ["apiversion", "token:", "password=", "-----begin", "client-key-data"];
  if (banned.some((marker) => lowered.includes(marker))) {
    return "Kubeconfig reference cannot contain credential material.";
  }
  return null;
}

export function providerLabel(provider: string): string {
  return providers.find((item) => item.value === provider)?.label ?? provider;
}

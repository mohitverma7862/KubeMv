import { useState } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { ApiError } from "../api/client";
import { fetchClusters, registerCluster, removeCluster } from "../api/resources";
import type { ClusterProvider } from "../api/types";
import { clustersKey } from "../app/query";
import { useSessionQuery } from "../app/session";
import { displayCluster } from "../domain/clusterView";
import {
  emptyClusterDraft,
  providerLabel,
  providers,
  validateClusterDraft,
  type ClusterDraft,
} from "../domain/clusterForm";

export function ClustersPage() {
  const session = useSessionQuery();
  const queryClient = useQueryClient();
  const clusters = useQuery({ queryKey: clustersKey, queryFn: fetchClusters });
  const [draft, setDraft] = useState<ClusterDraft>(emptyClusterDraft);
  const [formError, setFormError] = useState<string | null>(null);
  const [pendingRemove, setPendingRemove] = useState<string | null>(null);
  const csrf = session.data?.csrfToken ?? "";

  const create = useMutation({
    mutationFn: () =>
      registerCluster(
        {
          name: draft.name.trim(),
          provider: draft.provider,
          context: draft.context.trim(),
          ...(draft.kubeconfigRef.trim() ? { kubeconfigRef: draft.kubeconfigRef.trim() } : {}),
        },
        csrf,
      ),
    onSuccess: async () => {
      setDraft(emptyClusterDraft());
      setFormError(null);
      await queryClient.invalidateQueries({ queryKey: clustersKey });
    },
    onError: (error) => {
      setFormError(error instanceof ApiError ? error.message : "Could not register the cluster.");
    },
  });

  const remove = useMutation({
    mutationFn: (id: string) => removeCluster(id, csrf),
    onSuccess: async () => {
      setPendingRemove(null);
      await queryClient.invalidateQueries({ queryKey: clustersKey });
    },
  });

  const rows = (clusters.data ?? []).map(displayCluster);
  const selected = rows.find((row) => row.id === pendingRemove);

  return (
    <section>
      <header className="km-page-head">
        <div>
          <p className="km-kicker">Registry</p>
          <h1>Clusters</h1>
          <p className="km-lead">
            Register a name, provider, and context. The API keeps an optional server-side reference in memory and does not return it. Nothing is dialed.
          </p>
        </div>
      </header>
      <div className="km-form-grid">
        <div className="km-table-wrap" data-testid="cluster-list">
          {clusters.isError ? (
            <div className="km-alert" role="alert">
              Could not load the registry.
            </div>
          ) : null}
          {rows.length === 0 && !clusters.isLoading ? (
            <div className="km-empty">No clusters are registered in this process.</div>
          ) : (
            <table>
              <thead>
                <tr>
                  <th>Name</th>
                  <th>Provider</th>
                  <th>Context</th>
                  <th>Connection</th>
                  <th></th>
                </tr>
              </thead>
              <tbody>
                {rows.map((row) => (
                  <tr key={row.id}>
                    <td>{row.name}</td>
                    <td>{providerLabel(row.provider)}</td>
                    <td>{row.context}</td>
                    <td>
                      <span className="km-badge" data-tone="idle" title={row.connectionDetail}>
                        Not connected
                      </span>
                    </td>
                    <td>
                      <button type="button" className="km-btn-danger" onClick={() => setPendingRemove(row.id)}>
                        Remove
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </div>
        <form
          className="km-panel"
          data-testid="cluster-form"
          onSubmit={(event) => {
            event.preventDefault();
            const message = validateClusterDraft(draft);
            if (message) {
              setFormError(message);
              return;
            }
            create.mutate();
          }}
        >
          <header>
            <h2>Register</h2>
          </header>
          {formError ? (
            <div className="km-alert" role="alert" data-testid="cluster-error">
              {formError}
            </div>
          ) : null}
          <label className="km-field">
            <span>Name</span>
            <input
              name="name"
              value={draft.name}
              onChange={(event) => setDraft({ ...draft, name: event.target.value })}
            />
          </label>
          <label className="km-field">
            <span>Provider</span>
            <select
              name="provider"
              value={draft.provider}
              onChange={(event) =>
                setDraft({ ...draft, provider: event.target.value as ClusterProvider })
              }
            >
              {providers.map((provider) => (
                <option key={provider.value} value={provider.value}>
                  {provider.label}
                </option>
              ))}
            </select>
          </label>
          <label className="km-field">
            <span>Context</span>
            <input
              name="context"
              value={draft.context}
              onChange={(event) => setDraft({ ...draft, context: event.target.value })}
            />
          </label>
          <label className="km-field">
            <span>Server-side reference</span>
            <input
              name="kubeconfigRef"
              value={draft.kubeconfigRef}
              onChange={(event) => setDraft({ ...draft, kubeconfigRef: event.target.value })}
              autoComplete="off"
            />
          </label>
          <p className="km-help">
            Optional opaque reference. Do not paste a kubeconfig. The value is not shown again after registration.
          </p>
          <button className="km-btn" type="submit" disabled={create.isPending || csrf === ""}>
            Register cluster
          </button>
        </form>
      </div>
      {selected ? (
        <div className="km-modal-back" onMouseDown={() => setPendingRemove(null)}>
          <div
            className="km-modal"
            role="dialog"
            aria-modal="true"
            aria-labelledby="remove-title"
            data-testid="remove-dialog"
            onMouseDown={(event) => event.stopPropagation()}
          >
            <header id="remove-title">Remove {selected.name}</header>
            <div className="km-dialog">
              This deletes the registry record in this API process. It does not change Kubernetes.
            </div>
            <footer className="km-actions">
              <button type="button" className="km-btn-ghost" onClick={() => setPendingRemove(null)}>
                Cancel
              </button>
              <button
                type="button"
                className="km-btn-danger"
                data-testid="confirm-remove"
                onClick={() => remove.mutate(selected.id)}
              >
                Remove record
              </button>
            </footer>
          </div>
        </div>
      ) : null}
    </section>
  );
}

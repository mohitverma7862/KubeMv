import { useQuery } from "@tanstack/react-query";
import { fetchClusters, fetchHealth, fetchPlugins } from "../api/resources";
import { clustersKey, healthKey, pluginsKey } from "../app/query";
import { useSessionQuery } from "../app/session";

export function OverviewPage() {
  const session = useSessionQuery();
  const health = useQuery({ queryKey: healthKey, queryFn: fetchHealth });
  const clusters = useQuery({ queryKey: clustersKey, queryFn: fetchClusters });
  const plugins = useQuery({ queryKey: pluginsKey, queryFn: fetchPlugins });
  const principal = session.data?.principal;
  const apiUp = health.isSuccess;

  return (
    <section>
      <header className="km-page-head">
        <div>
          <p className="km-kicker">Command center</p>
          <h1>Foundation</h1>
          <p className="km-lead">
            Session, cluster registry, and plugin host are live. Workload operations, watches, logs, and AI are outside this phase.
          </p>
        </div>
      </header>
      <div className="km-stat-grid">
        <article className="km-panel km-stat">
          <span>API</span>
          <strong data-testid="stat-api">{apiUp ? "Up" : "Down"}</strong>
        </article>
        <article className="km-panel km-stat">
          <span>Phase</span>
          <strong>{health.data?.phase ?? "—"}</strong>
        </article>
        <article className="km-panel km-stat">
          <span>Clusters</span>
          <strong data-testid="stat-clusters">{clusters.data?.length ?? "—"}</strong>
        </article>
        <article className="km-panel km-stat">
          <span>Plugins</span>
          <strong data-testid="stat-plugins">{plugins.data?.length ?? "—"}</strong>
        </article>
      </div>
      <div className="km-split">
        <article className="km-panel">
          <header>
            <h2>Platform</h2>
          </header>
          <div className="km-rows">
            <div className="km-row">
              <i className="km-dot" data-state={apiUp ? "ok" : "down"} />
              <div>
                <strong>Operator API</strong>
                <div className="km-muted">Identity, registry, and plugin manifests</div>
              </div>
              <span className="km-badge" data-tone={apiUp ? "ok" : undefined}>
                {apiUp ? "Serving" : "Unreachable"}
              </span>
            </div>
            <div className="km-row">
              <i className="km-dot" data-state="idle" />
              <div>
                <strong>Cluster registry</strong>
                <div className="km-muted">Metadata only. No kube-apiserver client in this phase.</div>
              </div>
              <span className="km-badge" data-tone="idle">
                Not connected
              </span>
            </div>
            <div className="km-row">
              <i className="km-dot" data-state="ok" />
              <div>
                <strong>Plugin host</strong>
                <div className="km-muted">In-process manifests with a fixed permission catalog.</div>
              </div>
              <span className="km-badge" data-tone="ok">
                Boundary on
              </span>
            </div>
          </div>
          {health.data ? (
            <p className="km-help" style={{ marginTop: 12 }}>
              Capabilities: {health.data.capabilities.join(" · ")}
            </p>
          ) : null}
        </article>
        <article className="km-panel">
          <header>
            <h2>Session</h2>
          </header>
          {principal ? (
            <div className="km-rows">
              <div className="km-row">
                <i className="km-dot" data-state="ok" />
                <div>Operator</div>
                <strong>{principal.displayName}</strong>
              </div>
              <div className="km-row">
                <i className="km-dot" data-state="ok" />
                <div>Roles</div>
                <span>{principal.roles.join(", ")}</span>
              </div>
              <div className="km-row">
                <i className="km-dot" data-state="idle" />
                <div>Expires</div>
                <span>{new Date(session.data?.expiresAt ?? "").toLocaleString()}</span>
              </div>
            </div>
          ) : null}
          <p className="km-help" style={{ marginTop: 12 }}>
            A signed-in session is the only gate on these routes. Role enforcement and Kubernetes RBAC are not evaluated yet, because this process does not call the cluster.
          </p>
        </article>
      </div>
    </section>
  );
}

import { useQuery } from "@tanstack/react-query";
import { fetchPlugins } from "../api/resources";
import { pluginsKey } from "../app/query";

export function PluginsPage() {
  const plugins = useQuery({ queryKey: pluginsKey, queryFn: fetchPlugins });

  return (
    <section>
      <header className="km-page-head">
        <div>
          <p className="km-kicker">Host</p>
          <h1>Plugins</h1>
          <p className="km-lead">
            Plugins are registered in process by the host. The API does not accept plugin uploads. Each manifest must stay inside the permission catalog.
          </p>
        </div>
      </header>
      <div className="km-table-wrap" data-testid="plugin-list">
        {plugins.isError ? (
          <div className="km-alert" role="alert">
            Could not load plugins.
          </div>
        ) : null}
        {plugins.data && plugins.data.length === 0 ? (
          <div className="km-empty">No plugins are loaded.</div>
        ) : null}
        {plugins.data && plugins.data.length > 0 ? (
          <table>
            <thead>
              <tr>
                <th>Name</th>
                <th>Category</th>
                <th>Version</th>
                <th>Permissions</th>
              </tr>
            </thead>
            <tbody>
              {plugins.data.map((plugin) => (
                <tr key={plugin.id}>
                  <td>
                    <strong>{plugin.name}</strong>
                    <div className="km-muted">{plugin.description}</div>
                  </td>
                  <td>{plugin.category}</td>
                  <td>{plugin.version}</td>
                  <td>{plugin.permissions.join(", ")}</td>
                </tr>
              ))}
            </tbody>
          </table>
        ) : null}
      </div>
    </section>
  );
}

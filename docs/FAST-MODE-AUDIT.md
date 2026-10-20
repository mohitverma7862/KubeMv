# Fast Mode audit

Phase 0 inspection of the repository at `54d8762` (`origin/main` after the Fast Mode TUI design docs). No product code was changed.

The visual contract is `docs/FAST-MODE-UI-DESIGN.md` and `docs/assets/kubemv-fast-mode-reference.svg`. The current Fast Mode screen does not match that contract. It is a dashboard-style resource browser inside the application shell.

## What is the product today

The live application is the Go API in `cmd/server` plus the React app in `apps/web`.

```text
apps/web (React, Vite, TanStack Query, Zustand)
        │  REST + one exec WebSocket
        ▼
cmd/server → internal/api → internal/kubernetes
        │
        ├── stub.Connector          when kubeconfig is absent
        └── kube.Connector          client-go against local kubeconfig contexts
```

`cmd/kubemv` only prints a version string. It is not a terminal UI.

`backend/` and `desktop/` are an earlier foundation (session cookies, cluster registry metadata, Tauri shell). They are not imported by `cmd/server` or `apps/web`. Fast Mode must not grow a second Kubernetes client there.

`apps/desktop/README.md` is a placeholder. There is no `apps/web/src/fast/` tree. Fast Mode currently lives in `apps/web/src/features/fast/`.

## Frontend architecture

| Area | Where | Notes |
| --- | --- | --- |
| Routes | `apps/web/src/App.tsx` | `/login`, `/` overview, `/fast`, `/visual`, `/workloads` (same page as Fast), topology, observability, security, gitops, assist |
| Shell | `apps/web/src/layouts/AppShell.tsx` | Left nav, cluster and namespace selectors, Fast/Visual buttons, footer. This is the dashboard chrome Fast Mode must not inherit as its layout |
| Fast screen | `apps/web/src/features/fast/FastModeWorkspace.tsx` | Title, hint text, filter field, three rounded panels (table, details, YAML), then pod or workload operation cards |
| Visual screen | `apps/web/src/pages/VisualModePage.tsx` | Renders `TopologyExplorer`, not a separate visual dashboard |
| Overview | `apps/web/src/pages/OverviewPage.tsx` | Health-style dashboard. Keep this as Visual/secondary, do not restyle it into Fast Mode |
| Design system | `apps/web/src/design-system/KP*.tsx` | Buttons, grid, filter bar, YAML, logs, terminal, confirmation, command palette, charts, health score |
| State | `stores/sessionStore.ts`, `uiStore.ts`, `resourceStore.ts` | Token, cluster, namespace, mode, selected kind/row |
| API client | `apps/web/src/services/api.ts` | Typed fetch wrappers for the routes below |
| Commands | `apps/web/src/lib/commands.ts` | Resource aliases and a small `:` parser |
| Hotkeys | `apps/web/src/hooks/useHotkeys.ts` | Window listener. Shell binds `Ctrl+K` (palette) and `Ctrl+R` (full page reload) |

Fast Mode gaps against the reference layout:

- No compact header with Fast Mode label, production badge, connection state, or clock.
- No resource tab row (`Pods`, `Deployments`, `Services`, …).
- Table columns are Namespace / Name / Status / Age / Extra. Pods are not `NAME | READY | STATUS | RESTARTS | AGE | IP | NODE | CPU | MEM`.
- No right inspector with health, containers, labels, owner chain, and shortcuts.
- Logs, events, YAML, describe, terminal, rollout, and history are not one lower workspace. Logs, exec, and port-forward are a separate card. Events and YAML sit in side panels.
- No relationship tree.
- No always-visible command/status bar.
- Styling uses rounded borders, card gaps, and the shared dashboard shell.

## Kubernetes services to reuse

There is no `ActionManager` type. The approved mutation boundary is `kubernetes.ClusterClient` in `internal/kubernetes/client.go`, implemented by `internal/kubernetes/kube` (live) and `internal/kubernetes/stub` (fixtures), exposed by `internal/api/router.go`.

Fast Mode must call these HTTP routes. It must not add a second client-go stack in the browser or in `backend/`.

| Capability | API | Live behavior |
| --- | --- | --- |
| Clusters / contexts | `GET /api/v1/clusters` | One record per kubeconfig context (`internal/kubernetes/kube/connector.go`) |
| Overview / namespaces | `GET /api/v1/clusters/{id}/overview` | Namespace list and health |
| List | `GET .../resources/{kind}` | Request/response list. Kinds in `internal/kubernetes/resources.go`: pods, deployments, services, nodes, events, configmaps, secrets, crds, namespaces, statefulsets, daemonsets, jobs, cronjobs, hpa, pdb |
| Detail | `GET .../resources/{kind}/{namespace}/{name}` | Row, YAML, events, related rows |
| CRDs | `GET .../crds` | CRD catalog only. Custom resource instances are not listed |
| Containers | `GET .../pods/{ns}/{name}/containers` | Name, image, ready, restarts |
| Logs | `GET .../pods/{ns}/{name}/logs` | Container, previous, tail, search. Returns a string, not a stream |
| Exec | `GET /api/v1/ws/clusters/{id}/exec` | WebSocket. Live path uses SPDY exec of `/bin/sh` with TTY. No resize message. Token is a query parameter. Stub path echoes commands |
| Port-forward | `GET/POST/DELETE .../portforwards` | Opens a localhost listener. Accepted connections write a placeholder line. They are not attached to the pod (`internal/kubernetes/portforward/manager.go`) |
| Rollout / scale / restart / rollback | `.../workloads/...` | Live client-go for deployments and statefulsets. UI confirms restart and rollback only |
| Topology | `GET .../topology` | Includes ingress objects inside the graph, not as a list resource |
| Metrics / targets | `.../observability/*` | Prometheus when configured; otherwise stub/empty |
| Security summary | `.../security/summary` | Summary view, not an RBAC explorer inside Fast Mode |
| GitOps | `.../gitops/overview` | Overview payload for the GitOps page |
| Assist | `GET .../assist/bundle`, `POST .../assist/hooks/{id}/dry-run` | Rule-based triage in `internal/ai/triage.go`. Dry-run does not mutate |

Not implemented on this boundary:

- Watch, informers, resourceVersion resume, reconnect
- Delete, patch, edit, apply
- Attach
- Describe as its own projection (callers use YAML plus events)
- Node cordon, drain, or taint
- Ingress, PVC, NetworkPolicy, RBAC, and Gateway as list kinds
- Custom resources behind a CRD
- Server-side read-only mode
- Production-context policy
- An action registry that can disable a mutation and explain why

`CredentialManager` in `client.go` is an interface only. The live connector reads the operator kubeconfig on the server. Credentials are not sent to the browser as kubeconfig files. The exec WebSocket does put the session token in the query string.

## List, watch, and discovery

Lists use `client-go` `List` with an optional label selector and limit (`listOptions` in `internal/kubernetes/kube/client.go`). Search is an in-process name substring. There is no field selector API, regex filter, or watch.

The UI loads a kind with TanStack Query and refetches only when the query key changes, except pod logs, which poll every 3 seconds while the logs tab is open. A single object change rebuilds the whole `rows` array. `KPDataGrid` renders every row. It is not virtualized.

CRD discovery exists (`ListCRDs`). Selecting a CRD does not list its objects.

## Logs, exec, describe, YAML, delete

| Workflow | Status |
| --- | --- |
| Logs | Partial. One-shot text, previous flag, container picker, search string sent to the API. No follow/pause stream, wrap, `n`/`N`, clear, or download |
| Previous logs | Partial. `LogOptions.Previous` is honored by the live log call |
| Exec | Partial. WebSocket to `/bin/sh`. `KPTerminal` is a basic socket view. No resize, no shell detection, no attach |
| Describe | Missing as a view. Detail YAML exists for pods and deployments only; other kinds return a comment (`fetchYAML`) |
| YAML | Partial and read-only. No editor, search, copy, or download contract in Fast Mode |
| Edit / delete | Missing. No API |
| Events | Partial. Object events on the detail payload, and an events list kind. Not a lower workspace tab |
| Secrets | List shows type and key count only. Secret YAML is not implemented, so values are not returned today. Pod YAML marshals the pod object and can include env values that are already on the pod |

## Context, namespace, and routing

Cluster and namespace are Zustand UI state, shared by Fast and the other pages. The shell selectors and `:ns` inside the filter string (`parseColonCommand`) both set the namespace. There is no `:ctx` command. Switching cluster does not show LIVE / RECONNECTING / DISCONNECTED, because nothing watches.

`:pods`, `:deploy`, `:svc`, and the aliases in `resourceCommands` change the list kind when typed into the filter field. There is no command-mode overlay, Tab completion, or shared shortcut registry. `Ctrl+P` is documented on the pods command but the shell listens for `Ctrl+K`.

`/` is not a mode. A filter string that starts with `/` is treated as search text. Label filters are strings that contain `=`.

Navigation history does not restore selection, filter, or scroll. `Esc` does not go back. `j`/`k`, arrows, `g`/`G`, and `Enter` do not move a table cursor.

## Dashboard code that stays Visual Mode

Keep these out of the Fast Mode shell:

- `AppShell` left navigation and page padding
- `OverviewPage` and `KPHealthScore` / `KPMetricChart`
- `pages/ObservabilityPage.tsx`, `SecurityPage.tsx`, `GitOpsPage.tsx`, `TopologyPage.tsx`, `AssistPage.tsx`
- `VisualModePage` (today this is the topology explorer; it can remain the secondary surface)

`features/fast/FastModeWorkspace.tsx`, `PodOpsPanel.tsx`, and `WorkloadOpsPanel.tsx` are the screens to replace with `apps/web/src/fast/`. Reuse their API calls. Do not restyle the rounded grid into the reference.

## Keyboard and input

Existing abstractions:

- `useHotkeys` — global key string such as `Ctrl+K`
- `KPCommandPalette` / `KPCommandPaletteModal` — mouse-and-keyboard command list for resource kinds
- `KPFilterBar` — one text field that mixes filter text and `:` commands
- `KPConfirmationDialog` — used for rollout restart and rollback

Missing for the TUI: focus scopes (table, filter, command, logs, help), a single shortcut registry, `?` help, and a status bar that shows the same shortcuts.

## AI

`/assist` is a permanent page (`features/assist/AssistExplorer.tsx`). Fast Mode does not embed that page, which already satisfies “no permanent AI panel” only because AI is a different route.

`internal/ai/triage.go` builds hypotheses and kubectl-shaped commands from status, events, and a log snippet. It does not call Gemini or another provider. Hook dry-run does not execute scale, restart, or rollback. Commands are text. Nothing auto-runs them.

The bundle is not shaped as `Root Cause / Evidence / Recommended Fix / Commands`, and `a` / `Ctrl-A` is not bound.

## Safety

- Secret list rows omit values.
- Restart and rollback ask for confirmation in the workload card. Scale does not.
- There is no stronger production confirmation and no server read-only switch.
- Port-forward and exec are reachable by any authenticated session. Exec does not check a read-only flag.
- CORS allows any origin.

## Phase map

| Next prompt phase | Uses from this audit |
| --- | --- |
| 1 UI architecture | New `apps/web/src/fast/` contracts. Leave `internal/kubernetes` alone |
| 2 Reference shell | Replace `FastModeWorkspace` composition. Mock rows only at the edges |
| 3 Table | New virtualized table. `ResourceRow` is the row model; pod columns need `extra` or a richer row |
| 4 Command / filter | Extend `lib/commands.ts` into one registry. Do not keep a second shortcut list in components |
| 5 Live data | `api.resources` today. Watch does not exist yet; Phase 5 needs a backend watch or it cannot show LIVE honestly |
| 6 Inspector / tree | Detail payload has YAML, events, related. Owner chain is not a first-class tree |
| 7 Workspace | Logs, YAML, events, exec, rollout APIs exist in partial form |
| 8 Actions | Scale, restart, rollback, port-forward, exec. Delete, edit, attach, and node ops do not |
| 9 Context / readonly | Cluster list and namespace selector exist. Read-only enforcement does not |
| 10 AI overlay | `assist/bundle` can feed an overlay. It is not the reference output shape |
| 11 Scale | Grid is not virtualized and lists are full replacements |
| 12 Polish | Compare against the SVG, not against the current dashboard |
| 13 Acceptance | Keyboard path in the prompt is not possible on the current screen |

# Fast Mode K9s parity

Status of KubeMv against K9s operator workflows. Visual similarity is not parity.

Legend: `PASS` equivalent workflow exists. `PARTIAL` a piece exists but the operator path is incomplete. `MISSING` no workflow. `DIFFERENT-BETTER` KubeMv does something K9s does not, and the K9s path still needs to be judged separately.

Target phase is the Fast Mode prompt phase that should close the gap. Phase 0 is this document.

| Workflow | Status | Evidence | Target |
| --- | --- | --- | --- |
| Startup | PARTIAL | `cmd/server` serves the API. `apps/web` opens a login page and dashboard shell, not a terminal at 120×30. `cmd/kubemv` prints version only | 2 |
| Resource navigation | PARTIAL | List kinds in `internal/kubernetes/resources.go`. UI changes kind from the filter field or command palette (`apps/web/src/lib/commands.ts`). No resource tab row, no single-key resource jumps, no history stack | 2, 3 |
| `:` commands | PARTIAL | `parseColonCommand` accepts `:pods`, `:deploy`, `:svc`, and similar inside the filter input. No command mode, autocomplete, or Tab | 4 |
| Aliases | PARTIAL | Aliases live on `resourceCommands` (`:po`, `:dp`, `:sts`, `:ds`, `:cm`). Not shared with help or the palette shortcut labels (`Ctrl+K` vs documented `Ctrl+P`) | 4 |
| `/` filtering | PARTIAL | A leading `/` in the filter string becomes `q`. Label strings containing `=` become a label selector. No live filter mode, field selector, regex, or fuzzy match | 4 |
| `j` / `k` navigation | MISSING | Table selection is click-only (`KPDataGrid`). No arrow, `g`/`G`, or `Enter` cursor | 3 |
| Context / namespace | PARTIAL | Shell cluster selector lists kubeconfig contexts. Namespace selector and `:… -n name` set namespace. No `:ctx`. `all` is sent as `*` | 9 |
| Breadcrumbs | MISSING | No resource breadcrumb or back stack. `Esc` does not restore the previous view | 3 |
| Tables | PARTIAL | Sticky header and row select exist. Columns are generic. No virtualization, sort, resize, column picker, or pod column model | 3 |
| Watch | MISSING | Lists are HTTP GET. No watch, resourceVersion, or reconnect loop in `internal/kubernetes/kube` | 5 |
| Logs | PARTIAL | `GetPodLogs` plus `KPLogViewer`, polled every 3s. No follow/pause stream, wrap, match navigation, or Ctrl-C | 7 |
| Previous logs | PARTIAL | `previous=true` is implemented in the pod logs API and the pod ops checkbox | 7 |
| Exec | PARTIAL | WebSocket exec of `/bin/sh` (`internal/api/handlers/exec_ws.go`) and `KPTerminal`. No resize. Token is in the query string. Stub clusters echo input | 7 |
| Attach | MISSING | No attach API or UI | 8 |
| Describe | MISSING | No describe projection. Operators see a short detail card plus YAML for pods and deployments only | 6, 7 |
| YAML | PARTIAL | `fetchYAML` for pods and deployments. Other kinds return a placeholder comment. Read-only, no editor | 7 |
| Edit | MISSING | No patch/apply API | 8 |
| Delete | MISSING | No resource delete API. Port-forward delete only removes a local placeholder session | 8 |
| Port-forward | PARTIAL | Create/list/delete API and a pod-ops form. Listener does not dial the pod (`portforward/manager.go` placeholder write) | 8 |
| Events | PARTIAL | Events list kind and per-object events on detail. Not a workspace tab with severity timeline | 7 |
| Rollout | PARTIAL | Status, restart, rollback, and scale for deployments; scale/restart paths for statefulsets. Confirmation covers restart and rollback only. Scale has no confirm | 7, 8 |
| Node operations | MISSING | Nodes can be listed. No cordon, drain, or taint | 8 |
| Multi-select | MISSING | One selected row | 3 |
| Help / hotkeys | PARTIAL | Static hint text and `Ctrl+K` palette. No `?` overlay. `Ctrl+R` reloads the browser | 4 |
| Custom views | MISSING | No saved skins, columns, or K9s-style views | 4 |
| CRDs | PARTIAL | `ListCRDs` and a CRD table. Instances of those CRDs are not listed | 5 |
| Multi-cluster | PARTIAL | Multiple kubeconfig contexts appear as clusters and the selector switches the client. No cluster screen, health comparison, or reconnect | 5, 9 |
| Reconnect | MISSING | No watch connection to resume. A failed list surfaces as a query error | 5 |

## Different, not parity

These exist beside K9s and must stay out of the Fast Mode frame until their own phase:

| Area | Status | Evidence |
| --- | --- | --- |
| Topology graph | DIFFERENT-BETTER | `GetTopology` and `TopologyExplorer`. Fast Mode should use a compact owner tree, not this graph |
| Observability page | DIFFERENT-BETTER | Prometheus query, targets, dashboard endpoints |
| Security summary | DIFFERENT-BETTER | `GetSecuritySummary` page. Not an effective-permission explainer |
| GitOps overview | DIFFERENT-BETTER | `GetGitOpsOverview` page |
| Rule-based assist | DIFFERENT-BETTER | `internal/ai/triage.go` and `/assist`. Permanent page, not a Fast Mode overlay. No external model |

## Safety gaps that block calling an action done

- No `ActionManager`. Mutations are methods on `ClusterClient`. New Fast Mode actions must go through that interface and its HTTP handlers, not a new client.
- No server read-only mode. Hiding a button would not enforce it.
- No production confirmation policy.
- Secret values are not on the secret list. Full pod YAML is not redacted.
- Exec and port-forward are not audited as operator actions in the Fast Mode sense. The in-process audit recorder under `backend/internal/audit` is a different, unused process.

## Phase 0 result

Parity is not met. The next allowed change is Phase 1: component contracts under `apps/web/src/fast/` that match the reference regions, without a backend rewrite.

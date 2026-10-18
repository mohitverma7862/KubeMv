# KubeMv — Fast Mode Master Specification

## Status

Implementation contract for rebuilding KubeMv Fast Mode as a K9s-parity, keyboard-first Kubernetes operator console.

## Product objective

Fast Mode is the primary operator surface. It combines K9s-style terminal speed with a modern terminal IDE, Kubernetes resource intelligence, live observability and contextual AI.

It is not a dashboard redesign. The existing dashboard becomes Visual Mode and remains available separately.

## Non-negotiable rules

- Preserve existing working Kubernetes abstractions wherever possible.
- Do not duplicate Kubernetes API/watch/action logic.
- All mutations flow through one backend ActionManager.
- Fast Mode is keyboard-first; mouse is secondary.
- Minimum viewport: 100x25; target: 120x30.
- Dense monospace UI; no marketing-card layout.
- Existing Visual Mode remains available.
- AI is contextual and never auto-executes commands.
- Secret values and credentials are never exposed to AI or the UI.
- Production mutations require stronger confirmation.
- Read-only mode is enforced by the backend, not just hidden buttons.
- Work phase-by-phase with small diffs.

## Architecture target

### Backend

```text
internal/
├── k8s/
│   ├── service/
│   ├── discovery/
│   ├── watch/
│   └── actions/
└── session/
```

Responsibilities:

- Kubernetes clients and resource abstraction
- discovery and CRD support
- list/get/watch operations
- logs/exec/attach/port-forward
- namespace/context session state
- one ActionManager for mutations
- readonly and production safety

### Frontend

```text
apps/web/src/fast/
├── layout/
├── registry/
├── nav/
├── filter/
├── selection/
├── state/
├── resource/
├── terminal/
├── command/
├── views/
├── ai/
└── theme/
```

## Phase 0 — audit

No implementation changes.

Inspect `cmd/kubemv`, `cmd/server`, `internal/*`, `backend`, `api`, `apps/web`, `desktop`, tests, AI code and docs.

Create:

- `docs/FAST-MODE-AUDIT.md`
- `docs/FAST-MODE-K9S-PARITY.md`

Audit existing Kubernetes list/watch/logs/exec/attach/describe/edit/delete/port-forward/discovery/context/namespace abstractions and identify what can be reused.

K9s parity matrix must cover startup, resource navigation, command mode, aliases, filters, keyboard navigation, context/namespace switching, breadcrumbs, tables, live watch, logs, previous logs, exec, attach, describe, YAML, edit, delete, port-forward, events, rollout, node operations, multi-select, help, hotkeys, custom views, CRDs, multi-cluster and reconnect/error handling.

STOP.

## Phase 1 — architecture skeleton

Establish clean Fast Mode boundaries without rewriting existing Kubernetes logic.

Backend targets:

- `internal/k8s/service`
- `internal/k8s/discovery`
- `internal/k8s/watch`
- `internal/k8s/actions`
- `internal/session`

Frontend targets:

- registry
- nav
- filter
- selection
- state
- views

Add interfaces/types and registry tests. Avoid speculative abstractions.

STOP after tests/build.

## Phase 2 — terminal shell

Fast Mode is the default `/` operator console.

Layout:

```text
Header
Resource Navigation
┌───────────────────────────────┬──────────────────────┐
│ Main Resource Table           │ Resource Inspector   │
├───────────────────────────────┼──────────────────────┤
│ Logs / Events / YAML / ...    │ Resource Tree        │
└───────────────────────────────┴──────────────────────┘
Command / Status Bar
```

Header contains KubeMv, mode, cluster, context, namespace, production state, connection state and clock.

Implement dark, light and high-contrast themes. Use realistic mock Pods first. Visual Mode remains reachable at `/visual`.

STOP after tests/build.

## Phase 3 — command mode and resource registry

Implement resource registry and commands. Required commands include `:q`, `:ns`, `:ctx`, `:events` plus resource aliases. Add dynamic CRD discovery. `:` must provide autocomplete, Tab completion, history, Enter and Esc.

STOP after tests/build.

## Phase 4 — table, filtering and navigation

Implement virtualized resource table, sticky header, keyboard navigation and NavigationStack.

Keys:

- `j/k` or arrows: move
- `g/G`: first/last
- Enter: drill down
- Esc: back
- Ctrl-R: refresh
- `/`: filter
- `:`: command
- `?`: help
- Ctrl-P: command palette

FilterEngine should support text, regex/fuzzy, labels, status and field-aware filters. Preserve selection and scroll state when returning to a view.

STOP after tests/build.

## Phase 5 — live watch

Connect existing Kubernetes watch services through SSE/WebSocket or the repository's established streaming mechanism.

Requirements:

- live row updates
- keyed row patches
- connection indicator
- reconnect with backoff
- resourceVersion-safe recovery
- no full-table rerender per event

STOP after tests/build.

## Phase 6 — context and namespace

Implement context/namespace switching without losing the current session unexpectedly.

Support:

- context switch
- namespace switch
- all namespaces
- current context/ns visible globally
- production detection
- readonly state

Commands: `:ctx`, `:ns`.

STOP after tests/build.

## Phase 7 — pod workflow

Required pod shortcuts:

- `l`: logs
- `p`: previous logs
- `s`: shell/exec
- `a`: attach where supported
- `d`: describe
- `e`: edit
- `y`: YAML
- `f`: port-forward
- `c`: container selector
- Ctrl-D: delete

Logs require follow/pause, previous, timestamps, search, container selection and cancellation. Exec/attach must support terminal resize. Port-forward must have lifecycle/status in the status bar.

All mutations use ActionManager.

STOP after tests/build.

## Phase 8 — broader resources

Add deployments, replicasets, statefulsets, daemonsets, services, ingresses, nodes, events, jobs, cronjobs, configmaps, secrets metadata, PVCs and common workload resources.

Resource views share the same registry/table/inspector/command architecture. Secret values are never displayed.

Add deployment rollout/history/restart where supported by the backend.

STOP after tests/build.

## Phase 9 — safety and multi-select

Implement:

- multi-select
- bulk actions
- confirmation overlays
- production warnings
- readonly mode
- audit-friendly action history
- disabled-action reasons

Production destructive actions must identify cluster, namespace, resource and action before confirmation.

STOP after tests/build.

## Phase 10 — help and command palette

Build contextual help from a single shortcut registry. Ctrl-P searches actions and commands. `?` shows shortcuts applicable to the current resource. Avoid duplicated hardcoded help text.

STOP after tests/build.

## Phase 11 — AI enhancement

AI opens as an overlay from `a`/Ctrl-A. It is not a permanent panel.

Context gathering may include:

- selected resource YAML/metadata
- conditions
- current and previous logs
- events
- owner chain
- container states
- node conditions
- referenced object names

Redact secrets/tokens.

Response format:

1. Root Cause
2. Evidence
3. Recommended Fix
4. Commands

Commands are suggestions only. Never execute automatically.

STOP after tests/build.

## Phase 12 — Visual Mode and advanced views

Keep Visual Mode separate. It can provide topology, charts, dependency graphs and performance dashboards without changing Fast Mode's terminal-first contract.

STOP after tests/build.

## Phase 13 — tests and documentation

Add unit/integration/UI tests for:

- resource registry
- command registry
- navigation stack
- filters
- selection
- keyboard dispatch
- action safety
- readonly enforcement
- reconnect/watch updates
- AI redaction
- responsive 100x25 and 120x30 layouts

Update architecture docs and operator documentation.

## Definition of done

Fast Mode is complete when an operator can launch KubeMv and perform the complete lifecycle without leaving the keyboard-first console:

```text
Launch
→ choose/verify context and namespace
→ inspect Pods
→ filter
→ drill into Pod
→ logs / previous logs
→ events
→ describe / YAML
→ exec
→ port-forward
→ inspect owner chain
→ rollout/history
→ switch resource
→ change namespace/context
→ perform safe mutation
→ verify live watch result
→ use contextual AI diagnosis
```

## Required validation commands

```bash
go test ./...
cd apps/web && npm test && npm run build
```

Fix failures before reporting a phase.

## Phase execution protocol

At the end of each phase Cursor must report:

- files changed
- reused components
- deleted/replaced components
- shortcuts added
- tests/build output
- remaining parity gaps
- next phase recommendation

Then STOP.

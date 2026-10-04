# Changelog

## [0.7.0-phase6] - 2026-10-04

### Added

- Security summary API with RBAC explorer data and policy hints
- Security posture UI with grade, findings, and binding risk levels
- Live heuristics for privileged bindings, NetworkPolicy gaps, and runAsNonRoot

## [0.6.0-phase5] - 2026-10-04

### Added

- Observability dashboard, metrics query, and scrape target APIs
- Optional Prometheus HTTP proxy via `KUBEMV_PROMETHEUS_URL`
- Observability UI with metric charts and unified links to logs / external tools

## [0.5.0-phase4] - 2026-10-04

### Added

- Topology/XRay API and stub dependency graph for payment-api workload
- React Flow resource graph with health overlay and workload/network/dependency modes
- Topology navigation from Visual Mode, Topology page, and Fast Mode resource details

## [0.4.0-phase3] - 2026-10-04

### Added

- Advanced workload resource types (StatefulSet, DaemonSet, Job, CronJob, HPA, PDB)
- Deployment rollout status and workload mutations (scale, restart, rollback) with audit ids
- Fast Mode workload operations UI with confirmation dialogs

## [0.3.0-phase2] - 2026-10-04

### Added

- Pod logs API (current/previous, multi-container, search)
- WebSocket exec terminal (stub + live SPDY)
- Port-forward session API and Fast Mode pod operations UI (logs, exec, xterm)

## [0.2.0-phase1] - 2026-10-04

### Added

- `client-go` Kubernetes connector with kubeconfig context discovery
- Resource list/detail REST APIs for core K9s resource types
- Fast Mode explorer with split pane, command palette, and filters
- Namespace selector and command alias parsing (`:pods`, `:deploy`, etc.)

## [0.1.0-phase0] - 2026-10-04

### Added

- Monorepo layout (`apps/web`, `cmd/server`, `internal/*`, `docs/`)
- Go API server with health, meta, dev auth, and stub cluster overview endpoints
- Kubernetes and authentication abstraction layers
- React + Vite frontend shell with dark operations UI
- Design system primitives and Fast / Visual mode routing
- Phase 0 architecture and validation documentation

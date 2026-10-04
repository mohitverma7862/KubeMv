# Changelog

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

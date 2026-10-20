# KubeMv Architecture (Phase 0)

> Zero deviation from the Fast Mode reference: https://github.com/mohitverma7862/KubeMv/blob/main/docs/FAST-MODE-TUI-REFERENCE-SPEC.md and https://github.com/mohitverma7862/KubeMv/blob/main/docs/assets/kubemv-fast-mode-reference.svg. Match that specification and that image exactly: same logo, same colors, and same style.


KubeMv Fast Mode is the operator console. The API stays behind it. There is no dashboard UI.

## High-level topology

```text
Fast Mode console
        │
        │  REST / JSON (no kube credentials in browser)
        ▼
KubeMv API (Go)
        │
        ├── Auth abstraction (dev / future OIDC-SAML)
        ├── Credential manager (interface only in Phase 0)
        └── Kubernetes connector (stub → client-go in Phase 1)
                │
                ▼
           Kubernetes API
```

## Repository layout

```text
apps/
  web/          Fast Mode operator console
  desktop/      Tauri shell placeholder (future)
cmd/
  server/       HTTP API entrypoint
  kubemv/       CLI utilities
internal/
  api/          HTTP routes and handlers
  auth/         Application auth abstractions
  config/       Server configuration
  httputil/     JSON response helpers
  kubernetes/   Cluster client interfaces + stub connector
  server/       HTTP server wiring
docs/           Architecture and phase documentation
```

## Security rules (non-negotiable)

- Kubernetes credentials never reach the frontend.
- Application RBAC never bypasses Kubernetes RBAC.
- Secret values are not exposed by default.
- Dangerous mutations will require confirmation, policy, and audit in later phases.

## Phase roadmap

See `docs/PHASE-0.md` for validation criteria. Phase 1 implements K9s core resource navigation against real clusters via `client-go`.

# Phase 0 — Product Foundation

## Goals

- Monorepo structure for web, API, and future desktop shell
- Go API shell with health, meta, auth, and cluster overview endpoints
- Authentication abstraction (dev authenticator only)
- Kubernetes abstraction (`ClusterClient`, `Connector`, `CredentialManager` interface)
- React shell with dark operations UI, routing, TanStack Query, Zustand
- Design system primitives (`KPButton`, `KPStatusBadge`, `KPHealthScore`, etc.)
- Fast / Visual mode toggle with shared context (cluster, namespace)
- Unit tests (Go + Vitest) and documentation

## Explicitly out of scope

- AI providers, automation, cost management, marketplace
- Live Kubernetes watches, logs, exec, port-forward
- Production OIDC/SAML (interfaces only)

## Validation checklist

| Check | Command |
| --- | --- |
| Go unit tests | `go test ./...` |
| Go builds | `go build -o bin/kubemv-server ./cmd/server` |
| Frontend tests | `cd apps/web && npm test` |
| Frontend build | `cd apps/web && npm run build` |
| API health | `curl -s localhost:8080/healthz` |
| Login + overview | Sign in via UI, view cluster overview |

Phase 1 may begin only after this checklist passes.

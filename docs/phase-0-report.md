# Phase 0 completion report

Phase 0 is the architecture and operator-shell foundation. Kubernetes operations are not implemented. Do not start Phase 1 from this report; that waits for review.

## What shipped

- Go API on `127.0.0.1:8787` with health, session login, cluster metadata registry, and plugin manifest listing.
- OpenAPI contract in `api/openapi.yaml`.
- React operator UI: sign-in, overview, cluster registry, plugin list, settings, command palette, and keyboard navigation.
- Tauri 2 desktop shell with `core:default` only. The webview does not get a filesystem or shell permission.
- In-memory audit recorder that redacts credential-shaped details and is not exposed over HTTP.
- Closed role catalog. The bootstrap operator is `PLATFORM_ADMIN`. Routes check a session, not the later role matrix, and the process does not call Kubernetes.

## Checks

| Check | Result |
| --- | --- |
| `go test ./...` | Pass |
| `go vet ./...` and `gofmt -l` | Pass |
| `go build ./cmd/kubemv` | Pass |
| `npm run lint` | Pass |
| `npm run typecheck` | Pass |
| `npm test` (Vitest, 6 tests) | Pass |
| `npm run build` | Pass |
| `npm run e2e` (Playwright, 3 tests) | Pass: bad password, registry plus keyboard plus sign-out, 390px layout |
| `scripts/security-check.sh` | Pass |
| `cargo check` in `desktop/src-tauri` | Pass with Rust 1.99 and Tauri 2.12 |

Playwright covered sign-in failure, session redirect, command palette, cluster registration without echoing the server-side reference, removal, the empty plugin list, shortcut help, theme persistence, and sign-out. The narrow viewport shows every primary destination and does not widen the page.

Kind, k3d, and minikube are not used. This phase has no Kubernetes client, so a cluster harness would not exercise a code path.

## Security notes

- Bootstrap password is environment-only, bcrypt-hashed in memory, and never logged.
- Session id is an HttpOnly cookie. Mutations require a CSRF header and JSON.
- Cluster responses omit `kubeconfigRef`. Inline kubeconfigs and `token:` values are rejected.
- Unknown login fields are rejected. CORS does not reflect arbitrary origins.
- Plugin registration rejects permissions outside the catalog, including secret reads.
- `npm audit` reports moderate issues in transitive frontend tooling. They are not given a runtime path to cluster credentials. They were not force-upgraded.

## UI review

The shell is a dark graphite operator layout with a brass accent, IBM Plex type, and a light theme. Overview shows live API, session, cluster count, and plugin count. It does not invent workload health. Cluster rows say "Not connected". The command palette only runs navigation, theme, and sign-out.

## Left for later phases

Watches, logs, exec, port-forward, resource coverage, observability, RBAC evaluation, GitOps, and AI. The registry, authenticator, plugin host, and audit recorder are the extension points.

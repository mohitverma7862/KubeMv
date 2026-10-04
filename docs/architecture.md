# Architecture

KubeMv Phase 0 is two processes plus a desktop host.

```text
React operator UI  -- cookie session + CSRF -->  Go API
        ^
        |
   Tauri webview
```

The Tauri window has the core default permission only. It does not get a shell, filesystem, or HTTP plugin. The page calls the API with `fetch`.

## API

`api/openapi.yaml` is the contract. Routes:

| Method | Path | Auth |
| --- | --- | --- |
| GET | `/api/v1/health` | Public |
| POST | `/api/v1/auth/login` | Public |
| POST | `/api/v1/auth/logout` | Session + CSRF |
| GET | `/api/v1/auth/session` | Session |
| GET | `/api/v1/clusters` | Session |
| POST | `/api/v1/clusters` | Session + CSRF |
| GET | `/api/v1/clusters/{id}` | Session |
| DELETE | `/api/v1/clusters/{id}` | Session + CSRF |
| GET | `/api/v1/plugins` | Session |

## Identity

`auth.Authenticator` verifies an operator. The Phase 0 implementation is `LocalAuthenticator`: one bootstrap user from the environment, password retained only as a bcrypt hash in memory.

`auth.SessionStore` issues an opaque session id and a separate CSRF token. The id is stored in the `kubemv_session` HttpOnly cookie. The CSRF token is returned to the page and must be sent as `X-KubeMv-CSRF` on mutating requests.

The role catalog (`SUPER_ADMIN` through `AUDITOR`) is closed. The bootstrap operator is `PLATFORM_ADMIN`. Phase 0 authorizes by session presence. It does not evaluate the role matrix, and it does not call Kubernetes, so Kubernetes RBAC is not on this path yet.

## Clusters

`cluster.Registry` stores metadata: name, provider (`generic`, `eks`, `gke`, `aks`), context, and an optional server-side reference. Every record is `not_connected`. The memory registry does not dial an endpoint. `KubeconfigRef` is excluded from JSON.

## Plugins

`plugin.Plugin` exposes a manifest. `plugin.Registry.Register` rejects unknown categories, unknown permissions, and duplicate ids. There is no HTTP install endpoint. Secret reads and command execution are not in the permission catalog.

## Audit

`audit.Recorder` keeps a bounded in-process log of sign-in, sign-out, and registry changes. Details that look like credentials are replaced with `[redacted]`. Phase 0 does not expose this log over HTTP.

## UI state

TanStack Query owns server data: session, health, clusters, and plugins. Zustand owns theme and overlay state. The theme name is the only value written to `localStorage`.

## What is intentionally absent

No `client-go` client, no watches, no logs, no exec, no port-forward, no AI provider, and no cloud SDK. Those belong to later phases and should plug into the registry, plugin, and audit interfaces already here.

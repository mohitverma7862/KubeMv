# Development

## Prerequisites

- Go 1.22
- Node.js 22
- Rust 1.90 or newer, plus WebKitGTK 4.1 headers, to compile the Tauri shell
- A Chromium build for Playwright when running `make e2e`

## Environment

| Variable | Default | Purpose |
| --- | --- | --- |
| `KUBEMV_ADDR` | `127.0.0.1:8787` | API listen address |
| `KUBEMV_BOOTSTRAP_USERNAME` | required | Local operator name |
| `KUBEMV_BOOTSTRAP_PASSWORD` | required | Local operator password, 12–128 characters |
| `KUBEMV_ALLOWED_ORIGINS` | `http://127.0.0.1:1420,http://localhost:1420` | Browser origins allowed to call the API |
| `KUBEMV_SESSION_TTL` | `8h` | Session lifetime, from `1m` to `24h` |
| `KUBEMV_COOKIE_SECURE` | `false` | Mark the session cookie Secure |
| `VITE_API_BASE` | `http://127.0.0.1:8787` | API origin used by the UI |

The cluster registry and session store are in memory. Restarting the API drops them.

## Commands

```bash
make test
make lint
make typecheck
make security
make build
make e2e
```

Backend tests use `httptest` and do not need a cluster. Kind, k3d, and minikube are not part of Phase 0 because the API never contacts Kubernetes. Playwright covers sign-in, the command palette, cluster registration, theme persistence, and a narrow viewport.

## Desktop icons

```bash
python3 scripts/gen_icons.py
```

## Adding a later phase

Keep new Kubernetes calls behind a client that the registry does not grow by accident. Do not return `KubeconfigRef`. Do not add a plugin upload route. Extend `api/openapi.yaml` in the same change as the handler.

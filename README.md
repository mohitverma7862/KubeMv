# KubeMv

The intelligent command center for Kubernetes.

Phase 0 is the architecture and operator-shell foundation. It signs an operator in, stores cluster registry metadata, and enforces a plugin permission boundary. It does not call the Kubernetes API.

## Layout

```text
api/openapi.yaml          HTTP contract
backend/                  Go API
desktop/                  React operator UI
desktop/src-tauri/        Tauri desktop shell
docs/                     Architecture, security, and phase notes
```

The desktop shell renders the React UI. The UI talks to the local Go API. Kubernetes credentials, when a later phase accepts a server-side reference, stay in the API process. The browser and the Tauri webview do not receive kubeconfigs or secret values.

## Run

Use a bootstrap password of at least 12 characters. Do not commit it.

```bash
export KUBEMV_BOOTSTRAP_USERNAME=admin
export KUBEMV_BOOTSTRAP_PASSWORD='replace-with-a-long-password'
cd backend && go run ./cmd/kubemv
```

The API listens on `127.0.0.1:8787`.

```bash
cd desktop && npm install && npm run dev
```

Open `http://127.0.0.1:1420`.

The Tauri shell loads that same UI:

```bash
cd desktop && npm run tauri dev
```

## Phase 0 boundaries

Signed-in operators can:

- read process health and the capability list
- register, list, and remove cluster metadata
- list plugin manifests compiled into the process

They cannot:

- watch or mutate Kubernetes resources
- read secret values
- upload or execute plugins
- retrieve a kubeconfig reference after it is stored

Kubernetes watches, logs, exec, and the rest of the product phases are not in this build.

## Checks

```bash
make check
make e2e
```

See [docs/development.md](docs/development.md), [docs/architecture.md](docs/architecture.md), and [docs/security.md](docs/security.md).

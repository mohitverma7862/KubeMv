# KubeMv

Fast Mode is the Kubernetes operator console. Keyboard-first, dense, and modeled on the K9s workflow. There is no dashboard.

**Version:** `0.9.0-phase8`

## Operator console

The only screen is Fast Mode:

- header with cluster, context, namespace, and connection
- resource navigation
- resource table
- inspector and relationship tree
- logs, events, YAML, and describe
- command and status bar

Keys: `j`/`k`, `g`/`G`, Enter, Esc, `/`, `:`, `?`, Ctrl-P, Ctrl-R.

## Repository

See [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) and [docs/FAST-MODE-UI-DESIGN.md](docs/FAST-MODE-UI-DESIGN.md).

## Development

### Prerequisites

- Go 1.22+
- Node.js 20+
- `kubectl` (installed by `.cursor/install.sh` in Cloud Agents)

### Install

```bash
bash .cursor/install.sh
cd apps/web && npm ci
```

### Run the console

```bash
# Terminal 1 — API
go run ./cmd/server

# Terminal 2 — Fast Mode console
cd apps/web && npm run dev
```

Open `http://localhost:5173`. Sign in with any username. The session lands on the Fast Mode console.

### Test and build

```bash
go test ./...
cd apps/web && npm test && npm run lint && npm run build
```

## CLI

```bash
go build -o bin/kubemv ./cmd/kubemv
./bin/kubemv version
```

`kubemv` prints the version. It does not open a second interface.

## License

See [LICENSE](LICENSE).

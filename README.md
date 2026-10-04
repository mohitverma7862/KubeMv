# KubeMv

**Advanced Kubernetes Operations Platform** — K9s speed with a modern GUI, security, GitOps, and optional AI assistance.

**Status:** Phase 0 (foundation)  
**Version:** `0.1.0-phase0`

## Product modes

- **Fast Mode** — keyboard-first operator workflows (K9s-inspired)
- **Visual Mode** — dashboards, topology, and dense operational views

## Repository structure

See [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) and [docs/PHASE-0.md](docs/PHASE-0.md).

## Development

### Prerequisites

- Go 1.22+
- Node.js 20+ (for `apps/web`)
- `kubectl` (installed by `.cursor/install.sh` in Cloud Agents)

### Install

```bash
bash .cursor/install.sh
cd apps/web && npm ci
```

### Run API + UI

```bash
# Terminal 1
go run ./cmd/server

# Terminal 2
cd apps/web && npm run dev
```

Open `http://localhost:5173` and sign in with any username (dev auth).

### Test & build

```bash
go test ./...
cd apps/web && npm test && npm run build
```

## CLI

```bash
go build -o bin/kubemv ./cmd/kubemv
./bin/kubemv version
```

## License

See [LICENSE](LICENSE).

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

Keys: `j`/`k`, `g`/`G`, `1`-`9`, `/`, `?`, `q`.

## Repository

See [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) and [docs/FAST-MODE-UI-DESIGN.md](docs/FAST-MODE-UI-DESIGN.md).

## Start

From a terminal, the same way you start k9s:

```bash
go run ./cmd/kubemv
```

Or install it and run it by name:

```bash
go build -o bin/kubemv ./cmd/kubemv
./bin/kubemv
```

The console opens in that terminal. `q` quits. With a kubeconfig, it uses the current context. Without one, it opens the sample console and says why. A production context shows PRODUCTION. A nonprod or staging context does not.

`kubemv version` prints the version and does not open the console.

## Development

### Prerequisites

- Go 1.22+
- A kubeconfig, if you want the current cluster instead of the sample console

### Install

```bash
bash .cursor/install.sh
```

### Test and build

```bash
go test ./...
go build -o bin/kubemv ./cmd/kubemv
```

## License

See [LICENSE](LICENSE).

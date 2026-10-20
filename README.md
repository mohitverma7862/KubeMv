# KubeMv

> Zero deviation from the Fast Mode reference: https://github.com/mohitverma7862/KubeMv/blob/main/docs/FAST-MODE-TUI-REFERENCE-SPEC.md and https://github.com/mohitverma7862/KubeMv/blob/main/docs/assets/kubemv-fast-mode-reference.svg. Match that specification and that image exactly: same logo, same colors, and same style.


Fast Mode is the Kubernetes operator console. Start it in a terminal the same way you start k9s.

## Start

```bash
go run ./cmd/kubemv
```

Build it once, then run it by name:

```bash
go build -o bin/kubemv ./cmd/kubemv
./bin/kubemv
```

`:q` quits. The console uses the current kubeconfig context. Without a kubeconfig it opens the sample console and the status line says why.

A production context shows PRODUCTION. A nonprod or staging context does not.

```bash
./bin/kubemv version
```

`version` prints the version and does not open the console.

## Keys

| Key | Action |
| --- | --- |
| `j` / `k` | Select row |
| `g` / `G` | First / last |
| Enter | Drill down |
| Esc | Back |
| Ctrl-R | Refresh |
| `/` | Filter |
| `:` | Command mode (`:pods`, `:q`) |
| `?` | Help |
| Ctrl-P | Command palette |

## Requirements

- Go 1.22+
- A kubeconfig, when you want the current cluster instead of the sample console

## Test and build

```bash
go test ./...
go build -o bin/kubemv ./cmd/kubemv
```

## Docs

- [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md)
- [docs/FAST-MODE-UI-DESIGN.md](docs/FAST-MODE-UI-DESIGN.md)

## License

See [LICENSE](LICENSE).

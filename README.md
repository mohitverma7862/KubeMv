# KubeMv

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

`q` quits. The console uses the current kubeconfig context. Without a kubeconfig it opens the sample console and the status line says why.

A production context shows PRODUCTION. A nonprod or staging context does not.

```bash
./bin/kubemv version
```

`version` prints the version and does not open the console.

## Keys

| Key | Action |
| --- | --- |
| `j` / `k` | Next / previous row |
| `g` / `G` | First / last row |
| `1`–`9` | Switch resource |
| `/` | Filter |
| `Esc` | Clear filter or close help |
| `?` | Help |
| `q` | Quit |

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

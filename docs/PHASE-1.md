# Phase 1 — K9s Core

> Zero deviation from the Fast Mode reference: https://github.com/mohitverma7862/KubeMv/blob/main/docs/FAST-MODE-TUI-REFERENCE-SPEC.md and https://github.com/mohitverma7862/KubeMv/blob/main/docs/assets/kubemv-fast-mode-reference.svg. Match that specification and that image exactly: same logo, same colors, and same style.


## Delivered

- Live Kubernetes connector via kubeconfig contexts (`internal/kubernetes/kube`)
- Stub connector fallback for development without a cluster
- Resource list + detail APIs for pods, deployments, services, nodes, events, configmaps, secrets (metadata only), namespaces, CRDs
- Fast Mode split view: resource table, details, YAML/events
- Command palette (Ctrl+K) with K9s-style resource aliases
- Namespace selector and `/` / `:` filter bar parsing
- Server-side label and name filtering hooks

## Validation

```bash
go test ./...
cd apps/web && npm test && npm run build
go run ./cmd/server
# login, then:
curl -H "Authorization: Bearer $TOKEN" "http://localhost:8080/api/v1/clusters/local/resources/pods?namespace=payments"
```

## Next (Phase 2)

Logs, exec, port-forward, multi-container support.

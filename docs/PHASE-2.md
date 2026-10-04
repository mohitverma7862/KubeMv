# Phase 2 — Logs / Exec / Port Forward

## Delivered

- Pod log API with container selection, previous logs, tail lines, and search (`q`)
- Multi-container discovery endpoint
- WebSocket exec (`/api/v1/ws/clusters/{id}/exec`) with stub shell and live `client-go` SPDY when kubeconfig is available
- Server-managed port-forward session registry (start/list/stop)
- Fast Mode **Pod operations** panel: logs (auto-refresh), exec terminal (xterm.js), port-forward manager

## Validation

```bash
go test ./...
cd apps/web && npm test && npm run build
curl -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8080/api/v1/clusters/local/pods/payments/payment-api-7d9c8/logs?container=app&q=OOM"
```

## Next (Phase 3)

Rollouts, rollback, scaling, HPA, PDB, advanced workload actions.

# Phase 5 — Observability

## Delivered

- Observability API:
  - `GET /api/v1/clusters/{id}/observability/dashboard`
  - `GET /api/v1/clusters/{id}/observability/metrics` (Prometheus proxy when `KUBEMV_PROMETHEUS_URL` is set; stub series otherwise)
  - `GET /api/v1/clusters/{id}/observability/targets`
- Stub dashboards for `payment-api` (CPU, memory, RPS, error rate) with scrape target health
- Live scrape target discovery from Services with `prometheus.io/*` annotations
- Observability UI with metric sparklines, target table, Grafana/Loki links, Fast Mode log deep link

## Configuration

| Variable | Purpose |
|----------|---------|
| `KUBEMV_PROMETHEUS_URL` | Base URL for Prometheus HTTP API (optional) |
| `KUBEMV_GRAFANA_URL` | Reserved for future deep links |
| `KUBEMV_LOKI_URL` | Reserved for future log explorer integration |

## Validation

```bash
go test ./...
cd apps/web && npm test && npm run build
TOKEN=... # from POST /api/v1/auth/login
curl -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8080/api/v1/clusters/local/observability/dashboard?namespace=payments&name=payment-api"
```

## Next (Phase 6)

See [PHASE-6.md](PHASE-6.md).

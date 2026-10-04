# Phase 4 — XRay & Topology

## Delivered

- Topology API: `GET /api/v1/clusters/{id}/topology` with `mode=workload|network|dependency`
- Stub graph: Ingress → Service → Deployment → ReplicaSet → Pods + ConfigMap/Secret/ServiceAccount
- Live graph: deployments, pods (owner refs), services, ingresses (network mode)
- React Flow graph UI with pan/zoom/minimap and health-colored node borders
- Visual Mode + dedicated Topology route; Fast Mode link to open selected resource in graph

## Validation

```bash
go test ./...
cd apps/web && npm test && npm run build
curl -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8080/api/v1/clusters/local/topology?mode=workload&namespace=payments"
```

## Next (Phase 5)

Metrics, Prometheus integration, unified observability views.

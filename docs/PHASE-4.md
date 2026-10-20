# Phase 4 — XRay & Topology

> Zero deviation from the Fast Mode reference: https://github.com/mohitverma7862/KubeMv/blob/main/docs/FAST-MODE-TUI-REFERENCE-SPEC.md and https://github.com/mohitverma7862/KubeMv/blob/main/docs/assets/kubemv-fast-mode-reference.svg. Match that specification and that image exactly: same logo, same colors, and same style.


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

See [PHASE-5.md](PHASE-5.md) — delivered on `cursor/phase-5-observability-8762`.

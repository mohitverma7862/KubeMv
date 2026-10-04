# Phase 7 — GitOps

## Delivered

- GitOps API: `GET /api/v1/clusters/{id}/gitops/overview?namespace=`
- Stub: Argo CD / Flux style apps for `payments`, manifest drift (replicas, ConfigMap), CI pipeline runs
- Live: GitOps-labeled Deployments (Argo/Flux annotations or labels), rollout generation / replica drift heuristics
- GitOps UI at `/gitops` with applications, drift, and pipeline hooks tables

## Validation

```bash
go test ./...
cd apps/web && npm test && npm run build
curl -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8080/api/v1/clusters/local/gitops/overview?namespace=payments"
```

## Next (Phase 8)

See [PHASE-8.md](PHASE-8.md).

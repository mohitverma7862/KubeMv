# Phase 6 — Security & RBAC

## Delivered

- Security API: `GET /api/v1/clusters/{id}/security/summary?namespace=`
- Stub posture for `payments`: score/grade, RBAC bindings table, policy findings (RBAC, workload, network, secrets)
- Live cluster: RoleBindings + relevant ClusterRoleBindings, heuristic findings (privileged roles, missing NetworkPolicy, runAsNonRoot)
- Security UI replacing the placeholder route at `/security`

## Validation

```bash
go test ./...
cd apps/web && npm test && npm run build
curl -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8080/api/v1/clusters/local/security/summary?namespace=payments"
```

## Next (Phase 7)

GitOps views, manifest drift, and deployment pipeline hooks (per product roadmap).

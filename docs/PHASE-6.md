# Phase 6 — Security & RBAC

> Zero deviation from the Fast Mode reference: https://github.com/mohitverma7862/KubeMv/blob/main/docs/FAST-MODE-TUI-REFERENCE-SPEC.md and https://github.com/mohitverma7862/KubeMv/blob/main/docs/assets/kubemv-fast-mode-reference.svg. Match that specification and that image exactly: same logo, same colors, and same style.


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

See [PHASE-7.md](PHASE-7.md).

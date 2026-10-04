# Phase 8 — AI Assist

## Delivered

- Assist API:
  - `GET /api/v1/clusters/{id}/assist/bundle?namespace=&kind=&name=`
  - `POST /api/v1/clusters/{id}/assist/hooks/{hookId}/dry-run`
- Rule-based triage in `internal/ai` (on-cluster signals; no external LLM in this phase)
- Stub CrashLoop scenario for `payment-api` with hypotheses, kubectl runbook steps, automation hooks
- Live: collects pod/deployment status, events, log tail for triage input
- UI at `/assist`; `aiEnabled` in meta when `KUBEMV_AI_ENABLED=true` (default)

## Validation

```bash
go test ./...
cd apps/web && npm test && npm run build
curl -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8080/api/v1/clusters/local/assist/bundle?namespace=payments&kind=Deployment&name=payment-api"
```

## Next (Phase 9)

Hardening: multi-cluster RBAC, audit export, and production auth (OIDC) integration.

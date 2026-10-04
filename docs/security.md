# Security

Phase 0 enforces the parts of the trust model that exist before any Kubernetes call.

## Session

- Bootstrap password comes from `KUBEMV_BOOTSTRAP_PASSWORD` and must be 12 to 128 characters.
- The password is hashed with bcrypt at startup and is not logged.
- Failed sign-ins share one error message. The response and audit detail do not echo the password.
- Five failures for a username inside a minute return `429`.
- The session cookie is HttpOnly, `SameSite=Lax`, and scoped to `/`. Set `KUBEMV_COOKIE_SECURE=true` when the API is served over HTTPS.
- Mutating routes require `Content-Type: application/json` and a matching CSRF token.
- CORS allows only the configured origins and does not reflect arbitrary origins.
- Login bodies reject unknown JSON fields.

## Credentials

- The API binds to `127.0.0.1:8787` unless `KUBEMV_ADDR` is set.
- Cluster registration rejects inline kubeconfig documents, private keys, `token:` values, and path traversal.
- A stored reference uses `json:"-"` and is absent from list, get, and create responses.
- The UI renders a fixed set of cluster fields and does not print unknown payload keys.
- AI is not called, so no cluster data leaves the process for a model.

## Plugins and desktop

- Plugin permissions are a closed catalog. `secret.read` and command execution are rejected.
- Plugins cannot be posted to the API.
- The Tauri capability grants `core:default` only.

## Audit

Sign-in, sign-out, CSRF denial, and cluster register/remove are recorded in memory. The recorder redacts details that contain credential markers. The log is not an HTTP resource in this phase.

## Checks

`scripts/security-check.sh` runs the security-focused Go tests and scans the UI sources for HTML injection, credential literals, and `localStorage` use outside the theme key.

#!/usr/bin/env bash
set -euo pipefail

if ! command -v kubectl >/dev/null 2>&1; then
  KUBECTL_VERSION="$(curl -fsSL https://dl.k8s.io/release/stable.txt)"
  curl -fsSL "https://dl.k8s.io/release/${KUBECTL_VERSION}/bin/linux/amd64/kubectl" -o /tmp/kubectl
  chmod +x /tmp/kubectl
  sudo install -m 0755 /tmp/kubectl /usr/local/bin/kubectl
fi

if [[ -f go.mod ]]; then
  go mod download
  go build -o /dev/null ./...
fi

if [[ -f apps/web/package.json ]]; then
  (cd apps/web && npm ci)
  (cd apps/web && npm run build)
fi

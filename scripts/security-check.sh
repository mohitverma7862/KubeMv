#!/usr/bin/env bash
set -euo pipefail
root="$(cd "$(dirname "$0")/.." && pwd)"
cd "$root"

echo "backend security tests"
(
  cd backend
  go test ./... -count=1 -run 'Security|Kubeconfig|Password|Permission|CORS|Redact|RejectsSecrets|RoleCatalog'
)

echo "frontend secret handling"
python3 - <<'PY'
import pathlib, sys
root = pathlib.Path("desktop")
bad = []
for path in root.rglob("*"):
    if not path.is_file() or "node_modules" in path.parts or "dist" in path.parts or "src-tauri/target" in str(path):
        continue
    if path.suffix not in {".ts", ".tsx", ".js", ".html", ".css"}:
        continue
    text = path.read_text(errors="ignore")
    if "dangerouslySetInnerHTML" in text:
        bad.append(f"{path}: dangerouslySetInnerHTML")
    if path.parts[:2] == ("desktop", "src") or path.name == "index.html":
        if "localStorage" in text and "kubemv.theme" not in text:
            bad.append(f"{path}: localStorage outside theme")
        if "foundation-test-password" in text or "BEGIN PRIVATE" in text:
            bad.append(f"{path}: embedded credential")
if bad:
    print("\n".join(bad))
    sys.exit(1)
print("frontend patterns ok")
PY

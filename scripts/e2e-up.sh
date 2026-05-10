#!/usr/bin/env bash
# scripts/e2e-up.sh
# POSIX twin of e2e-up.ps1. Bootstraps the dashboard for e2e runs without
# touching the user's progress.json.
#
# Usage:
#   ./scripts/e2e-up.sh           # foreground (blocks)
#   ./scripts/e2e-up.sh --detach  # background, prints PID

set -euo pipefail

DETACH=false
[[ "${1:-}" == "--detach" ]] && DETACH=true

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
PROGRESS_FILE="$REPO_ROOT/e2e/.tmp-progress.json"
SERVER_SCRIPT="$REPO_ROOT/dashboard/server.js"

mkdir -p "$REPO_ROOT/e2e"
echo '{"version":3,"problems":{}}' > "$PROGRESS_FILE"
echo "[e2e-up] Reset $PROGRESS_FILE"

export PROGRESS_JSON_PATH=e2e/.tmp-progress.json

if ! $DETACH; then
  echo "[e2e-up] Starting server in foreground (Ctrl+C to stop)..."
  cd "$REPO_ROOT"
  exec node "$SERVER_SCRIPT"
fi

cd "$REPO_ROOT"
node "$SERVER_SCRIPT" >/dev/null 2>&1 &
PID=$!
echo "[e2e-up] Server started (PID $PID) — detached"

# Poll /api/runner-status until ready (max 10s)
for _ in $(seq 1 40); do
  if status=$(curl -fsS http://localhost:3000/api/runner-status 2>/dev/null); then
    echo "[e2e-up] Runner status: $status"
    echo "[e2e-up] Stop with: kill $PID"
    exit 0
  fi
  sleep 0.25
done

echo "[e2e-up] Server did not respond on :3000 within 10s. Killing PID $PID." >&2
kill "$PID" 2>/dev/null || true
exit 1

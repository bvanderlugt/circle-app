#!/usr/bin/env bash
# smoke test
set -euo pipefail

BASE_URL="${BASE_URL:-http://localhost:8080}"

wait_for_app() {
  echo "Waiting for app..."
  for i in $(seq 1 60); do
    curl -sf "$BASE_URL/" > /dev/null && echo "Ready after ${i}s" && return
    sleep 1
  done
  echo "App failed to start" && exit 1
}

wait_for_app

echo "--- frontend ---"
curl -sf "$BASE_URL/" | grep -q 'id="root"'
echo "HTML OK"

echo "--- api ---"
curl -sf "$BASE_URL/api/circles" | python3 -c "import sys,json; assert isinstance(json.load(sys.stdin), list)"
echo "GET /api/circles OK"

ID=$(curl -sf -X POST "$BASE_URL/api/circles" \
  -H "Content-Type: application/json" \
  -d '{"name":"smoke-test"}' | python3 -c "import sys,json; print(json.load(sys.stdin)['id'])")
echo "POST created id=$ID"

curl -sf -o /dev/null -w "%{http_code}" -X DELETE "$BASE_URL/api/circles/$ID" | grep -q 204
echo "DELETE OK"

echo "--- all checks passed ---"

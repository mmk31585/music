#!/usr/bin/env bash
set -euo pipefail

# Muse Watchdog v2 — Autonomous site reliability monitor
# Starts both servers, monitors health, captures errors, writes structured status.
# Run: bash scripts/watchdog.sh
# Stop: Ctrl+C

INTERVAL=10
FRONTEND_PORT=3000
API_PORT=8080
FRONTEND_PID=""
API_PID=""
LOG_FILE=".watchdog.log"
STATUS_FILE=".watchdog-status.json"
FRONTEND_LOG=".watchdog-frontend.log"
API_LOG=".watchdog-api.log"
FAILURES=0
MAX_FAILURES=5

init_status() {
  cat > "$STATUS_FILE" <<'EOF'
{
  "frontend": "unknown",
  "api": "unknown",
  "postgres": "unknown",
  "redis": "unknown",
  "last_errors": [],
  "restart_count": 0,
  "last_check": ""
}
EOF
}

update_status() {
  local key="$1" value="$2"
  local tmp=$(mktemp)
  python3 -c "
import json
with open('$STATUS_FILE') as f:
  s = json.load(f)
s['$key'] = '$value'
s['last_check'] = '$(date -Iseconds)'
with open('$STATUS_FILE', 'w') as f:
  json.dump(s, f, indent=2)
" 2>/dev/null || true
  rm -f "$tmp"
}

log_error() {
  local msg="$1"
  echo "[$(date '+%H:%M:%S')] ERROR: $msg" | tee -a "$LOG_FILE"
  local tmp=$(mktemp)
  python3 -c "
import json
with open('$STATUS_FILE') as f:
  s = json.load(f)
s['last_errors'].append({'time': '$(date -Iseconds)', 'msg': '$msg'})
s['last_errors'] = s['last_errors'][-10:]
with open('$STATUS_FILE', 'w') as f:
  json.dump(s, f, indent=2)
" 2>/dev/null || true
  rm -f "$tmp"
  FAILURES=$((FAILURES + 1))
}

log() {
  echo "[$(date '+%H:%M:%S')] $*" | tee -a "$LOG_FILE"
}

cleanup() {
  log "Shutting down watchdog..."
  [ -n "$FRONTEND_PID" ] && kill "$FRONTEND_PID" 2>/dev/null || true
  [ -n "$API_PID" ] && kill "$API_PID" 2>/dev/null || true
  update_status "frontend" "stopped"
  update_status "api" "stopped"
  exit 0
}
trap cleanup SIGINT SIGTERM

check_deps() {
  command -v node >/dev/null 2>&1 || { log "ERROR: node not found"; exit 1; }
  command -v go >/dev/null 2>&1 || { log "ERROR: go not found"; exit 1; }
  command -v curl >/dev/null 2>&1 || { log "ERROR: curl not found"; exit 1; }
}

start_frontend() {
  log "Starting frontend on :$FRONTEND_PORT..."
  cd frontend
  npm run dev > "../$FRONTEND_LOG" 2>&1 &
  FRONTEND_PID=$!
  cd ..
  update_status "frontend" "starting"
}

start_api() {
  log "Starting API on :$API_PORT..."
  go run cmd/api/main.go > "$API_LOG" 2>&1 &
  API_PID=$!
  update_status "api" "starting"
}

check_frontend() {
  local code
  code=$(curl -s -o /dev/null -w "%{http_code}" "http://localhost:$FRONTEND_PORT" 2>/dev/null || echo "000")
  if [ "$code" = "000" ] || [ "$code" = "502" ]; then
    log_error "FRONTEND DOWN (HTTP $code)"
    update_status "frontend" "down"
    # Capture last errors from log
    tail -20 "$FRONTEND_LOG" 2>/dev/null | while IFS= read -r line; do
      echo "$line" | grep -iE "error|fail|panic|crash|not found|conflict" >> "$LOG_FILE" 2>/dev/null || true
    done
    # Restart
    kill "$FRONTEND_PID" 2>/dev/null || true
    sleep 2
    start_frontend
    return 1
  fi
  log "Frontend OK (HTTP $code)"
  update_status "frontend" "healthy"
  return 0
}

check_api() {
  local code
  code=$(curl -s -o /dev/null -w "%{http_code}" "http://localhost:$API_PORT/api/v1/health" 2>/dev/null || echo "000")
  if [ "$code" = "000" ]; then
    log_error "API DOWN (HTTP $code)"
    update_status "api" "down"
    tail -20 "$API_LOG" 2>/dev/null | while IFS= read -r line; do
      echo "$line" | grep -iE "error|fail|panic|crash|not found|conflict|fatal" >> "$LOG_FILE" 2>/dev/null || true
    done
    kill "$API_PID" 2>/dev/null || true
    sleep 2
    start_api
    return 1
  fi
  log "API OK (HTTP $code)"
  update_status "api" "healthy"
  return 0
}

check_infra() {
  # Postgres runs in Docker on TCP port 5432 (not Unix socket)
  if pg_isready -h localhost -p 5432 -q 2>/dev/null; then
    update_status "postgres" "healthy"
  elif command -v pg_isready &>/dev/null && pg_isready -q 2>/dev/null; then
    update_status "postgres" "healthy"
  else
    update_status "postgres" "down"
    log_error "Postgres is not reachable"
  fi
  # Redis runs in Docker; redis-cli may not be on host, use nc as fallback
  if command -v redis-cli &>/dev/null; then
    if redis-cli ping 2>/dev/null | grep -q "PONG"; then
      update_status "redis" "healthy"
    else
      update_status "redis" "down"
      log_error "Redis is not reachable"
    fi
  elif echo "PING" | nc -w2 localhost 6379 2>/dev/null | grep -q "PONG"; then
    update_status "redis" "healthy"
  else
    update_status "redis" "down"
    log_error "Redis is not reachable"
  fi
}

# === MAIN ===
init_status
check_deps

# Clean slate
log "Cleaning up stale processes..."
lsof -ti ":$FRONTEND_PORT" 2>/dev/null | xargs kill -9 2>/dev/null || true
lsof -ti ":$API_PORT" 2>/dev/null | xargs kill -9 2>/dev/null || true
rm -f "$FRONTEND_LOG" "$API_LOG"
sleep 1

start_frontend
sleep 3
start_api

log "Watchdog v2 active — checking every ${INTERVAL}s (Ctrl+C to stop)"

while true; do
  sleep "$INTERVAL"
  check_infra || true
  check_frontend || true
  check_api || true

  # If too many failures in a row, escalate
  if [ "$FAILURES" -gt "$MAX_FAILURES" ]; then
    log "ESCALATING: $FAILURES consecutive failures — check .watchdog-status.json"
    update_status "escalated" "true"
    FAILURES=0
  fi
done

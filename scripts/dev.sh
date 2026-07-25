#!/usr/bin/env bash
# ═══════════════════════════════════════════════════════════════════════════
# Muse Dev — Unified development server manager
# ═══════════════════════════════════════════════════════════════════════════
# Starts/stops all services in parallel with consolidated log viewing.
#
# Services managed:
#   infra    → Postgres, Redis, MinIO, OpenSearch (Docker)
#   api      → Go/Gin backend (:8080)
#   web      → Vite frontend dev server (:3000)
#   ml       → FastAPI ML service (:8000)
#   worker   → Celery worker for ML (Whisper inference)
#   goworker → Go worker (cmd/worker)
#
# Usage:
#   bash scripts/dev.sh start        Start everything
#   bash scripts/dev.sh stop         Stop everything gracefully
#   bash scripts/dev.sh restart      Restart everything
#   bash scripts/dev.sh status       Show health of all services
#   bash scripts/dev.sh logs         Tail all logs (follow)
#   bash scripts/dev.sh logs api     Tail specific service log
#
# Log files (in PROJECT_ROOT/logs/):
#   api.log      Go backend
#   web.log      Vite dev server
#   ml-api.log   FastAPI ML service
#   worker.log   Celery worker
#   goworker.log Go worker
# ═══════════════════════════════════════════════════════════════════════════

set -uo pipefail
# Note: NO set -e. We handle errors explicitly with || checks.

# ── Paths ──────────────────────────────────────────────────────────────
PROJECT_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
LOG_DIR="$PROJECT_ROOT/logs"
PID_DIR="$PROJECT_ROOT/.pids"
ML_DIR="$PROJECT_ROOT/moja-ml-service"

# Services and their ports (for health checks)
declare -A SERVICE_PORT=(
  [web]="3000"
  [api]="8080"
  [ml]="8000"
)

# ── Colors ─────────────────────────────────────────────────────────────
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

info()  { echo -e "${CYAN}[dev]${NC} $*"; }
ok()    { echo -e "${GREEN}[✓]${NC} $*"; }
warn()  { echo -e "${YELLOW}[!]${NC} $*"; }
err()   { echo -e "${RED}[✗]${NC} $*" >&2; }

kill_port() {
  local port="$1"
  local pids
  pids=$(lsof -ti :"$port" 2>/dev/null) || true
  if [ -n "$pids" ]; then
    warn "Port $port in use by PID(s) $pids — killing..."
    echo "$pids" | xargs kill -9 2>/dev/null || true
    sleep 1
  fi
}

# ── Helpers ────────────────────────────────────────────────────────────

cleanup_logs() {
  mkdir -p "$LOG_DIR" "$PID_DIR"
}

is_running() {
  local name="$1"
  local pid_file="$PID_DIR/$name.pid"
  # Check PID file first
  if [ -f "$pid_file" ]; then
    local pid
    pid=$(cat "$pid_file")
    if kill -0 "$pid" 2>/dev/null; then
      return 0
    fi
  fi
  # Fallback for Docker-based workers
  if [ "$name" = "worker" ]; then
    local status
    status=$(docker inspect moja-ml-service-worker-1 --format='{{.State.Status}}' 2>/dev/null || echo "")
    [ "$status" = "running" ] && return 0
  fi
  # Fallback: check port-based health (handles stale PIDs after exec)
  health_check "$name" 2>/dev/null && return 0
  return 1
}

save_pid() {
  local name="$1"
  local pid="$2"
  mkdir -p "$PID_DIR"
  echo "$pid" > "$PID_DIR/$name.pid"
}

remove_pid() {
  rm -f "$PID_DIR/$1.pid"
}

health_check() {
  local name="$1"
  case "$name" in
    web)   curl -sfo /dev/null --max-time 2 "http://localhost:3000/" 2>/dev/null ;;
    api)   curl -sfo /dev/null --max-time 2 "http://localhost:8080/api/v1/health" 2>/dev/null ;;
    ml)    curl -sfo /dev/null --max-time 2 "http://localhost:8000/api/v1/health/live" 2>/dev/null ;;
    *)     return 1 ;;
  esac
}

# ── Service Start Functions ────────────────────────────────────────────

start_infra() {
  info "Starting Docker infrastructure..."
  (cd "$PROJECT_ROOT" && make infra-up 2>&1) || true
  ok "Infrastructure containers started"
}

_bg() {
  # Run a command in background with log redirection, save PID, disown.
  # Usage: _bg <pid_name> <log_file> <command...>
  local name="$1"
  local logfile="$2"
  shift 2
  (
    exec &>"$logfile"
    exec "$@"
  ) &
  local pid=$!
  save_pid "$name" "$pid"
  disown "$pid" 2>/dev/null || true
}

start_api() {
  if is_running "api"; then
    warn "API is already running"
    return 0
  fi
  kill_port 8080
  info "Starting Go API (:8080)..."
  export ML_SERVICE_WEBHOOK_HMAC_SECRET="${ML_SERVICE_WEBHOOK_HMAC_SECRET:-test-secret}"
  _bg "api" "$LOG_DIR/api.log" go run "$PROJECT_ROOT/cmd/api/main.go"
  sleep 6
  if curl -sfo /dev/null --max-time 3 "http://localhost:8080/api/v1/health" 2>/dev/null; then
    ok "API is running"
  else
    warn "API may still be starting up — check $LOG_DIR/api.log"
  fi
}

start_web() {
  if is_running "web"; then
    warn "Frontend is already running"
    return 0
  fi
  info "Starting Vite frontend (:3000)..."
  _bg "web" "$LOG_DIR/web.log" npm run dev --prefix "$PROJECT_ROOT/frontend" -- --host
  sleep 6
  if curl -sfo /dev/null --max-time 3 "http://localhost:3000/" 2>/dev/null; then
    ok "Frontend is running"
  else
    warn "Frontend may still be starting up — check $LOG_DIR/web.log"
  fi
}

start_ml() {
  if is_running "ml"; then
    warn "ML API is already running"
    return 0
  fi
  kill_port 8000
  info "Starting ML FastAPI (:8000)..."
  _bg "ml" "$LOG_DIR/ml-api.log" \
    bash -c "cd '$ML_DIR' && .venv/bin/uvicorn app.main:app --host 0.0.0.0 --port 8000"
  sleep 5
  if curl -sfo /dev/null --max-time 3 "http://localhost:8000/api/v1/health/live" 2>/dev/null; then
    ok "ML API is running"
  else
    warn "ML API may still be starting up — check $LOG_DIR/ml-api.log"
  fi
}

start_goworker() {
  if is_running "goworker"; then
    warn "Go worker is already running"
    return 0
  fi
  info "Starting Go worker..."
  _bg "goworker" "$LOG_DIR/goworker.log" go run "$PROJECT_ROOT/cmd/worker/main.go"
  sleep 3
  if is_running "goworker"; then
    ok "Go worker is running"
  else
    warn "Go worker may still be starting up — check $LOG_DIR/goworker.log"
  fi
}

start_worker() {
  if is_running "worker"; then
    warn "Celery worker is already running"
    return 0
  fi
  info "Starting Celery worker (Whisper)..."
  _bg "worker" "$LOG_DIR/worker.log" \
    env "HF_HOME=$HOME/.cache/huggingface" \
    bash -c "cd '$ML_DIR' && \
      .venv/bin/celery -A app.workers.celery_app worker \
      -Q lyrics_queue -l info --concurrency=1"
  # Wait for model load (can take 15-30s)
  local i=0
  while [ $i -lt 30 ]; do
    if grep -q "Whisper model loaded successfully" "$LOG_DIR/worker.log" 2>/dev/null; then
      ok "Celery worker is running (Whisper model loaded)"
      return 0
    fi
    if grep -q "failed to load\|Traceback\|ERROR.*whisper" "$LOG_DIR/worker.log" 2>/dev/null; then
      err "Celery worker failed — check $LOG_DIR/worker.log"
      return 1
    fi
    sleep 2
    i=$((i + 1))
  done
  warn "Celery worker may still be loading model — check $LOG_DIR/worker.log"
}

# ── Service Stop Functions ─────────────────────────────────────────────

stop_service() {
  local name="$1"
  local pid_file="$PID_DIR/$name.pid"
  if [ -f "$pid_file" ]; then
    local pid
    pid=$(cat "$pid_file")
    info "Stopping $name (PID $pid)..."
    kill "$pid" 2>/dev/null || true
    sleep 1
    # Force if still alive
    kill -0 "$pid" 2>/dev/null && kill -9 "$pid" 2>/dev/null || true
    remove_pid "$name"
    ok "$name stopped"
  fi
}

stop_all() {
  info "Stopping all services..."
  for svc in goworker worker ml web api; do
    stop_service "$svc"
  done
  ok "All services stopped"
}

# ── Status ─────────────────────────────────────────────────────────────

cmd_status() {
  echo ""
  echo -e "${CYAN}═══════════════════════════════════════════${NC}"
  echo -e "${CYAN}        Muse — Service Status              ${NC}"
  echo -e "${CYAN}═══════════════════════════════════════════${NC}"
  echo ""

  # Docker infrastructure
  echo -e "${YELLOW}Infrastructure:${NC}"
  for container in musicapp-db musicapp_redis musicapp_minio; do
    local display_name
    case "$container" in
      musicapp-db)     display_name="PostgreSQL" ;;
      musicapp_redis)  display_name="Redis" ;;
      musicapp_minio)  display_name="MinIO" ;;
      *)               display_name="$container" ;;
    esac
    local status
    status=$(docker inspect "$container" --format='{{.State.Status}}' 2>/dev/null || echo "not found")
    if [ "$status" = "running" ]; then
      echo -e "  ${GREEN}●${NC} $display_name"
    elif [ "$status" = "not found" ]; then
      echo -e "  ${RED}○${NC} $display_name (not found)"
    else
      echo -e "  ${RED}○${NC} $display_name ($status)"
    fi
  done

  echo ""
  echo -e "${YELLOW}Application:${NC}"

  # Web
  if is_running "web"; then
    echo -e "  ${GREEN}●${NC} Frontend :3000  $(curl -s -o /dev/null -w "HTTP %{http_code}" --max-time 2 http://localhost:3000/ 2>/dev/null)"
  else
    echo -e "  ${RED}○${NC} Frontend :3000  — not running"
  fi

  # API
  if is_running "api"; then
    local api_status
    api_status=$(curl -s -o /dev/null -w "HTTP %{http_code}" --max-time 2 http://localhost:8080/api/v1/health 2>/dev/null)
    echo -e "  ${GREEN}●${NC} API      :8080  $api_status"
  else
    echo -e "  ${RED}○${NC} API      :8080  — not running"
  fi

  # ML
  if is_running "ml"; then
    local ml_status
    ml_status=$(curl -s -o /dev/null -w "HTTP %{http_code}" --max-time 2 http://localhost:8000/api/v1/health/live 2>/dev/null)
    echo -e "  ${GREEN}●${NC} ML API   :8000  $ml_status"
  else
    echo -e "  ${RED}○${NC} ML API   :8000  — not running"
  fi

  # Celery Worker
  if is_running "worker"; then
    echo -e "  ${GREEN}●${NC} Celery Worker    (queue: lyrics_queue)"
  else
    echo -e "  ${RED}○${NC} Celery Worker    — not running"
  fi

  # Go Worker
  if is_running "goworker"; then
    echo -e "  ${GREEN}●${NC} Go Worker"
  else
    echo -e "  ${RED}○${NC} Go Worker        — not running"
  fi

  echo ""
}

# ── Logs ───────────────────────────────────────────────────────────────

cmd_logs() {
  local service="${1:-all}"
  mkdir -p "$LOG_DIR"

  case "$service" in
    all)
      echo -e "${CYAN}Tailing all logs (Ctrl+C to stop)${NC}"
      echo -e "${YELLOW}── web.log ─────────────────────${NC}"
      tail -f "$LOG_DIR/web.log" &
      local pid_web=$!
      echo -e "${YELLOW}── api.log ─────────────────────${NC}"
      tail -f "$LOG_DIR/api.log" &
      local pid_api=$!
      echo -e "${YELLOW}── ml-api.log ──────────────────${NC}"
      tail -f "$LOG_DIR/ml-api.log" &
      local pid_ml=$!
      echo -e "${YELLOW}── worker.log ──────────────────${NC}"
      tail -f "$LOG_DIR/worker.log" &
      local pid_worker=$!
      echo -e "${YELLOW}── goworker.log ────────────────${NC}"
      tail -f "$LOG_DIR/goworker.log" &
      local pid_goworker=$!
      # Trap Ctrl+C to kill all tails
      trap 'kill $pid_web $pid_api $pid_ml $pid_worker $pid_goworker 2>/dev/null; exit 0' INT TERM
      wait
      ;;
    web)      tail -f "$LOG_DIR/web.log" ;;
    api)      tail -f "$LOG_DIR/api.log" ;;
    ml)       tail -f "$LOG_DIR/ml-api.log" ;;
    worker)   tail -f "$LOG_DIR/worker.log" ;;
    goworker) tail -f "$LOG_DIR/goworker.log" ;;
    *)
      echo "Unknown service: $service"
      echo "Usage: bash scripts/dev.sh logs [web|api|ml|worker|goworker]"
      exit 1
      ;;
  esac
}

# ── Start / Stop ───────────────────────────────────────────────────────

cmd_start() {
  cleanup_logs

  echo -e "${CYAN}═══════════════════════════════════════════${NC}"
  echo -e "${CYAN}      Muse — Starting all services         ${NC}"
  echo -e "${CYAN}═══════════════════════════════════════════${NC}"
  echo ""

  start_infra
  echo ""

  # Start all app services in parallel
  start_api &
  local pid_api=$!
  start_web &
  local pid_web=$!
  start_ml &
  local pid_ml=$!
  start_worker &
  local pid_worker=$!
  start_goworker &
  local pid_goworker=$!
  wait $pid_api $pid_web $pid_ml $pid_worker $pid_goworker 2>/dev/null || true
  echo ""

  cmd_status

  echo -e "${CYAN}═══════════════════════════════════════════${NC}"
  echo -e "Logs:  ${GREEN}bash scripts/dev.sh logs${NC}"
  echo -e "       ${GREEN}bash scripts/dev.sh logs api${NC}   (single service)"
  echo -e "Stop:  ${GREEN}bash scripts/dev.sh stop${NC}"
  echo -e "Check: ${GREEN}bash scripts/dev.sh status${NC}"
  echo -e "${CYAN}═══════════════════════════════════════════${NC}"
}

cmd_stop() {
  echo -e "${CYAN}═══════════════════════════════════════════${NC}"
  echo -e "${CYAN}      Muse — Stopping all services         ${NC}"
  echo -e "${CYAN}═══════════════════════════════════════════${NC}"
  stop_all
  cmd_status
  ok "Done. Use 'bash scripts/dev.sh start' to restart."
}

# ═══════════════════════════════════════════════════════════════════════════
# Main
# ═══════════════════════════════════════════════════════════════════════════

cd "$PROJECT_ROOT"

case "${1:-help}" in
  start)
    cmd_start
    ;;
  stop)
    cmd_stop
    ;;
  restart)
    cmd_stop
    sleep 2
    cmd_start
    ;;
  status)
    cmd_status
    ;;
  logs)
    cmd_logs "${2:-all}"
    ;;
  help|*)
    echo ""
    echo -e "${CYAN}Muse Dev — Unified development server manager${NC}"
    echo ""
    echo "Usage: bash scripts/dev.sh <command> [args]"
    echo ""
    echo "Commands:"
    echo "  start          Start all services (infra, API, frontend, ML)"
    echo "  stop           Stop all services"
    echo "  restart        Restart everything"
    echo "  status         Show health of all services"
    echo "  logs           Tail all logs (follow)"
    echo "  logs api       Tail a specific service log"
    echo ""
    echo "Services managed:"
    echo "  infra    → Postgres, Redis, MinIO, OpenSearch (Docker)"
    echo "  api      → Go/Gin backend      (:8080)"
    echo "  web      → Vite frontend        (:3000)"
    echo "  ml       → FastAPI ML service   (:8000)"
    echo "  worker   → Celery (Whisper)     (lyrics_queue)"
    echo "  goworker → Go worker            (cmd/worker)"
    echo ""
    ;;
esac

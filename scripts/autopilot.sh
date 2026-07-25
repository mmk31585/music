#!/usr/bin/env bash
set -euo pipefail

# Muse Autopilot — Autonomous development team for Muse
# Starts frontend, API, and watchdog in a tmux session.
# The watchdog monitors health and captures errors.
# Invoke @muse-sre or @muse in opencode for autonomous fixes.
#
# Usage:
#   bash scripts/autopilot.sh          Start the autonomous team
#   bash scripts/autopilot.sh status   Check current status
#   bash scripts/autopilot.sh attach   Attach to tmux session
#   bash scripts/autopilot.sh stop     Stop everything
#   bash scripts/autopilot.sh fix      Run SRE fix cycle (if servers are up)

SESSION="muse"
STATUS_FILE=".watchdog-status.json"

ensure_tmux() {
  if ! command -v tmux &>/dev/null; then
    echo "tmux not found. Install it:"
    echo "  sudo apt install tmux  # Debian/Ubuntu"
    echo "  sudo pacman -S tmux    # Arch"
    echo "  brew install tmux      # macOS"
    exit 1
  fi
}

cmd_status() {
  if [ -f "$STATUS_FILE" ]; then
    echo "=== Muse Autopilot Status ==="
    python3 -m json.tool "$STATUS_FILE" 2>/dev/null || cat "$STATUS_FILE"
    echo ""
    echo "=== Recent Watchdog Log ==="
    tail -10 .watchdog.log 2>/dev/null || echo "(no log)"
  else
    echo "Autopilot is not running (no status file)"
  fi
  if tmux has-session -t "$SESSION" 2>/dev/null; then
    echo "tmux session '$SESSION' is active"
  else
    echo "tmux session '$SESSION' is not running"
  fi
}

cmd_start() {
  if tmux has-session -t "$SESSION" 2>/dev/null; then
    echo "Autopilot already running. Use:"
    echo "  bash scripts/autopilot.sh attach  # View logs"
    echo "  bash scripts/autopilot.sh stop    # Stop"
    exit 0
  fi

  ensure_tmux

  # Clean stale processes
  lsof -ti :5173 2>/dev/null | xargs kill -9 2>/dev/null || true
  lsof -ti :8080 2>/dev/null | xargs kill -9 2>/dev/null || true

  tmux new-session -d -s "$SESSION" -n "muse" \; \
    send-keys "cd frontend && npm run dev" C-m \; \
    split-window -h \; \
    send-keys "go run cmd/api/main.go 2>&1 | tee .watchdog-api.log" C-m \; \
    split-window -v \; \
    send-keys "bash scripts/watchdog.sh" C-m \; \
    select-pane -t 0

  sleep 3
  echo "=== Muse Autopilot Started ==="
  echo "Attach:  bash scripts/autopilot.sh attach"
  echo "Status:  bash scripts/autopilot.sh status"
  echo "Stop:    bash scripts/autopilot.sh stop"
  echo ""
  echo "To auto-fix issues, in opencode say:"
  echo "  @muse-sre run the auto-pilot cycle"
  echo "  @muse scan and fix everything"
}

cmd_stop() {
  echo "Stopping autopilot..."
  if tmux has-session -t "$SESSION" 2>/dev/null; then
    tmux kill-session -t "$SESSION"
  fi
  lsof -ti :5173 2>/dev/null | xargs kill -9 2>/dev/null || true
  lsof -ti :8080 2>/dev/null | xargs kill -9 2>/dev/null || true
  echo "Done. Run 'bash scripts/autopilot.sh start' to restart."
}

cmd_attach() {
  ensure_tmux
  if tmux has-session -t "$SESSION" 2>/dev/null; then
    tmux attach-session -t "$SESSION"
  else
    echo "Autopilot is not running. Start it first:"
    echo "  bash scripts/autopilot.sh start"
    exit 1
  fi
}

cmd_fix() {
  echo "=== SRE Fix Cycle ==="
  echo "Checking server health..."
  echo ""

  local fe_code api_code
  fe_code=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:5173 2>/dev/null || echo "000")
  api_code=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/api/v1/health 2>/dev/null || echo "000")

  echo "Frontend (:5173): HTTP $fe_code"
  echo "API      (:8080): HTTP $api_code"
  echo ""

  if [ -f "$STATUS_FILE" ]; then
    echo "=== Last Errors ==="
    python3 -c "
import json
with open('$STATUS_FILE') as f:
  s = json.load(f)
for e in s.get('last_errors', [])[-5:]:
  print(f\"  [{e.get('time','?')[-8:]}] {e.get('msg','')}\")
" 2>/dev/null || echo "  (could not read)"
    echo ""
  fi

  echo "=== Watchdog Log (last 15 lines) ==="
  tail -15 .watchdog.log 2>/dev/null || echo "  (empty)"

  echo ""
  echo "To fix issues, in opencode run:  @muse-sre"
}

# === MAIN ===
cd "$(dirname "$0")/.."

case "${1:-start}" in
  start)   cmd_start ;;
  status)  cmd_status ;;
  attach)  cmd_attach ;;
  stop)    cmd_stop ;;
  fix)     cmd_fix ;;
  *)
    echo "Usage: bash scripts/autopilot.sh {start|status|attach|stop|fix}"
    exit 1
    ;;
esac

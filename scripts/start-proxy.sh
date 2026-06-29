#!/bin/bash
# Start xray SOCKS5 proxy for Muse import pipeline.
# This bypasses the network's transparent proxy (10.10.34.36) to reach
# SoundCloud, YouTube, Deezer, etc.
#
# Usage:  ./scripts/start-proxy.sh
#         ./scripts/start-proxy.sh stop

set -e

ACTION="${1:-start}"
XRAY_BIN="/tmp/xray/xray"
XRAY_CONFIG="$HOME/.config/xray/config.json"
XRAY_LOG="/tmp/xray.log"
XRAY_PIDFILE="/tmp/xray.pid"

case "$ACTION" in
  start)
    # Kill any existing instance
    pkill -f "xray.*config.json" 2>/dev/null || true
    sleep 1

    if [ ! -f "$XRAY_BIN" ]; then
      echo "Installing xray..."
      curl -sL "https://github.com/XTLS/Xray-core/releases/latest/download/Xray-linux-64.zip" -o /tmp/xray.zip
      unzip -o /tmp/xray.zip -d /tmp/xray/ 2>/dev/null
      chmod +x "$XRAY_BIN"
    fi

    # Ensure config exists
    if [ ! -f "$XRAY_CONFIG" ]; then
      echo "Config not found at $XRAY_CONFIG"
      echo "Run the setup or create the config manually."
      exit 1
    fi

    nohup "$XRAY_BIN" run -c "$XRAY_CONFIG" > "$XRAY_LOG" 2>&1 &
    XRAY_PID=$!
    echo "$XRAY_PID" > "$XRAY_PIDFILE"

    # Wait and verify
    sleep 3
    if curl -s --socks5-hostname 127.0.0.1:1081 --connect-timeout 3 -o /dev/null -w "%{http_code}" "https://www.youtube.com" 2>/dev/null | grep -q 200; then
      echo "✅ xray proxy ready (PID $XRAY_PID)"
      echo "   SOCKS5: 127.0.0.1:1081"
      echo "   Server: sw.eorqen.ir:2050"
      exit 0
    else
      echo "❌ xray proxy failed to start"
      tail -5 "$XRAY_LOG"
      exit 1
    fi
    ;;

  stop)
    if [ -f "$XRAY_PIDFILE" ]; then
      kill "$(cat "$XRAY_PIDFILE")" 2>/dev/null && echo "xray stopped" || echo "xray not running"
      rm -f "$XRAY_PIDFILE"
    else
      pkill -f "xray.*config.json" 2>/dev/null && echo "xray stopped" || echo "xray not running"
    fi
    ;;

  status)
    if pgrep -f "xray.*config.json" > /dev/null; then
      echo "✅ xray is running (PID: $(pgrep -f 'xray.*config.json'))"
      echo "   Log: $XRAY_LOG"
      echo "   Proxy: 127.0.0.1:1081"
    else
      echo "❌ xray is not running"
    fi
    ;;

  *)
    echo "Usage: $0 {start|stop|status}"
    exit 1
    ;;
esac

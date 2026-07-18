#!/bin/bash
set -e

ROOT="$(cd "$(dirname "$0")" && pwd)"
BACKEND_DIR="$ROOT/backend"
FRONTEND_DIR="$ROOT/frontend"

APP_NAME="sheet-server"
BIN="$BACKEND_DIR/bin/$APP_NAME"
PID_FILE="$BACKEND_DIR/bin/$APP_NAME.pid"
LOG_FILE="$BACKEND_DIR/bin/$APP_NAME.log"
CONFIG_FILE="$BACKEND_DIR/config.yaml"

# ---------- build ----------
build_backend() {
    echo "=== building backend ==="
    cd "$BACKEND_DIR"
    go build -o "$BIN" ./cmd/server
    echo "backend done: $BIN"
}

build_frontend() {
    echo "=== building frontend ==="
    cd "$FRONTEND_DIR"
    if [ ! -d "node_modules" ]; then
        npm install
    fi
    npm run build
    echo "frontend done: $FRONTEND_DIR/dist"
}

build() {
    build_frontend
    build_backend
}

# ---------- run ----------
start() {
    if [ ! -f "$CONFIG_FILE" ]; then
        echo "config not found: $CONFIG_FILE"
        echo "copy config.example.yaml to config.yaml and fill in your config"
        exit 1
    fi

    if [ -f "$PID_FILE" ] && kill -0 "$(cat "$PID_FILE")" 2>/dev/null; then
        echo "$APP_NAME is already running (pid $(cat "$PID_FILE"))"
        return
    fi

    if [ ! -f "$BIN" ]; then
        echo "binary not found, building ..."
        build
    fi

    echo "starting $APP_NAME ..."
    mkdir -p "$BACKEND_DIR/bin"
    nohup "$BIN" --config "$CONFIG_FILE" >> "$LOG_FILE" 2>&1 &
    echo $! > "$PID_FILE"
    echo "started (pid $!)"
}

stop() {
    if [ ! -f "$PID_FILE" ]; then
        echo "$APP_NAME is not running"
        return
    fi
    PID=$(cat "$PID_FILE")
    if kill -0 "$PID" 2>/dev/null; then
        echo "stopping $APP_NAME (pid $PID) ..."
        kill "$PID"
        sleep 2
        if kill -0 "$PID" 2>/dev/null; then
            kill -9 "$PID"
        fi
        echo "stopped"
    fi
    rm -f "$PID_FILE"
}

restart() {
    stop
    sleep 1
    start
}

status() {
    if [ -f "$PID_FILE" ] && kill -0 "$(cat "$PID_FILE")" 2>/dev/null; then
        echo "$APP_NAME is running (pid $(cat "$PID_FILE"))"
    else
        echo "$APP_NAME is not running"
    fi
}

logs() {
    tail -f "$LOG_FILE"
}

case "${1:-start}" in
    build)        build ;;
    build-fe)     build_frontend ;;
    build-be)     build_backend ;;
    start)        start ;;
    stop)         stop ;;
    restart)      restart ;;
    status)       status ;;
    logs)         logs ;;
    *)
        echo "usage: $0 {build|start|stop|restart|status|logs}"
        exit 1
        ;;
esac

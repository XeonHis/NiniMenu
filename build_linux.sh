#!/usr/bin/env bash
set -Eeuo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
FRONTEND_DIR="$ROOT_DIR/frontend"
BACKEND_DIR="$ROOT_DIR/backend"
DIST_DIR="$ROOT_DIR/dist"
ARCHIVE="$ROOT_DIR/ninimenu-linux-amd64.tar.gz"

log() {
    printf '%s\n' "$*"
}

fail() {
    printf '[ERROR] %s\n' "$*" >&2
    exit 1
}

require_command() {
    command -v "$1" >/dev/null 2>&1 || fail "Required command not found: $1"
}

copy_dir_contents() {
    local source_dir="$1"
    local target_dir="$2"

    mkdir -p "$target_dir"
    cp -a "$source_dir"/. "$target_dir"/
}

require_command npm
require_command go
require_command tar

[[ -d "$FRONTEND_DIR" ]] || fail "Missing frontend directory: $FRONTEND_DIR"
[[ -d "$BACKEND_DIR" ]] || fail "Missing backend directory: $BACKEND_DIR"

log "============================================"
log "  NiniMenu Build Script (Linux amd64)"
log "============================================"
log

log "[1/4] Building frontend..."
(
    cd "$FRONTEND_DIR"
    npm install
    npm run build
)
[[ -d "$FRONTEND_DIR/dist" ]] || fail "Frontend build output not found: $FRONTEND_DIR/dist"
log "[1/4] Frontend build done!"
log

log "[2/4] Copying frontend dist to backend/static..."
rm -rf "$BACKEND_DIR/static"
copy_dir_contents "$FRONTEND_DIR/dist" "$BACKEND_DIR/static"
log "[2/4] Frontend dist copied!"
log

log "[3/4] Building Go backend (linux/amd64)..."
rm -rf "$DIST_DIR"
mkdir -p "$DIST_DIR"
(
    cd "$BACKEND_DIR"
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o "$DIST_DIR/ninimenu" ./cmd/server
)
chmod +x "$DIST_DIR/ninimenu"
log "[3/4] Go backend build done!"
log

log "[4/4] Packaging distribution files..."
copy_dir_contents "$BACKEND_DIR/static" "$DIST_DIR/static"
mkdir -p "$DIST_DIR/data"

if [[ -f "$ROOT_DIR/.env" ]]; then
    cp "$ROOT_DIR/.env" "$DIST_DIR/env.bak"
fi

if [[ -d "$BACKEND_DIR/uploads" ]]; then
    copy_dir_contents "$BACKEND_DIR/uploads" "$DIST_DIR/uploads"
else
    mkdir -p "$DIST_DIR/uploads"
fi

log "Creating archive..."
tar -czf "$ARCHIVE" -C "$DIST_DIR" .

log
log "============================================"
log "  Build complete!"
log "  Output: ninimenu-linux-amd64.tar.gz"
log "============================================"
log
log "  Deploy instructions:"
log "  1. Upload ninimenu-linux-amd64.tar.gz to your Linux server"
log "  2. Extract: tar -xzf ninimenu-linux-amd64.tar.gz -C /opt/ninimenu"
log "  3. chmod +x ninimenu"
log "  4. Run: ./ninimenu"
log "  5. Open browser: http://SERVER_IP:8080"
log
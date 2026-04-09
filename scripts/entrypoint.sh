#!/bin/sh
set -e

# Helper function for formatted logging
log() {
    local message="$1"
    echo "[ENTRYPOINT] $(date '+%Y/%m/%d - %H:%M:%S') | $message"
    return 0
}

log "Starting container..."
log "Ready!"
log "Launching app..."
log "API endpoints | http://localhost:9000"
log "Swagger UI    | http://localhost:9000/swagger/index.html"
exec "$@"

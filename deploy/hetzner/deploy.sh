#!/bin/bash

# Deploy Script für Hetzner Server
# Wird vom CI/CD Workflow aufgerufen

set -e  # Exit on error

echo "🚀 Starting deployment..."

# Konfiguration
PROJECT_DIR="/home/deploy/trip-manager"
COMPOSE_FILE="$PROJECT_DIR/docker-compose.yaml"
LOG_FILE="$PROJECT_DIR/logs/deploy-$(date +%Y%m%d_%H%M%S).log"

# Erstelle Log Verzeichnis falls nicht existent
mkdir -p "$(dirname "$LOG_FILE")"

# Logging Function
log() {
    echo "[$(date +'%Y-%m-%d %H:%M:%S')] $1" | tee -a "$LOG_FILE"
}

log "📋 Deployment started"

# 1. Pull latest images
log "📦 Pulling latest Docker images..."
cd "$PROJECT_DIR"
docker-compose pull 2>&1 | tee -a "$LOG_FILE"

# 2. Stop old containers
log "🛑 Stopping old containers..."
docker-compose down 2>&1 | tee -a "$LOG_FILE"

# 3. Start new containers
log "🟢 Starting new containers..."
docker-compose up -d 2>&1 | tee -a "$LOG_FILE"

# 4. Health check
log "🏥 Running health checks..."
sleep 5

# Check Backend
if ! curl -f http://localhost:8000/health > /dev/null 2>&1; then
    log "❌ Backend health check failed!"
    log "📋 Logs:"
    docker-compose logs backend >> "$LOG_FILE"
    exit 1
fi
log "✅ Backend is healthy"

# Check Frontend (optional - needs to be configured)
if ! curl -f http://frontend:3000 > /dev/null 2>&1; then
    log "⚠️  Frontend not responding (might still be building)"
else
    log "✅ Frontend is healthy"
fi

# 5. Reload Caddy configuration
log "🔄 Reloading Caddy configuration..."
docker-compose exec -T caddy caddy reload -c /etc/caddy/Caddyfile 2>&1 | tee -a "$LOG_FILE" || log "⚠️  Caddy reload warning (might be first start)"

# 6. Cleanup old images
log "🧹 Cleaning up old images..."
docker image prune -f 2>&1 | tee -a "$LOG_FILE"

log "✅ Deployment completed successfully!"
echo ""
echo "📊 Current containers:"
docker-compose ps

exit 0


#!/bin/bash

# Manual Deployment Script for Hetzner Server
# Performs the same steps as the GitHub Actions pipeline
#
# Usage:
#   ./manual-deploy.sh <server_ip> <username> <port>
#
# Example:
#   ./manual-deploy.sh 192.168.1.100 deploy 22
#
# Requirements:
#   - SSH key configured for passwordless login
#   - .env file in the same directory with all required variables
#   - SSH_AUTH_SOCK should be set for ssh-agent forwarding

set -e

# Color output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Logging function
log() {
    echo -e "${BLUE}[$(date +'%Y-%m-%d %H:%M:%S')]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1" >&2
    exit 1
}

success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

# Check arguments
if [ $# -ne 3 ]; then
    echo "Usage: $0 <server_ip> <username> <port>"
    echo "Example: $0 192.168.1.100 deploy 22"
    exit 1
fi

SERVER_IP="$1"
USERNAME="$2"
SSH_PORT="$3"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(cd "$SCRIPT_DIR/../.." && pwd)"
DEPLOYMENT_DIR="$SCRIPT_DIR/deployment"
ENV_FILE="$SCRIPT_DIR/.env"
REMOTE_PATH="/home/$USERNAME/app/cloud"
REMOTE_HOST="$USERNAME@$SERVER_IP"

# Check if .env file exists
if [ ! -f "$ENV_FILE" ]; then
    error ".env file not found at $ENV_FILE"
fi

log "📂 Loading environment variables from $ENV_FILE..."
set -a
source "$ENV_FILE"
set +a

# Validate required variables
required_vars=("GHCR_USERNAME" "GHCR_PAT" "DB_USER" "DB_PASSWORD" "DB_NAME" "JWT_SECRET" "DOMAIN")
for var in "${required_vars[@]}"; do
    if [ -z "${!var}" ]; then
        error "Missing required environment variable: $var"
    fi
done

log "🔐 Connecting to server $SERVER_IP ($USERNAME) on port $SSH_PORT..."

# 1. Copy deploy files to server
log "📤 Copying deploy files to server..."
log "📁 Creating remote directory: $REMOTE_PATH"
ssh -p "$SSH_PORT" "$REMOTE_HOST" "mkdir -p $REMOTE_PATH" || error "Failed to create remote directory $REMOTE_PATH"

scp -P "$SSH_PORT" "$DEPLOYMENT_DIR/deploy.sh" "$REMOTE_HOST:$REMOTE_PATH/" || error "Failed to copy deploy.sh"
scp -P "$SSH_PORT" "$PROJECT_DIR/docker-compose.yaml" "$REMOTE_HOST:$REMOTE_PATH/" || error "Failed to copy docker-compose.yaml"

log "✅ Deploy files copied"

# 2. Generate and copy Caddyfile with environment variable substitution
log "📤 Generating and copying Caddyfile to server..."
CADDYFILE_TMP=$(mktemp)
envsubst < "$DEPLOYMENT_DIR/Caddyfile" > "$CADDYFILE_TMP" || error "Failed to generate Caddyfile"
scp -P "$SSH_PORT" "$CADDYFILE_TMP" "$REMOTE_HOST:$REMOTE_PATH/Caddyfile" || error "Failed to copy Caddyfile"
rm -f "$CADDYFILE_TMP"

log "✅ Caddyfile copied"

# 3. Generate .env file on server
log "📝 Generating .env file on server..."
ssh -p "$SSH_PORT" "$REMOTE_HOST" \
    tee $REMOTE_PATH/.env << EOF > /dev/null
GHCR_USERNAME=$GHCR_USERNAME
GHCR_PAT=$GHCR_PAT
DB_USER=$DB_USER
DB_PASSWORD=$DB_PASSWORD
DB_NAME=$DB_NAME
JWT_SECRET=$JWT_SECRET
DOMAIN=$DOMAIN
EOF

log "✅ .env file generated"

# 4. Execute deploy script on server
log "🚀 Executing deploy script on server..."
ssh -p "$SSH_PORT" "$REMOTE_HOST" \
    "chmod +x $REMOTE_PATH/deploy.sh && cd $REMOTE_PATH && ./deploy.sh" || error "Deployment failed"

log "✅ Deploy script executed"

# 5. Verify deployment
log "🔍 Verifying deployment..."
ssh -p "$SSH_PORT" "$REMOTE_HOST" \
    "docker-compose -f $REMOTE_PATH/docker-compose.yaml ps" || error "Verification failed"

success "🎉 Deployment completed successfully!"


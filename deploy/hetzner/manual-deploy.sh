#!/bin/bash

# Manual Deployment Script for Trip Manager on Remote Server
# Performs the same steps as the GitHub Actions pipeline
#
# This script supports multiple ways to specify the server:
#
# 1. SSH Config Alias (recommended):
#    ./manual-deploy.sh myserver
#
# 2. SSH Config Alias with Port Override:
#    ./manual-deploy.sh myserver 48222
#
# 3. Direct IP + Username + Port:
#    ./manual-deploy.sh deploy@167.235.66.0 48222
#
# 4. IP + Port only (uses 'deploy' as default username):
#    ./manual-deploy.sh 167.235.66.0 48222
#
# Requirements:
#   - SSH key configured for passwordless login
#   - .env file in the same directory with all required variables
#   - SSH access to the remote server

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
if [ $# -lt 1 ] || [ $# -gt 2 ]; then
    echo "Usage: $0 <ssh_host> [port]"
    echo ""
    echo "Examples:"
    echo ""
    echo "  # Using SSH config alias (recommended):"
    echo "  $0 myserver"
    echo "  $0 myserver 48222         # Override port if needed"
    echo ""
    echo "  # Using IP + Port:"
    echo "  $0 deploy@192.168.1.100 22"
    echo "  $0 deploy@167.235.66.0 48222"
    echo ""
    echo "  # Using just IP (assumes 'deploy' user, SSH config will handle port):"
    echo "  $0 167.235.66.0"
    echo ""
    echo "Requirements:"
    echo "  - .env file in this directory with required variables"
    echo "  - SSH key configured (via ssh-agent or SSH config)"
    echo ""
    exit 1
fi

SSH_HOST="$1"
SSH_PORT="${2:-}"  # Optional second argument for port

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(cd "$SCRIPT_DIR/../.." && pwd)"
DEPLOYMENT_DIR="$SCRIPT_DIR/deployment"
ENV_FILE="$SCRIPT_DIR/.env"

# Parse SSH_HOST to extract username, server, and optional port
# Handles multiple formats:
# - "user@host" or "host" (SSH config)
# - "user@ip" or "ip"
# - With or without port in SSH_PORT

if [[ "$SSH_HOST" =~ @ ]]; then
    # Format: user@host or user@ip
    USERNAME="${SSH_HOST%@*}"
    SERVER="${SSH_HOST#*@}"
else
    # Format: host or ip (no user specified)
    SERVER="$SSH_HOST"
    # Try to get username from SSH config
    USERNAME=$(ssh -G "$SERVER" 2>/dev/null | grep "^user " | awk '{print $2}' || echo "")
    if [ -z "$USERNAME" ]; then
        USERNAME="deploy"  # Fallback to 'deploy'
    fi
fi

REMOTE_PATH="/home/$USERNAME/app/cloud"
REMOTE_HOST="$SSH_HOST"

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

log "🔐 Connecting to server..."
log "   Host: $REMOTE_HOST"
log "   User: $USERNAME"
if [ -n "$SSH_PORT" ]; then
    log "   Port: $SSH_PORT (override)"
else
    log "   Port: from SSH config or default"
fi
log "   Remote path: $REMOTE_PATH"

# Helper function for SSH commands with optional port
ssh_cmd() {
    if [ -n "$SSH_PORT" ]; then
        ssh -p "$SSH_PORT" "$REMOTE_HOST" "$@"
    else
        ssh "$REMOTE_HOST" "$@"
    fi
}

# Helper function for SCP commands with optional port
scp_cmd() {
    if [ -n "$SSH_PORT" ]; then
        scp -P "$SSH_PORT" "$@"
    else
        scp "$@"
    fi
}

# 1. Copy deploy files to server
log "📤 Copying deploy files to server..."
log "📁 Creating remote directory: $REMOTE_PATH"
ssh_cmd "mkdir -p $REMOTE_PATH" || error "Failed to create remote directory $REMOTE_PATH"

scp_cmd "$DEPLOYMENT_DIR/deploy.sh" "$REMOTE_HOST:$REMOTE_PATH/" || error "Failed to copy deploy.sh"
scp_cmd "$PROJECT_DIR/docker-compose.yaml" "$REMOTE_HOST:$REMOTE_PATH/" || error "Failed to copy docker-compose.yaml"

# Copy docker directory with all initialization scripts
log "📤 Copying docker directory to server..."
scp_cmd -r "$PROJECT_DIR/docker" "$REMOTE_HOST:$REMOTE_PATH/" || error "Failed to copy docker directory"

# Make minio-init.sh executable on the server
log "🔧 Making docker scripts executable..."
ssh_cmd "chmod +x $REMOTE_PATH/docker/*.sh" || error "Failed to make docker scripts executable"

log "✅ Deploy files copied"

# 2. Generate and copy Caddyfile with environment variable substitution
log "📤 Generating and copying Caddyfile to server..."
CADDYFILE_TMP=$(mktemp)
envsubst < "$DEPLOYMENT_DIR/Caddyfile" > "$CADDYFILE_TMP" || error "Failed to generate Caddyfile"
scp_cmd "$CADDYFILE_TMP" "$REMOTE_HOST:$REMOTE_PATH/Caddyfile" || error "Failed to copy Caddyfile"
rm -f "$CADDYFILE_TMP"

log "✅ Caddyfile copied"

# 3. Generate .env file on server
log "📝 Generating .env file on server..."
ssh_cmd \
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
ssh_cmd \
    "chmod +x $REMOTE_PATH/deploy.sh && cd $REMOTE_PATH && ./deploy.sh" || error "Deployment failed"

log "✅ Deploy script executed"

# 5. Verify deployment
log "🔍 Verifying deployment..."
ssh_cmd \
    "docker-compose -f $REMOTE_PATH/docker-compose.yaml ps" || error "Verification failed"

success "🎉 Deployment completed successfully!"


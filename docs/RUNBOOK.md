# Runbook: Deployment Guide

Complete guide for setting up and deploying the Trip Manager application to a Hetzner server.

---

## Current Server Setup

### Specifications

The Trip Manager application is currently deployed on the following infrastructure:

| Component | Details         |
|-----------|-----------------|
| **Hoster** | Hetzner Cloud   |
| **Operating System** | Ubuntu (latest LTS) |
| **vCPU** | 2 cores         |
| **RAM** | 4 GB            |
| **Disk Storage** | 40 GB SSD       |
| **Access URL** | https://www.travel-nugget.duckdns.org |
| **API Endpoint** | https://www.travel-nugget.duckdns.org/api |

### Infrastructure Independence

The deployment uses **Docker Compose** for containerization. This means:

- ✅ **Distribution Independent** - Can run on any Linux distribution, can also run on Windows and MacOS (with small modifications depending on processor architecture) that supports Docker
- ✅ **Reproducible** - Same containers run identically across different servers
- ✅ **Portable** - Easy to migrate to different hosting providers (AWS, DigitalOcean, etc.)
- ✅ **Scalable** - Can upgrade server specs without changing deployment process

### Why Docker Compose?

- All services (Frontend, Backend, Database, Reverse Proxy) are containerized
- No dependency on system packages or Python/Node versions
- Services communicate via Docker internal networking
- Environment variables manage configuration across different deployments
- Single `docker-compose.yaml` file describes entire stack

### Current Deployment Method

All deployments use **one of two methods**:

1. **Automated**: GitHub Actions Pipeline (main branch pushes)
2. **Manual**: `deploy/hetzner/manual-deploy.sh` script (via SSH)

Both methods perform identical steps, ensuring consistent deployments regardless of method used.

---

1. [SSH Key Setup](#ssh-key-setup)
2. [DuckDNS Domain Setup](#duckdns-domain-setup)
3. [Server Setup](#server-setup)
4. [GitHub Secrets Configuration](#github-secrets-configuration)
5. [Automated Deployment (GitHub Actions)](#automated-deployment-github-actions)
6. [Manual Deployment](#manual-deployment)
7. [Architecture Overview](#architecture-overview)
8. [Health Checks](#health-checks)
9. [Logs](#logs)
10. [Troubleshooting](#troubleshooting)
11. [Rollback](#rollback)

---

## SSH Key Setup

### 1. Generate SSH Key (if not already present)

```bash
ssh-keygen -t ed25519 -f ~/.ssh/hetzner_deploy -C "github-deploy@trip-manager"
```

### 2. Copy Public Key to Hetzner Server

```bash
ssh-copy-id -i ~/.ssh/hetzner_deploy.pub deploy@<hetzner-ip>
```

### 3. Add Private Key to GitHub Secrets

```bash
# Copy the private key content
cat ~/.ssh/hetzner_deploy
```

Then in GitHub: **Settings → Secrets and variables → Actions**

Add the following secret:
- `HETZNER_DEPLOY_KEY` = *contents of private key*

---

## DuckDNS Domain Setup

### What is DuckDNS?

DuckDNS is a free dynamic DNS service that allows you to:
- Get a free `.duckdns.org` subdomain
- Automatically update your IP address if it changes
- Access your server by domain name instead of IP address
- Use Let's Encrypt SSL certificates with your domain

### 1. Create DuckDNS Account

1. Go to https://www.duck.dns.org/
2. Sign in with any OAuth provider (GitHub, Google, etc.)
3. Note your **DuckDNS Token** (you'll need it later)

### 2. Create Your Domain

1. On the DuckDNS dashboard, enter your desired subdomain name (e.g., `travel-nugget`)
2. Click "Add Domain"
3. The domain `travel-nugget.duckdns.org` is now registered to your account

### 3. Update Your IP Address

Option A: **Manual Update** (once)

```bash
# Replace YOUR_TOKEN with your DuckDNS token from step 1
curl "https://www.duck.dns.org/update?domains=travel-nugget&token=YOUR_TOKEN&ip="
```

Option B: **Automatic Updates** (recommended) - Add to crontab

```bash
# Open crontab editor
crontab -e

# Add this line to update every 5 minutes:
*/5 * * * * curl "https://www.duck.dns.org/update?domains=travel-nugget&token=YOUR_TOKEN&ip=" > /dev/null 2>&1

# Save and exit (Ctrl+X, then Y, then Enter)
```

### 4. Verify Domain Resolution

```bash
# Test that your domain resolves to your server's IP
nslookup travel-nugget.duckdns.org

# Should output something like:
# Server:         8.8.8.8
# Address:        8.8.8.8#53
# Non-authoritative answer:
# Name:   travel-nugget.duckdns.org
# Address: 123.45.67.89  (your server IP)
```

### 5. Update GitHub Secrets

Once your domain is working, update GitHub repository secrets:

```bash
DOMAIN=travel-nugget.duckdns.org
NEXT_PUBLIC_API_URL=https://travel-nugget.duckdns.org/api
```

### 6. Deploy with Your Domain

The deployment will now:
- Use your DuckDNS domain
- Automatically provision Let's Encrypt SSL certificate
- Serve HTTPS with automatic redirect from HTTP

### What Happens If IP Changes?

DuckDNS will continue to work because:
1. Your server updates the IP address every 5 minutes via cron job
2. DuckDNS DNS records update to point to new IP
3. Your application continues to work with zero downtime
4. SSL certificate remains valid (it's tied to the domain, not IP)

### Troubleshooting DuckDNS

**Domain not resolving?**
```bash
# Check current IP on server
curl https://checkip.amazonaws.com

# Manually trigger update
curl "https://www.duck.dns.org/update?domains=travel-nugget&token=YOUR_TOKEN&ip="

# Wait 5-10 minutes for DNS propagation
nslookup travel-nugget.duckdns.org
```

**Want to change your domain?**
1. Delete current domain from DuckDNS dashboard
2. Create new domain
3. Update GitHub secrets with new domain
4. Redeploy

---

## Server Setup

### Prerequisites on Hetzner Server

Run these commands as root or with sudo:

```bash
# 1. Create project directory
sudo mkdir -p /app/cloud/logs
sudo chown -R deploy:deploy /app/cloud

# 2. Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker deploy

# 3. Install Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# 4. Install Caddy (Reverse Proxy + HTTPS)
sudo apt-get update && sudo apt-get install -y debian-keyring debian-archive-keyring apt-transport-https
curl https://dl.filippo.io/caddy/stable?plugins=dns.providers.duckdns -o /tmp/caddy
sudo mv /tmp/caddy /usr/local/bin/caddy
sudo chmod +x /usr/local/bin/caddy

# 5. Setup Duck DNS (for dynamic IP tracking)
# Complete guide: See [DuckDNS Domain Setup](#duckdns-domain-setup) section above
# Get your token from https://www.duck.dns.org/
crontab -e
# Add the following line:
*/5 * * * * curl "https://www.duck.dns.org/update?domains=yoursubdomain&token=DUCK_TOKEN&ip=" > /dev/null 2>&1

# 6. GitHub Container Registry Login
docker login ghcr.io -u <your-github-username> -p <your-pat>
```

---

## GitHub Secrets Configuration

Configure the following secrets in your GitHub repository settings:

### Build & Push Secrets

- `REGISTRY` = `ghcr.io`
- `NEXT_PUBLIC_API_URL` = `https://yoursubdomain.duck.dns/api`

### Container Registry Credentials

- `GHCR_USERNAME` = Your GitHub username
- `GHCR_PAT` = Personal Access Token with `read:packages` scope

### Database Secrets

- `DB_USER` = `trip_user`
- `DB_PASSWORD` = *secure random password*
- `DB_NAME` = `trip_manager`

### Application Secrets

- `JWT_SECRET` = *random secure string for JWT signing*
- `DOMAIN` = `yoursubdomain.duck.dns`

### Server Connection Secrets

- `HETZNER_HOST` = Your Hetzner server IP
- `HETZNER_USER` = `deploy` (or your deployment user)
- `HETZNER_PORT` = SSH port (usually `22`)
- `HETZNER_DEPLOY_KEY` = *private SSH key content*

---

## Automated Deployment (GitHub Actions)

### Deployment Flow

```
1. Merge PR to main/dev branch
   ↓
2. GitHub Actions: "Build & Push Docker Images"
   - Backend image built and pushed to ghcr.io
   - Frontend image built and pushed to ghcr.io
   ↓
3. GitHub Actions: "Deploy" job starts automatically
   ↓
4. SSH connects to /app/cloud on Hetzner
   ↓
5. Server deployment script (deploy.sh) executes:
   - Environment file generated with secret substitution
   - Caddyfile generated with domain substitution
   - docker-compose pulls latest images
   - Old containers stopped
   - New containers started
   - Health checks performed
   ↓
6. Server running with latest containers ✅
   - Frontend: https://yoursubdomain.duck.dns
   - API: https://yoursubdomain.duck.dns/api
   - Caddy: Reverse proxy with automatic HTTPS
```

### Trigger a Deployment

Simply merge a PR to the main branch. The pipeline automatically:
1. Detects changes in backend/frontend
2. Builds and pushes Docker images (if changed)
3. Deploys to Hetzner (if images were built)

---

## Manual Deployment

For manual deployments without GitHub Actions, use the provided script.

### Prerequisites

```bash
# 1. Ensure SSH key is set up (see SSH Key Setup section)
# 2. Navigate to the deployment directory
cd deploy/hetzner

# 3. Create .env file with all required variables
cat > .env << 'EOF'
GHCR_USERNAME=your-github-username
GHCR_PAT=your-personal-access-token
DB_USER=trip_user
DB_PASSWORD=your-secure-password
DB_NAME=trip_manager
JWT_SECRET=your-jwt-secret
DOMAIN=yoursubdomain.duck.dns
EOF

# 4. Ensure the script is executable
chmod +x manual-deploy.sh
```

### Run Manual Deployment

```bash
# Basic usage
./manual-deploy.sh <server-ip> <username> <port>

# Or via config
./manual-deploy.sh <config-entry-name>

# Example
./manual-deploy.sh 192.168.1.100 deploy 22

# With SSH key
./manual-deploy.sh 192.168.1.100 deploy 22
```

### What the Script Does

1. **Validates environment**
   - Checks for required .env file
   - Verifies all required variables are set

2. **Creates remote directory**
   - SSH connects to server
   - Creates `/home/USERNAME/app/cloud` if needed

3. **Copies deployment files**
   - `deploy.sh` script
   - `docker-compose.yaml`
   - `Caddyfile` (with variable substitution)

4. **Generates .env on server**
   - Substitutes all secrets into remote .env file

5. **Executes deployment**
   - Runs deploy.sh on remote server
   - Starts containers with docker-compose

6. **Verifies deployment**
   - Checks docker-compose status
   - Confirms all services are running

### Example Full Workflow

```bash
cd deploy/hetzner

# Create .env from template
cat > .env << 'EOF'
GHCR_USERNAME=octocat
GHCR_PAT=ghp_xxxxxxxxxxxxxxxxxxxx
DB_USER=trip_user
DB_PASSWORD=MySecurePassword123!
DB_NAME=trip_manager
JWT_SECRET=your-jwt-secret-key-here
DOMAIN=trips.duck.dns
EOF

# Deploy to server
./manual-deploy.sh 203.0.113.42 deploy 22

# Output will show progress:
# [2026-04-01 14:21:05] 📂 Loading environment variables...
# [2026-04-01 14:21:06] 🔐 Connecting to server 203.0.113.42 (deploy) on port 22...
# [2026-04-01 14:21:07] 📤 Copying deploy files to server...
# [2026-04-01 14:21:08] ✅ Deploy files copied
# [2026-04-01 14:21:09] 📤 Generating and copying Caddyfile to server...
# [2026-04-01 14:21:10] ✅ Caddyfile copied
# [2026-04-01 14:21:11] 📝 Generating .env file on server...
# [2026-04-01 14:21:12] ✅ .env file generated
# [2026-04-01 14:21:13] 🚀 Executing deploy script on server...
# [2026-04-01 14:21:25] ✅ Deployment completed successfully!
```

---

## Architecture Overview

### Components

```
┌─────────────────────────────────────────────────┐
│         GitHub Actions Pipeline                 │
├─────────────────────────────────────────────────┤
│ 1. Detect Changes                               │
│    ↓                                             │
│ 2. Build & Push Docker Images (if changed)      │
│    - Backend → ghcr.io/.../backend:tag          │
│    - Frontend → ghcr.io/.../frontend:tag        │
│    ↓                                             │
│ 3. Deploy to Hetzner (if images built)          │
│    ↓                                             │
│ 4. Manual Deploy Script (alternative)           │
└─────────────────────────────────────────────────┘
                     ↓
┌─────────────────────────────────────────────────┐
│      Hetzner Server (/app/cloud)                │
├─────────────────────────────────────────────────┤
│ docker-compose:                                 │
│                                                 │
│  ┌──────────────┐  ┌──────────────┐             │
│  │  Caddy       │  │  Frontend    │             │
│  │  (Reverse    │→→│  (Next.js)   │             │
│  │   Proxy)     │  │  :3000       │             │
│  │  :80/443     │  └──────────────┘             │
│  └──────────────┘                               │
│       ↓                                          │
│  ┌──────────────┐  ┌──────────────┐             │
│  │  Backend     │  │  Database    │             │
│  │  (Go API)    │→→│  (PostgreSQL)│             │
│  │  :8000       │  │  :5432       │             │
│  └──────────────┘  └──────────────┘             │
└─────────────────────────────────────────────────┘
```

### Caddy Reverse Proxy Configuration

The Caddyfile configures:
- ✅ HTTPS with Let's Encrypt (automatic)
- ✅ `/api/*` → Backend (port 8000)
- ✅ `/` → Frontend (port 3000)
- ✅ Security headers
- ✅ Gzip compression
- ✅ Access logging
- See: [Caddyfile](../deploy/hetzner/deployment/Caddyfile) for full configuration

### Networking

- **Public**: Only Caddy (port 80, 443)
- **Internal (Docker network)**: Backend, Frontend, Database communicate via Docker DNS
- **Database**: Only accessible from Backend container

---

## Health Checks

The deployment process performs automatic health checks:

### Backend Health Check
```bash
curl http://backend:8000/health
```

### Frontend Health Check
```bash
curl http://frontend:3000
```

These checks run from within the Docker network and verify that:
- All services are reachable
- Services respond to basic requests
- Database connections are working

---

## Logs

### Deployment Logs

Manual deployment logs are printed to stdout with timestamps.

### Container Logs

Access container logs on the server:

```bash
cd /app/cloud

# All containers
docker-compose logs

# Specific service
docker-compose logs backend
docker-compose logs frontend
docker-compose logs database

# Follow logs in real-time
docker-compose logs -f backend

# Last N lines
docker-compose logs --tail 50 backend
```

### Server Logs

Check system logs for SSH and deployment issues:

```bash
# SSH related
sudo journalctl -u ssh -n 50

# Docker related
sudo journalctl -u docker -n 50

# System messages
sudo tail -f /var/log/syslog
```

---

## Troubleshooting

### Issue: "manifest unknown" when pulling images

**Cause**: Images not built or not pushed to registry.

**Solution**:
```bash
# Check if images were built
docker images | grep trip-manager

# Manually trigger GitHub Actions build
# Or rebuild manually:
cd backend && docker build -t ghcr.io/username/trip-manager/backend:latest .
docker push ghcr.io/username/trip-manager/backend:latest
```

### Issue: SSH connection refused

**Cause**: SSH key not configured or incorrect port.

**Solution**:
```bash
# Test SSH connection
ssh -i ~/.ssh/hetzner_deploy -p 22 deploy@<server-ip> "echo OK"

# Check server SSH service
sudo systemctl status ssh

# Verify key permissions
chmod 600 ~/.ssh/hetzner_deploy
chmod 644 ~/.ssh/hetzner_deploy.pub
```

### Issue: Containers won't start - database connection fails

**Cause**: Database not initialized or password mismatch.

**Solution**:
```bash
# Check environment variables
docker-compose config | grep DB_

# Verify database is running
docker-compose ps database

# Check database logs
docker-compose logs database

# Reset database (WARNING: deletes data)
docker-compose down -v
docker-compose up -d database
sleep 10
docker-compose up -d
```

### Issue: Domain not resolving

**Cause**: Duck DNS not updated or DNS propagation delay.

**Solution**:
```bash
# Test DNS resolution
nslookup yoursubdomain.duck.dns

# Manually update Duck DNS
curl "https://www.duck.dns.org/update?domains=yoursubdomain&token=YOUR_TOKEN&ip="

# Check Caddy configuration
docker exec caddy caddy list-modules
docker logs $(docker ps | grep caddy | awk '{print $1}')
```

### Issue: pnpm command not found during build

**Cause**: Monorepo structure not properly configured.

**Solution**:
- Ensure `pnpm-workspace.yaml` exists in root
- Check `package.json` has proper scripts in backend/ and frontend/
- Verify `pnpm.lock` is up to date: `pnpm install`

---

## Rollback

### Quick Rollback to Previous Containers

```bash
cd /app/cloud

# Stop current containers
docker-compose down

# Pull and start previous images
# (Docker will use locally cached images if available)
docker-compose up -d

# Or manually specify previous tag
# Edit docker-compose.yaml to point to previous tag, then:
docker-compose up -d
```

### Full Rollback with Version Tags

```bash
cd /app/cloud

# Check available image tags
docker images | grep trip-manager

# Update docker-compose.yaml to use specific tag
# Example: backend:sha-abc123def456 instead of backend:latest

# Redeploy
docker-compose down
docker-compose pull
docker-compose up -d
```

### Database Rollback

```bash
# Backup current database
docker exec trip-manager-database-1 pg_dump -U trip_user trip_manager > backup_$(date +%Y%m%d_%H%M%S).sql

# Connect to database
docker exec -it trip-manager-database-1 psql -U trip_user -d trip_manager

# Inside psql:
-- List tables
\dt

-- Check migrations
SELECT version FROM schema_migrations ORDER BY version DESC;

-- Rollback specific migration (if using migration system)
```

---

## Environment File (.env)

### Required Variables

```bash
# GitHub Container Registry
GHCR_USERNAME=your-github-username
GHCR_PAT=ghp_xxxxxxxxxxxxxxxxxxxx

# Database
DB_USER=trip_user
DB_PASSWORD=secure_password_here
DB_NAME=trip_manager

# Application
JWT_SECRET=your-jwt-secret-key
DOMAIN=yoursubdomain.duck.dns

# Optional
NEXT_PUBLIC_API_URL=https://yoursubdomain.duck.dns/api
```

### Note

- **In GitHub Actions**: Secrets are automatically substituted
- **In Manual Deployment**: Create `.env` file in `deploy/hetzner/` directory
- **Never commit `.env` to git**: It contains sensitive data
- **Always use strong passwords**: Especially `DB_PASSWORD` and `JWT_SECRET`

---

## Quick Reference Commands

### Server Access
```bash
ssh -p 22 deploy@<server-ip>
cd /app/cloud
```

### Container Management
```bash
docker-compose ps              # Show running containers
docker-compose logs -f         # Follow all logs
docker-compose restart         # Restart all services
docker-compose pull            # Update images
docker-compose down            # Stop all services
```

### Database Access
```bash
docker exec -it <db-container> psql -U trip_user -d trip_manager
```

### Check Disk Space
```bash
df -h                          # Overall space
docker system df               # Docker storage
```

---

👈 **[Back to README](../README.md)**

---

## Related Documentation

- [CI/CD Pipeline](./CI_CD.md) - GitHub Actions workflow details
- [Setup Guide](./SETUP.md) - Initial project setup
- [Makefile Reference](./MAKEFILE.md) - Available make commands


# Manual Deployment Guide

## Overview

The `manual-deploy.sh` script automates the deployment of the Trip Manager application to a remote server. It performs the same steps as the GitHub Actions pipeline, but can be run manually from your local machine.

## Prerequisites

1. **SSH Access**: Configured SSH key for passwordless login to the server
2. **.env File**: Located in this directory with required variables (see `.env.example`)
3. **Docker Compose**: Available on the remote server
4. **Required Environment Variables**:
   - `GHCR_USERNAME` - GitHub Container Registry username
   - `GHCR_PAT` - GitHub Personal Access Token
   - `DB_USER` - Database user
   - `DB_PASSWORD` - Database password
   - `DB_NAME` - Database name
   - `JWT_SECRET` - JWT secret for authentication
   - `DOMAIN` - Domain name for the application

## Usage

### Option 1: SSH Config Alias (Recommended)

If you have an SSH config entry, simply use the alias:

```bash
./manual-deploy.sh myserver
```

Or with port override:
```bash
./manual-deploy.sh myserver 2222
```

### Option 2: Direct IP with User and Port

```bash
./manual-deploy.sh deploy@192.168.1.100 22
./manual-deploy.sh deploy@167.235.66.0 2222
```

### Option 3: Just IP (Uses Default User)

```bash
./manual-deploy.sh 192.168.1.100
```

This will try to extract the user from your SSH config, falling back to `deploy` user.

## SSH Configuration

### Setting Up SSH Config (Recommended)

Create or edit `~/.ssh/config`:

```
Host myserver
    HostName 192.168.1.100
    User deploy
    Port 22
    IdentityFile ~/.ssh/id_rsa
    StrictHostKeyChecking no
    UserKnownHostsFile /dev/null
```

Then deploy with:
```bash
./manual-deploy.sh myserver
```

**⚠️ Important**: Make sure the `IdentityFile` matches your actual SSH key path. If you don't have `id_rsa`, update it to your key name:

```
Host myserver
    HostName 192.168.1.100
    User deploy
    Port 22
    IdentityFile ~/.ssh/deploy_key    # Use your actual key name!
    StrictHostKeyChecking no
```

### Hetzner Example

```
Host hetzner
    HostName 167.235.66.0
    User deploy
    Port 2222
    IdentityFile ~/.ssh/github_deploy_key    # Update to your key name
    StrictHostKeyChecking no
```

Deploy:
```bash
./manual-deploy.sh hetzner
```

### Why SSH Config is Better

The script will automatically:
1. Read your SSH config
2. Extract the correct `IdentityFile` (SSH key)
3. Extract the port and user
4. Handle all authentication details

This is more reliable than passing parameters, especially if your SSH key is not named `id_rsa`.

## What the Script Does

1. **Validates environment** - Checks for .env file and required variables
2. **Creates remote directory** - Creates `/home/$USERNAME/app/cloud` on server
3. **Copies deploy files** - Transfers `deploy.sh` and `docker-compose.yaml`
4. **Generates Caddyfile** - Substitutes environment variables in Caddyfile
5. **Creates .env on server** - Generates environment file for Docker Compose
6. **Executes deployment** - Runs the deploy script on remote server
7. **Verifies deployment** - Checks Docker Compose status

## Troubleshooting

### "Permission denied (publickey)"

**Problem**: SSH key not found or not configured on server

**Solutions**:

**1. Add SSH Config Entry** (Recommended):
```
Host hetzner
    HostName 167.235.66.0
    User deploy
    Port 2222
    IdentityFile ~/.ssh/github_deploy_key   # Use your actual key name!
    StrictHostKeyChecking no
```

Then use:
```bash
./manual-deploy.sh hetzner
```

**2. If Key is Not on Server Yet**:
```bash
# Add your public key to server (replace with your key)
ssh-copy-id -i ~/.ssh/github_deploy_key -p 2222 deploy@167.235.66.0

# Or manually
cat ~/.ssh/github_deploy_key.pub | ssh -p 2222 deploy@167.235.66.0 \
  "mkdir -p ~/.ssh && cat >> ~/.ssh/authorized_keys && chmod 600 ~/.ssh/authorized_keys"

# Then test
ssh -p 2222 deploy@167.235.66.0 "echo Success"
```

**3. Debug Which Key SSH is Using**:
```bash
ssh -v deploy@167.235.66.0 -p 2222 2>&1 | grep -i "identity\|offering"
```

### "Missing required environment variable"

**Problem**: .env file missing or incomplete variables

**Solution**:
```bash
# Copy example and fill in your values
cp .env.example .env
# Edit .env with your actual values
nano .env
```

### "Failed to create remote directory"

**Problem**: SSH connection issue or permission denied on server

**Solution**:
```bash
# Test SSH connection
ssh deploy@167.235.66.0 -p 2222 "mkdir -p /home/deploy/app/cloud"

# If that fails, check permissions on server
ssh deploy@167.235.66.0 "ls -la /home/deploy/"
```

### Verbose SSH Debugging

```bash
# Run with verbose SSH output to see what's happening
ssh -vvv deploy@167.235.66.0 -p 2222 "echo test"
```

## Environment File Example

Create `.env` in this directory:

```bash
GHCR_USERNAME=your-github-username
GHCR_PAT=ghp_xxxxxxxxxxxxxxxxxxxx
DB_USER=postgres
DB_PASSWORD=your-secure-password
DB_NAME=trip_manager
JWT_SECRET=your-jwt-secret-key
DOMAIN=yourdomain.com
```

## Testing Deployment

Before full deployment, test SSH connection:

```bash
# Test connection
./manual-deploy.sh myserver  # Will show what it's doing

# Or just test SSH
ssh myserver "docker ps"
```

## Advanced: Port Forwarding

If your server is behind a firewall or NAT:

```bash
# Add to SSH config for port forwarding
Host myserver
    HostName 192.168.1.100
    User deploy
    Port 2222
    ProxyJump jumphost  # If using jump server
```

## Support

For issues, check:
1. SSH connectivity: `ssh myserver "echo test"`
2. .env file exists and is complete: `cat .env`
3. Docker on server: `ssh myserver "docker --version"`
4. Server logs: `ssh myserver "cd /home/deploy/app/cloud && docker-compose logs"`

## See Also

- `deployment/deploy.sh` - Remote deployment script
- `deployment/Caddyfile` - Reverse proxy configuration
- `../../docker-compose.yaml` - Docker Compose configuration
- `../../backend/README.md` - Backend documentation


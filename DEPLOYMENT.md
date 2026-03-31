# Deployment Setup

## SSH Key für Hetzner

### 1. SSH Key auf deinem Rechner generieren (falls nicht vorhanden)
```bash
ssh-keygen -t ed25519 -f ~/.ssh/hetzner_deploy -C "github-deploy@trip-manager"
```

### 2. Public Key auf Hetzner Server kopieren
```bash
ssh-copy-id -i ~/.ssh/hetzner_deploy.pub deploy@<hetzner-ip>
```

### 3. Private Key als GitHub Secret hinterlegen
```bash
# Inhalt des Private Keys kopieren
cat ~/.ssh/hetzner_deploy
```

Dann in GitHub: **Settings → Secrets and variables → Actions**

Neue Secrets:
```
HETZNER_HOST = <deine-hetzner-ip>
HETZNER_USER = deploy
HETZNER_DEPLOY_KEY = <inhalt-des-privaten-keys>
```

---

## Server Setup

### Auf deinem Hetzner Server:

```bash
# 1. Deploy User erstellen (falls nicht existent)
sudo useradd -m -s /bin/bash deploy

# 2. Verzeichnis erstellen
sudo mkdir -p /home/deploy/trip-manager/logs
sudo chown -R deploy:deploy /home/deploy/trip-manager

# 3. Docker installieren
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker deploy

# 4. Docker Compose installieren
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# 5. GitHub Container Registry Login
# Benötigt Personal Access Token (PAT) mit read:packages Scope
docker login ghcr.io -u <your-username> -p <your-pat>
```

---

## Deployment Flow

```
1. PR zu main/dev wird gemerged
   ↓
2. "Build & Push Docker Images" Workflow läuft
   - Baut Backend/Frontend Images
   - Pusht zu ghcr.io
   ↓
3. "Deploy to Hetzner" Workflow startet
   ↓
4. Change Detection:
   - Ist etwas in deploy/ geändert?
   - Ist Branch main?
   ↓
5. Falls ja → Deploy:
   - Kopiert deploy.sh zu Server
   - Führt deploy.sh aus
   - Script: docker-compose pull → down → up → health check
   ↓
6. Server läuft mit neuen Containern ✅
```

---

## Health Checks

Das Deploy Script checkt automatisch:
- Backend: `http://localhost:8000/health`
- Frontend: `http://localhost:3000` (optional)

Falls Health Check fehlschlägt:
- Deployment rollback
- Logs werden gesammelt
- GitHub Notification

---

## Rollback (manuell)

Falls was schiefgeht, auf dem Server:

```bash
cd /home/deploy/trip-manager
docker-compose down
docker-compose up -d --pull always
```

---

## Logs

Deploy Logs auf dem Server:
```bash
/home/deploy/trip-manager/logs/deploy-YYYYMMDD_HHMMSS.log
```

Container Logs:
```bash
cd /home/deploy/trip-manager
docker-compose logs backend
docker-compose logs frontend
```

---

## Troubleshooting

### ❌ "Permission denied" beim SSH
- Private Key Permissions: `chmod 600 ~/.ssh/hetzner_deploy`
- Public Key auf Server: `~/.ssh/authorized_keys`

### ❌ "docker: command not found"
- Docker Installation checken
- `docker ps` sollte funktionieren

### ❌ "ghcr.io: unauthorized"
- PAT Token in GitHub Secrets?
- Ist es noch gültig?

### ❌ Health Check schlägt fehl
- Logs checken: `docker-compose logs backend`
- Container startet? `docker-compose ps`


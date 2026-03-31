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

Neue Secrets (für automated Deployment):
```
HETZNER_HOST = <deine-hetzner-ip>
HETZNER_USER = deploy
HETZNER_DEPLOY_KEY = <inhalt-des-privaten-keys>

# Diese werden in die .env auf dem Server substituiert:
GHCR_USERNAME = <your-github-username>
GHCR_PAT = ghp_xxxxxxxxxxxx (read:packages scope!)
DB_USER = trip_user
DB_PASSWORD = <sichere-passwort>
DB_NAME = trip_manager
JWT_SECRET = <jwt-secret-key>
DOMAIN = yoursubdomain.duck.dns
```

---

## Server Setup

### Auf deinem Hetzner Server:

```bash
# 1. Projekt Verzeichnis erstellen
sudo mkdir -p /app/cloud/logs
sudo chown -R deploy:deploy /app/cloud

# 2. Docker installieren
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker deploy

# 3. Docker Compose installieren
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# 4. Caddy installieren (Reverse Proxy + HTTPS)
sudo apt-get update && sudo apt-get install -y debian-keyring debian-archive-keyring apt-transport-https
curl https://dl.filippo.io/caddy/stable?plugins=dns.providers.duckdns -o /tmp/caddy
sudo mv /tmp/caddy /usr/local/bin/caddy
sudo chmod +x /usr/local/bin/caddy

# 5. GitHub Container Registry Login
docker login ghcr.io -u <your-username> -p <your-pat>

# 6. Duck DNS aktualisieren (dynamic IP tracking)
# Von https://www.duck.dns.org/ Token besorgen
crontab -e
# Hinzufügen:
*/5 * * * * curl "https://www.duck.dns.org/update?domains=yoursubdomain&token=DUCK_TOKEN&ip=" > /dev/null 2>&1
```

---

## Environment Setup

Die `.env` wird **automatisch von der GitHub Pipeline** generiert! 🎉

Du brauchst **keine manuelle `.env` auf dem Server zu erstellen**.

Stattdessen:
1. GitHub Secrets setzen (siehe oben)
2. GitHub Actions ersetzt die Variablen automatisch
3. `.env` wird beim Deploy hochgeladen und substituiert

**Alte Anleitung (nicht mehr nötig):**
```bash
# NICHT MEHR NÖTIG:
cd /app/cloud
cp .env.example .env
nano .env  # nicht mehr nötig
```

---

## Deployment Flow

```
1. PR zu main/dev wird gemerged
   ↓
2. "Build & Push Docker Images" läuft
   - Baut Backend/Frontend Images
   - Pusht zu ghcr.io
   ↓
3. "Deploy" startet automatisch
   ↓
4. SSH zu /app/cloud auf Hetzner
   ↓
5. deploy.sh ausführen:
   - docker-compose pull (neue Images)
   - docker-compose down (alte Container stoppen)
   - docker-compose up -d (neue Container starten)
   - Health Checks
   - Caddy reload
   ↓
6. Server mit neuen Containern ✅
   - Frontend: https://yoursubdomain.duck.dns
   - API: https://yoursubdomain.duck.dns/api
   - Caddy: Reverse Proxy mit SSL
```

---

## Health Checks

Das Deploy Script checkt:
- Backend: `http://backend:8000/health` (Docker Netzwerk)
- Frontend: `http://frontend:3000` (Docker Netzwerk)

---

## Caddy Reverse Proxy

Die `Caddyfile` konfiguriert:
- ✅ HTTPS mit Let's Encrypt (automatisch)
- ✅ `/api/*` → Backend (8000)
- ✅ `/` → Frontend (3000)
- ✅ Security Headers
- ✅ Gzip Compression
- ✅ Logging

---

## Logs

Deploy Logs:
```bash
/app/cloud/logs/deploy-YYYYMMDD_HHMMSS.log
```

Container Logs:
```bash
cd /app/cloud
docker-compose logs backend
docker-compose logs frontend
docker-compose logs database
```

---

## Rollback

Falls etwas schiefgeht:
```bash
cd /app/cloud
docker-compose down
docker-compose up -d --pull always
```


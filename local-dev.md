# Trip Manager – Local Dev Guide

## Voraussetzungen

- Minikube
- Helm + Helmfile
- kubectl
- Docker
- Firebase CLI
- pnpm

---

## 1. Cluster starten

```bash
minikube start --cpus=4 --memory=8192
minikube addons enable ingress
```

---

## 2. Services bauen

```bash
eval $(minikube docker-env)
./scripts/build-dev.sh
```

---

## 3. Dev-Umgebung deployen

```bash
helmfile -e dev sync
```

---

## 4. Minikube Tunnel (für Ingress-Zugriff)

In einem separaten Terminal offen lassen:

```bash
minikube tunnel
```

Ingress-IP herausfinden:

```bash
kubectl get ingress -n trip-manager-dev
```

Diese IP in `frontend/.env.local` eintragen:

```
NEXT_PUBLIC_API_URL=http://<INGRESS-IP>
```

---

## 5. Firebase Emulator Port-Forward

In einem separaten Terminal offen lassen:

```bash
kubectl port-forward -n trip-manager-dev \
  $(kubectl get pod -n trip-manager-dev -l app.kubernetes.io/name=firebase-emulator -o jsonpath='{.items[0].metadata.name}') \
  9099:9099 4000:4000 8080:8080
```

---

## 6. Frontend starten

```bash
cd frontend
pnpm dev
```

Dann im Browser: `http://localhost:3000`

---

## 7. Prometheus Port-Forward (optional, für Metriken)

```bash
kubectl port-forward -n trip-manager-dev svc/prometheus-server 9090:80 &
```

Dann: `http://localhost:9090`

---

## 8. Pods prüfen

```bash
kubectl get pods -n trip-manager-dev
```

---

## Datenbank-Cleanup

### Users + Tenants zurücksetzen

```bash
kubectl exec -n trip-manager-dev users-postgres-0 -- psql -U users -d users -c "
  DELETE FROM users WHERE tenant_id != 'default';
  DELETE FROM tenants WHERE id != 'default';
  UPDATE users SET tenant_id = 'default', role = 'tenant_member' WHERE tenant_id = 'default';
"
```

### Nur Tenants zurücksetzen (User behalten)

```bash
kubectl exec -n trip-manager-dev users-postgres-0 -- psql -U users -d users -c "
  UPDATE users SET tenant_id = 'default', role = 'tenant_member';
  DELETE FROM tenants WHERE id != 'default';
"
```

### Tenant-Settings zurücksetzen

```bash
kubectl exec -n trip-manager-dev users-postgres-0 -- psql -U users -d users -c "
  UPDATE tenants SET settings = '{\"maxActiveTrips\": 3}' WHERE tier = 'free';
  UPDATE tenants SET settings = '{\"maxActiveTrips\": 0}' WHERE tier != 'free';
"
```

### Trips löschen

```bash
kubectl exec -n trip-manager-dev trips-postgres-0 -- psql -U trips -d trips -c "
  TRUNCATE TABLE trips CASCADE;
"
```

### Alles zurücksetzen (kompletter Neustart)

```bash
# Users + Tenants
kubectl exec -n trip-manager-dev users-postgres-0 -- psql -U users -d users -c "
  TRUNCATE TABLE users CASCADE;
  DELETE FROM tenants WHERE id != 'default';
"

# Trips
kubectl exec -n trip-manager-dev trips-postgres-0 -- psql -U trips -d trips -c "
  TRUNCATE TABLE trips CASCADE;
"

# Firebase Emulator neu starten (löscht alle Auth-User)
kubectl rollout restart deployment firebase-emulator -n trip-manager-dev
```

---

## Nützliche Befehle

### Logs eines Services anschauen

```bash
kubectl logs -n trip-manager-dev \
  $(kubectl get pod -n trip-manager-dev -l app.kubernetes.io/name=<service> -o jsonpath='{.items[0].metadata.name}') -f
```

Ersetze `<service>` mit z.B. `trips`, `users`, `social`, `feed`, `auth`, `travel-info`, `newsletter`.

### Service neu deployen (nach Code-Änderung)

```bash
eval $(minikube docker-env)
docker build -f backend/<service>/Dockerfile -t localhost/trip-manager/backend/<service>:dev backend/
kubectl rollout restart deployment <service> -n trip-manager-dev
```

### Helmfile einzelnen Release neu deployen

```bash
helmfile -e dev -l name=<service> sync
```

### Alten Release entfernen (bei Umbenennung)

```bash
helm uninstall <alter-release-name> -n trip-manager-dev
```

### Namespace komplett zurücksetzen

```bash
kubectl delete namespace trip-manager-dev
helmfile -e dev sync
```

---

## Kontext wechseln (Minikube ↔ GKE)

```bash
# Zu Minikube wechseln
kubectl config use-context minikube

# Zu GKE wechseln
kubectl config use-context gke_project-32c60644-299b-4b05-8cf_europe-west1_trip-manager-prod

# GKE Credentials neu holen (falls abgelaufen)
gcloud container clusters get-credentials trip-manager-prod \
  --region europe-west1 \
  --project project-32c60644-299b-4b05-8cf
```
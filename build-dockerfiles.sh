#!/bin/bash
set -e

TAG="${1:-dev}"
REGISTRY="localhost/trip-manager"

echo "🔨 Building backend services into minikube (tag: $TAG)"

# Minikube Docker Daemon verwenden
eval "$(minikube docker-env)"

services=$(find backend -mindepth 1 -maxdepth 1 -type d | xargs -I{} basename {})

for service in $services; do
  { [ "$service" = "shared" ] || [ "$service" = "nginx" ]; } && continue

  path="backend/$service"

  if [ ! -d "$path" ]; then
    echo "⏭️  Skipping: $service"
    continue
  fi

  if [ ! -f "$path/go.mod" ] || [ ! -f "$path/Dockerfile" ]; then
    echo "⏭️  Skipping $service (missing go.mod or Dockerfile)"
    continue
  fi

  IMAGE="$REGISTRY/backend/$service:$TAG"
  echo "🐳 Building $service → $IMAGE"

  docker build \
    -f "$path/Dockerfile" \
    -t "$IMAGE" \
    backend/

  echo "✅ $service done"
done

echo ""
echo "🎉 All services built into minikube"
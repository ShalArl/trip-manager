#!/bin/bash

DOCKERFILES_DIR="backend"

for service in "$DOCKERFILES_DIR"/*/; do
  service_path="${service%/}"
  service_name=$(basename "$service_path")

  { [ "$service_name" == "shared" ] || [ "$service_name" == "nginx" ]; } && continue

  dockerfile_path="$service_path/Dockerfile"

  if [ -f "$dockerfile_path" ]; then
    echo "--- Building Docker image for: $service_name ---"

    docker build -t "localhost/$service_name:latest" \
      -f "$dockerfile_path" \
      "$DOCKERFILES_DIR"
  else
    echo "Skipping $service_name: No Dockerfile found at $dockerfile_path"
  fi
done
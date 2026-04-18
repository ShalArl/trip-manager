#!/bin/bash
source ../.env
source lib.sh

echo "Setting up API service account and permissions"

ensure_service "run"
ensure_service "secretmanager"
ensure_service "cloudbuild"
ensure_service "sqladmin"
ensure_service "artifactregistry"
ensure_service "storage-api"


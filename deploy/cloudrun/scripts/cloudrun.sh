#!/bin/bash

PLACEHOLDER_IMAGE="us-docker.pkg.dev/cloudrun/container/hello"

deploy_backend_service() {
    local service_name=$1
    local region=$2

    local env_vars="^@@^"
    env_vars+="ENVIRONMENT=production"
    env_vars+="@@SERVER_PORT=8081"
    env_vars+="@@CORS_ALLOWED_ORIGINS=https://trip-manager-frontend-271566791555.europe-west3.run.app,https://trip-manager-frontend-rygwuplcya-ey.a.run.app"
    env_vars+="@@STORAGE_TYPE=gcs"
    env_vars+="@@GCS_BUCKET=${GCS_BUCKET}"
    env_vars+="@@GCS_SIGNER_SA=${SIGNED_URL_SA_EMAIL}"
    env_vars+="@@GCS_SIGNED_URL_TTL_SECONDS=900"
    env_vars+="@@FIREBASE_PROJECT_ID=${PROJECT_ID}"

    local secrets="DATABASE_URL=database-url:latest"
    secrets+=",JWT_SECRET=jwt-secret:latest"


    if gcloud run services describe "$service_name" \
         --region="$region" --project="$PROJECT_ID" &>/dev/null; then
        echo "Updating Cloud Run service $service_name..."
        gcloud run services update "$service_name" \
            --region="$region" \
            --project="$PROJECT_ID" \
            --service-account="$RUNTIME_SA_EMAIL" \
            --add-cloudsql-instances="$SQL_CONNECTION_NAME" \
            --set-env-vars="$env_vars" \
            --set-secrets="$secrets"
    else
        echo "Creating Cloud Run service $service_name (placeholder image)..."
        gcloud run deploy "$service_name" \
            --image="$PLACEHOLDER_IMAGE" \
            --region="$region" \
            --project="$PROJECT_ID" \
            --service-account="$RUNTIME_SA_EMAIL" \
            --add-cloudsql-instances="$SQL_CONNECTION_NAME" \
            --allow-unauthenticated \
            --port=8081 \
            --memory=512Mi \
            --set-env-vars="$env_vars" \
            --set-secrets="$secrets"
    fi
}

deploy_frontend_service() {
    local service_name=$1
    local region=$2

    if gcloud run services describe "$service_name" \
         --region="$region" --project="$PROJECT_ID" &>/dev/null; then
        echo "Cloud Run service $service_name already exists, skipping config update."
        # Env-Vars für Frontend sind im Image eingebaut (Build-Args), kein Runtime-Update nötig
    else
        echo "Creating Cloud Run service $service_name (placeholder image)..."
        gcloud run deploy "$service_name" \
            --image="$PLACEHOLDER_IMAGE" \
            --region="$region" \
            --project="$PROJECT_ID" \
            --allow-unauthenticated \
            --port=3000 \
            --memory=512Mi
    fi
}

get_service_url() {
    local service_name=$1
    local region=$2
    gcloud run services describe "$service_name" \
        --region="$region" \
        --project="$PROJECT_ID" \
        --format="value(status.url)" 2>/dev/null || echo ""
}
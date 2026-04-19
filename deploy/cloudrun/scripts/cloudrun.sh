#!/bin/bash

PLACEHOLDER_IMAGE="us-docker.pkg.dev/cloudrun/container/hello"

deploy_backend_service() {
    local service_name=$1
    local region=$2

    local env_vars="ENVIRONMENT=production"
    env_vars+=",SERVER_PORT=8081"

    if [ "$ENABLE_STORAGE" = "true" ]; then
        env_vars+=",STORAGE_TYPE=s3"
        env_vars+=",S3_ENDPOINT=https://storage.googleapis.com"
        env_vars+=",S3_BUCKET=${GCS_BUCKET}"
        env_vars+=",S3_REGION=auto"
        env_vars+=",S3_USE_SSL=true"
        env_vars+=",S3_PUBLIC_URL=https://storage.googleapis.com/${GCS_BUCKET}"
    else
        env_vars+=",STORAGE_TYPE=local"
        env_vars+=",UPLOAD_DIR=/tmp/uploads"
    fi

    local secrets="DATABASE_URL=database-url:latest"
    secrets+=",JWT_SECRET=jwt-secret:latest"

    if [ "$ENABLE_STORAGE" = "true" ]; then
        secrets+=",S3_ACCESS_KEY=s3-access-key:latest"
        secrets+=",S3_SECRET_KEY=s3-secret-key:latest"
    fi

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
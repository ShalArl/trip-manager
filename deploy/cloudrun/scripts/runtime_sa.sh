#!/bin/bash

setup_runtime_sa() {
    local sa_name="cloudrun-backend"
    local sa_email="${sa_name}@${PROJECT_ID}.iam.gserviceaccount.com"

    ensure_service_account "$sa_name" "Cloud Run Backend Runtime"

    add_iam_role "$sa_name" "roles/cloudsql.client"
    add_iam_role "$sa_name" "roles/secretmanager.secretAccessor"

    gcloud iam service-accounts add-iam-policy-binding "$sa_email" \
        --project="$PROJECT_ID" \
        --role="roles/iam.serviceAccountUser" \
        --member="serviceAccount:${DEPLOY_SA_EMAIL}"

    export RUNTIME_SA_EMAIL="$sa_email"
}
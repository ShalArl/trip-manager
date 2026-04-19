#!/bin/bash

setup_gcs_bucket() {
    local bucket=$1
    local region=$2

    if ! gcloud storage buckets describe "gs://${bucket}" --project="$PROJECT_ID" &>/dev/null; then
        echo "Creating GCS bucket $bucket..."
        gcloud storage buckets create "gs://${bucket}" \
            --project="$PROJECT_ID" \
            --location="$region" \
            --uniform-bucket-level-access
    else
        echo "GCS bucket $bucket already exists."
    fi
}

setup_storage_sa_and_hmac() {
    local sa_name="storage-app"
    local sa_email="${sa_name}@${PROJECT_ID}.iam.gserviceaccount.com"
    local bucket=$1

    ensure_service_account "$sa_name" "App Storage Access (S3-compat)"
    wait_for_service_account "$sa_email"

    gcloud storage buckets add-iam-policy-binding "gs://${bucket}" \
        --member="serviceAccount:${sa_email}" \
        --role="roles/storage.objectAdmin"

    if ! gcloud secrets describe "s3-access-key" --project="$PROJECT_ID" &>/dev/null; then
        echo "Creating HMAC key for ${sa_email}..."
        local hmac_output
        hmac_output=$(gcloud storage hmac create "$sa_email" \
            --project="$PROJECT_ID" --format=json)

        local hmac_access_id hmac_secret
        hmac_access_id=$(echo "$hmac_output" | jq -r '.metadata.accessId')
        hmac_secret=$(echo "$hmac_output" | jq -r '.secret')

        create_secret_if_missing "s3-access-key" "$hmac_access_id"
        create_secret_if_missing "s3-secret-key" "$hmac_secret"
        echo "HMAC credentials stored in Secret Manager."
    else
        echo "HMAC secrets already exist, skipping key creation."
    fi
}
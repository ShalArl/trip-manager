#!/bin/bash

setup_gcs_bucket() {
    local bucket=$1
    local region=$2

    if ! gcloud storage buckets describe "gs://${bucket}" --project="$PROJECT_ID" &>/dev/null; then
        echo "Creating GCS bucket $bucket..."
        gcloud storage buckets create "gs://${bucket}" \
            --project="$PROJECT_ID" \
            --location="$region" \
            --uniform-bucket-level-access \
            --public-access-prevention
    else
        echo "GCS bucket $bucket already exists."
    fi
}
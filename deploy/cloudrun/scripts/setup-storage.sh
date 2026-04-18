#!/bin/bash

source ../.env
source lib.sh

ensure_service_account "$STORAGE_SA_NAME" "Storage Ops Account"

add_iam_role "$STORAGE_SA_NAME" "roles/storage.objectAdmin"


if ! gcloud storage buckets describe gs://$BUCKET_NAME &>/dev/null; then
    echo "Creating Bucket..."
    gcloud storage buckets create gs://$BUCKET_NAME --location=$BUCKET_LOCATION
else
    echo "Bucket already exists."
fi
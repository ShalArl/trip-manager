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

setup_cors() {
  local GCS_BUCKET=$1

  if [[ -z "$GCS_BUCKET" ]]; then
    echo "Error: GCS_BUCKET parameter is required" >&2
    return 1
  fi

  local cors_file
  cors_file=$(mktemp) || return 1
  trap "rm -f '$cors_file'" RETURN

  cat >"$cors_file" <<EOF
  [
    {
      "origin": [
        "https://trip-manager-frontend-271566791555.europe-west3.run.app",
        "https://trip-manager-frontend-rygwuplcya-ey.a.run.app"
      ],
      "method": ["PUT", "GET"],
      "responseHeader": ["Content-Type"],
      "maxAgeSeconds": 3600
    }
  ]
EOF

  if ! gcloud storage buckets update "gs://${GCS_BUCKET}" --cors-file="$cors_file"; then
    echo "Error: Failed to update CORS configuration for bucket $GCS_BUCKET" >&2
    return 1
  fi
}
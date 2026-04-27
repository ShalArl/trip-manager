#!/bin/bash
set -euo pipefail

source .env

for f in lib/lib.sh lib/cloudsql.sh lib/artifactory.sh \
         lib/storage.sh lib/runtime_sa.sh lib/cloudrun.sh; do
    source "$f"
done

# gcloud auth login --quiet

# === Phase 1: Project & APIs ===
ensure_project "$PROJECT_ID"
ensure_services \
    artifactregistry.googleapis.com \
    run.googleapis.com \
    sqladmin.googleapis.com \
    storage.googleapis.com \
    secretmanager.googleapis.com \
    iam.googleapis.com \
    iamcredentials.googleapis.com \
    firebase.googleapis.com \
    identitytoolkit.googleapis.com \
    firestore.googleapis.com

PROJECT_NUMBER=$(gcloud projects describe "$PROJECT_ID" --format="value(projectNumber)")

# === Phase 2: WIF ===
ensure_wif "$WIF_POOL" "$WIF_PROVIDER"

# === Phase 3: Deploy Service Account ===
ensure_service_account "$DEPLOY_SA_NAME" "Deploy Service Account"
DEPLOY_SA_EMAIL="$DEPLOY_SA_NAME@$PROJECT_ID.iam.gserviceaccount.com"

gcloud iam service-accounts add-iam-policy-binding "$DEPLOY_SA_EMAIL" \
    --project="$PROJECT_ID" \
    --role="roles/iam.workloadIdentityUser" \
    --member="principalSet://iam.googleapis.com/projects/$PROJECT_NUMBER/locations/global/workloadIdentityPools/$WIF_POOL/attribute.repository/$GITHUB_REPO"

add_iam_role "$DEPLOY_SA_NAME" "roles/artifactregistry.writer"
add_iam_role "$DEPLOY_SA_NAME" "roles/run.admin"
add_iam_role "$DEPLOY_SA_NAME" "roles/iam.serviceAccountUser"

# === Phase 4: Artifact Registry ===
setup_artifact_registry

# === Phase 5: Cloud SQL + DB-Passwort ===
SQL_INSTANCE="${APP_NAME}-db"
SQL_DB_NAME="${DB_NAME:-tripmanager}"
SQL_DB_USER="${DB_USER:-app}"
SQL_REGION="${REGION:-europe-west3}"
SQL_TIER="db-f1-micro"

sql_setup_db_password   # setzt $DB_PASSWORD aus Secret Manager
sql_setup_instance "$SQL_INSTANCE" "$SQL_DB_NAME" "$SQL_DB_USER" "$DB_PASSWORD" "$SQL_REGION" "$SQL_TIER"
sql_create_db_user "$SQL_INSTANCE" "$SQL_DB_USER" "$DB_PASSWORD"
sql_create_database "$SQL_INSTANCE" "$SQL_DB_NAME"

SQL_CONNECTION_NAME="${PROJECT_ID}:${SQL_REGION}:${SQL_INSTANCE}"

# === Phase 6: Secrets ===
DATABASE_URL="postgres://${SQL_DB_USER}:${DB_PASSWORD}@/${SQL_DB_NAME}?host=/cloudsql/${SQL_CONNECTION_NAME}&sslmode=disable"
create_secret_if_missing "database-url" "$DATABASE_URL"

# === Phase 7: GCS Bucket ===
setup_gcs_bucket "$GCS_BUCKET" "$REGION"
ensure_service_account "signed-url-signer" "Service Account for signing GCS URLs"
SIGNED_URL_SA_EMAIL="signed-url-signer@$PROJECT_ID.iam.gserviceaccount.com"

gcloud storage buckets add-iam-policy-binding "gs://${GCS_BUCKET}" \
    --member="serviceAccount:${SIGNED_URL_SA_EMAIL}" \
    --role="roles/storage.objectAdmin"

setup_cors "$GCS_BUCKET"

# === Phase 8: Runtime SA ===
setup_runtime_sa   # exportiert $RUNTIME_SA_EMAIL

gcloud storage buckets add-iam-policy-binding "gs://${GCS_BUCKET}" \
    --member="serviceAccount:${RUNTIME_SA_EMAIL}" \
    --role="roles/storage.objectAdmin"

gcloud iam service-accounts add-iam-policy-binding "$SIGNED_URL_SA_EMAIL" \
    --project="$PROJECT_ID" \
    --member="serviceAccount:${RUNTIME_SA_EMAIL}" \
    --role="roles/iam.serviceAccountTokenCreator"


# === Phase 9: Setup Firebase ===
# firebase projects:addfirebase "$PROJECT_ID"

gcloud projects add-iam-policy-binding "$PROJECT_ID" \
    --member="serviceAccount:${RUNTIME_SA_EMAIL}" \
    --role="roles/firebaseauth.admin"

gcloud projects add-iam-policy-binding "$PROJECT_ID" \
    --member="serviceAccount:${RUNTIME_SA_EMAIL}" \
    --role="roles/datastore.user"


# === Phase 10: Cloud Run Services ===
BACKEND_SERVICE="${APP_NAME}-backend"
FRONTEND_SERVICE="${APP_NAME}-frontend"

deploy_backend_service "$BACKEND_SERVICE" "$SQL_REGION"
BACKEND_URL=$(get_service_url "$BACKEND_SERVICE" "$SQL_REGION")

deploy_frontend_service "$FRONTEND_SERVICE" "$SQL_REGION" "$BACKEND_URL"
FRONTEND_URL=$(get_service_url "$FRONTEND_SERVICE" "$SQL_REGION")

# === Output ===
echo ""
echo "========================================================"
echo "Setup complete!"
echo "========================================================"
echo "WIF Provider:    projects/$PROJECT_NUMBER/locations/global/workloadIdentityPools/$WIF_POOL/providers/$WIF_PROVIDER"
echo "Deploy SA:       $DEPLOY_SA_EMAIL"
echo "Runtime SA:      $RUNTIME_SA_EMAIL"
echo "Backend URL:     $BACKEND_URL"
echo "Frontend URL:    $FRONTEND_URL"
echo "SQL Connection:  $SQL_CONNECTION_NAME"
echo "Storage Bucket:  $GCS_BUCKET"
echo "Signed URL SA:   $SIGNED_URL_SA_EMAIL"
echo "========================================================"
#!/bin/bash
# scripts/cleanup.sh
set -e

source .env

# Security check
echo "==========================================================="
echo "WARNING: This will delete ALL trip-manager resources in the project"
echo "Project: $PROJECT_ID"
echo "==========================================================="
echo ""
read -p "Type 'DELETE' to continue: " confirmation
if [ "$confirmation" != "DELETE" ]; then
    echo "Aborted."
    exit 1
fi

# Set project before deletion to ensure nothing ends up in the wrong project
gcloud config set project "$PROJECT_ID"

echo ""
echo "=== Phase 1: Cloud Run Services ==="
gcloud run services delete "${APP_NAME}-backend" \
    --region="$REGION" --quiet 2>/dev/null || echo "Backend service not found, skipping"
gcloud run services delete "${APP_NAME}-frontend" \
    --region="$REGION" --quiet 2>/dev/null || echo "Frontend service not found, skipping"

echo ""
echo "=== Phase 2: Domain Mappings ==="
gcloud beta run domain-mappings delete "api.${DOMAIN}" \
    --region="$REGION" --quiet 2>/dev/null || echo "API domain mapping not found, skipping"
gcloud beta run domain-mappings delete "app.${DOMAIN}" \
    --region="$REGION" --quiet 2>/dev/null || echo "App domain mapping not found, skipping"

echo ""
echo "=== Phase 3: DNS Records and Zone ==="
ZONE_NAME=$(echo "$DOMAIN" | tr '.' '-')-zone

# Records individually (CNAME for subdomains) — SOA and NS records are deleted automatically with the zone
gcloud dns record-sets delete "api.${DOMAIN}." \
    --zone="$ZONE_NAME" --type=CNAME --quiet 2>/dev/null || echo "API CNAME not found, skipping"
gcloud dns record-sets delete "app.${DOMAIN}." \
    --zone="$ZONE_NAME" --type=CNAME --quiet 2>/dev/null || echo "App CNAME not found, skipping"

gcloud dns managed-zones delete "$ZONE_NAME" \
    --quiet 2>/dev/null || echo "DNS zone not found, skipping"

echo ""
echo "=== Phase 4: Cloud SQL ==="
# If deletion_protection is enabled, you need to disable it first
gcloud sql instances delete "${APP_NAME}-db" \
    --quiet 2>/dev/null || echo "SQL instance not found, skipping"

echo ""
echo "=== Phase 5: GCS Bucket ==="
# Force-delete including content
gsutil -m rm -r "gs://${GCS_BUCKET}" 2>/dev/null || echo "Bucket not found or already empty, skipping"

echo ""
echo "=== Phase 6: Secrets ==="
for secret in database-url jwt-secret; do
    gcloud secrets delete "$secret" --quiet 2>/dev/null || echo "Secret $secret not found, skipping"
done

echo ""
echo "=== Phase 7: Artifact Registry ==="
gcloud artifacts repositories delete "${AR_REPO:-trip-manager-images}" \
    --location="$REGION" --quiet 2>/dev/null || echo "Artifact registry not found, skipping"

echo ""
echo "=== Phase 8: Service Accounts ==="
# Removing bindings on the SAs first is not necessary - they will be deleted with the SA
SAs=(
    "${DEPLOY_SA_NAME}"
    "${APP_NAME}-runtime"
    "signed-url-signer"
)
for sa in "${SAs[@]}"; do
    sa_email="${sa}@${PROJECT_ID}.iam.gserviceaccount.com"
    gcloud iam service-accounts delete "$sa_email" --quiet 2>/dev/null \
        || echo "SA $sa_email not found, skipping"
done

echo ""
echo "=== Phase 9: Workload Identity Pool ==="
# Deleting the pool also automatically deletes all providers in it
gcloud iam workload-identity-pools delete "$WIF_POOL" \
    --location=global --quiet 2>/dev/null || echo "WIF pool not found, skipping"

echo ""
echo "=== Phase 10: Project-Level IAM Bindings ==="
# Bindings that point to deleted SAs are ineffective but technically remain in the policy.
# Let's clean up for a clean state:
echo "Cleaning up dead IAM bindings..."

DELETED_SAS=(
    "${DEPLOY_SA_NAME}@${PROJECT_ID}.iam.gserviceaccount.com"
    "${APP_NAME}-runtime@${PROJECT_ID}.iam.gserviceaccount.com"
    "signed-url-signer@${PROJECT_ID}.iam.gserviceaccount.com"
)



echo ""
echo "==========================================================="
echo "Cleanup completed!"
echo ""
echo "The following resources have not been deleted:"
echo "  - Firestore-Database"
echo "  - APIs"
echo "  - Project-Level IAM-Bindings for non SA users"
echo "  - Audit-Logs"
echo ""
echo "==========================================================="
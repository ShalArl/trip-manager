#!/bin/bash

# Service Account Setup Script for Microservices
PROJECT_ID=$(gcloud config get-value project)

# External Secrets SA
gcloud iam service-accounts create external-secrets-sa \
  --display-name="External Secrets Operator"

# Grant External Secrets Operator permissions to read secrets from Secret Manager
gcloud projects add-iam-policy-binding "$PROJECT_ID" \
  --member="serviceAccount:external-secrets-sa@${PROJECT_ID}.iam.gserviceaccount.com" \
  --role="roles/secretmanager.secretAccessor"

# WIF Binding for External Secrets Operator
gcloud iam service-accounts add-iam-policy-binding \
  external-secrets-sa@"${PROJECT_ID}".iam.gserviceaccount.com \
  --role="roles/iam.workloadIdentityUser" \
  --member="serviceAccount:${PROJECT_ID}.svc.id.goog[external-secrets/external-secrets]"

# Create Service Accounts
for sa in auth social presigner; do # Added a new service which requires a SA? Add service name here!
  gcloud iam service-accounts create ${sa}-sa \
    --display-name="${sa} service account"
done

# Add IAM Roles to Service Accounts
# auth → Firebase Auth
gcloud projects add-iam-policy-binding "$PROJECT_ID" \
  --member="serviceAccount:auth-sa@${PROJECT_ID}.iam.gserviceaccount.com" \
  --role="roles/firebase.sdkAdminServiceAgent"

# social → Firestore
gcloud projects add-iam-policy-binding "$PROJECT_ID" \
  --member="serviceAccount:social-sa@${PROJECT_ID}.iam.gserviceaccount.com" \
  --role="roles/datastore.user"

# presigner → GCS
gcloud projects add-iam-policy-binding "$PROJECT_ID" \
  --member="serviceAccount:presigner-sa@${PROJECT_ID}.iam.gserviceaccount.com" \
  --role="roles/storage.objectAdmin"

# Workload Identity Binding (GKE SA → GCP SA)
GKE_NAMESPACE=trip-manager-prod

for sa in auth social presigner; do
  gcloud iam service-accounts add-iam-policy-binding \
    ${sa}-sa@"${PROJECT_ID}".iam.gserviceaccount.com \
    --role="roles/iam.workloadIdentityUser" \
    --member="serviceAccount:${PROJECT_ID}.svc.id.goog[${GKE_NAMESPACE}/${sa}]"
done
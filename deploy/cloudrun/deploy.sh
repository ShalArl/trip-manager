#!/bin/bash
source .env.example
source scripts/lib.sh


ensure_wif "$WIF_POOL" "$WIF_PROVIDER"


#echo "Run auth login and set up project"
#gcloud auth login

#echo "Creating project and enable services"
#ensure_project "$PROJECT_ID"



#echo "Creating service account and granting permissions"
#PROJECT_NUMBER=$(gcloud projects describe $PROJECT_ID --format="value(projectNumber)")
#gcloud iam service-accounts create "$SERVICE_ACCOUNT_NAME" --display-name="$SERVICE_ACCOUNT_NAME"

#BUILD_SA=$(gcloud builds get-default-service-account)

#echo "Using build service account: $BUILD_SA"
# Grant required roles

#ensure_service_account "$SERVICE_ACCOUNT_NAME" "$SERVICE_ACCOUNT_NAME"

#add_iam_role "$SERVICE_ACCOUNT_NAME" "roles/cloudbuild.builds.builder"
#add_iam_role "$SERVICE_ACCOUNT_NAME" "roles/artifactregistry.writer"
#add_iam_role "$SERVICE_ACCOUNT_NAME" "roles/cloudsql.admin"
#add_iam_role "$SERVICE_ACCOUNT_NAME" "roles/storage.admin"

#gcloud run deploy "$SERVICE" \
#--source=./ \
#--platform=managed \
#--region=${REGION} \
#--allow-unauthenticated \
#--port=8081 \
#--memory=1Gi

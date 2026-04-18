#!/bin/bash

# Enable API services
ensure_service() {
    local service=$1
    if ! gcloud services list --enabled --filter="config.name:$service.googleapis.com" --format="value(config.name)" | grep -q "$service.googleapis.com"; then
        echo "Enabling $service API..."
        gcloud services enable "$service.googleapis.com"
    else
        echo "$service API is already enabled."
    fi
}

# Create a project if it doesn't exist and set it as the default
ensure_project() {
    local project_id=$1
    if ! gcloud projects describe "$project_id" &>/dev/null; then
        echo "Creating project $project_id..."
        gcloud projects create "$project_id" --name="$APP_NAME"
    else
        echo "Project $project_id already exists."
    fi
    gcloud config set project "$project_id"
}

# Create a Workload Identity Federation pool and provider if they don't exist
ensure_wif() {
    local pool_name=$1
    local provider_name=$2
    if ! gcloud iam workload-identity-pools describe "$pool_name" --location="global" &>/dev/null; then
        echo "Creating Workload Identity Pool $pool_name..."
        gcloud iam workload-identity-pools create "$pool_name" --location="global"
    else
        echo "Workload Identity Pool $pool_name already exists."
    fi
    if ! gcloud iam workload-identity-pools providers describe "$provider_name" --location="global" --workload-identity-pool="$pool_name" &>/dev/null; then
        echo "Creating Workload Identity Provider $provider_name..."
        gcloud iam workload-identity-pools providers create-oidc "$provider_name" \
            --location="global" \
            --workload-identity-pool="$pool_name" \
            --issuer-uri="https://accounts.google.com"
    else
        echo "Workload Identity Provider $provider_name already exists."
    fi
}

# Create a service account if it doesn't exist and assign a role
ensure_service_account() {
    local sa_name=$1
    local display_name=$2
    local full_sa="$sa_name@$MY_PROJECT_ID.iam.gserviceaccount.com"

    if ! gcloud iam service-accounts describe "$full_sa" &>/dev/null; then
        echo "Creating Service Account: $sa_name..."
        gcloud iam service-accounts create "$sa_name" --display-name="$display_name"
    else
        echo "Service Account $sa_name already exists."
    fi
}

# Create a service account if it doesn't exist and assign a role
ensure_sql_instance() {
    local instance_name=$1
    local tier=$2
    local region=$3

    if ! gcloud sql instances describe "$instance_name" &>/dev/null; then
        echo "Creating SQL Instance: $instance_name..."
        gcloud sql instances create "$instance_name" \
            --tier="$tier" \
            --region="$region"
    else
        echo "SQL Instance $instance_name already exists."
    fi
}


# Assign a specific IAM role to a service account
add_iam_role() {
    local sa_name=$1
    local role=$2
    local full_sa="serviceAccount:$sa_name@$MY_PROJECT_ID.iam.gserviceaccount.com"

    echo "Assigning role $role to $sa_name..."
    gcloud projects add-iam-policy-binding "$MY_PROJECT_ID" \
        --member="$full_sa" \
        --role="$role" --quiet > /dev/null
}




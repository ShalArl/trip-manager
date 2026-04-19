#!/bin/bash

# Enable API services
ensure_services() {
    local services=("$@")
    echo "Enabling services: ${services[*]}..."
    gcloud services enable "${services[@]}" --project="$PROJECT_ID"
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

    # Pool Check
    if ! gcloud iam workload-identity-pools describe "$pool_name" --location="global" &>/dev/null; then
        echo "Creating Workload Identity Pool $pool_name..."
        gcloud iam workload-identity-pools create "$pool_name" --location="global"
    else
        echo "Workload Identity Pool $pool_name already exists."
    fi

    # Provider Check
    if ! gcloud iam workload-identity-pools providers describe "$provider_name" --location="global" --workload-identity-pool="$pool_name" &>/dev/null; then
        echo "Creating Workload Identity Provider $provider_name..."
        gcloud iam workload-identity-pools providers create-oidc "$provider_name" \
            --location="global" \
            --workload-identity-pool="$pool_name" \
            --attribute-mapping="google.subject=assertion.sub,attribute.actor=assertion.actor,attribute.repository=assertion.repository" \
            --attribute-condition="assertion.repository == 'ShalArl/trip-manager'" \
            --issuer-uri="https://token.actions.githubusercontent.com"
    else
        echo "Workload Identity Provider $provider_name already exists."
    fi
}

# Create a service account if it doesn't exist and assign a role
ensure_service_account() {
    local sa_name=$1
    local display_name=$2
    local full_sa="$sa_name@$PROJECT_ID.iam.gserviceaccount.com"

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
    local full_sa="serviceAccount:$sa_name@$PROJECT_ID.iam.gserviceaccount.com"

    echo "Assigning role $role to $sa_name..."
    gcloud projects add-iam-policy-binding "$PROJECT_ID" \
        --member="$full_sa" \
        --role="$role" --quiet > /dev/null
}

# Create a secret if it doesn't exist
create_secret_if_missing() {
    local name=$1
    local value=$2
    if ! gcloud secrets describe "$name" --project="$PROJECT_ID" &>/dev/null; then
        echo -n "$value" | gcloud secrets create "$name" \
            --project="$PROJECT_ID" --replication-policy="automatic" --data-file=-
        echo "Created secret: $name"
    else
        echo "Secret $name already exists, skipping."
    fi
}

# Get the value of a secret
get_secret() {
    local name=$1
    gcloud secrets versions access latest --secret="$name" --project="$PROJECT_ID"
}

# Update the value of an existing secret
update_secret() {
    local name=$1
    local value=$2
    echo -n "$value" | gcloud secrets versions add "$name" \
        --project="$PROJECT_ID" --data-file=-
}

wait_for_service_account() {
    local sa_email=$1
    local max_attempts=30
    local attempt=0

    echo "Waiting for service account $sa_email to propagate..."
    while [ $attempt -lt $max_attempts ]; do
        if gcloud iam service-accounts describe "$sa_email" --project="$PROJECT_ID" &>/dev/null; then
            # Zusätzlich kurz warten, weil describe oft schneller konsistent ist als andere APIs
            sleep 5
            echo "Service account $sa_email is now available."
            return 0
        fi
        attempt=$((attempt + 1))
        sleep 2
    done
    echo "ERROR: Service account $sa_email did not become available after $((max_attempts * 2))s"
    return 1
}
#!/bin/bash

sql_setup_db_password() {
    if ! gcloud secrets describe "db-password" --project="$PROJECT_ID" &>/dev/null; then
        echo "Generating DB password..."
        DB_PASSWORD=$(openssl rand -base64 32 | tr -d '/+=' | cut -c1-32)
        create_secret_if_missing "db-password" "$DB_PASSWORD"
    else
        echo "DB password already in Secret Manager, reusing."
        DB_PASSWORD=$(get_secret "db-password")
    fi
    export DB_PASSWORD
}

sql_setup_instance() {
    local sql_instance=$1
    local db_name=$2
    local db_user=$3
    local db_password=$4
    local region=$5
    local sql_tier=$6

    if ! gcloud sql instances describe "$sql_instance" --project="$PROJECT_ID" &>/dev/null; then
        echo "Creating Cloud SQL instance $sql_instance (this takes ~5min)..."
        gcloud sql instances create "$sql_instance" \
            --project="$PROJECT_ID" \
            --database-version=POSTGRES_16 \
            --tier="$sql_tier" \
            --region="$region" \
            --edition=ENTERPRISE \
            --storage-size=10GB \
            --storage-type=HDD \
            --backup-start-time=03:00 \
            --availability-type=zonal
    else
        echo "Cloud SQL instance $sql_instance already exists."
    fi

    # Passwort generieren und in Secret Manager ablegen
    if ! gcloud secrets describe "db-password" --project="$PROJECT_ID" &>/dev/null; then
        echo "Generating DB password and storing in Secret Manager..."
        DB_PASSWORD=$(openssl rand -base64 32 | tr -d '/+=' | cut -c1-32)
        echo -n "$DB_PASSWORD" | gcloud secrets create "db-password" \
            --project="$PROJECT_ID" \
            --replication-policy="automatic" \
            --data-file=-
    else
        echo "Secret db-password already exists, reusing."
        DB_PASSWORD=$(gcloud secrets versions access latest --secret="db-password" --project="$PROJECT_ID")
    fi

    # DB-User anlegen (idempotent via || true, weil gcloud hier kein describe für User hat)
    gcloud sql users create "$db_user" \
        --instance="$sql_instance" \
        --password="$db_password" \
        --project="$PROJECT_ID" 2>/dev/null || echo "DB user $db_user already exists."

    # Database anlegen
    if ! gcloud sql databases describe "$db_name" \
        --instance="$sql_instance" --project="$PROJECT_ID" &>/dev/null; then
        gcloud sql databases create "$db_name" \
            --instance="$sql_instance" \
            --project="$PROJECT_ID"
    fi

    SQL_CONNECTION_NAME="${PROJECT_ID}:${region}:${sql_instance}"
    echo "Cloud SQL connection name: $SQL_CONNECTION_NAME"
}

sql_create_db_user() {
    local instance_name=$1
    local db_user=$2
    local db_password=$3

    # DB-User anlegen (idempotent via || true, weil gcloud hier kein describe für User hat)
    gcloud sql users create "$db_user" \
        --instance="$instance_name" \
        --password="$db_password" \
        --project="$PROJECT_ID" 2>/dev/null || echo "DB user $db_user already exists."
}

sql_create_database() {
    local instance_name=$1
    local db_name=$2

    # Database anlegen
    if ! gcloud sql databases describe "$db_name" \
        --instance="$instance_name" --project="$PROJECT_ID" &>/dev/null; then
        gcloud sql databases create "$db_name" \
            --instance="$instance_name" \
            --project="$PROJECT_ID"
    fi
}
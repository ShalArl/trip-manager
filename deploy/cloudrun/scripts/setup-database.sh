source ../.env
source lib.sh

ensure_sql_instance "$DATABASE_INSTANCE_NAME" "$DATABASE_TIER" "$DATABASE_REGION"

gcloud sql databases create "$DATABASE_NAME" --instance="$DATABASE_INSTANCE_NAME" || echo "DB exists"

gcloud sql users create "$DATABASE_USER" \
    --instance="$DATABASE_INSTANCE_NAME" \
    --password="$DATABASE_PASSWORD" || echo "User exists"
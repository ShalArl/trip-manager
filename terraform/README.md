# Terraform Infrastructure Deployment Guide

## 1. Prerequisites
- Terraform installed (version 1.14 or later)
- gcloud CLI installed and configured with appropriate permissions

## 2. Setup
1. Create the terraform state bucket in gcloud:
   ```bash
    gcloud storage buckets create gs://${PROJECT_ID}-terraform-state \
    --location=europe-west1 \
    --uniform-bucket-level-access
   ```
   **Important**: Activate versioning for the bucket to enable state versioning and recovery:
   ```bash
   gcloud storage buckets update gs://${PROJECT_ID}-terraform-state --versioning
   ```
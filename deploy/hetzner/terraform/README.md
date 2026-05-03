# Hetzner Terraform Deployment

Infrastructure-as-Code for managing the Trip Manager infrastructure on Hetzner Cloud.

## Prerequisites

- Terraform >= 1.0
- Hetzner Cloud API Token (`.env`)
- SSH Key for Root Access

## Structure

- `main.tf` - Main resources (server, networking, etc.)
- `variables.tf` - Variable definitions
- `outputs.tf` - Output important values
- `terraform.tfvars` - Configuration (not in Git)

## Usage

```bash
# Initialize
terraform init

# View plan
terraform plan

# Deploy
terraform apply

# Destroy
terraform destroy
```

## Environment

Required variables in `.env`:
```
HCLOUD_TOKEN=<api-token>
```

## Further Information

See `../manual-deploy.sh` for manual deployment steps.


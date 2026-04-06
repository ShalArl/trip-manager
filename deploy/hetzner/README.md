# Deployment

This directory contains deployment scripts and configurations for the Trip Manager application deployment on any remote Server with enabled SSH access.

## Contents
**deployment** contains:
- `deploy.sh`: A bash script to automate the deployment process. It handles building the application, transferring files to the server, and restarting services. (used in the GitHub Actions workflow)
- `Caddyfile`: Configuration file for Caddy web server, used to serve the frontend and reverse proxy API requests to the backend which is running as a container on the server. (see [docker-compose.yml](../../docker-compose.yaml))

in the root dir of `deploy/hetzner` you can find `manual-deploy.sh` which is a simplified version of `deploy.sh` for manual deployments without GitHub Actions  
please note a GitHub ghcr pat token as well as username is required for pulling the images on the server. (see [.env.example](./.env.example) for an overview of required environment variables and their purpose)


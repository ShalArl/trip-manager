#cloud-config
users:
  - name: deployer
    groups: [sudo, docker]
    shell: /bin/bash
    sudo: ['ALL=(ALL) NOPASSWD:ALL']
    ssh_authorized_keys:
      - ${ssh_public_key}

write_files:
  # SSH-Hardening
  - path: /etc/ssh/sshd_config.d/hardening.conf
    content: |
      PermitRootLogin no
      PasswordAuthentication no

  # Docker-Compose-File
  - path: /home/deployer/app/docker-compose.yml
    owner: deployer:deployer
    permissions: '0644'
    defer: true
    content: |
      services:
        backend:
          image: ${backend_image}
          restart: unless-stopped
          environment:
            ENVIRONMENT: production
            SERVER_PORT: "8081"
            CORS_ALLOWED_ORIGINS: "https://iaas-app.neatnode.xyz"
            FIREBASE_PROJECT_ID: ${firebase_project_id}
            DATABASE_URL: postgres://${postgres_user}:${postgres_password}@postgres:5432/${postgres_db}?sslmode=disable
            STORAGE_TYPE: s3
            S3_ENDPOINT: http://minio:9000
            S3_BUCKET: trip-manager
            S3_REGION: us-east-1
            S3_ACCESS_KEY: ${minio_access_key}
            S3_SECRET_KEY: ${minio_secret_key}
            S3_USE_SSL: "false"
            S3_PUBLIC_URL: https://iaas-storage.neatnode.xyz
          depends_on:
            postgres:
              condition: service_healthy
            minio:
              condition: service_healthy
          networks:
            - app-net

        frontend:
          image: ${frontend_image}
          restart: unless-stopped
          environment:
            NODE_ENV: production
          networks:
            - app-net

        postgres:
          image: ${postgres_image}
          restart: unless-stopped
          environment:
            POSTGRES_USER: ${postgres_user}
            POSTGRES_PASSWORD: ${postgres_password}
            POSTGRES_DB: ${postgres_db}
          volumes:
            - postgres-data:/var/lib/postgresql/data
          healthcheck:
            test: ["CMD-SHELL", "pg_isready -U ${postgres_user} -d ${postgres_db}"]
            interval: 10s
            timeout: 5s
            retries: 5
          networks:
            - app-net

        minio:
          image: quay.io/minio/minio:latest
          restart: unless-stopped
          command: server /data --console-address ":9001"
          environment:
            MINIO_ROOT_USER: ${minio_access_key}
            MINIO_ROOT_PASSWORD: ${minio_secret_key}
            MINIO_API_CORS_ALLOW_ORIGIN: "https://iaas-app.neatnode.xyz"
          volumes:
            - minio-data:/data
          healthcheck:
            test: ["CMD", "curl", "-f", "http://localhost:9000/minio/health/live"]
            interval: 10s
            timeout: 5s
            retries: 5
          networks:
            - app-net

        minio-init:
          image: quay.io/minio/minio:latest
          depends_on:
            minio:
              condition: service_healthy
          entrypoint: >
            /bin/sh -c "
            mc alias set myminio http://minio:9000 ${minio_access_key} ${minio_secret_key};
            mc mb --ignore-existing myminio/trip-manager;
            mc anonymous set download myminio/trip-manager;
            exit 0;
            "
          networks:
            - app-net

        caddy:
          image: caddy:2-alpine
          restart: unless-stopped
          ports:
            - "80:80"
            - "443:443"
          volumes:
            - ./Caddyfile:/etc/caddy/Caddyfile:ro
            - caddy-data:/data
            - caddy-config:/config
          depends_on:
            - backend
            - frontend
            - minio
          networks:
            - app-net

      volumes:
        postgres-data:
        caddy-data:
        caddy-config:
        minio-data:

      networks:
        app-net:


  # GHCR Credentials (für Skripte)
  - path: /home/deployer/scripts/ghcr-token
    owner: deployer:deployer
    permissions: '0600'
    defer: true
    content: |
      ${github_registry_token}

  # Caddyfile
  - path: /home/deployer/app/Caddyfile
    owner: deployer:deployer
    permissions: '0644'
    defer: true
    content: |
      iaas.neatnode.xyz {
        reverse_proxy backend:8081
        encode gzip

        log {
          output stdout
          format console
        }
      }

      iaas-app.neatnode.xyz {
        reverse_proxy frontend:3000
        encode gzip

        log {
          output stdout
          format console
        }
      }

      iaas-storage.neatnode.xyz {
        reverse_proxy minio:9000
        encode gzip

        request_body {
          max_size 100MB
        }
      }

  # Docker-Login-Skript für GHCR
  - path: /home/deployer/scripts/docker-login.sh
    owner: deployer:deployer
    permissions: '0755'
    defer: true
    content: |
      #!/bin/bash
      set -e
      if [ -f /home/deployer/scripts/ghcr-token ]; then
        cat /home/deployer/scripts/ghcr-token | docker login ghcr.io -u "${github_username}" --password-stdin
      fi

  # Update-Skript für CD
  - path: /home/deployer/scripts/update.sh
    permissions: '0755'
    content: |
      #!/bin/bash
      set -e
      cd /home/deployer/app
      docker compose pull
      docker compose up -d --remove-orphans
      docker image prune -f

package_update: true
# package_upgrade: true

packages:
  - docker.io
  - ufw
  - fail2ban

runcmd:
  # Install Docker-Compose
  - install -m 0755 -d /etc/apt/keyrings
  - curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
  - chmod a+r /etc/apt/keyrings/docker.asc
  - bash -c 'echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu $(. /etc/os-release && echo $VERSION_CODENAME) stable" > /etc/apt/sources.list.d/docker.list'
  - apt-get update
  - DEBIAN_FRONTEND=noninteractive apt-get install -y docker-compose-plugin

  # FS Permissions
  - chown -R deployer:deployer /home/deployer

  # Docker
  - systemctl enable --now docker
  - usermod -aG docker deployer

  # Permissions
  - chown -R deployer:deployer /home/deployer/app
  - chown -R deployer:deployer /home/deployer/scripts
  - chmod +x /home/deployer/scripts/*.sh

  # SSH-Restart
  - systemctl restart ssh

  # GHCR Login
  - su - deployer -c "/home/deployer/scripts/docker-login.sh"

  # App starten
  - su - deployer -c "cd /home/deployer/app && docker compose pull && docker compose up -d"
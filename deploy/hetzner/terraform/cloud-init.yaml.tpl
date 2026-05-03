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
    content: |
      services:
        backend:
          image: ${backend_image}
          restart: unless-stopped
          environment:
            ENVIRONMENT: production
            SERVER_PORT: "8081"
            CORS_ALLOWED_ORIGINS: ${cors_allowed_origins}
            FIREBASE_PROJECT_ID: ${firebase_project_id}
            DATABASE_URL: postgres://${postgres_user}:${postgres_password}@postgres:5432/${postgres_db}?sslmode=disable
            STORAGE_TYPE: local
          depends_on:
            postgres:
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
          networks:
            - app-net

      volumes:
        postgres-data:
        caddy-data:
        caddy-config:

      networks:
        app-net:

  # Caddyfile
  - path: /home/deployer/app/Caddyfile
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

  # Docker-Login-Skript für GHCR
  - path: /home/deployer/scripts/docker-login.sh
    permissions: '0700'
    content: |
      #!/bin/bash
      set -e
      if [ -n "${github_registry_token}" ]; then
        echo "${github_registry_token}" | docker login ghcr.io -u "${github_username}" --password-stdin
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
package_upgrade: true

packages:
  - docker.io
  - docker-compose-plugin
  - ufw
  - fail2ban

runcmd:
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
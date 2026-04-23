# Build stage
FROM node:24-alpine AS builder

WORKDIR /app

# Build arguments
ARG NEXT_PUBLIC_API_URL=http://localhost:8000
ARG NEXT_PUBLIC_S3_PUBLIC_URL=http://localhost/minio

ENV NEXT_PUBLIC_API_URL=$NEXT_PUBLIC_API_URL
ENV NEXT_PUBLIC_S3_PUBLIC_URL=$NEXT_PUBLIC_S3_PUBLIC_URL

# Firebase environment variables
ARG NEXT_PUBLIC_FIREBASE_API_KEY=api-key
ARG NEXT_PUBLIC_FIREBASE_APP_ID=app-id
ARG NEXT_PUBLIC_FIREBASE_AUTH_DOMAIN=auth-domain
ARG NEXT_PUBLIC_FIREBASE_PROJECT_ID=project-id

ENV NEXT_PUBLIC_FIREBASE_API_KEY=$NEXT_PUBLIC_FIREBASE_API_KEY
ENV NEXT_PUBLIC_FIREBASE_APP_ID=$NEXT_PUBLIC_FIREBASE_APP_ID
ENV NEXT_PUBLIC_FIREBASE_AUTH_DOMAIN=$NEXT_PUBLIC_FIREBASE_AUTH_DOMAIN
ENV NEXT_PUBLIC_FIREBASE_PROJECT_ID=$NEXT_PUBLIC_FIREBASE_PROJECT_ID


# Install pnpm & dependencies
RUN npm install -g pnpm@9.0.0

COPY package.json pnpm-lock.yaml pnpm-workspace.yaml ./
COPY api-spec ./api-spec
COPY frontend/package.json ./frontend/

RUN pnpm install --frozen-lockfile

# Source Code kopieren
COPY frontend ./frontend

# Types generieren & Build
RUN mkdir -p ./frontend/generated && npx openapi-typescript ./api-spec/openapi.yaml -o ./frontend/generated/types.ts
RUN pnpm --filter frontend build

# ---

# Runtime stage
FROM node:24-alpine
WORKDIR /app

RUN npm install -g pnpm@9.0.0

COPY package.json pnpm-lock.yaml pnpm-workspace.yaml ./
COPY frontend/package.json ./frontend/

# Nur Production-Abhängigkeiten installieren
RUN pnpm install --prod --frozen-lockfile --filter frontend

# Build-Artefakte kopieren (in den korrekten Unterordner!)
COPY --from=builder /app/frontend/.next ./frontend/.next
COPY --from=builder /app/frontend/public ./frontend/public
COPY --from=builder /app/frontend/next.config.ts ./frontend/

EXPOSE 3000

WORKDIR /app/frontend
CMD ["pnpm", "start"]
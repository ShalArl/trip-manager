# ==========================================
# STAGE 1: Build stage
# ==========================================
FROM node:24-alpine AS builder

WORKDIR /app

# Install pnpm global
RUN npm install -g pnpm@9.0.0

# Dependencies kopieren
COPY package.json pnpm-lock.yaml pnpm-workspace.yaml ./
COPY api-spec ./api-spec
COPY frontend/package.json ./frontend/

RUN pnpm install --frozen-lockfile

# Quellcode kopieren (MUSS vor dem Type-Gen und Build passieren!)
COPY frontend ./frontend

# Typen generieren aus der OpenAPI Spec
RUN mkdir -p ./frontend/generated && npx openapi-typescript ./api-spec/merged/openapi.yaml -o ./frontend/generated/types.ts

# JETZT ERST die Build-Argumente definieren (Direkt vor dem Next.js Build!)
ARG NEXT_PUBLIC_API_URL
ENV NEXT_PUBLIC_API_URL=$NEXT_PUBLIC_API_URL

ARG NEXT_PUBLIC_FIREBASE_API_KEY
ENV NEXT_PUBLIC_FIREBASE_API_KEY=$NEXT_PUBLIC_FIREBASE_API_KEY

ARG NEXT_PUBLIC_FIREBASE_APP_ID
ENV NEXT_PUBLIC_FIREBASE_APP_ID=$NEXT_PUBLIC_FIREBASE_APP_ID

ARG NEXT_PUBLIC_FIREBASE_AUTH_DOMAIN
ENV NEXT_PUBLIC_FIREBASE_AUTH_DOMAIN=$NEXT_PUBLIC_FIREBASE_AUTH_DOMAIN

ARG NEXT_PUBLIC_FIREBASE_PROJECT_ID
ENV NEXT_PUBLIC_FIREBASE_PROJECT_ID=$NEXT_PUBLIC_FIREBASE_PROJECT_ID

ARG NEXT_PUBLIC_USE_FIREBASE_EMULATOR=false
ENV NEXT_PUBLIC_USE_FIREBASE_EMULATOR=$NEXT_PUBLIC_USE_FIREBASE_EMULATOR

# Next.js Build ausführen - brennt die Variablen unumstößlich ein
RUN pnpm --filter frontend build

# ==========================================
# STAGE 2: Runtime stage (Standalone Mode)
# ==========================================
FROM node:24-alpine AS runner
WORKDIR /app

ENV NODE_ENV=production
ENV PORT=3000
ENV HOSTNAME="0.0.0.0"

COPY --from=builder /app/frontend/public ./frontend/public
COPY --from=builder /app/frontend/.next/standalone ./
COPY --from=builder /app/frontend/.next/static ./frontend/.next/static

EXPOSE 3000

# Next.js Standalone startet direkt über den node-Server (pnpm wird nicht mehr benötigt)
CMD ["node", "frontend/server.js"]
# Build stage
FROM node:20-alpine AS builder

WORKDIR /app

# Build arguments
ARG NEXT_PUBLIC_API_URL=http://localhost:8000
ENV NEXT_PUBLIC_API_URL=$NEXT_PUBLIC_API_URL

# Install pnpm
RUN npm install -g pnpm@9.0.0

# Copy only necessary monorepo files
COPY package.json pnpm-lock.yaml ./
COPY api-spec ./api-spec
COPY frontend ./frontend

# Install dependencies
RUN pnpm install --frozen-lockfile

# Generate TypeScript code from OpenAPI spec
RUN mkdir -p ./frontend/generated && npx openapi-typescript ./api-spec/openapi.yaml -o ./frontend/generated/types.ts

# Build Next.js app
WORKDIR /app/frontend
RUN pnpm build

# Deploy production dependencies using workspace root lockfile
WORKDIR /app
RUN pnpm --filter frontend deploy --prod /deploy

# Runtime stage
FROM node:20-alpine

WORKDIR /app

# Copy production dependencies (resolved from workspace root lockfile via pnpm deploy)
COPY --from=builder /deploy .

# Copy built application and public assets
COPY --from=builder /app/frontend/.next ./.next
COPY --from=builder /app/frontend/public ./public

EXPOSE 3000

CMD ["node_modules/.bin/next", "start"]

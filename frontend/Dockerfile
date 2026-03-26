# Build stage
FROM node:20-alpine AS builder

WORKDIR /app

# Install pnpm
RUN npm install -g pnpm@9.0.0

# Copy monorepo files
COPY . .

# Install dependencies
RUN CI=true pnpm install --frozen-lockfile

# Generate TypeScript code from OpenAPI spec
RUN pnpm gen:ts

# Build Next.js app
WORKDIR /app/frontend
RUN pnpm build

# Runtime stage
FROM node:20-alpine

WORKDIR /app

# Install pnpm in runtime
RUN npm install -g pnpm@9.0.0

# Copy package files
COPY --from=builder /app/frontend/package.json .
COPY --from=builder /app/frontend/pnpm-lock.yaml .

# Install production dependencies only
RUN pnpm install --frozen-lockfile --prod

# Copy built application and public assets
COPY --from=builder /app/frontend/.next ./.next
COPY --from=builder /app/frontend/public ./public

EXPOSE 3000

CMD ["pnpm", "start"]

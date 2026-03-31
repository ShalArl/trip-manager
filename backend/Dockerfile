# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install Node.js and pnpm for code generation
RUN apk add --no-cache nodejs npm
RUN npm install -g pnpm@9.0.0

# Copy only necessary monorepo files
COPY package.json pnpm-lock.yaml ./
COPY api-spec ./api-spec
COPY backend ./backend
COPY oapi-codegen.yaml ./

# Install oapi-codegen (Go tool for code generation)
RUN go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.4.1


RUN pnpm install

# Generate Go code from OpenAPI spec
RUN mkdir -p ./backend/internal/generated && rm -rf ./backend/internal/generated/*.go && oapi-codegen --config=oapi-codegen.yaml ./api-spec/openapi.yaml

# Build Go application
WORKDIR /app/backend
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o trip-manager ./cmd/api

# Runtime stage
FROM alpine:3.19

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/backend/trip-manager .

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8000/health || exit 1

EXPOSE 8000

CMD ["./trip-manager"]

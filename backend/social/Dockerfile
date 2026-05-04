FROM golang:1.25-alpine AS builder
WORKDIR /app

# go.work zuerst – damit Go den Workspace kennt
COPY go.work go.work.sum ./

# Nur die go.mod/go.sum der relevanten Module für Layer-Caching
COPY auth/go.mod auth/go.mod
COPY auth/go.sum auth/go.sum
COPY presigner/go.mod presigner/go.mod
COPY presigner/go.sum presigner/go.sum
COPY shared/authclient/go.mod shared/authclient/go.mod
COPY shared/authclient/go.sum shared/authclient/go.sum
COPY social/go.mod social/go.mod
COPY social/go.sum social/go.sum

RUN go work sync

# Source Code
COPY shared/authclient/ shared/authclient/
COPY social/ social/

# Aus dem social-Verzeichnis bauen
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/social ./social/cmd/api

# Runtime stage
FROM alpine:3.19
RUN addgroup -S app && adduser -S -G app app
USER app

COPY --from=builder /bin/social /bin/social

HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

EXPOSE 8080
CMD ["/bin/social"]
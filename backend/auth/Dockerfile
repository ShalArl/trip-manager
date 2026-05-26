FROM golang:1.25-alpine AS builder
WORKDIR /app

# go.work zuerst – damit Go den Workspace kennt
COPY go.work go.work.sum ./

# Nur die go.mod/go.sum der relevanten Module für Layer-Caching
COPY auth/go.mod auth/go.mod
COPY auth/go.sum auth/go.sum
COPY presigner/go.mod presigner/go.mod
COPY presigner/go.sum presigner/go.sum
COPY social/go.mod social/go.mod
COPY social/go.sum social/go.sum
COPY travel-warning/go.mod travel-warning/go.mod
COPY travel-warning/go.sum travel-warning/go.sum
COPY weather-info/go.mod weather-info/go.mod
COPY weather-info/go.sum weather-info/go.sum
COPY shared/authclient/go.mod shared/authclient/go.mod
COPY shared/authclient/go.sum shared/authclient/go.sum
COPY shared/middleware/go.mod shared/middleware/go.mod
COPY shared/middleware/go.sum shared/middleware/go.sum
COPY shared/userclient/go.mod shared/userclient/go.mod
COPY shared/userclient/go.sum shared/userclient/go.sum
COPY users/go.mod users/go.mod
COPY users/go.sum users/go.sum
COPY trips/go.mod trips/go.mod
COPY trips/go.sum trips/go.sum
COPY locations/go.mod locations/go.mod
COPY locations/go.sum locations/go.sum
COPY feed/go.mod feed/go.mod
COPY feed/go.sum feed/go.sum
COPY feed-generator/go.mod feed-generator/go.mod
COPY feed-generator/go.sum feed-generator/go.sum

RUN go work sync

# Source Code
COPY shared/authclient/ shared/authclient/
COPY shared/middleware/ shared/middleware/

COPY auth/ auth/

# Aus dem auth-Verzeichnis bauen
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/auth ./auth/cmd/api

# Runtime stage
FROM alpine:3.19
RUN addgroup -S app && adduser -S -G app app
USER app

WORKDIR /app

# Copy binary from builder
COPY --from=builder /bin/auth /bin/auth

HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

EXPOSE 8080
CMD ["/bin/auth"]

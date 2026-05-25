FROM golang:1.22-alpine AS builder

WORKDIR /app

# go.work aus dem backend/ root kopieren
COPY go.work go.work.sum* ./

# Nur die nötigen Module kopieren
COPY shared/ shared/
COPY feed-generator/ feed-generator/

WORKDIR /app/feed-generator
RUN go build -o /feed-generator ./cmd/worker

FROM alpine:3.19
RUN apk --no-cache add ca-certificates
COPY --from=builder /feed-generator /feed-generator
ENTRYPOINT ["/feed-generator"]
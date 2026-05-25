FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.work go.work.sum* ./
COPY shared/ shared/
COPY feed/ feed/

WORKDIR /app/feed
RUN go build -o /feed ./cmd/api

FROM alpine:3.19
RUN apk --no-cache add ca-certificates
COPY --from=builder /feed /feed
ENTRYPOINT ["/feed"]
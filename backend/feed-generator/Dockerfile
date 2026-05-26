FROM golang:1.25-alpine AS builder
WORKDIR /app

COPY go.work go.work.sum ./

COPY auth/go.mod auth/go.mod
COPY auth/go.sum auth/go.sum
COPY presigner/go.mod presigner/go.mod
COPY presigner/go.sum presigner/go.sum
COPY social/go.mod social/go.mod
COPY social/go.sum social/go.sum
COPY shared/authclient/go.mod shared/authclient/go.mod
COPY shared/authclient/go.sum shared/authclient/go.sum
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

COPY feed-generator/ feed-generator/

RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/feed-generator ./feed-generator/cmd/worker

FROM alpine:3.19
RUN addgroup -S app && adduser -S -G app app
USER app
WORKDIR /app
COPY --from=builder /bin/feed-generator /bin/feed-generator
CMD ["/bin/feed-generator"]
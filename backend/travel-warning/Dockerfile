FROM golang:1.25-alpine AS builder
WORKDIR /app

COPY go.work go.work.sum ./
COPY auth/go.mod auth/go.mod
COPY auth/go.sum auth/go.sum
COPY presigner/go.mod presigner/go.mod
COPY presigner/go.sum presigner/go.sum
COPY social/go.mod social/go.mod
COPY social/go.sum social/go.sum
COPY users/go.mod users/go.mod
COPY users/go.sum users/go.sum
COPY trips/go.mod trips/go.mod
COPY trips/go.sum trips/go.sum
COPY locations/go.mod locations/go.mod
COPY locations/go.sum locations/go.sum
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
COPY shared/tenantdb/go.mod shared/tenantdb/go.mod
COPY shared/tenantdb/go.sum shared/tenantdb/go.sum
COPY shared/firebaseclient/go.mod shared/firebaseclient/go.mod
COPY shared/firebaseclient/go.sum shared/firebaseclient/go.sum
COPY shared/otel/go.mod shared/otel/go.mod
COPY shared/otel/go.sum shared/otel/go.sum
COPY feed/go.mod feed/go.mod
COPY feed/go.sum feed/go.sum
COPY feed-generator/go.mod feed-generator/go.mod
COPY feed-generator/go.sum feed-generator/go.sum
COPY newsletter/go.mod newsletter/go.mod
COPY newsletter/go.sum newsletter/go.sum
COPY newsletter-worker/go.mod newsletter-worker/go.mod
COPY newsletter-worker/go.sum newsletter-worker/go.sum

RUN go work sync

COPY shared/ shared/
COPY travel-warning/ travel-warning/

# Beide Binaries bauen
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/api ./travel-warning/cmd/api
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/worker ./travel-warning/cmd/worker

FROM alpine:3.21
RUN addgroup -S app && adduser -S -G app app
USER app

COPY --from=builder /bin/api /bin/api
COPY --from=builder /bin/worker /bin/worker

EXPOSE 8090
CMD ["/bin/api"]
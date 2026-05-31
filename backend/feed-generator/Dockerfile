FROM golang:1.25.8-alpine AS builder

WORKDIR /app

# go.work aus dem backend/ root kopieren
COPY go.work go.work.sum* ./
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

# Nur die nötigen Module kopieren
COPY shared/ shared/
COPY feed-generator/ feed-generator/

WORKDIR /app/feed-generator
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/api ./cmd/worker

FROM alpine:3.21
RUN addgroup -S app && adduser -S -G app app
USER app

COPY --from=builder /bin/api /bin/api

EXPOSE 8091
CMD ["/bin/api"]
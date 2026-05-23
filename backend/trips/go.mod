module github.com/ShalArl/trip-manager/backend/trips

go 1.25.8

require (
	github.com/ShalArl/trip-manager/backend/shared/authclient v0.0.0-00010101000000-000000000000
	github.com/google/uuid v1.6.0
	github.com/jmoiron/sqlx v1.4.0
	github.com/lib/pq v1.12.3
	github.com/oapi-codegen/runtime v1.4.0
	github.com/segmentio/kafka-go v0.4.47
)

require (
	github.com/klauspost/compress v1.16.7 // indirect
	github.com/pierrec/lz4/v4 v4.1.15 // indirect
)

replace github.com/ShalArl/trip-manager/backend/shared/authclient => ../shared/authclient

module github.com/ShalArl/trip-manager/backend/trips

go 1.25.8

require (
	github.com/ShalArl/trip-manager/backend/shared/authclient v0.0.0
	github.com/ShalArl/trip-manager/backend/shared/middleware v0.0.0
	github.com/google/uuid v1.6.0
	github.com/jmoiron/sqlx v1.4.0
	github.com/lib/pq v1.12.3
	github.com/oapi-codegen/runtime v1.4.1
	github.com/segmentio/kafka-go v0.4.47
)

replace github.com/ShalArl/trip-manager/backend/shared/middleware => ../shared/middleware

require (
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
)

require (
	github.com/klauspost/compress v1.17.7 // indirect
	github.com/pierrec/lz4/v4 v4.1.15 // indirect
	golang.org/x/net v0.53.0 // indirect
	golang.org/x/text v0.37.0 // indirect
)

replace github.com/ShalArl/trip-manager/backend/shared/authclient => ../shared/authclient

module github.com/ShalArl/trip-manager/backend/trips

go 1.25.8

require (
	github.com/ShalArl/trip-manager/backend/shared/authclient v0.0.0
    github.com/ShalArl/trip-manager/backend/shared/middleware v0.0.0
    github.com/ShalArl/trip-manager/backend/shared/userclient v0.0.0
	github.com/google/uuid v1.6.0
	github.com/jmoiron/sqlx v1.4.0
	github.com/lib/pq v1.12.3
	github.com/oapi-codegen/runtime v1.4.0
)

replace (
	github.com/ShalArl/trip-manager/backend/shared/authclient => ../shared/authclient
	github.com/ShalArl/trip-manager/backend/shared/middleware => ../shared/middleware
	github.com/ShalArl/trip-manager/backend/shared/userclient => ../shared/userclient
)

require (
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
)
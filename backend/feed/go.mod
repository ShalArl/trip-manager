module github.com/ShalArl/trip-manager/backend/feed

go 1.25.8

require (
	github.com/ShalArl/trip-manager/backend/shared/authclient v0.0.0-00010101000000-000000000000
	github.com/google/uuid v1.6.0
	github.com/neo4j/neo4j-go-driver/v5 v5.26.0
	github.com/oapi-codegen/runtime v1.4.1
)

require (
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
)

replace github.com/ShalArl/trip-manager/backend/shared/authclient => ../shared/authclient

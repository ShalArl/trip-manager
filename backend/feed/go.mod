module github.com/ShalArl/trip-manager/backend/feed

go 1.25.8

require (
	github.com/ShalArl/trip-manager/backend/shared/authclient v0.0.0-00010101000000-000000000000
	github.com/neo4j/neo4j-go-driver/v5 v5.20.0
)

replace github.com/ShalArl/trip-manager/backend/shared/authclient => ../shared/authclient
module github.com/ShalArl/trip-manager/backend/newsletter

go 1.25.8

require (
	github.com/ShalArl/trip-manager/backend/shared/tenantdb v0.0.0
	github.com/ShalArl/trip-manager/backend/shared/authclient v0.0.0
	github.com/neo4j/neo4j-go-driver/v5 v5.26.0
)

replace (
	github.com/ShalArl/trip-manager/backend/shared/tenantdb => ../shared/tenantdb
	github.com/ShalArl/trip-manager/backend/shared/authclient => ../shared/authclient
)
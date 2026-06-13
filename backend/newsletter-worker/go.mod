module github.com/ShalArl/trip-manager/backend/newsletter-worker

go 1.25.8

require (
	github.com/google/uuid v1.6.0
	github.com/jmoiron/sqlx v1.4.0
	github.com/lib/pq v1.12.3
	github.com/neo4j/neo4j-go-driver/v5 v5.28.4
	github.com/robfig/cron/v3 v3.0.1

	github.com/ShalArl/trip-manager/backend/shared/authclient v0.0.0
	github.com/ShalArl/trip-manager/backend/shared/tenantdb v0.0.0
)

replace (
	github.com/ShalArl/trip-manager/backend/shared/authclient => ../shared/authclient
	github.com/ShalArl/trip-manager/backend/shared/tenantdb => ../shared/tenantdb
)

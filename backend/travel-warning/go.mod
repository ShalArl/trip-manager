module github.com/ShalArl/trip-manager/backend/travel-warning

go 1.25.8

require (
	github.com/ShalArl/trip-manager/backend/shared/authclient v0.0.0
	github.com/ShalArl/trip-manager/backend/shared/middleware v0.0.0
	github.com/redis/go-redis/v9 v9.19.0
)

replace (
	github.com/ShalArl/trip-manager/backend/shared/authclient => ../shared/authclient
	github.com/ShalArl/trip-manager/backend/shared/middleware => ../shared/middleware
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/kelseyhightower/envconfig v1.4.0 // indirect
	go.uber.org/atomic v1.11.0 // indirect
)

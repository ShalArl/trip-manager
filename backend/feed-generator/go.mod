module github.com/ShalArl/trip-manager/backend/feed-generator

go 1.25.8

require (
	github.com/neo4j/neo4j-go-driver/v5 v5.20.0
	github.com/segmentio/kafka-go v0.4.47
	github.com/ShalArl/trip-manager/backend/shared/middleware v0.0.0
)

replace (
	github.com/ShalArl/trip-manager/backend/shared/middleware => ../shared/middleware
)

require (
	github.com/klauspost/compress v1.15.9 // indirect
	github.com/pierrec/lz4/v4 v4.1.15 // indirect
)

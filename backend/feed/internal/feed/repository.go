package feed

import (
	"context"
	"time"

	generated "github.com/ShalArl/trip-manager/backend/feed/generated"
	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

type Repository interface {
	GetFeed(ctx context.Context, limit, offset int) ([]generated.FeedTrip, int, error)
}

type repository struct {
	driver neo4j.DriverWithContext
}

func NewRepository(driver neo4j.DriverWithContext) Repository {
	return &repository{driver: driver}
}

func (r *repository) GetFeed(ctx context.Context, limit, offset int) ([]generated.FeedTrip, int, error) {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	result, err := session.Run(ctx, `
		MATCH (t:Trip)
		OPTIONAL MATCH (t)<-[:CREATED]-(creator:User)
		OPTIONAL MATCH (t)<-[l:LIKED]-()
		OPTIONAL MATCH (t)<-[c:COMMENTED]-()
		WITH t, creator,
		     count(DISTINCT l) AS likes,
		     count(DISTINCT c) AS comments
		WITH t, creator, likes, comments,
		     (likes * 2 + comments * 1) AS score
		ORDER BY score DESC, t.createdAt DESC
		SKIP $offset
		LIMIT $limit
		RETURN t.id        AS tripId,
		       t.title     AS title,
		       t.createdAt AS createdAt,
		       coalesce(creator.id, '') AS creatorId,
		       likes, comments, score
	`, map[string]any{
		"limit":  limit,
		"offset": offset,
	})
	if err != nil {
		return nil, 0, err
	}

	var trips []generated.FeedTrip
	for result.Next(ctx) {
		rec := result.Record()
		trips = append(trips, generated.FeedTrip{
			TripId:    uuidVal(rec, "tripId"),
			Title:     stringVal(rec, "title"),
			CreatedAt: timeVal(rec, "createdAt"),
			CreatorId: uuidVal(rec, "creatorId"),
			Likes:     int64Val(rec, "likes"),
			Comments:  int64Val(rec, "comments"),
			Score:     float32Val(rec, "score"),
		})
	}
	if err := result.Err(); err != nil {
		return nil, 0, err
	}

	countResult, err := session.Run(ctx, `MATCH (t:Trip) RETURN count(t) AS total`, nil)
	if err != nil {
		return nil, 0, err
	}
	total := 0
	if countResult.Next(ctx) {
		if v, ok := countResult.Record().Get("total"); ok {
			if n, ok := v.(int64); ok {
				total = int(n)
			}
		}
	}

	return trips, total, nil
}

// ── Helpers ───────────────────────────────────────────────────────────────────

func stringVal(rec *neo4j.Record, key string) string {
	v, ok := rec.Get(key)
	if !ok || v == nil {
		return ""
	}
	s, _ := v.(string)
	return s
}

func int64Val(rec *neo4j.Record, key string) int64 {
	v, ok := rec.Get(key)
	if !ok || v == nil {
		return 0
	}
	n, _ := v.(int64)
	return n
}

func float32Val(rec *neo4j.Record, key string) float32 {
	v, ok := rec.Get(key)
	if !ok || v == nil {
		return 0
	}
	switch n := v.(type) {
	case float64:
		return float32(n)
	case int64:
		return float32(n)
	}
	return 0
}

func uuidVal(rec *neo4j.Record, key string) openapi_types.UUID {
	s := stringVal(rec, key)
	id, err := uuid.Parse(s)
	if err != nil {
		return openapi_types.UUID{}
	}
	return openapi_types.UUID(id)
}

func timeVal(rec *neo4j.Record, key string) time.Time {
	v, ok := rec.Get(key)
	if !ok || v == nil {
		return time.Time{}
	}
	switch t := v.(type) {
	case time.Time:
		return t
	case string:
		parsed, err := time.Parse(time.RFC3339, t)
		if err != nil {
			return time.Time{}
		}
		return parsed
	}
	return time.Time{}
}

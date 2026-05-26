package feed

import (
	"context"
	"math"
	"time"

	generated "github.com/ShalArl/trip-manager/backend/feed/generated"
	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

type Repository interface {
	GetFeed(ctx context.Context, limit, offset int) ([]generated.FeedTrip, int, error)
	GetPersonalizedFeed(ctx context.Context, userID string, limit, offset int) ([]generated.FeedTrip, int, error)
}

type repository struct {
	driver neo4j.DriverWithContext
}

func NewRepository(driver neo4j.DriverWithContext) Repository {
	return &repository{driver: driver}
}

// GetFeed – globaler Feed mit HackerNews Zeit-Decay
// Score = (likes * 2 + comments) / (alter_in_stunden + 2)^1.5
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
		     toFloat(likes * 2 + comments) / (duration.inSeconds(datetime(t.createdAt), datetime()).seconds / 3600.0 + 2)^1.5 AS score
		ORDER BY score DESC, t.createdAt DESC
		SKIP $offset
		LIMIT $limit
		RETURN t.id        AS tripId,
		       t.title     AS title,
		       t.createdAt AS createdAt,
		       coalesce(creator.id, '') AS creatorId,
		       likes, comments, toFloat(score) AS score
	`, map[string]any{
		"limit":  limit,
		"offset": offset,
	})
	if err != nil {
		return nil, 0, err
	}

	trips, err := collectTrips(ctx, result)
	if err != nil {
		return nil, 0, err
	}

	total, err := countTrips(ctx, session)
	if err != nil {
		return nil, 0, err
	}

	return trips, total, nil
}

// GetPersonalizedFeed – 3-stufiger hybrider Feed
//
// Score = creatorScore * 4 + collaborativeScore * 2 + globalScore * 1
//
// 1. creatorScore:       Trips von Creatorn deren Trips du geliked hast
// 2. collaborativeScore: Trips die ähnliche User geliked haben
// 3. globalScore:        HackerNews Zeit-Decay als Basis für alle Trips
func (r *repository) GetPersonalizedFeed(ctx context.Context, userID string, limit, offset int) ([]generated.FeedTrip, int, error) {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	result, err := session.Run(ctx, `
		MATCH (allTrip:Trip)
		OPTIONAL MATCH (allTrip)<-[:CREATED]-(creator:User)
		OPTIONAL MATCH (allTrip)<-[l:LIKED]-()
		OPTIONAL MATCH (allTrip)<-[c:COMMENTED]-()
		WITH allTrip, creator,
		     count(DISTINCT l) AS likes,
		     count(DISTINCT c) AS comments

		// Global HackerNews Score
		WITH allTrip, creator, likes, comments,
		     toFloat(likes * 2 + comments) / (duration.inSeconds(datetime(allTrip.createdAt), datetime()).seconds / 3600.0 + 2)^1.5 AS globalScore

		// Creator Score: wie oft hat der User Trips von diesem Creator geliked?
		OPTIONAL MATCH (me:User {id: $userId})-[:LIKED]->(:Trip)<-[:CREATED]-(creator)
		WITH allTrip, creator, likes, comments, globalScore,
		     count(me) AS creatorScore

		// Collaborative Score: ähnliche User die diesen Trip geliked haben
		OPTIONAL MATCH (me2:User {id: $userId})-[:LIKED]->(common:Trip)<-[:LIKED]-(similar:User)-[:LIKED]->(allTrip)
		WHERE NOT (me2)-[:LIKED]->(allTrip)
		WITH allTrip, creator, likes, comments, globalScore, creatorScore,
		     count(DISTINCT similar) AS collaborativeScore

		// Hybrider Score
		WITH allTrip, creator, likes, comments,
		     (creatorScore * 4.0 + collaborativeScore * 2.0 + globalScore * 1.0) AS score

		ORDER BY score DESC, allTrip.createdAt DESC
		SKIP $offset
		LIMIT $limit

		RETURN allTrip.id    AS tripId,
		       allTrip.title AS title,
		       allTrip.createdAt AS createdAt,
		       coalesce(creator.id, '') AS creatorId,
		       likes, comments, toFloat(score) AS score
	`, map[string]any{
		"userId": userID,
		"limit":  limit,
		"offset": offset,
	})
	if err != nil {
		// Fallback auf globalen Feed
		return r.GetFeed(ctx, limit, offset)
	}

	trips, err := collectTrips(ctx, result)
	if err != nil {
		return r.GetFeed(ctx, limit, offset)
	}

	total, err := countTrips(ctx, session)
	if err != nil {
		return nil, 0, err
	}

	return trips, total, nil
}

// ── Helpers ───────────────────────────────────────────────────────────────────

func collectTrips(ctx context.Context, result neo4j.ResultWithContext) ([]generated.FeedTrip, error) {
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
	return trips, result.Err()
}

func countTrips(ctx context.Context, session neo4j.SessionWithContext) (int, error) {
	countResult, err := session.Run(ctx, `MATCH (t:Trip) RETURN count(t) AS total`, nil)
	if err != nil {
		return 0, err
	}
	if countResult.Next(ctx) {
		if v, ok := countResult.Record().Get("total"); ok {
			if n, ok := v.(int64); ok {
				return int(n), nil
			}
		}
	}
	return 0, nil
}

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
		if math.IsNaN(n) || math.IsInf(n, 0) {
			return 0
		}
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

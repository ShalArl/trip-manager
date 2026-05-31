package feed

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Repository interface {
	GetFeed(ctx context.Context, limit, offset int) ([]FeedTrip, int, error)
}

type repository struct {
	driver neo4j.DriverWithContext
}

func NewRepository(driver neo4j.DriverWithContext) Repository {
	return &repository{driver: driver}
}

func (r *repository) GetFeed(ctx context.Context, limit, offset int) ([]FeedTrip, int, error) {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	// Score = likes * 2 + comments * 1
	// Neuere Trips bekommen einen kleinen Bonus über createdAt (optional ausbaubar)
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
		RETURN t.id AS tripId,
		       t.title AS title,
		       t.createdAt AS createdAt,
		       coalesce(creator.id, '') AS creatorId,
		       likes,
		       comments,
		       score
	`, map[string]any{
		"limit":  limit,
		"offset": offset,
	})
	if err != nil {
		return nil, 0, err
	}

	var trips []FeedTrip
	for result.Next(ctx) {
		rec := result.Record()
		trips = append(trips, FeedTrip{
			TripID:    stringVal(rec, "tripId"),
			Title:     stringVal(rec, "title"),
			CreatedAt: stringVal(rec, "createdAt"),
			CreatorID: stringVal(rec, "creatorId"),
			Likes:     int64Val(rec, "likes"),
			Comments:  int64Val(rec, "comments"),
			Score:     float64Val(rec, "score"),
		})
	}
	if err := result.Err(); err != nil {
		return nil, 0, err
	}

	// Total Count
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

func float64Val(rec *neo4j.Record, key string) float64 {
	v, ok := rec.Get(key)
	if !ok || v == nil {
		return 0
	}
	switch n := v.(type) {
	case float64:
		return n
	case int64:
		return float64(n)
	}
	return 0
}

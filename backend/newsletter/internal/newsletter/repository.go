package newsletter

import (
	"context"
	"log"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type TripNode struct {
	TripID          string
	Title           string
	CreatorID       string
	CreatorName     string
	LikeCount       int64
	CommentCount    int64
	RelevanceReason string
	CreatedAt       time.Time
}

type Repository interface {
	GetCreatorTrips(ctx context.Context, userID string, limit int) ([]TripNode, error)
	GetSocialGraphTrips(ctx context.Context, userID string, limit int) ([]TripNode, error)
	GetCollaborativeTrips(ctx context.Context, userID string, limit int) ([]TripNode, error)
}

type repository struct {
	driver neo4j.DriverWithContext
}

func NewRepository(driver neo4j.DriverWithContext) Repository {
	return &repository{driver: driver}
}

func (r *repository) GetCreatorTrips(ctx context.Context, userID string, limit int) ([]TripNode, error) {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer func(session neo4j.SessionWithContext, ctx context.Context) {
		err := session.Close(ctx)
		if err != nil {
			log.Printf("error closing session: %v", err)
		}
	}(session, ctx)

	result, err := session.Run(ctx, `
		MATCH (me:User {id: $userID})-[:LIKED|COMMENTED]->(:Trip)<-[:CREATED]-(creator:User)
		WITH me, creator, count(*) AS interactionCount
		MATCH (creator)-[:CREATED]->(t:Trip)
		WHERE NOT (me)-[:LIKED]->(t)
		AND NOT (me)-[:CREATED]->(t)
		OPTIONAL MATCH (t)<-[l:LIKED]-()
		OPTIONAL MATCH (t)<-[c:COMMENTED]-()
		WITH t, creator,
		     count(DISTINCT l) AS likes,
		     count(DISTINCT c) AS comments
		RETURN t.id AS tripId,
		       t.title AS title,
		       t.createdAt AS createdAt,
		       coalesce(creator.id, '') AS creatorId,
		       coalesce(creator.name, '') AS creatorName,
		       likes,
		       comments
		ORDER BY t.createdAt DESC
		LIMIT $limit
	`, map[string]any{
		"userID": userID,
		"limit":  limit,
	})
	if err != nil {
		return nil, err
	}
	return collectTrips(result, "creator_you_follow")
}

func (r *repository) GetSocialGraphTrips(ctx context.Context, userID string, limit int) ([]TripNode, error) {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer func(session neo4j.SessionWithContext, ctx context.Context) {
		err := session.Close(ctx)
		if err != nil {
			log.Printf("error closing session: %v", err)
		}
	}(session, ctx)

	result, err := session.Run(ctx, `
		MATCH (me:User {id: $userID})-[:LIKED]->(myTrip:Trip)<-[:LIKED]-(similar:User)
		WHERE similar.id <> $userID
		WITH me, similar, count(*) AS overlap
		MATCH (similar)-[:LIKED]->(t:Trip)
		WHERE NOT (me)-[:LIKED]->(t)
		AND NOT (me)-[:CREATED]->(t)
		OPTIONAL MATCH (t)<-[:CREATED]-(creator:User)
		OPTIONAL MATCH (t)<-[l:LIKED]-()
		OPTIONAL MATCH (t)<-[c:COMMENTED]-()
		WITH t, creator,
		     count(DISTINCT l) AS likes,
		     count(DISTINCT c) AS comments
		RETURN t.id AS tripId,
		       t.title AS title,
		       t.createdAt AS createdAt,
		       coalesce(creator.id, '') AS creatorId,
		       coalesce(creator.name, '') AS creatorName,
		       likes,
		       comments
		ORDER BY likes DESC, t.createdAt DESC
		LIMIT $limit
	`, map[string]any{
		"userID": userID,
		"limit":  limit,
	})
	if err != nil {
		return nil, err
	}
	return collectTrips(result, "liked_by_similar_users")
}

func (r *repository) GetCollaborativeTrips(ctx context.Context, userID string, limit int) ([]TripNode, error) {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer func(session neo4j.SessionWithContext, ctx context.Context) {
		err := session.Close(ctx)
		if err != nil {
			log.Printf("error closing session: %v", err)
		}
	}(session, ctx)

	result, err := session.Run(ctx, `
		MATCH (me:User {id: $userID})-[:LIKED]->(myTrip:Trip)<-[:LIKED]-(peer:User)
		WHERE peer.id <> $userID
		WITH me, peer
		MATCH (peer)-[:LIKED]->(t:Trip)
		WHERE NOT (me)-[:LIKED]->(t)
		AND NOT (me)-[:CREATED]->(t)
		OPTIONAL MATCH (t)<-[:CREATED]-(creator:User)
		OPTIONAL MATCH (t)<-[l:LIKED]-()
		OPTIONAL MATCH (t)<-[c:COMMENTED]-()
		WITH t, creator,
		     count(DISTINCT l) AS likes,
		     count(DISTINCT c) AS comments
		RETURN t.id AS tripId,
		       t.title AS title,
		       t.createdAt AS createdAt,
		       coalesce(creator.id, '') AS creatorId,
		       coalesce(creator.name, '') AS creatorName,
		       likes,
		       comments
		ORDER BY likes DESC, t.createdAt DESC
		LIMIT $limit
	`, map[string]any{
		"userID": userID,
		"limit":  limit,
	})
	if err != nil {
		return nil, err
	}
	return collectTrips(result, "trending_in_network")
}

func collectTrips(result neo4j.ResultWithContext, reason string) ([]TripNode, error) {
	var trips []TripNode
	for result.Next(context.Background()) {
		rec := result.Record()

		var createdAt time.Time
		if v, ok := rec.Get("createdAt"); ok && v != nil {
			switch t := v.(type) {
			case time.Time:
				createdAt = t
			case string:
				createdAt, _ = time.Parse(time.RFC3339, t)
			}
		}

		var likeCount, commentCount int64
		if v, ok := rec.Get("likes"); ok && v != nil {
			likeCount, _ = v.(int64)
		}
		if v, ok := rec.Get("comments"); ok && v != nil {
			commentCount, _ = v.(int64)
		}

		trips = append(trips, TripNode{
			TripID:          strVal(rec, "tripId"),
			Title:           strVal(rec, "title"),
			CreatorID:       strVal(rec, "creatorId"),
			CreatorName:     strVal(rec, "creatorName"),
			LikeCount:       likeCount,
			CommentCount:    commentCount,
			RelevanceReason: reason,
			CreatedAt:       createdAt,
		})
	}
	return trips, result.Err()
}

func strVal(rec *neo4j.Record, key string) string {
	v, ok := rec.Get(key)
	if !ok || v == nil {
		return ""
	}
	s, _ := v.(string)
	return s
}

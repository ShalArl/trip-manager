package db

import (
	"context"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type TripNode struct {
	TripID          string    `json:"tripId"`
	Title           string    `json:"title"`
	CreatorID       string    `json:"creatorId"`
	CreatorName     string    `json:"creatorName"`
	LikeCount       int64     `json:"likeCount"`
	CommentCount    int64     `json:"commentCount"`
	RelevanceReason string    `json:"relevanceReason"`
	CreatedAt       time.Time `json:"createdAt"`
}

type Section struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Trips       []TripNode `json:"trips"`
}

func GenerateForUser(ctx context.Context, driver neo4j.DriverWithContext, userID string) ([]Section, error) {
	creatorTrips, err := queryCreatorTrips(ctx, driver, userID, 5)
	if err != nil {
		return nil, err
	}

	socialTrips, err := querySocialTrips(ctx, driver, userID, 5)
	if err != nil {
		return nil, err
	}

	collaborativeTrips, err := queryCollaborativeTrips(ctx, driver, userID, 5)
	if err != nil {
		return nil, err
	}

	seen := make(map[string]bool)
	var sections []Section

	if s := buildSection("From Travellers You Follow", "New trips from creators you've interacted with", creatorTrips, seen); len(s.Trips) > 0 {
		sections = append(sections, s)
	}
	if s := buildSection("Popular in Your Network", "Trips liked by travellers with similar tastes", socialTrips, seen); len(s.Trips) > 0 {
		sections = append(sections, s)
	}
	if s := buildSection("Trending Among Your Peers", "Highly liked trips from your travel community", collaborativeTrips, seen); len(s.Trips) > 0 {
		sections = append(sections, s)
	}

	return sections, nil
}

func buildSection(title, description string, trips []TripNode, seen map[string]bool) Section {
	s := Section{Title: title, Description: description, Trips: []TripNode{}}
	for _, t := range trips {
		if seen[t.TripID] {
			continue
		}
		seen[t.TripID] = true
		s.Trips = append(s.Trips, t)
	}
	return s
}

func queryCreatorTrips(ctx context.Context, driver neo4j.DriverWithContext, userID string, limit int) ([]TripNode, error) {
	session := driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

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
	`, map[string]any{"userID": userID, "limit": limit})
	if err != nil {
		return nil, err
	}
	return collectTrips(result, "creator_you_follow")
}

func querySocialTrips(ctx context.Context, driver neo4j.DriverWithContext, userID string, limit int) ([]TripNode, error) {
	session := driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

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
	`, map[string]any{"userID": userID, "limit": limit})
	if err != nil {
		return nil, err
	}
	return collectTrips(result, "liked_by_similar_users")
}

func queryCollaborativeTrips(ctx context.Context, driver neo4j.DriverWithContext, userID string, limit int) ([]TripNode, error) {
	session := driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

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
	`, map[string]any{"userID": userID, "limit": limit})
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

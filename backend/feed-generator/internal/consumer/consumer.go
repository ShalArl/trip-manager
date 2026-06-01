package consumer

import (
	"context"
	"encoding/json"
	"log"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/segmentio/kafka-go"
)

const (
	TopicTripCreated   = "trip.created"
	TopicTripLiked     = "trip.liked"
	TopicTripCommented = "trip.commented"
)

// ── Event Types ───────────────────────────────────────────────────────────────

type TripCreatedEvent struct {
	TripID    string `json:"tripId"`
	UserID    string `json:"userId"`
	UserName  string `json:"userName"`
	Title     string `json:"title"`
	CreatedAt string `json:"createdAt"`
}

type TripLikedEvent struct {
	TripID    string `json:"tripId"`
	UserID    string `json:"userId"`
	CreatedAt string `json:"createdAt"`
}

type TripCommentedEvent struct {
	TripID    string `json:"tripId"`
	UserID    string `json:"userId"`
	CreatedAt string `json:"createdAt"`
}

// ── Consumer ──────────────────────────────────────────────────────────────────

type Consumer struct {
	driver  neo4j.DriverWithContext
	brokers []string
	groupID string
}

func New(driver neo4j.DriverWithContext, brokers []string, groupID string) *Consumer {
	return &Consumer{
		driver:  driver,
		brokers: brokers,
		groupID: groupID,
	}
}

// Start startet alle drei Topic-Consumer in eigenen Goroutinen.
// Blockiert bis ctx cancelled wird.
func (c *Consumer) Start(ctx context.Context) {
	topics := []string{TopicTripCreated, TopicTripLiked, TopicTripCommented}
	for _, topic := range topics {
		go c.consume(ctx, topic)
	}
	<-ctx.Done()
	log.Println("feed-generator: shutting down consumers")
}

func (c *Consumer) consume(ctx context.Context, topic string) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     c.brokers,
		GroupID:     c.groupID,
		Topic:       topic,
		MinBytes:    1,
		MaxBytes:    10e6,
		StartOffset: kafka.FirstOffset, // ← hinzufügen
	})
	defer r.Close()

	log.Printf("feed-generator: consuming topic %s", topic)

	for {
		msg, err := r.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			log.Printf("feed-generator: error reading from %s: %v", topic, err)
			continue
		}
		log.Printf("feed-generator: received message from %s: %s", topic, string(msg.Value))

		if err := c.handle(ctx, topic, msg.Value); err != nil {
			log.Printf("feed-generator: error handling message from %s: %v", topic, err)
		}

		r.CommitMessages(ctx, msg)
	}
}

func (c *Consumer) handle(ctx context.Context, topic string, payload []byte) error {
	switch topic {
	case TopicTripCreated:
		var e TripCreatedEvent
		if err := json.Unmarshal(payload, &e); err != nil {
			return err
		}
		return c.onTripCreated(ctx, e)

	case TopicTripLiked:
		var e TripLikedEvent
		if err := json.Unmarshal(payload, &e); err != nil {
			return err
		}
		return c.onTripLiked(ctx, e)

	case TopicTripCommented:
		var e TripCommentedEvent
		if err := json.Unmarshal(payload, &e); err != nil {
			return err
		}
		return c.onTripCommented(ctx, e)
	}
	return nil
}

// ── Neo4j Writes ──────────────────────────────────────────────────────────────

// onTripCreated schreibt:
//
//	(User {id, name}) -[:CREATED {createdAt}]-> (Trip {id, title, createdAt})
func (c *Consumer) onTripCreated(ctx context.Context, e TripCreatedEvent) error {
	_, err := neo4j.ExecuteQuery(ctx, c.driver,
		`MERGE (u:User {id: $userId})
		 SET u.name = $userName
		 MERGE (t:Trip {id: $tripId})
		 SET t.title = $title, t.createdAt = $createdAt
		 MERGE (u)-[r:CREATED]->(t)
		 SET r.createdAt = $createdAt`,
		map[string]any{
			"userId":    e.UserID,
			"userName":  e.UserName,
			"tripId":    e.TripID,
			"title":     e.Title,
			"createdAt": e.CreatedAt,
		},
		neo4j.EagerResultTransformer,
	)
	if err != nil {
		return err
	}
	log.Printf("feed-generator: CREATED %s -> %s", e.UserID, e.TripID)
	return nil
}

// onTripLiked schreibt:
//
//	(User {id}) -[:LIKED {createdAt}]-> (Trip {id})
//
// MERGE auf Trip damit auch Likes ankommen bevor trip.created verarbeitet wurde.
func (c *Consumer) onTripLiked(ctx context.Context, e TripLikedEvent) error {
	_, err := neo4j.ExecuteQuery(ctx, c.driver,
		`MERGE (u:User {id: $userId})
		 MERGE (t:Trip {id: $tripId})
		 MERGE (u)-[r:LIKED]->(t)
		 SET r.createdAt = $createdAt`,
		map[string]any{
			"userId":    e.UserID,
			"tripId":    e.TripID,
			"createdAt": e.CreatedAt,
		},
		neo4j.EagerResultTransformer,
	)
	if err != nil {
		return err
	}
	log.Printf("feed-generator: LIKED %s -> %s", e.UserID, e.TripID)
	return nil
}

// onTripCommented schreibt:
//
//	(User {id}) -[:COMMENTED {createdAt}]-> (Trip {id})
func (c *Consumer) onTripCommented(ctx context.Context, e TripCommentedEvent) error {
	_, err := neo4j.ExecuteQuery(ctx, c.driver,
		`MERGE (u:User {id: $userId})
		 MERGE (t:Trip {id: $tripId})
		 MERGE (u)-[r:COMMENTED]->(t)
		 SET r.createdAt = $createdAt`,
		map[string]any{
			"userId":    e.UserID,
			"tripId":    e.TripID,
			"createdAt": e.CreatedAt,
		},
		neo4j.EagerResultTransformer,
	)
	if err != nil {
		return err
	}
	log.Printf("feed-generator: COMMENTED %s -> %s", e.UserID, e.TripID)
	return nil
}

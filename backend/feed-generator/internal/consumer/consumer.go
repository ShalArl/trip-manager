package consumer

import (
	"context"
	"encoding/json"
	"log"

	"cloud.google.com/go/pubsub"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

const (
	AttrEventType          = "event_type"
	EventTypeTripCreated   = "trip.created"
	EventTypeTripLiked     = "trip.liked"
	EventTypeTripCommented = "trip.commented"
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
	driver         neo4j.DriverWithContext
	client         *pubsub.Client
	subscriptionID string
}

func New(driver neo4j.DriverWithContext, projectID, subscriptionID string) (*Consumer, error) {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		driver:         driver,
		client:         client,
		subscriptionID: subscriptionID,
	}, nil
}

func (c *Consumer) Start(ctx context.Context) {
	log.Printf("feed-generator: starting pubsub consumer for subscription %s", c.subscriptionID)

	sub := c.client.Subscription(c.subscriptionID)

	err := sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		eventType := msg.Attributes[AttrEventType]
		log.Printf("feed-generator: received message with type %s", eventType)

		if err := c.handle(ctx, eventType, msg.Data); err != nil {
			log.Printf("feed-generator: error handling event %s: %v", eventType, err)
			msg.Nack()
			return
		}

		msg.Ack()
	})

	if err != nil && ctx.Err() == nil {
		log.Fatalf("feed-generator: consumer runtime error: %v", err)
	}

	log.Println("feed-generator: shutting down consumer")
}

func (c *Consumer) Close() error {
	return c.client.Close()
}

func (c *Consumer) handle(ctx context.Context, eventType string, payload []byte) error {
	switch eventType {
	case EventTypeTripCreated:
		var e TripCreatedEvent
		if err := json.Unmarshal(payload, &e); err != nil {
			return err
		}
		return c.onTripCreated(ctx, e)

	case EventTypeTripLiked:
		var e TripLikedEvent
		if err := json.Unmarshal(payload, &e); err != nil {
			return err
		}
		return c.onTripLiked(ctx, e)

	case EventTypeTripCommented:
		var e TripCommentedEvent
		if err := json.Unmarshal(payload, &e); err != nil {
			return err
		}
		return c.onTripCommented(ctx, e)
	}

	log.Printf("feed-generator: warn: unknown event type %s", eventType)
	return nil
}

// ── Neo4j Writes ──────────────────────────────────────────────────────────────

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
	return err
}

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
	return err
}

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
	return err
}

package pubsub

import (
	"context"
	"encoding/json"
	"fmt"

	"cloud.google.com/go/pubsub"
)

const (
	AttrEventType          = "event_type"
	EventTypeTripLiked     = "trip.liked"
	EventTypeTripCommented = "trip.commented"
)

type Producer struct {
	client *pubsub.Client
	topic  *pubsub.Topic
}

func NewProducer(projectID, topicID string) (*Producer, error) {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to create pubsub client: %w", err)
	}

	return &Producer{
		client: client,
		topic:  client.Topic(topicID),
	}, nil
}

func (p *Producer) Close() error {
	p.topic.Stop()
	return p.client.Close()
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

func (p *Producer) PublishTripLiked(ctx context.Context, e TripLikedEvent) error {
	payload, err := json.Marshal(e)
	if err != nil {
		return err
	}

	result := p.topic.Publish(ctx, &pubsub.Message{
		Data: payload,
		Attributes: map[string]string{
			AttrEventType: EventTypeTripLiked,
		},
		OrderingKey: e.TripID,
	})

	_, err = result.Get(ctx)
	return err
}

func (p *Producer) PublishTripCommented(ctx context.Context, e TripCommentedEvent) error {
	payload, err := json.Marshal(e)
	if err != nil {
		return err
	}

	result := p.topic.Publish(ctx, &pubsub.Message{
		Data: payload,
		Attributes: map[string]string{
			AttrEventType: EventTypeTripCommented,
		},
		OrderingKey: e.TripID,
	})

	_, err = result.Get(ctx)
	return err
}

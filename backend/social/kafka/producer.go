package kafka

import (
	"context"
	"encoding/json"
	"time"

	"github.com/segmentio/kafka-go"
)

const (
	TopicTripLiked     = "trip.liked"
	TopicTripCommented = "trip.commented"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(brokers []string) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(brokers...),
			Balancer:     &kafka.LeastBytes{},
			WriteTimeout: 5 * time.Second,
			ReadTimeout:  5 * time.Second,
			Async:        true,
		},
	}
}

func (p *Producer) Close() error {
	return p.writer.Close()
}

// TripLikedEvent wird bei trip.liked gepublished
type TripLikedEvent struct {
	TripID    string `json:"tripId"`
	UserID    string `json:"userId"` // Firebase UID (wie in Firestore gespeichert)
	CreatedAt string `json:"createdAt"`
}

// TripCommentedEvent wird bei trip.commented gepublished
type TripCommentedEvent struct {
	TripID    string `json:"tripId"`
	UserID    string `json:"userId"` // Firebase UID
	CreatedAt string `json:"createdAt"`
}

func (p *Producer) PublishTripLiked(ctx context.Context, e TripLikedEvent) error {
	payload, err := json.Marshal(e)
	if err != nil {
		return err
	}
	return p.writer.WriteMessages(ctx, kafka.Message{
		Topic: TopicTripLiked,
		Key:   []byte(e.TripID + ":" + e.UserID),
		Value: payload,
	})
}

func (p *Producer) PublishTripCommented(ctx context.Context, e TripCommentedEvent) error {
	payload, err := json.Marshal(e)
	if err != nil {
		return err
	}
	return p.writer.WriteMessages(ctx, kafka.Message{
		Topic: TopicTripCommented,
		Key:   []byte(e.TripID + ":" + e.UserID),
		Value: payload,
	})
}

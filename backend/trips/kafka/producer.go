package kafka

import (
	"context"
	"encoding/json"
	"time"

	"github.com/segmentio/kafka-go"
)

const TopicTripCreated = "trip.created"

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
			// Async so der HTTP-Request nicht blockiert wird
			Async: true,
		},
	}
}

func (p *Producer) Close() error {
	return p.writer.Close()
}

// TripCreatedEvent wird bei trip.created gepublished
type TripCreatedEvent struct {
	TripID    string `json:"tripId"`
	UserID    string `json:"userId"` // PostgreSQL UUID
	UserName  string `json:"userName"`
	Title     string `json:"title"`
	CreatedAt string `json:"createdAt"`
}

func (p *Producer) PublishTripCreated(ctx context.Context, e TripCreatedEvent) error {
	payload, err := json.Marshal(e)
	if err != nil {
		return err
	}
	return p.writer.WriteMessages(ctx, kafka.Message{
		Topic: TopicTripCreated,
		Key:   []byte(e.TripID),
		Value: payload,
	})
}

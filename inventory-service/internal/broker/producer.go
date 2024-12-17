package broker

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"log"
	"time"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(address string, topic string) *Producer {
	writer := &kafka.Writer{
		Addr:         kafka.TCP(address),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		WriteTimeout: 10 * time.Second,
		RequiredAcks: kafka.RequireAll,
	}

	return &Producer{writer: writer}
}

func (p *Producer) publish(ctx context.Context, key string, value []byte) error {
	err := p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(key),
		Value: value,
	})

	if err != nil {
		log.Printf("Failed to publish message: %v", err)
		return err
	}

	log.Printf("Message published to Kafka: key=%s value=%s", key, string(value))
	return nil
}

func (p *Producer) PublishEvent(ctx context.Context, key string, event map[string]interface{}) error {
	eventBytes, err := json.Marshal(event)
	if err != nil {
		log.Printf("Failed to marshal event: %v", err)
		return err
	}

	err = p.publish(ctx, key, eventBytes)
	if err != nil {
		log.Printf("Failed to publish event: %v", err)
		return err
	}

	log.Printf("Event published to Kafka: key=%s event=%v", key, event)
	return nil
}

func (p *Producer) Close() {
	if err := p.writer.Close(); err != nil {
		log.Printf("Error closing Kafka writer: %v", err)
	}
}

package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
)

type Producer struct {
	writer *kafka.Writer
}

type OutOfStockEvent struct {
	Action  string `json:"action"`
	ItemID  int    `json:"item_id"`
	UserID  string `json:"user_id"`
	Message string `json:"message"`
}

func NewProducer(address string, topic string) *Producer {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(address),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}

	return &Producer{writer: writer}
}

func (p *Producer) Produce(ctx context.Context, itemID int, userID string) error {
	event := OutOfStockEvent{
		Action:  "out-of-stock",
		ItemID:  itemID,
		UserID:  userID,
		Message: fmt.Sprintf("Item %d is out of stock", itemID),
	}

	data, err := json.Marshal(event)
	if err != nil {
		log.Printf("Failed to marshal out-of-stock event: %v\n", err)
		return err
	}

	err = p.writer.WriteMessages(ctx, kafka.Message{Value: data})
	if err != nil {
		log.Printf("Failed to write message to Kafka: %v\n", err)
		return err
	}

	log.Printf("Out-of-stock message sent for item %d and user %s: %s\n", itemID, userID, string(data))
	return nil
}

func (p *Producer) Close() error {
	return p.writer.Close()
}

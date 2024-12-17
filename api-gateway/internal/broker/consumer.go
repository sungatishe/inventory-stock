package broker

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"log"
)

type OutOfStockEvent struct {
	Action  string `json:"action"`
	ItemID  int    `json:"item_id"`
	UserID  string `json:"user_id"`
	Message string `json:"message"`
}

type Consumer struct {
	reader *kafka.Reader
}

func NewConsumer(brokers []string, groupID, topic string) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		GroupID: groupID,
		Topic:   topic,
	})

	return &Consumer{reader: reader}
}

func (c *Consumer) Consume(ctx context.Context, handler func(event OutOfStockEvent)) error {
	for {
		log.Println("Waiting for a Kafka message...")
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				log.Println("Context cancelled, stopping consumer.")
				return nil
			}
			log.Printf("Error reading message: %v\n", err)
			continue
		}

		log.Printf("Message received: %s\n", string(msg.Value))

		var event OutOfStockEvent
		err = json.Unmarshal(msg.Value, &event)
		if err != nil {
			log.Printf("Failed to unmarshal message: %v\n", err)
			continue
		}

		log.Printf("Processing event: %+v\n", event)
		handler(event)
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}

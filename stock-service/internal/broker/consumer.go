package broker

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"log"
)

type StockUpdate struct {
	Action string    `json:"action"`
	Item   StockItem `json:"item"`
}

type StockItem struct {
	ID       int64   `json:"id"`
	Name     string  `json:"name"`
	Quantity int64   `json:"quantity"`
	Price    float64 `json:"price"`
	UserID   string  `json:"user_id"`
}

type Consumer struct {
	reader *kafka.Reader
}

func NewConsumer(brokers []string, topic, groupID string) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: groupID,
	})

	return &Consumer{reader: reader}
}

func (c *Consumer) Consume(ctx context.Context, handler func(StockUpdate)) error {
	defer c.reader.Close()

	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				log.Println("Context canceled, stopping Kafka consumer...")
				return nil
			}
			log.Printf("Failed to read message: %v\n", err)
			continue
		}

		log.Printf("Received message: %s\n", msg)
		var update StockUpdate
		err = json.Unmarshal(msg.Value, &update)
		if err != nil {
			log.Printf("Failed to unmarshal message: %v\n", err)
			continue
		}
		log.Printf("Parsed update: %+v\n", update)

		handler(update)
	}
}

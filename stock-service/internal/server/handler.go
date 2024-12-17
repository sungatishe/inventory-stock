package server

import (
	"context"
	"fmt"
	"log"
	"stock-service/internal/broker"
)

func (s *StockProcessorServer) processKafkaEvents() {
	err := s.consumer.Consume(s.ctx, func(update broker.StockUpdate) {
		fmt.Println("Message: ", update)
		s.processStockUpdate(update)
	})

	if err != nil {
		log.Printf("Error consuming Kafka messages: %v\n", err)
	}
}

func (s *StockProcessorServer) processStockUpdate(update broker.StockUpdate) {
	item := update.Item
	_, err := s.redis.GetItemQuantity(item.ID)
	if err != nil {
		if err.Error() == "redis: nil" {
			err = s.redis.SetItemQuantity(item.ID, int(item.Quantity))
			if err != nil {
				log.Printf("Failed to initialize quantity for item %v: %v\n", item.ID, err)
				return
			}
		} else {
			log.Printf("Failed to get current quantity for item %v: %v\n", item.ID, err)
			return
		}
	}

	if item.Quantity <= 0 {
		log.Printf("Warning: New quantity for item %v is less than one (%d). Skipping update.\n", item.ID, item.Quantity)

		err := s.producer.Produce(context.Background(), int(item.ID), item.UserID)
		if err != nil {
			log.Printf("Failed to publish out-of-stock event: %v\n", err)
			return
		}

		return
	}

	log.Printf("Updated quantity for item %v: %d\n", item.ID, item.Quantity)

}

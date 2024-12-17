package app

import (
	"log"
	"stock-service/internal/broker"
	"stock-service/internal/cache"
	"stock-service/internal/server"
)

func Run() {
	kafkaBrokers := []string{"kafka:9092"}
	kafkaAddress := "kafka:9092"
	kafkaTopicConsume := "stock-updates"
	kafkaTopicProduce := "out-of-stock"
	kafkaGroupID := "stock-processor-group"
	redisAddress := "redis:6379"

	redisClient := cache.NewRedisClient(redisAddress)

	consumer := broker.NewConsumer(kafkaBrokers, kafkaTopicConsume, kafkaGroupID)

	producer := broker.NewProducer(kafkaAddress, kafkaTopicProduce)

	defer func() {
		if err := producer.Close(); err != nil {
			log.Printf("Failed to close producer: %v\n", err)
		}
	}()

	stockProcessorServer := server.NewStockProcessorServer(consumer, producer, redisClient)
	stockProcessorServer.Run()
}

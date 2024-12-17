package app

import (
	db "inventory-service/db/config"
	"inventory-service/internal/broker"
	"inventory-service/internal/repository"
	"inventory-service/internal/server"
	"inventory-service/internal/service"
	"os"
)

func Run() {
	port := os.Getenv("PORT")
	kafkaAddress := os.Getenv("KAFKA_BROKER")
	topic := os.Getenv("KAFKA_TOPIC")

	producer := broker.NewProducer(kafkaAddress, topic)
	defer producer.Close()

	db.InitDB()

	inventoryRepo := repository.NewInventoryRepository(db.DB)
	inventoryService := service.NewItemService(inventoryRepo, producer)

	inventoryServer := server.NewInventoryServer(inventoryService)

	inventoryServer.Run(port)
}

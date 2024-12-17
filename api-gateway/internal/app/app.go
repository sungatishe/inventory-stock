package app

import (
	"api-gateway/internal/broker"
	"api-gateway/internal/clients"
	"api-gateway/internal/handlers"
	"api-gateway/internal/routes"
	"api-gateway/internal/server"
	"api-gateway/internal/ws"
	"context"
	"github.com/go-chi/chi/v5"
	"log"
	"os"
)

func Run() {
	kafkaBrokers := []string{"kafka:9092"}
	topic := "out-of-stock"
	groupID := "api-stock-processor-group"

	consumer := broker.NewConsumer(kafkaBrokers, groupID, topic)
	defer func() {
		if err := consumer.Close(); err != nil {
			log.Printf("Failed to close Kafka consumer: %v\n", err)
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wsManager := ws.NewManager()

	go func() {
		log.Println("Starting Kafka consumer...")
		err := consumer.Consume(ctx, broker.HandleOutOfStockEvent(wsManager))
		if err != nil {
			log.Printf("Error consuming Kafka messages: %v\n", err)
		}
	}()

	authClient := clients.NewAuthClient(os.Getenv("AUTH_SERVICE_URL"))
	inventoryClient := clients.NewInventoryClient(os.Getenv("INVENTORY_SERVICE_URL"))

	authHandler := handlers.NewAuthHandler(authClient)
	inventoryHandler := handlers.NewInventoryHandler(inventoryClient)
	stockHandler := handlers.NewStockHandlers(wsManager)

	r := chi.NewRouter()

	rt := routes.NewRouter(r)
	rt.SetupAuthRoutes(authHandler)
	rt.SetupInventoryRoutes(authClient, inventoryHandler)
	rt.SetupStockRoutes(stockHandler)

	apiServer := server.NewServer(os.Getenv("PORT"), r)

	apiServer.RunServer()
}

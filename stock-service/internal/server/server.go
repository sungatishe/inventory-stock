package server

import (
	"context"
	"log"
	"stock-service/internal/broker"
	"stock-service/internal/cache"
	"sync"
)

type StockProcessorServer struct {
	consumer   *broker.Consumer
	producer   *broker.Producer
	redis      *cache.Client
	ctx        context.Context
	cancelFunc context.CancelFunc
	wg         *sync.WaitGroup
}

func NewStockProcessorServer(consumer *broker.Consumer, producer *broker.Producer, redis *cache.Client) *StockProcessorServer {
	ctx, cancel := context.WithCancel(context.Background())
	return &StockProcessorServer{
		consumer:   consumer,
		producer:   producer,
		redis:      redis,
		ctx:        ctx,
		cancelFunc: cancel,
		wg:         &sync.WaitGroup{},
	}
}

func (s *StockProcessorServer) Run() {
	stopChan := SetupSignalHandler()

	log.Println("Starting Stock Processor Service...")

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		s.processKafkaEvents()
	}()

	<-stopChan
	log.Println("Shutting down Stock Processor Service...")

	s.Shutdown()
}

func (s *StockProcessorServer) Shutdown() {
	s.cancelFunc()
	s.wg.Wait()
	log.Println("All processes stopped.")
}

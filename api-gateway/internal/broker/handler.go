package broker

import (
	"api-gateway/internal/ws"
	"log"
	"strconv"
)

func HandleOutOfStockEvent(manager *ws.Manager) func(event OutOfStockEvent) {
	return func(event OutOfStockEvent) {
		log.Printf("Received out-of-stock event: ItemID=%d, UserID=%d, Message=%s\n", event.ItemID, event.UserID, event.Message)

		userIdStr, err := strconv.Atoi(event.UserID)
		if err != nil {
			panic(err)
		}
		// Отправляем уведомление через WebSocket
		manager.SendNotification(userIdStr, event.Message)
	}
}

package cache

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
)

var ctx = context.Background()

type Client struct {
	client *redis.Client
}

func NewRedisClient(addr string) *Client {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}
	return &Client{client: client}
}

func (c *Client) GetItemQuantity(itemID int64) (int, error) {
	key := fmt.Sprintf("inventory:%v", itemID)

	quantity, err := c.client.Get(ctx, key).Int()
	if err != nil {
		if err == redis.Nil {
			return 0, err
		}
		return 0, fmt.Errorf("failed to get quantity of items %v: %w", itemID, quantity)
	}

	return quantity, nil
}

func (c *Client) SetItemQuantity(itemID int64, quantity int) error {
	key := fmt.Sprintf("inventory:%v", itemID)

	err := c.client.Set(ctx, key, quantity, 0).Err()
	if err != nil {
		return fmt.Errorf("failed to set quantity for item %v: %w", itemID, err)
	}

	return nil
}

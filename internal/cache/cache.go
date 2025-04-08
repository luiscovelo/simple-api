package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	client *redis.Client
}

func New() *Cache {
	client := redis.NewClient(&redis.Options{
		Addr:     ":6379",
		Password: "",
		DB:       0,
	})

	return &Cache{client: client}
}

func (c *Cache) Ping() error {
	_, err := c.client.Ping(context.Background()).Result()
	if err != nil {
		return err
	}

	return nil
}

func (c *Cache) Set(ctx context.Context, key string, value any) error {
	_, err := c.client.Set(ctx, key, value, 5*time.Second).Result()
	if err != nil {
		return err
	}

	return nil
}

func (c *Cache) Get(ctx context.Context, key string) string {
	val, _ := c.client.Get(ctx, key).Result()

	if val == "" {
		return ""
	}

	return val
}

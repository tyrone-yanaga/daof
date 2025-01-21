package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type Client struct {
	client *redis.Client
}

type Config struct {
	Host     string
	Port     string
	Password string
	DB       int
}

func NewClient(config Config) (*Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.Host, config.Port),
		Password: config.Password,
		DB:       config.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return &Client{client: client}, nil
}

// Cache methods
func (c *Client) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	bytes, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	if err := c.client.Set(ctx, key, bytes, expiration).Err(); err != nil {
		return fmt.Errorf("failed to set cache: %w", err)
	}

	return nil
}

func (c *Client) Get(ctx context.Context, key string, dest interface{}) error {
	bytes, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil // Key doesn't exist
		}
		return fmt.Errorf("failed to get cache: %w", err)
	}

	if err := json.Unmarshal(bytes, dest); err != nil {
		return fmt.Errorf("failed to unmarshal value: %w", err)
	}

	return nil
}

func (c *Client) Delete(ctx context.Context, key string) error {
	result := c.client.Del(ctx, key)
	if err := result.Err(); err != nil {
		return fmt.Errorf("failed to delete key %s: %w", key, err)
	}
	return nil
}

// Session methods
func (c *Client) SetSession(ctx context.Context, sessionID string, data interface{}, expiration time.Duration) error {
	return c.Set(ctx, fmt.Sprintf("session:%s", sessionID), data, expiration)
}

func (c *Client) GetSession(ctx context.Context, sessionID string, dest interface{}) error {
	return c.Get(ctx, fmt.Sprintf("session:%s", sessionID), dest)
}

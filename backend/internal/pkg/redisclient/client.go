package redisclient

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// Client wraps the redis client to offer higher level helpers for caching sessions.
type Client struct {
	inner *redis.Client
}

// New builds a new redis helper around the provided go-redis client.
func New(inner *redis.Client) *Client {
	return &Client{inner: inner}
}

// Remember sets a key with ttl when value is not empty.
func (c *Client) Remember(ctx context.Context, key string, payload string, ttl time.Duration) error {
	if payload == "" {
		return nil
	}
	return c.inner.Set(ctx, key, payload, ttl).Err()
}

// Fetch returns the cached string if any.
func (c *Client) Fetch(ctx context.Context, key string) (string, error) {
	val, err := c.inner.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	return val, err
}

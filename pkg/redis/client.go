package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Client wraps the go-redis client with helper methods.
type Client struct {
	rdb *redis.Client
}

// Config holds Redis connection settings.
type Config struct {
	Host     string
	Port     string
	Password string
	DB       int
}

// New creates and returns a connected Redis Client.
// It performs a PING to verify the connection.
func New(cfg Config) (*Client, error) {
	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis: failed to connect to %s: %w", addr, err)
	}

	return &Client{rdb: rdb}, nil
}

// Set stores key-value with optional TTL.
// Pass 0 for ttl to store without expiry.
func (c *Client) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("redis: failed to marshal value: %w", err)
	}
	return c.rdb.Set(ctx, key, data, ttl).Err()
}

// Get retrieves a value and unmarshals it into dest.
// Returns ErrCacheMiss if the key does not exist.
func (c *Client) Get(ctx context.Context, key string, dest interface{}) error {
	data, err := c.rdb.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return ErrCacheMiss
		}
		return fmt.Errorf("redis: GET %s: %w", key, err)
	}
	if err := json.Unmarshal(data, dest); err != nil {
		return fmt.Errorf("redis: failed to unmarshal value for key %s: %w", key, err)
	}
	return nil
}

// Delete removes one or more keys.
func (c *Client) Delete(ctx context.Context, keys ...string) error {
	return c.rdb.Del(ctx, keys...).Err()
}

// Exists checks whether a key exists.
func (c *Client) Exists(ctx context.Context, key string) (bool, error) {
	n, err := c.rdb.Exists(ctx, key).Result()
	return n > 0, err
}

// SetNX sets key only if it does NOT exist (used for distributed locks).
func (c *Client) SetNX(ctx context.Context, key string, value interface{}, ttl time.Duration) (bool, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return false, fmt.Errorf("redis: failed to marshal value: %w", err)
	}
	return c.rdb.SetNX(ctx, key, data, ttl).Result()
}

// Ping verifies the connection is alive.
func (c *Client) Ping(ctx context.Context) error {
	return c.rdb.Ping(ctx).Err()
}

// Close gracefully shuts down the Redis connection.
func (c *Client) Close() error {
	return c.rdb.Close()
}

// ErrCacheMiss is returned when a key is not found in Redis.
var ErrCacheMiss = fmt.Errorf("redis: cache miss")

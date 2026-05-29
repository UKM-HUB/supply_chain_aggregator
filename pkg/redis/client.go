package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Client wraps the go-redis client dengan helper methods.
type Client struct {
	rdb *redis.Client
}

// Config menyimpan konfigurasi koneksi Redis.
type Config struct {
	Host     string
	Port     string
	Password string
	DB       int
}

// New membuat dan mengembalikan Redis Client yang sudah terhubung.
// Melakukan PING untuk verifikasi koneksi.
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
		return nil, fmt.Errorf("redis: gagal connect ke %s: %w", addr, err)
	}

	return &Client{rdb: rdb}, nil
}

// Set menyimpan key-value dengan optional TTL.
// Kirim 0 untuk TTL agar tersimpan tanpa expiry.
func (c *Client) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("redis: marshal error: %w", err)
	}
	return c.rdb.Set(ctx, key, data, ttl).Err()
}

// Get mengambil value dan unmarshal ke dest.
// Mengembalikan ErrCacheMiss jika key tidak ada.
func (c *Client) Get(ctx context.Context, key string, dest interface{}) error {
	data, err := c.rdb.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return ErrCacheMiss
		}
		return fmt.Errorf("redis: GET %s: %w", key, err)
	}
	if err := json.Unmarshal(data, dest); err != nil {
		return fmt.Errorf("redis: unmarshal error untuk key %s: %w", key, err)
	}
	return nil
}

// Delete menghapus satu atau lebih key.
// Key yang tidak ada diabaikan (tidak error).
func (c *Client) Delete(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}
	return c.rdb.Del(ctx, keys...).Err()
}

// Exists mengecek apakah key ada di Redis.
func (c *Client) Exists(ctx context.Context, key string) (bool, error) {
	n, err := c.rdb.Exists(ctx, key).Result()
	return n > 0, err
}

// SetNX menyimpan key hanya jika BELUM ada (dipakai untuk distributed lock).
// Mengembalikan true jika berhasil di-set.
func (c *Client) SetNX(ctx context.Context, key string, value interface{}, ttl time.Duration) (bool, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return false, fmt.Errorf("redis: marshal error: %w", err)
	}
	return c.rdb.SetNX(ctx, key, data, ttl).Result()
}

// Ping memverifikasi koneksi masih hidup.
func (c *Client) Ping(ctx context.Context) error {
	return c.rdb.Ping(ctx).Err()
}

// Close menutup koneksi Redis secara graceful.
func (c *Client) Close() error {
	return c.rdb.Close()
}

// ErrCacheMiss dikembalikan ketika key tidak ditemukan di Redis.
var ErrCacheMiss = fmt.Errorf("redis: cache miss")

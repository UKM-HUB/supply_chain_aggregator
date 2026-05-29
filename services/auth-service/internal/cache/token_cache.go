package cache

import (
	"context"
	"fmt"
	"time"

	pkgredis "supply-chain-aggregator/pkg/redis"
)

const (
	// tokenTTL is how long a validated token result is cached.
	tokenTTL = 5 * time.Minute
	keyPrefix = "auth:token:"
)

// TokenCache caches token-validation results in Redis.
type TokenCache struct {
	redis *pkgredis.Client
}

// TokenClaims holds cached claims data.
type TokenClaims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
}

// NewTokenCache creates a new TokenCache backed by Redis.
func NewTokenCache(redis *pkgredis.Client) *TokenCache {
	return &TokenCache{redis: redis}
}

// Set caches claims for a given token.
func (c *TokenCache) Set(ctx context.Context, token string, claims TokenClaims) error {
	key := keyPrefix + token
	return c.redis.Set(ctx, key, claims, tokenTTL)
}

// Get retrieves cached claims for a token.
// Returns (nil, nil) on cache miss.
func (c *TokenCache) Get(ctx context.Context, token string) (*TokenClaims, error) {
	key := keyPrefix + token
	var claims TokenClaims
	err := c.redis.Get(ctx, key, &claims)
	if err == pkgredis.ErrCacheMiss {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("token cache get: %w", err)
	}
	return &claims, nil
}

// Invalidate removes a token from the cache (e.g. on logout).
func (c *TokenCache) Invalidate(ctx context.Context, token string) error {
	return c.redis.Delete(ctx, keyPrefix+token)
}

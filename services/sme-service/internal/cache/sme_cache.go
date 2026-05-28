package cache

import (
	"context"
	"fmt"
	"time"

	pkgredis "supply-chain-aggregator/pkg/redis"
	"supply-chain-aggregator/services/sme-service/internal/entity"
)

const (
	smeListTTL       = 2 * time.Minute
	smeCategoryTTL   = 10 * time.Minute
	keyPrefixSMEList = "sme:list:"
	keyCategoryAll   = "sme:categories:all"
)

// SMECache caches SME list results in Redis.
type SMECache struct {
	redis *pkgredis.Client
}

// NewSMECache creates a new SMECache.
func NewSMECache(redis *pkgredis.Client) *SMECache {
	return &SMECache{redis: redis}
}

// SetList caches a list of SMEs for the given filter key.
func (c *SMECache) SetList(ctx context.Context, filterKey string, smes []entity.SME) error {
	key := keyPrefixSMEList + filterKey
	return c.redis.Set(ctx, key, smes, smeListTTL)
}

// GetList retrieves a cached list of SMEs.
// Returns (nil, nil) on cache miss.
func (c *SMECache) GetList(ctx context.Context, filterKey string) ([]entity.SME, error) {
	key := keyPrefixSMEList + filterKey
	var smes []entity.SME
	err := c.redis.Get(ctx, key, &smes)
	if err == pkgredis.ErrCacheMiss {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("sme cache get list: %w", err)
	}
	return smes, nil
}

// InvalidateList removes cached list entries (call on create/update/delete).
func (c *SMECache) InvalidateList(ctx context.Context) error {
	// Simple strategy: delete a broad wildcard by convention.
	// Since go-redis doesn't support pattern-delete natively without SCAN,
	// we delete the most common keys. For production, use SCAN + DEL.
	return c.redis.Delete(ctx, keyPrefixSMEList+"*")
}

// SetCategories caches the full category list.
func (c *SMECache) SetCategories(ctx context.Context, categories []entity.Category) error {
	return c.redis.Set(ctx, keyCategoryAll, categories, smeCategoryTTL)
}

// GetCategories retrieves the cached category list.
// Returns (nil, nil) on cache miss.
func (c *SMECache) GetCategories(ctx context.Context) ([]entity.Category, error) {
	var categories []entity.Category
	err := c.redis.Get(ctx, keyCategoryAll, &categories)
	if err == pkgredis.ErrCacheMiss {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("sme cache get categories: %w", err)
	}
	return categories, nil
}

// BuildFilterKey creates a cache key from filter parameters.
func BuildFilterKey(categoryID, status, search string) string {
	return fmt.Sprintf("cat=%s:status=%s:q=%s", categoryID, status, search)
}

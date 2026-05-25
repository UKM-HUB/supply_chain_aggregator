package repository

import (
	"context"
	"github.com/redis/go-redis/v9"
)

type GeoRepository struct {
	redisClient *redis.Client
}

func NewGeoRepository(client *redis.Client) *GeoRepository {
	return &GeoRepository{redisClient: client}
}

func (r *GeoRepository) GetNearbyUMKM(ctx context.Context, lng, lat, radiusKm float64) ([]redis.GeoLocation, error) {
	res, err := r.redisClient.GeoRadius(ctx, "umkm_locations", lng, lat, &redis.GeoRadiusQuery{
		Radius:      radiusKm,
		Unit:        "km",
		WithCoord:   true,
		WithDist:    true,
		Sort:        "ASC",
	}).Result()
	return res, err
}

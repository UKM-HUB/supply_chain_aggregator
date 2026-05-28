package usecase

import (
	"context"
	"errors"
	"math"
	"sort"
	"strings"

	"supply-chain-aggregator/services/nearby-service/internal/entity"
	"supply-chain-aggregator/services/nearby-service/internal/repository"
)

var ErrInvalidCoordinate = errors.New("invalid coordinate")

type SearchNearbyInput struct {
	Latitude   float64
	Longitude  float64
	CategoryID string
	RadiusKM   float64
	Limit      int
}

type NearbyUsecase struct {
	locationRepo repository.LocationRepository
}

func NewNearbyUsecase(locationRepo repository.LocationRepository) *NearbyUsecase {
	return &NearbyUsecase{locationRepo: locationRepo}
}

func (u *NearbyUsecase) Search(ctx context.Context, input SearchNearbyInput) ([]entity.NearbySME, error) {
	if input.Latitude < -90 || input.Latitude > 90 || input.Longitude < -180 || input.Longitude > 180 {
		return nil, ErrInvalidCoordinate
	}

	radiusKM := input.RadiusKM
	if radiusKM <= 0 {
		radiusKM = 10
	}

	limit := input.Limit
	if limit <= 0 {
		limit = 10
	}

	locations, err := u.locationRepo.List(ctx, repository.LocationFilter{
		CategoryID: strings.ToLower(strings.TrimSpace(input.CategoryID)),
		Status:     "active",
	})
	if err != nil {
		return nil, err
	}

	result := make([]entity.NearbySME, 0)
	for _, location := range locations {
		distanceKM := haversineKM(input.Latitude, input.Longitude, location.Latitude, location.Longitude)
		if distanceKM > radiusKM {
			continue
		}

		result = append(result, entity.NearbySME{
			SMELocation: location,
			DistanceKM:  distanceKM,
		})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].DistanceKM < result[j].DistanceKM
	})

	if len(result) > limit {
		result = result[:limit]
	}

	return result, nil
}

func haversineKM(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadiusKM = 6371

	lat1Rad := degreesToRadians(lat1)
	lon1Rad := degreesToRadians(lon1)
	lat2Rad := degreesToRadians(lat2)
	lon2Rad := degreesToRadians(lon2)

	dLat := lat2Rad - lat1Rad
	dLon := lon2Rad - lon1Rad

	a := math.Sin(dLat/2)*math.Sin(dLat/2) + math.Cos(lat1Rad)*math.Cos(lat2Rad)*math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadiusKM * c
}

func degreesToRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
}

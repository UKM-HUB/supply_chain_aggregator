package repository

import (
	"context"
	"strings"
	"sync"

	"supply-chain-aggregator/services/nearby-service/internal/entity"
)

type LocationFilter struct {
	CategoryID string
	Status     string
}

type LocationRepository interface {
	List(ctx context.Context, filter LocationFilter) ([]entity.SMELocation, error)
}

type InMemoryLocationRepository struct {
	mu        sync.RWMutex
	locations []entity.SMELocation
}

func NewInMemoryLocationRepository() *InMemoryLocationRepository {
	return &InMemoryLocationRepository{
		locations: []entity.SMELocation{
			{
				ID:          "sme-food-001",
				Name:        "UMKM Maju Food",
				Address:     "Jakarta Selatan",
				Description: "Food supplier for snacks and frozen food",
				CategoryIDs: []string{"food"},
				Latitude:    -6.2245,
				Longitude:   106.8099,
				Status:      "active",
			},
			{
				ID:          "sme-packaging-001",
				Name:        "Kemasan Nusantara",
				Address:     "Jakarta Barat",
				Description: "Packaging supplier for retail and manufacturing",
				CategoryIDs: []string{"packaging", "manufacturing-support"},
				Latitude:    -6.1683,
				Longitude:   106.7588,
				Status:      "active",
			},
			{
				ID:          "sme-textile-001",
				Name:        "Tekstil Mandiri",
				Address:     "Tangerang",
				Description: "Textile and garment production support",
				CategoryIDs: []string{"textile"},
				Latitude:    -6.1783,
				Longitude:   106.6319,
				Status:      "active",
			},
		},
	}
}

func (r *InMemoryLocationRepository) List(ctx context.Context, filter LocationFilter) ([]entity.SMELocation, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	categoryID := strings.ToLower(strings.TrimSpace(filter.CategoryID))
	status := strings.ToLower(strings.TrimSpace(filter.Status))
	result := make([]entity.SMELocation, 0)

	for _, location := range r.locations {
		if status != "" && strings.ToLower(location.Status) != status {
			continue
		}

		if categoryID != "" && !containsCategory(location.CategoryIDs, categoryID) {
			continue
		}

		result = append(result, location)
	}

	return result, nil
}

func containsCategory(values []string, target string) bool {
	for _, value := range values {
		if strings.ToLower(strings.TrimSpace(value)) == target {
			return true
		}
	}

	return false
}

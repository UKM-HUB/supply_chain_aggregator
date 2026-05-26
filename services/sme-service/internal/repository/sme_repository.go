package repository

import (
	"context"
	"sync"

	"supply-chain-aggregator/services/sme-service/internal/entity"
)

type ListFilter struct {
	CategoryID string
	Status     string
	Search     string
}

type SMERepository interface {
	Create(ctx context.Context, sme *entity.SME) error
	List(ctx context.Context, filter ListFilter) ([]entity.SME, error)
	ListCategories(ctx context.Context) ([]entity.Category, error)
}

type InMemorySMERepository struct {
	mu         sync.RWMutex
	smes       []entity.SME
	categories []entity.Category
}

func NewInMemorySMERepository() *InMemorySMERepository {
	return &InMemorySMERepository{
		categories: []entity.Category{
			{ID: "food", Name: "Food Supplier", Description: "Food and beverage suppliers"},
			{ID: "packaging", Name: "Packaging", Description: "Packaging material and services"},
			{ID: "textile", Name: "Textile", Description: "Textile and garment suppliers"},
			{ID: "raw-material", Name: "Raw Material", Description: "Raw material suppliers"},
			{ID: "logistics", Name: "Logistics", Description: "Logistics and delivery support"},
			{ID: "manufacturing-support", Name: "Manufacturing Support", Description: "Supporting vendors for manufacturing operations"},
		},
	}
}

func (r *InMemorySMERepository) Create(ctx context.Context, sme *entity.SME) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	r.smes = append(r.smes, *sme)

	return nil
}

func (r *InMemorySMERepository) List(ctx context.Context, filter ListFilter) ([]entity.SME, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	result := make([]entity.SME, 0)
	for _, sme := range r.smes {
		if filter.Status != "" && sme.Status != filter.Status {
			continue
		}

		if filter.CategoryID != "" && !contains(sme.CategoryIDs, filter.CategoryID) {
			continue
		}

		if filter.Search != "" && !containsText(sme.Name, filter.Search) && !containsText(sme.Description, filter.Search) {
			continue
		}

		result = append(result, sme)
	}

	return result, nil
}

func (r *InMemorySMERepository) ListCategories(ctx context.Context) ([]entity.Category, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	categories := make([]entity.Category, len(r.categories))
	copy(categories, r.categories)

	return categories, nil
}

func contains(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}

	return false
}

func containsText(value, search string) bool {
	value = stringsToLower(value)
	search = stringsToLower(search)

	return len(search) == 0 || stringsContains(value, search)
}

func stringsToLower(value string) string {
	result := []rune(value)
	for index, char := range result {
		if char >= 'A' && char <= 'Z' {
			result[index] = char + 32
		}
	}

	return string(result)
}

func stringsContains(value, search string) bool {
	if len(search) > len(value) {
		return false
	}

	for index := 0; index <= len(value)-len(search); index++ {
		if value[index:index+len(search)] == search {
			return true
		}
	}

	return false
}

package repository

import (
	"context"
	"errors"
	"strings"
	"sync"

	"supply-chain-aggregator/services/sme-service/internal/entity"
)

var ErrCategoryNotFound = errors.New("category not found")

type ListFilter struct {
	CategoryID string
	Status     string
	Search     string
}

type SMERepository interface {
	Create(ctx context.Context, sme *entity.SME) error
	List(ctx context.Context, filter ListFilter) ([]entity.SME, error)
	ListCategories(ctx context.Context) ([]entity.Category, error)
	GetCategoryByID(ctx context.Context, id string) (*entity.Category, error)
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
	categoryID := strings.ToLower(strings.TrimSpace(filter.CategoryID))
	status := strings.ToLower(strings.TrimSpace(filter.Status))
	search := strings.ToLower(strings.TrimSpace(filter.Search))

	for _, sme := range r.smes {
		if status != "" && strings.ToLower(sme.Status) != status {
			continue
		}

		if categoryID != "" && !containsCategory(sme.CategoryIDs, categoryID) {
			continue
		}

		if search != "" && !containsText(sme.Name, search) && !containsText(sme.Description, search) {
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

func (r *InMemorySMERepository) GetCategoryByID(ctx context.Context, id string) (*entity.Category, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	normalizedID := strings.ToLower(strings.TrimSpace(id))
	for _, category := range r.categories {
		if category.ID == normalizedID {
			return &category, nil
		}
	}

	return nil, ErrCategoryNotFound
}

func containsCategory(values []string, target string) bool {
	for _, value := range values {
		if strings.ToLower(strings.TrimSpace(value)) == target {
			return true
		}
	}

	return false
}

func containsText(value, search string) bool {
	return strings.Contains(strings.ToLower(value), search)
}

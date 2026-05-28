package usecase

import (
	"context"
	"errors"
	"strings"
	"time"

	"supply-chain-aggregator/services/sme-service/internal/entity"
	"supply-chain-aggregator/services/sme-service/internal/repository"

	"github.com/google/uuid"
)

var ErrInvalidCategory = errors.New("invalid category")

type CreateSMEInput struct {
	OwnerID     string
	Name        string
	Phone       string
	Address     string
	Description string
	CategoryIDs []string
	Products    []string
	Capacity    string
	Latitude    float64
	Longitude   float64
	Status      string
}

type ListSMEInput struct {
	CategoryID string
	Status     string
	Search     string
}

type SMEUsecase struct {
	smeRepo repository.SMERepository
}

func NewSMEUsecase(smeRepo repository.SMERepository) *SMEUsecase {
	return &SMEUsecase{smeRepo: smeRepo}
}

func (u *SMEUsecase) Create(ctx context.Context, input CreateSMEInput) (*entity.SME, error) {
	status := strings.TrimSpace(input.Status)
	if status == "" {
		status = "active"
	}

	categoryIDs, err := u.normalizeCategoryIDs(ctx, input.CategoryIDs)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	sme := &entity.SME{
		ID:          uuid.NewString(),
		OwnerID:     input.OwnerID,
		Name:        input.Name,
		Phone:       input.Phone,
		Address:     input.Address,
		Description: input.Description,
		CategoryIDs: categoryIDs,
		Products:    input.Products,
		Capacity:    input.Capacity,
		Latitude:    input.Latitude,
		Longitude:   input.Longitude,
		Status:      status,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := u.smeRepo.Create(ctx, sme); err != nil {
		return nil, err
	}

	return sme, nil
}

func (u *SMEUsecase) List(ctx context.Context, input ListSMEInput) ([]entity.SME, error) {
	categoryID := strings.ToLower(strings.TrimSpace(input.CategoryID))
	if categoryID != "" {
		if _, err := u.smeRepo.GetCategoryByID(ctx, categoryID); err != nil {
			return nil, ErrInvalidCategory
		}
	}

	return u.smeRepo.List(ctx, repository.ListFilter{
		CategoryID: categoryID,
		Status:     input.Status,
		Search:     input.Search,
	})
}

func (u *SMEUsecase) ListCategories(ctx context.Context) ([]entity.Category, error) {
	return u.smeRepo.ListCategories(ctx)
}

func (u *SMEUsecase) GetCategoryByID(ctx context.Context, id string) (*entity.Category, error) {
	category, err := u.smeRepo.GetCategoryByID(ctx, id)
	if err != nil {
		return nil, ErrInvalidCategory
	}

	return category, nil
}

func (u *SMEUsecase) normalizeCategoryIDs(ctx context.Context, categoryIDs []string) ([]string, error) {
	normalized := make([]string, 0, len(categoryIDs))
	seen := make(map[string]struct{})

	for _, categoryID := range categoryIDs {
		id := strings.ToLower(strings.TrimSpace(categoryID))
		if id == "" {
			continue
		}

		if _, exists := seen[id]; exists {
			continue
		}

		if _, err := u.smeRepo.GetCategoryByID(ctx, id); err != nil {
			return nil, ErrInvalidCategory
		}

		seen[id] = struct{}{}
		normalized = append(normalized, id)
	}

	return normalized, nil
}

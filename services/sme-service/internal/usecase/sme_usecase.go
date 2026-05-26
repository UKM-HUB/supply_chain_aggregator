package usecase

import (
	"context"
	"strings"
	"time"

	"supply-chain-aggregator/services/sme-service/internal/entity"
	"supply-chain-aggregator/services/sme-service/internal/repository"

	"github.com/google/uuid"
)

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

	now := time.Now()
	sme := &entity.SME{
		ID:          uuid.NewString(),
		OwnerID:     input.OwnerID,
		Name:        input.Name,
		Phone:       input.Phone,
		Address:     input.Address,
		Description: input.Description,
		CategoryIDs: input.CategoryIDs,
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
	return u.smeRepo.List(ctx, repository.ListFilter{
		CategoryID: input.CategoryID,
		Status:     input.Status,
		Search:     input.Search,
	})
}

func (u *SMEUsecase) ListCategories(ctx context.Context) ([]entity.Category, error) {
	return u.smeRepo.ListCategories(ctx)
}

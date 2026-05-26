package http

import (
	"net/http"

	"supply-chain-aggregator/services/sme-service/internal/entity"
	"supply-chain-aggregator/services/sme-service/internal/usecase"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	smeUsecase *usecase.SMEUsecase
}

func NewHandler(smeUsecase *usecase.SMEUsecase) *Handler {
	return &Handler{smeUsecase: smeUsecase}
}

func (h *Handler) Health(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status":  "ok",
		"service": "sme-service",
	})
}

func (h *Handler) CreateSME(c echo.Context) error {
	var request struct {
		OwnerID     string   `json:"owner_id"`
		Name        string   `json:"name"`
		Phone       string   `json:"phone"`
		Address     string   `json:"address"`
		Description string   `json:"description"`
		CategoryIDs []string `json:"category_ids"`
		Products    []string `json:"products"`
		Capacity    string   `json:"capacity"`
		Latitude    float64  `json:"latitude"`
		Longitude   float64  `json:"longitude"`
		Status      string   `json:"status"`
	}

	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "invalid request body"})
	}

	sme, err := h.smeUsecase.Create(c.Request().Context(), usecase.CreateSMEInput{
		OwnerID:     request.OwnerID,
		Name:        request.Name,
		Phone:       request.Phone,
		Address:     request.Address,
		Description: request.Description,
		CategoryIDs: request.CategoryIDs,
		Products:    request.Products,
		Capacity:    request.Capacity,
		Latitude:    request.Latitude,
		Longitude:   request.Longitude,
		Status:      request.Status,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "failed to create SME"})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "SME profile created successfully",
		"data":    formatSME(*sme),
	})
}

func (h *Handler) ListSMEs(c echo.Context) error {
	smes, err := h.smeUsecase.List(c.Request().Context(), usecase.ListSMEInput{
		CategoryID: c.QueryParam("category_id"),
		Status:     c.QueryParam("status"),
		Search:     c.QueryParam("search"),
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "failed to list SMEs"})
	}

	data := make([]map[string]interface{}, 0, len(smes))
	for _, sme := range smes {
		data = append(data, formatSME(sme))
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data":  data,
		"total": len(data),
	})
}

func (h *Handler) ListCategories(c echo.Context) error {
	categories, err := h.smeUsecase.ListCategories(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "failed to list categories"})
	}

	data := make([]map[string]interface{}, 0, len(categories))
	for _, category := range categories {
		data = append(data, map[string]interface{}{
			"id":          category.ID,
			"name":        category.Name,
			"description": category.Description,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data":  data,
		"total": len(data),
	})
}

func formatSME(sme entity.SME) map[string]interface{} {
	return map[string]interface{}{
		"id":           sme.ID,
		"owner_id":     sme.OwnerID,
		"name":         sme.Name,
		"phone":        sme.Phone,
		"address":      sme.Address,
		"description":  sme.Description,
		"category_ids": sme.CategoryIDs,
		"products":     sme.Products,
		"capacity":     sme.Capacity,
		"latitude":     sme.Latitude,
		"longitude":    sme.Longitude,
		"status":       sme.Status,
		"created_at":   sme.CreatedAt,
		"updated_at":   sme.UpdatedAt,
	}
}

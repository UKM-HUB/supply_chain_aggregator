package http

import (
	"errors"
	"net/http"
	"strconv"

	"supply-chain-aggregator/services/nearby-service/internal/entity"
	"supply-chain-aggregator/services/nearby-service/internal/usecase"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	nearbyUsecase *usecase.NearbyUsecase
}

func NewHandler(nearbyUsecase *usecase.NearbyUsecase) *Handler {
	return &Handler{nearbyUsecase: nearbyUsecase}
}

func (h *Handler) Health(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status":  "ok",
		"service": "nearby-service",
	})
}

func (h *Handler) SearchNearbySMEs(c echo.Context) error {
	lat, err := parseRequiredFloat(c.QueryParam("lat"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "lat query parameter is required and must be numeric"})
	}

	lng, err := parseRequiredFloat(c.QueryParam("lng"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "lng query parameter is required and must be numeric"})
	}

	radiusKM := parseOptionalFloat(c.QueryParam("radius_km"), 10)
	limit := parseOptionalInt(c.QueryParam("limit"), 10)

	results, err := h.nearbyUsecase.Search(c.Request().Context(), usecase.SearchNearbyInput{
		Latitude:   lat,
		Longitude:  lng,
		CategoryID: c.QueryParam("category_id"),
		RadiusKM:   radiusKM,
		Limit:      limit,
	})
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidCoordinate) {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "invalid latitude or longitude"})
		}

		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "failed to search nearby SMEs"})
	}

	data := make([]map[string]interface{}, 0, len(results))
	for _, result := range results {
		data = append(data, formatNearbySME(result))
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data":      data,
		"total":     len(data),
		"radius_km": radiusKM,
	})
}

func parseRequiredFloat(value string) (float64, error) {
	if value == "" {
		return 0, strconv.ErrSyntax
	}

	return strconv.ParseFloat(value, 64)
}

func parseOptionalFloat(value string, fallback float64) float64 {
	if value == "" {
		return fallback
	}

	parsed, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return fallback
	}

	return parsed
}

func parseOptionalInt(value string, fallback int) int {
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return parsed
}

func formatNearbySME(result entity.NearbySME) map[string]interface{} {
	return map[string]interface{}{
		"id":           result.ID,
		"name":         result.Name,
		"address":      result.Address,
		"description":  result.Description,
		"category_ids": result.CategoryIDs,
		"latitude":     result.Latitude,
		"longitude":    result.Longitude,
		"status":       result.Status,
		"distance_km":  result.DistanceKM,
	}
}

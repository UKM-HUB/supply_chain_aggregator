package http

import (
	"errors"
	"net/http"

	"supply-chain-aggregator/services/auth-service/internal/repository"
	"supply-chain-aggregator/services/auth-service/internal/usecase"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	authUsecase *usecase.AuthUsecase
}

func NewHandler(authUsecase *usecase.AuthUsecase) *Handler {
	return &Handler{authUsecase: authUsecase}
}

func (h *Handler) Health(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status":  "ok",
		"service": "auth-service",
	})
}

func (h *Handler) Register(c echo.Context) error {
	var request struct {
		Name      string  `json:"name"`
		Email     string  `json:"email"`
		Phone     string  `json:"phone"`
		Password  string  `json:"password"`
		Role      string  `json:"role"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	}

	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "invalid request body"})
	}

	user, err := h.authUsecase.Register(c.Request().Context(), usecase.RegisterInput{
		Name:      request.Name,
		Email:     request.Email,
		Phone:     request.Phone,
		Password:  request.Password,
		Role:      request.Role,
		Latitude:  request.Latitude,
		Longitude: request.Longitude,
	})
	if err != nil {
		if errors.Is(err, repository.ErrEmailAlreadyExists) {
			return c.JSON(http.StatusConflict, map[string]string{"message": "email already exists"})
		}

		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "failed to register user"})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "user registered successfully",
		"user": map[string]interface{}{
			"id":        user.ID,
			"name":      user.Name,
			"email":     user.Email,
			"phone":     user.Phone,
			"role":      user.Role,
			"latitude":  user.Latitude,
			"longitude": user.Longitude,
		},
	})
}

func (h *Handler) Login(c echo.Context) error {
	var request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "invalid request body"})
	}

	output, err := h.authUsecase.Login(c.Request().Context(), request.Email, request.Password)
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidCredential) {
			return c.JSON(http.StatusUnauthorized, map[string]string{"message": "invalid email or password"})
		}

		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "failed to login"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"token": output.Token,
		"user": map[string]interface{}{
			"id":    output.User.ID,
			"name":  output.User.Name,
			"email": output.User.Email,
			"role":  output.User.Role,
		},
	})
}

func (h *Handler) Refresh(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{
		"message": "refresh token flow is not implemented yet",
	})
}

func (h *Handler) Logout(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"message": "logout successful",
	})
}

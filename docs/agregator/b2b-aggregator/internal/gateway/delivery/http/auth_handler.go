package http

import (
	"net/http"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	// Di industri nyata, ini memanggil gRPC Client ke Auth Service
	// Untuk saat ini kita siapkan kerangkanya
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

func (h *AuthHandler) Login(c echo.Context) error {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Input tidak valid"})
	}

	// TODO: Panggil gRPC Client Login di sini

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "success",
		"token":  "mock_jwt_token_from_grpc", 
	})
}

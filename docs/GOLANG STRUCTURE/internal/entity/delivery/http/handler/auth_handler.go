package handler

import (
    "net/http"
    "project/internal/entity"
    "project/internal/usecase"

    "github.com/labstack/echo/v4"
)

type AuthHandler struct {
    Usecase *usecase.AuthUsecase
}

func NewAuthHandler(u *usecase.AuthUsecase) *AuthHandler {
    return &AuthHandler{Usecase: u}
}

func (h *AuthHandler) Register(c echo.Context) error {
    var user entity.User

    if err := c.Bind(&user); err != nil {
        return c.JSON(http.StatusBadRequest, err)
    }

    err := h.Usecase.Register(&user)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, err)
    }

    return c.JSON(http.StatusOK, user)
}
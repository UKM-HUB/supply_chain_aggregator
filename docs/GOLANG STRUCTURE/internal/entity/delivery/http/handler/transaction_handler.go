package handler

import (
    "net/http"
    "project/internal/entity"
    "project/internal/usecase"

    "github.com/labstack/echo/v4"
)

type TransactionHandler struct {
    Usecase *usecase.TransactionUsecase
}

func NewTransactionHandler(u *usecase.TransactionUsecase) *TransactionHandler {
    return &TransactionHandler{Usecase: u}
}

func (h *TransactionHandler) Create(c echo.Context) error {
    var tx entity.Transaction

    if err := c.Bind(&tx); err != nil {
        return c.JSON(http.StatusBadRequest, err)
    }

    err := h.Usecase.Create(&tx)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, err)
    }

    return c.JSON(http.StatusOK, tx)
}
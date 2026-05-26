package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	appName     string
	environment string
	openAPIPath string
}

func NewHandler(appName, environment, openAPIPath string) *Handler {
	return &Handler{
		appName:     appName,
		environment: environment,
		openAPIPath: openAPIPath,
	}
}

func (h *Handler) Health(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status":      "ok",
		"service":     h.appName,
		"environment": h.environment,
	})
}

func (h *Handler) OpenAPI(c echo.Context) error {
	return c.File(h.openAPIPath)
}

func (h *Handler) Login(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{
		"message": "auth login will be forwarded to auth-service through gRPC",
	})
}

func (h *Handler) Register(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{
		"message": "auth registration will be forwarded to auth-service through gRPC",
	})
}

func (h *Handler) ListUsers(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{
		"message": "user listing will be forwarded to user-service",
	})
}

func (h *Handler) ListSMEs(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{
		"message": "SME listing will be forwarded to sme-service",
	})
}

func (h *Handler) ListCategories(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{
		"message": "category listing will be forwarded to sme-service",
	})
}

func (h *Handler) GetCategory(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{
		"message": "category detail will be forwarded to sme-service",
	})
}

func (h *Handler) FindNearbySMEs(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{
		"message": "nearby SME search will be forwarded to nearby-service",
	})
}

func (h *Handler) CreateTransaction(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{
		"message": "transaction creation will be forwarded to transaction-service",
	})
}

func (h *Handler) ListTransactions(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{
		"message": "transaction listing will be forwarded to transaction-service",
	})
}

func (h *Handler) GetTransaction(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{
		"message": "transaction detail will be forwarded to transaction-service",
	})
}

func (h *Handler) CreateVirtualAccount(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{
		"message": "virtual account creation will be forwarded to payment-service",
	})
}

func (h *Handler) DailyReport(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{
		"message": "daily report will be forwarded to report-service",
	})
}

func (h *Handler) MonthlyReport(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{
		"message": "monthly report will be forwarded to report-service",
	})
}

func (h *Handler) ExportReport(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{
		"message": "report export will be forwarded to report-service",
	})
}

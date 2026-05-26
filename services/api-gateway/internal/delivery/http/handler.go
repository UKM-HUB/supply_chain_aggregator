package http

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	appName       string
	environment   string
	openAPIPath   string
	contractsPath string
}

func NewHandler(appName, environment, openAPIPath, contractsPath string) *Handler {
	return &Handler{
		appName:       appName,
		environment:   environment,
		openAPIPath:   openAPIPath,
		contractsPath: contractsPath,
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

// ServeContract serves a specific OpenAPI YAML file by name, e.g. /openapi/auth.yaml
func (h *Handler) ServeContract(c echo.Context) error {
	name := c.Param("file")
	if name == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "file name is required"})
	}
	path := filepath.Join(h.contractsPath, filepath.Base(name))
	return c.File(path)
}

// SwaggerUI serves the Swagger UI HTML page using the CDN build.
// A custom dropdown lets users switch between all 8 service contracts
// without using StandalonePreset (which hijacks the URL to /docs).
func (h *Handler) SwaggerUI(c echo.Context) error {
	html := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8"/>
  <meta name="viewport" content="width=device-width, initial-scale=1"/>
  <title>%s — API Docs</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css"/>
  <style>
    * { box-sizing: border-box; }
    body { margin: 0; padding: 0; font-family: sans-serif; }
    .topbar {
      background: #1b1b1b;
      padding: 10px 20px;
      display: flex;
      align-items: center;
      gap: 16px;
    }
    .topbar span { color: #fff; font-size: 15px; font-weight: 600; }
    .topbar select {
      padding: 6px 12px;
      font-size: 14px;
      border-radius: 4px;
      border: none;
      cursor: pointer;
      min-width: 280px;
    }
  </style>
</head>
<body>
<div class="topbar">
  <span>%s</span>
  <select id="contract-selector" onchange="loadContract(this.value)">
    <option value="/openapi/api-gateway.yaml">API Gateway (port 8080)</option>
    <option value="/openapi/auth.yaml">Auth Service (port 8081)</option>
    <option value="/openapi/sme.yaml">SME Service (port 8082)</option>
    <option value="/openapi/nearby.yaml">Nearby Service (port 8083)</option>
    <option value="/openapi/transactions.yaml">Transaction Service (port 8084)</option>
    <option value="/openapi/payments.yaml">Payment Service (port 8085)</option>
    <option value="/openapi/users.yaml">User Service (port 8088)</option>
    <option value="/openapi/reports.yaml">Report Service (port 8087)</option>
  </select>
</div>
<div id="swagger-ui"></div>
<script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
<script>
  var ui;

  function loadContract(url) {
    if (ui) {
      ui.specActions.updateUrl(url);
      ui.specActions.download(url);
      return;
    }
    ui = SwaggerUIBundle({
      url: url,
      dom_id: '#swagger-ui',
      deepLinking: false,
      presets: [SwaggerUIBundle.presets.apis],
      layout: 'BaseLayout'
    });
  }

  loadContract('/openapi/api-gateway.yaml');
</script>
</body>
</html>`, h.appName, h.appName)

	return c.HTML(http.StatusOK, html)
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

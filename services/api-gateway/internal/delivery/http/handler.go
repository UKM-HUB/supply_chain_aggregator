package http

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/labstack/echo/v4"
)

// Handler menyimpan konfigurasi untuk api-gateway HTTP handler.
type Handler struct {
	appName       string
	environment   string
	openAPIPath   string
	contractsPath string
}

// NewHandler membuat Handler baru.
func NewHandler(appName, environment, openAPIPath, contractsPath string) *Handler {
	return &Handler{
		appName:       appName,
		environment:   environment,
		openAPIPath:   openAPIPath,
		contractsPath: contractsPath,
	}
}

// ContractsPath mengembalikan path ke direktori YAML contracts (dipakai di router).
func (h *Handler) ContractsPath() string {
	return h.contractsPath
}

// Health endpoint.
func (h *Handler) Health(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status":      "ok",
		"service":     h.appName,
		"environment": h.environment,
	})
}

// OpenAPI menyajikan file api-gateway.yaml utama.
func (h *Handler) OpenAPI(c echo.Context) error {
	return c.File(filepath.Join(h.contractsPath, "api-gateway.yaml"))
}

// SwaggerUI menyajikan halaman Swagger UI lengkap dengan dropdown pemilih service.
// File YAML di-serve melalui route /openapi/* (echo.Static di router).
func (h *Handler) SwaggerUI(c echo.Context) error {
	contracts := []struct{ label, file string }{
		{"API Gateway (port 8080)", "api-gateway.yaml"},
		{"Auth Service (port 8081)", "auth.yaml"},
		{"SME Service (port 8082)", "sme.yaml"},
		{"Nearby Service (port 8083)", "nearby.yaml"},
		{"Transaction Service (port 8084)", "transactions.yaml"},
		{"Payment Service (port 8085)", "payments.yaml"},
		{"Report Service (port 8087)", "reports.yaml"},
		{"User Service (port 8088)", "users.yaml"},
	}

	options := ""
	for _, c := range contracts {
		options += fmt.Sprintf(`<option value="/openapi/%s">%s</option>`, c.file, c.label)
	}

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
      padding: 12px 20px;
      display: flex;
      align-items: center;
      gap: 16px;
    }
    .topbar span { color: #fff; font-size: 16px; font-weight: 700; }
    .topbar select {
      padding: 6px 12px;
      font-size: 14px;
      border-radius: 4px;
      border: none;
      cursor: pointer;
      min-width: 300px;
    }
    .topbar a {
      color: #aaa;
      font-size: 13px;
      text-decoration: none;
      margin-left: auto;
    }
    .topbar a:hover { color: #fff; }
  </style>
</head>
<body>
<div class="topbar">
  <span>%s</span>
  <select id="spec-selector" onchange="loadSpec(this.value)">
    %s
  </select>
  <a href="/health">health ✓</a>
</div>
<div id="swagger-ui"></div>
<script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
<script>
  var ui = null;

  function loadSpec(url) {
    if (ui) {
      ui.specActions.updateUrl(url);
      ui.specActions.download(url);
      return;
    }
    ui = SwaggerUIBundle({
      url: url,
      dom_id: '#swagger-ui',
      deepLinking: false,
      presets: [SwaggerUIBundle.presets.apis, SwaggerUIBundle.SwaggerUIStandalonePreset],
      layout: 'BaseLayout',
      validatorUrl: null,
      tryItOutEnabled: true,
      requestInterceptor: function(req) {
        // Tambah Authorization header jika sudah ada token di sessionStorage
        var token = sessionStorage.getItem('jwt_token');
        if (token && !req.headers['Authorization']) {
          req.headers['Authorization'] = 'Bearer ' + token;
        }
        return req;
      }
    });
  }

  // Mulai dengan spec pertama
  var initialSpec = document.getElementById('spec-selector').value;
  loadSpec(initialSpec);
</script>
</body>
</html>`, h.appName, h.appName, options)

	return c.HTML(http.StatusOK, html)
}

// ── Proxy stub handlers ────────────────────────────────────────────────────────
// Semua handler ini akan diteruskan ke masing-masing service melalui gRPC
// setelah integrasi gRPC client selesai. Saat ini mengembalikan 501.

func (h *Handler) Login(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{
		"message": "akan diteruskan ke auth-service melalui gRPC",
	})
}

func (h *Handler) Register(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{
		"message": "akan diteruskan ke auth-service melalui gRPC",
	})
}

func (h *Handler) ListUsers(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{
		"message": "akan diteruskan ke user-service",
	})
}

func (h *Handler) ListSMEs(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{
		"message": "akan diteruskan ke sme-service melalui gRPC",
	})
}

func (h *Handler) ListCategories(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{
		"message": "akan diteruskan ke sme-service",
	})
}

func (h *Handler) GetCategory(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{
		"message": "akan diteruskan ke sme-service",
	})
}

func (h *Handler) FindNearbySMEs(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{
		"message": "akan diteruskan ke nearby-service",
	})
}

func (h *Handler) CreateTransaction(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{
		"message": "akan diteruskan ke transaction-service melalui gRPC",
	})
}

func (h *Handler) ListTransactions(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{
		"message": "akan diteruskan ke transaction-service melalui gRPC",
	})
}

func (h *Handler) GetTransaction(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{
		"message": "akan diteruskan ke transaction-service melalui gRPC",
	})
}

func (h *Handler) CreateVirtualAccount(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{
		"message": "akan diteruskan ke payment-service",
	})
}

func (h *Handler) DailyReport(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{
		"message": "akan diteruskan ke report-service",
	})
}

func (h *Handler) MonthlyReport(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{
		"message": "akan diteruskan ke report-service",
	})
}

func (h *Handler) ExportReport(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{
		"message": "akan diteruskan ke report-service",
	})
}

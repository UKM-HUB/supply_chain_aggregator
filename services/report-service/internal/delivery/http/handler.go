package http

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"supply-chain-aggregator/services/report-service/internal/usecase"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	reportUsecase *usecase.ReportUsecase
}

func NewHandler(reportUsecase *usecase.ReportUsecase) *Handler {
	return &Handler{reportUsecase: reportUsecase}
}

func (h *Handler) Health(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status":  "ok",
		"service": "report-service",
	})
}

func (h *Handler) DailyReport(c echo.Context) error {
	date := time.Now()

	if raw := c.QueryParam("date"); raw != "" {
		parsed, err := time.ParseInLocation("2006-01-02", raw, time.Local)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "invalid date format, expected YYYY-MM-DD"})
		}
		date = parsed
	}

	result, err := h.reportUsecase.Daily(c.Request().Context(), date)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "failed to generate daily report"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": result,
	})
}

func (h *Handler) MonthlyReport(c echo.Context) error {
	now := time.Now()
	year := now.Year()
	month := int(now.Month())

	if raw := c.QueryParam("year"); raw != "" {
		y, err := strconv.Atoi(raw)
		if err != nil || y < 2000 {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "invalid year"})
		}
		year = y
	}

	if raw := c.QueryParam("month"); raw != "" {
		m, err := strconv.Atoi(raw)
		if err != nil || m < 1 || m > 12 {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "invalid month, expected 1-12"})
		}
		month = m
	}

	result, err := h.reportUsecase.Monthly(c.Request().Context(), year, month)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "failed to generate monthly report"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": result,
	})
}

func (h *Handler) ExportReport(c echo.Context) error {
	now := time.Now()
	from := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	to := now

	if raw := c.QueryParam("from"); raw != "" {
		parsed, err := time.ParseInLocation("2006-01-02", raw, time.Local)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "invalid from date format, expected YYYY-MM-DD"})
		}
		from = parsed
	}

	if raw := c.QueryParam("to"); raw != "" {
		parsed, err := time.ParseInLocation("2006-01-02", raw, time.Local)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "invalid to date format, expected YYYY-MM-DD"})
		}
		to = parsed.AddDate(0, 0, 1)
	}

	records, err := h.reportUsecase.Export(c.Request().Context(), from, to)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "failed to export report"})
	}

	var buf bytes.Buffer
	w := csv.NewWriter(&buf)

	w.Write([]string{"id", "invoice_number", "user_id", "amount", "status", "created_at"})
	for _, r := range records {
		w.Write([]string{
			r.ID,
			r.InvoiceNumber,
			r.UserID,
			fmt.Sprintf("%.0f", r.Amount),
			r.Status,
			r.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}
	w.Flush()

	filename := fmt.Sprintf("report_%s_%s.csv",
		from.Format("20060102"),
		to.AddDate(0, 0, -1).Format("20060102"),
	)

	c.Response().Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	return c.Blob(http.StatusOK, "text/csv", buf.Bytes())
}

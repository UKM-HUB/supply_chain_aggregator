package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	"supply-chain-aggregator/services/report-service/internal/entity"
	"supply-chain-aggregator/services/report-service/internal/repository"
)

type ReportUsecase struct {
	repo repository.ReportRepository
}

func NewReportUsecase(repo repository.ReportRepository) *ReportUsecase {
	return &ReportUsecase{repo: repo}
}

type DailyReportResult struct {
	Date   string `json:"date"`
	entity.ReportSummary
}

type MonthlyReportResult struct {
	Month string `json:"month"`
	entity.ReportSummary
}

func (u *ReportUsecase) Daily(ctx context.Context, date time.Time) (DailyReportResult, error) {
	from := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	to := from.AddDate(0, 0, 1)

	records, err := u.repo.ListByDateRange(ctx, repository.DateRange{From: from, To: to})
	if err != nil {
		return DailyReportResult{}, err
	}

	summary := summarise(records)
	return DailyReportResult{
		Date:          date.Format("2006-01-02"),
		ReportSummary: summary,
	}, nil
}

func (u *ReportUsecase) Monthly(ctx context.Context, year, month int) (MonthlyReportResult, error) {
	from := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
	to := from.AddDate(0, 1, 0)

	records, err := u.repo.ListByDateRange(ctx, repository.DateRange{From: from, To: to})
	if err != nil {
		return MonthlyReportResult{}, err
	}

	summary := summarise(records)
	return MonthlyReportResult{
		Month:         fmt.Sprintf("%04d-%02d", year, month),
		ReportSummary: summary,
	}, nil
}

func (u *ReportUsecase) Export(ctx context.Context, from, to time.Time) ([]entity.TransactionRecord, error) {
	return u.repo.ListByDateRange(ctx, repository.DateRange{From: from, To: to})
}

func summarise(records []entity.TransactionRecord) entity.ReportSummary {
	var totalPaid float64
	var totalPending int

	for _, r := range records {
		switch strings.ToLower(r.Status) {
		case "paid":
			totalPaid += r.Amount
		case "pending":
			totalPending++
		}
	}

	return entity.ReportSummary{
		TotalTransaction: len(records),
		TotalPaid:        totalPaid,
		TotalPending:     totalPending,
	}
}

package cron

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"supply-chain-aggregator/services/report-service/internal/usecase"

	"github.com/robfig/cron/v3"
)

// Scheduler wraps the cron library and holds references to usecases.
type Scheduler struct {
	c          *cron.Cron
	reportUC   *usecase.ReportUsecase
	outputDir  string
}

// NewScheduler creates a new Scheduler.
// outputDir is where generated report files will be written (e.g. "/tmp/reports").
func NewScheduler(reportUC *usecase.ReportUsecase, outputDir string) *Scheduler {
	c := cron.New(cron.WithLocation(time.Local))
	return &Scheduler{
		c:         c,
		reportUC:  reportUC,
		outputDir: outputDir,
	}
}

// Register mendaftarkan semua cron job.
func (s *Scheduler) Register() {
	// 1. Daily Report — setiap hari pukul 00:05 (laporan hari sebelumnya)
	s.c.AddFunc("5 0 * * *", s.runDailyReport)

	// 2. Monthly Report — setiap tanggal 1 pukul 01:00 (laporan bulan sebelumnya)
	s.c.AddFunc("0 1 1 * *", s.runMonthlyReport)

	// 3. Cleanup old reports — setiap hari pukul 03:00, hapus file >30 hari
	s.c.AddFunc("0 3 * * *", s.runCleanupOldReports)
}

// Start memulai scheduler di background goroutine.
func (s *Scheduler) Start() {
	s.c.Start()
	fmt.Println("[report-service] cron scheduler started")
	fmt.Println("[report-service] registered jobs:")
	for _, e := range s.c.Entries() {
		fmt.Printf("  - entry #%d | next run: %s\n", e.ID, e.Next.Format(time.RFC3339))
	}
}

// Stop menghentikan scheduler dan menunggu semua job selesai.
func (s *Scheduler) Stop() {
	ctx := s.c.Stop()
	<-ctx.Done()
	fmt.Println("[report-service] cron scheduler stopped")
}

// ─── Job Handlers ─────────────────────────────────────────────────────────────

func (s *Scheduler) runDailyReport() {
	yesterday := time.Now().AddDate(0, 0, -1)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fmt.Printf("[cron:daily-report] generating report for %s\n", yesterday.Format("2006-01-02"))

	result, err := s.reportUC.Daily(ctx, yesterday)
	if err != nil {
		fmt.Printf("[cron:daily-report] ERROR: %v\n", err)
		return
	}

	filename := fmt.Sprintf("%s/daily_%s.json", s.outputDir, yesterday.Format("2006-01-02"))
	if err := writeJSON(filename, result); err != nil {
		fmt.Printf("[cron:daily-report] failed to write file: %v\n", err)
		return
	}

	fmt.Printf("[cron:daily-report] done — total_tx=%d total_paid=%.0f file=%s\n",
		result.TotalTransaction, result.TotalPaid, filename)
}

func (s *Scheduler) runMonthlyReport() {
	now := time.Now()
	// Bulan sebelumnya
	firstOfLastMonth := time.Date(now.Year(), now.Month()-1, 1, 0, 0, 0, 0, now.Location())
	year := firstOfLastMonth.Year()
	month := int(firstOfLastMonth.Month())

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	fmt.Printf("[cron:monthly-report] generating report for %04d-%02d\n", year, month)

	result, err := s.reportUC.Monthly(ctx, year, month)
	if err != nil {
		fmt.Printf("[cron:monthly-report] ERROR: %v\n", err)
		return
	}

	filename := fmt.Sprintf("%s/monthly_%04d-%02d.json", s.outputDir, year, month)
	if err := writeJSON(filename, result); err != nil {
		fmt.Printf("[cron:monthly-report] failed to write file: %v\n", err)
		return
	}

	fmt.Printf("[cron:monthly-report] done — total_tx=%d total_paid=%.0f file=%s\n",
		result.TotalTransaction, result.TotalPaid, filename)
}

func (s *Scheduler) runCleanupOldReports() {
	threshold := time.Now().AddDate(0, 0, -30)
	fmt.Printf("[cron:cleanup] removing reports older than %s from %s\n",
		threshold.Format("2006-01-02"), s.outputDir)

	entries, err := os.ReadDir(s.outputDir)
	if err != nil {
		if os.IsNotExist(err) {
			return // directory belum ada, skip
		}
		fmt.Printf("[cron:cleanup] ERROR reading dir: %v\n", err)
		return
	}

	removed := 0
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			continue
		}
		if info.ModTime().Before(threshold) {
			path := fmt.Sprintf("%s/%s", s.outputDir, entry.Name())
			if err := os.Remove(path); err != nil {
				fmt.Printf("[cron:cleanup] failed to remove %s: %v\n", path, err)
			} else {
				removed++
			}
		}
	}

	fmt.Printf("[cron:cleanup] done — removed %d files\n", removed)
}

// ─── Helper ───────────────────────────────────────────────────────────────────

func writeJSON(path string, v interface{}) error {
	if err := os.MkdirAll(filepath(path), 0755); err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create file %s: %w", path, err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

func filepath(path string) string {
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '/' {
			return path[:i]
		}
	}
	return "."
}

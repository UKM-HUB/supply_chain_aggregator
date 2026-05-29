package cron

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"supply-chain-aggregator/services/report-service/internal/usecase"

	"github.com/robfig/cron/v3"
)

// Scheduler wraps the cron library and holds references to usecases.
type Scheduler struct {
	c         *cron.Cron
	reportUC  *usecase.ReportUsecase
	outputDir string
}

// NewScheduler creates a new Scheduler.
func NewScheduler(reportUC *usecase.ReportUsecase, outputDir string) *Scheduler {
	c := cron.New(cron.WithLocation(time.Local))
	return &Scheduler{c: c, reportUC: reportUC, outputDir: outputDir}
}

// Register mendaftarkan semua cron job.
func (s *Scheduler) Register() {
	// Daily report — setiap hari pukul 00:05
	s.c.AddFunc("5 0 * * *", s.runDailyReport)
	// Monthly report — setiap tanggal 1 pukul 01:00
	s.c.AddFunc("0 1 1 * *", s.runMonthlyReport)
	// Cleanup reports > 30 hari — setiap pukul 03:00
	s.c.AddFunc("0 3 * * *", s.runCleanupOldReports)
}

// Start memulai scheduler di background goroutine.
func (s *Scheduler) Start() {
	s.c.Start()
	fmt.Println("[report-service] cron scheduler started")
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
	if err := writeJSONFile(filename, result); err != nil {
		fmt.Printf("[cron:daily-report] failed to write file: %v\n", err)
		return
	}
	fmt.Printf("[cron:daily-report] done — total_tx=%d total_paid=%.0f file=%s\n",
		result.TotalTransaction, result.TotalPaid, filename)
}

func (s *Scheduler) runMonthlyReport() {
	now := time.Now()
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
	if err := writeJSONFile(filename, result); err != nil {
		fmt.Printf("[cron:monthly-report] failed to write file: %v\n", err)
		return
	}
	fmt.Printf("[cron:monthly-report] done — total_tx=%d total_paid=%.0f file=%s\n",
		result.TotalTransaction, result.TotalPaid, filename)
}

func (s *Scheduler) runCleanupOldReports() {
	threshold := time.Now().AddDate(0, 0, -30)
	fmt.Printf("[cron:cleanup] removing reports older than %s\n", threshold.Format("2006-01-02"))

	entries, err := os.ReadDir(s.outputDir)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		fmt.Printf("[cron:cleanup] ERROR: %v\n", err)
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
			path := s.outputDir + "/" + entry.Name()
			if err := os.Remove(path); err != nil {
				fmt.Printf("[cron:cleanup] failed to remove %s: %v\n", path, err)
			} else {
				removed++
			}
		}
	}
	fmt.Printf("[cron:cleanup] done — removed %d files\n", removed)
}

// writeJSONFile menulis value v sebagai JSON ke path file.
func writeJSONFile(path string, v interface{}) error {
	// Buat direktori jika belum ada
	dir := dirOf(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("mkdir %s: %w", dir, err)
	}

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create %s: %w", path, err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

// dirOf returns the directory portion of a file path.
func dirOf(path string) string {
	i := strings.LastIndex(path, "/")
	if i < 0 {
		return "."
	}
	return path[:i]
}

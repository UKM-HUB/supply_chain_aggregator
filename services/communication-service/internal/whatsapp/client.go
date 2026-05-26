package whatsapp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	apiURL     string
	token      string
	httpClient *http.Client
}

func NewClient(apiURL, token string) *Client {
	return &Client{
		apiURL: apiURL,
		token:  token,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) Send(ctx context.Context, phone, message string) error {
	if c.apiURL == "" {
		fmt.Printf("[whatsapp] (no API configured) to=%s\n%s\n", phone, message)
		return nil
	}

	payload, err := json.Marshal(map[string]string{
		"phone":   phone,
		"message": message,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal whatsapp payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.apiURL, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("failed to create whatsapp request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("whatsapp API request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("whatsapp API returned status %d", resp.StatusCode)
	}

	return nil
}

func FormatPaymentMessage(invoice string, amount float64) string {
	return strings.TrimSpace(fmt.Sprintf(`Pembayaran berhasil diterima.

Invoice: %s
Nominal: Rp%s

Terima kasih.`, invoice, formatRupiah(amount)))
}

func formatRupiah(amount float64) string {
	intAmount := int64(math.Round(amount))
	s := fmt.Sprintf("%d", intAmount)

	result := make([]byte, 0, len(s)+len(s)/3)
	for i, c := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			result = append(result, '.')
		}
		result = append(result, byte(c))
	}

	return string(result)
}

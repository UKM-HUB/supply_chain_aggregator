package xendit

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"supply-chain-aggregator/services/payment-service/internal/entity"
)

const xenditBaseURL = "https://api.xendit.co"

type Client struct {
	secretKey  string
	httpClient *http.Client
}

func NewClient(secretKey string) *Client {
	return &Client{
		secretKey: secretKey,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) CreateInvoice(ctx context.Context, req entity.XenditCreateInvoiceRequest) (entity.XenditCreateInvoiceResponse, error) {
	if c.secretKey == "" {
		return mockInvoiceResponse(req), nil
	}

	body, err := json.Marshal(req)
	if err != nil {
		return entity.XenditCreateInvoiceResponse{}, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, xenditBaseURL+"/v2/invoices", bytes.NewReader(body))
	if err != nil {
		return entity.XenditCreateInvoiceResponse{}, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.SetBasicAuth(c.secretKey, "")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return entity.XenditCreateInvoiceResponse{}, fmt.Errorf("xendit API request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return entity.XenditCreateInvoiceResponse{}, fmt.Errorf("xendit API returned status %d", resp.StatusCode)
	}

	var result entity.XenditCreateInvoiceResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return entity.XenditCreateInvoiceResponse{}, fmt.Errorf("failed to decode xendit response: %w", err)
	}

	return result, nil
}

func mockInvoiceResponse(req entity.XenditCreateInvoiceRequest) entity.XenditCreateInvoiceResponse {
	return entity.XenditCreateInvoiceResponse{
		ID:         fmt.Sprintf("mock-xendit-%d", time.Now().UnixNano()),
		ExternalID: req.ExternalID,
		Amount:     req.Amount,
		InvoiceURL: fmt.Sprintf("https://checkout.xendit.co/web/mock-%s", req.ExternalID),
		Status:     "PENDING",
	}
}

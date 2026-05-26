package entity

type XenditCreateInvoiceRequest struct {
	ExternalID  string  `json:"external_id"`
	Amount      float64 `json:"amount"`
	Currency    string  `json:"currency"`
	Description string  `json:"description,omitempty"`
}

type XenditCreateInvoiceResponse struct {
	ID         string  `json:"id"`
	ExternalID string  `json:"external_id"`
	Amount     float64 `json:"amount"`
	InvoiceURL string  `json:"invoice_url"`
	Status     string  `json:"status"`
}

type XenditWebhookPayload struct {
	ID         string  `json:"id"`
	ExternalID string  `json:"external_id"`
	Status     string  `json:"status"`
	Amount     float64 `json:"amount"`
	PaidAt     string  `json:"paid_at"`
}

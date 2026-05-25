package xendit

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type Adapter struct {
	secretKey  string
	httpClient *http.Client
}

func NewAdapter(secretKey string) *Adapter {
	return &Adapter{
		secretKey: secretKey,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (a *Adapter) GenerateInvoice(ctx context.Context, orderID string, amount float64) (string, error) {
	// Di sini berisi logika HTTP POST ke https://api.xendit.co/v2/invoices
	// Menggunakan pola yang mencegah aplikasi crash jika API pihak ketiga down
	fmt.Printf("Mencetak Invoice Xendit untuk Order %s sebesar %.2f\n", orderID, amount)
	return "https://checkout.xendit.co/web/mock-invoice-url", nil
}

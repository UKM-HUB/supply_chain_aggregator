package entity

type PaymentPaidEvent struct {
	Invoice string  `json:"invoice"`
	Amount  float64 `json:"amount"`
	Phone   string  `json:"phone"`
}

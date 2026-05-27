package services

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type CreateVARequest struct {
	ExternalID string  `json:"external_id"`
	BankCode   string  `json:"bank_code"`
	Name       string  `json:"name"`
	Amount     float64 `json:"amount"`
}

type XenditVAResponse struct {
	ID            string `json:"id"`
	ExternalID    string `json:"external_id"`
	BankCode      string `json:"bank_code"`
	AccountNumber string `json:"account_number"`
	Name          string `json:"name"`
	Status        string `json:"status"`
}

func CreateVirtualAccount(req CreateVARequest) (*XenditVAResponse, error) {

	url := "https://api.xendit.co/callback_virtual_accounts"

	payload, _ := json.Marshal(map[string]interface{}{
		"external_id":  req.ExternalID,
		"bank_code":    req.BankCode,
		"name":         req.Name,
		"expected_amount": req.Amount,
		"is_closed":    true,
	})

	request, _ := http.NewRequest("POST", url, bytes.NewBuffer(payload))

	secretKey := os.Getenv("XENDIT_SECRET_KEY")

	auth := base64.StdEncoding.EncodeToString([]byte(secretKey + ":"))

	request.Header.Set("Authorization", "Basic "+auth)
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	var result XenditVAResponse

	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	if response.StatusCode >= 400 {
		return nil, fmt.Errorf("xendit error")
	}

	return &result, nil
}
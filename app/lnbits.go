package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Invoice struct {
	PaymentHash    string `json:"payment_hash"`
	PaymentRequest string `json:"payment_request"`
}

func createInvoice(memo string) (*Invoice, error) {
	url := fmt.Sprintf("%s/api/v1/payments", os.Getenv("LNBITS_URL"))
	data := map[string]interface{}{
		"amount": os.Getenv("STAMP_PRICE_SATS"),
		"memo":   memo,
		"out":    false,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Api-Key", os.Getenv("LNBITS_API_KEY"))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var invoice Invoice
	if err := json.NewDecoder(resp.Body).Decode(&invoice); err != nil {
		return nil, err
	}

	return &invoice, nil
}

func checkLNbitsPayment(paymentHash string) (bool, error) {
	url := fmt.Sprintf("%s/api/v1/payments/%s", os.Getenv("LNBITS_URL"), paymentHash)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, err
	}

	req.Header.Set("X-Api-Key", os.Getenv("LNBITS_API_KEY"))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	var payment struct {
		Paid bool `json:"paid"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payment); err != nil {
		return false, err
	}

	return payment.Paid, nil
}

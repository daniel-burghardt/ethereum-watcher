package observer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type TransactionWebhookBody struct {
	Addresses []string `json:"addresses"`
	Event     string   `json:"event"`
}

func (s *Service) PostTransactionWebhook(addresses []string) error {
	webhookBody := TransactionWebhookBody{
		Addresses: addresses,
		Event:     "newTransaction",
	}

	body, err := json.Marshal(webhookBody)
	if err != nil {
		return fmt.Errorf("marshalling request body: %w", err)
	}

	_, err = http.Post(s.WebhookUrl, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("http post: %w", err)
	}

	return nil
}

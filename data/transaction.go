package data

import "strings"

type TransactionRepository interface {
	GetTransactions(address string) ([]Transaction, error)
	AddTransaction(transaction Transaction) error
}

type Transaction struct {
	BlockNumber      int64   `json:"blockNumber"`
	TransactionIndex int64   `json:"transactionIndex"`
	From             string  `json:"from"`
	To               string  `json:"to"`
	Value            float64 `json:"value"`
}

func (s *InMemoryStorage) GetTransactions(address string) ([]Transaction, error) {
	address = strings.ToLower(address)
	result := make([]Transaction, 0)
	for _, t := range s.transactions {
		if strings.EqualFold(address, t.From) || strings.EqualFold(address, t.To) {
			result = append(result, t)
		}
	}

	return result, nil
}

func (s *InMemoryStorage) AddTransaction(transaction Transaction) error {
	transaction.From = strings.ToLower(transaction.From)
	transaction.To = strings.ToLower(transaction.To)
	s.transactions = append(s.transactions, transaction)

	return nil
}

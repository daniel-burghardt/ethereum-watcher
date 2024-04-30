package data

import (
	"math"
	"testing"
)

func TestTransactionRepository(t *testing.T) {
	repo := NewInMemoryStorage()
	err := repo.AddTransaction(Transaction{
		BlockNumber:      math.MaxInt64,
		TransactionIndex: math.MaxInt64,
		From:             "0xE59280E580696F9532b9821Ef41E9981A929FEBF",
		To:               "0xFE9E3CDD06Ab5EBB65B408E8503DD67CDE74478E",
		Value:            0.5,
	})
	if err != nil {
		t.Errorf("Error = %v; want nil", err)
	}

	t.Run("get_transactions_case_insensitive", func(t *testing.T) {
		txs, err := repo.GetTransactions("0xe59280e580696f9532b9821Ef41E9981a929FeBf")
		if err != nil {
			t.Errorf("Error = %v; want nil", err)
		}
		if len(txs) != 1 {
			t.Errorf("Len = %d; want 1", len(txs))
		}
		if txs[0].BlockNumber != math.MaxInt64 {
			t.Errorf("BlockNumber = %d; want %d", txs[0].BlockNumber, math.MaxInt64)
		}
		if txs[0].TransactionIndex != math.MaxInt64 {
			t.Errorf("TransactionIndex = %d; want %d", txs[0].TransactionIndex, math.MaxInt64)
		}
		if txs[0].From != "0xe59280e580696f9532b9821ef41e9981a929febf" {
			t.Errorf("From = %s; want 0xe59280e580696f9532b9821ef41e9981a929febf", txs[0].From)
		}
		if txs[0].To != "0xfe9e3cdd06ab5ebb65b408e8503dd67cde74478e" {
			t.Errorf("To = %s; want 0xfe9e3cdd06ab5ebb65b408e8503dd67cde74478e", txs[0].To)
		}
		if txs[0].Value != 0.5 {
			t.Errorf("Value = %f; want %f", txs[0].Value, 0.5)
		}
	})
}

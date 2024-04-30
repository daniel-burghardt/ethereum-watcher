package observer

import (
	"fmt"
	"log"
	"slices"
	"strings"
	"time"

	"github.com/daniel-burghardt/ethereum-parser/data"
	"github.com/daniel-burghardt/ethereum-parser/ethrpc"
	"github.com/daniel-burghardt/ethereum-parser/util"
)

type Service struct {
	RPC        ethrpc.Service
	Repo       data.Repository
	WebhookUrl string
}

func (s *Service) Start() error {
	log.Printf("Starting transactions observer on node %s...", s.RPC.Url)

	latestBlock, err := s.RPC.InvokeBlockNumber()
	if err != nil {
		return fmt.Errorf("invoking blockNumber: %w", err)
	}

	for currentBlock := latestBlock; ; {
		log.Println("Checking for new block...")

		latestBlock, err := s.RPC.InvokeBlockNumber()
		if err != nil {
			return fmt.Errorf("invoking blockNumber: %w", err)
		}

		// If already caught up to latest, wait until next check
		if currentBlock >= latestBlock {
			time.Sleep(time.Second * 1)
			continue
		}

		log.Printf("Processing block %d", currentBlock)

		block, err := s.RPC.InvokeGetBlockByNumber(currentBlock)
		if err != nil {
			return fmt.Errorf("invoking getBlockByHash: %w", err)
		}

		subscribers, err := s.Repo.GetSubscribers()
		if err != nil {
			return fmt.Errorf("getting subscribers: %w", err)
		}

		// Ingest block transactions
		for _, ethTx := range block.Transactions {
			txSubs := FindSubcribers(subscribers, []string{ethTx.From, ethTx.To})
			if len(txSubs) == 0 {
				continue
			}

			log.Printf("Subscribed address found, processing transaction")

			tx, err := ParseTransaction(ethTx)
			if err != nil {
				return fmt.Errorf("parsing transaction: %w", err)
			}

			err = s.Repo.AddTransaction(tx)
			if err != nil {
				return fmt.Errorf("storing transaction: %w", err)
			}

			err = s.PostTransactionWebhook(txSubs)
			if err != nil {
				return fmt.Errorf("posting transaction webhook: %w", err)
			}
		}

		err = s.Repo.UpdateCurrentBlock(currentBlock)
		if err != nil {
			return fmt.Errorf("updating currentBlock: %w", err)
		}

		currentBlock++
	}
}

func FindSubcribers(subscribers []string, targets []string) []string {
	matches := []string{}
	for _, target := range targets {
		if slices.ContainsFunc(subscribers, func(sub string) bool {
			return strings.EqualFold(sub, target)
		}) {
			matches = append(matches, target)
		}
	}

	return matches
}

func ParseTransaction(tx ethrpc.EthTransaction) (data.Transaction, error) {
	blockNumber, err := util.HexToInt(tx.BlockNumber)
	if err != nil {
		return data.Transaction{}, fmt.Errorf("parsing blockNumber: %w", err)
	}

	txIndex, err := util.HexToInt(tx.TransactionIndex)
	if err != nil {
		return data.Transaction{}, fmt.Errorf("parsing transactionIndex: %w", err)
	}

	value, err := util.HexToInt(tx.Value)
	if err != nil {
		return data.Transaction{}, fmt.Errorf("parsing value: %w", err)
	}

	return data.Transaction{
		BlockNumber:      blockNumber,
		TransactionIndex: txIndex,
		From:             tx.From,
		To:               tx.To,
		Value:            util.WeiToEth(value),
	}, nil
}

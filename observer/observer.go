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
	RPC  ethrpc.Service
	Repo data.Repository
}

func (s *Service) Start() error {
	log.Println("Starting transactions observer...")

	filterID, err := s.RPC.InvokeNewBlockFilter()
	if err != nil {
		return fmt.Errorf("invoking newBlockFilter: %w", err)
	}

	for {
		log.Println("Checking for new block...")
		blockHashes, err := s.RPC.InvokeGetFilterChanges(filterID)
		if err != nil {
			return fmt.Errorf("invoking getFilterChanges: %w", err)
		}

		for _, hash := range blockHashes {
			log.Printf("Processing block %s", hash)

			subscribers, err := s.Repo.GetSubscribers()
			if err != nil {
				return fmt.Errorf("getting subscribers: %w", err)
			}

			block, err := s.RPC.InvokeGetBlockByHash(hash)
			if err != nil {
				return fmt.Errorf("invoking getBlockByHash: %w", err)
			}

			for _, ethTx := range block.Transactions {
				// Verify whether either From or To addresses match any of the subscribed addresses
				if !slices.ContainsFunc(subscribers, func(sub string) bool {
					return strings.EqualFold(sub, ethTx.From) || strings.EqualFold(sub, ethTx.To)
				}) {
					continue
				}

				log.Printf("Subscribed address found, storing transaction")

				tx, err := ParseTransaction(ethTx)
				if err != nil {
					return fmt.Errorf("parsing transaction: %w", err)
				}

				s.Repo.AddTransaction(tx)
			}
		}

		if len(blockHashes) == 0 {
			time.Sleep(time.Second * 1)
		}
	}
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

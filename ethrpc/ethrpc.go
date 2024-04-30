package ethrpc

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/daniel-burghardt/ethereum-parser/util"
)

type RequestBody struct {
	Version string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  []any  `json:"params"`
	Id      int    `json:"id"`
}

type ResponseBodyError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ResponseBody struct {
	Version string             `json:"jsonrpc"`
	Id      int                `json:"id"`
	Result  *any               `json:"result"`
	Error   *ResponseBodyError `json:"error"`
}

type Service struct {
	Url string
}

func (s *Service) invokeMethod(reqBody RequestBody) (any, error) {
	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshalling request body: %w", err)
	}

	resp, err := http.Post(s.Url, "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("http post: %w", err)
	}

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}

	var respBody ResponseBody
	err = json.Unmarshal(respBytes, &respBody)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling response body: %w", err)
	}

	if respBody.Error != nil {
		return nil, fmt.Errorf("invalid request: %s", respBody.Error.Message)
	}

	if respBody.Result == nil {
		return nil, errors.New("no result")
	}

	return *respBody.Result, nil
}

func (s *Service) InvokeBlockNumber() (int64, error) {
	result, err := s.invokeMethod(RequestBody{
		Version: "2.0",
		Method:  "eth_blockNumber",
		Params:  []any{},
	})
	if err != nil {
		return 0, fmt.Errorf("invoking method: %w", err)
	}

	blockNumberHex, ok := result.(string)
	if !ok {
		return 0, fmt.Errorf("unexpected response result type: %s", result)
	}

	blockNumber, err := util.HexToInt(blockNumberHex)
	if err != nil {
		return 0, fmt.Errorf("parsing block number: %w", err)
	}

	return blockNumber, nil
}

type EthBlock struct {
	Number       string           `json:"number"`
	Transactions []EthTransaction `json:"transactions"`
}

type EthTransaction struct {
	BlockNumber      string `json:"blockNumber"`
	TransactionIndex string `json:"transactionIndex"`
	From             string `json:"from"`
	To               string `json:"to"`
	Value            string `json:"value"`
}

func (s *Service) InvokeGetBlockByNumber(blockNumber int64) (EthBlock, error) {
	blockNumberHex := util.IntToHex(blockNumber)
	result, err := s.invokeMethod(RequestBody{
		Version: "2.0",
		Method:  "eth_getBlockByNumber",
		Params:  []any{blockNumberHex, true},
	})
	// Format of the response varies depending on the RPC node, but the following should indicate that the requested block is not yet available
	if err != nil && (strings.Contains(err.Error(), "Resource not found.") || strings.Contains(err.Error(), "no result")) {
		return EthBlock{}, nil
	}
	if err != nil {
		return EthBlock{}, fmt.Errorf("invoking method: %w", err)
	}

	block := EthBlock{
		Transactions: []EthTransaction{},
	}

	// Marshal and Unmarshal result object for struct mapping
	marshalled, err := json.Marshal(result)
	if err != nil {
		return EthBlock{}, fmt.Errorf("marshalling result: %s", result)
	}

	err = json.Unmarshal(marshalled, &block)
	if err != nil {
		return EthBlock{}, fmt.Errorf("unmarshalling result: %s", result)
	}

	return block, nil
}

package ethrpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

	return *respBody.Result, nil
}

func (s *Service) InvokeNewBlockFilter() (string, error) {
	result, err := s.invokeMethod(RequestBody{
		Version: "2.0",
		Method:  "eth_newBlockFilter",
		Params:  []any{},
	})
	if err != nil {
		return "", fmt.Errorf("invoking method: %w", err)
	}

	filterID, ok := result.(string)
	if !ok {
		return "", fmt.Errorf("unexpected response result type: %s", result)
	}

	return filterID, nil
}

func (s *Service) InvokeGetFilterChanges(filterID string) ([]string, error) {
	result, err := s.invokeMethod(RequestBody{
		Version: "2.0",
		Method:  "eth_getFilterChanges",
		Params:  []any{filterID},
	})
	if err != nil {
		return []string{}, fmt.Errorf("invoking method: %w", err)
	}

	resultArray, ok := result.([]any)
	if !ok {
		return []string{}, fmt.Errorf("unexpected response result type: %s", result)
	}

	blockHashes := []string{}
	for _, hash := range resultArray {
		hashStr, ok := hash.(string)
		if !ok {
			return []string{}, fmt.Errorf("unexpected response result type: %s", result)
		}
		blockHashes = append(blockHashes, hashStr)
	}

	return blockHashes, nil
}

type EthBlock struct {
	Transactions []EthTransaction `json:"transactions"`
}

type EthTransaction struct {
	BlockNumber      string `json:"blockNumber"`
	TransactionIndex string `json:"transactionIndex"`
	From             string `json:"from"`
	To               string `json:"to"`
	Value            string `json:"value"`
}

func (s *Service) InvokeGetBlockByHash(blockHash string) (EthBlock, error) {
	result, err := s.invokeMethod(RequestBody{
		Version: "2.0",
		Method:  "eth_getBlockByHash",
		Params:  []any{blockHash, true},
	})
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

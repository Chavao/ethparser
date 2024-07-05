package ethereum

import (
	"bytes"
	"encoding/json"
	"ethparser/utils"
	"fmt"
	"io"
)

type Transaction struct {
	TransactionHash  string `json:"transactionHash"`
	TransactionIndex string `json:"transactionIndex"`
}

type TransactionsResponse struct {
	JsonRPC string        `json:"jsonrpc"`
	Id      int           `json:"id"`
	Result  []Transaction `json:"result"`
}

func (rpc *RPCService) Subscribe(address string) bool {
	if !utils.IsValidHex(address) {
		fmt.Printf("Address is not valid")
		return false
	}

	err := rpc.storageClient.Put(address, address)
	if err != nil {
		fmt.Printf("Unable to store address: %s", err)
		return false
	}

	return true
}

func (rpc *RPCService) GetTransactions(address string) []Transaction {
	_, err := rpc.storageClient.Get(address)
	if err != nil {
		fmt.Printf("Unable fetch transactions for address: %s, error %s", address, err)
		return nil
	}

	fields := map[string]any{
		"jsonrpc": "2.0",
		"method":  "eth_getLogs",
		"params": []any{
			map[string]any{
				"address": []string{address},
				"toBlock": "latest",
			},
		},
		"id": 0,
	}

	data, _ := json.Marshal(fields)
	body := bytes.NewBuffer(data)

	r, err := rpc.httpClient.Post(rpc.url, body)
	if err != nil {
		fmt.Printf("Error process request: %v", err)
		return nil
	}
	defer r.Body.Close()

	bodyResponse, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("Error to read response body: %v", err)
		return nil
	}

	response := &TransactionsResponse{}

	err = json.Unmarshal(bodyResponse, response)
	if err != nil {
		fmt.Printf("Error to parse response: %v", err)
		return nil
	}

	return response.Result
}

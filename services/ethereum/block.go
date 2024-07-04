package ethereum

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
)

type BlockNumberResponse struct {
	JsonRPC string `json:"jsonrpc"`
	Id      int    `json:"id"`
	Result  string `json:"result"`
}

func (rpc *RPCService) GetCurrentBlock() int {
	fields := map[string]any{
		"jsonrpc": "2.0",
		"method":  "eth_blockNumber",
		"params":  []any{},
		"id":      0,
	}

	data, _ := json.Marshal(fields)
	body := bytes.NewBuffer(data)

	r, err := rpc.httpClient.Post(rpc.url, body)
	if err != nil {
		fmt.Printf("Error process request: %v", err)
		return 0
	}
	defer r.Body.Close()

	bodyResponse, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("Error to read response body: %v", err)
		return 0
	}

	response := &BlockNumberResponse{}

	err = json.Unmarshal(bodyResponse, response)
	if err != nil {
		fmt.Printf("Error to parse response: %v", err)
		return 0
	}

	return resultToInt(response.Result)
}

func resultToInt(hexValue string) int {
	// Remove "0x" prefix
	cleanHexValue := hexValue[2:]

	value, err := strconv.ParseInt(cleanHexValue, 16, 64)
	if err != nil {
		fmt.Printf("Error to convert from hex: %v", err)
		return 0
	}

	return int(value)
}

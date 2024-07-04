package ethereum

import (
	"ethparser/utils"
	"fmt"
)

type Transaction struct {
	TransactionHash  string `json:"transactionHash"`
	TransactionIndex string `json:"transactionIndex"`
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

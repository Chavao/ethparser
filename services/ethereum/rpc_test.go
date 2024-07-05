package ethereum

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRPCService(t *testing.T) {
	url := "https://rpc-test.com"
	httpClientMock := &mockPostRequester{}
	storageMock := &mockStorage{}

	var expectedHttpClientInterface postRequester
	var expectedStorageClientInterface storage

	rpc := NewRPCService(url, httpClientMock, storageMock)

	assert.NotNil(t, rpc)
	assert.IsType(t, &RPCService{}, rpc)
	assert.Implements(t, &expectedHttpClientInterface, rpc.httpClient)
	assert.Implements(t, &expectedStorageClientInterface, rpc.storageClient)
}

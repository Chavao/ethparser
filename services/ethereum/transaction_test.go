package ethereum

import (
	"bytes"
	"errors"
	"ethparser/utils/tests"
	"io"
	http "net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	validAddress = "0x388c818ca8b9251b393131c08a736a67ccb19297"
	urlAddress   = "https://api.example.com/eth"
)

func TestSubscribe_HappyPath(t *testing.T) {
	storageMock := &mockStorage{}

	storageMock.On("Put", mock.Anything, mock.Anything).Return(nil)

	service := &RPCService{
		url:           urlAddress,
		storageClient: storageMock,
	}

	subscribed := service.Subscribe(validAddress)

	assert.Equal(t, true, subscribed)
}

func TestSubscribe_InvalidAddress(t *testing.T) {
	service := &RPCService{}

	reader, writer, originalStdout := tests.CaptureOutput(t)
	defer func() {
		os.Stdout = originalStdout
	}()

	subscribed := service.Subscribe("invalid-address")

	output := tests.ReadOutput(reader, writer)

	expectedMessage := "Address is not valid"
	assert.Contains(t, output, expectedMessage)

	assert.Equal(t, false, subscribed)
}

func TestSubscribe_ErrorStoreAddress(t *testing.T) {
	storageMock := &mockStorage{}

	storageMock.On("Put", mock.Anything, mock.Anything).Return(errors.New(""))

	service := &RPCService{
		url:           urlAddress,
		storageClient: storageMock,
	}

	reader, writer, originalStdout := tests.CaptureOutput(t)
	defer func() {
		os.Stdout = originalStdout
	}()

	subscribed := service.Subscribe(validAddress)

	output := tests.ReadOutput(reader, writer)

	expectedMessage := "Unable to store address"
	assert.Contains(t, output, expectedMessage)

	assert.Equal(t, false, subscribed)
}

func TestGetTransactions_HappyPath(t *testing.T) {
	httpClientMock := &mockPostRequester{}
	storageMock := &mockStorage{}

	storageMock.On("Get", validAddress).Return(validAddress, nil)

	httpClientMock.On("Post", urlAddress, mock.Anything).Return(&http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewReader([]byte(`{"result": [{"transactionHash": "0xff792001","transactionIndex": "0x123"}]}`))),
	}, nil)

	service := &RPCService{
		url:           urlAddress,
		httpClient:    httpClientMock,
		storageClient: storageMock,
	}

	transactions := service.GetTransactions(validAddress)

	assert.NotNil(t, transactions)
}

func TestGetTransactions_ErrorGetSubscribedAddress(t *testing.T) {
	storageMock := &mockStorage{}

	storageMock.On("Get", validAddress).Return("", errors.New(""))

	service := &RPCService{
		storageClient: storageMock,
	}

	reader, writer, originalStdout := tests.CaptureOutput(t)
	defer func() {
		os.Stdout = originalStdout
	}()

	transactions := service.GetTransactions(validAddress)

	output := tests.ReadOutput(reader, writer)

	expectedMessage := "Unable fetch transactions for address"
	assert.Contains(t, output, expectedMessage)

	assert.Nil(t, transactions)
}

func TestGetTransactions_ErrorPost(t *testing.T) {
	httpClientMock := &mockPostRequester{}
	storageMock := &mockStorage{}

	storageMock.On("Get", validAddress).Return(validAddress, nil)

	httpClientMock.On("Post", urlAddress, mock.Anything).Return(nil, errors.New(""))

	service := &RPCService{
		url:           urlAddress,
		httpClient:    httpClientMock,
		storageClient: storageMock,
	}

	reader, writer, originalStdout := tests.CaptureOutput(t)
	defer func() {
		os.Stdout = originalStdout
	}()

	transactions := service.GetTransactions(validAddress)

	output := tests.ReadOutput(reader, writer)

	expectedMessage := "Error process request"
	assert.Contains(t, output, expectedMessage)
	assert.Nil(t, transactions)
}

func TestGetTransactions_ErrorReadBody(t *testing.T) {
	httpClientMock := &mockPostRequester{}
	storageMock := &mockStorage{}

	storageMock.On("Get", validAddress).Return(validAddress, nil)

	httpClientMock.On("Post", urlAddress, mock.Anything).Return(&http.Response{
		StatusCode: 504,
		Body:       errorReaderMock{},
	}, nil)

	service := &RPCService{
		url:           urlAddress,
		httpClient:    httpClientMock,
		storageClient: storageMock,
	}

	reader, writer, originalStdout := tests.CaptureOutput(t)
	defer func() {
		os.Stdout = originalStdout
	}()

	transactions := service.GetTransactions(validAddress)

	output := tests.ReadOutput(reader, writer)

	expectedMessage := "Error to read response body"
	assert.Contains(t, output, expectedMessage)
	assert.Nil(t, transactions)
}

func TestGetTransactions_ErrorParseResponse(t *testing.T) {
	httpClientMock := &mockPostRequester{}
	storageMock := &mockStorage{}

	storageMock.On("Get", validAddress).Return(validAddress, nil)

	httpClientMock.On("Post", urlAddress, mock.Anything).Return(&http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewReader([]byte("invalid-json"))),
	}, nil)

	service := &RPCService{
		url:           urlAddress,
		httpClient:    httpClientMock,
		storageClient: storageMock,
	}

	reader, writer, originalStdout := tests.CaptureOutput(t)
	defer func() {
		os.Stdout = originalStdout
	}()

	transactions := service.GetTransactions(validAddress)

	output := tests.ReadOutput(reader, writer)

	expectedMessage := "Error to parse response"
	assert.Contains(t, output, expectedMessage)
	assert.Nil(t, transactions)
}

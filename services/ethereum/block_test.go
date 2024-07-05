package ethereum

import (
	"bytes"
	"errors"
	"ethparser/utils/tests"
	"fmt"
	"io"
	http "net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type errorReaderMock struct{}

func (errorReaderMock) Read(p []byte) (n int, err error) {
	return 0, fmt.Errorf("forced read error")
}

func (errorReaderMock) Close() error {
	return nil
}

func TestGetCurrentBlock_HappyPath(t *testing.T) {
	httpClientMock := &mockPostRequester{}

	httpClientMock.On("Post", urlAddress, mock.Anything).Return(&http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewReader([]byte(`{"result": "0x134c63d"}`))),
	}, nil)

	service := &RPCService{
		url:        urlAddress,
		httpClient: httpClientMock,
	}

	block := service.GetCurrentBlock()

	assert.Equal(t, 20235837, block)
}

func TestGetCurrentBlock_ErrorPost(t *testing.T) {
	httpClientMock := &mockPostRequester{}

	httpClientMock.On("Post", urlAddress, mock.Anything).Return(nil, errors.New(""))

	service := &RPCService{
		url:        urlAddress,
		httpClient: httpClientMock,
	}

	reader, writer, originalStdout := tests.CaptureOutput(t)
	defer func() {
		os.Stdout = originalStdout
	}()

	block := service.GetCurrentBlock()

	output := tests.ReadOutput(reader, writer)

	expectedMessage := "Error process request"
	assert.Contains(t, output, expectedMessage)
	assert.Equal(t, 0, block)
}

func TestGetCurrentBlock_ErrorReadBody(t *testing.T) {
	httpClientMock := &mockPostRequester{}

	httpClientMock.On("Post", urlAddress, mock.Anything).Return(&http.Response{
		StatusCode: 504,
		Body:       errorReaderMock{},
	}, nil)

	service := &RPCService{
		url:        urlAddress,
		httpClient: httpClientMock,
	}

	reader, writer, originalStdout := tests.CaptureOutput(t)
	defer func() {
		os.Stdout = originalStdout
	}()

	block := service.GetCurrentBlock()

	output := tests.ReadOutput(reader, writer)

	expectedMessage := "Error to read response body"
	assert.Contains(t, output, expectedMessage)
	assert.Equal(t, 0, block)
}

func TestGetCurrentBlock_ErrorParseResponse(t *testing.T) {
	httpClientMock := &mockPostRequester{}

	httpClientMock.On("Post", urlAddress, mock.Anything).Return(&http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewReader([]byte("invalid-json"))),
	}, nil)

	service := &RPCService{
		url:        urlAddress,
		httpClient: httpClientMock,
	}

	reader, writer, originalStdout := tests.CaptureOutput(t)
	defer func() {
		os.Stdout = originalStdout
	}()

	block := service.GetCurrentBlock()

	output := tests.ReadOutput(reader, writer)

	expectedMessage := "Error to parse response"
	assert.Contains(t, output, expectedMessage)
	assert.Equal(t, 0, block)
}

func TestResultToInt(t *testing.T) {
	result := resultToInt("0x134c63d")

	assert.Equal(t, 20235837, result)
}

func TestResultToInt_ErrorToConvertFromHex(t *testing.T) {
	reader, writer, originalStdout := tests.CaptureOutput(t)
	defer func() {
		os.Stdout = originalStdout
	}()

	result := resultToInt("potato")

	output := tests.ReadOutput(reader, writer)

	expectedMessage := "Error to convert from hex"
	assert.Contains(t, output, expectedMessage)
	assert.Equal(t, 0, result)
}

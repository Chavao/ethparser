package ethereum

import (
	"errors"
	"ethparser/utils/tests"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSubscribe_HappyPath(t *testing.T) {
	storageMock := &mockStorage{}

	storageMock.On("Put", mock.Anything, mock.Anything).Return(nil)

	service := &RPCService{
		storageClient: storageMock,
	}

	subscribed := service.Subscribe("0x388c818ca8b9251b393131c08a736a67ccb19297")

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
		storageClient: storageMock,
	}

	reader, writer, originalStdout := tests.CaptureOutput(t)
	defer func() {
		os.Stdout = originalStdout
	}()

	subscribed := service.Subscribe("0x388c818ca8b9251b393131c08a736a67ccb19297")

	output := tests.ReadOutput(reader, writer)

	expectedMessage := "Unable to store address"
	assert.Contains(t, output, expectedMessage)

	assert.Equal(t, false, subscribed)
}

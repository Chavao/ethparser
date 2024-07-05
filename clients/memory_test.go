package clients

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMemoryStorage(t *testing.T) {
	storageClient := NewMemoryStorage()
	assert.NotNil(t, storageClient)
}

func TestStorageClient_Put(t *testing.T) {
	storageClient := NewMemoryStorage()

	err := storageClient.Put("foo", "bar")
	require.NoError(t, err)

	lock.Lock()
	assert.Equal(t, "bar", storage["foo"])
	lock.Unlock()
}

func TestStorageClient_Get_HappyPath(t *testing.T) {
	storageClient := NewMemoryStorage()

	err := storageClient.Put("foo", "bar")
	require.NoError(t, err)

	value, err := storageClient.Get("foo")

	require.NoError(t, err)
	assert.Equal(t, "bar", value)
}

func TestStorageClient_Get_ErrorRetriveValue(t *testing.T) {
	storageClient := NewMemoryStorage()

	value, err := storageClient.Get("nonexistent")

	assert.Empty(t, value)
	assert.Error(t, err)
	assert.Equal(t, errors.New("value not found"), err)
}

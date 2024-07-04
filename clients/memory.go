package clients

import (
	"errors"
	"sync"
)

type StorageClient struct {
}

func NewMemoryStorage() *StorageClient {
	return &StorageClient{}
}

var lock = &sync.Mutex{}
var storage = make(map[string]string)

func (s *StorageClient) Put(key string, value string) (err error) {
	lock.Lock()
	defer lock.Unlock()

	storage[key] = value

	return nil
}

func (s *StorageClient) Get(key string) (value string, err error) {
	if value, ok := storage[key]; ok {
		return value, nil
	}

	return "", errors.New("value not found")
}

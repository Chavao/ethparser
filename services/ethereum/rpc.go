package ethereum

import (
	"io"
	"net/http"
)

type postRequester interface {
	Post(url string, body io.Reader) (resp *http.Response, err error)
}

type storage interface {
	Put(key string, value string) (err error)
	Get(key string) (value string, err error)
}

type RPCService struct {
	url           string
	httpClient    postRequester
	storageClient storage
}

func NewRPCService(url string, httpClient postRequester, storageClient storage) *RPCService {
	return &RPCService{
		url:           url,
		httpClient:    httpClient,
		storageClient: storageClient,
	}
}

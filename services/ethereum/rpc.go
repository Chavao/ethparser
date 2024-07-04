package ethereum

import (
	"io"
	"net/http"
)

type postRequester interface {
	Post(url string, body io.Reader) (resp *http.Response, err error)
}

type RPCService struct {
	url        string
	httpClient postRequester
}

func NewRPCService(url string, httpClient postRequester) *RPCService {
	return &RPCService{
		url:        url,
		httpClient: httpClient,
	}
}

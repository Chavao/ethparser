package clients

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewHttpClient(t *testing.T) {
	httpClient := NewHttpClient()
	assert.NotNil(t, httpClient)
}

func TestHttpClient_Post_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/test", r.URL.Path)

		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		assert.Equal(t, "body-mock", string(body))

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("response-mock"))
	}))
	defer server.Close()

	httpClient := NewHttpClient()
	r, err := httpClient.Post(server.URL+"/test", strings.NewReader("body-mock"))

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, r.StatusCode)

	responseBody, err := io.ReadAll(r.Body)
	require.NoError(t, err)
	assert.Equal(t, "response-mock", string(responseBody))
}

func TestHttpClient_Post_ErrorCreatingRequest(t *testing.T) {
	httpClient := NewHttpClient()

	r, err := httpClient.Post("http://invalid-url", nil)

	assert.Error(t, err)
	assert.Nil(t, r)
}

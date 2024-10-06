package httptesting

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/georgepsarakis/go-httpclient"
	"github.com/georgepsarakis/go-httpclient/httpassert"
)

func TestClient_Get(t *testing.T) {
	c := NewClient(t)
	requestURL := "http://localhost/p123"
	c.Client.WithDefaultHeaders(map[string]string{"Content-Type": "application/json"})
	c.NewMockRequest(http.MethodGet, requestURL+"?test=1").
		RespondWithJSON(http.StatusOK, `{"name": "hello", "surname": "world"}`).Register()
	resp, err := c.Get(context.Background(),
		requestURL,
		httpclient.WithQueryParameters(map[string]string{"test": "1"}))
	require.NoError(t, err)
	reqHeaders := http.Header{}
	reqHeaders.Set("Content-Type", "application/json")
	httpassert.ResponseEqual(t, resp, &http.Response{
		StatusCode: http.StatusOK,
		Request: &http.Request{
			Method: http.MethodGet,
			Header: reqHeaders,
			URL:    httpassert.URLFromString(t, requestURL+"?test=1"),
			Proto:  "HTTP/1.1",
		},
		Header: http.Header{"Content-Type": []string{"application/json"}},
		Body:   io.NopCloser(strings.NewReader(`{"name": "hello", "surname": "world"}`)),
	})
	httpassert.SuccessfulJSONResponseEqual(t, resp, []byte(`{"name": "hello", "surname": "world"}`))
}

func TestClient_Head(t *testing.T) {
	c := NewClient(t)
	requestURL := "http://localhost/p123"
	c.Client.WithDefaultHeaders(map[string]string{"Content-Type": "application/json"})
	c.NewMockRequest(http.MethodHead, requestURL+"?test=1",
		httpclient.WithHeaders(map[string]string{"Content-Type": "application/json"})).
		RespondWithJSON(http.StatusOK, `{"name": "hello", "surname": "world"}`).
		RespondWithHeaders(map[string]string{"Content-Type": "application/json"}).
		Register()
	resp, err := c.Head(context.Background(),
		requestURL,
		httpclient.WithQueryParameters(map[string]string{"test": "1"}))
	require.NoError(t, err)
	reqHeaders := http.Header{}
	reqHeaders.Set("Content-Type", "application/json")
	httpassert.ResponseEqual(t, resp, &http.Response{
		StatusCode: http.StatusOK,
		Request: &http.Request{
			Method: http.MethodHead,
			Header: reqHeaders,
			URL:    httpassert.URLFromString(t, requestURL+"?test=1"),
			Proto:  "HTTP/1.1",
		},
		Header: http.Header{"Content-Type": []string{"application/json"}},
	})
	httpassert.ResponseEqual(t, resp, &http.Response{
		StatusCode: http.StatusOK,
	})
}

func TestClient_WithBaseURL(t *testing.T) {
	c := NewClient(t)
	c, err := c.WithBaseURL("http://www.example.com/test")
	require.NoError(t, err)
	require.Equal(t, "http://www.example.com/test", c.BaseURL())
}

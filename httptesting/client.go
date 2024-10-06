package httptesting

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/georgepsarakis/go-httpclient"
)

type HttpTestRequestParameter httpclient.RequestParameter

type Client struct {
	*httpclient.Client
	mock *httpmock.MockTransport
	t    *testing.T
}

func NewClient(t *testing.T) *Client {
	mt := httpmock.NewMockTransport()
	return &Client{
		Client: httpclient.NewWithTransport(mt),
		mock:   mt,
		t:      t,
	}
}

func (c *Client) Get(ctx context.Context, url string, parameters ...httpclient.RequestParameter) (*http.Response, error) {
	return c.Client.Get(ctx, url, parameters...)
}

func (c *Client) Head(ctx context.Context, url string, parameters ...httpclient.RequestParameter) (*http.Response, error) {
	return c.Client.Head(ctx, url, parameters...)
}

func (c *Client) Post(ctx context.Context, url string, body io.Reader, parameters ...httpclient.RequestParameter) (*http.Response, error) {
	return c.Client.Post(ctx, url, body, parameters...)
}

func (c *Client) Patch(ctx context.Context, url string, body io.Reader, parameters ...httpclient.RequestParameter) (*http.Response, error) {
	return c.Client.Patch(ctx, url, body, parameters...)
}

func (c *Client) Delete(ctx context.Context, url string, parameters ...httpclient.RequestParameter) (*http.Response, error) {
	return c.Client.Delete(ctx, url, parameters...)
}

// WithBaseURL sets the base URL setting for the underlying `httpclient.Client`.
func (c *Client) WithBaseURL(baseURL string) (*Client, error) {
	_, err := c.Client.WithBaseURL(baseURL)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Client) WithDefaultHeaders(headers map[string]string) *Client {
	c.Client = c.Client.WithDefaultHeaders(headers)
	return c
}

// HTTPMock exposes the httpmock.MockTransport instance for advanced usage.
func (c *Client) HTTPMock() *httpmock.MockTransport {
	return c.mock
}

type MockRequest struct {
	req            *http.Request
	requestMatcher httpmock.Matcher
	responder      httpmock.Responder
	t              *testing.T
	mockTransport  *httpmock.MockTransport
}

type MockResponse httpmock.Responder

func (c *Client) NewMockRequest(method, url string, params ...httpclient.RequestParameter) *MockRequest {
	c.t.Helper()

	req, err := http.NewRequest(method, url, nil)
	require.NoError(c.t, err)

	opts := httpclient.NewRequestParameters()
	if len(params) > 0 {
		opts = httpclient.NewRequestParameters(params...)
	}

	matcherName := fmt.Sprintf("%s_%s", c.t.Name(), url)
	mReq := &MockRequest{
		t:             c.t,
		req:           req,
		mockTransport: c.mock,
		requestMatcher: httpmock.NewMatcher(matcherName, func(r *http.Request) bool {
			return r.Method == method &&
				r.URL.String() == url &&
				(opts.Headers() == nil || assert.ObjectsAreEqual(opts.Headers(), r.Header))
		}),
		responder: httpmock.NewStringResponder(http.StatusOK, "OK"),
	}
	return mReq
}

func (r *MockRequest) Register() {
	r.mockTransport.RegisterMatcherResponder(
		r.req.Method,
		r.req.URL.String(),
		r.requestMatcher,
		r.responder)
}

// String provides a representation of the mock request. Only used for debugging purposes.
func (r *MockRequest) String() string {
	return fmt.Sprintf("MockRequest: [%s] %s", r.req.Method, r.req.URL.String())
}

// Responder provides access to the current responder for inspection or direct operations.
func (r *MockRequest) Responder(resp httpmock.Responder) *MockRequest {
	r.responder = resp
	return r
}

// RespondWithJSON will configure a JSON response with the given status code.
func (r *MockRequest) RespondWithJSON(statusCode int, body string) *MockRequest {
	responder, err := httpmock.NewJsonResponder(statusCode, json.RawMessage(body))
	require.NoError(r.t, err)
	r.responder = responder
	return r
}

// RespondWithHeaders configures the response headers. It can be used multiple times in order to pass different headers.
// If the header key already exists it will be overwritten.
func (r *MockRequest) RespondWithHeaders(respHeaders map[string]string) *MockRequest {
	h := http.Header{}
	for k, v := range respHeaders {
		h.Set(k, v)
	}
	r.responder.HeaderSet(h)
	return r
}

func (c *Client) NewJSONBodyMatcher(body string) httpmock.MatcherFunc {
	c.t.Helper()

	return func(r *http.Request) bool {
		var m1, m2 map[string]any
		require.NoError(c.t, json.Unmarshal(interceptBody(c.t, r), &m1))
		require.NoError(c.t, json.Unmarshal([]byte(body), &m2))
		return assert.ObjectsAreEqual(m1, m2)
	}
}

func interceptBody(t *testing.T, req *http.Request) []byte {
	t.Helper()
	body, err := io.ReadAll(req.Body)
	require.NoError(t, err)
	req.Body.Close()
	req.Body = io.NopCloser(bytes.NewReader(body))
	return body
}

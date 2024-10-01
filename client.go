package httpclient

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	base           http.RoundTripper
	timeout        time.Duration
	defaultHeaders map[string]string
	baseURL        *url.URL
	networkClient  *http.Client
}

const DefaultTimeout = 30 * time.Second

func New() *Client {
	return NewWithTransport(http.DefaultTransport)
}

func NewWithTransport(transport http.RoundTripper) *Client {
	if transport == nil {
		panic("transport must be non-nil")
	}
	return &Client{
		timeout: DefaultTimeout,
		base:    transport,
		networkClient: &http.Client{
			Timeout:   DefaultTimeout,
			Transport: transport,
		},
	}
}

func (c *Client) WithTimeout(timeout time.Duration) *Client {
	c.timeout = timeout
	return c
}

func (c *Client) WithBaseTransport(base http.RoundTripper) *Client {
	c.base = base
	return c
}

func (c *Client) WithDefaultHeaders(headers map[string]string) *Client {
	if c.defaultHeaders == nil {
		c.defaultHeaders = make(map[string]string)
	}
	for k, v := range headers {
		c.defaultHeaders[k] = v
	}
	return c
}

func (c *Client) WithBaseURL(baseURL string) (*Client, error) {
	base, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	c.baseURL = base
	return c, nil
}

func (c *Client) WithJSONContentType() *Client {
	return c.WithDefaultHeaders(map[string]string{"Content-Type": "application/json"})
}

func (c *Client) Get(ctx context.Context, url string, parameters ...RequestParameter) (*http.Response, error) {
	req, err := c.prepareRequest(ctx, http.MethodGet, url, nil, parameters...)
	if err != nil {
		return nil, err
	}
	return c.networkClient.Do(req)
}

func (c *Client) prepareRequest(ctx context.Context, method string, rawURL string, body io.Reader, parameters ...RequestParameter) (*http.Request, error) {
	var reqParams []RequestParameter
	reqParams = append(reqParams, WithHeaders(c.defaultHeaders))
	reqParams = append(reqParams, parameters...)
	params := NewRequestParameters(reqParams...)
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}
	var fullURL *url.URL
	if c.baseURL != nil {
		fullURL = c.baseURL.ResolveReference(parsedURL)
	} else {
		fullURL = parsedURL
	}
	req, err := http.NewRequestWithContext(ctx, method, fullURL.String(), body)
	if err != nil {
		return nil, err
	}
	req.Header = params.headers
	req.URL.RawQuery += params.queryParams.Encode()
	return req, nil
}

func (c *Client) Head(ctx context.Context, url string, parameters ...RequestParameter) (*http.Response, error) {
	req, err := c.prepareRequest(ctx, http.MethodHead, url, nil, parameters...)
	if err != nil {
		return nil, err
	}
	return c.networkClient.Do(req)
}

func (c *Client) Post(ctx context.Context, url string, body io.Reader, parameters ...RequestParameter) (*http.Response, error) {
	req, err := c.prepareRequest(ctx, http.MethodPost, url, body, parameters...)
	if err != nil {
		return nil, err
	}
	return c.networkClient.Do(req)
}

func (c *Client) Patch(ctx context.Context, url string, body io.Reader, parameters ...RequestParameter) (*http.Response, error) {
	req, err := c.prepareRequest(ctx, http.MethodPatch, url, body, parameters...)
	if err != nil {
		return nil, err
	}
	return c.networkClient.Do(req)
}

func (c *Client) Delete(ctx context.Context, url string, parameters ...RequestParameter) (*http.Response, error) {
	req, err := c.prepareRequest(ctx, http.MethodDelete, url, nil, parameters...)
	if err != nil {
		return nil, err
	}
	return c.networkClient.Do(req)
}

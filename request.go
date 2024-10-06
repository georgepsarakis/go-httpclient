package httpclient

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
)

type RequestParameter func(opts *RequestParameters)

type RequestParameters struct {
	headers     http.Header
	queryParams url.Values
	// Convert response with the following status code to error return values
	errorCodes []int
}

// QueryParameters returns a clone of the currently configured query parameters.
// Multiple calls will override already existing keys.
func (rp *RequestParameters) QueryParameters() url.Values {
	if rp.queryParams == nil {
		return nil
	}
	qp := make(url.Values, len(rp.queryParams))
	for k, v := range rp.queryParams {
		for _, qv := range v {
			qp.Add(k, qv)
		}
	}
	return qp
}

func (rp *RequestParameters) Headers() http.Header {
	return rp.headers.Clone()
}

func (rp *RequestParameters) ErrorCodes() []int {
	return rp.errorCodes
}

// WithQueryParameters configures the given name-value pairs as Query String parameters for the request.
// Multiple calls will override values for existing keys.
func WithQueryParameters(params map[string]string) RequestParameter {
	return func(opts *RequestParameters) {
		if opts.queryParams == nil {
			opts.queryParams = url.Values{}
		}
		for name, value := range params {
			opts.queryParams.Set(name, value)
		}
	}
}

// WithHeaders allows headers to be set on the request. Multiple calls using the same header name
// will overwrite existing header values.
func WithHeaders(headers map[string]string) RequestParameter {
	return func(opts *RequestParameters) {
		if opts.headers == nil {
			opts.headers = http.Header{}
		}
		for name, value := range headers {
			opts.headers.Set(name, value)
		}
	}
}

func NewRequestParameters(opts ...RequestParameter) *RequestParameters {
	rp := &RequestParameters{}
	for _, o := range opts {
		o(rp)
	}
	return rp
}

func InterceptRequestBody(r *http.Request) ([]byte, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	if err := r.Body.Close(); err != nil {
		return nil, err
	}
	r.Body = io.NopCloser(bytes.NewReader(body))
	return body, nil
}

func MustInterceptRequestBody(r *http.Request) []byte {
	b, err := InterceptRequestBody(r)
	if err != nil {
		panic(err)
	}
	return b
}

// NewRequest builds a new request based on the given Method, full URL, body and optional functional option parameters.
func NewRequest(ctx context.Context, method string, rawURL string, body io.Reader, parameters ...RequestParameter) (*http.Request, error) {
	reqParams := NewRequestParameters(parameters...)
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}
	if encodedQP := reqParams.queryParams.Encode(); encodedQP != "" {
		parsedURL.RawQuery += encodedQP
	}
	req, err := http.NewRequestWithContext(ctx, method, parsedURL.String(), body)
	if err != nil {
		return nil, err
	}
	req.Header = reqParams.headers
	return req, nil
}

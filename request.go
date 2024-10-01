package httpclient

import (
	"bytes"
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

func (rp *RequestParameters) QueryParameters() url.Values {
	if rp.queryParams == nil {
		return nil
	}
	qp := make(url.Values, len(rp.queryParams))
	for k, v := range rp.queryParams {
		qp[k] = v
	}
	return qp
}

func (rp *RequestParameters) Headers() http.Header {
	return rp.headers.Clone()
}

func (rp *RequestParameters) ErrorCodes() []int {
	return rp.errorCodes
}

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
	r.Body.Close()
	r.Body = io.NopCloser(bytes.NewReader(body))
	return body, nil
}

func MustInterceptRequestBody(r *http.Request, body []byte) []byte {
	b, err := InterceptRequestBody(r)
	if err != nil {
		panic(err)
	}
	return b
}

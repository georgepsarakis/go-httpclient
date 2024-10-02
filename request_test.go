package httpclient

import (
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWithQueryParameters(t *testing.T) {
	tests := []struct {
		name    string
		params  map[string]string
		params2 map[string]string
		want    url.Values
	}{
		{
			name: "set query parameters with a non-nil map",
			params: map[string]string{
				"q1": "2",
			},
			want: url.Values{
				"q1": []string{"2"},
			},
		},
		{
			name:   "query parameters are nil",
			params: nil,
			want:   url.Values{},
		},
		{
			name: "query parameters are already set",
			params: map[string]string{
				"q1": "2",
				"q2": "3",
			},
			params2: map[string]string{
				"q1": "4",
				"q4": "5",
			},
			want: url.Values{
				"q1": []string{"4"},
				"q2": []string{"3"},
				"q4": []string{"5"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qp := []RequestParameter{WithQueryParameters(tt.params)}
			if tt.params2 != nil {
				qp = append(qp, WithQueryParameters(tt.params2))
			}
			opts := NewRequestParameters(qp...)

			assert.Equal(t, tt.want, opts.QueryParameters())
		})
	}
}

func TestMustInterceptRequestBody(t *testing.T) {
	require.Panics(t, func() {
		MustInterceptRequestBody(&http.Request{Body: failureOnReadReader{}})
	})
	require.Panics(t, func() {
		MustInterceptRequestBody(&http.Request{Body: failureOnCloseReader{}})
	})

	req := &http.Request{Body: io.NopCloser(strings.NewReader("test"))}
	require.Equal(t, []byte("test"), MustInterceptRequestBody(req))
	b, err := io.ReadAll(req.Body)
	require.NoError(t, err)
	require.Equal(t, []byte("test"), b)
}

type failureOnReadReader struct {
	io.ReadCloser
}

func (f failureOnReadReader) Read(_ []byte) (n int, err error) {
	return 0, io.ErrUnexpectedEOF
}

func (f failureOnReadReader) Close() error {
	return nil
}

type failureOnCloseReader struct {
	io.ReadCloser
}

func (f failureOnCloseReader) Read(_ []byte) (n int, err error) {
	return 0, io.EOF
}

func (f failureOnCloseReader) Close() error {
	return io.ErrClosedPipe
}

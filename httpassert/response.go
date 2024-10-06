package httpassert

import (
	"bytes"
	"io"
	"mime"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func ResponseEqual(t *testing.T, actual, expected *http.Response) {
	t.Helper()

	assert.Equal(t, expected.StatusCode, actual.StatusCode)
	if expected.Header != nil {
		assert.Equal(t, expected.Header, actual.Header)
	}
	if expected.Body != nil {
		expectedBody, err := io.ReadAll(expected.Body)
		require.NoError(t, err)
		expected.Body.Close()
		actualBody, err := io.ReadAll(actual.Body)
		require.NoError(t, err)
		actual.Body.Close()
		// Restore the body stream in order to allow multiple assertions
		actual.Body = io.NopCloser(bytes.NewBuffer(actualBody))
		mediatype, _, err := mime.ParseMediaType(actual.Header.Get("Content-Type"))
		if mediatype == "application/json" {
			assert.JSONEq(t, string(expectedBody), string(actualBody))
		} else {
			assert.Equal(t, string(expectedBody), string(actualBody))
		}
	}
	if expected.Request != nil {
		assert.Equal(t, expected.Request.URL, actual.Request.URL)
		assert.Equal(t, expected.Request.Method, actual.Request.Method)
		assert.Equal(t, expected.Request.Proto, actual.Request.Proto)
		assert.Equal(t, expected.Request.Header, actual.Request.Header)
	}
}

func SuccessfulJSONResponseEqual(t *testing.T, actual *http.Response, body []byte) {
	t.Helper()
	ResponseEqual(t, actual, &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader(body)),
	})
}

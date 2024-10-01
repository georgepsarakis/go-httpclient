package httpassert

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func URLFromString(t *testing.T, fullURL string) *url.URL {
	t.Helper()
	u, err := url.Parse(fullURL)
	require.NoError(t, err)
	return u
}

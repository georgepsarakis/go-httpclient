package integration

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/georgepsarakis/go-httpclient"
	"github.com/georgepsarakis/go-httpclient/httpassert"
)

func TestClient_Get_JSON(t *testing.T) {
	c := githubClient(t)
	resp, err := c.Get(context.Background(), "/repos/georgepsarakis/go-httpclient")
	require.NoError(t, err)
	v := map[string]interface{}{}
	require.NoError(t, httpclient.DeserializeJSON(resp, &v))
	httpassert.PrintJSON(t, v)
}

func githubClient(t *testing.T) *httpclient.Client {
	t.Helper()
	c, err := httpclient.New().
		WithJSONContentType().
		WithBaseURL("https://api.github.com")
	require.NoError(t, err)
	return c
}

func TestClient_Get_JSON_ContextDeadline(t *testing.T) {
	c := githubClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 0)
	defer cancel()
	resp, err := c.Get(ctx, "/url/doesnt/matter")

	require.ErrorIs(t, err, context.DeadlineExceeded)
	require.Nil(t, resp)
}

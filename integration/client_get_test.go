package integration

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/georgepsarakis/go-httpclient"
)

func TestClient_Get_JSON(t *testing.T) {
	c := githubClient(t)
	resp, err := c.Get(context.Background(), "/repos/georgepsarakis/go-httpclient")
	require.NoError(t, err)
	v := map[string]interface{}{}
	require.NoError(t, httpclient.DeserializeJSON(resp, &v))
	printJSON(t, v)
}

func githubClient(t *testing.T) *httpclient.Client {
	t.Helper()
	c, err := httpclient.New().
		WithJSONContentType().
		WithBaseURL("https://api.github.com")
	require.NoError(t, err)
	return c
}

func printJSON(t *testing.T, v any) {
	b, err := json.MarshalIndent(v, "", "  ")
	require.NoError(t, err)
	t.Log(string(b))
}

func TestClient_Get_JSON_ContextDeadline(t *testing.T) {
	c := githubClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 0)
	defer cancel()
	resp, err := c.Get(ctx, "/url/doesnt/matter")

	require.ErrorIs(t, err, context.DeadlineExceeded)
	require.Nil(t, resp)
}

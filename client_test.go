package httpclient

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

func TestClient_Get_JSON(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	c, err := New().
		WithJSONContentType().
		WithBaseURL("https://api.github.com")
	require.NoError(t, err)

	responseBody := `{
	  "allow_forking": true,
	  "archive_url": "https://api.github.com/repos/georgepsarakis/go-httpclient/{archive_format}{/ref}",
	  "archived": false,
	  "assignees_url": "https://api.github.com/repos/georgepsarakis/go-httpclient/assignees{/user}",
	  "blobs_url": "https://api.github.com/repos/georgepsarakis/go-httpclient/git/blobs{/sha}",
	  "branches_url": "https://api.github.com/repos/georgepsarakis/go-httpclient/branches{/branch}",
	  "clone_url": "https://github.com/georgepsarakis/go-httpclient.git",
	  "collaborators_url": "https://api.github.com/repos/georgepsarakis/go-httpclient/collaborators{/collaborator}",
	  "comments_url": "https://api.github.com/repos/georgepsarakis/go-httpclient/comments{/number}",
	  "commits_url": "https://api.github.com/repos/georgepsarakis/go-httpclient/commits{/sha}",
	  "compare_url": "https://api.github.com/repos/georgepsarakis/go-httpclient/compare/{base}...{head}",
	  "contents_url": "https://api.github.com/repos/georgepsarakis/go-httpclient/contents/{+path}",
	  "contributors_url": "https://api.github.com/repos/georgepsarakis/go-httpclient/contributors",
	  "created_at": "2024-06-07T09:59:37Z",
	  "default_branch": "main"
	}`
	responder, err := httpmock.NewJsonResponder(http.StatusOK, json.RawMessage(responseBody))
	require.NoError(t, err)
	httpmock.RegisterMatcherResponder(
		http.MethodGet,
		"https://api.github.com/repos/georgepsarakis/go-httpclient",
		httpmock.HeaderIs("Content-Type", "application/json; charset=utf-8"),
		responder,
	)
	url := "/repos/georgepsarakis/go-httpclient"
	resp, err := c.Get(context.Background(), url)
	require.NoError(t, err)

	v := map[string]any{}
	require.NoError(t, DeserializeJSON(resp, &v))
	require.JSONEq(t, responseBody, string(MustInterceptResponseBody(resp)))
}

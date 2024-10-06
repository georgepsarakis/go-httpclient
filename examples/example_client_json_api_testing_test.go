package examples_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/georgepsarakis/go-httpclient"
	"github.com/georgepsarakis/go-httpclient/examples/githubsdk"
	"github.com/georgepsarakis/go-httpclient/httpassert"
	"github.com/georgepsarakis/go-httpclient/httptesting"
)

func TestGitHubSDK_GetUserByUsername(t *testing.T) {
	testClient := httptesting.NewClient(t)
	sdk := githubsdk.NewWithClient(testClient.Client)
	testClient, err := testClient.WithBaseURL("https://test-api-github-com")
	require.NoError(t, err)

	testClient.
		NewMockRequest(
			http.MethodGet,
			"https://test-api-github-com/users/georgepsarakis",
			httpclient.WithHeaders(map[string]string{
				"x-github-api-version": "2022-11-28",
				"accept":               "application/vnd.github+json",
			})).
		RespondWithJSON(http.StatusOK, `
	{
	   "id": 963304,
	   "bio": "Test 123",
	   "blog": "https://test.blog/",
	   "created_at": "2025-09-16T16:57:12Z",
       "login": "georgepsarakis",
	   "name": "Test Name"
	}`).Register()

	user, err := sdk.GetUserByUsername(context.Background(), "georgepsarakis")
	require.NoError(t, err)

	httpassert.PrintJSON(t, user)
	// Output:
	//{
	//   "id": 963304,
	//   "bio": "Test 123",
	//   "blog": "https://test.blog/",
	//   "created_at": "2025-09-16T16:57:12Z",
	//   "login": "georgepsarakis",
	//   "name": "Test Name"
	//}
	require.Equal(t, githubsdk.User{
		ID:        963304,
		Bio:       "Test 123",
		Blog:      "https://test.blog/",
		CreatedAt: time.Date(2025, time.September, 16, 16, 57, 12, 0, time.UTC),
		Login:     "georgepsarakis",
		Name:      "Test Name",
	}, user)
}

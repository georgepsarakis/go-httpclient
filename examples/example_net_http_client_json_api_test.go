package examples_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"

	"github.com/georgepsarakis/go-httpclient/httpassert"
)

type GitHubSDK struct {
	Client         *http.Client
	DefaultHeaders http.Header
	BaseURL        string
}

func New() GitHubSDK {
	client := http.DefaultClient
	return NewWithClient(client)
}

func NewWithClient(c *http.Client) GitHubSDK {
	headers := http.Header{}
	headers.Set("X-GitHub-Api-Version", "2022-11-28")
	headers.Set("Accept", "application/vnd.github+json")
	return GitHubSDK{Client: c, DefaultHeaders: headers, BaseURL: "https://api.github.com"}
}

type User struct {
	ID        int       `json:"id"`
	Bio       string    `json:"bio"`
	Blog      string    `json:"blog"`
	CreatedAt time.Time `json:"created_at"`
	Login     string    `json:"login"`
	Name      string    `json:"name"`
}

// GetUserByUsername retrieves a user based on their public username.
// See https://docs.github.com/en/rest/users/users
func (g GitHubSDK) GetUserByUsername(ctx context.Context, username string) (User, error) {
	path, err := url.JoinPath("/users", username)
	if err != nil {
		return User{}, err
	}
	fullURL := fmt.Sprintf("%s%s", g.BaseURL, path)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL, nil)
	if err != nil {
		return User{}, err
	}
	req.Header = g.DefaultHeaders.Clone()
	resp, err := g.Client.Do(req)
	if err != nil {
		return User{}, err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return User{}, err
	}
	u := User{}
	if err := json.Unmarshal(b, &u); err != nil {
		return User{}, err
	}
	return u, nil
}

func TestGitHubSDK_NetHTTP_GetUserByUsername(t *testing.T) {
	var err error
	mt := httpmock.NewMockTransport()
	testClient := &http.Client{
		Transport: mt,
	}
	sdk := NewWithClient(testClient)
	sdk.BaseURL = "https://test-api-github-com"

	responder, err := httpmock.NewJsonResponder(http.StatusOK, json.RawMessage(`
	{
	   "id": 963304,
	   "bio": "Test 123",
	   "blog": "https://test.blog/",
	   "created_at": "2025-09-16T16:57:12Z",
       "login": "georgepsarakis",
	   "name": "Test Name"
	}`))
	require.NoError(t, err)
	defaultHeaderMatcher := func(req *http.Request) bool {
		return req.Header.Get("Accept") == "application/vnd.github+json" &&
			req.Header.Get("X-GitHub-Api-Version") == "2022-11-28"
	}
	mt.RegisterMatcherResponder(http.MethodGet,
		"https://test-api-github-com/users/georgepsarakis",
		httpmock.NewMatcher("GetUserByUsername", func(req *http.Request) bool {
			if !defaultHeaderMatcher(req) {
				return false
			}
			return true
		}),
		responder)

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
	require.Equal(t, User{
		ID:        963304,
		Bio:       "Test 123",
		Blog:      "https://test.blog/",
		CreatedAt: time.Date(2025, time.September, 16, 16, 57, 12, 0, time.UTC),
		Login:     "georgepsarakis",
		Name:      "Test Name",
	}, user)
}

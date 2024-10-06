# go-httpclient

## Summary

**go-httpclient** aims to reduce the boilerplate of HTTP request/response
setup for Go, along with providing out-of-the-box testing in a standardized way.

The Go standard library [net/http](https://pkg.go.dev/net/http) has already an excellent & powerful API.
However, the complexity of reliably & securely composing an HTTP request or reading back the HTTP response
cannot be easily avoided without a higher-level abstraction layer.

**go-httpclient** also tries to enforce best practices such as:
- Non-zero request timeout
- Passing `context.Context` to `net/http`
- URL-encoded Query Parameters
- Safe URL Path Joining
- Always closing the response body

Furthermore, testing is facilitated by the `httptesting` & `httpassert` libraries. 
`httptesting` provides a 100% compatible API with `httpclient.Client` and exposes a `httpclient.Client` instance that 
can be injected directly as a drop-in replacement of the regular production code `Client`.
The testing abstraction layer is using a [httpmock](https://github.com/jarcoal/httpmock) Mock Transport under the hood
which allows registration of custom matcher/responders when required.

## Key Features

- Offers an intuitive and ergonomic API based on HTTP verb names `Get, Post, Patch, Delete` and functional option parameters.
- All request emitter methods accept `context.Context` as their first parameter. 
- Uses plain `map[string]string` structures for passing Query Parameters and Headers which should cover the majority of cases.
- Always URL-encodes query parameters.
- Ensures Response body is read when streaming is not required.
- Separate testing `Client` that implements the exact same API.
- Utilizes the powerful [httpmock](https://github.com/jarcoal/httpmock) under the hood in order to allow fine-grained and 
  scoped request mocking and assertions.

## Examples

Example implementation of a GitHub REST API SDK using `httpclient.Client`:

```go
package githubsdk

import (
	"context"
	"net/url"
	"time"

	"github.com/georgepsarakis/go-httpclient"
)

type GitHubSDK struct {
	Client *httpclient.Client
}

func New() GitHubSDK {
	client := httpclient.New()
	return NewWithClient(client)
}

func NewWithClient(c *httpclient.Client) GitHubSDK {
	c.WithDefaultHeaders(map[string]string{
		"X-GitHub-Api-Version": "2022-11-28",
		"Accept":               "application/vnd.github+json",
	})
	c, _ = c.WithBaseURL("https://api.github.com")
	return GitHubSDK{Client: c}
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
	// Note: `httpclient.Client.Get` allows header parameterization, for example changing an API version:
	// resp, err := g.Client.Get(ctx, path, httpclient.WithHeaders(map[string]string{"x-github-api-version": "2023-11-22"}))
	resp, err := g.Client.Get(ctx, path)
	if err != nil {
		return User{}, err
	}
	u := User{}
	if err := httpclient.DeserializeJSON(resp, &u); err != nil {
		return u, err
	}
	return u, nil
}
```

Here is how using our SDK looks like:

```go
package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/georgepsarakis/go-httpclient/examples/githubsdk"
)

func main() {
	sdk := githubsdk.New()

	user, err := sdk.GetUserByUsername(context.Background(), "georgepsarakis")
	panicOnError(err)

	m, err := json.MarshalIndent(user, "", "  ")
	panicOnError(err)

	fmt.Println(string(m))
	// Output:
	//{
	//   "id": 963304,
	//   "bio": "Software Engineer",
	//   "blog": "https://controlflow.substack.com/",
	//   "created_at": "2011-08-06T16:57:12Z",
	//   "login": "georgepsarakis",
	//   "name": "George Psarakis"
	//}
}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}
```

Testing our GitHub SDK as well as code that depends on it is straightforward thanks to the `httptesting` package:

```go
package githubsdk_test

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
	var err error

	testClient := httptesting.NewClient(t)
	sdk := githubsdk.NewWithClient(testClient.Client)
	testClient, err = testClient.WithBaseURL("https://test-api-github-com")
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
	require.Equal(t, githubsdk.User{
		ID:        963304,
		Bio:       "Test 123",
		Blog:      "https://test.blog/",
		CreatedAt: time.Date(2025, time.September, 16, 16, 57, 12, 0, time.UTC),
		Login:     "georgepsarakis",
		Name:      "Test Name",
	}, user)
}
```

For comparison, below is an alternative implementation using the `net/http` & `httpmock` packages:

```go
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
	require.Equal(t, User{
		ID:        963304,
		Bio:       "Test 123",
		Blog:      "https://test.blog/",
		CreatedAt: time.Date(2025, time.September, 16, 16, 57, 12, 0, time.UTC),
		Login:     "georgepsarakis",
		Name:      "Test Name",
	}, user)
}
```

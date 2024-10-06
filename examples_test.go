package httpclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"time"
)

func ExampleClient_Get() {
	sdk := NewSDK()
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

type GitHubSDK struct {
	Client *Client
}

func NewSDK() GitHubSDK {
	return NewSDKWithClient(New())
}

func NewSDKWithClient(c *Client) GitHubSDK {
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
	resp, err := g.Client.Get(ctx, path)
	if err != nil {
		return User{}, err
	}
	u := User{}
	if err := DeserializeJSON(resp, &u); err != nil {
		return u, err
	}
	return u, nil
}

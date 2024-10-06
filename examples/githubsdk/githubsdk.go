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

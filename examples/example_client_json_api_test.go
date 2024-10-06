package examples

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/georgepsarakis/go-httpclient/examples/githubsdk"
)

func Example() {
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

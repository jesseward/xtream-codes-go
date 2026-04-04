## Xtream Codes GO


Fork of [pbergman/xtream-codes-go](https://github.com/pbergman/xtream-codes-go)

golang [xtream codes api](https://github.com/engenex/xtream-codes-api-v2/blob/main/%5BHow-To%5D%20Player%20API%20v2%20-%20Tutorials%20-%20Xtream%20Codes.pdf) client for fetching data from server and creting stream urls.

```go

package main

import (
	"context"
	"fmt"

	xtream_codes "github.com/jesseward/xtream-codes-go"
)

func main() {

	config, err := xtream_codes.NewApiClientConfig("http://example.com", "username", "password")

	if err != nil {
		panic(err)
	}

	client, err := xtream_codes.NewApiClient(config, nil, nil, nil)

	if err != nil {
		panic(err)
	}

	categories, err := client.GetLiveCategories(context.Background())

	if err != nil {
		panic(err)
	}

	for _, category := range categories {
		fmt.Printf("%-4d %s", category.GetId(), category.GetName())
	}
}
```


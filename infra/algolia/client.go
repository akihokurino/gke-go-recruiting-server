package algolia

import (
	"fmt"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
)

func NewClient(appID string, apiKey string, envName string) *Client {
	return &Client{
		envName: envName,
		core:    search.NewClient(appID, apiKey),
	}
}

type Client struct {
	envName string
	core    *search.Client
}

func (c *Client) Index(name string) *search.Index {
	return c.core.InitIndex(fmt.Sprintf("%s_%s", c.envName, name))
}

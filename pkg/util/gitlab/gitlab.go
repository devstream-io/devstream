package gitlab

import (
	"fmt"
	"os"

	"github.com/xanzy/go-gitlab"

	"github.com/merico-dev/stream/pkg/util/log"
)

var client *Client

type Client struct {
	*gitlab.Client
}

func NewClient() (*Client, error) {
	if client != nil {
		log.Debug("Use a cached client.")
		return client, nil
	}

	token := os.Getenv("GITLAB_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("failed to read GITLAB_TOKEN from environment variable")
	}
	client, err := gitlab.NewClient(token)
	if err != nil {
		return nil, err
	}

	return &Client{Client: client}, nil
}

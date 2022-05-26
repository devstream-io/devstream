package gitlab

import (
	"errors"
	"os"

	"github.com/xanzy/go-gitlab"

	"github.com/devstream-io/devstream/pkg/util/log"
)

var client *Client

const (
	DefaultGitlabHost = "https://gitlab.com"
)

type OptionFunc func(*Client)

func WithBaseURL(baseURL string) OptionFunc {
	return func(c *Client) {
		c.baseURL = baseURL
	}
}

type Client struct {
	baseURL string
	*gitlab.Client
}

func NewClient(opts ...OptionFunc) (*Client, error) {
	if client != nil {
		log.Debug("Use a cached client.")
		return client, nil
	}

	token := os.Getenv("GITLAB_TOKEN")
	if token == "" {
		return nil, errors.New("failed to read GITLAB_TOKEN from environment variable")
	}

	c := &Client{}

	for _, opt := range opts {
		opt(c)
	}

	var err error

	if c.baseURL == "" {
		c.Client, err = gitlab.NewClient(token)
	} else {
		c.Client, err = gitlab.NewClient(token, gitlab.WithBaseURL(c.baseURL))

	}

	if err != nil {
		return nil, err
	}

	return c, nil

}

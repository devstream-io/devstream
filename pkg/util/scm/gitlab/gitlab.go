package gitlab

import (
	"errors"
	"os"

	"github.com/xanzy/go-gitlab"

	"github.com/devstream-io/devstream/pkg/util/scm/git"
)

const (
	DefaultGitlabHost = "https://gitlab.com"
	TokenEnvKey       = "GITLAB_TOKEN"
)

type Client struct {
	*gitlab.Client
	*git.RepoInfo
}

func NewClient(options *git.RepoInfo) (*Client, error) {
	token := os.Getenv(TokenEnvKey)
	if token == "" {
		return nil, errors.New("failed to read GITLAB_TOKEN from environment variable")
	}

	c := &Client{}

	var err error

	if options.BaseURL == "" {
		c.Client, err = gitlab.NewClient(token)
	} else {
		c.Client, err = gitlab.NewClient(token, gitlab.WithBaseURL(options.BaseURL))
	}
	c.RepoInfo = options

	if err != nil {
		return nil, err
	}

	return c, nil

}

package gitlab

import (
	"errors"

	"github.com/xanzy/go-gitlab"

	"github.com/devstream-io/devstream/pkg/util/scm/git"
)

const (
	DefaultGitlabHost = "https://gitlab.com"
)

type Client struct {
	*gitlab.Client
	*git.RepoInfo
}

func NewClient(options *git.RepoInfo) (*Client, error) {
	if options.Token == "" {
		return nil, errors.New("config field scm.token is not setted")
	}

	c := &Client{}

	var err error

	if options.BaseURL == "" {
		c.Client, err = gitlab.NewClient(options.Token)
	} else {
		c.Client, err = gitlab.NewClient(options.Token, gitlab.WithBaseURL(options.BaseURL))
	}
	c.RepoInfo = options

	if err != nil {
		return nil, err
	}

	return c, nil

}

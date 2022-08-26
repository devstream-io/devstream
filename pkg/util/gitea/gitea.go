package gitea

import (
	"errors"
	"os"
	"code.gitea.io/sdk/gitea"

	"github.com/devstream-io/devstream/pkg/util/git"
)

const (
	DefaultGiteaHost = "https://gitea.com/"
)

type Client struct {
	*gitea.Client
	*git.RepoInfo
}
func NewClient(options *git.RepoInfo) (*Client, error) {
	token := os.Getenv("GITLEA_TOKEN")
	if token == "" {
		return nil, errors.New("failed to read GITEA_TOKEN from environment variable")
	}
	c := &Client{}
	var err error	
	if options.BaseURL == "" {
		c.Client, err = gitea.NewClient(DefaultGiteaHost,gitea.SetToken(token))
	} else {
		c.Client, err = gitea.NewClient(options.BaseURL,gitea.SetToken(token))
	}
	c.RepoInfo = options
	if err != nil {
		return nil, err
	}
	return c, nil

}
package github

import (
	"context"
	"net/http"

	"github.com/google/go-github/v42/github"

	"github.com/merico-dev/stream/pkg/util/log"
)

func (c *Client) CreateRepo() error {
	repo := &github.Repository{
		Name: &c.Repo,
	}

	_, _, err := c.Repositories.Create(context.TODO(), "", repo)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) Delete() error {
	_, err := c.Client.Repositories.Delete(
		c.Context,
		c.Owner,
		c.Repo)

	if err != nil {
		return err
	}
	log.Successf("GitHub repo %s removed.", c.Repo)
	return nil
}

func (c *Client) GetRepoDescription() (*github.Repository, error) {
	repo, resp, err := c.Client.Repositories.Get(
		c.Context,
		c.Owner,
		c.Repo)

	if repo == nil && resp.StatusCode == http.StatusNotFound {
		return repo, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, err
	}

	return repo, nil
}

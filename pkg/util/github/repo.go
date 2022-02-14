package github

import (
	"context"
	"github.com/merico-dev/stream/internal/pkg/log"

	"github.com/google/go-github/v42/github"
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

func (c *Client) DeleteRepo() error {
	_, err := c.Client.Repositories.Delete(
		c.Context,
		c.Option.Owner,
		c.Option.Repo)

	if err != nil {
		return err
	}
	log.Successf("GitHub repo %s removed.", c.Repo)
	return nil
}

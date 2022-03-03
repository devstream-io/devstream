package github

import (
	"context"
	"fmt"

	"github.com/go-errors/errors"
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
	rps, resp, err := c.Client.Repositories.Get(
		c.Context,
		c.Owner,
		c.Repo)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New("response status is not 200 OK, but is " + fmt.Sprintf("%d", resp.StatusCode))
	}

	return rps, nil
}

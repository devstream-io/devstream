package github

import (
	"context"
	"fmt"

	"github.com/go-errors/errors"
	"github.com/google/go-github/v42/github"

	"github.com/merico-dev/stream/internal/pkg/log"
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

func (c *Client) IsRepoExists() error {
	rps, resp, err := c.Client.Repositories.Get(
		c.Context,
		c.Owner,
		c.Repo)

	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return errors.New("response status is not 200 OK, but is " + fmt.Sprintf("%d", resp.StatusCode))
	}

	log.Successf("GitHub repo exists, repo is %s, owner is %s.", *rps.Name, *rps.Owner.Login)
	return nil
}

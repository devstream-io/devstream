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

func (c *Client) CheckRepo() error {
	rps, rsp, err := c.Client.Repositories.Get(
		c.Context,
		c.Option.Owner,
		c.Option.Repo)

	if err != nil {
		return err
	}

	if rsp.StatusCode != 200 {
		return errors.New("response status is not 200 OK, but is " + fmt.Sprintf("%d", rsp.StatusCode))
	}

	log.Successf("GitHub repo exists, repo is %s, owner is %s.", *rps.Name, *rps.Owner.Login)
	return nil
}

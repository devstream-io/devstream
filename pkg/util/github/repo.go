package github

import (
	"net/http"
	"strings"

	"github.com/google/go-github/v42/github"

	"github.com/devstream-io/devstream/pkg/util/log"
)

func (c *Client) CreateRepo(org string) error {
	repo := &github.Repository{
		Name: &c.Repo,
	}

	if org != "" {
		log.Infof("Prepare to create an organization repository: %s/%s", org, repo.GetName())
	}
	_, _, err := c.Repositories.Create(c.Context, org, repo)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) DeleteRepo() error {
	var owner = c.Owner
	if c.Org != "" {
		owner = c.Org
	}

	response, err := c.Client.Repositories.Delete(c.Context, owner, c.Repo)

	// error reason is not 404
	if err != nil && !strings.Contains(err.Error(), "404") {
		log.Errorf("Delete repo failed: %s.", err)
		return err
	}

	if response.StatusCode == http.StatusNotFound {
		log.Infof("GitHub repo %s was not found. Nothing to do here.", c.Repo)
		return nil
	}

	log.Successf("GitHub repo %s removed.", c.Repo)
	return nil
}

func (c *Client) GetRepoDescription() (*github.Repository, error) {
	var owner = c.Owner
	if c.Org != "" {
		owner = c.Org
	}

	repo, resp, err := c.Client.Repositories.Get(
		c.Context,
		owner,
		c.Repo)

	if err != nil {
		return nil, err
	}

	if repo == nil && resp.StatusCode == http.StatusNotFound {
		return repo, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, err
	}

	return repo, nil
}

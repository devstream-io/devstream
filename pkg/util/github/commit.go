package github

import (
	"fmt"

	"github.com/google/go-github/v42/github"

	"github.com/devstream-io/devstream/pkg/util/log"
)

func (c *Client) GetLastCommit() (*github.RepositoryCommit, error) {
	var owner = c.Owner
	if c.Org != "" {
		owner = c.Org
	}

	commits, _, err := c.Client.Repositories.ListCommits(c.Context, owner, c.Repo, &github.CommitsListOptions{})
	if err != nil {
		log.Debugf("Failed to get RepositoryCommits: %s.", err)
		return nil, err
	}

	if len(commits) == 0 {
		msg := "no commits was found"
		log.Info(msg)
		return nil, fmt.Errorf(msg)
	}

	return commits[0], nil
}

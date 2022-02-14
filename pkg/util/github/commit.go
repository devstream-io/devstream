package github

import (
	"fmt"

	"github.com/google/go-github/v42/github"

	"github.com/merico-dev/stream/internal/pkg/log"
)

func (c *Client) GetLastCommit() (*github.RepositoryCommit, error) {
	commits, _, err := c.Client.Repositories.ListCommits(c.Context, c.Owner, c.Repo, &github.CommitsListOptions{})
	if err != nil {
		log.Debugf("failed to get RepositoryCommits: %s", err)
		return nil, err
	}

	if len(commits) == 0 {
		msg := "no commits was found"
		log.Info(msg)
		return nil, fmt.Errorf(msg)
	}

	return commits[0], nil
}

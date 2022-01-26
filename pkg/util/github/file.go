package github

import (
	"context"

	"github.com/google/go-github/v42/github"
)

func (c *Client) CreateFile(content []byte, filePath string) error {
	defaultMsg := "initialize"
	defaultBranch := "main"

	opt := &github.RepositoryContentFileOptions{
		Message: &defaultMsg,
		Content: content,
		Branch:  &defaultBranch,
	}

	_, _, err := c.Repositories.CreateFile(context.TODO(), c.Owner, c.Repo, filePath, opt)
	return err
}

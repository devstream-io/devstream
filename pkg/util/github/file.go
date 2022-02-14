package github

import (
	"github.com/google/go-github/v42/github"
)

func (c *Client) CreateFile(content []byte, filePath, targetBranch string) error {
	defaultMsg := "Initialize the repository"

	opt := &github.RepositoryContentFileOptions{
		Message: &defaultMsg,
		Content: content,
		Branch:  &targetBranch,
	}

	_, _, err := c.Repositories.CreateFile(c.Context, c.Owner, c.Repo, filePath, opt)
	return err
}

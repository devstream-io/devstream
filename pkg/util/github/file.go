package github

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/google/go-github/v42/github"

	"github.com/devstream-io/devstream/pkg/util/log"
)

func (c *Client) CreateFile(content []byte, filePath, targetBranch string) error {
	defaultMsg := "Initialize the repository"

	opt := &github.RepositoryContentFileOptions{
		Message: &defaultMsg,
		Content: content,
		Branch:  &targetBranch,
	}

	var owner = c.Owner
	if c.Org != "" {
		owner = c.Org
	}

	_, _, err := c.Repositories.CreateFile(c.Context, owner, c.Repo, filePath, opt)
	return err
}

func (c *Client) PushLocalPath(repoPath, branch string) error {
	if err := filepath.Walk(repoPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			log.Debugf("Walk error: %s.", err)
			return err
		}

		if info.IsDir() {
			log.Debugf("Found dir: %s.", path)
			return nil
		}

		log.Debugf("Found file: %s.", path)

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		githubPath, _ := filepath.Rel(repoPath, path)
		return c.CreateFile(content, githubPath, branch)
	}); err != nil {
		return err
	}

	return nil
}

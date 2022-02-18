package gitlab

import (
	"github.com/xanzy/go-gitlab"
)

// https://docs.gitlab.com/ee/api/repository_files.html
// https://github.com/xanzy/go-gitlab/blob/master/repository_files.go

func (c *Client) FileExists(project, branch, filename string) error {
	getFileOptions := &gitlab.GetFileOptions{
		Ref: gitlab.String(branch),
	}

	_, _, err := c.RepositoryFiles.GetFile(project, filename, getFileOptions)
	if err != nil {
		return err
	}

	return nil
}

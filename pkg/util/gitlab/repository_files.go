package gitlab

import (
	"github.com/xanzy/go-gitlab"
)

// https://docs.gitlab.com/ee/api/repository_files.html
// https://github.com/xanzy/go-gitlab/blob/master/repository_files.go

func (c *Client) FileExists(project, branch, filename string) (bool, error) {
	getFileOptions := &gitlab.GetFileOptions{
		Ref: gitlab.String(branch),
	}

	_, response, err := c.RepositoryFiles.GetFile(project, filename, getFileOptions)
	if response.StatusCode == 404 {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

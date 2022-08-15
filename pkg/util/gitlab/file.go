package gitlab

import (
	"net/http"

	"github.com/xanzy/go-gitlab"
)

// https://docs.gitlab.com/ee/api/repository_files.html
// https://github.com/xanzy/go-gitlab/blob/master/repository_files.go

func (c *Client) FileExists(filename string) (bool, error) {
	getFileOptions := &gitlab.GetFileOptions{
		Ref: gitlab.String(c.Branch),
	}

	_, response, err := c.RepositoryFiles.GetFile(c.GetRepoPath(), filename, getFileOptions)
	for _, v := range []int{http.StatusBadRequest, http.StatusUnauthorized, http.StatusNotFound} {
		if response.StatusCode == v {
			return false, nil
		}
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

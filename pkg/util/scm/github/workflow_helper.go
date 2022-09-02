package github

import (
	"fmt"
	"net/http"

	"github.com/google/go-github/v42/github"
)

func generateGitHubWorkflowFileByName(f string) string {
	return fmt.Sprintf(".github/workflows/%s", f)
}

// getFileSHA will try to collect the SHA hash value of the file, then return it. the return values will be:
// 1. If file exists without error -> string(SHA), nil
// 2. If some errors occurred -> return "", err
// 3. If file not found without error -> return "", nil
func (c *Client) getFileSHA(filename string) (string, error) {
	content, _, resp, err := c.Client.Repositories.GetContents(
		c.Context,
		c.GetRepoOwner(),
		c.Repo,
		generateGitHubWorkflowFileByName(filename),
		&github.RepositoryContentGetOptions{},
	)

	if resp != nil && resp.StatusCode == http.StatusNotFound {
		return "", nil
	}

	// error reason is not 404
	if err != nil {
		return "", err
	}

	return *content.SHA, nil
}

package github

import (
	"fmt"
	"net/http"
	"strings"

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
	var owner = c.Owner
	if c.Org != "" {
		owner = c.Org
	}

	content, _, resp, err := c.Client.Repositories.GetContents(
		c.Context,
		owner,
		c.Option.Repo,
		generateGitHubWorkflowFileByName(filename),
		&github.RepositoryContentGetOptions{},
	)

	// error reason is not 404
	if err != nil && !strings.Contains(err.Error(), "404") {
		return "", err
	}

	// error reason is 404
	if resp.StatusCode == http.StatusNotFound {
		return "", nil
	}

	// no error occurred
	if resp.StatusCode == http.StatusOK {
		return *content.SHA, nil
	}
	return "", fmt.Errorf("got some error is not expected")
}

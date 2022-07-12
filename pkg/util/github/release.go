package github

import (
	"context"
	"fmt"
	"net/http"

	"github.com/devstream-io/devstream/pkg/util/log"
)

func (c *Client) LatestReleaseTagName() (string, error) {
	ltstRelease, resp, err := c.Repositories.GetLatestRelease(context.Background(), c.Org, c.Repo)
	if err != nil {
		return "", err
	}

	log.Debugf("Response status: %s.", resp.Status)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("got response status not expected: %s", resp.Status)
	}

	return *ltstRelease.TagName, nil
}

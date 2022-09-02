package github

import (
	"context"
)

func (c *Client) GetLatestReleaseTagName() (string, error) {
	release, _, err := c.Repositories.GetLatestRelease(context.Background(), c.Org, c.Repo)
	if err != nil {
		return "", err
	}

	return *release.TagName, nil
}

package gitlab

import (
	"github.com/xanzy/go-gitlab"

	"github.com/devstream-io/devstream/pkg/util/log"
)

// https://docs.gitlab.com/ee/api/projects.html
// https://github.com/xanzy/go-gitlab/blob/master/projects.go

func (c *Client) CreateProject(repoName, branch string) error {
	log.Debugf("Repo to be created: %s", repoName)

	p := &gitlab.CreateProjectOptions{
		Name:                 gitlab.String(repoName),
		Description:          gitlab.String("Bootstrapped by DevStream."),
		MergeRequestsEnabled: gitlab.Bool(true),
		SnippetsEnabled:      gitlab.Bool(true),
		Visibility:           gitlab.Visibility(gitlab.PublicVisibility),
		DefaultBranch:        gitlab.String(branch),
	}

	_, _, err := c.Projects.CreateProject(p)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) DeleteProject(project string) error {
	_, err := c.Projects.DeleteProject(project)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) DescribeProject(project string) (*gitlab.Project, error) {
	p := &gitlab.GetProjectOptions{}
	res, _, err := c.Projects.GetProject(project, p)

	if err != nil {
		return nil, err
	}

	return res, nil
}

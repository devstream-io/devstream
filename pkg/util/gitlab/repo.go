package gitlab

import (
	"github.com/xanzy/go-gitlab"

	"github.com/devstream-io/devstream/pkg/util/log"
)

func (c *Client) InitRepo() error {
	log.Debugf("Repo to be created: %s", c.Repo)

	var err error
	var res *gitlab.Group
	var groupId int

	gitlabGetGroupOptions := &gitlab.GetGroupOptions{}

	if c.Namespace != "" {
		res, _, err = c.Groups.GetGroup(c.Namespace, gitlabGetGroupOptions)
		if err != nil {
			return err
		}
		groupId = res.ID
	}

	log.Debugf("Group: %#v\n", res)

	p := &gitlab.CreateProjectOptions{
		Name:                 gitlab.String(c.Repo),
		Description:          gitlab.String("Bootstrapped by DevStream."),
		MergeRequestsEnabled: gitlab.Bool(true),
		SnippetsEnabled:      gitlab.Bool(true),
		DefaultBranch:        gitlab.String(c.Branch),
	}

	switch c.Visibility {
	case "public":
		p.Visibility = gitlab.Visibility(gitlab.PublicVisibility)
	case "internal":
		p.Visibility = gitlab.Visibility(gitlab.InternalVisibility)
	case "private":
		p.Visibility = gitlab.Visibility(gitlab.PrivateVisibility)
	default:
		p.Visibility = gitlab.Visibility(gitlab.PublicVisibility)
	}

	if groupId != 0 {
		p.NamespaceID = gitlab.Int(groupId)
	}
	_, _, err = c.Projects.CreateProject(p)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) DeleteRepo() error {
	_, err := c.Projects.DeleteProject(c.GetRepoPath())
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) DescribeRepo() (*gitlab.Project, error) {
	p := &gitlab.GetProjectOptions{}
	res, _, err := c.Projects.GetProject(c.GetRepoPath(), p)

	if err != nil {
		log.Debugf("gitlab project: get [%s] info error %s", c.GetRepoPath(), err)
		return nil, err
	}

	return res, nil
}

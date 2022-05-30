package gitlab

import (
	"github.com/xanzy/go-gitlab"

	"github.com/devstream-io/devstream/pkg/util/log"
)

// https://docs.gitlab.com/ee/api/projects.html
// https://github.com/xanzy/go-gitlab/blob/master/projects.go

type CreateProjectOptions struct {
	Name       string
	Namespace  string
	Branch     string
	Visibility string
}

// func (c *Client) CreateProject(repoName, branch string) error {
func (c *Client) CreateProject(opts *CreateProjectOptions) error {
	log.Debugf("Repo to be created: %s", opts.Name)

	var err error
	var res *gitlab.Group
	var groupId int

	gitlabGetGroupOptions := &gitlab.GetGroupOptions{}

	if opts.Namespace != "" {
		res, _, err = c.Groups.GetGroup(opts.Namespace, gitlabGetGroupOptions)
		if err != nil {
			return err
		}
		groupId = res.ID
	}

	log.Debugf("Group: %#v\n", res)

	p := &gitlab.CreateProjectOptions{
		Name:                 gitlab.String(opts.Name),
		Description:          gitlab.String("Bootstrapped by DevStream."),
		MergeRequestsEnabled: gitlab.Bool(true),
		SnippetsEnabled:      gitlab.Bool(true),
		DefaultBranch:        gitlab.String(opts.Branch),
	}

	switch opts.Visibility {
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

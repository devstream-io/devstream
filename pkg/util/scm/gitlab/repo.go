package gitlab

import (
	"strings"

	"github.com/xanzy/go-gitlab"

	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/scm/git"
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
		if strings.Contains(err.Error(), "has already been taken") {
			log.Debugf("gitlab repo %s already exist, not create...", c.Repo)
			return nil
		}
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

func (c *Client) AddWebhook(webhookConfig *git.WebhookConfig) error {
	projectHook, err := c.getWebhook(webhookConfig)
	if err != nil {
		return err
	}
	if projectHook == nil {
		p := &gitlab.AddProjectHookOptions{
			PushEvents: gitlab.Bool(true),
			Token:      gitlab.String(webhookConfig.SecretToken),
			URL:        gitlab.String(webhookConfig.Address),
		}
		_, _, err := c.Projects.AddProjectHook(c.GetRepoPath(), p)
		return err
	}
	log.Debugf("gitlab AddWebhook already exist")
	return nil
}

func (c *Client) DeleteWebhook(webhookConfig *git.WebhookConfig) error {
	projectHook, err := c.getWebhook(webhookConfig)
	if err != nil {
		log.Debugf("gitlab DeleteWebhook list hooks failed: %s", err)
		return err
	}
	if projectHook != nil {
		_, err := c.Projects.DeleteProjectHook(c.GetRepoPath(), projectHook.ID)
		return err
	}
	log.Infof("gitlab DeleteWebhook not found")
	return nil
}

func (c *Client) getWebhook(webhookConfig *git.WebhookConfig) (*gitlab.ProjectHook, error) {
	p := &gitlab.ListProjectHooksOptions{}
	hooks, _, err := c.Projects.ListProjectHooks(c.GetRepoPath(), p)
	if err != nil {
		log.Debugf("gitlab DeleteWebhook lsit hooks failed: %s", err)
		return nil, err
	}
	for _, hook := range hooks {
		if hook.URL == webhookConfig.Address {
			return hook, nil
		}
	}
	return nil, nil
}

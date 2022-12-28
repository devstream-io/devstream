package gitlab

import (
	"fmt"

	"github.com/xanzy/go-gitlab"

	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/pkgerror"
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
		AutoDevopsEnabled:    gitlab.Bool(false),
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
	if err != nil && !pkgerror.CheckErrorMatchByMessage(err, errRepoNotFound, errRepoExist) {
		return err
	}
	return nil
}

func (c *Client) DeleteRepo() error {
	_, err := c.Projects.DeleteProject(c.GetRepoPath())
	if err != nil && !pkgerror.CheckErrorMatchByMessage(err, errRepoNotFound) {
		return err
	}
	return nil
}

func (c *Client) DescribeRepo() (*git.RepoInfo, error) {
	p := &gitlab.GetProjectOptions{}
	project, _, err := c.Projects.GetProject(c.GetRepoPath(), p)

	if err != nil {
		log.Debugf("gitlab project: [%s] error %+v", c.GetRepoPath(), err)
		return nil, c.newModuleError(err)
	}

	log.Debugf("GitLab Project is: %#v\n", project)
	repoInfo := &git.RepoInfo{
		Repo:     project.Name,
		CloneURL: git.ScmURL(project.HTTPURLToRepo),
	}
	if project.Owner != nil {
		log.Debugf("GitLab project owner is: %#v.\n", project.Owner)
		repoInfo.Owner = project.Owner.Username
		repoInfo.Org = project.Owner.Organization
	}
	return repoInfo, nil
}

// AddWebhook will update webhook when it exists
// else create a webbhook
func (c *Client) AddWebhook(webhookConfig *git.WebhookConfig) error {
	projectHook, err := c.getWebhook(webhookConfig)
	if err != nil {
		return err
	}
	if projectHook != nil {
		log.Debugf("gitlab AddWebhook already exist, update this webhook")
		p := &gitlab.EditProjectHookOptions{
			PushEvents:          gitlab.Bool(true),
			Token:               gitlab.String(webhookConfig.SecretToken),
			URL:                 gitlab.String(webhookConfig.Address),
			MergeRequestsEvents: gitlab.Bool(true),
		}
		_, _, err = c.Projects.EditProjectHook(c.GetRepoPath(), projectHook.ID, p)
	} else {
		p := &gitlab.AddProjectHookOptions{
			PushEvents:          gitlab.Bool(true),
			Token:               gitlab.String(webhookConfig.SecretToken),
			URL:                 gitlab.String(webhookConfig.Address),
			MergeRequestsEvents: gitlab.Bool(true),
		}
		_, _, err = c.Projects.AddProjectHook(c.GetRepoPath(), p)
	}
	if err != nil {
		return c.newModuleError(err)
	}
	return nil
}

func (c *Client) DeleteWebhook(webhookConfig *git.WebhookConfig) error {
	projectHook, err := c.getWebhook(webhookConfig)
	if err != nil && !pkgerror.CheckErrorMatchByMessage(err, errRepoNotFound) {
		return err
	}
	if projectHook == nil {
		log.Debugf("gitlab DeleteWebhook not found")
		return nil
	}
	_, err = c.Projects.DeleteProjectHook(c.GetRepoPath(), projectHook.ID)
	if err != nil {
		return c.newModuleError(err)
	}
	return nil
}

func (c *Client) getWebhook(webhookConfig *git.WebhookConfig) (*gitlab.ProjectHook, error) {
	p := &gitlab.ListProjectHooksOptions{}
	hooks, _, err := c.Projects.ListProjectHooks(c.GetRepoPath(), p)
	if err != nil {
		log.Debugf("gitlab get webhook list hooks failed: %s", err)
		return nil, c.newModuleError(err)
	}
	for _, hook := range hooks {
		if hook.URL == webhookConfig.Address {
			return hook, nil
		}
	}
	return nil, nil
}

// TODO(steinliber): support gtlab later
func (c *Client) DownloadRepo() (string, error) {
	return "", fmt.Errorf("gitlab doesn't support download repo for now")
}

func (c *Client) AddRepoSecret(secretKey, secretValue string) error {
	var err error
	createOpts := gitlab.CreateProjectVariableOptions{
		Key:    gitlab.String(secretKey),
		Value:  gitlab.String(secretValue),
		Masked: gitlab.Bool(true),
	}
	_, _, err = c.ProjectVariables.CreateVariable(c.GetRepoPath(), &createOpts)
	// if secret already exist, just update this secret
	if err != nil && pkgerror.CheckErrorMatchByMessage(err, errVariableExist) {
		updateOpts := gitlab.UpdateProjectVariableOptions{
			Value:  gitlab.String(secretValue),
			Masked: gitlab.Bool(true),
		}
		_, _, err = c.ProjectVariables.UpdateVariable(c.GetRepoPath(), secretKey, &updateOpts)
	}
	return err
}

func (c *Client) ListRepoRunner() ([]*gitlab.Runner, error) {
	listOpts := &gitlab.ListProjectRunnersOptions{
		Status:  gitlab.String("online"),
		TagList: &[]string{"ci"},
	}
	runners, _, err := c.Runners.ListProjectRunners(c.GetRepoPath(), listOpts)
	return runners, err

}

func (c *Client) ResetRepoRunnerToken() (string, error) {
	token, _, err := c.Runners.ResetProjectRunnerRegistrationToken(c.GetRepoPath())
	if err != nil {
		return "", err
	}
	return *token.Token, nil
}

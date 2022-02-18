package gitlab

// https://docs.gitlab.com/ee/api/templates/gitlab_ci_ymls.html
// https://github.com/xanzy/go-gitlab/blob/master/ci_yml_templates.go

func (c *Client) GetGitLabCIGolangTemplate() (string, error) {
	tpl, _, err := c.CIYMLTemplate.GetTemplate("Go")
	if err != nil {
		return "", err
	}

	return tpl.Content, nil
}

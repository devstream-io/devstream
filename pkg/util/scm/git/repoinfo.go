package git

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/devstream-io/devstream/pkg/util/config/configGetter"
	"github.com/devstream-io/devstream/pkg/util/log"
)

type RepoInfo struct {
	Owner    string `validate:"required_without=Org" mapstructure:"owner"`
	Org      string `validate:"required_without=Owner" mapstructure:"org"`
	Repo     string `validate:"required" mapstructure:"repo"`
	Branch   string `mapstructure:"branch"`
	RepoType string `validate:"oneof=gitlab github" mapstructure:"repoType"`
	// This is config for gitlab
	CloneURL      string `mapstructure:"cloneURL"`
	SSHPrivateKey string `mapstructure:"sshPrivateKey"`

	// used for gitlab
	Namespace  string
	BaseURL    string `validate:"omitempty,url" mapstructure:"baseURL"`
	Visibility string `validate:"omitempty,oneof=public private internal" mapstructure:"visibility"`

	// used for GitHub
	WorkPath string
	NeedAuth bool
}

func (r *RepoInfo) GetRepoOwner() string {
	if r.Org != "" {
		return r.Org
	}
	return r.Owner
}

func (r *RepoInfo) GetRepoPath() string {
	return fmt.Sprintf("%s/%s", r.GetRepoOwner(), r.Repo)
}

// BuildRepoRenderConfig will generate template render variables
func (r *RepoInfo) BuildRepoRenderConfig() map[string]interface{} {
	renderConfigMap := map[string]interface{}{
		"AppName": r.Repo,
		"Repo": map[string]string{
			"Name":  r.Repo,
			"Owner": r.GetRepoOwner(),
		},
	}
	return renderConfigMap
}

func (r *RepoInfo) GetRepoNameWithBranch() string {
	return fmt.Sprintf("%s-%s", r.Repo, r.GetBranchWithDefault())
}

func (r *RepoInfo) GetBranchWithDefault() string {
	branch := r.Branch
	if branch != "" {
		return branch
	}
	switch r.RepoType {
	case "github":
		branch = "main"
	case "gitlab":
		branch = "master"
	}
	return branch
}

func (r *RepoInfo) UpdateRepoPathByCloneURL(cloneURL string) error {
	var paths string
	c, err := url.ParseRequestURI(cloneURL)
	if err != nil {
		if strings.Contains(cloneURL, "git@") {
			gitSSHLastIndex := strings.LastIndex(cloneURL, ":")
			if gitSSHLastIndex == -1 {
				return fmt.Errorf("git ssh repo not valid")
			}
			paths = strings.Trim(cloneURL[gitSSHLastIndex:], ":")
		} else {
			return fmt.Errorf("git repo transport not support for now")
		}
	} else {
		paths = c.Path
	}
	projectPaths := strings.Split(strings.Trim(paths, "/"), "/")
	if len(projectPaths) != 2 {
		return fmt.Errorf("git repo path is not valid")
	}
	r.Owner = projectPaths[0]
	r.Repo = strings.TrimSuffix(projectPaths[1], ".git")
	return nil
}

// BuildURL return url build from repo struct
func (r *RepoInfo) BuildScmURL() string {
	switch r.RepoType {
	case "github":
		return fmt.Sprintf("https://github.com/%s/%s", r.GetRepoOwner(), r.Repo)
	case "gitlab":
		var gitlabURL string
		if r.BaseURL != "" {
			gitlabURL = r.BaseURL
		} else {
			gitlabURL = "https://gitlab.com"
		}
		return fmt.Sprintf("%s/%s/%s.git", gitlabURL, r.GetRepoOwner(), r.Repo)
	default:
		return ""
	}
}

func (r *RepoInfo) CheckValid() error {
	switch r.RepoType {
	case "gitlab":
		return configGetter.CheckItemExist(configGetter.NewEnvGetter("GITLAB_TOKEN"))
	case "github":
		return configGetter.CheckItemExist(configGetter.NewEnvGetter("GITHUB_TOKEN"))
	}
	return nil
}

func (r *RepoInfo) BuildWebhookInfo(baseURL, appName, token string) *WebhookConfig {
	var webHookURL string
	switch r.RepoType {
	case "gitlab":
		webHookURL = fmt.Sprintf("%s/project/%s", baseURL, appName)
	case "github":
		webHookURL = fmt.Sprintf("%s/github-webhook/", baseURL)
	}
	log.Debugf("jenkins config webhook is %s", webHookURL)
	return &WebhookConfig{
		Address:     webHookURL,
		SecretToken: token,
	}
}

func (r *RepoInfo) IsGitlab() bool {
	return r.RepoType == "gitlab"
}

func (r *RepoInfo) IsGithub() bool {
	return r.RepoType == "github"
}

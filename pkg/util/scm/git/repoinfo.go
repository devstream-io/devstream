package git

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/mapz"
)

type RepoInfo struct {
	Owner    string `validate:"required_without=Org" mapstructure:"owner,omitempty"`
	Org      string `validate:"required_without=Owner" mapstructure:"org,omitempty"`
	Repo     string `validate:"required" mapstructure:"repo,omitempty"`
	Branch   string `mapstructure:"branch,omitempty"`
	RepoType string `validate:"oneof=gitlab github" mapstructure:"repoType,omitempty"`
	// This is config for gitlab
	CloneURL      string `mapstructure:"url,omitempty"`
	SSHPrivateKey string `mapstructure:"sshPrivateKey,omitempty"`

	// used for gitlab
	Namespace  string `mapstructure:"nameSpace,omitempty"`
	BaseURL    string `mapstructure:"baseURL,omitempty"`
	Visibility string `mapstructure:"visibility,omitempty"`

	// used for GitHub
	WorkPath string `mapstructure:"workPath,omitempty"`
	NeedAuth bool   `mapstructure:"needAuth,omitempty"`
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
				return fmt.Errorf("scm git ssh repo not valid")
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
		if os.Getenv("GITLAB_TOKEN") == "" {
			return fmt.Errorf("pipeline gitlab should set env GITLAB_TOKEN")
		}
	case "github":
		if os.Getenv("GITHUB_TOKEN") == "" {
			return fmt.Errorf("pipeline github should set env GITHUB_TOKEN")
		}
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

func (r *RepoInfo) Encode() map[string]any {
	m, err := mapz.DecodeStructToMap(r)
	if err != nil {
		log.Errorf("gitRepo [%+v] decode to map failed: %+v", r, err)
	}
	return m
}

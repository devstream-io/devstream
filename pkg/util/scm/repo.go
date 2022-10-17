package scm

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/devstream-io/devstream/pkg/util/scm/git"
	"github.com/devstream-io/devstream/pkg/util/scm/gitlab"
)

type Repo struct {
	Owner    string `validate:"required_without=Org" mapstructure:"owner"`
	Org      string `validate:"required_without=Owner" mapstructure:"org"`
	Repo     string `validate:"required" mapstructure:"repo"`
	Branch   string `mapstructure:"branch"`
	RepoType string `validate:"oneof=gitlab github" mapstructure:"repoType"`
	// This is config for gitlab
	BaseURL    string `validate:"omitempty,url" mapstructure:"baseURL"`
	Visibility string `validate:"omitempty,oneof=public private internal" mapstructure:"visibility"`
}

// BuildRepoRenderConfig will generate template render variables
func (d *Repo) BuildRepoRenderConfig() map[string]interface{} {
	renderConfigMap := map[string]interface{}{
		"AppName": d.Repo,
		"Repo": map[string]string{
			"Name":  d.Repo,
			"Owner": d.BuildRepoInfo().GetRepoOwner(),
		},
	}
	return renderConfigMap
}

// BuildRepoInfo is used to return git.RepoInfo for init scm client
func (d *Repo) BuildRepoInfo() *git.RepoInfo {
	return &git.RepoInfo{
		Repo:       d.Repo,
		Owner:      d.Owner,
		Org:        d.Org,
		Visibility: d.Visibility,
		NeedAuth:   true,
		Branch:     d.getBranch(),
		BaseURL:    d.BaseURL,
		Type:       d.RepoType,
	}
}

func (d *Repo) GetRepoNameWithBranch() string {
	return fmt.Sprintf("%s-%s", d.Repo, d.getBranch())
}

func (d *Repo) getBranch() string {
	branch := d.Branch
	if branch != "" {
		return branch
	}
	switch d.RepoType {
	case "github":
		branch = "main"
	case "gitlab":
		branch = "master"
	}
	return branch
}

func (d *Repo) UpdateRepoPathByCloneURL(cloneURL string) error {
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
	d.Owner = projectPaths[0]
	d.Repo = strings.TrimSuffix(projectPaths[1], ".git")
	return nil
}

// BuildURL return url build from repo struct
func (d *Repo) BuildScmURL() string {
	repoInfo := d.BuildRepoInfo()
	switch d.RepoType {
	case "github":
		return fmt.Sprintf("https://github.com/%s/%s", repoInfo.GetRepoOwner(), d.Repo)
	case "gitlab":
		var gitlabURL string
		if d.BaseURL != "" {
			gitlabURL = d.BaseURL
		} else {
			gitlabURL = gitlab.DefaultGitlabHost
		}
		return fmt.Sprintf("%s/%s/%s.git", gitlabURL, repoInfo.GetRepoOwner(), d.Repo)
	default:
		return ""
	}
}

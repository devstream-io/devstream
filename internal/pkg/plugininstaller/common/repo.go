package common

import (
	"fmt"
	"path/filepath"

	"github.com/devstream-io/devstream/pkg/util/file"
	"github.com/devstream-io/devstream/pkg/util/github"
	"github.com/devstream-io/devstream/pkg/util/gitlab"
	"github.com/devstream-io/devstream/pkg/util/log"
)

const (
	transitBranch    = "init-with-devstream"
	defaultCommitMsg = "init with devstream"
)

// Repo is the repo info of github or gitlab
type Repo struct {
	Owner             string `validate:"required_without=Org" mapstructure:"owner"`
	Org               string `validate:"required_without=Owner" mapstructure:"org"`
	Repo              string `validate:"required" mapstructure:"repo"`
	Branch            string `mapstructure:"branch"`
	PathWithNamespace string
	RepoType          string `validate:"oneof=gitlab github" mapstructure:"repo_type"`
	// This is config for gitlab
	BaseURL    string `validate:"omitempty,url" mapstructure:"base_url"`
	Visibility string `validate:"omitempty,oneof=public private internal" mapstructure:"visibility"`
}

// BuildRepoRenderConfig will generate template render variables
func (d *Repo) BuildRepoRenderConfig() map[string]interface{} {
	renderConfigMap := map[string]interface{}{
		"AppName": d.Repo,
		"Repo": map[string]string{
			"Name":  d.Repo,
			"Owner": d.getRepoOwner(),
		},
	}
	return renderConfigMap
}

// CreateGithubClient build github client connection info
func (d *Repo) CreateGithubClient(needAuth bool) (*github.Client, error) {
	ghOptions := &github.Option{
		Owner:    d.Owner,
		Org:      d.Org,
		Repo:     d.Repo,
		NeedAuth: needAuth,
	}
	ghClient, err := github.NewClient(ghOptions)
	if err != nil {
		return nil, err
	}
	return ghClient, nil
}

// CreateGitlabClient build gitlab connection info
func (d *Repo) CreateGitlabClient() (*gitlab.Client, error) {
	return gitlab.NewClient(gitlab.WithBaseURL(d.BaseURL))
}

// buildgitlabOpts build gitlab connection options
func (d *Repo) buildgitlabOpts() *gitlab.CreateProjectOptions {
	return &gitlab.CreateProjectOptions{
		Name:       d.Repo,
		Branch:     d.getBranch(),
		Namespace:  d.Org,
		Visibility: d.Visibility,
	}
}

// CreateAndRenderLocalRepo will download repo from source repo and render it locally
func (d *Repo) CreateAndRenderLocalRepo(appName string, vars map[string]interface{}) (string, error) {
	//TODO(steinliber) support gtlab later
	if d.RepoType != "github" {
		return "", fmt.Errorf("the download target repo is currently only supported github")
	}
	// 1. download zip file and unzip this file then render folders
	downloadURL := d.GetGithubDownloadURL()
	projectDir, err := file.NewTemplate().FromRemote(downloadURL).UnzipFile().RenderRepoDir(
		appName, vars,
	).Run()
	if err != nil {
		log.Debugf("reposcaffolding process file error: %s", err)
		return "", err
	}
	// 2. join download path and repo name to get repo path
	repoDirName := d.getRepoNameWithBranch()
	return filepath.Join(projectDir, repoDirName), nil
}

// This Func will push repo to remote base on repoType
func (d *Repo) Push(repoPath string) error {
	switch d.RepoType {
	case "github":
		ghClient, err := d.CreateGithubClient(true)
		if err != nil {
			return err
		}
		return ghClient.PushInitRepo(transitBranch, d.getBranch(), repoPath)
	case "gitlab":
		c, err := d.CreateGitlabClient()
		if err != nil {
			return err
		}
		return c.PushInitRepo(d.buildgitlabOpts(), d.PathWithNamespace, repoPath, defaultCommitMsg)
	default:
		return fmt.Errorf("scaffolding not support repo destination: %s", d.RepoType)
	}
}

func (d *Repo) Delete() error {
	switch d.RepoType {
	case "github":
		// 1. create ghClient
		ghClient, err := d.CreateGithubClient(true)
		if err != nil {
			return err
		}
		// 2. delete github repo
		return ghClient.DeleteRepo()
	case "gitlab":
		gLclient, err := d.CreateGitlabClient()
		if err != nil {
			return err
		}
		return gLclient.DeleteProject(d.PathWithNamespace)
	default:
		return fmt.Errorf("scaffolding not support repo destination: %s", d.RepoType)
	}
}

func (d *Repo) getBranch() string {
	if d.Branch != "" {
		return d.Branch
	}
	return "main"
}

func (d *Repo) getRepoNameWithBranch() string {
	return fmt.Sprintf("%s-%s", d.Repo, d.getBranch())
}

func (d *Repo) GetGithubDownloadURL() string {
	latestCodeZipfileDownloadURL := fmt.Sprintf(
		github.DefaultLatestCodeZipfileDownloadUrlFormat, d.getRepoOwner(), d.Repo, d.getBranch(),
	)
	log.Debugf("LatestCodeZipfileDownloadUrl: %s.", latestCodeZipfileDownloadURL)
	return latestCodeZipfileDownloadURL
}

func (d *Repo) getRepoOwner() string {
	if d.Org != "" {
		return d.Org
	}
	return d.Owner
}

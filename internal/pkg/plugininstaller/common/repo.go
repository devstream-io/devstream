package common

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/devstream-io/devstream/pkg/util/file"
	"github.com/devstream-io/devstream/pkg/util/git"
	"github.com/devstream-io/devstream/pkg/util/github"
	"github.com/devstream-io/devstream/pkg/util/gitlab"
	"github.com/devstream-io/devstream/pkg/util/log"
)

const (
	defaultMainBranch = "main"
	transitBranch     = "init-with-devstream"
	defaultCommitMsg  = "init with devstream"
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
			"Owner": d.BuildRepoInfo().GetRepoOwner(),
		},
	}
	return renderConfigMap
}

// CreateGithubClient build github client connection info
func (d *Repo) NewClient() (git.ClientOperation, error) {
	repoInfo := d.BuildRepoInfo()
	switch d.RepoType {
	case "github":
		return github.NewClient(repoInfo)
	case "gitlab":
		return gitlab.NewClient(repoInfo)
	}
	return nil, fmt.Errorf("scaffolding not support repo destination: %s", d.RepoType)
}

// CreateAndRenderLocalRepo will download repo from source repo and render it locally
func (d *Repo) CreateAndRenderLocalRepo(appName string, vars map[string]interface{}) (git.GitFileContentMap, error) {
	//TODO(steinliber) support gtlab later
	if d.RepoType != "github" {
		return nil, fmt.Errorf("the download target repo is currently only supported github")
	}
	// 1. download zip file and unzip this file then render folders
	downloadURL := d.getRepoDownloadURL()
	zipFilesDir, err := file.DownloadAndUnzipFile(downloadURL)
	if err != nil {
		log.Debugf("reposcaffolding process files error: %s", err)
		return nil, err
	}
	return file.WalkDir(
		zipFilesDir, filterGitFiles,
		getRepoFileNameFunc(appName, d.BuildRepoInfo().GetRepoNameWithBranch()),
		processRepoFileFunc(appName, vars),
	)
}

// This Func will push repo to remote base on repoType
func (d *Repo) CreateAndPush(gitMap git.GitFileContentMap) error {
	client, err := d.NewClient()
	if err != nil {
		return err
	}
	commitInfo := &git.CommitInfo{
		CommitMsg:    defaultCommitMsg,
		CommitBranch: transitBranch,
		GitFileMap:   gitMap,
	}
	return git.PushInitRepo(client, commitInfo)
}

func (d *Repo) Delete() error {
	client, err := d.NewClient()
	if err != nil {
		return err
	}
	return client.DeleteRepo()
}

func (d *Repo) BuildRepoInfo() *git.RepoInfo {
	branch := d.Branch
	if branch == "" {
		branch = defaultMainBranch
	}
	return &git.RepoInfo{
		Repo:       d.Repo,
		Owner:      d.Owner,
		Org:        d.Org,
		Visibility: d.Visibility,
		NeedAuth:   true,
		Branch:     branch,
		BaseURL:    d.BaseURL,
	}
}

func (d *Repo) getRepoDownloadURL() string {
	repoInfo := d.BuildRepoInfo()
	latestCodeZipfileDownloadURL := fmt.Sprintf(
		github.DefaultLatestCodeZipfileDownloadUrlFormat, repoInfo.GetRepoOwner(), repoInfo.Repo, repoInfo.Branch,
	)
	log.Debugf("LatestCodeZipfileDownloadUrl: %s.", latestCodeZipfileDownloadURL)
	return latestCodeZipfileDownloadURL
}

func (d *Repo) GetRepoToken() string {
	var envKey string
	switch d.RepoType {
	case "github":
		envKey = "GITHUB_TOKEN"
	case "gitlab":
		envKey = "GITLAB_TOKEN"
	default:
		return ""
	}
	return os.Getenv(envKey)
}

func (d *Repo) BuildURL() string {
	repoInfo := d.BuildRepoInfo()
	switch d.RepoType {
	case "github":
		return fmt.Sprintf("https://github.com/%s/%s.git", repoInfo.GetRepoOwner(), d.Repo)
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

func (d *Repo) AddWebHook(webhookInfo *git.WebhookConfig) error {
	client, err := d.NewClient()
	if err != nil {
		return err
	}
	return client.AddWebhook(webhookInfo)
}

func (d *Repo) DeleteWebhook(webhookInfo *git.WebhookConfig) error {
	client, err := d.NewClient()
	if err != nil {
		return err
	}
	return client.DeleteWebhook(webhookInfo)
}

func (d *Repo) setDefaultBranch() {
	if d.Branch == "" {
		switch d.RepoType {
		case "github":
			d.Branch = "main"
		case "gitlab":
			d.Branch = "master"
		}
	}
}

func NewRepoFromURL(repoURL, branch string) (*Repo, error) {
	repo := &Repo{
		Branch: branch,
	}
	u, err := url.Parse(repoURL)
	if err != nil {
		return nil, err
	}
	//config repo type
	if strings.Contains(u.Host, "github") {
		repo.RepoType = "github"
	} else if strings.Contains(u.Host, "gitlab.com") {
		repo.RepoType = "gitlab"
	} else {
		repo.RepoType = "gitlab"
		repo.BaseURL = fmt.Sprintf("%s://%s", u.Scheme, u.Host)
	}
	path := u.Path
	// config repo owner org and
	pathPart := strings.Split(strings.TrimPrefix(path, "/"), "/")
	if len(pathPart) != 2 {
		return nil, fmt.Errorf("git repo path is not valid")
	}
	repo.Owner = pathPart[0]
	repo.Repo = pathPart[1]
	repo.setDefaultBranch()
	return repo, nil
}

package common

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/devstream-io/devstream/pkg/util/file"
	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/scm"
	"github.com/devstream-io/devstream/pkg/util/scm/git"
	"github.com/devstream-io/devstream/pkg/util/scm/github"
	"github.com/devstream-io/devstream/pkg/util/scm/gitlab"
)

const (
	defaultMainBranch = "main"
)

// Repo is the repo info of github or gitlab
type Repo struct {
	Owner             string `validate:"required_without=Org" mapstructure:"owner"`
	Org               string `validate:"required_without=Owner" mapstructure:"org"`
	Repo              string `validate:"required" mapstructure:"repo"`
	Branch            string `mapstructure:"branch"`
	PathWithNamespace string
	RepoType          string `validate:"oneof=gitlab github" mapstructure:"repoType"`
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
		Type:       d.RepoType,
	}
}

// NewRepoFromURL build repo struct from scm url
func NewRepoFromURL(repoType, apiURL, cloneURL, branch string) (*Repo, error) {
	repo := &Repo{
		Branch: branch,
	}

	if scm.IsGithubRepo(repoType, cloneURL) {
		repo.RepoType = "github"
	} else {
		repo.RepoType = "gitlab"
		// extract gitlab baseURL from url string
		if apiURL == "" {
			apiURL = cloneURL
		}
		gitlabBaseURL, err := gitlab.ExtractBaseURLfromRaw(apiURL)
		if err != nil {
			return nil, fmt.Errorf("gitlab repo extract baseURL failed: %w", err)
		}
		repo.BaseURL = gitlabBaseURL
	}

	if err := repo.updateRepoPathByCloneURL(cloneURL); err != nil {
		return nil, fmt.Errorf("git extract repo info failed: %w", err)
	}
	repo.Branch = repo.getBranch()
	return repo, nil
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

func (d *Repo) getRepoDownloadURL() string {
	repoInfo := d.BuildRepoInfo()
	latestCodeZipfileDownloadURL := fmt.Sprintf(
		github.DefaultLatestCodeZipfileDownloadUrlFormat, repoInfo.GetRepoOwner(), repoInfo.Repo, repoInfo.Branch,
	)
	log.Debugf("LatestCodeZipfileDownloadUrl: %s.", latestCodeZipfileDownloadURL)
	return latestCodeZipfileDownloadURL
}

func (d *Repo) updateRepoPathByCloneURL(cloneURL string) error {
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
func (d *Repo) BuildURL() string {
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

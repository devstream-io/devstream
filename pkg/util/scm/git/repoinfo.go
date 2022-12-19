package git

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/mapz"
)

type ScmURL string

const (
	privateSSHKeyEnv = "REPO_SSH_PRIVATEKEY"
)

type RepoInfo struct {
	// repo detail fields
	Owner    string `yaml:"owner" mapstructure:"owner,omitempty"`
	Org      string `yaml:"org" mapstructure:"org,omitempty"`
	Repo     string `yaml:"name" mapstructure:"name,omitempty"`
	Branch   string `yaml:"branch" mapstructure:"branch,omitempty"`
	RepoType string `yaml:"scmType" mapstructure:"scmType,omitempty"`

	// url fields
	APIURL   string `yaml:"apiURL" mapstructure:"apiURL,omitempty"`
	CloneURL ScmURL `yaml:"url" mapstructure:"url,omitempty"`

	// used for gitlab
	Namespace     string `mapstructure:"nameSpace,omitempty"`
	Visibility    string `mapstructure:"visibility,omitempty"`
	BaseURL       string `yaml:"baseURL" mapstructure:"baseURL,omitempty"`
	SSHPrivateKey string `yaml:"sshPrivateKey" mapstructure:"sshPrivateKey,omitempty"`

	// used for GitHub
	WorkPath string `mapstructure:"workPath,omitempty"`
	NeedAuth bool   `mapstructure:"needAuth,omitempty"`
}

func (r *RepoInfo) SetDefault() error {
	// set repoInfo default field
	if r.checkNeedUpdateFromURL() {
		// opts has only CloneURL is configured, update other fields from url
		if err := r.updateFieldsFromURLField(); err != nil {
			return err
		}
		// else build ScmURL from RepoInfo other fields
	} else {
		r.CloneURL = r.buildScmURL()
		r.Branch = r.getBranchWithDefault()
	}
	return r.checkValid()
}

func (r *RepoInfo) GetRepoOwner() string {
	// if org or owner is not empty, return org/owner
	if r.Org != "" {
		return r.Org
	}
	if r.Owner != "" {
		return r.Owner
	}
	// else return owner extract from url
	if r.CloneURL != "" {
		owner, _, err := r.CloneURL.extractRepoOwnerAndName()
		if err != nil {
			log.Warnf("git GetRepoName failed %s", err)
		}
		return owner
	}
	return ""
}

func (r *RepoInfo) GetRepoPath() string {
	return fmt.Sprintf("%s/%s", r.GetRepoOwner(), r.GetRepoName())
}

func (r *RepoInfo) GetRepoName() string {
	var repoName string
	var err error
	if r.Repo != "" {
		repoName = r.Repo
	} else if r.CloneURL != "" {
		_, repoName, err = r.CloneURL.extractRepoOwnerAndName()
		if err != nil {
			log.Warnf("git GetRepoName failed %s", err)
		}
	}
	return repoName
}

func (r *RepoInfo) GetCloneURL() string {
	var cloneURL string
	if r.CloneURL != "" {
		cloneURL = string(r.CloneURL.addGithubURLScheme())
	} else {
		cloneURL = string(r.buildScmURL())
	}
	return cloneURL
}

func (r *RepoInfo) Encode() map[string]any {
	m, err := mapz.DecodeStructToMap(r)
	if err != nil {
		log.Errorf("gitRepo [%+v] decode to map failed: %+v", r, err)
	}
	return m
}

// IsGithubRepo return ture if repo is github
func (r *RepoInfo) IsGithubRepo() bool {
	return r.RepoType == "github" || strings.Contains(string(r.CloneURL), "github")
}

func (r *RepoInfo) getBranchWithDefault() string {
	branch := r.Branch
	if branch != "" {
		return branch
	}
	if r.IsGithubRepo() {
		branch = "main"
	} else {
		branch = "master"
	}
	return branch
}

// BuildURL return url build from repo struct
func (r *RepoInfo) buildScmURL() ScmURL {
	switch r.RepoType {
	case "github":
		return ScmURL(fmt.Sprintf("https://github.com/%s/%s", r.GetRepoOwner(), r.Repo))
	case "gitlab":
		var gitlabURL string
		if r.BaseURL != "" {
			gitlabURL = r.BaseURL
		} else {
			gitlabURL = "https://gitlab.com"
		}
		return ScmURL(fmt.Sprintf("%s/%s/%s.git", gitlabURL, r.GetRepoOwner(), r.Repo))
	default:
		log.Warnf("git repo buildScmURL get invalid repo type: %s", r.RepoType)
		return ""
	}
}

func (r *RepoInfo) checkNeedUpdateFromURL() bool {
	return r.CloneURL != "" && (r.RepoType == "" || r.Repo == "")
}

func (r *RepoInfo) updateFieldsFromURLField() error {
	// 1. config basic info for different repo type
	if r.IsGithubRepo() {
		r.RepoType = "github"
	} else {
		r.RepoType = "gitlab"
		// extract gitlab baseURL from url string
		apiURL := r.APIURL
		if apiURL == "" {
			apiURL = string(r.CloneURL)
		}
		gitlabBaseURL, err := extractBaseURLfromRaw(apiURL)
		if err != nil {
			return err
		}
		r.BaseURL = gitlabBaseURL
	}

	// 2. get repoOwner and Name from CloneURL field
	repoOwner, repoName, err := r.CloneURL.extractRepoOwnerAndName()
	if err != nil {
		return err
	}
	if r.Org == "" && r.Owner == "" {
		r.Owner = repoOwner
	}
	r.Repo = repoName
	// 3. if scm.branch is not configured, just use repo's default branch
	r.Branch = r.getBranchWithDefault()
	if r.SSHPrivateKey == "" {
		r.SSHPrivateKey = os.Getenv(privateSSHKeyEnv)
	}
	return nil
}

func (r *RepoInfo) checkValid() error {
	// basic check
	if r.Org != "" && r.Owner != "" {
		return fmt.Errorf("git org and owner can't be configured at the same time")
	}
	if r.RepoType != "github" && r.RepoType != "gitlab" {
		return fmt.Errorf("git scmType only support gitlab and github, current scm type is [%s]", r.RepoType)
	}
	if r.Repo == "" {
		return fmt.Errorf("git name field must be configured")
	}
	// check token is configured
	if r.NeedAuth {
		switch r.RepoType {
		case "gitlab":
			if os.Getenv("GITLAB_TOKEN") == "" {
				return fmt.Errorf("gitlab repo should set env GITLAB_TOKEN")
			}
		case "github":
			if os.Getenv("GITHUB_TOKEN") == "" {
				return fmt.Errorf("github repo should set env GITHUB_TOKEN")
			}
		}
	}
	return nil
}

// extractRepoOwnerAndName will get repoOwner and repoName from ScmURL
func (u ScmURL) extractRepoOwnerAndName() (string, string, error) {
	var paths string
	ScmURLStr := string(u.addGithubURLScheme())
	c, err := url.ParseRequestURI(ScmURLStr)
	if err != nil {
		if strings.Contains(ScmURLStr, "git@") {
			gitSSHLastIndex := strings.LastIndex(ScmURLStr, ":")
			if gitSSHLastIndex == -1 {
				err = fmt.Errorf("git url ssh repo not valid")
				return "", "", err
			}
			paths = strings.Trim(ScmURLStr[gitSSHLastIndex:], ":")
		} else {
			err = fmt.Errorf("git url repo transport not support for now")
			return "", "", err
		}
	} else {
		paths = c.Path
	}
	projectPaths := strings.Split(strings.Trim(paths, "/"), "/")
	if len(projectPaths) != 2 {
		err = fmt.Errorf("git url repo path is not valid")
		return "", "", err
	}
	repoOwner := projectPaths[0]
	repoName := strings.TrimSuffix(projectPaths[1], ".git")
	return repoOwner, repoName, nil
}

// addGithubURLScheme will add "https://"  in github url config if it doesn't contain schme
func (u ScmURL) addGithubURLScheme() ScmURL {
	cloneURL := string(u)
	if !strings.Contains(cloneURL, "git@") && !strings.HasPrefix(cloneURL, "http") {
		return ScmURL(fmt.Sprintf("https://%s", cloneURL))
	}
	return u
}

func extractBaseURLfromRaw(repoURL string) (string, error) {
	u, err := url.ParseRequestURI(repoURL)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s://%s", u.Scheme, u.Host), nil
}

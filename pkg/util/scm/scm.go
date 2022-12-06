package scm

import (
	"fmt"
	"os"
	"strings"

	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/mapz"
	"github.com/devstream-io/devstream/pkg/util/scm/git"
	"github.com/devstream-io/devstream/pkg/util/scm/gitlab"
)

const (
	privateSSHKeyEnv = "REPO_SSH_PRIVATEKEY"
)

// SCMInfo can use BasicRepo and URLRepo to get scm info
type SCMInfo struct {
	// SCM URL related info
	CloneURL string `mapstructure:"url,omitempty" yaml:"url"`
	APIURL   string `mapstructure:"apiURL,omitempty" yaml:"apiURL"`
	// SCM Basic related info
	Owner string `yaml:"owner" mapstructure:"owner,omitempty"`
	Org   string `yaml:"org" mapstructure:"org,omitempty"`
	Name  string `yaml:"name" mapstructure:"name,omitempty"`
	Type  string `mapstructure:"scmType,omitempty" yaml:"scmType"`
	// common fields
	Branch        string `mapstructure:"branch,omitempty" yaml:"branch"`
	SSHPrivateKey string `mapstructure:"sshPrivateKey,omitempty"`
}

// BuildRepoInfo will return RepoInfo from SCMInfo
func (s *SCMInfo) BuildRepoInfo() (*git.RepoInfo, error) {
	var repo *git.RepoInfo
	var err error
	if s.CloneURL != "" {
		// get repo from url
		repo, err = s.getRepoInfoFromURL()
	} else {
		// get repo from repo basic fields
		repo = s.getRepoInfoFromBasic()
	}
	if err != nil {
		return nil, err
	}
	// config default branch
	repo.Branch = repo.GetBranchWithDefault()
	return repo, nil
}

func (s *SCMInfo) getRepoInfoFromURL() (*git.RepoInfo, error) {
	repo := &git.RepoInfo{
		Branch:        s.Branch,
		SSHPrivateKey: s.SSHPrivateKey,
		CloneURL:      s.CloneURL,
	}

	if isGithubRepo(s.Type, s.CloneURL) {
		repo.RepoType = "github"
		s.CloneURL = formatGithubCloneURL(s.CloneURL)
	} else {
		repo.RepoType = "gitlab"
		// extract gitlab baseURL from url string
		apiURL := s.APIURL
		if apiURL == "" {
			apiURL = s.CloneURL
		}
		gitlabBaseURL, err := gitlab.ExtractBaseURLfromRaw(apiURL)
		if err != nil {
			return nil, fmt.Errorf("gitlab repo extract baseURL failed: %w", err)
		}
		repo.BaseURL = gitlabBaseURL
	}

	if err := repo.UpdateRepoPathByCloneURL(s.CloneURL); err != nil {
		return nil, fmt.Errorf("git extract repo info failed: %w", err)
	}
	// if scm.branch is not configured, just use repo's default branch
	repo.Branch = repo.GetBranchWithDefault()
	if repo.SSHPrivateKey == "" {
		repo.SSHPrivateKey = os.Getenv(privateSSHKeyEnv)
	}
	repo.NeedAuth = true
	// update org info
	if s.Org != "" {
		repo.Org = s.Org
		repo.Owner = ""
	}
	return repo, nil

}

func (s *SCMInfo) getRepoInfoFromBasic() *git.RepoInfo {
	repoInfo := &git.RepoInfo{
		Owner:         s.Owner,
		Org:           s.Org,
		Repo:          s.Name,
		Branch:        s.Branch,
		RepoType:      s.Type,
		SSHPrivateKey: s.SSHPrivateKey,
		NeedAuth:      true,
	}
	repoInfo.CloneURL = repoInfo.BuildScmURL()
	return repoInfo
}

func (s *SCMInfo) Encode() map[string]any {
	m, err := mapz.DecodeStructToMap(s)
	if err != nil {
		log.Errorf("scmInfo [%+v] decode to map failed: %+v", s, err)
	}
	return m
}

func (s *SCMInfo) NewClientWithAuthFromScm() (ClientOperation, error) {
	repo, err := s.BuildRepoInfo()
	if err != nil {
		return nil, err
	}
	return NewClientWithAuth(repo)
}

func formatGithubCloneURL(cloneURL string) string {
	if !strings.Contains(cloneURL, "git@") && !strings.HasPrefix(cloneURL, "http") {
		return fmt.Sprintf("https://%s", cloneURL)
	}
	return cloneURL
}

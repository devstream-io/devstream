package scm

import (
	"fmt"
	"os"

	"github.com/devstream-io/devstream/pkg/util/scm/git"
	"github.com/devstream-io/devstream/pkg/util/scm/gitlab"
)

const (
	privateSSHKeyEnv = "REPO_SSH_PRIVATEKEY"
)

type SCMInfo struct {
	CloneURL string `mapstructure:"cloneURL" validate:"required"`
	APIURL   string `mapstructure:"apiURL"`
	Branch   string `mapstructure:"branch"`
	Type     string `mapstructure:"type"`

	// used in package
	SSHPrivateKey string `mapstructure:"sshPrivateKey"`
}

func (s *SCMInfo) BuildRepoInfo() (*git.RepoInfo, error) {
	repo := &git.RepoInfo{
		Branch:        s.Branch,
		SSHPrivateKey: s.SSHPrivateKey,
		CloneURL:      s.CloneURL,
	}

	if isGithubRepo(s.Type, s.CloneURL) {
		repo.RepoType = "github"
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
	return repo, nil
}

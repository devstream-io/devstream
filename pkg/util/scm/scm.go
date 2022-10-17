package scm

import (
	"fmt"

	"github.com/devstream-io/devstream/pkg/util/scm/gitlab"
)

type SCMInfo struct {
	CloneURL string `mapstructure:"cloneURL" validate:"required"`
	APIURL   string `mapstructure:"apiURL"`
	Branch   string `mapstructure:"branch"`
	Type     string `mapstructure:"type"`

	// used in package
	SSHprivateKey string `mapstructure:"sshPrivateKey"`
}

func (s *SCMInfo) NewRepo() (*Repo, error) {
	repo := &Repo{
		Branch: s.Branch,
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
	repo.Branch = repo.getBranch()
	if s.Branch == "" {
		s.Branch = repo.Branch
	}
	return repo, nil
}

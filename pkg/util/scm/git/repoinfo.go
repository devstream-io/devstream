package git

import (
	"fmt"
)

type RepoInfo struct {
	Repo     string
	Branch   string
	Owner    string
	Org      string
	Type     string
	CloneURL string

	// used for gitlab
	Visibility string
	Namespace  string
	BaseURL    string

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

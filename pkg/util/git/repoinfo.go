package git

import (
	"fmt"

	
)

type RepoInfo struct {
	Repo   string
	Branch string
	Owner  string
	Org    string

	// used for gitlab
	Visibility string
	Namespace  string
	BaseURL    string

	// used for github
	WorkPath string
	NeedAuth bool


	// used for gitea
	Description string
	Labels string

}



func (r *RepoInfo) GetRepoNameWithBranch() string {
	return fmt.Sprintf("%s-%s", r.Repo, r.Branch)
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

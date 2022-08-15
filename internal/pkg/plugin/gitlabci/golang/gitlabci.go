package golang

import (
	"errors"
	"strings"

	"github.com/devstream-io/devstream/pkg/util/git"
	"github.com/devstream-io/devstream/pkg/util/gitlab"
)

const (
	ciFileName    string = ".gitlab-ci.yml"
	commitMessage string = "managed by DevStream"
)

type Options struct {
	PathWithNamespace string `validate:"required"`
	Branch            string `validate:"required"`
	BaseURL           string `validate:"omitempty,url"`
}

func buildState(opts *Options) map[string]interface{} {
	return map[string]interface{}{
		"pathWithNamespace": opts.PathWithNamespace,
		"branch":            opts.Branch,
	}
}

func (opts *Options) newGitlabClient() (*gitlab.Client, error) {
	pathSplit := strings.Split(opts.PathWithNamespace, "/")
	if len(pathSplit) != 2 {
		return nil, errors.New("gitlabci generic not valid PathWithNamespace params")
	}
	repoInfo := &git.RepoInfo{
		Owner:   pathSplit[0],
		Repo:    pathSplit[1],
		Branch:  opts.Branch,
		BaseURL: opts.BaseURL,
	}
	return gitlab.NewClient(repoInfo)
}

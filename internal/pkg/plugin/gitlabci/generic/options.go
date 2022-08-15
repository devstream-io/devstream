package generic

import (
	"errors"
	"strings"

	"github.com/devstream-io/devstream/pkg/util/git"
	"github.com/devstream-io/devstream/pkg/util/gitlab"
)

// Options is the struct for configurations of the gitlabci-generic plugin.
type Options struct {
	PathWithNamespace string `validate:"required"`
	Branch            string `validate:"required"`
	TemplateURL       string `validate:"required"`
	BaseURL           string `validate:"omitempty,url"`
	TemplateVariables map[string]interface{}
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

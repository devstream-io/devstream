package reposcaffolding

import (
	"github.com/devstream-io/devstream/pkg/util/scm/git"
)

type options struct {
	SourceRepo      *git.RepoInfo `validate:"required" mapstructure:"sourceRepo"`
	DestinationRepo *git.RepoInfo `validate:"required" mapstructure:"destinationRepo"`
	Vars            map[string]interface{}
}

func (opts *options) renderTplConfig() map[string]interface{} {
	// default render value from repo
	renderConfig := map[string]any{
		"AppName": opts.DestinationRepo.GetRepoName(),
		"Repo": map[string]string{
			"Name":  opts.DestinationRepo.GetRepoName(),
			"Owner": opts.DestinationRepo.GetRepoOwner(),
		},
	}
	for k, v := range opts.Vars {
		renderConfig[k] = v
	}
	return renderConfig
}

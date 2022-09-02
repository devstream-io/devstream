package golang

import (
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/ci"
	"github.com/devstream-io/devstream/pkg/util/scm/gitlab"
)

func setCIContent(options plugininstaller.RawOptions) (plugininstaller.RawOptions, error) {
	opts, err := ci.NewOptions(options)
	if err != nil {
		return nil, err
	}
	gitlabClient, err := gitlab.NewClient(opts.ProjectRepo.BuildRepoInfo())
	if err != nil {
		return nil, err
	}

	ciContent, err := gitlabClient.GetGitLabCIGolangTemplate()
	if err != nil {
		return nil, err
	}
	ciConfig := opts.CIConfig
	if ciConfig == nil {
		ciConfig = &ci.CIConfig{}
	}
	ciConfig.Content = ciContent
	opts.CIConfig = ciConfig
	return opts.Encode()
}

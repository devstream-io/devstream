package golang

import (
	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/ci"
	"github.com/devstream-io/devstream/pkg/util/scm/gitlab"
	"github.com/devstream-io/devstream/pkg/util/types"
)

func setCIContent(options configmanager.RawOptions) (configmanager.RawOptions, error) {
	opts, err := ci.NewOptions(options)
	if err != nil {
		return nil, err
	}
	gitlabClient, err := gitlab.NewClient(opts.ProjectRepo)
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
	ciConfig.SetContent(ciContent)
	opts.CIConfig = ciConfig
	return types.EncodeStruct(opts)
}

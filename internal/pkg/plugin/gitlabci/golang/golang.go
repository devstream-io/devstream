package golang

import (
	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer/ci/cifile"
	"github.com/devstream-io/devstream/pkg/util/scm/gitlab"
	"github.com/devstream-io/devstream/pkg/util/types"
)

func setCIContent(options configmanager.RawOptions) (configmanager.RawOptions, error) {
	opts, err := cifile.NewOptions(options)
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
	CIFileConfig := opts.CIFileConfig
	if CIFileConfig == nil {
		CIFileConfig = &cifile.CIFileConfig{}
	}
	CIFileConfig.SetContent(ciContent)
	opts.CIFileConfig = CIFileConfig
	return types.EncodeStruct(opts)
}

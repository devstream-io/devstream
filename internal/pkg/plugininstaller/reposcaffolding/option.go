package reposcaffolding

import (
	"github.com/mitchellh/mapstructure"

	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/pkg/util/file"
	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/scm"
	"github.com/devstream-io/devstream/pkg/util/scm/git"
)

type Options struct {
	SourceRepo      *scm.Repo `validate:"required" mapstructure:"sourceRepo"`
	DestinationRepo *scm.Repo `validate:"required" mapstructure:"destinationRepo"`
	Vars            map[string]interface{}
}

func NewOptions(options plugininstaller.RawOptions) (*Options, error) {
	var opts Options
	if err := mapstructure.Decode(options, &opts); err != nil {
		return nil, err
	}
	return &opts, nil
}

func (opts *Options) renderTplConfig() map[string]interface{} {
	renderConfig := opts.DestinationRepo.BuildRepoRenderConfig()
	for k, v := range opts.Vars {
		renderConfig[k] = v
	}
	return renderConfig
}

// downloadAndRenderScmRepo will download repo from source repo and render it locally
func (opts *Options) downloadAndRenderScmRepo(scmClient scm.ClientOperation) (git.GitFileContentMap, error) {
	// 1. download zip file and unzip this file then render folders
	zipFilesDir, err := scmClient.DownloadRepo()
	if err != nil {
		log.Debugf("reposcaffolding process files error: %s", err)
		return nil, err
	}
	appName := opts.DestinationRepo.Repo
	return file.GetFileMapByWalkDir(
		zipFilesDir, filterGitFiles,
		getRepoFileNameFunc(appName, opts.SourceRepo.GetRepoNameWithBranch()),
		processRepoFileFunc(appName, opts.renderTplConfig()),
	)
}

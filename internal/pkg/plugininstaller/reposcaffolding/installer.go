package reposcaffolding

import (
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
)

// InstallRepo will install repo by opts config
func InstallRepo(options plugininstaller.RawOptions) error {
	opts, err := NewOptions(options)
	if err != nil {
		return err
	}

	// 1. Create and render repo get from given url
	gitMap, err := opts.SourceRepo.CreateAndRenderLocalRepo(
		opts.DestinationRepo.Repo, opts.renderTplConfig(),
	)
	if err != nil {
		return err
	}

	// 2. Push local repo to remote
	return opts.DestinationRepo.CreateAndPush(gitMap)
}

// DeleteRepo will delete repo by options
func DeleteRepo(options plugininstaller.RawOptions) error {
	var err error
	opts, err := NewOptions(options)
	if err != nil {
		return err
	}

	return opts.DestinationRepo.Delete()
}

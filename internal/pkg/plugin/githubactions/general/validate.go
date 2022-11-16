package general

import (
	"fmt"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/plugin/plugininstaller/ci"
)

// validate validates the options provided by the core.
func validate(options configmanager.RawOptions) (configmanager.RawOptions, error) {
	opts, err := ci.NewCIOptions(options)
	if err != nil {
		return nil, err
	}
	// check basic ci error
	_, err = ci.Validate(options)
	if err != nil {
		return nil, err
	}
	// check repo is valid
	if opts.ProjectRepo.RepoType != "github" {
		return nil, fmt.Errorf("github action don't support other repo")
	}
	return options, nil
}

package reposcaffolding

import (
	"fmt"
	"path/filepath"

	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/pkg/util/validator"
)

// Validate validates the options provided by the core.
func Validate(options plugininstaller.RawOptions) (plugininstaller.RawOptions, error) {
	opts, err := NewOptions(options)
	if err != nil {
		return nil, err
	}
	err = validator.StructAllError(opts)
	if err != nil {
		return nil, err
	}
	return options, nil
}

// SetDefaultTemplateRepo set default value for repo
func SetDefaultTemplateRepo(options plugininstaller.RawOptions) (plugininstaller.RawOptions, error) {
	opts, err := NewOptions(options)
	if err != nil {
		return nil, err
	}
	// set dstRepo default value
	dstRepo := opts.DestinationRepo
	// set PathWithNamespace for GitLab. GitHub won't need to use this
	// opts.PathWithNamespace = fmt.Sprintf("%s/%s", opts.Owner, opts.Repo)
	if dstRepo.Org != "" {
		dstRepo.PathWithNamespace = filepath.Clean(fmt.Sprintf("%s/%s", dstRepo.Org, dstRepo.Repo))
	} else {
		dstRepo.PathWithNamespace = filepath.Clean(fmt.Sprintf("%s/%s", dstRepo.Owner, dstRepo.Repo))
	}
	return opts.Encode()
}

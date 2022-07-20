package github

import (
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/pkg/util/validator"
)

func Validate(options plugininstaller.RawOptions) (plugininstaller.RawOptions, error) {
	opts, err := NewGithubActionOptions(options)
	if err != nil {
		return nil, err
	}
	return options, validator.StructAllError(opts)
}

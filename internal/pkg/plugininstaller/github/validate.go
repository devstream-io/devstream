package github

import (
	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/pkg/util/validator"
)

func Validate(options configmanager.RawOptions) (configmanager.RawOptions, error) {
	opts, err := NewGithubActionOptions(options)
	if err != nil {
		return nil, err
	}
	return options, validator.StructAllError(opts)
}

package gitlabcedocker

import (
	"fmt"
	"path"

	"github.com/devstream-io/devstream/pkg/util/validator"
)

// validate validates the options provided by the core.
func validate(options *Options) []error {
	errs := validator.Struct(options)

	// volume directory must be absolute path
	if !path.IsAbs(options.GitLabHome) {
		errs = append(errs, fmt.Errorf("GitLabHome must be an absolute path"))
	}

	return errs
}

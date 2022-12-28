package gitlabcedocker

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"

	"github.com/mitchellh/mapstructure"

	"github.com/devstream-io/devstream/pkg/util/validator"
)

func validateAndDefault(options configmanager.RawOptions) (*Options, error) {
	var opts *Options
	if err := mapstructure.Decode(options, &opts); err != nil {
		return nil, err
	}

	opts.Defaults()

	// validate
	errs := validator.CheckStructError(opts)
	// volume directory must be absolute path
	if !filepath.IsAbs(opts.GitLabHome) {
		errs = append(errs, fmt.Errorf("field gitLabHome must be an absolute path"))
	}
	if len(errs) != 0 {
		return nil, errs.Combine()
	}

	if err := os.MkdirAll(opts.GitLabHome, 0755); err != nil {
		return nil, err
	}

	opts.setGitLabURL()

	return opts, nil
}

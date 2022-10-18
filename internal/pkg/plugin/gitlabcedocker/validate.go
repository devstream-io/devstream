package gitlabcedocker

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"

	"github.com/mitchellh/mapstructure"

	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/validator"
)

func validateAndDefault(options configmanager.RawOptions) (*Options, error) {
	var opts *Options
	if err := mapstructure.Decode(options, &opts); err != nil {
		return nil, err
	}

	opts.Defaults()

	// validate
	errs := validator.Struct(opts)
	// volume directory must be absolute path
	if !filepath.IsAbs(opts.GitLabHome) {
		errs = append(errs, fmt.Errorf("GitLabHome must be an absolute path"))
	}
	if len(errs) > 0 {
		for _, e := range errs {
			log.Errorf("Options error: %s.", e)
		}
		return nil, fmt.Errorf("opts are illegal")
	}

	if err := os.MkdirAll(opts.GitLabHome, 0755); err != nil {
		return nil, err
	}

	opts.setGitLabURL()

	return opts, nil
}

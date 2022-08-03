package gitlabcedocker

import (
	"fmt"
	"path/filepath"

	"github.com/mitchellh/mapstructure"

	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/validator"
)

func preHandleOptions(options map[string]interface{}) (*Options, error) {
	var opts *Options
	if err := mapstructure.Decode(options, &opts); err != nil {
		return nil, err
	}

	defaults(opts)

	if err := validate(opts); err != nil {
		return nil, err
	}

	return opts, nil
}

func defaults(opts *Options) {
	if opts.ImageTag == "" {
		opts.ImageTag = defaultImageTag
	}
}

// validate validates the options provided by the core.
func validate(options *Options) error {
	errs := validator.Struct(options)
	// volume directory must be absolute path
	if !filepath.IsAbs(options.GitLabHome) {
		errs = append(errs, fmt.Errorf("GitLabHome must be an absolute path"))
	}

	if len(errs) > 0 {
		for _, e := range errs {
			log.Errorf("Options error: %s.", e)
		}
		return fmt.Errorf("opts are illegal")
	}

	return nil
}

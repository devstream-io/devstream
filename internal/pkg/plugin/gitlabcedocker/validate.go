package gitlabcedocker

import (
	"fmt"
	"path/filepath"

	"github.com/mitchellh/mapstructure"

	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/validator"
)

func validateAndDefault(options map[string]interface{}) (*Options, error) {
	var opts *Options
	if err := mapstructure.Decode(options, &opts); err != nil {
		return nil, err
	}

	// validate
	errs := validator.Struct(options)
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

	// default
	if opts.ImageTag == "" {
		opts.ImageTag = defaultGitlabImageTag
	}

	opts.setGitLabURL()

	return opts, nil
}

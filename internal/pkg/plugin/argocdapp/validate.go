package argocdapp

import (
	"fmt"
	"strings"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/pkg/util/mapz"
	"github.com/devstream-io/devstream/pkg/util/validator"
)

// validate validates the options provided by the core.
func validate(options configmanager.RawOptions) (configmanager.RawOptions, error) {
	opts, err := newOptions(options)
	if err != nil {
		return nil, err
	}
	if err = validator.StructAllError(opts); err != nil {
		return nil, err
	}
	return options, nil
}

func setDefault(options configmanager.RawOptions) (configmanager.RawOptions, error) {
	const defaultImageTag = "0.0.1"
	opts, err := newOptions(options)
	if err != nil {
		return nil, err
	}
	// imageRepo initalTag use latest by default
	if opts.ImageRepo != nil {
		if opts.ImageRepo.InitalTag == "" {
			opts.ImageRepo.InitalTag = defaultImageTag
		}
		if opts.ImageRepo.URL != "" && strings.HasSuffix(opts.ImageRepo.URL, "") {
			opts.ImageRepo.URL = fmt.Sprintf("%s/", opts.ImageRepo.URL)
		}
	}

	// set ci file config
	return mapz.DecodeStructToMap(opts)
}

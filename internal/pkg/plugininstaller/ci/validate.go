package ci

import (
	"errors"
	"net/url"
	"os"

	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/pkg/util/validator"
)

// Validate validates the options provided by the dtm-core.
func Validate(options plugininstaller.RawOptions) (plugininstaller.RawOptions, error) {
	opts, err := NewOptions(options)
	if err != nil {
		return nil, err
	}
	fieldErr := validator.StructAllError(opts)
	if fieldErr != nil {
		return nil, fieldErr
	}
	// check CI config
	config := opts.CIConfig

	if config.RemoteURL != "" {
		_, err := url.ParseRequestURI(config.RemoteURL)
		if err != nil {
			return nil, err
		}
	} else if config.LocalPath != "" {
		_, err := os.Stat(config.LocalPath)
		if err != nil {
			return nil, err
		}
	} else if config.Content == "" {
		return nil, errors.New("ci.locaPath, ci.remoteURL, ci.content can't all be empty at the same time")
	}
	return options, nil
}

// SetDefaultConfig will update options empty values base on import options
func SetDefaultConfig(defaultConfig *Options) plugininstaller.MutableOperation {
	return func(options plugininstaller.RawOptions) (plugininstaller.RawOptions, error) {
		opts, err := NewOptions(options)
		if err != nil {
			return nil, err
		}
		opts.FillDefaultValue(defaultConfig)
		return opts.Encode()
	}
}

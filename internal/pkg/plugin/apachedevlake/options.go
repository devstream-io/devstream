package apachedevlake

import (
	"github.com/mitchellh/mapstructure"

	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
)

// Options is the struct for configurations of the apache-devlake plugin.
type Options struct {
	// TODO(dtm): Add your params here.
	Foo string
}

// NewOptions create options by raw options
func NewOptions(options plugininstaller.RawOptions) (Options, error) {
	var opts Options
	if err := mapstructure.Decode(options, &opts); err != nil {
		return opts, err
	}
	return opts, nil
}

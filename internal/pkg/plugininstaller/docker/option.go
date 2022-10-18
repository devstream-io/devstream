package docker

import (
	"fmt"

	"github.com/mitchellh/mapstructure"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/pkg/util/docker"
	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/validator"
)

type Options struct {
	ImageName         string `validate:"required"`
	ImageTag          string `validate:"required"`
	ContainerName     string `validate:"required"`
	RmDataAfterDelete *bool

	RunParams     []string
	Hostname      string
	PortPublishes []docker.PortPublish
	Volumes       docker.Volumes
	RestartAlways bool
}

// NewOptions create options by raw options
func NewOptions(options configmanager.RawOptions) (Options, error) {
	var opts Options
	if err := mapstructure.Decode(options, &opts); err != nil {
		return opts, err
	}
	return opts, nil
}

func (opts *Options) GetImageNameWithTag() string {
	return fmt.Sprintf("%s:%s", opts.ImageName, opts.ImageTag)
}

func (opts *Options) GetRunOpts() *docker.RunOptions {
	return &docker.RunOptions{
		ImageName:     opts.ImageName,
		ImageTag:      opts.ImageTag,
		ContainerName: opts.ContainerName,
		Hostname:      opts.Hostname,
		PortPublishes: opts.PortPublishes,
		Volumes:       opts.Volumes,
		RestartAlways: opts.RestartAlways,
	}
}
func Validate(options configmanager.RawOptions) (configmanager.RawOptions, error) {
	opts, err := NewOptions(options)
	if err != nil {
		return nil, err
	}
	log.Debugf("Options: %v.", opts)
	errs := validator.Struct(opts)
	if len(errs) > 0 {
		for _, e := range errs {
			log.Errorf("Options error: %s.", e)
		}
		return nil, fmt.Errorf("opts are illegal")
	}
	return options, nil
}

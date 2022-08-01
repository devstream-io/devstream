package docker

import (
	"fmt"
	"strings"

	"go.uber.org/multierr"
)

// RunOptions is used to pass options to ContainerRunWithOptions
type (
	RunOptions struct {
		ImageName     string
		ImageTag      string
		Hostname      string
		ContainerName string
		PortPublishes []PortPublish
		Volumes       Volumes
		RestartAlways bool
		RunParams     []string
	}

	Volume struct {
		HostPath      string
		ContainerPath string
	}
	Volumes []Volume

	PortPublish struct {
		HostPort      uint
		ContainerPort uint
	}
)

func (opts *RunOptions) Validate() error {
	var errs []error
	if strings.TrimSpace(opts.ImageName) == "" {
		errs = append(errs, fmt.Errorf("image name is required"))
	}
	if strings.TrimSpace(opts.ImageTag) == "" {
		errs = append(errs, fmt.Errorf("image tag is required"))
	}
	if strings.TrimSpace(opts.ContainerName) == "" {
		errs = append(errs, fmt.Errorf("container name is required"))
	}
	for _, volume := range opts.Volumes {
		if volume.HostPath == "" {
			errs = append(errs, fmt.Errorf("HostPath can not be empty"))
		}
		if volume.ContainerPath == "" {
			errs = append(errs, fmt.Errorf("ContainerPath can not be empty"))
		}
	}

	return multierr.Combine(errs...)
}

func CombineImageNameAndTag(imageName, tag string) string {
	return imageName + ":" + tag
}

func (volumes Volumes) ExtractHostPaths() []string {
	hostPaths := make([]string, len(volumes))
	for i, volume := range volumes {
		hostPaths[i] = volume.HostPath
	}
	return hostPaths
}

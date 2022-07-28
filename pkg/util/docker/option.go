package docker

import (
	"fmt"
	"strings"
)

// RunOptions is used to pass options to ContainerRunWithOptions
type (
	RunOptions struct {
		ImageName     string
		ImageTag      string
		Hostname      string
		ContainerName string
		PortPublishes []PortPublish
		Volumes       []Volume
		RestartAlways bool
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

	return CombineErrs(errs)
}

func CombineImageNameAndTag(imageName, tag string) string {
	return imageName + ":" + tag
}

func CombineErrs(errs []error) error {
	if len(errs) == 0 {
		return nil
	}

	errsString := make([]string, len(errs))
	for _, err := range errs {
		errsString = append(errsString, err.Error())
	}

	return fmt.Errorf(strings.Join(errsString, ";"))
}

func (volumes Volumes) ExtractHostPaths() []string {
	hostPaths := make([]string, len(volumes))
	for i, volume := range volumes {
		hostPaths[i] = volume.HostPath
	}
	return hostPaths
}

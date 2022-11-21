package docker

import (
	"fmt"

	"go.uber.org/multierr"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/pkg/util/docker"
	"github.com/devstream-io/devstream/pkg/util/docker/dockersh"
	"github.com/devstream-io/devstream/pkg/util/log"
)

var op docker.Operator = &dockersh.ShellOperator{}

// Install runs the docker container
func Install(options configmanager.RawOptions) error {
	opts, err := NewOptions(options)
	if err != nil {
		return err
	}

	// 1. pull the image if it not exists
	if !op.ImageIfExist(opts.GetImageNameWithTag()) {
		if err := op.ImagePull(opts.GetImageNameWithTag()); err != nil {
			log.Debugf("Failed to pull the image: %s.", err)
			return err
		}
	}

	// 2. try to run the container
	log.Infof("Running container as the name <%s>", opts.ContainerName)
	if err := op.ContainerRun(opts.GetRunOpts()); err != nil {
		return fmt.Errorf("failed to run container: %v", err)
	}

	// 3. check if the volume is created successfully
	mounts, err := op.ContainerListMounts(opts.ContainerName)
	if err != nil {
		return fmt.Errorf("failed to get container mounts: %v", err)
	}
	volumes := mounts.ExtractSources()
	if docker.IfVolumesDiffer(volumes, opts.Volumes.ExtractHostPaths()) {
		return fmt.Errorf("failed to create volumes")
	}

	return nil
}

// ClearWhenInterruption will delete the container if the container fails to run
func ClearWhenInterruption(options configmanager.RawOptions) error {
	opts, err := NewOptions(options)
	if err != nil {
		return err
	}

	// 1. stop the container if it is running
	if ok := op.ContainerIfRunning(opts.ContainerName); ok {
		if err := op.ContainerStop(opts.ContainerName); err != nil {
			log.Errorf("Failed to stop container %s: %v", opts.ContainerName, err)
		}
	}

	// 2. remove the container if it exists
	if ok := op.ContainerIfExist(opts.ContainerName); ok {
		if err := op.ContainerRemove(opts.ContainerName); err != nil {
			log.Errorf("failed to remove container %v: %v", opts.ContainerName, err)
		}
	}

	var errs []error

	// 3. check if the container is stopped
	if ok := op.ContainerIfRunning(opts.ContainerName); ok {
		errs = append(errs, fmt.Errorf("failed to stop container %s", opts.ContainerName))
	}

	// 4. check if the container is removed
	if ok := op.ContainerIfExist(opts.ContainerName); ok {
		errs = append(errs, fmt.Errorf("failed to delete container %s", opts.ContainerName))
	}

	return multierr.Combine(errs...)
}

// DeleteAll will delete the container/volumes
func DeleteAll(options configmanager.RawOptions) error {
	opts, err := NewOptions(options)
	if err != nil {
		return err
	}

	// 1. stop the container if it is running
	if ok := op.ContainerIfRunning(opts.ContainerName); ok {
		if err := op.ContainerStop(opts.ContainerName); err != nil {
			log.Errorf("Failed to stop container %s: %v", opts.ContainerName, err)
		}
	}

	// 2. remove the container if it exists
	if ok := op.ContainerIfExist(opts.ContainerName); ok {
		if err := op.ContainerRemove(opts.ContainerName); err != nil {
			log.Errorf("failed to remove container %v: %v", opts.ContainerName, err)
		}
	}

	// 3. remove the volume if it exists
	if *opts.RmDataAfterDelete {
		volumesDirFromOptions := opts.Volumes.ExtractHostPaths()
		for _, err := range RemoveDirs(volumesDirFromOptions) {
			log.Error(err)
		}
	}

	var errs []error

	// 4. check if the container is stopped and deleted
	if ok := op.ContainerIfRunning(opts.ContainerName); ok {
		errs = append(errs, fmt.Errorf("failed to stop container %s", opts.ContainerName))
	}

	// 5. check if the container is removed
	if ok := op.ContainerIfExist(opts.ContainerName); ok {
		errs = append(errs, fmt.Errorf("failed to delete container %s", opts.ContainerName))
	}

	// 6. check if the volume is removed
	if *opts.RmDataAfterDelete {
		volumesDirFromOptions := opts.Volumes.ExtractHostPaths()
		for _, volume := range volumesDirFromOptions {
			if exist := PathExist(volume); exist {
				errs = append(errs, fmt.Errorf("failed to delete volume %s", volume))
			}
		}
	}

	return multierr.Combine(errs...)
}

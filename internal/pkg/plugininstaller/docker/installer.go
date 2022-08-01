package docker

import (
	"fmt"

	"go.uber.org/multierr"

	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/pkg/util/docker"
	"github.com/devstream-io/devstream/pkg/util/docker/dockersh"
	"github.com/devstream-io/devstream/pkg/util/log"
)

var op docker.Operator

func init() {
	// default to shell operator
	op = &dockersh.ShellOperator{}
}

func UseShellOperator() {
	op = &dockersh.ShellOperator{}
}

// InstallOrUpdate runs or updates the docker container
// note: any update will stop and remove the container, then run the new container
func InstallOrUpdate(options plugininstaller.RawOptions) error {
	opts, err := NewOptions(options)
	if err != nil {
		return err
	}

	// 1. try to pull the image
	// always pull the image because docker will check the image existence
	if err := op.ImagePull(opts.GetImageNameWithTag()); err != nil {
		return err
	}

	// 2. try to run the container
	log.Infof("Running container as the name <%s>", opts.ContainerName)
	if err := op.ContainerRun(opts.GetRunOpts()); err != nil {
		return fmt.Errorf("failed to run container: %v", err)
	}

	// 3. check if the container is started successfully
	if ok := op.ContainerIfRunning(opts.ContainerName); !ok {
		return fmt.Errorf("failed to run container")
	}

	// 4. check if the volume is created successfully
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

// HandleRunFailure will delete the container if the container fails to run
func HandleRunFailure(options plugininstaller.RawOptions) error {
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

// Delete will delete the container/image/volumes
func Delete(options plugininstaller.RawOptions) error {
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

	// 3. remove the image if it exists
	if ok := op.ImageIfExist(opts.GetImageNameWithTag()); ok {
		if err := op.ImageRemove(opts.GetImageNameWithTag()); err != nil {
			log.Errorf("failed to remove image %v: %v", opts.GetImageNameWithTag(), err)
		}
	}

	// 4. remove the volume if it exists
	if opts.RmDataAfterDelete {
		volumesDirFromOptions := opts.Volumes.ExtractHostPaths()
		for _, err := range RemoveDirs(volumesDirFromOptions) {
			log.Error(err)
		}
	}

	var errs []error

	// 5. check if the container is stopped and deleted
	if ok := op.ContainerIfRunning(opts.ContainerName); ok {
		errs = append(errs, fmt.Errorf("failed to stop container %s", opts.ContainerName))
	}

	// 6. check if the container is removed
	if ok := op.ContainerIfExist(opts.ContainerName); ok {
		errs = append(errs, fmt.Errorf("failed to delete container %s", opts.ContainerName))
	}

	// 7. check if the image is removed
	if ok := op.ImageIfExist(opts.GetImageNameWithTag()); ok {
		errs = append(errs, fmt.Errorf("failed to delete image %s", opts.GetImageNameWithTag()))
	}

	// 8. check if the volume is removed
	if opts.RmDataAfterDelete {
		volumesDirFromOptions := opts.Volumes.ExtractHostPaths()
		for _, volume := range volumesDirFromOptions {
			if exist := pathExist(volume); exist {
				errs = append(errs, fmt.Errorf("failed to delete volume %s", volume))
			}
		}
	}

	return multierr.Combine(errs...)
}

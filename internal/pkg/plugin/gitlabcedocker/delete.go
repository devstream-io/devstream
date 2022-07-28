package gitlabcedocker

import (
	"fmt"
	"os"

	"github.com/mitchellh/mapstructure"

	"github.com/devstream-io/devstream/pkg/util/log"
)

func Delete(options map[string]interface{}) (bool, error) {
	var opts Options
	if err := mapstructure.Decode(options, &opts); err != nil {
		return false, err
	}

	defaults(&opts)

	if errs := validate(&opts); len(errs) != 0 {
		for _, e := range errs {
			log.Errorf("Options error: %s.", e)
		}
		return false, fmt.Errorf("opts are illegal")
	}

	op := GetDockerOperator(opts)

	// 1. stop the container if it is running
	if ok := op.ContainerIfRunning(gitlabContainerName); ok {
		if err := op.ContainerStop(gitlabContainerName); err != nil {
			log.Errorf("Failed to stop container: %v", err)
		}
	}

	// 2. remove the container if it exists
	if ok := op.ContainerIfExist(gitlabContainerName); ok {
		if err := op.ContainerRemove(gitlabContainerName); err != nil {
			log.Errorf("failed to remove container %v: %v", gitlabContainerName, err)
		}
	}

	// 3. remove the image if it exists
	if ok := op.ImageIfExist(getImageNameWithTag(opts)); ok {
		if err := op.ImageRemove(getImageNameWithTag(opts)); err != nil {
			log.Errorf("failed to remove image %v: %v", getImageNameWithTag(opts), err)
		}
	}

	// 4. remove the volume if it exists
	volumesDirFromOptions := getVolumesDirFromOptions(opts)
	if opts.RmDataAfterDelete {
		for _, volume := range volumesDirFromOptions {
			if err := os.RemoveAll(volume); err != nil {
				log.Errorf("failed to remove data %v: %v", volume, err)
			}
		}
	}

	var errs []error

	// 1. check if the container is stopped and deleted
	if ok := op.ContainerIfRunning(gitlabContainerName); ok {
		errs = append(errs, fmt.Errorf("failed to delete/stop container %s", gitlabContainerName))
	}
	if ok := op.ContainerIfExist(gitlabContainerName); ok {
		errs = append(errs, fmt.Errorf("failed to delete container %s", gitlabContainerName))
	}

	// 2. check if the image is removed
	if ok := op.ImageIfExist(getImageNameWithTag(opts)); ok {
		errs = append(errs, fmt.Errorf("failed to delete image %s", getImageNameWithTag(opts)))
	}

	// 3. check if the data volume is removed
	if opts.RmDataAfterDelete {
		errs = append(errs, RemoveDirs(volumesDirFromOptions)...)
	}

	// splice the errors
	if len(errs) != 0 {
		errsString := ""
		for _, e := range errs {
			errsString += e.Error() + "; "
		}
		return false, fmt.Errorf(errsString)
	}

	return true, nil
}

// RemoveDirs removes the all the directories in the given list recursively
func RemoveDirs(dirs []string) []error {
	var errs []error
	for _, dir := range dirs {
		if err := os.RemoveAll(dir); err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}

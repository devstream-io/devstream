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

	if errs := validate(&opts); len(errs) != 0 {
		for _, e := range errs {
			log.Errorf("Options error: %s.", e)
		}
		return false, fmt.Errorf("opts are illegal")
	}

	op := getDockerOperator(opts)

	// 1. stop the container if it is running
	if ok := op.IfContainerRunning(gitlabContainerName); ok {
		if err := op.StopContainer(gitlabContainerName); err != nil {
			log.Errorf("Failed to stop container: %v", err)
		}
	}

	// 2. remove the container if it exists
	if ok := op.IfContainerExists(gitlabContainerName); ok {
		if err := op.RemoveContainer(gitlabContainerName); err != nil {
			log.Errorf("failed to remove container %v: %v", gitlabContainerName, err)
		}
	}

	// 3. remove the image if it exists
	if ok := op.IfImageExists(getImageNameWithTag(opts)); ok {
		if err := op.RemoveImage(getImageNameWithTag(opts)); err != nil {
			log.Errorf("failed to remove image %v: %v", getImageNameWithTag(opts), err)
		}
	}

	// 4. remove the volume if it exists
	if opts.RmDataAfterDelete {
		if err := os.RemoveAll(opts.GitLabHome); err != nil {
			log.Errorf("failed to remove data %v: %v", opts.GitLabHome, err)
		}
	}

	var errs []error

	// 1. check if the container is stopped and deleted
	if ok := op.IfContainerRunning(gitlabContainerName); ok {
		errs = append(errs, fmt.Errorf("failed to delete/stop container %s", gitlabContainerName))
	}
	if ok := op.IfContainerExists(gitlabContainerName); ok {
		errs = append(errs, fmt.Errorf("failed to delete container %s", gitlabContainerName))
	}

	// 2. check if the image is removed
	if ok := op.IfImageExists(getImageNameWithTag(opts)); ok {
		errs = append(errs, fmt.Errorf("failed to delete image %s", getImageNameWithTag(opts)))
	}

	// TODO: 3. check if data is removed successfully(if opts.RmDataAfterDelete is true)
	// note: the data is defined by opts.GitLabHome

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

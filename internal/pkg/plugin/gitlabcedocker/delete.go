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

	// 1. remove the image if it exists
	if ok := op.IfImageExists(gitlabImageName); ok {
		if err := op.RemoveImage(gitlabImageName); err != nil {
			log.Warnf("failed to remove image %v: %v", gitlabImageName, err)
		}
	}

	// 2. stop the container if it exists
	if ok := op.IfContainerRunning(gitlabContainerName); ok {
		if err := op.StopContainer(gitlabContainerName); err != nil {
			log.Errorf("Failed to stop container: %v", err)
		}
	}

	// 3. remove the container if it exists
	if ok := op.IfContainerExists(gitlabContainerName); ok {
		if err := op.RemoveContainer(gitlabContainerName); err != nil {
			log.Errorf("failed to remove container %v: %v", gitlabContainerName, err)
		}
	}

	// 4. remove the volume if it exists
	if opts.RmDataAfterDelete {
		// TODO: read user input to ask if the user really want to remove the data
		// you can use the /internal/pkg/pluginengine/helper.go/readUserInput function
		if err := os.RemoveAll(opts.GitLabHome); err != nil {
			log.Warnf("failed to remove data %v: %v", opts.GitLabHome, err)
		}
	}

	var errs []error

	// judge if the deletion is successful by checking if the container is running
	if ok := op.IfContainerRunning(gitlabContainerName); ok {
		errs = append(errs, fmt.Errorf("failed to delete/stop container %s", gitlabContainerName))
	}

	if ok := op.IfImageExists(gitlabImageName); ok {
		errs = append(errs, fmt.Errorf("failed to delete image %s", gitlabImageName))
	}

	if ok := op.IfContainerExists(gitlabContainerName); ok {
		errs = append(errs, fmt.Errorf("failed to delete container %s", gitlabContainerName))
	}

	// TODO: check if data is removed successfully(if opts.RmDataAfterDelete is true)
	// this is related to "read user input to ask if the user really want to remove the data"

	if len(errs) != 0 {
		errsString := ""
		for _, e := range errs {
			errsString += e.Error() + "; "
		}
		return false, fmt.Errorf(errsString)
	}

	return true, nil
}

package gitlabcedocker

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/mitchellh/mapstructure"

	"github.com/devstream-io/devstream/pkg/util/docker"
	"github.com/devstream-io/devstream/pkg/util/log"
)

func Update(options map[string]interface{}) (map[string]interface{}, error) {
	var opts Options
	if err := mapstructure.Decode(options, &opts); err != nil {
		return nil, err
	}

	defaults(&opts)

	if errs := validate(&opts); len(errs) != 0 {
		for _, e := range errs {
			log.Errorf("Options error: %s.", e)
		}
		return nil, fmt.Errorf("opts are illegal")
	}

	op := GetDockerOperator(opts)

	// 0. check if the volumes are the same
	mounts, err := op.ContainerListMounts(gitlabContainerName)
	if err != nil {
		return nil, fmt.Errorf("failed to get container mounts: %v", err)
	}
	volumesFromRunningContainer := mounts.ExtractSources()

	volumesDirFromOptions := getVolumesDirFromOptions(opts)

	if docker.IfVolumesDiffer(volumesFromRunningContainer, volumesDirFromOptions) {
		log.Warnf("You changed volumes of the container or change the gitlab home directory")
		log.Infof("Your volumes of the current container were: %v", strings.Join(volumesFromRunningContainer, " "))
		return nil, fmt.Errorf("sorry, you can't change the gitlab_home of the container once it's already been created")
	}

	// 1. stop the container
	if ok := op.ContainerIfRunning(gitlabContainerName); ok {
		if err := op.ContainerStop(gitlabContainerName); err != nil {
			log.Warnf("Failed to stop container: %v", err)
		}
	}

	// 2. remove the container if it exists
	if exists := op.ContainerIfExist(gitlabContainerName); exists {
		if err := op.ContainerRemove(gitlabContainerName); err != nil {
			return nil, fmt.Errorf("failed to remove container: %v", err)
		}
	}

	// 3. run the container with the new options
	if err := op.ContainerRun(buildDockerRunOptions(opts), dockerRunShmSizeParam); err != nil {
		return nil, fmt.Errorf("failed to run container: %v", err)
	}

	resource := gitlabResource{
		ContainerRunning: true,
		Volumes:          volumesFromRunningContainer,
		Hostname:         opts.Hostname,
		SSHPort:          strconv.Itoa(int(opts.SSHPort)),
		HTTPPort:         strconv.Itoa(int(opts.HTTPPort)),
		HTTPSPort:        strconv.Itoa(int(opts.HTTPSPort)),
	}

	return resource.toMap(), nil
}

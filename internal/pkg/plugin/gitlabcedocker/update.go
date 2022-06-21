package gitlabcedocker

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/mitchellh/mapstructure"

	"github.com/devstream-io/devstream/pkg/util/log"
)

func Update(options map[string]interface{}) (map[string]interface{}, error) {
	var opts Options
	if err := mapstructure.Decode(options, &opts); err != nil {
		return nil, err
	}

	if errs := validate(&opts); len(errs) != 0 {
		for _, e := range errs {
			log.Errorf("Options error: %s.", e)
		}
		return nil, fmt.Errorf("opts are illegal")
	}

	op := getDockerOperator(opts)

	// 0. check if the volumes are the same
	volumesDockerNow, err := op.ListContainerMounts(gitlabContainerName)
	if err != nil {
		return nil, fmt.Errorf("failed to get container mounts: %v", err)
	}

	volumesFromOptions := []string{
		filepath.Join(opts.GitLabHome, "config"),
		filepath.Join(opts.GitLabHome, "data"),
		filepath.Join(opts.GitLabHome, "logs"),
	}

	if ifVolumesDiffer(volumesDockerNow, volumesFromOptions) {
		log.Warnf("You changed volumes of the container or change the gitlab home directory")
		log.Infof("Your volumes of the current container were: %v", strings.Join(volumesDockerNow, " "))
		return nil, fmt.Errorf("sorry, you can't change the gitlab_home of the container once it's already been created")
	}

	// 1. stop the container
	if ok := op.IfContainerRunning(gitlabContainerName); ok {
		if err := op.StopContainer(gitlabContainerName); err != nil {
			log.Warnf("Failed to stop container: %v", err)
		}
	}

	// 2. remove the container if it exists
	if exists := op.IfContainerExists(gitlabContainerName); exists {
		if err := op.RemoveContainer(gitlabContainerName); err != nil {
			return nil, fmt.Errorf("failed to remove container: %v", err)
		}
	}

	// 3. run the container with the new options
	if err := op.RunContainer(opts); err != nil {
		return nil, fmt.Errorf("failed to run container: %v", err)
	}

	resource := gitlabResource{
		ContainerRunning: true,
		Volumes:          volumesDockerNow,
		Hostname:         opts.Hostname,
		SSHPort:          strconv.Itoa(int(opts.SSHPort)),
		HTTPPort:         strconv.Itoa(int(opts.HTTPPort)),
		HTTPSPort:        strconv.Itoa(int(opts.HTTPSPort)),
	}

	return resource.toMap(), nil
}

func ifVolumesDiffer(volumesBefore, volumesCurrent []string) bool {
	if len(volumesBefore) != len(volumesCurrent) {
		return true
	}

	beforeMap := make(map[string]struct{})
	currentMap := make(map[string]struct{})

	for _, v := range volumesBefore {
		beforeMap[v] = struct{}{}
	}
	for _, v := range volumesCurrent {
		currentMap[v] = struct{}{}
	}

	for _, v := range volumesBefore {
		if _, ok := currentMap[v]; !ok {
			return true
		}
	}

	return false
}

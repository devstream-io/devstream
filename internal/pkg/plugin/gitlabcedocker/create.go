package gitlabcedocker

import (
	"fmt"
	"strconv"

	"github.com/mitchellh/mapstructure"

	"github.com/devstream-io/devstream/pkg/util/log"
)

func Create(options map[string]interface{}) (map[string]interface{}, error) {
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

	// 1. try to pull the image
	// always pull the image because docker will check the image existence
	if err := op.PullImage(getImageNameWithTag(opts)); err != nil {
		return nil, err
	}

	// 2. try to run the container
	log.Info("Running container as the name <gitlab>")
	if err := op.RunContainer(opts); err != nil {
		return nil, fmt.Errorf("failed to run container: %v", err)
	}

	// 3. check if the container is started successfully
	if ok := op.IfContainerRunning(gitlabContainerName); !ok {
		return nil, fmt.Errorf("failed to run container")
	}

	volumes, err := op.ListContainerMounts(gitlabContainerName)
	if err != nil {
		return nil, fmt.Errorf("failed to get container mounts: %v", err)
	}

	resource := gitlabResource{
		ContainerRunning: true,
		Volumes:          volumes,
		Hostname:         opts.Hostname,
		SSHPort:          strconv.Itoa(int(opts.SSHPort)),
		HTTPPort:         strconv.Itoa(int(opts.HTTPPort)),
		HTTPSPort:        strconv.Itoa(int(opts.HTTPSPort)),
	}

	return resource.toMap(), nil
}

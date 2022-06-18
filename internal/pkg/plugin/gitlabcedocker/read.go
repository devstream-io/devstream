package gitlabcedocker

import (
	"fmt"

	"github.com/mitchellh/mapstructure"

	"github.com/devstream-io/devstream/pkg/util/log"
)

func Read(options map[string]interface{}) (map[string]interface{}, error) {
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

	running := op.IfContainerRunning(gitlabContainerName)

	volumes, err := op.ListContainerMounts(gitlabContainerName)
	if err != nil {
		return nil, fmt.Errorf("failed to get container mounts: %v", err)
	}

	return buildState(running, volumes), nil
}

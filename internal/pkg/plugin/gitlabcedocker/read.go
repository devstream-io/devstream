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

	// 1. get running status
	running := op.IfContainerRunning(gitlabContainerName)
	if !running {
		return (&gitlabResource{}).toMap(), nil
	}

	// 2. get volumes(gitlab_home)
	volumes, err := op.ListContainerMounts(gitlabContainerName)
	if err != nil {
		// `Read` shouldn't return errors even if failed to read ports, volumes, hostname.
		// because:
		// 1. when the docker is stopped it could cause these errors.
		// 2. if Read failed, the following steps contain the docker's restart will be aborted.
		log.Errorf("failed to get container mounts: %v", err)
	}

	// 3. get hostname
	hostname, err := op.GetContainerHostname(gitlabContainerName)
	if err != nil {
		log.Errorf("failed to get container hostname: %v", err)
	}

	// 4. get port bindings
	SSHPort, err := op.GetContainerPortBinding(gitlabContainerName, "22", tcp)
	if err != nil {
		log.Errorf("failed to get container ssh port: %v", err)
	}
	HTTPPort, err := op.GetContainerPortBinding(gitlabContainerName, "80", tcp)
	if err != nil {
		log.Errorf("failed to get container http port: %v", err)
	}
	HTTPSPort, err := op.GetContainerPortBinding(gitlabContainerName, "443", tcp)
	if err != nil {
		log.Errorf("failed to get container https port: %v", err)
	}

	// if the previous steps failed, the parameters will be empty
	// so dtm will find the resource is drifted and restart docker
	resource := gitlabResource{
		ContainerRunning: running,
		Volumes:          volumes,
		Hostname:         hostname,
		SSHPort:          SSHPort,
		HTTPPort:         HTTPPort,
		HTTPSPort:        HTTPSPort,
	}

	return resource.toMap(), nil
}

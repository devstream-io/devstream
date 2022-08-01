package docker

import (
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/pkg/util/docker"
	"github.com/devstream-io/devstream/pkg/util/log"
)

type State struct {
	ContainerRunning bool
	Volumes          []string
	Hostname         string
	PortPublishes    []docker.PortPublish
}

func (s *State) toMap() map[string]interface{} {
	return map[string]interface{}{
		"ContainerRunning": s.ContainerRunning,
		"Volumes":          s.Volumes,
		"Hostname":         s.Hostname,
		"PortPublishes":    s.PortPublishes,
	}
}

func GetStaticStateFromOptions(options plugininstaller.RawOptions) (map[string]interface{}, error) {
	opts, err := NewOptions(options)
	if err != nil {
		return nil, err
	}

	staticState := &State{
		ContainerRunning: true,
		Volumes:          opts.Volumes.ExtractHostPaths(),
		Hostname:         opts.Hostname,
		PortPublishes:    opts.PortPublishes,
	}

	return staticState.toMap(), nil
}

func GetRunningState(options plugininstaller.RawOptions) (map[string]interface{}, error) {
	opts, err := NewOptions(options)
	if err != nil {
		return nil, err
	}

	// 1. get running status
	running := op.ContainerIfRunning(opts.ContainerName)
	if !running {
		return map[string]interface{}{}, nil
	}

	// 2. get volumes
	mounts, err := op.ContainerListMounts(opts.ContainerName)
	if err != nil {
		// `Read` shouldn't return errors even if failed to read ports, volumes, hostname.
		// because:
		// 1. when the docker is stopped it could cause these errors.
		// 2. if Read failed, the following steps contain the docker's restart will be aborted.
		log.Errorf("failed to get container mounts: %v", err)
	}
	volumes := mounts.ExtractSources()

	// 3. get hostname
	hostname, err := op.ContainerGetHostname(opts.ContainerName)
	if err != nil {
		log.Errorf("failed to get container hostname: %v", err)
	}

	// 4. get port bindings
	PortPublishes, err := op.ContainerListPortPublishes(opts.ContainerName)
	if err != nil {
		log.Errorf("failed to get container port publishes: %v", err)
	}

	// if the previous steps failed, the parameters will be empty
	// so dtm will find the resource is drifted and restart docker
	resource := &State{
		ContainerRunning: running,
		Volumes:          volumes,
		Hostname:         hostname,
		PortPublishes:    PortPublishes,
	}

	return resource.toMap(), nil
}

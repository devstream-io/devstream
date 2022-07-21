package docker

import (
	"sort"

	mapset "github.com/deckarep/golang-set/v2"
)

// Operator is an interface for docker operations
// It is implemented by shDockerOperator
// in the future, we can add other implementations such as SdkDockerOperator
type Operator interface {
	ImageIfExist(imageNameWithTag string) bool
	ImagePull(imageNameWithTag string) error
	ImageRemove(imageNameWithTag string) error

	ContainerIfExist(containerName string) bool
	ContainerIfRunning(containerName string) bool
	// ContainerRun runs a container with the given options
	// params is a list of additional parameters for docker run
	// params will be appended to the end of the command
	ContainerRun(opts RunOptions, params ...string) error
	ContainerStop(containerName string) error
	ContainerRemove(containerName string) error

	// ContainerListMounts lists container mounts
	ContainerListMounts(containerName string) (Mounts, error)

	ContainerGetHostname(containerName string) (string, error)
	ContainerGetPortBinding(containerName, containerPort, protocol string) (hostPort string, err error)
}

type MountPoint struct {
	Type        string `json:"Type"`
	Source      string `json:"Source"`
	Destination string `json:"Destination"`
	Mode        string `json:"Mode"`
	Rw          bool   `json:"RW"`
	Propagation string `json:"Propagation"`
}

type Mounts []MountPoint

// ExtractSources returns a list of sources for the given mounts
func (ms Mounts) ExtractSources() []string {
	sources := make([]string, 0)
	for _, mount := range ms {
		sources = append(sources, mount.Source)
	}
	sort.Slice(sources, func(i, j int) bool {
		return sources[i] < sources[j]
	})

	return sources
}

func IfVolumesDiffer(volumesBefore, volumesCurrent []string) bool {
	beforeSet := mapset.NewSet[string](volumesBefore...)
	currentSet := mapset.NewSet[string](volumesCurrent...)

	return !beforeSet.Equal(currentSet)
}

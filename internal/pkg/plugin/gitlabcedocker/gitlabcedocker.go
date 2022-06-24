package gitlabcedocker

import (
	"fmt"
	"strings"
)

const (
	gitlabImageName = "gitlab/gitlab-ce"
	// TODO expose image tag to user in config file to customize
	gitlabImageTag = "rc"

	tcp = "tcp"
)

var (
	gitlabImageNameWithTag = fmt.Sprintf("%v:%v", gitlabImageName, gitlabImageTag)
	gitlabContainerName    = "gitlab"
)

// dockerOperator is an interface for docker operations
// It is implemented by sshDockerOperator
// in the future, we can add other implementations such as sshDockerOperator
type dockerOperator interface {
	IfImageExists(imageName string) bool
	PullImage(image string) error
	RemoveImage(image string) error

	IfContainerExists(container string) bool
	IfContainerRunning(container string) bool
	RunContainer(options Options) error
	StopContainer(container string) error
	RemoveContainer(container string) error

	ListContainerMounts(container string) ([]string, error)

	GetContainerHostname(container string) (string, error)
	GetContainerPortBinding(container, containerPort, protocol string) (hostPort string, err error)
}

func getDockerOperator(_ Options) dockerOperator {
	// just return a sshDockerOperator for now
	return &sshDockerOperator{}
}

type gitlabResource struct {
	ContainerRunning bool
	Volumes          []string
	Hostname         string
	SSHPort          string
	HTTPPort         string
	HTTPSPort        string
}

func (res *gitlabResource) toMap() map[string]interface{} {
	return map[string]interface{}{
		"containerRunning": res.ContainerRunning,
		"volumes":          strings.Join(res.Volumes, ","),
		"hostname":         res.Hostname,
		"SSHPort":          res.SSHPort,
		"HTTPPort":         res.HTTPPort,
		"HTTPSPort":        res.HTTPSPort,
	}
}

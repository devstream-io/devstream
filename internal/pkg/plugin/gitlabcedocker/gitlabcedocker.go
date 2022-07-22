package gitlabcedocker

import (
	"strings"

	"github.com/devstream-io/devstream/pkg/util/docker"
	"github.com/devstream-io/devstream/pkg/util/docker/dockersh"
)

const (
	gitlabImageName       = "gitlab/gitlab-ce"
	defaultImageTag       = "rc"
	gitlabContainerName   = "gitlab"
	tcp                   = "tcp"
	dockerRunShmSizeParam = "--shm-size 256m"
)

func getImageNameWithTag(opt Options) string {
	return gitlabImageName + ":" + opt.ImageTag
}

func defaults(opts *Options) {
	if opts.ImageTag == "" {
		opts.ImageTag = defaultImageTag
	}
}

func GetDockerOperator(_ Options) docker.Operator {
	// just return a ShellOperator for now
	return &dockersh.ShellOperator{}
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

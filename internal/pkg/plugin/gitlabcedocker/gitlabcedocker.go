package gitlabcedocker

import "github.com/devstream-io/devstream/pkg/util/types"

const Name = "gitlab-ce-docker"

const (
	defaultHostname       = "gitlab.example.com"
	defaultGitlabHome     = "/srv/gitlab"
	defaultSSHPort        = 22
	defaultHTTPPort       = 80
	defaultHTTPSPort      = 443
	defaultImageTag       = "rc"
	gitlabImageName       = "gitlab/gitlab-ce"
	gitlabContainerName   = "gitlab"
	dockerRunShmSizeParam = "--shm-size 256m"
)

var (
	defaultRMDataAfterDelete = types.Bool(false)
)

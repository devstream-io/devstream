package gitlabcedocker

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
	rmDataAfterDelete        = false
	defaultRMDataAfterDelete = &rmDataAfterDelete
)

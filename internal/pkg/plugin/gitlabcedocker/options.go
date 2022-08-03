package gitlabcedocker

import (
	"path/filepath"

	dockerInstaller "github.com/devstream-io/devstream/internal/pkg/plugininstaller/docker"
	"github.com/devstream-io/devstream/pkg/util/docker"
)

// Options is the struct for configurations of the gitlab-ce-docker plugin.
type Options struct {
	// GitLab home directory, we assume the path set by user is always correct.
	GitLabHome        string `validate:"required" mapstructure:"gitlab_home"`
	Hostname          string `validate:"required,hostname" mapstructure:"hostname"`
	SSHPort           uint   `validate:"required" mapstructure:"ssh_port"`
	HTTPPort          uint   `validate:"required" mapstructure:"http_port"`
	HTTPSPort         uint   `validate:"required" mapstructure:"https_port"`
	RmDataAfterDelete bool   `mapstructure:"rm_data_after_delete"`
	ImageTag          string `mapstructure:"image_tag"`
}

func buildDockerOptions(opts *Options) *dockerInstaller.Options {
	dockerOpts := &dockerInstaller.Options{}
	dockerOpts.RmDataAfterDelete = opts.RmDataAfterDelete
	dockerOpts.ImageName = gitlabImageName
	dockerOpts.ImageTag = opts.ImageTag
	dockerOpts.Hostname = opts.Hostname
	dockerOpts.ContainerName = gitlabContainerName
	dockerOpts.RestartAlways = true

	portPublishes := []docker.PortPublish{
		{HostPort: opts.SSHPort, ContainerPort: 22},
		{HostPort: opts.HTTPPort, ContainerPort: 80},
		{HostPort: opts.HTTPSPort, ContainerPort: 443},
	}
	dockerOpts.PortPublishes = portPublishes

	dockerOpts.Volumes = buildDockerVolumes(opts)

	dockerOpts.RunParams = []string{dockerRunShmSizeParam}

	return dockerOpts
}

func buildDockerVolumes(opts *Options) docker.Volumes {
	volumes := []docker.Volume{
		{HostPath: filepath.Join(opts.GitLabHome, "config"), ContainerPath: "/etc/gitlab"},
		{HostPath: filepath.Join(opts.GitLabHome, "data"), ContainerPath: "/var/opt/gitlab"},
		{HostPath: filepath.Join(opts.GitLabHome, "logs"), ContainerPath: "/var/log/gitlab"},
	}

	return volumes
}

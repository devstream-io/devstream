package gitlabcedocker

import (
	"fmt"
	"path/filepath"

	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	dockerInstaller "github.com/devstream-io/devstream/internal/pkg/plugininstaller/docker"
	"github.com/devstream-io/devstream/pkg/util/docker"
	"github.com/devstream-io/devstream/pkg/util/log"
)

// Options is the struct for configurations of the gitlab-ce-docker plugin.
type Options struct {
	Hostname string `validate:"hostname" mapstructure:"hostname"`
	// GitLab home directory, we assume the path set by user is always correct.
	GitLabHome        string `mapstructure:"gitlab_home"`
	SSHPort           uint   `mapstructure:"ssh_port"`
	HTTPPort          uint   `mapstructure:"http_port"`
	HTTPSPort         uint   `mapstructure:"https_port"`
	RmDataAfterDelete *bool  `mapstructure:"rm_data_after_delete"`
	ImageTag          string `mapstructure:"image_tag"`
}

func (opts *Options) Defaults() {
	if opts.Hostname == "" {
		opts.Hostname = defaultHostname
	}
	if opts.GitLabHome == "" {
		opts.GitLabHome = defaultGitlabHome
	}
	if opts.SSHPort == 0 {
		opts.SSHPort = defaultSSHPort
	}
	if opts.HTTPPort == 0 {
		opts.HTTPPort = defaultHTTPPort
	}
	if opts.HTTPSPort == 0 {
		opts.HTTPSPort = defaultHTTPSPort
	}
	if opts.RmDataAfterDelete == nil {
		opts.RmDataAfterDelete = defaultRMDataAfterDelete
	}
	if opts.ImageTag == "" {
		opts.ImageTag = defaultImageTag
	}
}

// gitlabURL is the access URL of GitLab.
var gitlabURL string

func (opts *Options) setGitLabURL() {
	if gitlabURL != "" {
		return
	}
	gitlabURL = fmt.Sprintf("http://%s:%d", opts.Hostname, opts.HTTPPort)
}

func showHelpMsg(options plugininstaller.RawOptions) error {
	log.Infof("GitLab access URL: %s", gitlabURL)
	log.Infof("GitLab initial root password: execute the command -> docker exec -it gitlab grep 'Password:' /etc/gitlab/initial_root_password")

	return nil
}

func buildDockerOptions(opts *Options) *dockerInstaller.Options {
	portPublishes := []docker.PortPublish{
		{HostPort: opts.SSHPort, ContainerPort: 22},
		{HostPort: opts.HTTPPort, ContainerPort: 80},
		{HostPort: opts.HTTPSPort, ContainerPort: 443},
	}

	dockerOpts := &dockerInstaller.Options{
		RmDataAfterDelete: opts.RmDataAfterDelete,
		ImageName:         gitlabImageName,
		ImageTag:          opts.ImageTag,
		Hostname:          opts.Hostname,
		ContainerName:     gitlabContainerName,
		RestartAlways:     true,
		Volumes:           buildDockerVolumes(opts),
		RunParams:         []string{dockerRunShmSizeParam},
		PortPublishes:     portPublishes,
	}

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

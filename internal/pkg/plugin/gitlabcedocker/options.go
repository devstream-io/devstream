package gitlabcedocker

import "path/filepath"

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

func getVolumesDirFromOptions(opts Options) []string {
	volumesDirFromOptions := []string{
		filepath.Join(opts.GitLabHome, "config"),
		filepath.Join(opts.GitLabHome, "data"),
		filepath.Join(opts.GitLabHome, "logs"),
	}

	return volumesDirFromOptions
}

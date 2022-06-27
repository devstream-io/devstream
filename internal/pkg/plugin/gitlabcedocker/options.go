package gitlabcedocker

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

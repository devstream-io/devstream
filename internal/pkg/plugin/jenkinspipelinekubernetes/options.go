package jenkinspipelinekubernetes

const (
	defaultJenkinsUser               = "admin"
	defaultJenkinsPipelineScriptPath = "Jenkinsfile"
)

// Options is the struct for configurations of the jenkins-pipeline-kubernetes plugin.
type Options struct {
	J             *JenkinsOptions `mapstructure:"jenkins"`
	GitHubToken   string          `mapstructure:"githubToken"`
	GitHubRepoURL string          `mapstructure:"githubRepoUrl" validate:"required"`
}

type JenkinsOptions struct {
	URL                string `mapstructure:"url" validate:"required,hostname_port"`
	User               string `mapstructure:"user" validate:"required"`
	Password           string `mapstructure:"password"`
	JobName            string `mapstructure:"jobName"`
	PipelineScriptPath string `mapstructure:"pipelineScriptPath"`
}

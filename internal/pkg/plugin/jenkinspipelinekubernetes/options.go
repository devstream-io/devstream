package jenkinspipelinekubernetes

const (
	defaultJenkinsUser               = "admin"
	defaultJenkinsPipelineScriptPath = "Jenkinsfile"
)

// Options is the struct for configurations of the jenkins-pipeline-kubernetes plugin.
type Options struct {
	J             *JenkinsOption `mapstructure:"jenkins"`
	GitHubToken   string         `mapstructure:"githubToken"`
	GitHubRepoURL string         `mapstructure:"githubRepoUrl" validateAndHandleOptions:"required"`
}

type JenkinsOption struct {
	URL                string `mapstructure:"url" validateAndHandleOptions:"required,hostname_port"`
	User               string `mapstructure:"user" validateAndHandleOptions:"required"`
	Password           string `mapstructure:"password"`
	JobName            string `mapstructure:"jobName"`
	PipelineScriptPath string `mapstructure:"pipelineScriptPath"`
}

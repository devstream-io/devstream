package jenkinspipelinekubernetes

const (
	defaultJenkinsUser               = "admin"
	defaultJenkinsPipelineScriptPath = "Jenkinsfile"
)

// Options is the struct for configurations of the jenkins-pipeline-kubernetes plugin.
type Options struct {
	JenkinsURL        string `mapstructure:"jenkinsURL" validate:"required,hostname_port"`
	JenkinsUser       string `mapstructure:"jenkinsUser" validate:"required"`
	JobName           string `mapstructure:"jobName"`
	JenkinsfilePath   string `mapstructure:"jenkinsfilePath"`
	JenkinsfileScmURL string `mapstructure:"jenkinsfileScmURL" validate:"required"`
}

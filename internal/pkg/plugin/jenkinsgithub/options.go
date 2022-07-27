package jenkinsgithub

const (
	defaultJenkinsUser               = "admin"
	defaultJenkinsPipelineScriptPath = "Jenkinsfile-pr"
)

// Options is the struct for configurations of the jenkins-pipeline-kubernetes plugin.
type (
	Options struct {
		J             *JenkinsOptions `mapstructure:"jenkins"`
		Helm          *HelmOptions    `mapstructure:"helm" validate:"required"`
		GitHubToken   string          `mapstructure:"githubToken"`
		GitHubRepoURL string          `mapstructure:"githubRepoUrl" validate:"required"`
		AdminList     []string        `mapstructure:"adminList"`
	}

	JenkinsOptions struct {
		URL                string `mapstructure:"url" validate:"required,hostname_port"`
		URLOverride        string `mapstructure:"urlOverride"`
		User               string `mapstructure:"user" validate:"required"`
		Password           string `mapstructure:"password"`
		JobName            string `mapstructure:"jobName"`
		PipelineScriptPath string `mapstructure:"pipelineScriptPath"`
	}

	HelmOptions struct {
		ReleaseName string `mapstructure:"releaseName" validate:"required"`
		Namespace   string `mapstructure:"namespace"`
	}

	GitHubIntegOptions struct {
		AdminList          []string // GitHub admin list
		CredentialsID      string
		JenkinsURLOverride string
		GithubAuthID       string
	}

	JobOptions struct {
		JobName              string `mapstructure:"jobName"`
		PrPipelineScriptPath string
		GitHubRepoURL        string
		AdminList            []string
		CredentialsID        string
	}
)

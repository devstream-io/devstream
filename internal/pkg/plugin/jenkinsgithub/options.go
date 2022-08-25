package jenkinsgithub

const (
	defaultJenkinsUser               = "admin"
	defaultJenkinsPipelineScriptPath = "Jenkinsfile-pr"
)

// Options is the struct for configurations of the jenkins-pipeline-kubernetes plugin.
type (
	Options struct {
		J    *JenkinsOptions `mapstructure:"jenkins"`
		Helm *HelmOptions    `mapstructure:"helm" validate:"required"`
		// TODO(aFlyBird0): maybe these three options can be put into a struct called "GitHubOptions"
		GitHubToken   string   `mapstructure:"githubToken"`
		GitHubRepoURL string   `mapstructure:"githubRepoUrl" validate:"required"`
		AdminList     []string `mapstructure:"adminList"`
	}

	JenkinsOptions struct {
		URL         string `mapstructure:"url" validate:"required,hostname_port"`
		URLOverride string `mapstructure:"urlOverride"`
		User        string `mapstructure:"user" validate:"required"`
		Password    string `mapstructure:"password"`
		// now we only have job-pr,
		// we should change the config name if there is job-main in the future
		JobName string `mapstructure:"jobName"`
		// as same as the job name
		PipelineScriptPath string `mapstructure:"pipelineScriptPath"`
	}

	HelmOptions struct {
		ReleaseName string `mapstructure:"releaseName" validate:"required"`
		Namespace   string `mapstructure:"namespace"`
	}

	// GitHubIntegOptions is the struct for configurations of the GitHub pull request builder plugin
	GitHubIntegOptions struct {
		AdminList          []string // GitHub admin list
		CredentialsID      string
		JenkinsURLOverride string
		GithubAuthID       string
	}

	// JobOptions is the struct to render job xml.
	// TODO(aFlyBird0): it seems that GitHubProjectURL and GitSCMURL could be set separately,
	// now they are set together. Need more investigation.
	// TODO(aFlyBird0): figure out what does diffent format of GitHub repo JenkinsURL mean and handle it.
	// such as https://.../, https://.../repo.git, git@......, etc.
	JobOptions struct {
		JobName              string
		PrPipelineScriptPath string
		GitHubRepoURL        string
		AdminList            []string
		CredentialsID        string
	}
)

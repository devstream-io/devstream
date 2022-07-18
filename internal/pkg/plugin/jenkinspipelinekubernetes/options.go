package jenkinspipelinekubernetes

import "fmt"

// Options is the struct for configurations of the jenkins-pipeline-kubernetes plugin.
type Options struct {
	JenkinsURL     string `mapstructure:"jenkinsUrl" validate:"required,hostname_port"`
	JenkinsUser    string `mapstructure:"jenkinsUser" validate:"required"`
	JenkinsToken   string `mapstructure:"jenkinsToken"`
	GitHubToken    string `mapstructure:"githubToken"`
	GitHubRepoURL  string `mapstructure:"githubRepoUrl" validate:"required"`
	JenkinsJobName string `mapstructure:"jenkinsJobName" validate:"required"`
	// TODO(aFlyBird0): add options to configure the script path in GitHub repo, now it is hardcoded to "Jenkinsfile"
}

func (options *Options) GetJenkinsAccessURL() string {
	return fmt.Sprintf("http://%s:%s@%s", options.JenkinsUser, options.JenkinsToken, options.JenkinsURL)
}

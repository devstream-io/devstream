package jenkinspipelinekubernetes

import (
	"fmt"
	"os"

	"github.com/devstream-io/devstream/pkg/util/validator"
)

// validateAndHandleOptions validates and pre handle the options provided by the core.
func validateAndHandleOptions(options *Options) []error {
	validateErrs := validate(options)

	defaults(options)

	envErrs := handleEnv(options)

	return append(validateErrs, envErrs...)
}

func validate(options *Options) []error {
	return validator.Struct(options)
}

func defaults(options *Options) {
	// pre handle the options
	if options.J.PipelineScriptPath == "" {
		options.J.PipelineScriptPath = defaultJenkinsPipelineScriptPath
	}

	if options.J.User == "" {
		options.J.User = defaultJenkinsUser
	}

	options.J.URL = "http://" + options.J.URL
}

func handleEnv(options *Options) []error {
	var errs []error

	options.GitHubToken = os.Getenv("GITHUB_TOKEN")
	if options.GitHubToken == "" {
		errs = append(errs, fmt.Errorf("env GITHUB_TOKEN is required"))
	}

	// TODO(aFlyBird0): now jenkins token should be provided by the user,
	// so, user should install jenkins first and stop to set the token in env, then install this pipeline plugin.
	// could we generate a token automatically in "jenkins" plugin?
	// and put it into .outputs of "jenkins" plugin,
	// so that user could run "jenkins" and "jenkins-pipeline-kubernetes"  in the same tool file.(using depends on).
	//options.J.Token = os.Getenv("JENKINS_TOKEN")
	//if options.J.Token == "" {
	//	errs = append(errs, fmt.Errorf("env JENKINS_TOKEN is required"))
	//}

	// read the password from config file(including the outputs from last plugin) first, then from env
	if options.J.Password == "" {
		options.J.Password = os.Getenv("JENKINS_PASSWORD")
	}

	if options.J.Password == "" {
		errs = append(errs, fmt.Errorf("env JENKINS_PASSWORD is required"))
	}

	return errs
}

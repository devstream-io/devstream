package jenkinsgithub

import (
	"fmt"
	"os"
	"strings"

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
	if options.Helm.Namespace == "" {
		options.Helm.Namespace = "jenkins"
	}

	if options.J.PipelineScriptPath == "" {
		options.J.PipelineScriptPath = defaultJenkinsPipelineScriptPath
	}

	if options.J.User == "" {
		options.J.User = defaultJenkinsUser
	}

	options.J.URL = "http://" + options.J.URL

	if !strings.HasSuffix(options.J.URLOverride, "/") {
		options.J.URLOverride += "/"
	}
}

func handleEnv(options *Options) []error {
	var errs []error

	options.GitHubToken = os.Getenv("GITHUB_TOKEN")
	if options.GitHubToken == "" {
		errs = append(errs, fmt.Errorf("env GITHUB_TOKEN is required"))
	}

	// read the password from config file(including the outputs from last plugin) first, then from env
	if options.J.Password == "" {
		options.J.Password = os.Getenv("JENKINS_PASSWORD")
	}

	if options.J.Password == "" {
		errs = append(errs, fmt.Errorf("jenkins password is required"))
	}

	return errs
}

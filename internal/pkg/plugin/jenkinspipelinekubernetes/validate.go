package jenkinspipelinekubernetes

import (
	"fmt"
	"os"
	"strings"

	"github.com/devstream-io/devstream/pkg/util/validator"
)

// validate validates the options provided by the core.
func validate(options *Options) []error {

	// pre-handle options to remove "http://" from JenkinsURL
	preHandleOptions(options)

	var retErrs []error

	if errs := validator.Struct(options); len(errs) != 0 {
		retErrs = append(retErrs, errs...)
	}

	options.GitHubToken = os.Getenv("GITHUB_TOKEN")
	options.JenkinsToken = os.Getenv("JENKINS_TOKEN")
	if options.GitHubToken == "" {
		retErrs = append(retErrs, fmt.Errorf("GITHUB_TOKEN is required"))
	}
	if options.JenkinsToken == "" {
		retErrs = append(retErrs, fmt.Errorf("JENKINS_TOKEN is required"))
	}

	// todo: check if the jenkins url is valid (try to connect to jenkins)

	return retErrs
}

func preHandleOptions(options *Options) {
	options.JenkinsURL = strings.Replace(options.JenkinsURL, "http://", "", 1)
}

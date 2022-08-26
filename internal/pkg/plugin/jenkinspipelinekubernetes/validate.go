package jenkinspipelinekubernetes

import (
	"fmt"
	"os"
	"strings"

	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/validator"
)

func ValidateAndDefaults(options plugininstaller.RawOptions) (plugininstaller.RawOptions, error) {
	opts, err := NewOptions(options)
	if err != nil {
		return nil, err
	}
	defaults(opts)

	if errs := validate(opts); len(errs) != 0 {
		for _, e := range errs {
			log.Errorf("Validate error: %s.", e)
		}
		return nil, fmt.Errorf("validate failed")
	}

	if errs := initPasswdFromEnvVars(); len(errs) != 0 {
		for _, e := range errs {
			log.Errorf("Init password from env vars failed: %s.", e)
		}
		return nil, fmt.Errorf("init failed")
	}

	return opts.Encode()
}

func validate(options *Options) []error {
	return validator.Struct(options)
}

func defaults(options *Options) {
	// pre handle the options
	if options.JenkinsfilePath == "" {
		options.JenkinsfilePath = defaultJenkinsPipelineScriptPath
	}

	if options.JenkinsUser == "" {
		options.JenkinsUser = defaultJenkinsUser
	}

	if !strings.Contains(options.JenkinsURL, "http") {
		options.JenkinsURL = "http://" + options.JenkinsURL
	}
}

func initPasswdFromEnvVars() []error {
	var errs []error

	githubToken = os.Getenv("GITHUB_TOKEN")
	if githubToken == "" {
		errs = append(errs, fmt.Errorf("env GITHUB_TOKEN is required"))
	}

	// TODO(aFlyBird0): now jenkins token should be provided by the user,
	// so, user should install jenkins first and stop to set the token in env, then install this pipeline plugin.
	// could we generate a token automatically in "jenkins" plugin?
	// and put it into .outputs of "jenkins" plugin,
	// so that user could run "jenkins" and "jenkins-pipeline-kubernetes"  in the same tool file.(using depends on).
	//options.Token = os.Getenv("JENKINS_TOKEN")
	//if options.Token == "" {
	//	errs = append(errs, fmt.Errorf("env JENKINS_TOKEN is required"))
	//}

	// read the password from config file(including the outputs from last plugin) first, then from env
	if jenkinsPassword == "" {
		jenkinsPassword = os.Getenv("JENKINS_PASSWORD")
	}

	if jenkinsPassword == "" {
		errs = append(errs, fmt.Errorf("jenkins password is required"))
	}

	return errs
}

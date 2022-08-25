package jenkinspipelinekubernetes

import (
	"github.com/mitchellh/mapstructure"

	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
)

const (
	defaultJenkinsUser               = "admin"
	defaultJenkinsPipelineScriptPath = "Jenkinsfile"
)

var (
	jenkinsPassword string
	githubToken     string
)

// Options is the struct for configurations of the jenkins-pipeline-kubernetes plugin.
type Options struct {
	JenkinsURL        string `mapstructure:"jenkinsURL" validate:"required,url"`
	JenkinsUser       string `mapstructure:"jenkinsUser" validate:"required"`
	JobName           string `mapstructure:"jobName"`
	JenkinsfilePath   string `mapstructure:"jenkinsfilePath"`
	JenkinsfileScmURL string `mapstructure:"jenkinsfileScmURL" validate:"required"`
}

// NewOptions create options by raw options
func NewOptions(options plugininstaller.RawOptions) (*Options, error) {
	var opts Options
	if err := mapstructure.Decode(options, &opts); err != nil {
		return nil, err
	}
	return &opts, nil
}

func (opts *Options) Encode() (map[string]interface{}, error) {
	var options map[string]interface{}
	if err := mapstructure.Decode(opts, &options); err != nil {
		return nil, err
	}
	return options, nil
}

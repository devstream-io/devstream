package jenkins

import (
	"reflect"
	"strings"

	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/ci"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/jenkins/plugins"
	"github.com/devstream-io/devstream/pkg/util/log"
)

type Pipeline struct {
	JobName         string                          `mapstructure:"jobName" validate:"required"`
	JenkinsfilePath string                          `mapstructure:"jenkinsfilePath" validate:"required"`
	ImageRepo       *plugins.ImageRepoJenkinsConfig `mapstructure:"imageRepo"`
	Dingtalk        *plugins.DingtalkJenkinsConfig  `mapstructure:"dingTalk"`
	Sonarqube       *plugins.SonarQubeJenkinsConfig `mapstructure:"sonarqube"`
}

func (p *Pipeline) getJobName() string {
	if strings.Contains(p.JobName, "/") {
		return strings.Split(p.JobName, "/")[1]
	}
	return p.JobName
}

func (p *Pipeline) getJobFolder() string {
	if strings.Contains(p.JobName, "/") {
		return strings.Split(p.JobName, "/")[0]
	}
	return ""
}

func (p *Pipeline) extractPipelinePlugins() []pluginConfigAPI {
	var pluginsConfigs []pluginConfigAPI
	v := reflect.ValueOf(p).Elem()
	for i := 0; i < v.NumField(); i++ {
		valueField := v.Field(i)
		if valueField.Kind() == reflect.Ptr && !valueField.IsNil() {
			fieldVal, ok := valueField.Interface().(pluginConfigAPI)
			if ok {
				pluginsConfigs = append(pluginsConfigs, fieldVal)
			} else {
				log.Warnf("jenkins extract pipeline plugins failed: %+v", valueField)
			}
		}

	}
	return pluginsConfigs
}

func (p *Pipeline) setDefaultValue(repoName, jenkinsNamespace string) {
	if p.JobName == "" {
		p.JobName = repoName
	}
	if p.ImageRepo != nil && p.ImageRepo.AuthNamespace == "" {
		p.ImageRepo.AuthNamespace = jenkinsNamespace
	}
}

func (p *Pipeline) buildCIConfig() *ci.CIConfig {
	ciConfig := &ci.CIConfig{
		Type: "jenkins",
	}
	ciConfig.ConfigLocation = p.JenkinsfilePath
	return ciConfig
}

package jenkins

import (
	"net/url"
	"strings"
)

type Pipeline struct {
	JobName         string    `mapstructure:"jobName" validate:"required"`
	JenkinsfilePath string    `mapstructure:"jenkinsfilePath" validate:"required"`
	ImageRepo       ImageRepo `mapstructure:"imageRepo"`
}

type ImageRepo struct {
	URL  string `mapstructure:"url" validate:"url"`
	User string `mapstructure:"user"`
}

func (p *Pipeline) getImageHost() string {
	harborAddress := p.ImageRepo.URL
	harborURL, err := url.ParseRequestURI(harborAddress)
	if err != nil {
		return harborAddress
	}
	return harborURL.Host
}

func (p *Pipeline) getJobName() string {
	if strings.Contains(p.JobName, "/") {
		return strings.Split(p.JobName, "/")[1]
	}
	return p.JobName
}

func (p *Pipeline) getJobPath() string {
	return p.JobName
}

func (p *Pipeline) getJobFolder() string {
	if strings.Contains(p.JobName, "/") {
		return strings.Split(p.JobName, "/")[0]
	}
	return ""
}

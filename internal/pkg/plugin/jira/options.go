package jira

import (
	"github.com/mitchellh/mapstructure"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer/ci/cifile"
	"github.com/devstream-io/devstream/pkg/util/scm/git"
)

type options struct {
	Scm  *git.RepoInfo `mapstructure:"scm" validate:"required"`
	Jira *jiraInfo     `mapstructure:"jira" validate:"required"`
	// used in package
	CIFileConfig *cifile.CIFileConfig `mapstructure:"ci"`
}

type jiraInfo struct {
	BaseUrl    string `validate:"required" mapstructure:"baseURL"`
	UserEmail  string `validate:"required" mapstructure:"userEmail"`
	ProjectKey string `validate:"required" mapstructure:"projectKey"`
	Token      string `validate:"required" mapstructure:"token"`
}

func newOptions(rawOptions configmanager.RawOptions) (*options, error) {
	var opts options
	if err := mapstructure.Decode(rawOptions, &opts); err != nil {
		return nil, err
	}
	return &opts, nil
}

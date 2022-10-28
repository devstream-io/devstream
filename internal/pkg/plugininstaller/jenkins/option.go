package jenkins

import (
	"github.com/mitchellh/mapstructure"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/ci/cifile"
	"github.com/devstream-io/devstream/pkg/util/scm"
	"github.com/devstream-io/devstream/pkg/util/scm/git"
)

const (
	ciType = "jenkins"
)

type JobOptions struct {
	Jenkins  jenkinsOption `mapstructure:"jenkins"`
	SCM      scm.SCMInfo   `mapstructure:"scm"`
	Pipeline pipeline      `mapstructure:"pipeline"`

	// used in package
	CIConfig    *cifile.CIConfig `mapstructure:"ci"`
	ProjectRepo *git.RepoInfo    `mapstructure:"projectRepo"`
}

func newJobOptions(options configmanager.RawOptions) (*JobOptions, error) {
	var opts JobOptions
	if err := mapstructure.Decode(options, &opts); err != nil {
		return nil, err
	}
	return &opts, nil
}

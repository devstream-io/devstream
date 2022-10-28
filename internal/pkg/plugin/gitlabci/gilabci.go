package gitlabci

import (
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/ci/cifile"
	"github.com/devstream-io/devstream/pkg/util/scm/git"
)

var DefaultCIOptions = &cifile.Options{
	CIConfig: &cifile.CIConfig{
		Type: "gitlab",
	},
	ProjectRepo: &git.RepoInfo{
		RepoType: "gitlab",
	},
}

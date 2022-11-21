package gitlabci

import (
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer/ci/cifile"
	"github.com/devstream-io/devstream/pkg/util/scm/git"
)

var DefaultCIOptions = &cifile.Options{
	CIFileConfig: &cifile.CIFileConfig{
		Type: "gitlab",
	},
	ProjectRepo: &git.RepoInfo{
		RepoType: "gitlab",
	},
}

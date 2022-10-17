package gitlabci

import (
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/ci"
	"github.com/devstream-io/devstream/pkg/util/scm"
)

var DefaultCIOptions = &ci.Options{
	CIConfig: &ci.CIConfig{
		Type: "gitlab",
	},
	ProjectRepo: &scm.Repo{
		RepoType: "gitlab",
	},
}

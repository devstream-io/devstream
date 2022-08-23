package gitlabci

import (
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/ci"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/common"
)

var DefaultCIOptions = &ci.Options{
	CIConfig: &ci.CIConfig{
		Type: "gitlab",
	},
	ProjectRepo: &common.Repo{
		RepoType: "gitlab",
	},
}

package devlake

import (
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/kubectl"
	"github.com/devstream-io/devstream/pkg/util/file"
	kubectlUtil "github.com/devstream-io/devstream/pkg/util/kubectl"
)

func Delete(options map[string]interface{}) (bool, error) {
	// Initialize Operator with Operations
	operator := &plugininstaller.Operator{
		ExecuteOperations: plugininstaller.ExecuteOperations{
			kubectl.ProcessByContent(
				kubectlUtil.Delete, file.NewTemplate().FromRemote(devLakeInstallYAMLDownloadURL),
			),
		},
	}

	// Execute all Operations in Operator
	_, err := operator.Execute(plugininstaller.RawOptions(options))
	if err != nil {
		return false, err
	}
	return true, nil
}

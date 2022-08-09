package devlake

import (
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/kubectl"
	"github.com/devstream-io/devstream/pkg/util/file"
)

func Delete(options map[string]interface{}) (bool, error) {
	// 1. config delete operations
	runner := &plugininstaller.Operator{
		ExecuteOperations: []plugininstaller.BaseOperation{
			kubectl.ProcessByContent(
				"delete", file.NewTemplate().FromRemote(devLakeInstallYAMLDownloadURL),
			),
		},
	}

	// 2. execute installer get status and error
	_, err := runner.Execute(plugininstaller.RawOptions(options))
	if err != nil {
		return false, err
	}
	return true, nil
}

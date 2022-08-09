package devlake

import (
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/kubectl"
	"github.com/devstream-io/devstream/pkg/util/file"
	"github.com/devstream-io/devstream/pkg/util/log"
)

func Create(options map[string]interface{}) (map[string]interface{}, error) {
	// 1. config install operations
	runner := &plugininstaller.Runner{
		ExecuteOperations: []plugininstaller.BaseOperation{
			kubectl.ProcessByContent(
				"create", file.NewTemplate().FromRemote(devLakeInstallYAMLDownloadURL),
			),
		},
		GetStatusOperation: getStaticState,
	}

	// 2. execute installer get status and error
	status, err := runner.Execute(plugininstaller.RawOptions(options))
	if err != nil {
		return nil, err
	}
	log.Debugf("Return map: %v", status)
	return status, nil
}

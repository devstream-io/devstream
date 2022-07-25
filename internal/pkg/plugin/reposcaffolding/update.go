package reposcaffolding

import (
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/reposcaffolding"
	"github.com/devstream-io/devstream/pkg/util/log"
)

func Update(options map[string]interface{}) (map[string]interface{}, error) {
	runner := &plugininstaller.Runner{
		PreExecuteOperations: []plugininstaller.MutableOperation{
			reposcaffolding.Validate,
			reposcaffolding.SetDefaultTemplateRepo,
		},
		ExecuteOperations: []plugininstaller.BaseOperation{
			reposcaffolding.DeleteRepo,
			reposcaffolding.InstallRepo,
		},
		GetStatusOperation: reposcaffolding.GetStaticState,
	}

	// 2. execute installer get status and error
	status, err := runner.Execute(plugininstaller.RawOptions(options))
	if err != nil {
		return nil, err
	}
	log.Debugf("Return map: %v", status)
	return status, nil

}

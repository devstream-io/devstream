package jenkinspipelinekubernetes

import (
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/pkg/util/log"
)

func Update(options map[string]interface{}) (map[string]interface{}, error) {
	// Initialize Operator with Operations
	operator := &plugininstaller.Operator{
		PreExecuteOperations: plugininstaller.PreExecuteOperations{
			ValidateAndDefaults,
		},
		ExecuteOperations: plugininstaller.ExecuteOperations{
			Todo,
		},
		TerminateOperations: nil,
		GetStateOperation:   GetState,
	}

	// Execute all Operations in Operator
	status, err := operator.Execute(plugininstaller.RawOptions(options))
	if err != nil {
		return nil, err
	}
	log.Debugf("Return map: %v", status)
	return status, nil
}

func Todo(options plugininstaller.RawOptions) error {
	// TODO(aFlyBird0): determine how to update the resource, such as:
	// if some config/resource are changed, we should restart the Jenkins
	// some, we should only call some update function
	// others, we just ignore them

	// now we just use the same way as create,
	// because the logic is the same: "if not exists, create; if exists, do nothing"
	// if it changes in the future, we should change the way to update
	return CreateJob(options)
}

package jenkinspipelinekubernetes

import (
	_ "embed"

	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/pkg/util/log"
)

const (
	jenkinsCredentialID       = "credential-jenkins-pipeline-kubernetes-by-devstream"
	jenkinsCredentialDesc     = "Jenkins Pipeline secret, created by devstream/jenkins-pipeline-kubernetes"
	jenkinsCredentialUsername = "foo-useless-username"
)

func Create(options map[string]interface{}) (map[string]interface{}, error) {
	// Initialize Operator with Operations
	operator := &plugininstaller.Operator{
		PreExecuteOperations: plugininstaller.PreExecuteOperations{
			ValidateAndDefaults,
		},
		ExecuteOperations: plugininstaller.ExecuteOperations{
			CreateJob,
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

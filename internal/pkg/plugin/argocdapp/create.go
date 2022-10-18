package argocdapp

import (
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/kubectl"
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
	kubectlUtil "github.com/devstream-io/devstream/pkg/util/kubectl"
	"github.com/devstream-io/devstream/pkg/util/log"
)

// Create creates an ArgoCD app YAML and applys it.
func Create(options map[string]interface{}) (statemanager.ResourceStatus, error) {
	// Initialize Operator with Operations
	operator := &plugininstaller.Operator{
		PreExecuteOperations: plugininstaller.PreExecuteOperations{
			validate,
		},
		ExecuteOperations: plugininstaller.ExecuteOperations{
			kubectl.ProcessByContent(kubectlUtil.Create, templateFileLoc),
		},
		GetStatusOperation: getStaticStatus,
	}

	// Execute all Operations in Operator
	status, err := operator.Execute(plugininstaller.RawOptions(options))
	if err != nil {
		return nil, err
	}
	log.Debugf("Return map: %v", status)
	return status, nil
}

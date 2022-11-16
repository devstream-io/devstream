package argocdapp

import (
	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/plugin/plugininstaller"
	"github.com/devstream-io/devstream/internal/pkg/plugin/plugininstaller/kubectl"
	kubectlUtil "github.com/devstream-io/devstream/pkg/util/kubectl"
)

func Delete(options configmanager.RawOptions) (bool, error) {
	// Initialize Operator with Operations
	operator := &plugininstaller.Operator{
		PreExecuteOperations: plugininstaller.PreExecuteOperations{
			validate,
		},
		ExecuteOperations: plugininstaller.ExecuteOperations{
			kubectl.ProcessByContent(kubectlUtil.Delete, templateFileLoc),
		},
	}

	// Execute all Operations in Operator
	_, err := operator.Execute(options)
	if err != nil {
		return false, err
	}
	return true, nil
}

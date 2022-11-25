package argocdapp

import (
	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	. "github.com/devstream-io/devstream/internal/pkg/plugin/common"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer/kubectl"
	kubectlUtil "github.com/devstream-io/devstream/pkg/util/kubectl"
)

func Delete(options configmanager.RawOptions) (bool, error) {
	var err error
	defer func() {
		HandleErrLogsWithPlugin(err, Name)
	}()

	// Initialize Operator with Operations
	operator := &installer.Operator{
		PreExecuteOperations: installer.PreExecuteOperations{
			validate,
		},
		ExecuteOperations: installer.ExecuteOperations{
			kubectl.ProcessByContent(kubectlUtil.Delete, templateFileLoc),
		},
	}

	// Execute all Operations in Operator
	_, err = operator.Execute(options)
	if err != nil {
		return false, err
	}
	return true, nil
}

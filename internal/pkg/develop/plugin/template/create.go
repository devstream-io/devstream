package template

var createGoNameTpl = "create.go"
var createGoDirTpl = "internal/pkg/plugin/[[ .Name | dirFormat ]]/"
var createGoContentTpl = `package [[ .Name | format ]]

import (
	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer"
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
	"github.com/devstream-io/devstream/pkg/util/log"
)

func Create(options configmanager.RawOptions) (statemanager.ResourceStatus, error) {
	// Initialize Operator with Operations
	operator := &installer.Operator{
		PreExecuteOperations: installer.PreExecuteOperations{
			// TODO(dtm): Add your PreExecuteOperations here.
		},
		ExecuteOperations: installer.ExecuteOperations{
			// TODO(dtm): Add your ExecuteOperations here.
		},
		TerminateOperations: installer.TerminateOperations{
			// TODO(dtm): Add your TerminateOperations here.
		},
		GetStatusOperation: func(options configmanager.RawOptions) (statemanager.ResourceStatus, error) {
			// TODO(dtm): Add your GetStatusOperation here.
			return nil, nil
		},
	}

	// Execute all Operations in Operator
	status, err := operator.Execute(options)
	if err != nil {
		return nil, err
	}
	log.Debugf("Return map: %v", status)
	return status, nil
}
`
var create_go_mustExistFlag = true

func init() {
	TplFiles = append(TplFiles, TplFile{
		NameTpl:       createGoNameTpl,
		DirTpl:        createGoDirTpl,
		ContentTpl:    createGoContentTpl,
		MustExistFlag: create_go_mustExistFlag,
	})
}
